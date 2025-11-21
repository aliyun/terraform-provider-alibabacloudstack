package alibabacloudstack

import (
	"encoding/json"
	"fmt"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackSchedulerx2AppGroup() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"app_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"max_jobs": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"max_concurrency": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"monitor_config": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"send_channel": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"ding", "mail", "mail,ding", ""}, false),
						},
						"alarm_type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "CustomContacts",
						},
					},
				},
			},
			"contacts": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"user_email": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"dingding_ak": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"metrics_threshold": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"load5": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
						"heap5_usage": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  90,
						},
						"disk_usage": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  95,
						},
					},
				},
			},
			"accept_language": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"app_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"app_group_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackSchedulerx2AppGroupCreate, resourceAlibabacloudStackSchedulerx2AppGroupRead, resourceAlibabacloudStackSchedulerx2AppGroupUpdate, resourceAlibabacloudStackSchedulerx2AppGroupDelete)
	return resource
}

func resourceAlibabacloudStackSchedulerx2AppGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	reqBody := make(map[string]interface{})
	reqBody["Namespace"] = d.Get("namespace")
	reqBody["GroupId"] = d.Get("group_id")
	reqBody["AppName"] = d.Get("app_name")
	if v, ok := d.GetOk("description"); ok {
		reqBody["Description"] = v
	}
	if v, ok := d.GetOk("max_jobs"); ok {
		reqBody["MaxJobs"] = v
	}
	if v, ok := d.GetOk("max_concurrency"); ok {
		reqBody["MaxConcurrency"] = v
	}
	if v, ok := d.GetOk("monitor_config"); ok {
		monitorConfig := v.(*schema.Set).List()[0].(map[string]interface{})
		monitor := map[string]interface{}{
			"sendChannel": monitorConfig["send_channel"],
			"alarmType":   monitorConfig["alarm_type"],
		}
		monitorConfigJson, err := json.Marshal(monitor)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		reqBody["MonitorConfigJson"] = string(monitorConfigJson)
	}
	if v, ok := d.GetOk("contacts"); ok {
		contacts := v.(*schema.Set).List()
		contactlist := make([]map[string]interface{}, len(contacts))
		for _, v := range contacts {
			contact := v.(map[string]interface{})
			contactlist = append(contactlist, map[string]interface{}{
				"userName":  contact["username"],
				"userEmail": contact["user_email"],
				"dingding":  contact["dingding_ak"],
			})
		}
		contactsJson, err := json.Marshal(contactlist)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		reqBody["ContactsJson"] = string(contactsJson)
	}
	if v, ok := d.GetOk("metrics_threshold"); ok {
		metricsThreshold := v.(*schema.Set).List()[0].(map[string]interface{})
		metrics_threshold := map[string]interface{}{
			"load5":      metricsThreshold["load5"],
			"heap5Usage": metricsThreshold["heap5_usage"],
			"diskUsage":  metricsThreshold["disk_usage"],
		}
		metrics_threshold_json, err := json.Marshal(metrics_threshold)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		reqBody["MetricsThresholdJson"] = string(metrics_threshold_json)
	}
	if v, ok := d.GetOk("accept_language"); ok {
		reqBody["AcceptLanguage"] = v
	}
	reqBody["Action"] = "CreateAppGroup"
	reqBody["AccessKeyId"] = client.AccessKey
	reqBody["SignatureMethod"] = "HMAC-SHA1"
	resp, err := client.DoTeaRequest("POST", "schedulerx2", "2019-04-30", "CreateAppGroup", "/openapi/v2/group/create", nil, reqBody, nil)
	if err != nil {
		return err
	}

	data, ok := resp["Data"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("failed to get Data from response")
	}

	appGroupId, ok := data["AppGroupId"].(float64)
	if !ok {
		return fmt.Errorf("failed to get AppGroupId from response data")
	}

	d.SetId(fmt.Sprintf("%d", int64(appGroupId)))

	return nil
}

func resourceAlibabacloudStackSchedulerx2AppGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	schedulerx2Service := Schedulerx2Service{client}
	object, err := schedulerx2Service.DescribeSchedulerx2AppGroup(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}
	d.Set("group_id", object["GroupId"])
	d.Set("app_name", object["AppName"])
	d.Set("description", object["Description"])
	d.Set("max_jobs", object["MaxJobs"])
	d.Set("max_concurrency", object["MaxConcurrency"])
	if v, ok := object["MonitorConfig"]; ok {
		monitorConfig := make(map[string]interface{})
		_ = json.Unmarshal([]byte(v.(string)), &monitorConfig)
		monitor_config := map[string]interface{}{
			"send_channel": monitorConfig["sendChannel"],
			"alarm_type":   monitorConfig["alarmType"],
		}
		if v, ok := monitorConfig["sendChannel"]; ok {
			monitor_config["send_channel"] = v
		}
		d.Set("monitor_config", []interface{}{monitor_config})
	}
	if v, ok := object["MetricsThresholdJson"]; ok {
		metricsThresholdJson := make(map[string]interface{})
		_ = json.Unmarshal([]byte(v.(string)), &metricsThresholdJson)
		metricsThreshold := map[string]interface{}{
			"load5":       metricsThresholdJson["load5"],
			"heap5_usage": metricsThresholdJson["heap5Usage"],
			"disk_usage":  metricsThresholdJson["diskUsage"],
		}
		d.Set("metrics_threshold", []interface{}{metricsThreshold})
	}
	d.Set("accept_language", object["AcceptLang"])
	d.Set("app_key", object["AppKey"])

	return nil
}

func resourceAlibabacloudStackSchedulerx2AppGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	if d.IsNewResource() {
		return nil
	}
	if d.HasChanges("namespace", "app_name", "description", "max_jobs",
		"max_concurrency", "monitor_config", "contacts", "metrics_threshold", "accept_language") {
		request := make(map[string]interface{})
		request["Namespace"] = d.Get("namespace")
		request["GroupId"] = d.Get("group_id")
		request["AppName"] = d.Get("app_name")
		if v, ok := d.GetOk("description"); ok {
			request["Description"] = v
		}
		if v, ok := d.GetOk("max_jobs"); ok {
			request["MaxJobs"] = v
		}
		if v, ok := d.GetOk("max_concurrency"); ok {
			request["MaxConcurrency"] = v
		}
		if v, ok := d.GetOk("monitor_config"); ok {
			monitorConfig := v.(*schema.Set).List()[0].(map[string]interface{})
			monitor := map[string]interface{}{
				"sendChannel": monitorConfig["send_channel"],
				"alarmType":   monitorConfig["alarm_type"],
			}
			monitorConfigJson, err := json.Marshal(monitor)
			if err != nil {
				return errmsgs.WrapError(err)
			}
			request["MonitorConfigJson"] = string(monitorConfigJson)
		}
		if v, ok := d.GetOk("contacts"); ok {
			contacts := v.(*schema.Set).List()
			contactlist := make([]map[string]interface{}, len(contacts))
			for _, v := range contacts {
				contact := v.(map[string]interface{})
				contactlist = append(contactlist, map[string]interface{}{
					"userName":  contact["username"],
					"userEmail": contact["user_email"],
					"dingding":  contact["dingding_ak"],
				})
			}
			contactsJson, err := json.Marshal(contactlist)
			if err != nil {
				return errmsgs.WrapError(err)
			}
			request["ContactsJson"] = string(contactsJson)
		}
		if v, ok := d.GetOk("metrics_threshold"); ok {
			metricsThreshold := v.(*schema.Set).List()[0].(map[string]interface{})
			metrics_threshold := map[string]interface{}{
				"load5":      metricsThreshold["load5"],
				"heap5Usage": metricsThreshold["heap5_usage"],
				"diskUsage":  metricsThreshold["disk_usage"],
			}
			metrics_threshold_json, err := json.Marshal(metrics_threshold)
			if err != nil {
				return errmsgs.WrapError(err)
			}
			request["MetricsThresholdJson"] = string(metrics_threshold_json)
		}
		if v, ok := d.GetOk("accept_language"); ok {
			request["AcceptLanguage"] = v
		}
		request["Action"] = "UpdateAppGroup"
		request["AccessKeyId"] = client.AccessKey
		_, err := client.DoTeaRequest("POST", "schedulerx2", "2019-04-30", "UpdateAppGroup", "/openapi/v1/group/update", nil, request, nil)
		if err != nil {
			return fmt.Errorf("failed to update schedulerx2 app group: %v", err)
		}
	}

	return nil
}

func resourceAlibabacloudStackSchedulerx2AppGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	reqQuery := map[string]interface{}{
		"Namespace": d.Get("namespace"),
		"GroupId":   d.Get("group_id"),
	}
	reqQuery["Action"] = "DeleteAppGroup"
	_, err := client.DoTeaRequest("POST", "schedulerx2", "2019-04-30", "DeleteAppGroup", "/openapi/v2/group/delete", nil, reqQuery, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "DeleteAppGroup", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return nil
}
