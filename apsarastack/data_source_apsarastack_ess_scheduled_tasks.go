package apsarastack

import (
	"regexp"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceApsaraStackEssScheduledTasks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApsaraStackEssScheduledTasksRead,
		Schema: map[string]*schema.Schema{
			"scheduled_task_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scheduled_action": {
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
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tasks": {
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
						"scheduled_action": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"launch_expiration_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"launch_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"max_value": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"min_value": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"recurrence_end_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"recurrence_value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"recurrence_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"task_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceApsaraStackEssScheduledTasksRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	request := ess.CreateDescribeScheduledTasksRequest()
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ess", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.RegionId = client.RegionId
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)

	if id, ok := d.GetOk("scheduled_task_id"); ok && id.(string) != "" {
		request.ScheduledTaskId1 = id.(string)
	}
	if a, ok := d.GetOk("scheduled_action"); ok && a.(string) != "" {
		request.ScheduledAction1 = a.(string)
	}

	var allScheduledTasks []ess.ScheduledTask

	for {
		raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.DescribeScheduledTasks(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "apsarastack_ess_scheduled_tasks", request.GetActionName(), ApsaraStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response := raw.(*ess.DescribeScheduledTasksResponse)
		if len(response.ScheduledTasks.ScheduledTask) < 1 {
			break
		}
		allScheduledTasks = append(allScheduledTasks, response.ScheduledTasks.ScheduledTask...)
		if len(response.ScheduledTasks.ScheduledTask) < PageSizeLarge {
			break
		}
		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	var filteredScheduledtasks = make([]ess.ScheduledTask, 0)

	nameRegex, okNameRegex := d.GetOk("name_regex")
	idsMap := make(map[string]string)
	ids, okIds := d.GetOk("ids")
	if okIds {
		for _, i := range ids.([]interface{}) {
			idsMap[i.(string)] = i.(string)
		}
	}

	if okNameRegex || okIds {
		for _, task := range allScheduledTasks {
			if okNameRegex && nameRegex != "" {
				var r = regexp.MustCompile(nameRegex.(string))
				if r != nil && !r.MatchString(task.ScheduledTaskName) {
					continue
				}
			}
			if okIds && len(idsMap) > 0 {
				if _, ok := idsMap[task.ScheduledTaskId]; !ok {
					continue
				}
			}
			filteredScheduledtasks = append(filteredScheduledtasks, task)
		}
	} else {
		filteredScheduledtasks = allScheduledTasks
	}

	return scheduledTasksDescriptionAttribute(d, filteredScheduledtasks, meta)
}

func scheduledTasksDescriptionAttribute(d *schema.ResourceData, tasks []ess.ScheduledTask, meta interface{}) error {
	var ids []string
	var names []string
	var s = make([]map[string]interface{}, 0)
	for _, t := range tasks {
		mapping := map[string]interface{}{
			"id":                     t.ScheduledTaskId,
			"name":                   t.ScheduledTaskName,
			"scheduled_action":       t.ScheduledAction,
			"description":            t.Description,
			"launch_expiration_time": t.LaunchExpirationTime,
			"launch_time":            t.LaunchTime,
			"max_value":              t.MaxValue,
			"min_value":              t.MinValue,
			"recurrence_end_time":    t.RecurrenceEndTime,
			"recurrence_value":       t.RecurrenceValue,
			"recurrence_type":        t.RecurrenceType,
			"task_enabled":           t.TaskEnabled,
		}
		ids = append(ids, t.ScheduledTaskId)
		names = append(names, t.ScheduledTaskName)
		s = append(s, mapping)
	}
	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("tasks", s); err != nil {
		return WrapError(err)
	}

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
