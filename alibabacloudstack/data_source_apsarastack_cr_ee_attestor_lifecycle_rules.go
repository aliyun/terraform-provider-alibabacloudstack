package alibabacloudstack

import (
	"regexp"
	"strconv"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackCrEEAttestorLifecycleRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackCrEEAttestorLifecycleRulesRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A list of rule IDs to filter results by.",
			},
			"namespace_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  "A regex string to filter results by namespace name.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The instance ID associated with the lifecycle rules.",
			},
			"enable_delete_tag": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Filter rules by whether tag deletion is enabled.",
			},
			"enable_delete_untagged_manifest": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Filter rules by whether untagged manifest deletion is enabled.",
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A list of matched namespace names.",
			},
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"namespace_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"repo_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"retention_tag_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"tag_regexp": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enable_delete_tag": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"auto": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"enable_delete_untagged_manifest": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"modified_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"schedule_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"next_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackCrEEAttestorLifecycleRulesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// Build request parameters
	request := map[string]interface{}{
		"InstanceId": d.Get("instance_id").(string),
		"PageSize":   30,
	}

	if v, ok := d.GetOk("enable_delete_tag"); ok {
		request["EnableDeleteTag"] = v.(bool)
	}

	if v, ok := d.GetOk("enable_delete_untagged_manifest"); ok {
		request["EnableDeleteUntaggedManifest"] = v.(bool)
	}

	// Initialize filters
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
	if v, ok := d.GetOk("namespace_regex"); ok {
		nameRegex = regexp.MustCompile(v.(string))
	}

	var allRules []interface{}
	pageNumber := 1

	// Paginate through all results
	for {
		request["PageNo"] = pageNumber

		response, err := client.DoTeaRequest("POST", "cr-ee", "2018-12-01", "ListArtifactLifecycleRule", "", nil, request, nil)
		if err != nil {
			return errmsgs.WrapError(err)
		}

		// Check if request was successful
		isSuccess, ok := response["IsSuccess"]
		if !ok || !isSuccess.(bool) {
			return errmsgs.WrapErrorf(err, "Failed to list artifact lifecycle rules")
		}

		// Extract rules from response
		rules, ok := response["Rules"]
		if !ok {
			break
		}

		rulesList := rules.([]interface{})
		if len(rulesList) == 0 {
			break
		}

		// Filter rules based on ids and namespace_regex
		for _, rule := range rulesList {
			ruleMap := rule.(map[string]interface{})

			// Check if rule id is in the ids filter
			ruleId := ""
			if id, ok := ruleMap["RuleId"]; ok {
				ruleId = id.(string)
			}

			if len(idsMap) > 0 {
				if _, ok := idsMap[ruleId]; !ok {
					continue
				}
			}

			// Check if namespace name matches the name_regex filter
			namespaceName := ""
			if name, ok := ruleMap["NamespaceName"]; ok {
				namespaceName = name.(string)
			}

			if nameRegex != nil && !nameRegex.MatchString(namespaceName) {
				continue
			}

			allRules = append(allRules, ruleMap)
		}

		// Check if we've reached the end of results
		totalCount := 0
		if count, ok := response["TotalCount"]; ok {
			totalCount = int(count.(float64))
		}

		if len(allRules) >= totalCount {
			break
		}

		pageNumber++
	}

	// Set computed values
	ids := make([]string, 0)
	names := make([]string, 0)
	rules := make([]map[string]interface{}, 0)

	for _, rule := range allRules {
		ruleMap := rule.(map[string]interface{})

		mapping := map[string]interface{}{}

		if v, ok := ruleMap["RuleId"]; ok {
			mapping["rule_id"] = v
			ids = append(ids, v.(string))
		}

		if v, ok := ruleMap["InstanceId"]; ok {
			mapping["instance_id"] = v
		}

		if v, ok := ruleMap["NamespaceName"]; ok {
			mapping["namespace_name"] = v
			names = append(names, v.(string))
		}

		if v, ok := ruleMap["RepoName"]; ok {
			mapping["repo_name"] = v
		}

		if v, ok := ruleMap["RetentionTagCount"]; ok {
			if count, err := strconv.Atoi(v.(string)); err == nil {
				mapping["retention_tag_count"] = count
			} else {
				mapping["retention_tag_count"] = int(v.(float64))
			}
		}

		if v, ok := ruleMap["TagRegexp"]; ok {
			mapping["tag_regexp"] = v
		}

		if v, ok := ruleMap["EnableDeleteTag"]; ok {
			mapping["enable_delete_tag"] = v
		}

		if v, ok := ruleMap["Auto"]; ok {
			mapping["auto"] = v
		}

		if v, ok := ruleMap["EnableDeleteUntaggedManifest"]; ok {
			mapping["enable_delete_untagged_manifest"] = v
		}

		if v, ok := ruleMap["ModifiedTime"]; ok {
			if time, err := strconv.Atoi(v.(string)); err == nil {
				mapping["modified_time"] = time
			} else {
				mapping["modified_time"] = int(v.(float64))
			}
		}

		if v, ok := ruleMap["CreateTime"]; ok {
			if time, err := strconv.Atoi(v.(string)); err == nil {
				mapping["create_time"] = time
			} else {
				mapping["create_time"] = int(v.(float64))
			}
		}

		if v, ok := ruleMap["ScheduleTime"]; ok {
			mapping["schedule_time"] = v
		}

		if v, ok := ruleMap["NextTime"]; ok {
			if time, err := strconv.Atoi(v.(string)); err == nil {
				mapping["next_time"] = time
			} else {
				mapping["next_time"] = int(v.(float64))
			}
		}

		rules = append(rules, mapping)
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return errmsgs.WrapError(err)
	}

	if err := d.Set("rules", rules); err != nil {
		return errmsgs.WrapError(err)
	}
	return nil
}
