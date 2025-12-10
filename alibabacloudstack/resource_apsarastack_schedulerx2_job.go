package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackSchedulerx2Job() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "system_namespace",
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"job_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"execute_mode": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"standalone", "broadcast", "parallel", "grid", "batch", "sharding"}, false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"priority": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5,
			},
			"parameters": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"max_attempt": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  3,
			},
			"attempt_interval": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  30,
			},
			"max_concurrency": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"time_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"time_expression": {
				Type:     schema.TypeString,
				Required: true,
			},
			"content": {
				Type:     schema.TypeString,
				Required: true,
			},
			"monitor_timeout_enable": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"monitor_timeout_kill_enable": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"monitor_fail_enable": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"monitor_miss_worker_enable": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"monitor_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackSchedulerx2JobCreate, resourceAlibabacloudStackSchedulerx2JobRead, resourceAlibabacloudStackSchedulerx2JobUpdate, resourceAlibabacloudStackSchedulerx2JobDelete)
	return resource
}

func resourceAlibabacloudStackSchedulerx2JobCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := make(map[string]interface{})
	request["Namespace"] = d.Get("namespace")
	request["GroupId"] = d.Get("group_id")
	request["Name"] = d.Get("name")
	request["JobType"] = d.Get("job_type")
	request["ExecuteMode"] = d.Get("execute_mode")
	request["Description"] = d.Get("description")
	request["Priority"] = d.Get("priority")
	request["Parameters"] = d.Get("parameters")
	request["MaxAttempt"] = d.Get("max_attempt")
	request["AttemptInterval"] = d.Get("attempt_interval")
	request["MaxConcurrency"] = d.Get("max_concurrency")
	request["TimeType"] = d.Get("time_type")
	request["TimeExpression"] = d.Get("time_expression")
	request["Content"] = d.Get("content")
	monitorConfig := make(map[string]interface{})
	monitorConfig["timeoutEnable"] = d.Get("monitor_timeout_enable")
	monitorConfig["timeoutKillEnable"] = d.Get("monitor_timeout_kill_enable")
	monitorConfig["failEnable"] = d.Get("monitor_fail_enable")
	monitorConfig["missWorkerEnable"] = d.Get("monitor_miss_worker_enable")
	monitorConfig["timeout"] = d.Get("monitor_timeout")
	monitorConfigStr, _ := json.Marshal(monitorConfig)
	request["MonitorConfig"] = string(monitorConfigStr)

	resp, err := client.DoTeaRequest("POST", "schedulerx2", "2019-04-30", "CreateJob", "", nil, request, nil)
	if err != nil {
		return err
	}

	jobIdStr, ok := resp["Data"].(string)
	if !ok {
		return fmt.Errorf("failed to get job id from response")
	}

	d.SetId(jobIdStr)

	return nil
}

func resourceAlibabacloudStackSchedulerx2JobRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	schedulerx2Service := Schedulerx2Service{client}

	object, err := schedulerx2Service.DescribeSchedulerx2Job(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("name", object["Name"])
	d.Set("group_id", object["GroupId"])
	d.Set("job_type", object["JobType"])
	execute_mode := strings.ToLower(d.Get("execute_mode").(string))
	d.Set("execute_mode", execute_mode)
	d.Set("description", object["Description"])
	d.Set("priority", object["Priority"])
	d.Set("parameters", object["Parameters"])
	d.Set("max_attempt", object["MaxAttempt"])
	d.Set("attempt_interval", object["AttemptInterval"])
	d.Set("max_concurrency", object["MaxConcurrency"])
	d.Set("time_type", object["TimeType"])
	d.Set("time_expression", object["TimeExpression"])
	d.Set("content", object["Content"])
	monitorConfig, ok := object["MonitorConfig"]
	if ok {
		monitor := make(map[string]interface{})
		err = json.Unmarshal([]byte(monitorConfig.(string)), &monitor)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		d.Set("monitor_timeout_enable", monitor["timeoutEnable"])
		d.Set("monitor_timeout_kill_enable", monitor["timeoutKillEnable"])
		d.Set("monitor_fail_enable", monitor["failEnable"])
		d.Set("monitor_miss_worker_enable", monitor["missWorkerEnable"])
		d.Set("monitor_timeout", monitor["timeout"])
	}

	return nil
}

func resourceAlibabacloudStackSchedulerx2JobUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	if d.IsNewResource() {
		return nil
	}

	request := make(map[string]interface{})
	request["Namespace"] = d.Get("namespace")
	request["GroupId"] = d.Get("group_id")
	request["Name"] = d.Get("name")
	request["ExecuteMode"] = d.Get("execute_mode")
	request["Description"] = d.Get("description")
	request["Priority"] = d.Get("priority")
	request["Parameters"] = d.Get("parameters")
	request["MaxAttempt"] = d.Get("max_attempt")
	request["AttemptInterval"] = d.Get("attempt_interval")
	request["MaxConcurrency"] = d.Get("max_concurrency")
	request["TimeType"] = d.Get("time_type")
	request["TimeExpression"] = d.Get("time_expression")
	request["Content"] = d.Get("content")
	monitorConfig := make(map[string]interface{})
	monitorConfig["timeoutEnable"] = d.Get("monitor_timeout_enable")
	monitorConfig["timeoutKillEnable"] = d.Get("monitor_timeout_kill_enable")
	monitorConfig["failEnable"] = d.Get("monitor_fail_enable")
	monitorConfig["missWorkerEnable"] = d.Get("monitor_miss_worker_enable")
	monitorConfig["timeout"] = d.Get("monitor_timeout")
	monitorConfigStr, _ := json.Marshal(monitorConfig)
	request["MonitorConfig"] = string(monitorConfigStr)
	_, err := client.DoTeaRequest("POST", "schedulerx2", "2019-04-30", "UpdateJob", "", nil, request, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_schedulerx2_job", "UpdateJob", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return nil
}

func resourceAlibabacloudStackSchedulerx2JobDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	reqQuery := map[string]interface{}{
		"Namespace":      d.Get("namespace").(string),
		"AcceptLanguage": "zh",
		"GroupId":        d.Get("group_id").(string),
		"JobId":          d.Id(),
	}

	_, err := client.DoTeaRequest("POST", "schedulerx2", "2019-04-30", "DeleteJob", "", nil, reqQuery, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "DeleteJob", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return nil
}
