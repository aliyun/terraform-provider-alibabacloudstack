package alibabacloudstack

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackRouteEntry() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackRouteEntryCreate,
		Read:   resourceAlibabacloudStackRouteEntryRead,
		Delete: resourceAlibabacloudStackRouteEntryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"router_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Deprecated:  "Attribute router_id has been deprecated and suggest removing it from your template.",
			},
			"route_table_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"destination_cidrblock": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"nexthop_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"nexthop_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
			},
		},
	}
}

func resourceAlibabacloudStackRouteEntryCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}

	rtId := d.Get("route_table_id").(string)
	cidr := d.Get("destination_cidrblock").(string)
	nt := d.Get("nexthop_type").(string)
	ni := d.Get("nexthop_id").(string)

	table, err := vpcService.QueryRouteTableById(rtId)

	if err != nil {
		return errmsgs.WrapError(err)
	}
	request := vpc.CreateCreateRouteEntryRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.RouteTableId = rtId
	request.DestinationCidrBlock = cidr
	request.NextHopType = nt
	request.NextHopId = ni
	request.ClientToken = buildClientToken(request.GetActionName())
	request.RouteEntryName = d.Get("name").(string)
	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		if err := vpcService.WaitForAllRouteEntriesAvailable(rtId, DefaultTimeout); err != nil {
			return resource.NonRetryableError(err)
		}
		args := *request
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.CreateRouteEntry(&args)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"TaskConflict", "IncorrectRouteEntryStatus", errmsgs.Throttling, "IncorrectVpcStatus"}) {
				return resource.RetryableError(err)
			}
			errmsg := ""
			if response, ok := raw.(*vpc.CreateRouteEntryResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_route_entry", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"RouterEntryConflict.Duplicated"}) {
			en, err := vpcService.DescribeRouteEntry(rtId + ":" + table.VRouterId + ":" + cidr + ":" + nt + ":" + ni)
			if err != nil {
				return errmsgs.WrapError(err)
			}
			return errmsgs.WrapError(errmsgs.Error("The route entry %s has already existed. "+
				"Please import it using ID '%s:%s:%s:%s:%s' or specify a new 'destination_cidrblock' and try again.",
				en.DestinationCidrBlock, en.RouteTableId, table.VRouterId, en.DestinationCidrBlock, en.NextHopType, ni))
		}
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_route_entry", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
	}

	d.SetId(rtId + ":" + table.VRouterId + ":" + cidr + ":" + nt + ":" + ni)

	if err := vpcService.WaitForRouteEntry(d.Id(), Available, DefaultTimeout); err != nil {
		return errmsgs.WrapError(err)
	}
	return resourceAlibabacloudStackRouteEntryRead(d, meta)
}

func resourceAlibabacloudStackRouteEntryRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}
	parts, err := ParseResourceId(d.Id(), 5)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	object, err := vpcService.DescribeRouteEntry(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("router_id", parts[1])
	d.Set("route_table_id", object.RouteTableId)
	d.Set("destination_cidrblock", object.DestinationCidrBlock)
	d.Set("nexthop_type", object.NextHopType)
	d.Set("nexthop_id", object.InstanceId)
	d.Set("name", object.RouteEntryName)
	return nil
}

func resourceAlibabacloudStackRouteEntryDelete(d *schema.ResourceData, meta interface{}) error {
	request, err := buildAlibabacloudStackRouteEntryDeleteArgs(d, meta)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	client := meta.(*connectivity.AlibabacloudStackClient)
	client.InitRpcRequest(*request.RpcRequest)
	vpcService := VpcService{client}
	parts, err := ParseResourceId(d.Id(), 5)
	rtId := parts[0]
	if err := vpcService.WaitForAllRouteEntriesAvailable(rtId, DefaultTimeout); err != nil {
		return errmsgs.WrapError(err)
	}
	retryTimes := 7
	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteRouteEntry(request)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"IncorrectVpcStatus", "TaskConflict", "IncorrectRouteEntryStatus", "Forbbiden", "UnknownError"}) {
				time.Sleep(time.Duration(retryTimes) * time.Second)
				retryTimes += 7
				return resource.RetryableError(err)
			}
			errmsg := ""
			if response, ok := raw.(*vpc.DeleteRouteEntryResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_route_entry", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidRouteEntry.NotFound"}) {
			return nil
		}
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return errmsgs.WrapError(vpcService.WaitForRouteEntry(d.Id(), Deleted, DefaultTimeout))
}

func buildAlibabacloudStackRouteEntryDeleteArgs(d *schema.ResourceData, meta interface{}) (*vpc.DeleteRouteEntryRequest, error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := vpc.CreateDeleteRouteEntryRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.RouteTableId = d.Get("route_table_id").(string)
	request.DestinationCidrBlock = d.Get("destination_cidrblock").(string)

	if v := d.Get("destination_cidrblock").(string); v != "" {
		request.DestinationCidrBlock = v
	}

	if v := d.Get("nexthop_id").(string); v != "" {
		request.NextHopId = v
	}

	return request, nil
}
