package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackDataWorksRemind() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackDataWorksRemindCreate,
		Read:   resourceAlibabacloudStackDataWorksRemindRead,
		Update: resourceAlibabacloudStackDataWorksRemindUpdate,
		Delete: resourceAlibabacloudStackDataWorksRemindDelete,
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

func resourceAlibabacloudStackDataWorksRemindCreate(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	action := "CreateRemind"
	request := buildRemindArgs(d)

	response, err := client.DoTeaRequest("POST", "dataworks-public", "2020-05-18", action, "", nil, nil, request)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprint(response["Data"]))

	return resourceAlibabacloudStackDataWorksRemindRead(d, meta)
}

func resourceAlibabacloudStackDataWorksRemindRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	dataworksPublicService := DataworksPublicService{client}
	object, err := dataworksPublicService.DescribeDataWorksRemind(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_data_works_remind dataworksPublicService.DescribeDataWorksRemind Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		errmsg := ""
		if object != nil {
			errmsg = errmsgs.GetAsapiErrorMessage(object)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_data_works_remind", "DescribeDataWorksRemind", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
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

func resourceAlibabacloudStackDataWorksRemindUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := make(map[string]interface{})

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
	_, err = client.DoTeaRequest("POST", "dataworks-public", "2020-05-18", action, "", nil, nil, request)
	if err != nil {
		return err
	}
	return resourceAlibabacloudStackDataWorksRemindRead(d, meta)
}

func resourceAlibabacloudStackDataWorksRemindDelete(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*connectivity.AlibabacloudStackClient)

	action := "DeleteRemind"
	request := map[string]interface{}{
		"RemindId": d.Id(),
		"RegionId": "default",
	}
	_, err = client.DoTeaRequest("POST", "dataworks-public", "2020-05-18", action, "", nil, nil, request)
	if err != nil {
		return err
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
