package alibabacloudstack

import (
	"encoding/json"
	"fmt"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackDnsRecursorAcl() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"policy": {
				Type:     schema.TypeString,
				Required: true,
			},
			"line_ids": {
				Type:     schema.TypeSet,
				MinItems: 1,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
			"remark": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}

	setResourceFunc(resource, resourceAlibabacloudStackDnsRecursorAclCreate, resourceAlibabacloudStackDnsRecursorAclRead, resourceAlibabacloudStackDnsRecursorAclUpdate, resourceAlibabacloudStackDnsRecursorAclDelete)
	return resource
}

func resourceAlibabacloudStackDnsRecursorAclCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := make(map[string]interface{})
	request["Name"] = d.Get("name")
	request["Policy"] = d.Get("policy")
	lineIds := d.Get("line_ids").(*schema.Set).List()
	line_ids, _ := json.Marshal(lineIds)
	request["LineIds"] = string(line_ids)

	if remark, ok := d.GetOk("remark"); ok {
		request["Remark"] = remark
	}

	response, err := client.DoTeaRequest("POST", "CloudDns", "2022-06-24", "AddRecursorAcl", "", nil, nil, request)
	if err != nil {
		return err
	}

	resourceId, ok := response["Id"].(string)
	if !ok || resourceId == "" {
		return fmt.Errorf("failed to get resource id from response")
	}

	d.SetId(resourceId)

	return nil
}

func resourceAlibabacloudStackDnsRecursorAclRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	dnsService := DnsService{client}

	object, err := dnsService.DescribeDnsRecursorAcl(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", object["Name"])
	d.Set("policy", object["Policy"])
	d.Set("remark", object["Remark"])
	d.Set("line_ids", object["LineIds"])

	return nil
}

func resourceAlibabacloudStackDnsRecursorAclUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	if d.IsNewResource() {
		return nil
	}

	if d.HasChanges("name", "policy", "line_ids") {
		request := make(map[string]interface{})
		request["Id"] = d.Id()
		request["Name"] = d.Get("name")
		request["Policy"] = d.Get("policy")
		lineIds := d.Get("line_ids").(*schema.Set).List()
		line_ids, _ := json.Marshal(lineIds)
		request["LineIds"] = string(line_ids)

		_, err := client.DoTeaRequest("POST", "CloudDns", "2022-06-24", "UpdateRecursorAcl", "", nil, nil, request)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg,
				"alibabacloudstack_universal_dns_recursor_acl", "UpdateRecursorAcl", errmsgs.AlibabacloudStackSdkGoERROR)
		}
	}
	if d.HasChange("remark") {
		request := make(map[string]interface{})
		request["Id"] = d.Id()
		request["Remark"] = d.Get("remark")

		_, err := client.DoTeaRequest("POST", "CloudDns", "2022-06-24", "UpdateRecursorAclRemark", "", nil, nil, request)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg,
				"alibabacloudstack_universal_dns_recursor_acl", "UpdateRecursorAclRemark", errmsgs.AlibabacloudStackSdkGoERROR)
		}
	}

	return nil
}

func resourceAlibabacloudStackDnsRecursorAclDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	reqBody := map[string]interface{}{
		"Id": d.Id(),
	}

	_, err := client.DoTeaRequest("POST", "CloudDns", "2022-06-24", "DeleteRecursorAcl", "", nil, nil, reqBody)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "DeleteRecursorAcl", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return nil
}
