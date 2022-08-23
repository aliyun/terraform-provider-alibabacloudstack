package apsarastack

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"strconv"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alibabacloudstack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApsaraStackDataWorksRemind() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackDataWorksRemindCreate,
		Read:   resourceApsaraStackDataWorksRemindRead,
		Update: resourceApsaraStackDataWorksRemindUpdate,
		Delete: resourceApsaraStackDataWorksRemindDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"remind_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
			},
			"alert_unit": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"OWNER", "OTHER"}, false),
			},
			"remind_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"remind_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"FINISHED", "UNFINISHED", "ERROR", "CYCLE_UNFINISHED", "TIMEOUT"}, false),
			},
			"remind_unit": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"NODE", "BASELINE", "PROJECT", "BIZPROCESS"}, false),
			},
			"dnd_end": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "00:00",
			},
			"node_ids": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"baseline_ids": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"biz_process_ids": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"max_alert_times": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      3,
				ValidateFunc: validation.IntBetween(1, 10),
			},
			"alert_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1800,
				ValidateFunc: validation.IntAtLeast(1200),
			},
			"detail": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"alert_methods": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"alert_targets": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"robot_urls": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"use_flag": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

var RemindUnit string
var RemindType string

func resourceApsaraStackDataWorksRemindCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	var response map[string]interface{}
	action := "CreateRemind"
	conn, err := client.NewDataworkspublicClient()
	if err != nil {
		return WrapError(err)
	}

	request := buildRemindArgs(d)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-05-18"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_data_works_remind", action, ApsaraStackSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["Data"]))

	return resourceApsaraStackDataWorksRemindRead(d, meta)
}
func resourceApsaraStackDataWorksRemindRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	dataworksPublicService := DataworksPublicService{client}
	object, err := dataworksPublicService.DescribeDataWorksRemind(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource apsarastack_data_works_remind dataworksPublicService.DescribeDataWorksRemind Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	if err != nil {
		return WrapError(err)
	}

	for key, value := range object {
		fmt.Println(key, value)
	}

	d.Set("remind_id", d.Id())
	d.Set("alert_unit", object["AlertUnit"].(string))
	d.Set("remind_name", object["RemindName"].(string))
	d.Set("remind_type", object["RemindType"].(string))
	d.Set("remind_unit", object["RemindUnit"].(string))
	d.Set("dnd_end", object["DndEnd"].(string))

	d.Set("baseline_ids", getObjectListToString(object, "Baselines", "BaselineId"))
	d.Set("node_ids", getObjectListToString(object, "Nodes", "NodeId"))
	d.Set("biz_process_ids", getObjectListToString(object, "BizProcesses", "BizId"))

	if len(object["Projects"].([]interface{})) > 0 {
		projectId := object["Projects"].([]interface{})[0].(map[string]interface{})["ProjectId"].(json.Number)
		d.Set("project_id", fmt.Sprint(projectId))
	}

	d.Set("max_alert_times", object["MaxAlertTimes"].(json.Number))
	d.Set("alert_interval", object["AlertInterval"].(json.Number))
	if RemindType == "TIMEOUT" {
		n, _ := strconv.Atoi(object["Detail"].(string))
		d.Set("detail", fmt.Sprintf("%d", n*60))
	} else {
		d.Set("detail", object["Detail"].(string))
	}

	d.Set("alert_methods", getObjectListToString(object, "AlertMethods", ""))
	d.Set("alert_targets", getObjectListToString(object, "AlertTargets", ""))
	d.Set("robot_urls", getObjectListToString(object, "Robots", "WebUrl"))
	d.Set("use_flag", object["Useflag"].(bool))

	return nil
}

func resourceApsaraStackDataWorksRemindUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	request := make(map[string]interface{})
	var response map[string]interface{}

	request["RemindId"] = d.Id()

	if d.HasChange("alert_interval") {
		request["AlertInterval"] = d.Get("alert_interval").(int)
	}

	if d.HasChange("remind_name") {
		request["RemindName"] = d.Get("remind_name").(string)
	}

	if d.HasChange("dnd_end") {
		request["DndEnd"] = d.Get("dnd_end").(string)
	}

	if d.HasChange("remind_unit") {
		RemindUnit = d.Get("remind_unit").(string)
		request["RemindUnit"] = RemindUnit

		buildRemindUnitArgs(d, request)
	}

	if d.HasChange("remind_type") {
		RemindType = d.Get("remind_type").(string)
		request["RemindType"] = RemindType
		request["Detail"] = d.Get("detail").(string)
	}

	if d.HasChange("max_alert_times") {
		request["MaxAlertTimes"] = d.Get("max_alert_times").(int)
	}

	if d.HasChange("alert_unit") {
		request["AlertUnit"] = d.Get("alert_unit").(string)
		request["AlertTargets"] = d.Get("alert_targets").(string)
	}

	if d.HasChange("alert_methods") {
		request["AlertMethods"] = d.Get("alert_methods").(string)
	}

	if d.HasChange("use_flag") {
		request["UseFlag"] = d.Get("use_flag").(bool)
	}

	if d.HasChange("robot_urls") {
		request["RobotUrls"] = d.Get("robot_urls").(string)
	}

	action := "UpdateRemind"
	request["RegionId"] = "default"
	conn, err := client.NewDataworkspublicClient()
	if err != nil {
		return WrapError(err)
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-05-18"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, ApsaraStackSdkGoERROR)
	}
	return resourceApsaraStackDataWorksRemindRead(d, meta)
}

func resourceApsaraStackDataWorksRemindDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)

	action := "DeleteRemind"
	var response map[string]interface{}
	conn, err := client.NewDataworkspublicClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"RemindId": d.Id(),
		"RegionId": "default",
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-05-18"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, ApsaraStackSdkGoERROR)
	}
	return nil
}

func buildRemindArgs(d *schema.ResourceData) map[string]interface{} {
	request := make(map[string]interface{})
	if v, ok := d.GetOk("alert_unit"); ok {
		request["AlertUnit"] = v.(string)
		request["AlertTargets"] = d.Get("alert_targets").(string)
	}

	if v, ok := d.GetOk("remind_name"); ok {
		request["RemindName"] = v.(string)
	}

	if v, ok := d.GetOk("remind_type"); ok {
		RemindType = v.(string)
		request["RemindType"] = RemindType
	}

	if v, ok := d.GetOk("remind_unit"); ok {
		RemindUnit = v.(string)
		request["RemindUnit"] = RemindUnit
	}

	if v, ok := d.GetOk("dnd_end"); ok {
		request["DndEnd"] = v.(string)
	}

	buildRemindUnitArgs(d, request)

	if v, ok := d.GetOk("max_alert_times"); ok {
		request["MaxAlertTimes"] = v.(int)
	}

	if v, ok := d.GetOk("alert_interval"); ok {
		request["AlertInterval"] = v.(int)
	}

	if v, ok := d.GetOk("detail"); ok {
		request["Detail"] = v.(string)
	}

	if v, ok := d.GetOk("alert_methods"); ok {
		request["AlertMethods"] = v.(string)
	}

	if v, ok := d.GetOk("robot_urls"); ok {
		request["RobotUrls"] = v.(string)
	}

	request["RegionId"] = "default"
	return request
}

func buildRemindUnitArgs(d *schema.ResourceData, request map[string]interface{}) {
	if RemindUnit == "NODE" {
		request["NodeIds"] = d.Get("node_ids").(string)
	} else if RemindUnit == "BASELINE" {
		request["BaselineIds"] = d.Get("baseline_ids").(string)
	} else if RemindUnit == "PROJECT" {
		request["ProjectId"] = d.Get("project_id").(string)
	} else if RemindUnit == "BIZPROCESS" {
		request["BizProcessIds"] = d.Get("biz_process_ids").(string)
	}
}

func getObjectListToString(object map[string]interface{}, listName string, mapKey string) string {
	var s string
	if len(mapKey) == 0 {
		for i, k := range object[listName].([]interface{}) {
			if i == 0 {
				s = fmt.Sprintf("%s", k)
			} else {
				s = fmt.Sprintf("%s,%s", s, k)
			}
		}
	} else {
		for i, k := range object[listName].([]interface{}) {
			if i == 0 {
				s = fmt.Sprintf("%s", k.(map[string]interface{})[mapKey])
			} else {
				s = fmt.Sprintf("%s,%s", s, k.(map[string]interface{})[mapKey])
			}
		}
	}

	return s
}
