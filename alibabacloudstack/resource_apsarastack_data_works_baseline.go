package alibabacloudstack

import (
	"fmt"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackDataWorksBaseline() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"baseline_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the baseline.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the DataWorks project.",
			},
			"owner": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The owner of the baseline.",
			},
			"priority": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 8),
				Description:  "The priority of the baseline. Valid values: 1 to 8.",
			},
			"alert_margin_threshold": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The alert margin threshold in minutes.",
			},
			"baseline_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"DAILY", "HOURLY", "WEEKLY", "MONTHLY"}, false),
				Description:  "The type of the baseline. Valid values: DAILY, HOURLY, WEEKLY, MONTHLY.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether the baseline is enabled.",
			},
			"alert_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether alert is enabled for the baseline.",
			},
			"overtime_settings": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "The overtime settings for the baseline.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cycle": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The cycle number.",
						},
						"time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The time in HH:mm format.",
						},
					},
				},
			},
			// Computed attributes
			"baseline_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The ID of the created baseline.",
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackDataWorksBaselineCreate, resourceAlibabacloudStackDataWorksBaselineRead, resourceAlibabacloudStackDataWorksBaselineUpdate, resourceAlibabacloudStackDataWorksBaselineDelete)
	return resource
}

func resourceAlibabacloudStackDataWorksBaselineCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := make(map[string]interface{})

	// Map parameters using tagName from Create API
	request["BaselineName"] = d.Get("baseline_name").(string)
	request["ProjectId"] = d.Get("project_id").(int)
	request["Owner"] = convertAscmUid2MemberUid(d.Get("owner").(string))
	request["Priority"] = d.Get("priority").(int)
	request["BaselineType"] = d.Get("baseline_type").(string)

	if v, ok := d.GetOk("alert_margin_threshold"); ok {
		request["AlertMarginThreshold"] = v
	}

	// Handle OvertimeSettings (RepeatList with Json invokeDataType)
	overtimeSettings := d.Get("overtime_settings").([]interface{})
	if len(overtimeSettings) > 0 {
		var settings []map[string]interface{}
		for _, setting := range overtimeSettings {
			if settingMap, ok := setting.(map[string]interface{}); ok {
				item := make(map[string]interface{})
				if v, exists := settingMap["cycle"]; exists && v.(int) != 0 {
					item["Cycle"] = v.(int)
				}
				if v, exists := settingMap["time"]; exists && v.(string) != "" {
					item["Time"] = v.(string)
				}
				if len(item) > 0 {
					settings = append(settings, item)
				}
			}
		}
		request["OvertimeSettings"] = settings
	}

	reqQuery := map[string]interface{}{
		"regionId": client.RegionId,
	}

	response, err := client.DoTeaRequest("POST", "dataworks-public", "2020-05-18", "CreateBaseline", "", nil, reqQuery, request)
	if err != nil {
		return err
	}

	baselineId, err := toInt(response["Data"])
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg,
			"alibabacloudstack_dataworks_baseline", "CreateBaseline",
			"Failed to get baseline ID from response")
	}

	d.SetId(fmt.Sprintf("%d:%d", d.Get("project_id"), baselineId))
	return nil
}

func resourceAlibabacloudStackDataWorksBaselineRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	dataworksService := DataworksService{client}
	object, err := dataworksService.DescribeDataWorksBaseline(d.Id())
	if err != nil {
		return err
	}

	// Set computed attributes using tagName mappings
	d.Set("baseline_id", object["BaselineId"])
	d.Set("baseline_name", object["BaselineName"])
	d.Set("project_id", object["ProjectId"])
	d.Set("owner", object["Owner"].(string)[1:])
	d.Set("priority", object["Priority"])
	d.Set("baseline_type", object["BaselineType"])
	d.Set("enabled", object["Enabled"])
	d.Set("alert_enabled", object["AlertEnabled"])

	if alertMarginThreshold := object["AlertMarginThreshold"]; alertMarginThreshold != 0 {
		d.Set("alert_margin_threshold", alertMarginThreshold)
	}

	// Handle OverTimeSettings (Array)
	overtimeSettings := make([]interface{}, 0)
	if overTimeSettings, exists := object["OverTimeSettings"]; exists {
		if settingsList, ok := overTimeSettings.([]interface{}); ok {
			for _, setting := range settingsList {
				if settingMap, ok := setting.(map[string]interface{}); ok {
					overtimeSettings = append(overtimeSettings, map[string]interface{}{
						"cycle": settingMap["Cycle"],
						"time":  settingMap["Time"],
					})
				}
			}
		}
	}
	d.Set("overtime_settings", overtimeSettings)

	return nil
}

func resourceAlibabacloudStackDataWorksBaselineUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	update := false
	if _, ok := d.GetOk("enabled"); ok {
		update = true
	}

	if _, ok := d.GetOk("alert_enabled"); ok {
		update = true
	}

	if d.IsNewResource() && !update {
		return nil
	}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	request := map[string]interface{}{
		"BaselineId": parts[1],
		"ProjectId":  parts[0],
	}

	// Optional parameters - only include if changed
	if d.HasChange("baseline_name") {
		request["BaselineName"] = d.Get("baseline_name").(string)
	}
	if d.HasChange("owner") {
		request["Owner"] = convertAscmUid2MemberUid(d.Get("owner").(string))
	}
	if d.HasChange("priority") {
		request["Priority"] = d.Get("priority").(int)
	}
	if d.HasChange("alert_margin_threshold") {
		if v, ok := d.GetOk("alert_margin_threshold"); ok {
			request["AlertMarginThreshold"] = v.(int)
		}
	}
	if d.HasChange("baseline_type") {
		request["BaselineType"] = d.Get("baseline_type").(string)
	}
	if d.HasChange("enabled") {
		request["Enabled"] = d.Get("enabled").(bool)
	}
	if d.HasChange("alert_enabled") {
		request["AlertEnabled"] = d.Get("alert_enabled").(bool)
	}

	// Handle OvertimeSettings if changed
	if d.HasChange("overtime_settings") {
		overtimeSettings := d.Get("overtime_settings").([]interface{})
		if len(overtimeSettings) > 0 {
			var settings []map[string]interface{}
			for _, setting := range overtimeSettings {
				if settingMap, ok := setting.(map[string]interface{}); ok {
					item := make(map[string]interface{})
					if v, exists := settingMap["cycle"]; exists && v.(int) != 0 {
						item["Cycle"] = v.(int)
					}
					if v, exists := settingMap["time"]; exists && v.(string) != "" {
						item["Time"] = v.(string)
					}
					if len(item) > 0 {
						settings = append(settings, item)
					}
				}
			}
			request["OvertimeSettings"] = settings
		}
	}

	reqQuery := map[string]interface{}{
		"regionId": client.RegionId,
	}

	if _, err := client.DoTeaRequest("POST", "dataworks-public", "2020-05-18", "UpdateBaseline", "", nil, reqQuery, request); err != nil {
		return err
	}

	return nil
}

func resourceAlibabacloudStackDataWorksBaselineDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	request := map[string]interface{}{
		"BaselineId": parts[1],
		"ProjectId":  parts[0],
	}

	reqQuery := map[string]interface{}{
		"regionId": client.RegionId,
	}

	_, err = client.DoTeaRequest("POST", "dataworks-public", "2020-05-18", "DeleteBaseline", "", nil, reqQuery, request)
	return err
}
