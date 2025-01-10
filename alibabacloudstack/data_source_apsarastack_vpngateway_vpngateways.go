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

func dataSourceAlibabacloudStackVpnGateways() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackVpnsRead,

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

			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				ForceNew:     true,
			},

			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Init", "Provisioning", "Active", "Updating", "Deleting"}, false),
			},

			"business_status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Normal", "FinancialLocked"}, false),
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
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internet_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"specification": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"business_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enable_ipsec": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enable_ssl": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ssl_connections": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackVpnsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := vpc.CreateDescribeVpnGatewaysRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)

	var allVpns []vpc.VpnGateway

	if v, ok := d.GetOk("vpc_id"); ok && v.(string) != "" {
		request.VpcId = v.(string)
	}

	if v, ok := d.GetOk("status"); ok && v.(string) != "" {
		request.Status = strings.ToLower(v.(string))
	}

	if v, ok := d.GetOk("business_status"); ok && v.(string) != "" {
		request.BusinessStatus = v.(string)
	}

	for {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeVpnGateways(request)
		})
		response, ok := raw.(*vpc.DescribeVpnGatewaysResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_vpn_gateways", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)

		if len(response.VpnGateways.VpnGateway) < 1 {
			break
		}

		allVpns = append(allVpns, response.VpnGateways.VpnGateway...)

		if len(response.VpnGateways.VpnGateway) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return errmsgs.WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	var filteredVpns []vpc.VpnGateway
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

	for _, vpn := range allVpns {
		if reg != nil {
			if !reg.MatchString(vpn.Name) {
				continue
			}
		}
		if ids != nil && len(ids) != 0 {
			for _, id := range ids {
				if vpn.VpnGatewayId == id {
					filteredVpns = append(filteredVpns, vpn)
				}
			}
		} else {
			filteredVpns = append(filteredVpns, vpn)
		}

	}

	return vpnsDecriptionAttributes(d, filteredVpns, meta)
}

func convertStatus(lower string) string {
	upStr := strings.ToUpper(lower)

	wholeStr := string(upStr[0]) + lower[1:]
	return wholeStr
}

func convertChargeType(originType string) string {
	if string("PostpayByFlow") == originType {
		return string(PostPaid)
	} else {
		return string(PrePaid)
	}
}

func vpnsDecriptionAttributes(d *schema.ResourceData, vpnSetTypes []vpc.VpnGateway, meta interface{}) error {
	var ids []string
	var names []string
	var s []map[string]interface{}
	for _, vpn := range vpnSetTypes {
		mapping := map[string]interface{}{
			"id":                vpn.VpnGatewayId,
			"vpc_id":            vpn.VpcId,
			"internet_ip":       vpn.InternetIp,
			"create_time":       TimestampToStr(vpn.CreateTime),
			"end_time":          TimestampToStr(vpn.EndTime),
			"specification":     vpn.Spec,
			"name":              vpn.Name,
			"description":       vpn.Description,
			"status":            convertStatus(vpn.Status),
			"business_status":   vpn.BusinessStatus,
			"instance_charge_type": convertChargeType(vpn.ChargeType),
			"enable_ipsec":      vpn.IpsecVpn,
			"enable_ssl":        vpn.SslVpn,
			"ssl_connections":   vpn.SslMaxConnections,
		}

		ids = append(ids, vpn.VpnGatewayId)
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
		if err := writeToFile(output.(string), s); err != nil {
			return err
		}
	}
	return nil
}
