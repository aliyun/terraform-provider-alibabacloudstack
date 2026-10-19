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

func resourceAlibabacloudStackPolarDBClusterParameterGroup() *schema.Resource {
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
	setResourceFunc(resource, resourceAlibabacloudStackPolarDBClusterParameterGroupCreate, resourceAlibabacloudStackPolarDBClusterParameterGroupRead, resourceAlibabacloudStackPolarDBClusterParameterGroupUpdate, resourceAlibabacloudStackPolarDBClusterParameterGroupDelete)
	return resource
}

func resourceAlibabacloudStackPolarDBClusterParameterGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := map[string]interface{}{
		"DBType":             d.Get("engine").(string),
		"DBVersion":          d.Get("engine_version").(string),
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

	// Call the CreateParameterGroup API (polardb/2017-08-01, cluster edition)
	response, err := client.DoTeaRequest("POST", "polardb", "2017-08-01", "CreateParameterGroup", "", nil, request, nil)
	addDebug("CreateParameterGroup", response, request)
	if err != nil {
		return err
	}

	// The CreateParameterGroup API returns ParameterGroupId directly
	paramGroupId := ""
	if id, ok := response["ParameterGroupId"].(string); ok && id != "" {
		paramGroupId = id
	}

	if paramGroupId == "" {
		return fmt.Errorf("CreateParameterGroup did not return ParameterGroupId, response: %v", response)
	}

	d.SetId(paramGroupId)
	d.Set("parameter_group_id", paramGroupId)

	return nil
}

func resourceAlibabacloudStackPolarDBClusterParameterGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// DescribeParameterGroup (singular) returns empty data on ApsaraStack,
	// so we use DescribeParameterGroups (plural) and find by ID.
	response, err := client.DoTeaRequest("POST", "polardb", "2017-08-01", "DescribeParameterGroups", "", nil, map[string]interface{}{}, nil)
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_polardb_cluster_parameter_group DescribeParameterGroups Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	// Find the parameter group by ID in the flat list
	var object map[string]interface{}
	if groups, ok := response["ParameterGroups"].([]interface{}); ok {
		for _, group := range groups {
			if groupMap, ok := group.(map[string]interface{}); ok {
				if id, ok := groupMap["ParameterGroupId"].(string); ok && id == d.Id() {
					object = groupMap
					break
				}
			}
		}
	}

	if object == nil {
		log.Printf("[DEBUG] Resource alibabacloudstack_polardb_cluster_parameter_group not found in DescribeParameterGroups: %s", d.Id())
		d.SetId("")
		return nil
	}

	// Set computed fields from the response
	d.Set("parameter_group_id", firstNonEmptyStr(object, "ParameterGroupId", "ParamGroupId"))
	d.Set("engine", firstNonEmptyStr(object, "DBType", "Engine"))
	d.Set("engine_version", firstNonEmptyStr(object, "DBVersion", "EngineVersion"))
	d.Set("parameter_group_type", object["ParameterGroupType"])
	d.Set("param_counts", firstNonEmptyStr(object, "ParameterCounts", "ParamCounts"))
	d.Set("force_restart", object["ForceRestart"])
	d.Set("created", firstNonEmptyStr(object, "CreateTime", "Created"))
	d.Set("modified", firstNonEmptyStr(object, "ModifyTime", "Modified"))
	d.Set("parameter_group_name", object["ParameterGroupName"])
	d.Set("parameter_group_desc", object["ParameterGroupDesc"])

	return nil
}

func resourceAlibabacloudStackPolarDBClusterParameterGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAlibabacloudStackPolarDBClusterParameterGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	reqQuery := map[string]interface{}{
		"ParameterGroupId": d.Id(),
	}

	err := resource.Retry(10*time.Second, func() *resource.RetryError {
		_, err := client.DoTeaRequest("POST", "polardb", "2017-08-01", "DeleteParameterGroup", "", nil, reqQuery, nil)
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

// firstNonEmptyStr returns the first non-empty string value from the map for the given keys.
func firstNonEmptyStr(m map[string]interface{}, keys ...string) interface{} {
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
