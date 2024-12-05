package alibabacloudstack

import (
	"regexp"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackVpnCustomerGateways() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackVpnCgwsRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
				MinItems: 1,
			},

			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				ForceNew:     true,
			},

			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"gateways": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackVpnCgwsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := vpc.CreateDescribeCustomerGatewaysRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	var allCgws []vpc.CustomerGateway

	for {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeCustomerGateways(request)
		})
		response, ok := raw.(*vpc.DescribeCustomerGatewaysResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_vpn_customer_gateways", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		if len(response.CustomerGateways.CustomerGateway) < 1 {
			break
		}
		allCgws = append(allCgws, response.CustomerGateways.CustomerGateway...)
		if len(response.CustomerGateways.CustomerGateway) < PageSizeLarge {
			break
		}
		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		request.PageNumber = page
	}

	var filteredCgws []vpc.CustomerGateway
	var reg *regexp.Regexp
	var ids []string
	if v, ok := d.GetOk("ids"); ok && len(v.([]interface{})) > 0 {
		for _, item := range v.([]interface{}) {
			if item == nil {
				continue
			}
			ids = append(ids, strings.Trim(item.(string), " "))
		}
	}

	if nameRegex, ok := d.GetOk("name_regex"); ok && nameRegex.(string) != "" {
		if r, err := regexp.Compile(nameRegex.(string)); err == nil {
			reg = r
		} else {
			return errmsgs.WrapError(err)
		}
	}

	for _, cgw := range allCgws {
		if reg != nil {
			if !reg.MatchString(cgw.Name) {
				continue
			}
		}

		if ids != nil && len(ids) != 0 {
			for _, id := range ids {
				if cgw.CustomerGatewayId == id {
					filteredCgws = append(filteredCgws, cgw)
				}
			}
		} else {
			filteredCgws = append(filteredCgws, cgw)
		}
	}

	return vpnCgwsDecriptionAttributes(d, filteredCgws, meta)
}

func vpnCgwsDecriptionAttributes(d *schema.ResourceData, vpnSetTypes []vpc.CustomerGateway, meta interface{}) error {
	var ids []string
	var names []string
	var s []map[string]interface{}
	for _, vpn := range vpnSetTypes {
		mapping := map[string]interface{}{
			"id":           vpn.CustomerGatewayId,
			"name":         vpn.Name,
			"ip_address":   vpn.IpAddress,
			"description":  vpn.Description,
			"create_time":  TimestampToStr(vpn.CreateTime),
		}
		ids = append(ids, vpn.CustomerGatewayId)
		names = append(names, vpn.Name)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("gateways", s); err != nil {
		return errmsgs.WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return errmsgs.WrapError(err)
	}

	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
