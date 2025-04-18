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

func dataSourceAlibabacloudStackRouterInterfaces() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackRouterInterfacesRead,

		Schema: map[string]*schema.Schema{
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Active", "Inactive", "Idle"}, false),
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				ForceNew:     true,
			},
			"specification": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"router_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"router_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{string(VRouter), string(VBR)}, false),
			},
			"role": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"opposite_interface_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"opposite_interface_owner_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"interfaces": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
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
						"role": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"specification": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"router_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"router_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"access_point_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"opposite_region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"opposite_interface_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"opposite_router_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"opposite_router_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"opposite_interface_owner_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check_source_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check_target_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackRouterInterfacesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := vpc.CreateDescribeRouterInterfacesRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)

	var filters []vpc.DescribeRouterInterfacesFilter
	for _, key := range []string{"status", "router_id", "router_type", "opposite_interface_id", "opposite_interface_owner_id"} {
		if v, ok := d.GetOk(key); ok && v.(string) != "" {
			value := []string{v.(string)}
			filters = append(filters, vpc.DescribeRouterInterfacesFilter{
				Key:   terraformToAPI(key),
				Value: &value,
			})
		}
	}

	request.Filter = &filters

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			idsMap[Trim(vv.(string))] = Trim(vv.(string))
		}
	}

	var allRouterInterfaces []vpc.RouterInterfaceType
	invoker := NewInvoker()

	for {
		var raw interface{}
		var err error
		err = invoker.Run(func() error {
			raw, err = client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
				return vpcClient.DescribeRouterInterfaces(request)
			})
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		})
		response, ok := raw.(*vpc.DescribeRouterInterfacesResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_router_interfaces", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		if len(response.RouterInterfaceSet.RouterInterfaceType) < 1 {
			break
		}

		allRouterInterfaces = append(allRouterInterfaces, response.RouterInterfaceSet.RouterInterfaceType...)

		if len(response.RouterInterfaceSet.RouterInterfaceType) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		request.PageNumber = page
	}

	var filteredRouterInterfaces []vpc.RouterInterfaceType
	var r *regexp.Regexp
	if nameRegex, ok := d.GetOk("name_regex"); ok && nameRegex.(string) != "" {
		r = regexp.MustCompile(nameRegex.(string))
	}

	for _, v := range allRouterInterfaces {
		if len(idsMap) > 0 {
			if _, ok := idsMap[v.RouterInterfaceId]; !ok {
				continue
			}
		}
		if r != nil && !r.MatchString(v.Name) {
			continue
		}
		if role := d.Get("role").(string); role != "" && role != v.Role {
			continue
		}
		if spec := d.Get("specification").(string); spec != "" && spec != v.Spec {
			continue
		}
		filteredRouterInterfaces = append(filteredRouterInterfaces, v)
	}

	return riDecriptionAttributes(d, filteredRouterInterfaces, meta)
}

func riDecriptionAttributes(d *schema.ResourceData, riSetTypes []vpc.RouterInterfaceType, meta interface{}) error {
	var ids []string
	var names []string
	var s []map[string]interface{}
	for _, ri := range riSetTypes {
		mapping := map[string]interface{}{
			"id":                          ri.RouterInterfaceId,
			"status":                      ri.Status,
			"name":                        ri.Name,
			"description":                 ri.Description,
			"role":                        ri.Role,
			"specification":               ri.Spec,
			"router_id":                   ri.RouterId,
			"router_type":                 ri.RouterType,
			"vpc_id":                      ri.VpcInstanceId,
			"access_point_id":             ri.AccessPointId,
			"creation_time":               ri.CreationTime,
			"opposite_region_id":          ri.OppositeRegionId,
			"opposite_interface_id":       ri.OppositeInterfaceId,
			"opposite_router_id":          ri.OppositeRouterId,
			"opposite_router_type":        ri.OppositeRouterType,
			"opposite_interface_owner_id": ri.OppositeInterfaceOwnerId,
			"health_check_source_ip":      ri.HealthCheckSourceIp,
			"health_check_target_ip":      ri.HealthCheckTargetIp,
		}
		ids = append(ids, ri.RouterInterfaceId)
		names = append(names, ri.Name)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return errmsgs.WrapError(err)
	}

	if err := d.Set("interfaces", s); err != nil {
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
