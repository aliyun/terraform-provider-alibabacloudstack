package alibabacloudstack

import (
	"fmt"
	"log"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackArmsAlertContactGroup() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"alert_contact_group_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"contact_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackArmsAlertContactGroupCreate,
		resourceAlibabacloudStackArmsAlertContactGroupRead, resourceAlibabacloudStackArmsAlertContactGroupUpdate, resourceAlibabacloudStackArmsAlertContactGroupDelete)
	return resource
}

func resourceAlibabacloudStackArmsAlertContactGroupCreate(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var response map[string]interface{}
	action := "CreateAlertContactGroup"
	request := make(map[string]interface{})
	request["ContactGroupName"] = d.Get("alert_contact_group_name")
	if v, ok := d.GetOk("contact_ids"); ok {
		request["ContactIds"] = convertArrayToString(v.(*schema.Set).List(), " ")
	}
	response, err = client.DoTeaRequest("POST", "ARMS", "2019-08-08", action, "", nil, nil, request)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprint(response["ContactGroupId"]))

	return nil
}

func resourceAlibabacloudStackArmsAlertContactGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	armsService := ArmsService{client}
	object, err := armsService.DescribeArmsAlertContactGroup(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_arms_alert_contact_group armsService.DescribeArmsAlertContactGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	d.Set("alert_contact_group_name", object["ContactGroupName"])
	contactIdsItems := make([]string, 0)
	if contacts, ok := object["Contacts"]; ok && contacts != nil {
		for _, contactsItem := range contacts.([]interface{}) {
			if contactId, ok := contactsItem.(map[string]interface{})["ContactId"]; ok && contactId != nil {
				contactIdsItems = append(contactIdsItems, fmt.Sprint(contactId))
			}
		}
	}
	d.Set("contact_ids", contactIdsItems)
	return nil
}

func resourceAlibabacloudStackArmsAlertContactGroupUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	update := false
	request := map[string]interface{}{
		"ContactGroupId": d.Id(),
	}
	if d.HasChange("alert_contact_group_name") {
		update = true
	}
	request["ContactGroupName"] = d.Get("alert_contact_group_name")
	if d.HasChange("contact_ids") {
		update = true
	}
	if v, ok := d.GetOk("contact_ids"); ok {
		request["ContactIds"] = convertArrayToString(v.(*schema.Set).List(), " ")
	}
	if update {
		action := "UpdateAlertContactGroup"
		_, err = client.DoTeaRequest("POST", "ARMS", "2019-08-08", action, "", nil, nil, request)
		if err != nil {
			return err
		}
	}
	return nil
}

func resourceAlibabacloudStackArmsAlertContactGroupDelete(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	action := "DeleteAlertContactGroup"
	request := map[string]interface{}{
		"ContactGroupId": d.Id(),
	}
	_, err = client.DoTeaRequest("POST", "ARMS", "2019-08-08", action, "", nil, nil, request)
	if err != nil {
		return err
	}
	return nil
}
