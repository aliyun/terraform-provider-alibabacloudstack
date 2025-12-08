package alibabacloudstack

import (
	"encoding/json"
	"fmt"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackPrometheusV2Alert() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"notify_recovered": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"is_check_all": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"trigger_clusters": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"trigger_promql": {
				Type:     schema.TypeString,
				Required: true,
			},
			"trigger_period": {
				Type:     schema.TypeString,
				Required: true,
			},
			"trigger_severity": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"warning", "serious", "fatal"}, false),
			},
			"trigger_cron": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tag_set": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"notify_types": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"notify_group_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"notify_interval": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"notification": {
				Type:     schema.TypeString,
				Required: true,
			},
			"recover_notification": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackPrometheusV2AlertCreate, resourceAlibabacloudStackPrometheusV2AlertRead, resourceAlibabacloudStackPrometheusV2AlertUpdate, resourceAlibabacloudStackPrometheusV2AlertDelete)
	return resource
}

func resourceAlibabacloudStackPrometheusV2AlertCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// Prepare request body
	if !d.Get("is_check_all").(bool) && d.Get("trigger_clusters").(*schema.Set).Len() == 0 {
		return fmt.Errorf("trigger_clusters must be set when is_check_all is false")
	}
	if d.Get("notify_recovered").(bool) && d.Get("recover_notification").(string) == "" {
		return fmt.Errorf("recover_notification must be set when notify_recovered is true")
	}
	body := map[string]interface{}{
		"name":            d.Get("name").(string),
		"alertType":       "PROMETHEUS",
		"notifyRecovered": d.Get("notify_recovered"),
		"isCheckAll":      d.Get("is_check_all"),
		"triggerRule": map[string]interface{}{
			"clusterIds": d.Get("trigger_clusters").(*schema.Set).List(),
			"promql":     d.Get("trigger_promql").(string),
			"period":     d.Get("trigger_period").(string),
			"severity":   d.Get("trigger_severity").(string),
			"cron":       d.Get("trigger_cron").(string),
			"timeType":   1,
		},
		"notify_keys": []int{0},
	}
	if v, ok := d.GetOk("tag_set"); ok {
		body["tagSet"] = v.(*schema.Set).List()
	}
	alertNotifyParam := make(map[string]interface{})
	if v, ok := d.GetOk("notify_types"); ok {
		alertNotifyParam["notifyTypes"] = v.(*schema.Set).List()
	}
	if v, ok := d.GetOk("notify_group_ids"); ok {
		alertNotifyParam["notifyGroupIds"] = v.(*schema.Set).List()
	}
	if v, ok := d.GetOk("notify_interval"); ok {
		alertNotifyParam["notifyInterval"] = v.(string)
	}
	body["alertNotifyParams"] = []interface{}{alertNotifyParam}
	if v, ok := d.GetOk("notification"); ok {
		body["notification"] = map[string]interface{}{
			"message": v,
		}
	}
	if v, ok := d.GetOk("recover_notification"); ok {
		body["recoverNotification"] = v.(string)
	}

	resp, err := client.DoTeaRequest("POST", "prometheus2", "2023-04-13", "CreateAlert", "/log/api/v2/alert/save", nil, nil, body)
	if err != nil {
		return err
	}

	success, ok := resp["success"].(bool)
	if !ok || !success {
		return fmt.Errorf("failed to create alert: %v", resp)
	}

	// Retrieve Alert ID by querying PageAlert with the name filter
	queryBody := map[string]interface{}{
		"keyword":    "",
		"pageNumber": 1,
		"pageSize":   10,
	}
	pageResp, err := client.DoTeaRequest("POST", "prometheus2", "2023-04-13", "PageAlert", "/log/api/v2/alert/page", nil, nil, queryBody)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	pageData, err := jsonpath.Get("$.data.data", pageResp)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, "PageAlert", "$.data.data", pageResp)
	}

	resourceList, ok := pageData.([]interface{})
	if !ok || len(resourceList) == 0 {
		return fmt.Errorf("no alert found after creation")
	}

	var alertId string
	for _, res := range resourceList {
		if r, ok := res.(map[string]interface{}); ok {
			if n, ok := r["name"].(string); ok && n == d.Get("name").(string) {
				if idVal, ok := r["id"].(json.Number); ok {
					alertId = fmt.Sprint(idVal)
					break
				}
			}
		}
	}

	if alertId == "" {
		return fmt.Errorf("could not retrieve alert ID after creation")
	}

	d.SetId(alertId)
	return nil
}

func resourceAlibabacloudStackPrometheusV2AlertRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	prometheusService := PrometheusService{client}
	object, err := prometheusService.DescribePrometheusV2Alert(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}
	d.Set("notify_recovered", object["notifyRecovered"])
	d.Set("name", object["name"])
	d.Set("tag_set", object["tagSet"])
	d.Set("recover_notification", object["recoverNotification"])

	if triggerRule, ok := object["triggerRule"].(map[string]interface{}); ok {
		is_check_all := false
		if v, ok := triggerRule["clusterIds"]; ok {
			d.Set("trigger_clusters", v)
			if len(v.([]interface{})) == 0 {
				is_check_all = true
			}
		}
		d.Set("is_check_all", is_check_all)
		if v, ok := triggerRule["promql"]; ok {
			d.Set("trigger_promql", v)
		}
		if v, ok := triggerRule["period"]; ok {
			d.Set("trigger_period", v)
		}
		if v, ok := triggerRule["severity"]; ok {
			d.Set("trigger_severity", v)
		}
		if v, ok := triggerRule["cron"]; ok {
			d.Set("trigger_cron", v)
		}
	}

	if notification, ok := object["notification"].(map[string]interface{}); ok {
		d.Set("notification", notification["message"])
	}

	if alertNotifies, ok := object["alertNotifies"].([]interface{}); ok && len(alertNotifies) > 0 {
		if notify, ok := alertNotifies[0].(map[string]interface{}); ok {
			d.Set("notify_types", notify["notifyTypes"])
			d.Set("notify_group_ids", notify["notifyGroupIds"])
			d.Set("notify_interval", notify["notifyInterval"])
		}
	}

	return nil
}

func resourceAlibabacloudStackPrometheusV2AlertUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	if d.IsNewResource() {
		return nil
	}
	if !d.Get("is_check_all").(bool) && d.Get("trigger_clusters").(*schema.Set).Len() == 0 {
		return fmt.Errorf("trigger_clusters must be set when is_check_all is false")
	}
	if d.Get("notify_recovered").(bool) && d.Get("recover_notification").(string) == "" {
		return fmt.Errorf("recover_notification must be set when notify_recovered is true")
	}
	body := map[string]interface{}{
		"id":              d.Id(),
		"name":            d.Get("name").(string),
		"alertType":       "PROMETHEUS",
		"notifyRecovered": d.Get("notify_recovered"),
		"isCheckAll":      d.Get("is_check_all"),
		"triggerRule": map[string]interface{}{
			"clusterIds": d.Get("trigger_clusters").(*schema.Set).List(),
			"promql":     d.Get("trigger_promql").(string),
			"period":     d.Get("trigger_period").(string),
			"severity":   d.Get("trigger_severity").(string),
			"cron":       d.Get("trigger_cron").(string),
			"timeType":   1,
		},
		"notify_keys": []int{0},
	}
	if v, ok := d.GetOk("tag_set"); ok {
		body["tagSet"] = v.(*schema.Set).List()
	}
	alertNotifyParam := make(map[string]interface{})
	if v, ok := d.GetOk("notify_types"); ok {
		alertNotifyParam["notifyTypes"] = v.(*schema.Set).List()
	}
	if v, ok := d.GetOk("notify_group_ids"); ok {
		alertNotifyParam["notifyGroupIds"] = v.(*schema.Set).List()
	}
	if v, ok := d.GetOk("notify_interval"); ok {
		alertNotifyParam["notifyInterval"] = v.(string)
	}
	body["alertNotifyParams"] = []interface{}{alertNotifyParam}
	if v, ok := d.GetOk("notification"); ok {
		body["notification"] = map[string]interface{}{
			"message": v,
		}
	}
	if v, ok := d.GetOk("recover_notification"); ok {
		body["recoverNotification"] = v.(string)
	}

	_, err := client.DoTeaRequest("POST", "prometheus2", "2023-04-13", "UpdateAlert", "/log/api/v2/alert/update", nil, nil, body)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_prometheus_v2_alert", "UpdateAlert", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return nil
}

func resourceAlibabacloudStackPrometheusV2AlertDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// Prepare the request body with the alert ID
	reqBody := map[string]interface{}{
		"ids": []string{d.Id()},
	}

	// Call the DeleteAlert API
	raw, err := client.DoTeaRequest("POST", "prometheus2", "2023-04-13", "DeleteAlert", "/log/api/v2/alert/delete", nil, nil, reqBody)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "DeleteAlert", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	// Check if the deletion was successful based on the response
	successRaw, ok := raw["success"]
	if !ok {
		return errmsgs.WrapErrorf(fmt.Errorf("missing 'success' field in response"), errmsgs.DefaultErrorMsg, d.Id(), "DeleteAlert", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	success, ok := successRaw.(bool)
	if !ok || !success {
		messageRaw, msgOk := raw["message"]
		message := ""
		if msgOk {
			message, _ = messageRaw.(string)
		}
		return errmsgs.WrapErrorf(fmt.Errorf("delete failed: %s", message), errmsgs.DefaultErrorMsg, d.Id(), "DeleteAlert", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return nil
}
