package alibabacloudstack

import (
	"encoding/json"
	"fmt"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackDnsPrivateLine() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"priority": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"v4_addresses": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"v6_addresses": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackDnsPrivateLineCreate, resourceAlibabacloudStackDnsPrivateLineRead, resourceAlibabacloudStackDnsPrivateLineUpdate, resourceAlibabacloudStackDnsPrivateLineDelete)
	return resource
}

func resourceAlibabacloudStackDnsPrivateLineCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := make(map[string]interface{})
	request["Name"] = d.Get("name")
	request["Priority"] = "0"
	v4Addresses := d.Get("v4_addresses").(*schema.Set).List()
	v6Addresses := d.Get("v6_addresses").(*schema.Set).List()
	if len(v4Addresses) == 0 && len(v6Addresses) == 0 {
		return fmt.Errorf("only one of v4_addresses and v6_addresses can be set")
	}
	if len(v4Addresses) > 0 {
		stringV4Addresses, _ := json.Marshal(v4Addresses)
		request["V4Addresses"] = string(stringV4Addresses)
	}
	if len(v6Addresses) > 0 {
		stringV6Addresses, _ := json.Marshal(v6Addresses)
		request["V6Addresses"] = string(stringV6Addresses)
	}
	response, err := client.DoTeaRequest("POST", "CloudDns", "2021-06-24", "AddPrivateLine", "", nil, nil, request)
	if err != nil {
		return err
	}

	resourceId, ok := response["Id"].(string)
	if !ok {
		return fmt.Errorf("failed to get Id from response")
	}

	d.SetId(resourceId)

	return nil
}

func resourceAlibabacloudStackDnsPrivateLineRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	dnsservice := DnsService{client}
	object, err := dnsservice.DescribePrivateLine(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	d.Set("name", object["Name"])
	d.Set("v4_addresses", object["V4Addresses"])
	d.Set("v6_addresses", object["V6Addresses"])
	d.Set("priority", object["Priority"])

	return nil
}

func resourceAlibabacloudStackDnsPrivateLineUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// If it is a new resource, there is no need to update.
	if d.IsNewResource() {
		return nil
	}

	if d.HasChanges("name", "v4_addresses", "v6_addresses") {
		updateReq := make(map[string]interface{})
		updateReq["Id"] = d.Id()
		updateReq["Name"] = d.Get("name")
		updateReq["Priority"] = d.Get("priority")
		v4Addresses := d.Get("v4_addresses").(*schema.Set).List()
		v6Addresses := d.Get("v6_addresses").(*schema.Set).List()
		if len(v4Addresses) == 0 && len(v6Addresses) == 0 {
			return fmt.Errorf("only one of v4_addresses and v6_addresses can be set")
		}
		if len(v4Addresses) > 0 {
			stringV4Addresses, _ := json.Marshal(v4Addresses)
			updateReq["V4Addresses"] = string(stringV4Addresses)
		}
		if len(v6Addresses) > 0 {
			stringV6Addresses, _ := json.Marshal(v6Addresses)
			updateReq["V6Addresses"] = string(stringV6Addresses)
		}
		_, err := client.DoTeaRequest("POST", "CloudDns", "2021-06-24", "UpdatePrivateLine", "", nil, nil, updateReq)
		if err != nil {
			return fmt.Errorf("failed to update DNS Private Line: %v", err)
		}
	}

	return nil
}

func resourceAlibabacloudStackDnsPrivateLineDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	reqBody := map[string]interface{}{
		"Id": d.Id(),
	}

	raw, err := client.DoTeaRequest("POST", "CloudDns", "2021-06-24", "DeletePrivateLine", "", nil, nil, reqBody)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "DeletePrivateLine", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	if !raw["success"].(bool) {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "DeletePrivateLine", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return nil
}
