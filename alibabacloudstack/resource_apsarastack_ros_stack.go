package alibabacloudstack

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackRosStack() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackRosStackCreate,
		Read:   resourceAlibabacloudStackRosStackRead,
		Update: resourceAlibabacloudStackRosStackUpdate,
		Delete: resourceAlibabacloudStackRosStackDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(11 * time.Minute),
			Delete: schema.DefaultTimeout(6 * time.Minute),
			Update: schema.DefaultTimeout(11 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_option": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"deletion_protection": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Disabled", "Enabled"}, false),
				Default:      "Disabled",
			},
			"disable_rollback": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"notification_urls": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ForceNew: true,
			},
			"parameters": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"parameter_key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"parameter_value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"ram_role_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"replacement_option": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"retain_all_resources": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"retain_resources": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"stack_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"stack_policy_body": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsJSON,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"stack_policy_during_update_body": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"stack_policy_during_update_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"stack_policy_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"template_body": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsJSON,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"template_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"template_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"timeout_in_minutes": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  60,
			},
			"use_previous_parameters": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceAlibabacloudStackRosStackCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	rosService := RosService{client}
	action := "CreateStack"
	request := make(map[string]interface{})

	if v, ok := d.GetOk("create_option"); ok {
		request["CreateOption"] = v
	}

	if v, ok := d.GetOk("deletion_protection"); ok {
		request["DeletionProtection"] = v
	}

	if v, ok := d.GetOkExists("disable_rollback"); ok {
		request["DisableRollback"] = v
	}

	if v, ok := d.GetOk("notification_urls"); ok {
		request["NotificationURLs"] = v.(*schema.Set).List()
	}

	if v, ok := d.GetOk("parameters"); ok {
		parameters := make([]map[string]interface{}, len(v.(*schema.Set).List()))
		for i, j := range v.(*schema.Set).List() {
			parameters[i] = make(map[string]interface{})
			parameters[i]["ParameterKey"] = j.(map[string]interface{})["parameter_key"]
			parameters[i]["ParameterValue"] = j.(map[string]interface{})["parameter_value"]
		}
		request["Parameters"] = parameters
	}

	if v, ok := d.GetOk("ram_role_name"); ok {
		request["RamRoleName"] = v
	}

	request["StackName"] = d.Get("stack_name")
	if v, ok := d.GetOk("stack_policy_body"); ok {
		request["StackPolicyBody"] = v
	}

	if v, ok := d.GetOk("stack_policy_url"); ok {
		request["StackPolicyURL"] = v
	}

	if v, ok := d.GetOk("template_body"); ok {
		request["TemplateBody"] = v
	}

	if v, ok := d.GetOk("template_url"); ok {
		request["TemplateURL"] = v
	}

	if v, ok := d.GetOk("template_version"); ok {
		request["TemplateVersion"] = v
	}

	if v, ok := d.GetOk("timeout_in_minutes"); ok {
		request["TimeoutInMinutes"] = v
	}

	response, err := client.DoTeaRequest("POST", "ROS", "2019-09-10", action, "", nil, request)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprint(response["StackId"]))
	stateConf := BuildStateConf([]string{}, []string{"CREATE_COMPLETE"}, d.Timeout(schema.TimeoutCreate), 100*time.Second, rosService.RosStackStateRefreshFunc(d.Id(), []string{"CREATE_FAILED", "CREATE_ROLLBACK_COMPLETE", "CREATE_ROLLBACK_FAILED"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}

	return resourceAlibabacloudStackRosStackRead(d, meta)
}

func resourceAlibabacloudStackRosStackRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	rosService := RosService{client}
	object, err := rosService.DescribeRosStack(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_ros_stack rosService.DescribeRosStack Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	d.Set("deletion_protection", object["DeletionProtection"])
	d.Set("disable_rollback", object["DisableRollback"])

	parameters := make([]map[string]interface{}, 0)
	if parametersList, ok := object["Parameters"].([]interface{}); ok {
		for _, v := range parametersList {
			if m1, ok := v.(map[string]interface{}); ok {
				temp1 := map[string]interface{}{
					"parameter_key":   m1["ParameterKey"],
					"parameter_value": m1["ParameterValue"],
				}
				if !strings.HasPrefix(v.(map[string]interface{})["ParameterKey"].(string), "ALIYUN::") {
					parameters = append(parameters, temp1)
				}
			}
		}
	}
	if err := d.Set("parameters", parameters); err != nil {
		return errmsgs.WrapError(err)
	}
	d.Set("ram_role_name", object["RamRoleName"])
	d.Set("stack_name", object["StackName"])
	d.Set("status", object["Status"])
	d.Set("timeout_in_minutes", formatInt(object["TimeoutInMinutes"]))

	return nil
}

func resourceAlibabacloudStackRosStackUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	rosService := RosService{client}
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"StackId": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("disable_rollback") {
		update = true
		request["DisableRollback"] = d.Get("disable_rollback")
	}
	if !d.IsNewResource() && d.HasChange("parameters") {
		update = true
		parameters := make([]map[string]interface{}, len(d.Get("parameters").(*schema.Set).List()))
		for i, v := range d.Get("parameters").(*schema.Set).List() {
			parameters[i] = make(map[string]interface{})
			parameters[i]["ParameterKey"] = v.(map[string]interface{})["parameter_key"]
			parameters[i]["ParameterValue"] = v.(map[string]interface{})["parameter_value"]
		}
		request["Parameters"] = parameters
	}
	if !d.IsNewResource() && d.HasChange("ram_role_name") {
		update = true
		request["RamRoleName"] = d.Get("ram_role_name")
	}
	if !d.IsNewResource() && d.HasChange("stack_policy_body") {
		update = true
		request["StackPolicyBody"] = d.Get("stack_policy_body")
	}
	if !d.IsNewResource() && d.HasChange("timeout_in_minutes") {
		update = true
		request["TimeoutInMinutes"] = d.Get("timeout_in_minutes")
	}
	if update {
		if _, ok := d.GetOk("replacement_option"); ok {
			request["ReplacementOption"] = d.Get("replacement_option")
		}
		if _, ok := d.GetOk("stack_policy_during_update_body"); ok {
			request["StackPolicyDuringUpdateBody"] = d.Get("stack_policy_during_update_body")
		}
		if _, ok := d.GetOk("stack_policy_during_update_url"); ok {
			request["StackPolicyDuringUpdateURL"] = d.Get("stack_policy_during_update_url")
		}
		if _, ok := d.GetOk("stack_policy_url"); ok {
			request["StackPolicyURL"] = d.Get("stack_policy_url")
		}
		if _, ok := d.GetOk("template_body"); ok {
			request["TemplateBody"] = d.Get("template_body")
		}
		if _, ok := d.GetOk("template_url"); ok {
			request["TemplateURL"] = d.Get("template_url")
		}
		if _, ok := d.GetOk("template_version"); ok {
			request["TemplateVersion"] = d.Get("template_version")
		}
		if _, ok := d.GetOkExists("use_previous_parameters"); ok {
			request["UsePreviousParameters"] = d.Get("use_previous_parameters")
		}

		_, err := client.DoTeaRequest("POST", "ROS", "2019-09-10", "UpdateStack", "", nil, request)
		if err != nil {
			return err
		}
		stateConf := BuildStateConf([]string{}, []string{"UPDATE_COMPLETE"}, d.Timeout(schema.TimeoutUpdate), 100*time.Second, rosService.RosStackStateRefreshFunc(d.Id(), []string{"UPDATE_FAILED", "ROLLBACK_FAILED"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
	}
	d.Partial(false)
	return resourceAlibabacloudStackRosStackRead(d, meta)
}

func resourceAlibabacloudStackRosStackDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	rosService := RosService{client}
	action := "DeleteStack"
	request := map[string]interface{}{
		"StackId": d.Id(),
	}

	if v, ok := d.GetOk("ram_role_name"); ok {
		request["RamRoleName"] = v
	}
	if v, ok := d.GetOkExists("retain_all_resources"); ok {
		request["RetainAllResources"] = v
	}
	if v, ok := d.GetOk("retain_resources"); ok {
		request["RetainResources"] = v.(*schema.Set).List()
	}

	_, err := client.DoTeaRequest("POST", "ROS", "2019-09-10", action, "", nil, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"StackNotFound"}) {
			return nil
		}
		return err
	}
	stateConf := BuildStateConf([]string{}, []string{"DELETE_COMPLETE"}, d.Timeout(schema.TimeoutDelete), 100*time.Second, rosService.RosStackStateRefreshFunc(d.Id(), []string{"DELETE_FAILED"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}
	return nil
}
