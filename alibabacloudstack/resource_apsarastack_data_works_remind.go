package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackDataWorksRemind() *schema.Resource {
	resource := &schema.Resource{
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
				Type:                  schema.TypeSet,
				Elem:                  &schema.Schema{Type: schema.TypeString},
				MinItems:              1,
				Optional:              true,
				ConflictsWith:         []string{"baseline_ids", "project_id", "biz_process_ids"},
				AtLeastOneOf:          []string{"baseline_ids", "project_id", "biz_process_ids"},
				DiffSuppressFunc:      remindTypeDiffSuppressFunc("NODE"),
				DiffSuppressOnRefresh: true,
			},
			"baseline_ids": {
				Type:                  schema.TypeSet,
				Elem:                  &schema.Schema{Type: schema.TypeString},
				MinItems:              1,
				Optional:              true,
				ConflictsWith:         []string{"node_ids", "project_id", "biz_process_ids"},
				AtLeastOneOf:          []string{"node_ids", "project_id", "biz_process_ids"},
				DiffSuppressFunc:      remindTypeDiffSuppressFunc("BASELINE"),
				DiffSuppressOnRefresh: true,
			},
			"project_id": {
				Type:                  schema.TypeString,
				Optional:              true,
				ConflictsWith:         []string{"node_ids", "baseline_ids", "biz_process_ids"},
				AtLeastOneOf:          []string{"node_ids", "baseline_ids", "biz_process_ids"},
				DiffSuppressFunc:      remindTypeDiffSuppressFunc("PROJECT"),
				DiffSuppressOnRefresh: true,
			},
			"biz_process_ids": {
				Type:                  schema.TypeSet,
				Elem:                  &schema.Schema{Type: schema.TypeString},
				MinItems:              1,
				Optional:              true,
				ConflictsWith:         []string{"node_ids", "baseline_ids", "project_id"},
				AtLeastOneOf:          []string{"node_ids", "baseline_ids", "project_id"},
				DiffSuppressFunc:      remindTypeDiffSuppressFunc("BIZPROCESS"),
				DiffSuppressOnRefresh: true,
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
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"MAIL", "SMS", "PHONE", "WEBHOOKS", "DINGROBOTS"}, false),
				},
				MinItems: 1,
				Required: true,
			},
			"alert_targets": {
				Type:                  schema.TypeSet,
				Elem:                  &schema.Schema{Type: schema.TypeString},
				Optional:              true,
				Computed:              true,
				DiffSuppressFunc:      diffWithRemindMethod([]string{"MAIL", "SMS", "PHONE"}),
				DiffSuppressOnRefresh: true,
			},
			"webhooks": {
				Type:                  schema.TypeSet,
				Elem:                  &schema.Schema{Type: schema.TypeString},
				MinItems:              1,
				Optional:              true,
				DiffSuppressFunc:      diffWithRemindMethod([]string{"WEBHOOKS"}),
				DiffSuppressOnRefresh: true,
			},
			"robot_urls": {
				Type:                  schema.TypeSet,
				Elem:                  &schema.Schema{Type: schema.TypeString},
				MinItems:              1,
				Optional:              true,
				DiffSuppressFunc:      diffWithRemindMethod([]string{"DINGROBOTS"}),
				DiffSuppressOnRefresh: true,
			},
			"use_flag": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackDataWorksRemindCreate, resourceAlibabacloudStackDataWorksRemindRead, resourceAlibabacloudStackDataWorksRemindUpdate, resourceAlibabacloudStackDataWorksRemindDelete)
	return resource
}

func diffWithRemindMethod(methods []string) func(k, oldValue, newValue string, d *schema.ResourceData) bool {
	return func(k, oldValue, newValue string, d *schema.ResourceData) bool {
		alertMethod := []string{}
		hited := false
		for _, v := range d.Get("alert_methods").(*schema.Set).List() {
			alertMethod = append(alertMethod, v.(string))
		}
		for _, method := range methods {
			if slices.Contains(alertMethod, method) {
				hited = true
				break
			}
		}

		field := strings.Split(k, ".")[0]
		if hited {
			o, n := d.GetChange(field)
			ov := []string{}
			for _, v := range o.(*schema.Set).List() {
				ov = append(ov, v.(string))
			}
			nv := []string{}
			for _, v := range n.(*schema.Set).List() {
				nv = append(nv, v.(string))
			}
			slices.Sort(ov)
			slices.Sort(nv)
			return slices.Equal(ov, nv)
		} else {
			return true
		}
	}
}

func resourceAlibabacloudStackDataWorksRemindCreate(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	action := "CreateRemind"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("alert_unit"); ok {
		request["AlertUnit"] = v.(string)
		if request["AlertUnit"] == "OTHER" {
			var list []string
			for _, vv := range d.Get("alert_targets").(*schema.Set).List() {
				list = append(list, vv.(string))
			}
			request["AlertTargets"] = strings.Join(list, ",")
		}
	}

	if v, ok := d.GetOk("remind_name"); ok {
		request["RemindName"] = v.(string)
	}

	if v, ok := d.GetOk("remind_type"); ok {
		request["RemindType"] = v.(string)
	}

	if v, ok := d.GetOk("remind_unit"); ok {
		request["RemindUnit"] = v.(string)
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
		var list []string
		for _, vv := range v.(*schema.Set).List() {
			list = append(list, vv.(string))
		}
		request["AlertMethods"] = strings.Join(list, ",")
	}

	if v, ok := d.GetOk("robot_urls"); ok {
		var list []string
		for _, vv := range v.(*schema.Set).List() {
			list = append(list, vv.(string))
		}
		request["RobotUrls"] = strings.Join(list, ",")
	}

	if v, ok := d.GetOk("webhooks"); ok {
		var list []string
		for _, vv := range v.(*schema.Set).List() {
			list = append(list, vv.(string))
		}
		request["Webhooks"] = strings.Join(list, ",")
	}

	response, err := client.DoTeaRequest("POST", "dataworks-public", "2020-05-18", action, "", nil, nil, request)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprint(response["Data"]))

	return nil
}

func resourceAlibabacloudStackDataWorksRemindRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	dataworksPublicService := DataworksService{client}
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
		log.Println(key, value)
	}

	d.Set("remind_id", d.Id())
	d.Set("alert_unit", object["AlertUnit"].(string))
	d.Set("remind_name", object["RemindName"].(string))
	d.Set("remind_type", object["RemindType"].(string))
	d.Set("remind_unit", object["RemindUnit"].(string))
	d.Set("dnd_end", object["DndEnd"].(string))

	d.Set("node_ids", nil)
	d.Set("baseline_ids", nil)
	d.Set("project_id", nil)
	d.Set("biz_process_ids", nil)
	switch object["RemindUnit"].(string) {
	case "NODE":
		var nodes []string
		for _, v := range object["Nodes"].([]interface{}) {
			item := v.(map[string]interface{})
			nodes = append(nodes, item["NodeId"].(string))
		}
		d.Set("node_ids", nodes)
	case "BASELINE":
		var baselines []string
		for _, v := range object["Baselines"].([]interface{}) {
			item := v.(map[string]interface{})
			if id, err := toInt(item["BaselineId"]); err == nil {
				baselines = append(baselines, strconv.Itoa(id))

			}
		}
		d.Set("baseline_ids", baselines)
	case "PROJECT":
		if len(object["Projects"].([]interface{})) > 0 {
			projectId := object["Projects"].([]interface{})[0].(map[string]interface{})["ProjectId"].(json.Number)
			d.Set("project_id", fmt.Sprint(projectId))
		}
	case "BIZPROCESS":
		var bizIds []string
		for _, v := range object["BizProcesses"].([]interface{}) {
			item := v.(map[string]interface{})
			bidId, err := toInt(item["BizId"])
			if err != nil {
				return err
			}
			bizIds = append(bizIds, strconv.Itoa(bidId))
		}
		d.Set("biz_process_ids", bizIds)
	}

	d.Set("max_alert_times", object["MaxAlertTimes"].(json.Number))
	d.Set("alert_interval", object["AlertInterval"].(json.Number))
	if object["RemindType"].(string) == "TIMEOUT" {
		n, _ := strconv.Atoi(object["Detail"].(string))
		d.Set("detail", fmt.Sprintf("%d", n*60))
	} else {
		d.Set("detail", object["Detail"].(string))
	}

	d.Set("alert_methods", object["AlertMethods"])
	d.Set("alert_targets", object["AlertTargets"])
	var ruls []string
	for _, v := range object["Robots"].([]interface{}) {
		item := v.(map[string]interface{})
		ruls = append(ruls, item["WebUrl"].(string))
	}
	d.Set("robot_urls", ruls)
	d.Set("webhooks", object["Webhooks"])
	d.Set("use_flag", object["Useflag"].(bool))

	return nil
}

func resourceAlibabacloudStackDataWorksRemindUpdate(d *schema.ResourceData, meta interface{}) (err error) {

	if d.IsNewResource() {
		return nil
	}

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
		request["RemindUnit"] = d.Get("remind_unit").(string)

		buildRemindUnitArgs(d, request)
	}

	if d.HasChange("remind_type") {
		request["RemindType"] = d.Get("remind_type").(string)
		request["Detail"] = d.Get("detail").(string)
	}

	if d.HasChange("max_alert_times") {
		request["MaxAlertTimes"] = d.Get("max_alert_times").(int)
	}

	if d.HasChanges("alert_unit", "alert_targets") {
		request["AlertUnit"] = d.Get("alert_unit").(string)
		if request["AlertUnit"] == "OTHER" {
			var list []string
			for _, vv := range d.Get("alert_targets").(*schema.Set).List() {
				list = append(list, vv.(string))
			}
			request["AlertTargets"] = strings.Join(list, ",")
		}
	}

	if d.HasChange("alert_methods") {
		var list []string
		for _, vv := range d.Get("alert_methods").(*schema.Set).List() {
			list = append(list, vv.(string))
		}
		request["AlertMethods"] = strings.Join(list, ",")
	}

	if d.HasChange("use_flag") {
		request["UseFlag"] = d.Get("use_flag").(bool)
	}

	if d.HasChange("robot_urls") {
		var list []string
		for _, vv := range d.Get("robot_urls").(*schema.Set).List() {
			list = append(list, vv.(string))
		}
		request["RobotUrls"] = strings.Join(list, ",")
	}

	if d.HasChange("webhooks") {
		var list []string
		for _, vv := range d.Get("webhooks").(*schema.Set).List() {
			list = append(list, vv.(string))
		}
		request["Webhooks"] = strings.Join(list, ",")
	}

	action := "UpdateRemind"
	_, err = client.DoTeaRequest("POST", "dataworks-public", "2020-05-18", action, "", nil, nil, request)
	if err != nil {
		return err
	}
	return nil
}

func resourceAlibabacloudStackDataWorksRemindDelete(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*connectivity.AlibabacloudStackClient)

	action := "DeleteRemind"
	request := map[string]interface{}{
		"RemindId": d.Id(),
	}
	_, err = client.DoTeaRequest("POST", "dataworks-public", "2020-05-18", action, "", nil, nil, request)
	if err != nil {
		return err
	}
	return nil
}

func buildRemindUnitArgs(d *schema.ResourceData, request map[string]interface{}) {
	switch d.Get("remind_unit").(string) {
	case "NODE":
		var list []string
		for _, v := range d.Get("node_ids").(*schema.Set).List() {
			list = append(list, v.(string))
		}
		request["NodeIds"] = strings.Join(list, ",")
	case "BASELINE":
		var list []string
		for _, v := range d.Get("baseline_ids").(*schema.Set).List() {
			list = append(list, v.(string))
		}
		request["BaselineIds"] = strings.Join(list, ",")
	case "PROJECT":
		request["ProjectId"] = d.Get("project_id").(string)
	case "BIZPROCESS":
		var list []string
		for _, v := range d.Get("biz_process_ids").(*schema.Set).List() {
			list = append(list, v.(string))
		}
		request["BizProcessIds"] = strings.Join(list, ",")
	}
}
