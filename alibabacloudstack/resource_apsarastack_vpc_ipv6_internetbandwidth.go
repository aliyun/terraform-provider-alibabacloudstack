package alibabacloudstack

import (
	"fmt"
	"log"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
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

func resourceAlibabacloudStackVpcIpv6InternetBandwidthCreate(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var response map[string]interface{}
	action := "AllocateIpv6InternetBandwidth"
	request := make(map[string]interface{})
	request["Bandwidth"] = d.Get("bandwidth")
	if v, ok := d.GetOk("internet_charge_type"); ok {
		request["InternetChargeType"] = v
	}
	request["Ipv6AddressId"] = d.Get("ipv6_address_id")
	request["Ipv6GatewayId"] = d.Get("ipv6_gateway_id")

	request["ClientToken"] = buildClientToken("AllocateIpv6InternetBandwidth")
	response, err = client.DoTeaRequest("POST", "Vpc", "2016-04-28", action, "", nil, nil, request)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprint(response["InternetBandwidthId"]))

	return resourceAlibabacloudStackVpcIpv6InternetBandwidthRead(d, meta)
}

func resourceAlibabacloudStackVpcIpv6InternetBandwidthRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}
	object, err := vpcService.DescribeVpcIpv6InternetBandwidth(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_vpc_ipv6_internet_bandwidth vpcService.DescribeVpcIpv6InternetBandwidth Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
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

func resourceAlibabacloudStackVpcIpv6InternetBandwidthUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := map[string]interface{}{
		"Ipv6InternetBandwidthId": d.Id(),
	}
	if d.HasChange("bandwidth") {
		request["Bandwidth"] = d.Get("bandwidth")
	}

	action := "ModifyIpv6InternetBandwidth"

	request["ClientToken"] = buildClientToken("ModifyIpv6InternetBandwidth")
	_, err = client.DoTeaRequest("POST", "Vpc", "2016-04-28", action, "", nil, nil, request)
	if err != nil {
		return err
	}
	return resourceAlibabacloudStackVpcIpv6InternetBandwidthRead(d, meta)
}

func resourceAlibabacloudStackVpcIpv6InternetBandwidthDelete(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	action := "DeleteIpv6InternetBandwidth"
	request := map[string]interface{}{
		"Ipv6InternetBandwidthId": d.Id(),
	}
	request["Ipv6AddressId"] = d.Get("ipv6_address_id")

	_, err = client.DoTeaRequest("POST", "Vpc", "2016-04-28", action, "", nil, nil, request)
	if err != nil {
		return err
	}
	return nil
}
