package alibabacloudstack

import (
	"fmt"
	"log"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackPrometheusV2Contact() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mobile": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mail": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
	setResourceFunc(resource,
		resourceAlibabacloudStackPrometheusV2ContactCreate,
		resourceAlibabacloudStackPrometheusV2ContactRead,
		resourceAlibabacloudStackPrometheusV2ContactUpdate,
		resourceAlibabacloudStackPrometheusV2ContactDelete,
	)
	return resource
}

func resourceAlibabacloudStackPrometheusV2ContactCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// Prepare request body
	body := make(map[string]interface{})
	body["username"] = d.Get("username")
	body["mobile"] = d.Get("mobile")
	body["mail"] = d.Get("mail")
	body["groupIds"] = []string{}

	resp, err := client.DoTeaRequest("POST", "prometheus2", "2023-04-13", "CreateContact", "/log/api/v2/alert/contact/add", nil, nil, body)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "CreateContact", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	success, ok := resp["success"].(bool)
	if !ok || !success {
		return fmt.Errorf("failed to create contact: %v", resp)
	}

	pageResp, err := client.DoTeaRequest("POST", "prometheus2", "2023-04-13", "PageContact", "/log/api/v2/alert/contact/page", nil, nil, map[string]interface{}{
		"keyword":         "",
		"pageNumber":      1,
		"pageSize":        10,
		"direction":       "desc",
		"page":            1,
		"selectedRowKeys": []string{},
	})
	if err != nil {
		return err
	}

	data, ok := pageResp["data"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("unexpected response format when retrieving contact list: %v", pageResp)
	}

	contactList, ok := data["data"].([]interface{})
	if !ok || len(contactList) == 0 {
		return fmt.Errorf("no contacts found after creation")
	}
	var contactId string
	username := d.Get("username").(string)
	for _, item := range contactList {
		contact, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		if contact["username"] == username {
			idVal, ok := contact["id"]
			if !ok {
				continue
			}
			contactId = fmt.Sprint(idVal)
			break
		}
	}

	if contactId == "" {
		return fmt.Errorf("failed to find created contact with username: %s", username)
	}
	log.Printf("[DEBUG] ============================= Created contact with ID: %s", contactId)
	d.SetId(contactId)
	return nil
}

func resourceAlibabacloudStackPrometheusV2ContactRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	prometheusService := PrometheusService{client}
	object, err := prometheusService.DescribePrometheusV2Contact(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("username", object["username"])
	d.Set("mobile", object["mobile"])
	d.Set("mail", object["mail"])

	return nil
}

func resourceAlibabacloudStackPrometheusV2ContactUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// If this is a new resource, no update is needed
	if d.IsNewResource() {
		return nil
	}

	updateReq := make(map[string]interface{})
	updateReq["id"] = d.Id()

	if d.HasChanges("username", "mobile", "mail") {

		prometheusService := PrometheusService{client}
		object, err := prometheusService.DescribePrometheusV2Contact(d.Id())
		if err != nil {
			return errmsgs.WrapError(err)
		}
		groups := object["groups"].([]interface{})
		var groupIdStrs []string
		for _, v := range groups {
			group := v.(map[string]interface{})
			groupIdStrs = append(groupIdStrs, fmt.Sprint(group["id"]))
		}
		updateReq["groupIds"] = groupIdStrs
		updateReq["username"] = d.Get("username")
		updateReq["mobile"] = d.Get("mobile")
		updateReq["mail"] = d.Get("mail")
		// Call UpdateContact API
		_, err = client.DoTeaRequest("POST", "prometheus2", "2023-04-13", "UpdateContact", "/log/api/v2/alert/contact/update", nil, nil, updateReq)
		if err != nil {
			return fmt.Errorf("failed to update contact: %v", err)
		}
	}
	return nil
}

func resourceAlibabacloudStackPrometheusV2ContactDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	contactId := d.Id()

	reqBody := map[string]interface{}{
		"ids": []string{contactId},
	}

	_, err := client.DoTeaRequest("POST", "prometheus2", "2023-04-13", "DeleteContact", "/log/api/v2/alert/contact/delete", nil, nil, reqBody)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, contactId, "DeleteContact", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return nil
}
