package alibabacloudstack

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackCommonBandwidthPackages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackCommonBandwidthPackagesRead,

		Schema: map[string]*schema.Schema{
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
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			// Computed values
			"packages": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bandwidth": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
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
						"business_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"isp": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_ip_addresses": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip_address": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"allocation_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
							MinItems: 0,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackCommonBandwidthPackagesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := vpc.CreateDescribeCommonBandwidthPackagesRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	request.ResourceGroupId = d.Get("resource_group_id").(string)
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			idsMap[Trim(vv.(string))] = Trim(vv.(string))
		}
	}

	var allCommonBandwidthPackages []vpc.CommonBandwidthPackage
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		if r, err := regexp.Compile(Trim(v.(string))); err == nil {
			nameRegex = r
		} else {
			errmsgs.WrapError(err)
		}
	}
	invoker := NewInvoker()
	for {
		var raw interface{}
		var err error
		err = invoker.Run(func() error {
			raw, err = client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
				return vpcClient.DescribeCommonBandwidthPackages(request)
			})
			return err
		})
		response, ok := raw.(*vpc.DescribeCommonBandwidthPackagesResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_common_bandwidth_packages", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		if !ok {
			return errmsgs.WrapErrorf(nil, errmsgs.DataDefaultErrorMsg, "alibabacloudstack_common_bandwidth_packages", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, "Failed to cast response to DescribeCommonBandwidthPackagesResponse")
		}
		if len(response.CommonBandwidthPackages.CommonBandwidthPackage) < 1 {
			break
		}

		for _, cbwp := range response.CommonBandwidthPackages.CommonBandwidthPackage {
			if nameRegex != nil {
				if !nameRegex.MatchString(cbwp.Name) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[cbwp.BandwidthPackageId]; !ok {
					continue
				}
			}
			allCommonBandwidthPackages = append(allCommonBandwidthPackages, cbwp)
		}

		if len(response.CommonBandwidthPackages.CommonBandwidthPackage) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return errmsgs.WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	return CommonBandwidthPackagesDecriptionAttributes(d, allCommonBandwidthPackages, meta)
}

func CommonBandwidthPackagesDecriptionAttributes(d *schema.ResourceData, cbwps []vpc.CommonBandwidthPackage, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}
	var ids []string
	var names []string
	var s []map[string]interface{}
	for _, cbwp := range cbwps {
		mapping := map[string]interface{}{
			"id":                  cbwp.BandwidthPackageId,
			"bandwidth":           cbwp.Bandwidth,
			"description":         cbwp.Description,
			"status":              cbwp.Status,
			"business_status":     cbwp.BusinessStatus,
			"isp":                 cbwp.ISP,
			"name":                cbwp.Name,
			"creation_time":       cbwp.CreationTime,
			"public_ip_addresses": vpcService.FlattenPublicIpAddressesMappings(cbwp.PublicIpAddresses.PublicIpAddresse),
		}
		names = append(names, cbwp.Name)
		ids = append(ids, cbwp.BandwidthPackageId)
		s = append(s, mapping)
	}
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("packages", s); err != nil {
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
