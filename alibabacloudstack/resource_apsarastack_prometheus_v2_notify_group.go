package alibabacloudstack

import (
	"fmt"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackPrometheusV2NotifyGroup() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"DINGDING", "WECHAT_ROBOT", "CONTACT", "WEBHOOK"}, false),
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"webhook_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"webhook_header_params": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"im": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"contact_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackPrometheusV2NotifyGroupCreate, resourceAlibabacloudStackPrometheusV2NotifyGroupRead, resourceAlibabacloudStackPrometheusV2NotifyGroupUpdate, resourceAlibabacloudStackPrometheusV2NotifyGroupDelete)
	return resource
}

func resourceAlibabacloudStackPrometheusV2NotifyGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	groupType := d.Get("type").(string)
	body := map[string]interface{}{
		"type": d.Get("type"),
		"name": d.Get("name"),
	}
	if v, ok := d.GetOk("description"); ok {
		body["description"] = v
	}
	if groupType == "DINGDING" || groupType == "WECHAT_ROBOT" {
		im := d.Get("im").(string)
		if im == "" {
			return fmt.Errorf("im is required when type is %s", groupType)
		}
		body["im"] = im
	}
	if groupType == "CONTACT" {
		contactIds := d.Get("contact_ids").(*schema.Set).List()
		if len(contactIds) == 0 {
			return fmt.Errorf("contactIds is required when type is %s", groupType)
		}
		body["contactIds"] = contactIds
	}
	if groupType == "WEBHOOK" {
		webhook := make(map[string]interface{})
		webhookUrl := d.Get("webhook_url").(string)
		if webhookUrl == "" {
			return fmt.Errorf("webhook_url is required when type is %s", groupType)
		}
		webhook["url"] = webhookUrl
		if v, ok := d.GetOk("webhook_header_params"); ok && v.(*schema.Set).Len() > 0 {
			header_params := make(map[string]interface{})
			for _, headerParam := range v.(*schema.Set).List() {
				header := headerParam.(map[string]interface{})
				header_params[header["key"].(string)] = header["value"]
			}
			webhook["headerParams"] = header_params
		}
		body["webhook"] = webhook
	}

	resp, err := client.DoTeaRequest("POST", "prometheus2", "2023-04-13", "CreateNotifyGroup", "/log/api/v2/alert/group/add", nil, nil, body)
	if err != nil {
		return err
	}

	success, ok := resp["success"].(bool)
	if !ok || !success {
		return fmt.Errorf("failed to create notify group: %v", resp)
	}

	// Since the create API does not return the ID, we need to retrieve it by listing and filtering
	listResp, err := client.DoTeaRequest("GET", "prometheus2", "2023-04-13", "PageNotifyGroup", "/log/api/v2/alert/group/list", nil, map[string]interface{}{
		"keyword":   d.Get("name"),
		"page":      1,
		"pageSize":  10,
		"direction": "desc",
	}, nil)
	if err != nil {
		return err
	}

	data, ok := listResp["data"].([]interface{})
	if !ok {
		return fmt.Errorf("unexpected response format when retrieving notify group")
	}
	var groupId string
	for _, item := range data {
		if group, ok := item.(map[string]interface{}); ok {
			if name, ok := group["name"].(string); ok && name == d.Get("name") {
				groupId = fmt.Sprint(group["id"])
				break
			}
		}
	}
	if groupId == "" {
		return fmt.Errorf("failed to find notify group with name %s", d.Get("name"))
	}
	d.SetId(groupId)

	return nil
}

func resourceAlibabacloudStackPrometheusV2NotifyGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	prometheusService := PrometheusService{client}
	object, err := prometheusService.DescribePrometheusV2NotifyGroup(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("type", object["type"])
	d.Set("name", object["name"])
	d.Set("description", object["description"])
	d.Set("im", object["im"])
	if object["type"].(string) == "CONTACT" {
		contactIds, err := prometheusService.ListGroupContactIds(d.Id())
		if err != nil {
			return errmsgs.WrapError(err)
		}
		d.Set("contact_ids", contactIds)
	} else {
		d.Set("contact_ids", []string{})
	}
	if webhookData, ok := object["webhook"].(map[string]interface{}); ok {
		if url, ok := webhookData["url"].(string); ok {
			d.Set("webhook_url", url)
		}
		if headerParams, ok := webhookData["headerParams"].(map[string]interface{}); ok {
			headers := make([]map[string]interface{}, 0)
			for k, v := range headerParams {
				headers = append(headers, map[string]interface{}{
					"key":   k,
					"value": v,
				})
			}
			d.Set("webhook_header_params", headers)
		}
	}

	return nil
}

func resourceAlibabacloudStackPrometheusV2NotifyGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// If it is a new resource, do nothing
	if d.IsNewResource() {
		return nil
	}
	groupType := d.Get("type").(string)
	body := map[string]interface{}{
		"type":    d.Get("type"),
		"name":    d.Get("name"),
		"id":      d.Id(),
		"groupId": d.Id(),
	}
	if v, ok := d.GetOk("description"); ok {
		body["description"] = v
	}
	if groupType == "DINGDING" || groupType == "WECHAT_ROBOT" {
		im := d.Get("im").(string)
		if im == "" {
			return fmt.Errorf("im is required when type is %s", groupType)
		}
		body["im"] = im
	}
	if groupType == "CONTACT" {
		contactIds := d.Get("contact_ids").(*schema.Set).List()
		if len(contactIds) == 0 {
			return fmt.Errorf("contactIds is required when type is %s", groupType)
		}
		body["contactIds"] = contactIds
	}
	if groupType == "WEBHOOK" {
		webhook := make(map[string]interface{})
		webhookUrl := d.Get("webhook_url").(string)
		if webhookUrl == "" {
			return fmt.Errorf("webhook_url is required when type is %s", groupType)
		}
		webhook["url"] = webhookUrl
		if v, ok := d.GetOk("webhook_header_params"); ok && v.(*schema.Set).Len() > 0 {
			header_params := make(map[string]interface{})
			for _, headerParam := range v.(*schema.Set).List() {
				header := headerParam.(map[string]interface{})
				header_params[header["key"].(string)] = header["value"]
			}
			webhook["headerParams"] = header_params
		}
		body["webhook"] = webhook
	}
	response, err := client.DoTeaRequest("POST", "prometheus2", "2023-04-13", "UpdateNotifyGroup", "/log/api/v2/alert/group/update", nil, nil, body)
	if err != nil {
		return fmt.Errorf("failed to update notify group: %v", err)
	}
	if success, ok := response["Success"].(bool); ok && !success {
		return fmt.Errorf("failed to update notify group: %v", response)
	}

	return nil
}

func resourceAlibabacloudStackPrometheusV2NotifyGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	reqBody := map[string]interface{}{
		"ids": []string{d.Id()},
	}

	_, err := client.DoTeaRequest("POST", "prometheus2", "2023-04-13", "DeleteNotifyGroup", "/log/api/v2/alert/group/delete", nil, nil, reqBody)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "DeleteNotifyGroup", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return nil
}
