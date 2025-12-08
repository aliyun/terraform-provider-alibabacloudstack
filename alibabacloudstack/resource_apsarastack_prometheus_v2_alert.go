package alibabacloudstack

import (
	"fmt"
	"strconv"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackPrometheusV2Alert() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"alert_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"notify_recovered": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"is_check_all": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"trigger_rule": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     schema.TypeString,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tag_set": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"notify_keys": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"alert_notify_params": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
					},
				},
			},
			"notification": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     schema.TypeString,
			},
			"recover_notification": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackPrometheusV2AlertCreate, resourceAlibabacloudStackPrometheusV2AlertRead, resourceAlibabacloudStackPrometheusV2AlertUpdate, resourceAlibabacloudStackPrometheusV2AlertDelete)
	return resource
}

func resourceAlibabacloudStackPrometheusV2AlertCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// Prepare request body
	body := make(map[string]interface{})

	if v, ok := d.GetOk("alert_type"); ok {
		body["alertType"] = v.(string)
	}
	if v, ok := d.GetOk("notify_recovered"); ok {
		body["notifyRecovered"] = v.(bool)
	}
	if v, ok := d.GetOk("is_check_all"); ok {
		body["isCheckAll"] = v.(bool)
	}
	if v, ok := d.GetOk("trigger_rule"); ok {
		triggerRule := make(map[string]interface{})
		for key, value := range v.(map[string]interface{}) {
			triggerRule[key] = value
		}
		body["triggerRule"] = triggerRule
	}
	if v, ok := d.GetOk("name"); ok {
		body["name"] = v.(string)
	}
	if v, ok := d.GetOk("tag_set"); ok {
		tagSet := make([]interface{}, 0)
		for _, item := range v.([]interface{}) {
			tagSet = append(tagSet, item.(string))
		}
		body["tagSet"] = tagSet
	}
	if v, ok := d.GetOk("notify_keys"); ok {
		notifyKeys := make([]interface{}, 0)
		for _, item := range v.([]interface{}) {
			notifyKeys = append(notifyKeys, item)
		}
		body["notifyKeys"] = notifyKeys
	}
	if v, ok := d.GetOk("alert_notify_params"); ok {
		alertNotifyParams := make([]map[string]interface{}, 0)
		for _, item := range v.(*schema.Set).List() {
			param := item.(map[string]interface{})
			mappedParam := make(map[string]interface{})
			if notifyTypes, ok := param["notify_types"]; ok {
				types := make([]interface{}, 0)
				for _, t := range notifyTypes.(*schema.Set).List() {
					types = append(types, t.(string))
				}
				mappedParam["notifyTypes"] = types
			}
			if groupIds, ok := param["notify_group_ids"]; ok {
				ids := make([]interface{}, 0)
				for _, id := range groupIds.(*schema.Set).List() {
					ids = append(ids, id)
				}
				mappedParam["notifyGroupIds"] = ids
			}
			if interval, ok := param["notify_interval"]; ok {
				mappedParam["notifyInterval"] = interval.(string)
			}
			alertNotifyParams = append(alertNotifyParams, mappedParam)
		}
		body["alertNotifyParams"] = alertNotifyParams
	}
	if v, ok := d.GetOk("notification"); ok {
		notification := make(map[string]interface{})
		for key, value := range v.(map[string]interface{}) {
			notification[key] = value
		}
		body["notification"] = notification
	}
	if v, ok := d.GetOk("recover_notification"); ok {
		body["recoverNotification"] = v.(string)
	}

	// Call CreateAlert API
	resp, err := client.DoTeaRequest("POST", "prometheus2", "2023-04-13", "CreateAlert", "/log/api/v2/alert/save", nil, nil, body)
	if err != nil {
		return err
	}

	success, ok := resp["success"].(bool)
	if !ok || !success {
		return fmt.Errorf("failed to create alert: %v", resp)
	}

	// Retrieve Alert ID by querying PageAlert with the name filter
	name := d.Get("name").(string)
	pageResp, err := client.DoTeaRequest("GET", "prometheus2", "2023-04-13", "PageAlert", "/log/api/v2/alert/page", nil, map[string]interface{}{"name": name}, nil)
	if err != nil {
		return err
	}

	pageData, ok := pageResp["data"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("unexpected response format when retrieving alert list")
	}

	resourceList, ok := pageData["resource"].([]interface{})
	if !ok || len(resourceList) == 0 {
		return fmt.Errorf("no alert found after creation")
	}

	var alertId string
	for _, res := range resourceList {
		if r, ok := res.(map[string]interface{}); ok {
			if n, ok := r["name"].(string); ok && n == name {
				if idVal, ok := r["id"].(float64); ok {
					alertId = strconv.Itoa(int(idVal))
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

	d.Set("alert_type", object["alertType"])
	d.Set("notify_recovered", object["notifyRecovered"])
	d.Set("name", object["name"])
	d.Set("tag_set", object["tagSet"])
	d.Set("recover_notification", object["recoverNotification"])

	if triggerRule, ok := object["triggerRule"].(map[string]interface{}); ok {
		d.Set("trigger_rule", triggerRule)
	}

	if notification, ok := object["notification"].(map[string]interface{}); ok {
		d.Set("notification", notification)
	}

	if alertNotifies, ok := object["alertNotifies"].([]interface{}); ok && len(alertNotifies) > 0 {
		alertNotifyParams := make([]map[string]interface{}, 0, len(alertNotifies))
		for _, item := range alertNotifies {
			if notify, ok := item.(map[string]interface{}); ok {
				param := map[string]interface{}{
					"notify_types":     notify["notifyTypes"],
					"notify_group_ids": notify["notifyGroupIds"],
					"notify_interval":  notify["notifyInterval"],
				}
				alertNotifyParams = append(alertNotifyParams, param)
			}
		}
		d.Set("alert_notify_params", alertNotifyParams)
	}

	return nil
}

func resourceAlibabacloudStackPrometheusV2AlertUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// If it is a new resource, do nothing
	if d.IsNewResource() {
		return nil
	}

	updateFields := []string{
		"alert_type",
		"notify_recovered",
		"name",
		"tag_set",
		"notify_keys",
		"alert_notify_params",
		"trigger_rule",
		"notification",
		"recover_notification",
	}

	if d.HasChanges(updateFields...) {
		reqBody := make(map[string]interface{})
		reqBody["id"] = d.Id()

		if v, ok := d.GetOk("alert_type"); ok {
			reqBody["alertType"] = v
		}
		if v, ok := d.GetOkExists("notify_recovered"); ok {
			reqBody["notifyRecovered"] = v
		}
		if v, ok := d.GetOk("name"); ok {
			reqBody["name"] = v
		}
		if v, ok := d.GetOk("tag_set"); ok {
			tagSet := v.([]interface{})
			tags := make([]string, len(tagSet))
			for i, t := range tagSet {
				tags[i] = t.(string)
			}
			reqBody["tagSet"] = tags
		}
		if v, ok := d.GetOk("notify_keys"); ok {
			keys := v.([]interface{})
			intKeys := make([]int, len(keys))
			for i, k := range keys {
				intKeys[i] = k.(int)
			}
			reqBody["notifyKeys"] = intKeys
		}
		if v, ok := d.GetOk("alert_notify_params"); ok {
			params := v.(*schema.Set).List()
			alertNotifyParams := make([]map[string]interface{}, len(params))
			for i, p := range params {
				param := p.(map[string]interface{})
				item := make(map[string]interface{})
				if notifyTypes, ok := param["notify_types"]; ok {
					types := notifyTypes.(*schema.Set).List()
					strTypes := make([]string, len(types))
					for j, t := range types {
						strTypes[j] = t.(string)
					}
					item["notifyTypes"] = strTypes
				}
				if groupIds, ok := param["notify_group_ids"]; ok {
					ids := groupIds.(*schema.Set).List()
					intIds := make([]int, len(ids))
					for j, id := range ids {
						intIds[j] = id.(int)
					}
					item["notifyGroupIds"] = intIds
				}
				if interval, ok := param["notify_interval"]; ok && interval.(string) != "" {
					item["notifyInterval"] = interval
				}
				alertNotifyParams[i] = item
			}
			reqBody["alertNotifyParams"] = alertNotifyParams
		}
		if v, ok := d.GetOk("trigger_rule"); ok {
			triggerRule := v.(map[string]interface{})
			rule := make(map[string]interface{})
			for key, value := range triggerRule {
				rule[key] = value
			}
			reqBody["triggerRule"] = rule
		}
		if v, ok := d.GetOk("notification"); ok {
			notification := v.(map[string]interface{})
			note := make(map[string]interface{})
			for key, value := range notification {
				note[key] = value
			}
			reqBody["notification"] = note
		}
		if v, ok := d.GetOk("recover_notification"); ok {
			reqBody["recoverNotification"] = v
		}

		_, err := client.DoTeaRequest("POST", "prometheus2", "2023-04-13", "UpdateAlert", "/log/api/v2/alert/update", nil, nil, reqBody)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_prometheus_v2_alert", "UpdateAlert", errmsgs.AlibabacloudStackSdkGoERROR)
		}
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

	// Optionally wait for the resource to be fully deleted if needed
	// Since there's no specific state check method provided, we assume deletion is immediate
	return nil
}
