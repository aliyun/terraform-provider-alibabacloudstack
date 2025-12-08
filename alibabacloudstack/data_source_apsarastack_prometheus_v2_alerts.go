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

func dataSourceAlibabacloudStackPrometheusV2Alerts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackPrometheusV2AlertsRead,

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
			"alerts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"notify_recovered": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_check_all": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"trigger_clusters": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"trigger_promql": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"trigger_period": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"trigger_severity": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"trigger_cron": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tag_set": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"notify_types": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"notify_group_ids": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
						"notify_interval": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"notification": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"recover_notification": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackPrometheusV2AlertsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// Build request body for PageAlert API
	requestBody := map[string]interface{}{
		"keyword":    "",
		"pageNumber": 1,
		"pageSize":   100,
	}

	// Call PageAlert API to list all alerts
	response, err := client.DoTeaRequest("POST", "prometheus2", "2023-04-13", "PageAlert", "/log/api/v2/alert/page", nil, nil, requestBody)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	// Parse response data
	dataRaw, err := jsonpath.Get("$.data.data", response)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, "PageAlert", "$.data.data", response)
	}

	// Filter by IDs if provided
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	// Compile name regex if provided
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		nameRegex = regexp.MustCompile(v.(string))
	}

	alerts := make([]map[string]interface{}, 0)
	ids := make([]string, 0)
	prometheusService := PrometheusService{client}
	for _, v := range dataRaw.([]interface{}) {
		item := v.(map[string]interface{})
		idStr := fmt.Sprint(item["id"].(json.Number))
		nameStr := item["name"].(string)

		// Apply ID filter
		if len(idsMap) > 0 {
			if _, exists := idsMap[idStr]; !exists {
				continue
			}
		}

		// Apply name regex filter
		if nameRegex != nil && !nameRegex.MatchString(nameStr) {
			continue
		}
		object, err := prometheusService.DescribePrometheusV2Alert(idStr)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		mapping := make(map[string]interface{})

		mapping["id"] = idStr
		mapping["name"] = nameStr

		if val, ok := object["notifyRecovered"].(bool); ok {
			mapping["notify_recovered"] = val
		}
		if val, ok := object["tagSet"].([]interface{}); ok {
			mapping["tag_set"] = val
		}
		if val, ok := object["recoverNotification"].(string); ok {
			mapping["recover_notification"] = val
		}

		// Trigger rule fields
		if triggerRule, ok := object["triggerRule"].(map[string]interface{}); ok {
			isCheckAll := false
			if clusterIds, ok := triggerRule["clusterIds"].([]interface{}); ok {
				mapping["trigger_clusters"] = clusterIds
				if len(clusterIds) == 0 {
					isCheckAll = true
				}
			}
			mapping["is_check_all"] = isCheckAll

			if val, ok := triggerRule["promql"].(string); ok {
				mapping["trigger_promql"] = val
			}
			if val, ok := triggerRule["period"].(string); ok {
				mapping["trigger_period"] = val
			}
			if val, ok := triggerRule["severity"].(string); ok {
				mapping["trigger_severity"] = val
			}
			if val, ok := triggerRule["cron"].(string); ok {
				mapping["trigger_cron"] = val
			}
		}

		// Notification message
		if notification, ok := object["notification"].(map[string]interface{}); ok {
			if msg, ok := notification["message"].(string); ok {
				mapping["notification"] = msg
			}
		}

		// Alert notifies (first entry only)
		if alertNotifies, ok := object["alertNotifies"].([]interface{}); ok && len(alertNotifies) > 0 {
			if notify, ok := alertNotifies[0].(map[string]interface{}); ok {
				if val, ok := notify["notifyTypes"].([]interface{}); ok {
					mapping["notify_types"] = val
				}
				if val, ok := notify["notifyGroupIds"].([]interface{}); ok {
					mapping["notify_group_ids"] = val
				}
				if val, ok := notify["notifyInterval"].(string); ok {
					mapping["notify_interval"] = val
				}
			}
		}
		ids = append(ids, idStr)
		alerts = append(alerts, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("alerts", alerts); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
