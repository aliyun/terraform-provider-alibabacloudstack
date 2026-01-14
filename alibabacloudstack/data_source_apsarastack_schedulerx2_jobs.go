package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackSchedulerx2Jobs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackSchedulerx2JobsRead,

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
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"jobs": {
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
						"group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"job_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"execute_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"priority": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"parameters": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"max_attempt": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"attempt_interval": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_concurrency": {
							Type:     schema.TypeInt,
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
						"content": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"monitor_timeout_enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"monitor_timeout_kill_enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"monitor_fail_enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"monitor_miss_worker_enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"monitor_timeout": {
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
						"version": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"app_group_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"data_offset": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"contact": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"calendar": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"clean_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gmt_create": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gmt_modified": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"template": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"timezone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"workflow_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"content_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"xattrs": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackSchedulerx2JobsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// Build request parameters
	pageNum := 1
	pageSize := 100
	reqQuery := map[string]interface{}{
		"Action":      "ListJobs",
		"AccessKeyId": client.AccessKey,
		"Namespace":   d.Get("namespace"),
		"PageNum":     pageNum,
		"PageSize":    100,
	}
	result := make([]interface{}, 0)
	for true {
		response, err := client.DoTeaRequest("POST", "schedulerx2", "2019-04-30", "ListJobs", "", nil, reqQuery, nil)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		total, err := jsonpath.Get("$.Data.Total", response)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		records, err := jsonpath.Get("$.Data.Records", response)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		count, _ := total.(json.Number).Int64()
		result = append(result, records.([]interface{})...)
		if pageNum*pageSize >= int(count) {
			break
		}
		pageNum++
	}

	// Prepare filters
	idsMap := getIdsStringFilter(d)

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return fmt.Errorf("invalid name_regex pattern: %v", err)
		}
		nameRegex = r
	}

	var filteredJobs []map[string]interface{}
	var jobIds []string

	for _, record := range result {
		job, ok := record.(map[string]interface{})
		if !ok {
			continue
		}
		jobId := fmt.Sprint(job["JobId"])

		jobName := job["Name"].(string)

		// Apply name regex filter
		if nameRegex != nil && !nameRegex.MatchString(jobName) {
			continue
		}

		// Apply ids filter
		if len(idsMap) > 0 {
			if _, exists := idsMap[jobId]; !exists {
				continue
			}
		}
		execute_mode := strings.ToLower(job["ExecuteMode"].(string))
		// Construct job info map
		jobInfo := map[string]interface{}{
			"id":               jobId,
			"name":             jobName,
			"group_id":         job["GroupId"],
			"job_type":         job["JobType"],
			"execute_mode":     execute_mode,
			"description":      job["Description"],
			"priority":         job["Priority"],
			"parameters":       job["Parameters"],
			"max_attempt":      job["MaxAttempt"],
			"attempt_interval": job["AttemptInterval"],
			"max_concurrency":  job["MaxConcurrency"],
			"time_type":        job["TimeType"],
			"time_expression":  job["TimeExpression"],
			"content":          job["Content"],
			"creator":          job["Creator"],
			"updater":          job["Updater"],
			"version":          job["Version"],
			"app_group_id":     job["AppGroupId"],
			"data_offset":      job["DataOffset"],
			"contact":          job["Contact"],
			"calendar":         job["Calendar"],
			"resource":         job["Resource"],
			"clean_mode":       job["CleanMode"],
			"gmt_create":       job["GmtCreate"],
			"gmt_modified":     job["GmtModified"],
			"template":         job["Template"],
			"timezone":         job["Timezone"],
			"workflow_id":      job["WorkFlowId"],
			"content_type":     job["ContentType"],
			"xattrs":           job["XAttrs"],
		}

		// Parse MonitorConfig
		if monitorConfigStr, ok := job["MonitorConfig"].(string); ok && monitorConfigStr != "" {
			var monitorConfig map[string]interface{}
			if err := json.Unmarshal([]byte(monitorConfigStr), &monitorConfig); err == nil {
				jobInfo["monitor_timeout_enable"] = monitorConfig["timeoutEnable"]
				jobInfo["monitor_timeout_kill_enable"] = monitorConfig["timeoutKillEnable"]
				jobInfo["monitor_fail_enable"] = monitorConfig["failEnable"]
				jobInfo["monitor_miss_worker_enable"] = monitorConfig["missWorkerEnable"]
				if timeout, ok := monitorConfig["timeout"].(float64); ok {
					jobInfo["monitor_timeout"] = int(timeout)
				}
			}
		}

		filteredJobs = append(filteredJobs, jobInfo)
		jobIds = append(jobIds, jobId)
	}

	d.SetId(dataResourceIdHash(jobIds))
	if err := d.Set("ids", jobIds); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("jobs", filteredJobs); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
