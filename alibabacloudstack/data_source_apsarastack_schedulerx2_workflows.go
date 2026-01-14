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

func dataSourceAlibabacloudStackSchedulerx2Workflows() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackSchedulerx2WorkflowsRead,

		Schema: map[string]*schema.Schema{
			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "system_namespace",
			},
			"group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
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
			"workflows": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"workflow_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"time_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"time_expression": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"max_concurrency": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"creator": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updater": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"app_group_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackSchedulerx2WorkflowsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// Prepare filters
	idsMap := getIdsStringFilter(d)

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return fmt.Errorf("invalid name_regex: %v", err)
		}
		nameRegex = r
	}

	pageNum := 1
	requestQuery := map[string]interface{}{
		"Namespace": d.Get("namespace"),
		"PageSize":  100,
	}

	if groupId, ok := d.GetOk("group_id"); ok {
		requestQuery["GroupId"] = groupId
	}
	result := make([]interface{}, 0)
	for true {
		requestQuery["PageNum"] = pageNum
		resp, err := client.DoTeaRequest("POST", "schedulerx2", "2019-04-30", "ListWorkflows", "", nil, requestQuery, nil)
		if err != nil {
			return errmsgs.WrapError(err)
		}

		if fmt.Sprint(resp["Code"]) != "200" {
			return errmsgs.WrapError(fmt.Errorf("list workflows failed: %v", resp))
		}
		total, _ := jsonpath.Get("$.Data.Total", resp)
		count, _ := total.(json.Number).Int64()
		records, _ := jsonpath.Get("$.Data.Records", resp)
		result = append(result, records.([]interface{})...)
		if total == 0 || int(count) <= pageNum*100 {
			break
		}

	}

	filteredRecords := make([]map[string]interface{}, 0)
	for _, v := range result {
		record := v.(map[string]interface{})
		if len(idsMap) > 0 {
			idStr := fmt.Sprint(record["WorkflowId"])
			if _, ok := idsMap[idStr]; !ok {
				continue
			}
		}

		if nameRegex != nil {
			if name, ok := record["Name"].(string); ok {
				if !nameRegex.MatchString(name) {
					continue
				}
			} else {
				continue
			}
		}

		filteredRecords = append(filteredRecords, record)
	}

	// Construct output
	ids := make([]string, 0)
	workflows := make([]map[string]interface{}, 0)

	for _, record := range filteredRecords {
		mapping := map[string]interface{}{
			"id":              fmt.Sprintf("%v", record["WorkflowId"]),
			"workflow_id":     record["WorkflowId"],
			"group_id":        record["GroupId"],
			"name":            record["Name"],
			"description":     record["Description"],
			"time_type":       record["TimeType"],
			"time_expression": record["TimeExpression"],
			"max_concurrency": record["MaxConcurrency"],
			"creator":         record["Creator"],
			"updater":         record["Updater"],
			"app_group_id":    record["AppGroupId"],
		}

		ids = append(ids, fmt.Sprintf("%v", record["WorkflowId"]))
		workflows = append(workflows, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("workflows", workflows); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}
	return nil
}
