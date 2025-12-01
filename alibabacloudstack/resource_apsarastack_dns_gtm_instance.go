package alibabacloudstack

import (
	"fmt"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackDnsGtmInstance() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ttl": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"prefix": {
				Type:     schema.TypeString,
				Required: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zone_name": {
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

	setResourceFunc(resource, resourceAlibabacloudStackDnsGtmInstanceCreate, resourceAlibabacloudStackDnsGtmInstanceRead, resourceAlibabacloudStackDnsGtmInstanceUpdate, resourceAlibabacloudStackDnsGtmInstanceDelete)
	return resource
}

func resourceAlibabacloudStackDnsGtmInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := make(map[string]interface{})
	request["Name"] = d.Get("name")
	request["ZoneId"] = d.Get("zone_id")
	request["Ttl"] = d.Get("ttl")
	request["Prefix"] = d.Get("prefix")

	resp, err := client.DoTeaRequest("POST", "CloudDns", "2021-06-24", "AddDnsGtmInstance", "", nil, nil, request)
	if err != nil {
		return err
	}

	id, ok := resp["Id"].(string)
	if !ok || id == "" {
		return fmt.Errorf("failed to get instance id from response")
	}

	d.SetId(id)

	return nil
}

func resourceAlibabacloudStackDnsGtmInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	dnsService := DnsService{client}
	object, err := dnsService.DescribeDnsGtmInstance(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", object["Name"])
	d.Set("zone_id", object["ZoneId"])
	d.Set("ttl", object["Ttl"])
	d.Set("prefix", object["Prefix"])
	d.Set("id", object["Id"])
	d.Set("zone_name", object["ZoneName"])
	d.Set("create_timestamp", object["CreateTimestamp"])
	d.Set("update_timestamp", object["UpdateTimestamp"])

	return nil
}

func resourceAlibabacloudStackDnsGtmInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	if d.IsNewResource() {
		return nil
	}

	if d.HasChanges("name", "ttl", "prefix") {
		reqBody := map[string]interface{}{
			"Id":     d.Id(),
			"Name":   d.Get("name"),
			"ZoneId": d.Get("zone_id"),
			"Ttl":    d.Get("ttl"),
			"Prefix": d.Get("prefix"),
		}

		_, err := client.DoTeaRequest("POST", "CloudDns", "2021-06-24", "UpdateDnsGtmInstance", "", nil, nil, reqBody)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg,
				"alibabacloudstack_dns_gtm_instance", "UpdateDnsGtmInstance", errmsgs.AlibabacloudStackSdkGoERROR)
		}
	}

	return nil
}

func resourceAlibabacloudStackDnsGtmInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	reqBody := map[string]interface{}{
		"Id": d.Id(),
	}

	_, err := client.DoTeaRequest("POST", "CloudDns", "2021-06-24", "DeleteDnsGtmInstance", "", nil, nil, reqBody)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "DeleteDnsGtmInstance", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return nil
}
