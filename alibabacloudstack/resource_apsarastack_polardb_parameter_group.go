package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"MySQL"}, true),
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
		"Engine":             strings.ToLower(d.Get("engine").(string)),
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

	// Try to get ParameterGroupId from the response directly
	paramGroupId := ""
	if id, ok := response["ParameterGroupId"].(string); ok && id != "" {
		paramGroupId = id
	} else if id, ok := response["ParamGroupId"].(string); ok && id != "" {
		paramGroupId = id
	}

	// If the API didn't return the ID, look it up by name via DescribeParameterGroups
	if paramGroupId == "" {
		targetName := d.Get("parameter_group_name").(string)
		listResp, listErr := client.DoTeaRequest("POST", "polardb", "2024-01-30", "DescribeParameterGroups", "", nil, map[string]interface{}{}, nil)
		if listErr != nil {
			return fmt.Errorf("CreateParameterGroup succeeded but DescribeParameterGroups failed: %v", listErr)
		}
		var createGroupList []interface{}
		if pg, ok := listResp["ParameterGroups"]; ok {
			if nested, ok := pg.(map[string]interface{}); ok {
				if arr, ok := nested["ParameterGroup"].([]interface{}); ok {
					createGroupList = arr
				}
			} else if arr, ok := pg.([]interface{}); ok {
				createGroupList = arr
			}
		}
		for _, group := range createGroupList {
			if groupMap, ok := group.(map[string]interface{}); ok {
				if name, _ := groupMap["ParameterGroupName"].(string); name == targetName {
					if id := firstNonEmptyString(groupMap, "ParamGroupId", "ParameterGroupId"); id != "" {
						paramGroupId = id
						break
					}
				}
			}
		}
	}

	if paramGroupId == "" {
		return fmt.Errorf("CreateParameterGroup succeeded but could not determine ParameterGroupId, response: %v", response)
	}

	d.SetId(paramGroupId)
	d.Set("parameter_group_id", paramGroupId)

	return nil
}

func resourceAlibabacloudStackPolarDBParameterGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// DescribeParameterGroup (singular) returns empty data on ApsaraStack,
	// so we use DescribeParameterGroups (plural) and find by ID.
	reqQuery := map[string]interface{}{}
	response, err := client.DoTeaRequest("POST", "polardb", "2024-01-30", "DescribeParameterGroups", "", nil, reqQuery, nil)
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_polardb_parameter_group DescribeParameterGroups Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	// Find the parameter group by ID in the list
	// Response may be nested (ParameterGroups.ParameterGroup[]) or flat (ParameterGroups[])
	var object map[string]interface{}
	var groupList []interface{}
	if pg, ok := response["ParameterGroups"]; ok {
		if nested, ok := pg.(map[string]interface{}); ok {
			if arr, ok := nested["ParameterGroup"].([]interface{}); ok {
				groupList = arr
			}
		} else if arr, ok := pg.([]interface{}); ok {
			groupList = arr
		}
	}
	for _, group := range groupList {
		if groupMap, ok := group.(map[string]interface{}); ok {
			id := firstNonEmptyString(groupMap, "ParamGroupId", "ParameterGroupId")
			if id == d.Id() {
				object = groupMap
				break
			}
		}
	}

	if object == nil {
		log.Printf("[DEBUG] Resource alibabacloudstack_polardb_parameter_group not found in DescribeParameterGroups: %s", d.Id())
		d.SetId("")
		return nil
	}

	// Set computed fields from the response (handle both old and new field names)
	d.Set("parameter_group_id", firstNonEmpty(object, "ParameterGroupId", "ParamGroupId"))
	d.Set("engine", firstNonEmptyString(object, "Engine", "DBType"))
	d.Set("engine_version", firstNonEmptyString(object, "EngineVersion", "DBVersion"))
	d.Set("parameter_group_type", object["ParameterGroupType"])
	d.Set("param_counts", firstNonEmpty(object, "ParameterCounts", "ParamCounts"))
	d.Set("force_restart", object["ForceRestart"])
	d.Set("created", firstNonEmpty(object, "CreateTime", "Created"))
	d.Set("modified", firstNonEmpty(object, "ModifyTime", "Modified"))
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
			if errmsgs.IsExpectedErrors(err, "ParameterGroup.NotFound", "InvalidParameterGroupId.NotFound") {
				return resource.NonRetryableError(err)
			}
			err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "DeleteParameterGroup", errmsgs.AlibabacloudStackSdkGoERROR, err.Error())
			return resource.RetryableError(err)
		}
		return nil
	})

	return errmsgs.WrapError(err)
}

// firstNonEmpty returns the first non-empty string value from the map for the given keys.
func firstNonEmpty(m map[string]interface{}, keys ...string) interface{} {
	for _, k := range keys {
		if v, ok := m[k]; ok {
			if s, ok := v.(string); ok && s != "" {
				return v
			}
			if v != nil {
				return v
			}
		}
	}
	return nil
}
