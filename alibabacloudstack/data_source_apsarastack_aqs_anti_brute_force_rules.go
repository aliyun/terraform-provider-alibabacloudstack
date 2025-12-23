package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackAqsAntiBruteForceRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackAqsAntiBruteForceRulesRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A list of rule IDs to filter results by.",
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  "A regex string to filter results by rule name.",
			},
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"span": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"fail_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"forbidden_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"default_rule": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"instance_ids": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"enable_smart_rule": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"machine_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"create_timestamp": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackAqsAntiBruteForceRulesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
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
	pageNum := 1
	request := map[string]interface{}{
		"From":        "sas",
		"CurrentPage": pageNum,
		"PageSize":    "100",
	}
	rules := make([]interface{}, 0)
	for true {
		response, err := client.DoTeaRequest("POST", "aegis", "2016-11-11", "DescribeAntiBruteForceRules", "", nil, request, nil)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_aqs_anti_brute_force_rules", "DescribeAntiBruteForceRules", errmsgs.AlibabacloudStackSdkGoERROR)
		}
		r, ok := response["Rules"].([]interface{})
		if !ok || len(r) == 0 {
			break
		}
		rules = append(rules, r...)
		count, err := jsonpath.Get("$.PageInfo.TotalCount", response)
		if err != nil {
			return errmsgs.Error("DescribeAntiBruteForceRules Failed! %v", response)
		}
		totalCount, _ := count.(json.Number).Int64()
		if totalCount <= int64(pageNum*100) {
			break
		}
	}
	if len(rules) == 0 {
		d.SetId("")
		return nil
	}
	var s []map[string]interface{}
	var ids []string
	for _, v := range rules {
		rule := v.(map[string]interface{})
		idStr := fmt.Sprint(rule["Id"])
		if len(idsMap) > 0 {
			if _, ok := idsMap[idStr]; !ok {
				continue
			}
		}
		if nameRegex != nil {
			name := rule["Name"].(string)
			if !nameRegex.MatchString(name) {
				continue
			}
		}
		mapping := map[string]interface{}{}
		mapping["id"] = idStr
		mapping["name"] = rule["Name"]
		mapping["span"] = rule["Span"]
		mapping["fail_count"] = rule["FailCount"]
		mapping["forbidden_time"] = rule["ForbiddenTime"]
		mapping["default_rule"] = rule["DefaultRule"]
		mapping["enable_smart_rule"] = rule["EnableSmartRule"]
		mapping["machine_count"] = rule["MachineCount"]
		mapping["create_timestamp"] = rule["CreateTimestamp"]

		if uuidList, ok := rule["UuidList"].([]interface{}); ok {
			instance_ids, err := InstanceIdsHandler(uuidList, "uuid", meta)
			if err != nil {
				return errmsgs.WrapError(err)
			}
			mapping["instance_ids"] = instance_ids
		} else {
			mapping["instance_ids"] = []string{}
		}

		s = append(s, mapping)
		ids = append(ids, mapping["id"].(string))
	}
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("rules", s); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
