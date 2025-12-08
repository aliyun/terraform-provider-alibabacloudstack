package alibabacloudstack

import (
	"fmt"

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
			"group_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}

	resource.Create = resourceAlibabacloudStackPrometheusV2ContactCreate
	resource.Read = resourceAlibabacloudStackPrometheusV2ContactRead
	resource.Update = resourceAlibabacloudStackPrometheusV2ContactUpdate
	resource.Delete = resourceAlibabacloudStackPrometheusV2ContactDelete

	return resource
}

func resourceAlibabacloudStackPrometheusV2ContactCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// Prepare request body
	body := make(map[string]interface{})
	body["username"] = d.Get("username")
	body["mobile"] = d.Get("mobile")
	body["mail"] = d.Get("mail")

	groupIds := d.Get("group_ids").([]interface{})
	if len(groupIds) > 0 {
		groupIdList := make([]string, 0, len(groupIds))
		for _, id := range groupIds {
			groupIdList = append(groupIdList, id.(string))
		}
		body["groupIds"] = groupIdList
	} else {
		body["groupIds"] = []string{}
	}

	// Call CreateContact API
	resp, err := client.DoTeaRequest("POST", "prometheus2", "2023-04-13", "CreateContact", "/log/api/v2/alert/contact/add", nil, nil, body)
	if err != nil {
		return err
	}

	// Check success field in response
	success, ok := resp["success"].(bool)
	if !ok || !success {
		return fmt.Errorf("failed to create contact: %v", resp)
	}

	// Since the create API does not return the contact ID, we need to retrieve it by listing contacts
	// and filtering based on username.
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
			contactId = fmt.Sprintf("%v", idVal)
			break
		}
	}

	if contactId == "" {
		return fmt.Errorf("failed to find created contact with username: %s", username)
	}

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

	if groupIds, ok := object["groupIds"].([]interface{}); ok {
		groupIdStrings := make([]string, len(groupIds))
		for i, gid := range groupIds {
			groupIdStrings[i] = fmt.Sprintf("%v", gid)
		}
		d.Set("group_ids", groupIdStrings)
	}

	if groups, ok := object["groups"].([]interface{}); ok {
		groupStrings := make([]string, len(groups))
		for i, g := range groups {
			groupStrings[i] = fmt.Sprintf("%v", g)
		}
		d.Set("groups", groupStrings)
	}

	d.Set("id", fmt.Sprintf("%.0f", object["id"].(float64)))

	return nil
}

func resourceAlibabacloudStackPrometheusV2ContactUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// If this is a new resource, no update is needed
	if d.IsNewResource() {
		return nil
	}

	// Prepare update request body
	updateReq := make(map[string]interface{})
	updateReq["id"] = d.Id()

	if d.HasChange("username") {
		updateReq["username"] = d.Get("username")
	}
	if d.HasChange("mobile") {
		updateReq["mobile"] = d.Get("mobile")
	}
	if d.HasChange("mail") {
		updateReq["mail"] = d.Get("mail")
	}
	if d.HasChange("group_ids") {
		groupIds := d.Get("group_ids").([]interface{})
		var groupIdStrs []string
		for _, gid := range groupIds {
			groupIdStrs = append(groupIdStrs, gid.(string))
		}
		updateReq["groupIds"] = groupIdStrs
	}

	// Call UpdateContact API
	_, err := client.DoTeaRequest("POST", "prometheus2", "2023-04-13", "UpdateContact", "/log/api/v2/alert/contact/update", nil, nil, updateReq)
	if err != nil {
		return fmt.Errorf("failed to update contact: %v", err)
	}

	return resourceAlibabacloudStackPrometheusV2ContactRead(d, meta)
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
