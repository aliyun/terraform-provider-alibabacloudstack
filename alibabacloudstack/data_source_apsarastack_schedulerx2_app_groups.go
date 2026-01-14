package alibabacloudstack

import (
	"fmt"
	"regexp"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackSchedulerx2AppGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackSchedulerx2AppGroupsRead,

		Schema: map[string]*schema.Schema{
			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "system_namespace",
			},
			"department": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"app_name": {
				Type:     schema.TypeString,
				Optional: true,
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
						"app_group_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"app_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"app_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"max_jobs": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_concurrency": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"xattrs": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"monitor_config": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"metrics_threshold_json": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parent_group_id": {
							Type:     schema.TypeString,
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
						"unique_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"global_max_jobs": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"accept_lang": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enable_log": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"log_config_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"contact_group_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cur_jobs": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"leader": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"alive_workers": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"auto_scale": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"app_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"alarm_json": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackSchedulerx2AppGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// Build request parameters
	reqQuery := make(map[string]interface{})
	if v, ok := d.GetOk("namespace"); ok {
		reqQuery["Namespace"] = v
	}
	if v, ok := d.GetOk("app_name"); ok {
		reqQuery["AppName"] = v
	}

	// Pagination settings
	pageSize := 10
	pageNum := 1
	var allGroups []interface{}

	for {
		reqQuery["PageSize"] = pageSize
		reqQuery["PageNum"] = pageNum

		resp, err := client.DoTeaRequest("POST", "schedulerx2", "2019-04-30", "ListGroups", "", nil, reqQuery, nil)
		if err != nil {
			return fmt.Errorf("failed to list schedulerx2 app groups: %v", err)
		}

		data, ok := resp["Data"].(map[string]interface{})
		if !ok || data == nil {
			break
		}

		totalRaw, _ := data["Total"].(float64)
		total := int(totalRaw)
		recordsRaw, _ := data["Records"].([]interface{})
		if recordsRaw == nil {
			recordsRaw = []interface{}{}
		}

		allGroups = append(allGroups, recordsRaw...)

		// Check pagination
		if len(allGroups) >= total {
			break
		}
		pageNum++
	}

	// Filter by ids
	idsMap := getIdsStringFilter(d)

	// Filter by name_regex
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		nameRegex = regexp.MustCompile(v.(string))
	}

	filteredGroups := make([]interface{}, 0)
	for _, groupRaw := range allGroups {
		group, ok := groupRaw.(map[string]interface{})
		if !ok {
			continue
		}
		idStr := fmt.Sprintf("%v", group["Id"])

		// Apply id filter
		if len(idsMap) > 0 {
			if _, exists := idsMap[idStr]; !exists {
				continue
			}
		}

		// Apply name regex filter
		if nameRegex != nil {
			appName := group["AppName"].(string)
			if !nameRegex.MatchString(appName) {
				continue
			}
		}

		filteredGroups = append(filteredGroups, group)
	}

	// Prepare results
	ids := make([]string, 0)
	groups := make([]map[string]interface{}, 0)

	for _, groupRaw := range filteredGroups {
		group := groupRaw.(map[string]interface{})
		idStr := fmt.Sprintf("%v", group["Id"])

		mapping := map[string]interface{}{
			"id":                     idStr,
			"app_group_id":           group["Id"],
			"app_name":               group["AppName"],
			"app_key":                group["AppKey"],
			"description":            group["Description"],
			"group_id":               group["GroupId"],
			"max_jobs":               group["MaxJobs"],
			"max_concurrency":        group["MaxConcurrency"],
			"xattrs":                 group["XAttr"],
			"version":                group["Version"],
			"monitor_config":         group["MonitorConfig"],
			"metrics_threshold_json": group["MetricsThresholdJson"],
			"parent_group_id":        group["ParentGroupId"],
			"creator":                group["Creator"],
			"updater":                group["Updater"],
			"unique_id":              group["UniqueId"],
			"global_max_jobs":        group["GlobalMaxJobs"],
			"accept_lang":            group["AcceptLang"],
			"enable_log":             group["EnableLog"],
			"log_config_id":          group["LogConfigId"],
			"contact_group_id":       group["ContactGroupId"],
			"cur_jobs":               group["CurJobs"],
			"leader":                 group["Leader"],
			"alive_workers":          group["AliveWorkers"],
			"auto_scale":             group["AutoScale"],
			"app_type":               group["AppType"],
			"alarm_json":             group["AlarmJson"],
		}

		ids = append(ids, idStr)
		groups = append(groups, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("groups", groups); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
