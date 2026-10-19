package alibabacloudstack

import (
	"regexp"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackPolardbClusterParameterGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackPolardbClusterParameterGroupsRead,

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

func dataSourceAlibabacloudStackPolardbClusterParameterGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	response, err := client.DoTeaRequest("POST", "polardb", "2017-08-01", "DescribeParameterGroups", "", nil, map[string]interface{}{}, nil)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	// Build filters
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

	// Parse the flat list response
	groups, ok := response["ParameterGroups"].([]interface{})
	if !ok {
		d.SetId(dataResourceIdHash([]string{}))
		d.Set("ids", []string{})
		d.Set("names", []string{})
		d.Set("groups", []map[string]interface{}{})
		return nil
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	groupList := make([]map[string]interface{}, 0)

	for _, g := range groups {
		groupMap, ok := g.(map[string]interface{})
		if !ok {
			continue
		}

		groupId := firstNonEmptyStr(groupMap, "ParameterGroupId", "ParamGroupId")
		groupIdStr, _ := groupId.(string)
		groupName, _ := groupMap["ParameterGroupName"].(string)

		// Apply name regex filter
		if nameRegex != nil && !nameRegex.MatchString(groupName) {
			continue
		}

		// Apply ids filter
		if len(idsMap) > 0 {
			if _, ok := idsMap[groupIdStr]; !ok {
				continue
			}
		}

		mapping := map[string]interface{}{
			"id":                   groupIdStr,
			"parameter_group_id":   groupIdStr,
			"parameter_group_name": groupName,
			"parameter_group_desc": firstNonEmptyStr(groupMap, "ParameterGroupDesc", ""),
			"engine":               firstNonEmptyStr(groupMap, "DBType", "Engine"),
			"engine_version":       firstNonEmptyStr(groupMap, "DBVersion", "EngineVersion"),
			"parameter_group_type": mustToInt(groupMap["ParameterGroupType"]),
			"force_restart":        mustToInt(groupMap["ForceRestart"]),
			"param_counts":         mustToInt(firstNonEmptyStr(groupMap, "ParameterCounts", "ParamCounts")),
			"created":              firstNonEmptyStr(groupMap, "CreateTime", "Created"),
			"modified":             firstNonEmptyStr(groupMap, "ModifyTime", "Modified"),
			"parameters":           []map[string]interface{}{},
		}

		// Add parameters if available in the response
		if paramDetail, ok := groupMap["ParamDetail"]; ok {
			if detailMap, ok := paramDetail.(map[string]interface{}); ok {
				if paramDetails, ok := detailMap["ParameterDetail"].([]interface{}); ok {
					parameters := make([]map[string]interface{}, 0)
					for _, paramItem := range paramDetails {
						if paramMap, ok := paramItem.(map[string]interface{}); ok {
							parameters = append(parameters, map[string]interface{}{
								"param_name":  paramMap["ParamName"],
								"param_value": paramMap["ParamValue"],
							})
						}
					}
					mapping["parameters"] = parameters
				}
			}
		}

		ids = append(ids, groupIdStr)
		names = append(names, groupName)
		groupList = append(groupList, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("groups", groupList); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}

