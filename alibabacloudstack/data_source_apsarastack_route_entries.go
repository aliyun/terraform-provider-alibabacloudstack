package alibabacloudstack

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackRouteEntries() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceAlibabacloudStackRouteEntriesRead,
		Schema: map[string]*schema.Schema{
			"route_table_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cidr_block": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"entries": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"next_hop_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"route_table_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cidr_block": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackRouteEntriesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := vpc.CreateDescribeRouteTablesRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.RouteTableId = d.Get("route_table_id").(string)
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)

	var allRouteEntries []vpc.RouteEntry
	invoker := NewInvoker()
	for {
		var raw interface{}
		var err error
		err = invoker.Run(func() error {
			raw, err = client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
				return vpcClient.DescribeRouteTables(request)
			})
			return err
		})
		response, ok := raw.(*vpc.DescribeRouteTablesResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			} else {
				errmsg = "Invalid response type"
			}

			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_route_entries", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		if len(response.RouteTables.RouteTable) < 1 {
			break
		}

		for _, entries := range response.RouteTables.RouteTable[0].RouteEntrys.RouteEntry {
			if instance_id, ok := d.GetOk("instance_id"); ok && entries.InstanceId != instance_id.(string) {
				continue
			}
			if route_entry_type, ok := d.GetOk("type"); ok && entries.Type != route_entry_type.(string) {
				continue
			}
			if cidr_block, ok := d.GetOk("cidr_block"); ok && entries.DestinationCidrBlock != cidr_block.(string) {
				continue
			}
			allRouteEntries = append(allRouteEntries, entries)
		}

		if len(response.RouteTables.RouteTable) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return errmsgs.WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	return RouteEntriesDecriptionAttributes(d, allRouteEntries, meta)
}

func RouteEntriesDecriptionAttributes(d *schema.ResourceData, entries []vpc.RouteEntry, meta interface{}) error {
	var ids []string
	var s []map[string]interface{}
	for _, entry := range entries {
		mapping := map[string]interface{}{
			"route_table_id": entry.RouteTableId,
			"instance_id":    entry.InstanceId,
			"status":         entry.Status,
			"next_hop_type":  entry.NextHopType,
			"type":           entry.Type,
			"cidr_block":     entry.DestinationCidrBlock,
		}
		ids = append(ids, entry.RouteTableId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("entries", s); err != nil {
		return errmsgs.WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), s); err != nil {
			return err
		}
	}
	return nil
}
