package apsarastack

import (
	"fmt"
	"regexp"
	"strings"

	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceApsaraStackVpcs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApsaraStackVpcsRead,

		Schema: map[string]*schema.Schema{
			"cidr_block": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Available", "Pending"}, false),
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"is_default": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
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
			"vpcs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cidr_block": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipv6_cidr_block": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_default": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"route_table_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"router_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vrouter_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"secondary_cidr_blocks": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"user_cidrs": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"vswitch_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"dhcp_options_set_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"tags": tagsSchema(),
			"vpc_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpc_owner_id": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}
func dataSourceApsaraStackVpcsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)

	request := vpc.CreateDescribeVpcsRequest()
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.RegionId = string(client.Region)
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)

	if v, ok := d.GetOk("dhcp_options_set_id"); ok {
		request.DhcpOptionsSetId = v.(string)
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request.DryRun = requests.NewBoolean(v.(bool))
	}

	if v, ok := d.GetOkExists("is_default"); ok {
		request.IsDefault = requests.NewBoolean(v.(bool))
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request.ResourceGroupId = v.(string)
	}

	if v, ok := d.GetOk("tags"); ok {
		var tags []vpc.DescribeVpcsTag

		for key, value := range v.(map[string]interface{}) {
			tags = append(tags, vpc.DescribeVpcsTag{
				Key:   key,
				Value: value.(string),
			})
		}
		request.Tag = &tags
	}

	if v, ok := d.GetOk("vpc_name"); ok {
		request.VpcName = v.(string)
	}

	if v, ok := d.GetOk("vpc_owner_id"); ok {
		request.VpcOwnerId = requests.NewInteger(v.(int))
	}

	var allVpcs []vpc.Vpc
	invoker := NewInvoker()
	for {
		var raw interface{}
		var err error
		err = invoker.Run(func() error {
			raw, err = client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
				return vpcClient.DescribeVpcs(request)
			})
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return err
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "apsarastack_vpcs", request.GetActionName(), ApsaraStackSdkGoERROR)
		}
		response, _ := raw.(*vpc.DescribeVpcsResponse)
		if len(response.Vpcs.Vpc) < 1 {
			break
		}

		allVpcs = append(allVpcs, response.Vpcs.Vpc...)

		if len(response.Vpcs.Vpc) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		request.PageNumber = page
	}

	var filteredVpcs []vpc.Vpc
	var route_tables []string
	var r *regexp.Regexp
	if nameRegex, ok := d.GetOk("name_regex"); ok && nameRegex.(string) != "" {
		r = regexp.MustCompile(nameRegex.(string))
	}

	// ids
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			idsMap[Trim(vv.(string))] = Trim(vv.(string))
		}
	}

	for _, v := range allVpcs {
		if r != nil && !r.MatchString(v.VpcName) {
			continue
		}

		if len(idsMap) > 0 {
			if _, ok := idsMap[v.VpcId]; !ok {
				continue
			}
		}

		if cidrBlock, ok := d.GetOk("cidr_block"); ok && v.CidrBlock != cidrBlock.(string) {
			continue
		}

		if status, ok := d.GetOk("status"); ok && string(v.Status) != status.(string) {
			continue
		}

		if isDefault, ok := d.GetOk("is_default"); ok && v.IsDefault != isDefault.(bool) {
			continue
		}

		if vswitchId, ok := d.GetOk("vswitch_id"); ok && !vpcVswitchIdListContains(v.VSwitchIds.VSwitchId, vswitchId.(string)) {
			continue
		}
		request := vpc.CreateDescribeVRoutersRequest()
		request.RegionId = client.RegionId
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
		request.VRouterId = v.VRouterId
		request.RegionId = string(client.Region)

		var response *vpc.DescribeVRoutersResponse
		wait := incrementalWait(1*time.Second, 1*time.Second)
		err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
				return vpcClient.DescribeVRouters(request)
			})
			if err != nil {
				if IsThrottling(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			response, _ = raw.(*vpc.DescribeVRoutersResponse)
			return nil
		})

		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "apsarastack_vpcs", request.GetActionName(), ApsaraStackSdkGoERROR)
		}

		if len(response.VRouters.VRouter) > 0 {
			route_tables = append(route_tables, response.VRouters.VRouter[0].RouteTableIds.RouteTableId[0])
		} else {
			route_tables = append(route_tables, "")
		}

		filteredVpcs = append(filteredVpcs, v)
	}

	return vpcsDecriptionAttributes(d, filteredVpcs, route_tables, meta)
}
func vpcVswitchIdListContains(vswitchIdList []string, vswitchId string) bool {
	for _, idListItem := range vswitchIdList {
		if idListItem == vswitchId {
			return true
		}
	}
	return false
}
func vpcsDecriptionAttributes(d *schema.ResourceData, vpcSetTypes []vpc.Vpc, route_tables []string, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	vpcService := VpcService{client}
	var ids []string
	var names []string
	var s []map[string]interface{}
	for index, vpc := range vpcSetTypes {
		mapping := map[string]interface{}{
			"id":                    vpc.VpcId,
			"region_id":             vpc.RegionId,
			"status":                vpc.Status,
			"vpc_name":              vpc.VpcName,
			"vswitch_ids":           vpc.VSwitchIds.VSwitchId,
			"cidr_block":            vpc.CidrBlock,
			"vrouter_id":            vpc.VRouterId,
			"route_table_id":        route_tables[index],
			"description":           vpc.Description,
			"is_default":            vpc.IsDefault,
			"creation_time":         vpc.CreationTime,
			"ipv6_cidr_block":       vpc.Ipv6CidrBlock,
			"resource_group_id":     vpc.ResourceGroupId,
			"router_id":             vpc.VRouterId,
			"secondary_cidr_blocks": vpc.SecondaryCidrBlocks.SecondaryCidrBlock,
			"user_cidrs":            vpc.UserCidrs.UserCidr,
			"vpc_id":                fmt.Sprint(vpc.VpcId),
			"tags":                  vpcService.tagToMap(vpc.Tags.Tag),
		}

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(vpc.VpcId))
			names = append(names, vpc.VpcName)
			s = append(s, mapping)
			continue
		}

		ids = append(ids, vpc.VpcId)
		names = append(names, vpc.VpcName)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("vpcs", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
