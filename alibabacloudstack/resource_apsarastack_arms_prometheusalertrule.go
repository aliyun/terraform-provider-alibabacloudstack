package alibabacloudstack

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackArmsPrometheusAlertRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackArmsPrometheusAlertRuleCreate,
		Read:   resourceAlibabacloudStackArmsPrometheusAlertRuleRead,
		Update: resourceAlibabacloudStackArmsPrometheusAlertRuleUpdate,
		Delete: resourceAlibabacloudStackArmsPrometheusAlertRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"annotations": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dispatch_rule_id": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("notify_type"); ok && v.(string) == "DISPATCH_RULE" {
						return false
					}
					return true
				},
			},
			"duration": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"expression": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"labels": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"message": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"notify_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ALERT_MANAGER", "DISPATCH_RULE"}, false),
			},
			"prometheus_alert_rule_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
		},
	}
}

func resourceAlibabacloudStackArmsPrometheusAlertRuleCreate(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var response map[string]interface{}
	action := "CreatePrometheusAlertRule"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("annotations"); ok {
		annotationsMaps := make([]map[string]interface{}, 0)
		for _, annotations := range v.(*schema.Set).List() {
			annotationsMap := annotations.(map[string]interface{})
			annotationsMaps = append(annotationsMaps, annotationsMap)
		}
		if v, err := convertArrayObjectToJsonString(annotationsMaps); err == nil {
			request["Annotations"] = v
		} else {
			return errmsgs.WrapError(err)
		}
	}
	request["ClusterId"] = d.Get("cluster_id")
	if v, ok := d.GetOk("dispatch_rule_id"); ok {
		request["DispatchRuleId"] = v
	}
	request["Duration"] = d.Get("duration")
	request["Expression"] = d.Get("expression")
	if v, ok := d.GetOk("labels"); ok {
		labelsMaps := make([]map[string]interface{}, 0)
		for _, labels := range v.(*schema.Set).List() {
			labelsMap := labels.(map[string]interface{})
			labelsMaps = append(labelsMaps, labelsMap)
		}
		if v, err := convertArrayObjectToJsonString(labelsMaps); err == nil {
			request["Labels"] = v
		} else {
			return errmsgs.WrapError(err)
		}
	}
	request["Message"] = d.Get("message")
	if v, ok := d.GetOk("notify_type"); ok {
		request["NotifyType"] = v
	}
	request["AlertName"] = d.Get("prometheus_alert_rule_name")
	if v, ok := d.GetOk("type"); ok {
		request["Type"] = v
	}
	response, err = client.DoTeaRequest("POST", "ARMS", "2019-08-08", action, "", nil, request)
	if err != nil {
		return err
	}
	responsePrometheusAlertRule := response["PrometheusAlertRule"].(map[string]interface{})
	d.SetId(fmt.Sprint(request["ClusterId"], ":", responsePrometheusAlertRule["AlertId"]))

	return resourceAlibabacloudStackArmsPrometheusAlertRuleRead(d, meta)
}

func resourceAlibabacloudStackArmsPrometheusAlertRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	armsService := ArmsService{client}
	object, err := armsService.DescribeArmsPrometheusAlertRule(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_arms_prometheus_alert_rule armsService.DescribeArmsPrometheusAlertRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	d.Set("cluster_id", parts[0])
	d.Set("prometheus_alert_rule_id", parts[1])
	if v, ok := object["Annotations"].([]interface{}); ok {
		annotations := make([]map[string]interface{}, 0)
		for _, val := range v {
			item := val.(map[string]interface{})
			if item["Name"] == "message" {
				continue
			}

			temp := map[string]interface{}{
				"name":  item["Name"],
				"value": item["Value"],
			}

			annotations = append(annotations, temp)
		}
		if err := d.Set("annotations", annotations); err != nil {
			return errmsgs.WrapError(err)
		}
	}
	d.Set("dispatch_rule_id", fmt.Sprint(formatInt(object["DispatchRuleId"])))
	d.Set("duration", object["Duration"])
	d.Set("expression", object["Expression"])
	if v, ok := object["Labels"].([]interface{}); ok {
		labels := make([]map[string]interface{}, 0)
		for _, val := range v {
			item := val.(map[string]interface{})
			temp := map[string]interface{}{
				"name":  item["Name"],
				"value": item["Value"],
			}

			labels = append(labels, temp)
		}
		if err := d.Set("labels", labels); err != nil {
			return errmsgs.WrapError(err)
		}
	}
	d.Set("message", object["Message"])
	d.Set("notify_type", object["NotifyType"])
	d.Set("prometheus_alert_rule_name", object["AlertName"])
	d.Set("status", fmt.Sprint(formatInt(object["Status"])))
	d.Set("type", object["Type"])
	return nil
}

func resourceAlibabacloudStackArmsPrometheusAlertRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	update := false
	request := map[string]interface{}{
		"ClusterId": parts[0],
		"AlertId":   parts[1],
	}

	request["Duration"] = d.Get("duration")
	if d.HasChange("duration") {
		update = true
	}
	request["Expression"] = d.Get("expression")
	if d.HasChange("expression") {
		update = true
	}
	request["Message"] = d.Get("message")
	if d.HasChange("message") {
		update = true
	}
	request["AlertName"] = d.Get("prometheus_alert_rule_name")
	if d.HasChange("prometheus_alert_rule_name") {
		update = true
	}
	if d.HasChange("annotations") {
		update = true
		if v, ok := d.GetOk("annotations"); ok {
			annotationsMaps := make([]map[string]interface{}, 0)
			for _, annotations := range v.(*schema.Set).List() {
				annotationsMap := annotations.(map[string]interface{})
				annotationsMaps = append(annotationsMaps, annotationsMap)
			}
			if v, err := convertArrayObjectToJsonString(annotationsMaps); err == nil {
				request["Annotations"] = v
			} else {
				return errmsgs.WrapError(err)
			}
		}
	}
	if d.HasChange("dispatch_rule_id") {
		update = true
		if v, ok := d.GetOk("dispatch_rule_id"); ok {
			request["DispatchRuleId"] = v
		}
	}
	if d.HasChange("labels") {
		update = true
		if v, ok := d.GetOk("labels"); ok {
			labelsMaps := make([]map[string]interface{}, 0)
			for _, labels := range v.(*schema.Set).List() {
				labelsMap := labels.(map[string]interface{})
				labelsMaps = append(labelsMaps, labelsMap)
			}
			if v, err := convertArrayObjectToJsonString(labelsMaps); err == nil {
				request["Labels"] = v
			} else {
				return errmsgs.WrapError(err)
			}
		}
	}
	if d.HasChange("notify_type") {
		update = true
		if v, ok := d.GetOk("notify_type"); ok {
			request["NotifyType"] = v
		}
	}

	if update {
		action := "UpdatePrometheusAlertRule"
		_, err = client.DoTeaRequest("POST", "ARMS", "2019-08-08", action, "", nil, request)
		if err != nil {
			return err
		}
	}
	return resourceAlibabacloudStackArmsPrometheusAlertRuleRead(d, meta)
}

func resourceAlibabacloudStackArmsPrometheusAlertRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	action := "DeletePrometheusAlertRule"
	request := map[string]interface{}{
		"AlertId": parts[1],
	}
	_, err = client.DoTeaRequest("POST", "ARMS", "2019-08-08", action, "", nil, request)
	if err != nil {
		return err
	}
	return nil
}
