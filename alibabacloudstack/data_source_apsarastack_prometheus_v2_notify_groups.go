package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackPrometheusV2NotifyGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackPrometheusV2NotifyGroupsRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A list of notify group IDs to filter results by.",
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  "A regex string to filter results by notify group name.",
			},
			"notify_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
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
						"webhook_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"webhook_header_params": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
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
						"im": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"contact_ids": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackPrometheusV2NotifyGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// Build request parameters
	reqBody := map[string]interface{}{
		"keyword":         "",
		"pageNumber":      1,
		"pageSize":        10,
		"direction":       "desc",
		"page":            1,
		"selectedRowKeys": []interface{}{},
	}

	// Call PageNotifyGroup API
	resp, err := client.DoTeaRequest("POST", "prometheus2", "2023-04-13", "PageNotifyGroup", "/log/api/v2/alert/group/page", nil, nil, reqBody)
	if err != nil {
		return fmt.Errorf("failed to call PageNotifyGroup API: %v", err)
	}

	// Parse response data
	dataRaw, ok := resp["data"]
	if !ok {
		return fmt.Errorf("missing data field in response")
	}

	dataBytes, err := json.Marshal(dataRaw)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}

	var pageResult struct {
		Total int `json:"total"`
		Data  []struct {
			ID          interface{} `json:"id"`
			Name        string      `json:"name"`
			Type        string      `json:"type"`
			Description string      `json:"description"`
			Webhook     struct {
				URL          string                 `json:"url"`
				HeaderParams map[string]interface{} `json:"headerParams"`
			} `json:"webhook"`
			Im         string        `json:"im"`
			ContactIds []interface{} `json:"contactIds"`
		} `json:"data"`
	}
	if err := json.Unmarshal(dataBytes, &pageResult); err != nil {
		return fmt.Errorf("failed to unmarshal data: %v", err)
	}

	// Prepare filters
	idsMap := getIdsStringFilter(d)

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		nameRegex = regexp.MustCompile(v.(string))
	}

	// Filter results
	var filteredResults []map[string]interface{}
	var ids []string
	for _, item := range pageResult.Data {
		idStr := fmt.Sprintf("%v", item.ID)
		// Apply id filter
		if len(idsMap) > 0 {
			if _, exists := idsMap[idStr]; !exists {
				continue
			}
		}

		// Apply name regex filter
		if nameRegex != nil && !nameRegex.MatchString(item.Name) {
			continue
		}

		// Construct webhook header params
		webhookHeaders := make([]map[string]interface{}, 0)
		for k, v := range item.Webhook.HeaderParams {
			webhookHeaders = append(webhookHeaders, map[string]interface{}{
				"key":   k,
				"value": v,
			})
		}

		// Convert contactIds to integers
		contactIds := make([]string, 0)
		for _, cid := range item.ContactIds {
			contactIds = append(contactIds, fmt.Sprint(cid))
		}

		result := map[string]interface{}{
			"id":                    idStr,
			"name":                  item.Name,
			"type":                  item.Type,
			"description":           item.Description,
			"webhook_url":           item.Webhook.URL,
			"webhook_header_params": webhookHeaders,
			"im":                    item.Im,
			"contact_ids":           contactIds,
		}
		ids = append(ids, idStr)
		filteredResults = append(filteredResults, result)
	}

	// Set computed values
	if err := d.Set("notify_groups", filteredResults); err != nil {
		return fmt.Errorf("error setting notify_groups: %v", err)
	}
	if err := d.Set("ids", ids); err != nil {
		return fmt.Errorf("error setting notify_groups: %v", err)
	}
	d.SetId(dataResourceIdHash(ids))

	return nil
}
