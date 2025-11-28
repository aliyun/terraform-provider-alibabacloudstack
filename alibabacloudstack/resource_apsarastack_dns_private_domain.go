package alibabacloudstack

import (
	"fmt"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackDnsPrivateDomain() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"remark": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"caller_uid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"record_count": {
				Type:     schema.TypeInt,
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
	setResourceFunc(resource, resourceAlibabacloudStackDnsPrivateDomainCreate, resourceAlibabacloudStackDnsPrivateDomainRead, resourceAlibabacloudStackDnsPrivateDomainUpdate, resourceAlibabacloudStackDnsPrivateDomainDelete)
	return resource
}

func resourceAlibabacloudStackDnsPrivateDomainCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := make(map[string]interface{})
	request["Name"] = d.Get("name")

	action := "AddPrivateZone"
	response, err := client.DoTeaRequest("POST", "CloudDns", "2021-06-24", action, "", nil, nil, request)
	if err != nil {
		return err
	}
	id, ok := response["Id"].(string)
	if !ok || id == "" {
		return fmt.Errorf("unable to find 'Id' in response for %s", action)
	}
	d.SetId(id)
	return nil
}

func resourceAlibabacloudStackDnsPrivateDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	dnsService := DnsService{client}
	object, err := dnsService.DescribePrivateZone(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", object["Name"])
	d.Set("remark", object["Remark"])
	d.Set("caller_uid", object["CallerUid"])
	d.Set("record_count", object["RecordCount"])
	d.Set("create_timestamp", object["CreateTimestamp"])
	d.Set("update_timestamp", object["UpdateTimestamp"])

	if vpcs, ok := object["RegionAndVpcs"].([]interface{}); ok && len(vpcs) > 0 {
		var vpcList []string
		for _, regionVpc := range vpcs {
			rv := regionVpc.(map[string]interface{})
			if rvpcs, ok := rv["Vpcs"].([]interface{}); ok {
				for _, vpc := range rvpcs {
					v := vpc.(map[string]interface{})
					vpcList = append(vpcList, v["Id"].(string))
				}
			}
		}
		d.Set("vpc_ids", vpcList)
	}

	return nil
}

func resourceAlibabacloudStackDnsPrivateDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	if d.HasChange("remark") {
		request := map[string]interface{}{
			"Id":     d.Id(),
			"Name":   d.Get("name"),
			"Remark": d.Get("remark"),
		}
		_, err := client.DoTeaRequest("POST", "CloudDns", "2021-06-24", "UpdatePrivateZoneRemark", "", nil, nil, request)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_dns_private_domain", "UpdatePrivateZoneRemark", errmsgs.AlibabacloudStackSdkGoERROR)
		}
	}

	// Update VPCs if changed
	if d.HasChange("vpc_ids") {
		vpcservice := VpcService{client}
		request := make(map[string]interface{})
		request["Id"] = d.Id()
		request["Name"] = d.Get("name")

		vpcs := d.Get("vpc_ids").(*schema.Set).List()
		vpcMap := make([]map[string]interface{}, 0)
		for _, v := range vpcs {
			vpc, err := vpcservice.DescribeVpc(v.(string))
			if err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_dns_private_domain", "DescribeVpc", errmsgs.AlibabacloudStackSdkGoERROR)
			}
			vpcMap = append(vpcMap, map[string]interface{}{
				"Id":   vpc.VpcId,
				"Name": vpc.VpcName,
			})
		}
		request["Vpcs"] = vpcMap
		_, err := client.DoTeaRequest("POST", "CloudDns", "2021-06-24", "BindZoneVpc", "", nil, nil, request)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_dns_private_domain", "BindZoneVpc", errmsgs.AlibabacloudStackSdkGoERROR)
		}
	}

	return nil
}

func resourceAlibabacloudStackDnsPrivateDomainDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	reqQuery := map[string]interface{}{
		"Id": d.Id(),
	}

	raw, err := client.DoTeaRequest("POST", "CloudDns", "2021-06-24", "DeletePrivateZone", "", nil, nil, reqQuery)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "DeletePrivateZone", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	if !raw["success"].(bool) {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "DeletePrivateZone", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return nil
}
