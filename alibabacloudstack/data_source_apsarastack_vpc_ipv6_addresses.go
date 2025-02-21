package alibabacloudstack

import (
	"fmt"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackVpcIpv6Addresses() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackVpcIpv6AddressesRead,
		Schema: map[string]*schema.Schema{
			"associated_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Available", "Pending"}, false),
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"associated_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"associated_instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipv6_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipv6_address_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipv6_address_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipv6_gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackVpcIpv6AddressesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	action := "DescribeIpv6Addresses"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("associated_instance_id"); ok {
		request["AssociatedInstanceId"] = v
	}
	if v, ok := d.GetOk("vswitch_id"); ok {
		request["VSwitchId"] = v
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		request["VpcId"] = v
	}
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}

	status, statusOk := d.GetOk("status")
	for {
		response, err := client.DoTeaRequest("POST", "Vpc", "2016-04-28", action, "", nil, nil, request)
		if err != nil {
			return err
		}
		resp, err := jsonpath.Get("$.Ipv6Addresses.Ipv6Address", response)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, action, "$.Ipv6Addresses.Ipv6Address", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if statusOk && status.(string) != "" && status.(string) != item["Status"].(string) {
				continue
			}
			objects = append(objects, item)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"associated_instance_id":   object["AssociatedInstanceId"],
			"associated_instance_type": object["AssociatedInstanceType"],
			"ipv6_address":             object["Ipv6Address"],
			"id":                       fmt.Sprint(object["Ipv6AddressId"]),
			"ipv6_address_id":          fmt.Sprint(object["Ipv6AddressId"]),
			"ipv6_address_name":        object["Ipv6AddressName"],
			"ipv6_gateway_id":          object["Ipv6GatewayId"],
			"network_type":             object["NetworkType"],
			"status":                   object["Status"],
			"vswitch_id":               object["VSwitchId"],
			"vpc_id":                   object["VpcId"],
			"create_time":              object["AllocationTime"],
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["Ipv6AddressName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return errmsgs.WrapError(err)
	}

	if err := d.Set("addresses", s); err != nil {
		return errmsgs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), s); err != nil {
			return err
		}
	}

	return nil
}
