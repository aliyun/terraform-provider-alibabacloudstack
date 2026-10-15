package alibabacloudstack

import (
	"regexp"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackPolardbParameterGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackPolardbParameterGroupsRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: false,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"engine": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: false,
			},
			"engine_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: false,
			},
			"parameter_group_type": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: false,
			},
			// Computed top-level attributes (only the list of instances)
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parameter_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parameter_group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parameter_group_desc": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parameter_group_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"force_restart": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"param_counts": {
							Type:     schema.TypeInt,
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
						"parameters": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"param_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"param_value": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceAlibabacloudStackPolardbParameterGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	response, err := client.DoTeaRequest("POST", "polardb", "2017-08-01", "DescribeParameterGroups", "", nil, map[string]interface{}{}, nil)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	// Create filters based on schema
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		nameRegex = regexp.MustCompile(v.(string))
	}

	// The response is a flat array: response["ParameterGroups"] = []interface{}
	var rawGroups []interface{}
	if groups, ok := response["ParameterGroups"].([]interface{}); ok {
		rawGroups = groups
	}

	// Filter results
	var filteredGroups []map[string]interface{}
	for _, g := range rawGroups {
		group, ok := g.(map[string]interface{})
		if !ok {
			continue
		}

		// Get ID using firstNonEmptyString for field name compatibility
		groupId := firstNonEmptyString(group, "ParameterGroupId", "ParamGroupId")

		// Get name
		groupName := ""
		if v, ok := group["ParameterGroupName"].(string); ok {
			groupName = v
		}

		// Apply name regex filter
		if nameRegex != nil && !nameRegex.MatchString(groupName) {
			continue
		}

		// Apply ids filter
		if len(idsMap) > 0 {
			if _, ok := idsMap[groupId]; !ok {
				continue
			}
		}

		filteredGroups = append(filteredGroups, group)
	}

	// Prepare result data
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	groups := make([]map[string]interface{}, 0)

	for _, group := range filteredGroups {
		groupId := firstNonEmptyString(group, "ParameterGroupId", "ParamGroupId")
		engine := firstNonEmptyString(group, "DBType", "Engine")
		engineVersion := firstNonEmptyString(group, "DBVersion", "EngineVersion")
		created := firstNonEmptyString(group, "CreateTime", "Created")
		modified := firstNonEmptyString(group, "ModifyTime", "Modified")

		mapping := map[string]interface{}{
			"id":                   groupId,
			"parameter_group_id":   groupId,
			"parameter_group_name": group["ParameterGroupName"],
			"parameter_group_desc": group["ParameterGroupDesc"],
			"engine":               engine,
			"engine_version":       engineVersion,
			"parameter_group_type": mustToInt(group["ParameterGroupType"]),
			"force_restart":        mustToInt(group["ForceRestart"]),
			"param_counts":         mustToInt(group["ParameterCounts"]),
			"created":              created,
			"modified":             modified,
			"parameters":           []map[string]interface{}{},
		}

		// Add parameters if available in the response
		if paramDetail, ok := group["ParamDetail"]; ok {
			if detailMap, ok := paramDetail.(map[string]interface{}); ok {
				if paramDetails, ok := detailMap["ParameterDetail"]; ok {
					if paramArray, ok := paramDetails.([]interface{}); ok {
						parameters := make([]map[string]interface{}, 0)
						for _, paramItem := range paramArray {
							if paramMap, ok := paramItem.(map[string]interface{}); ok {
								paramMap2 := map[string]interface{}{
									"param_name":  paramMap["ParamName"],
									"param_value": paramMap["ParamValue"],
								}
								parameters = append(parameters, paramMap2)
							}
						}
						mapping["parameters"] = parameters
					}
				}
			}
		}

		ids = append(ids, groupId)
		names = append(names, group["ParameterGroupName"])
		groups = append(groups, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("groups", groups); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}

// firstNonEmptyString returns the first non-empty string value from the map for the given keys.
func firstNonEmptyString(m map[string]interface{}, keys ...string) string {
	for _, k := range keys {
		if v, ok := m[k]; ok {
			if s, ok := v.(string); ok && s != "" {
				return s
			}
		}
	}
	return ""
}

// mustToInt converts an interface{} to int, returning 0 on error.
func mustToInt(v interface{}) int {
	n, _ := toInt(v)
	return n
}
