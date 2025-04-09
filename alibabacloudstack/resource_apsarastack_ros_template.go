package alibabacloudstack

import (
	"fmt"
	"log"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackRosTemplate() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
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
			"template_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"template_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackRosTemplateCreate, resourceAlibabacloudStackRosTemplateRead, resourceAlibabacloudStackRosTemplateUpdate, resourceAlibabacloudStackRosTemplateDelete)
	return resource
}

func resourceAlibabacloudStackRosTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var response map[string]interface{}
	action := "CreateTemplate"
	request := make(map[string]interface{})

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if v, ok := d.GetOk("template_body"); ok {
		request["TemplateBody"] = v
	}

	request["TemplateName"] = d.Get("template_name")
	if v, ok := d.GetOk("template_url"); ok {
		request["TemplateURL"] = v
	}

	response, err := client.DoTeaRequest("POST", "ROS", "2019-09-10", action, "", nil, nil, request)
	if err != nil {
		errmsg := ""
		if response != nil {
			errmsg = errmsgs.GetAsapiErrorMessage(response)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ros_template", action, errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	d.SetId(fmt.Sprint(response["TemplateId"]))

	return nil
}

func resourceAlibabacloudStackRosTemplateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	rosService := RosService{client}
	object, err := rosService.DescribeRosTemplate(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_ros_template rosService.DescribeRosTemplate Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	d.Set("template_body", object["TemplateBody"])
	return nil
}

func resourceAlibabacloudStackRosTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"TemplateId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}
	if !d.IsNewResource() && d.HasChange("template_body") {
		update = true
		request["TemplateBody"] = d.Get("template_body")
	}
	if !d.IsNewResource() && d.HasChange("template_name") {
		update = true
		request["TemplateName"] = d.Get("template_name")
	}
	if update {
		if _, ok := d.GetOk("template_url"); ok {
			request["TemplateURL"] = d.Get("template_url")
		}
		action := "UpdateTemplate"
		response, err := client.DoTeaRequest("POST", "ROS", "2019-09-10", action, "", nil, nil, request)
		if err != nil {
			errmsg := ""
			if response != nil {
				errmsg = errmsgs.GetAsapiErrorMessage(response)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), action, errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(action, response, request)
	}
	d.Partial(false)
	return nil
}

func resourceAlibabacloudStackRosTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	action := "DeleteTemplate"
	var response map[string]interface{}
	request := map[string]interface{}{
		"TemplateId": d.Id(),
	}
	response, err := client.DoTeaRequest("POST", "ROS", "2019-09-10", action, "", nil, nil, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"ChangeSetNotFound", "StackNotFound", "TemplateNotFound"}) {
			return nil
		}
		errmsg := ""
		if response != nil {
			errmsg = errmsgs.GetAsapiErrorMessage(response)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), action, errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	return nil
}