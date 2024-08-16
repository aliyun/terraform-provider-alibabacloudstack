package alibabacloudstack

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackVpcIpv6InternetBandwidth() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackVpcIpv6InternetBandwidthCreate,
		Read:   resourceAlibabacloudStackVpcIpv6InternetBandwidthRead,
		Update: resourceAlibabacloudStackVpcIpv6InternetBandwidthUpdate,
		Delete: resourceAlibabacloudStackVpcIpv6InternetBandwidthDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"bandwidth": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 5000),
			},
			"internet_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PayByBandwidth", "PayByTraffic"}, false),
			},
			"ipv6_address_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ipv6_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlibabacloudStackVpcIpv6InternetBandwidthCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var response map[string]interface{}
	action := "AllocateIpv6InternetBandwidth"
	request := make(map[string]interface{})
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	request["Bandwidth"] = d.Get("bandwidth")
	if v, ok := d.GetOk("internet_charge_type"); ok {
		request["InternetChargeType"] = v
	}
	request["Ipv6AddressId"] = d.Get("ipv6_address_id")
	request["Ipv6GatewayId"] = d.Get("ipv6_gateway_id")
	request["RegionId"] = client.RegionId
	runtime := util.RuntimeOptions{IgnoreSSL: tea.Bool(client.Config.Insecure)}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken("AllocateIpv6InternetBandwidth")
		request["Product"] = "Vpc"
		request["OrganizationId"] = client.Department
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_vpc_ipv6_internet_bandwidth", action, AlibabacloudStackSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["InternetBandwidthId"]))

	return resourceAlibabacloudStackVpcIpv6InternetBandwidthRead(d, meta)
}
func resourceAlibabacloudStackVpcIpv6InternetBandwidthRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}
	object, err := vpcService.DescribeVpcIpv6InternetBandwidth(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_vpc_ipv6_internet_bandwidth vpcService.DescribeVpcIpv6InternetBandwidth Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("ipv6_address_id", object["Ipv6AddressId"])
	d.Set("ipv6_gateway_id", object["Ipv6GatewayId"])
	if ipv6InternetBandwidth, ok := object["Ipv6InternetBandwidth"]; ok {
		if v, ok := ipv6InternetBandwidth.(map[string]interface{}); ok {
			d.Set("bandwidth", formatInt(v["Bandwidth"]))
			d.Set("internet_charge_type", v["InternetChargeType"])
			d.Set("status", v["BusinessStatus"])
		}
	}
	return nil
}
func resourceAlibabacloudStackVpcIpv6InternetBandwidthUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var response map[string]interface{}
	request := map[string]interface{}{
		"Ipv6InternetBandwidthId": d.Id(),
	}
	if d.HasChange("bandwidth") {
		request["Bandwidth"] = d.Get("bandwidth")
	}
	request["RegionId"] = client.RegionId
	action := "ModifyIpv6InternetBandwidth"
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	runtime := util.RuntimeOptions{IgnoreSSL: tea.Bool(client.Config.Insecure)}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken("ModifyIpv6InternetBandwidth")
		request["Product"] = "Vpc"
		request["OrganizationId"] = client.Department
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabacloudStackSdkGoERROR)
	}
	return resourceAlibabacloudStackVpcIpv6InternetBandwidthRead(d, meta)
}
func resourceAlibabacloudStackVpcIpv6InternetBandwidthDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	action := "DeleteIpv6InternetBandwidth"
	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"Ipv6InternetBandwidthId": d.Id(),
	}

	request["Ipv6AddressId"] = d.Get("ipv6_address_id")
	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		request["Product"] = "Vpc"
		request["OrganizationId"] = client.Department
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{IgnoreSSL: tea.Bool(client.Config.Insecure)})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabacloudStackSdkGoERROR)
	}
	return nil
}
