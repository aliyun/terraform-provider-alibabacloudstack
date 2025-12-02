package alibabacloudstack

import (
	"fmt"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackDnsGtmAddressPool() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"A", "AAAA", "CNAME"}, false),
			},
			"lba_strategy": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"ALL_RR", "RATIO"}, false),
			},
			"addrs": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"mode": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"SMART", "ONLINE", "OFFLINE"}, false),
						},
						"lba_weight": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 100),
						},
					},
				},
			},
		},
	}

	setResourceFunc(resource, resourceAlibabacloudStackDnsGtmAddressPoolCreate, resourceAlibabacloudStackDnsGtmAddressPoolRead, resourceAlibabacloudStackDnsGtmAddressPoolUpdate, resourceAlibabacloudStackDnsGtmAddressPoolDelete)
	return resource
}

func resourceAlibabacloudStackDnsGtmAddressPoolCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := make(map[string]interface{})
	request["Name"] = d.Get("name")
	request["Type"] = d.Get("type")
	request["LbaStrategy"] = d.Get("lba_strategy")

	addrs := d.Get("addrs").([]interface{})
	addrMaps := make([]map[string]interface{}, 0)
	for _, v := range addrs {
		addr := v.(map[string]interface{})
		addrMap := map[string]interface{}{
			"Value": addr["value"],
			"Mode":  addr["mode"],
		}
		if v, ok := addr["lba_weight"]; ok {
			addrMap["LbaWeight"] = v
		}
		addrMaps = append(addrMaps, addrMap)
	}
	request["Addr"] = addrMaps

	resp, err := client.DoTeaRequest("POST", "CloudDns", "2021-06-24", "AddDnsGtmAddressPool", "", nil, nil, request)
	if err != nil {
		return fmt.Errorf("creating DnsGtmAddressPool failed: %v", err)
	}

	id, ok := resp["Id"].(string)
	if !ok || id == "" {
		return fmt.Errorf("failed to get Id from response")
	}

	d.SetId(id)

	return nil
}

func resourceAlibabacloudStackDnsGtmAddressPoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	dnsService := DnsService{client}

	object, err := dnsService.DescribeDnsGtmAddressPool(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", object["Name"])
	d.Set("type", object["Type"])
	d.Set("lba_strategy", object["LbaStrategy"])

	if addrs, ok := object["Addrs"].([]interface{}); ok {
		addrList := make([]map[string]interface{}, 0, len(addrs))
		for _, addr := range addrs {
			if addrMap, ok := addr.(map[string]interface{}); ok {
				addrList = append(addrList, map[string]interface{}{
					"value":      addrMap["Value"],
					"mode":       addrMap["Mode"],
					"lba_weight": addrMap["LbaWeight"],
				})
			}
		}
		d.Set("addrs", addrList)
	}

	return nil
}

func resourceAlibabacloudStackDnsGtmAddressPoolUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	if d.IsNewResource() {
		return resourceAlibabacloudStackDnsGtmAddressPoolRead(d, meta)
	}

	if d.HasChanges("name", "type", "lba_strategy", "addrs") {
		request := make(map[string]interface{})
		request["Id"] = d.Id()
		request["Name"] = d.Get("name")
		request["Type"] = d.Get("type")
		request["LbaStrategy"] = d.Get("lba_strategy")

		addrs := d.Get("addrs").([]interface{})
		addrMaps := make([]map[string]interface{}, 0)
		for _, v := range addrs {
			addr := v.(map[string]interface{})
			addrMap := map[string]interface{}{
				"Value": addr["value"],
				"Mode":  addr["mode"],
			}
			if v, ok := addr["lba_weight"]; ok {
				addrMap["LbaWeight"] = v
			}
			addrMaps = append(addrMaps, addrMap)
		}
		request["Addr"] = addrMaps

		_, err := client.DoTeaRequest("POST", "CloudDns", "2021-06-24", "UpdateDnsGtmAddressPool", "", nil, request, nil)
		if err != nil {
			return fmt.Errorf("updating DnsGtmAddressPool failed: %v", err)
		}
	}

	return nil
}

func resourceAlibabacloudStackDnsGtmAddressPoolDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	reqQuery := map[string]interface{}{
		"Id": d.Id(),
	}

	raw, err := client.DoTeaRequest("POST", "CloudDns", "2021-06-24", "DeleteDnsGtmAddressPool", "", nil, reqQuery, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "DeleteDnsGtmAddressPool", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	if !raw["success"].(bool) {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "DeleteDnsGtmAddressPool", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return nil
}
