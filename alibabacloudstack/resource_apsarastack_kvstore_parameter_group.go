package alibabacloudstack

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackKvstoreParameterGroup() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"character_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"engine": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Redis",
				ForceNew: true,
			},
			"parameter_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"engine_version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"parameter_group_desc": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"parameters": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"param_name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"is_dynamic": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackKvstoreParameterGroupCreate, resourceAlibabacloudStackKvstoreParameterGroupRead, nil, resourceAlibabacloudStackKvstoreParameterGroupDelete)
	return resource
}

func resourceAlibabacloudStackKvstoreParameterGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := map[string]interface{}{
		"CharacterType":      d.Get("character_type").(string),
		"Engine":             d.Get("engine").(string),
		"ParameterGroupName": d.Get("parameter_group_name").(string),
		"EngineVersion":      d.Get("engine_version").(string),
		"ParameterGroupDesc": d.Get("parameter_group_desc").(string),
	}

	// Add parameters if provided
	if v, ok := d.GetOk("parameters"); ok {
		parameters := v.(*schema.Set).List()
		for i, param := range parameters {
			paramMap := param.(map[string]interface{})
			paramNameKey := fmt.Sprintf("Parameters.%d.ParamName", i+1)
			valueKey := fmt.Sprintf("Parameters.%d.Value", i+1)
			request[paramNameKey] = paramMap["param_name"].(string)
			request[valueKey] = paramMap["value"].(string)
		}
	}

	raw, err := client.DoTeaRequest("POST", "R-kvstore", "2015-01-01", "CreateParameterGroup", "", nil, request, nil)
	if err != nil {
		return err
	}

	// Extract parameter group id from response
	parameterGroupId, ok := raw["ParameterGroupId"]
	if !ok {
		return fmt.Errorf("failed to get ParameterGroupId from response")
	}

	// Set the resource ID temporarily
	d.SetId(parameterGroupId.(string))

	return nil
}

func resourceAlibabacloudStackKvstoreParameterGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	kvstoreService := KvstoreService{client}

	object, err := kvstoreService.DescribeParameterGroup(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_kvstore_parameter_group kvstoreService.DescribeParameterGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	d.Set("engine_version", object["EngineVersion"])
	d.Set("parameter_group_name", object["ParameterGroupName"])
	d.Set("parameter_group_desc", object["ParameterGroupDesc"])
	d.Set("create_time", object["CreateTime"])
	d.Set("character_type", object["CharacterType"])
	d.Set("type", object["Type"])
	d.Set("is_dynamic", object["IsDynamic"])

	parameters, err := kvstoreService.ListParameterGroupParms(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.Set("parameters", nil)
		} else {
			return errmsgs.WrapError(err)
		}
	}
	parametersList := make([]map[string]interface{}, 0)
	for _, v := range parameters {
		param := v.(map[string]interface{})
		parametersList = append(parametersList, map[string]interface{}{
			"param_name": param["ParamName"],
			"value":      param["Value"],
		})
	}
	d.Set("parameters", parametersList)

	return nil
}

func resourceAlibabacloudStackKvstoreParameterGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	reqQuery := map[string]interface{}{
		"ParameterGroupId": d.Id(),
	}

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.DoTeaRequest("POST", "R-kvstore", "2015-01-01", "DeleteParameterGroup", "", nil, reqQuery, nil)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"InvalidParameterGroupId.NotFound"}) {
				return resource.NonRetryableError(err)
			}
			err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "DeleteParameterGroup", errmsgs.AlibabacloudStackSdkGoERROR, err.Error())
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidParameterGroupId.NotFound"}) {
			return nil
		}
		return errmsgs.WrapError(err)
	}
	return nil
}
