package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackPolarDBParameterGroup() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"engine_version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"engine": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"parameters": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"parameter_group_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"parameter_group_desc": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"parameter_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parameter_group_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"param_counts": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"force_restart": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modified": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackPolarDBParameterGroupCreate, resourceAlibabacloudStackPolarDBParameterGroupRead, resourceAlibabacloudStackPolarDBParameterGroupUpdate, resourceAlibabacloudStackPolarDBParameterGroupDelete)
	return resource
}

func resourceAlibabacloudStackPolarDBParameterGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := map[string]interface{}{
		"Engine":             d.Get("engine").(string),
		"EngineVersion":      d.Get("engine_version").(string),
		"ParameterGroupName": d.Get("parameter_group_name").(string),
	}
	if v, ok := d.GetOk("parameter_group_desc"); ok {
		request["ParameterGroupDesc"] = v.(string)
	}
	if v, ok := d.GetOk("parameters"); ok {
		paramsMap := make(map[string]string)
		for key, value := range v.(map[string]interface{}) {
			paramsMap[key] = value.(string)
		}

		paramsJSON, err := json.Marshal(paramsMap)
		if err != nil {
			return fmt.Errorf("failed to marshal parameters to JSON: %v", err)
		}
		request["Parameters"] = string(paramsJSON)
	}

	// Call the CreateParameterGroup API
	response, err := client.DoTeaRequest("POST", "polardb", "2024-01-30", "CreateParameterGroup", "", nil, request, nil)
	addDebug("CreateParameterGroup", response, request)
	if err != nil {
		return err
	}

	// Since the create API doesn't return the parameter group ID, we need to retrieve it by querying
	// Get all parameter groups and find the one that matches our criteria
	describeReq := map[string]interface{}{}
	response, err = client.DoTeaRequest("POST", "polardb", "2024-01-30", "DescribeParameterGroups", "", nil, describeReq, nil)
	if err != nil {
		return err
	}

	// Find the parameter group with matching name
	paramGroupId := ""
	if parameterGroups, ok := response["ParameterGroups"].(map[string]interface{}); ok {
		if groups, ok := parameterGroups["ParameterGroup"].([]interface{}); ok {
			for _, group := range groups {
				if groupMap, ok := group.(map[string]interface{}); ok {
					if groupName, ok := groupMap["ParameterGroupName"].(string); ok && groupName == d.Get("parameter_group_name").(string) {
						if id, ok := groupMap["ParamGroupId"].(string); ok {
							paramGroupId = id
							break
						}
					}
				}
			}
		}
	}

	if paramGroupId == "" {
		return fmt.Errorf("failed to find ParameterGroupId after creation")
	}

	// Set the temporary ID before any potential status checks
	d.SetId(paramGroupId)

	// Set the computed fields from the response
	d.Set("parameter_group_id", paramGroupId)

	return nil
}

func resourceAlibabacloudStackPolarDBParameterGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	polardbService := PolardbService{client}

	object, err := polardbService.DescribeParameterGroup(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_polardb_parameter_group polardbService.DescribeParameterGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	// Set computed fields from the response
	d.Set("parameter_group_id", object["ParamGroupId"])
	d.Set("engine", object["Engine"])
	d.Set("engine_version", object["EngineVersion"])
	d.Set("parameter_group_type", object["ParameterGroupType"])
	d.Set("param_counts", object["ParamCounts"])
	d.Set("force_restart", object["ForceRestart"])
	d.Set("created", object["Created"])
	d.Set("modified", object["Modified"])
	d.Set("parameter_group_name", object["ParameterGroupName"])
	d.Set("parameter_group_desc", object["ParameterGroupDesc"])

	// Handle parameters field - convert ParamDetail back to JSON string format
	if paramDetail, ok := object["ParamDetail"]; ok {
		if detailMap, ok := paramDetail.(map[string]interface{}); ok {
			if paramDetails, ok := detailMap["ParameterDetail"]; ok {
				if paramArray, ok := paramDetails.([]interface{}); ok {
					paramsMap := make(map[string]string)
					for _, paramItem := range paramArray {
						if paramMap, ok := paramItem.(map[string]interface{}); ok {
							paramName := paramMap["ParamName"].(string)
							paramValue := paramMap["ParamValue"].(string)
							paramsMap[paramName] = paramValue
						}
					}

					// paramsJSON, err := json.Marshal(paramsMap)
					if err == nil {
						d.Set("parameters", paramsMap)
					}
				}
			}
		}
	}

	return nil
}

func resourceAlibabacloudStackPolarDBParameterGroupUpdate(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceAlibabacloudStackPolarDBParameterGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	reqQuery := map[string]interface{}{
		"ParameterGroupId": d.Id(),
	}

	err := resource.Retry(10*time.Second, func() *resource.RetryError {
		_, err := client.DoTeaRequest("POST", "polardb", "2024-01-30", "DeleteParameterGroup", "", nil, reqQuery, nil)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"ParameterGroup.NotFound", "InvalidParameterGroupId.NotFound"}) {
				return resource.NonRetryableError(err)
			}
			err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "DeleteParameterGroup", errmsgs.AlibabacloudStackSdkGoERROR, err.Error())
			return resource.RetryableError(err)
		}
		return nil
	})

	return errmsgs.WrapError(err)
}
