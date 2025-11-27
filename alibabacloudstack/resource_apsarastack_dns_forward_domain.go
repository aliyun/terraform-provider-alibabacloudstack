package alibabacloudstack

import (
	"encoding/json"
	"fmt"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackDnsForwardDomain() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"forward_mode": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"FORWARD_FIRST", "FORWARD_ONLY"}, false),
			},
			"forwarders": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"remark": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"caller_uid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_timestamp": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"update_timestamp": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}

	setResourceFunc(resource, resourceAlibabacloudStackDnsForwardDomainCreate, resourceAlibabacloudStackDnsForwardDomainRead, resourceAlibabacloudStackDnsForwardDomainUpdate, resourceAlibabacloudStackDnsForwardDomainDelete)
	return resource
}

func resourceAlibabacloudStackDnsForwardDomainCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := make(map[string]interface{})
	request["Name"] = d.Get("name")
	request["ForwardMode"] = d.Get("forward_mode")

	forwarderlist := d.Get("forwarders").(*schema.Set).List()
	forwarders, _ := json.Marshal(forwarderlist)
	request["Forwarders"] = string(forwarders)

	resp, err := client.DoTeaRequest("POST", "CloudDns", "2021-06-24", "AddGlobalForwardZone", "", nil, nil, request)
	if err != nil {
		return err
	}

	resourceId, ok := resp["Id"].(string)
	if !ok || resourceId == "" {
		return fmt.Errorf("failed to get resource id from response")
	}

	d.SetId(resourceId)

	return nil
}

func resourceAlibabacloudStackDnsForwardDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	dnsService := DnsService{client}
	object, err := dnsService.DescribeDnsForwardDomain(d.Id())

	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("name", object["Name"])
	d.Set("forward_mode", object["ForwardMode"])
	d.Set("forwarders", object["Forwarders"])
	d.Set("remark", object["Remark"])
	d.Set("caller_uid", object["CallerUid"])
	d.Set("create_timestamp", object["CreateTimestamp"])
	d.Set("update_timestamp", object["UpdateTimestamp"])

	return nil
}

func resourceAlibabacloudStackDnsForwardDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	if d.HasChange("remark") {
		request := make(map[string]interface{})
		request["Id"] = d.Id()
		request["Name"] = d.Get("name")
		request["Remark"] = d.Get("remark")

		_, err := client.DoTeaRequest("POST", "CloudDns", "2021-06-24", "UpdateGlobalForwardZoneRemark", "", nil, nil, request)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg,
				"alibabacloudstack_dns_forward_domain", "UpdateGlobalForwardZoneRemark", errmsgs.AlibabacloudStackSdkGoERROR)
		}
	}

	if d.IsNewResource() {
		return nil
	}

	if d.HasChanges("forward_mode", "forwarders") {
		request := make(map[string]interface{})
		request["Id"] = d.Id()
		request["Name"] = d.Get("name")
		request["ForwardMode"] = d.Get("forward_mode")

		forwarderlist := d.Get("forwarders").(*schema.Set).List()
		forwarders, _ := json.Marshal(forwarderlist)
		request["Forwarders"] = string(forwarders)

		_, err := client.DoTeaRequest("POST", "CloudDns", "2021-06-24", "UpdateGlobalForwardZone", "", nil, nil, request)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg,
				"alibabacloudstack_dns_forward_domain", "UpdateGlobalForwardZone", errmsgs.AlibabacloudStackSdkGoERROR)
		}
	}
	return nil
}

func resourceAlibabacloudStackDnsForwardDomainDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	reqBody := map[string]interface{}{
		"Id": d.Id(),
	}

	_, err := client.DoTeaRequest("POST", "CloudDns", "2021-06-24", "DeleteGlobalForwardZone", "", nil, nil, reqBody)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "DeleteGlobalForwardZone", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return nil
}
