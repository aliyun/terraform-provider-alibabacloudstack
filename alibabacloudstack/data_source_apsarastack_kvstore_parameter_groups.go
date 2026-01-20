package alibabacloudstack

import (
	"regexp"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackKvstoreParameterGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackKvstoreParameterGroupsRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"engine_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"character_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"logic", "normal"}, false),
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
						"engine_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"character_type": {
							Type:     schema.TypeString,
							Computed: true,
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
						"parameters": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"param_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackKvstoreParameterGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := make(map[string]interface{})
	if v, ok := d.GetOk("engine_version"); ok {
		request["EngineVersion"] = v.(string)
	}
	if v, ok := d.GetOk("character_type"); ok {
		request["CharacterType"] = v.(string)
	}

	raw, err := client.DoTeaRequest("POST", "R-kvstore", "2015-01-01", "DescribeParameterGroups", "", nil, request, nil)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	// Parse response
	parameterGroupsData, ok := raw["ParameterGroups"]
	if !ok {
		return nil
	}

	groupsList, ok := parameterGroupsData.(map[string]interface{})["ParameterGroups"]
	if !ok {
		return nil
	}

	groups, ok := groupsList.([]interface{})
	if !ok {
		return nil
	}

	// Create idsMap for filtering if ids are specified
	idsMap := getIdsStringFilter(d)
	// Create nameRegex for filtering if name_regex is specified
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		nameRegex = regexp.MustCompile(v.(string))
	}

	var filteredGroups []interface{}
	for _, group := range groups {
		groupMap := group.(map[string]interface{})

		// Check if the group matches the id filter
		groupId := groupMap["ParameterGroupId"].(string)
		if len(idsMap) > 0 {
			if _, exists := idsMap[groupId]; !exists {
				continue
			}
		}

		// Check if the group matches the name regex filter
		groupName := groupMap["ParameterGroupName"].(string)
		if nameRegex != nil {
			if !nameRegex.MatchString(groupName) {
				continue
			}
		}

		filteredGroups = append(filteredGroups, groupMap)
	}

	// Prepare result data
	ids := make([]string, 0)
	names := make([]string, 0)
	groupsResult := make([]map[string]interface{}, 0)
	kvstoreService := KvstoreService{client}
	for _, group := range filteredGroups {
		groupMap := group.(map[string]interface{})

		mapping := map[string]interface{}{}
		mapping["id"] = groupMap["ParameterGroupId"]
		mapping["parameter_group_id"] = groupMap["ParameterGroupId"]
		mapping["parameter_group_name"] = groupMap["ParameterGroupName"]
		mapping["parameter_group_desc"] = groupMap["ParameterGroupDesc"]
		mapping["engine_version"] = groupMap["EngineVersion"]
		mapping["character_type"] = groupMap["CharacterType"]
		mapping["create_time"] = groupMap["CreateTime"]
		mapping["type"] = groupMap["Type"]
		mapping["is_dynamic"] = groupMap["IsDynamic"]
		parameters, err := kvstoreService.ListParameterGroupParms(groupMap["ParameterGroupId"].(string))
		if err != nil && !errmsgs.NotFoundError(err) {
			return errmsgs.WrapError(err)
		}
		parametersList := make([]map[string]interface{}, 0)
		for _, v := range parameters {
			param := v.(map[string]interface{})
			parametersList = append(parametersList, map[string]interface{}{
				"param_name": param["ParamName"],
				"value":      param["Value"],
			})
		}
		mapping["parameters"] = parametersList
		ids = append(ids, groupMap["ParameterGroupId"].(string))
		names = append(names, groupMap["ParameterGroupName"].(string))
		groupsResult = append(groupsResult, mapping)
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return errmsgs.WrapError(err)
	}

	if err := d.Set("groups", groupsResult); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
