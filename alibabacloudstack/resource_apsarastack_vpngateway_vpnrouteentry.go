package alibabacloudstack

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackVpnRouteEntry() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"vpn_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"next_hop": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"route_dest": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"weight": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 100}),
			},

			"publish_vpc": {
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackVpnRouteEntryCreate, resourceAlibabacloudStackVpnRouteEntryRead, resourceAlibabacloudStackVpnRouteEntryUpdate, resourceAlibabacloudStackVpnRouteEntryDelete)
	return resource
}

func resourceAlibabacloudStackVpnRouteEntryCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpnRouteEntryService := VpnGatewayService{client}
	request := vpc.CreateCreateVpnRouteEntryRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.VpnGatewayId = d.Get("vpn_gateway_id").(string)
	request.RouteDest = d.Get("route_dest").(string)
	request.NextHop = d.Get("next_hop").(string)
	request.Weight = requests.NewInteger(d.Get("weight").(int))
	request.PublishVpc = requests.NewBoolean(d.Get("publish_vpc").(bool))
	request.ClientToken = buildClientToken(request.GetActionName())

	var raw interface{}
	wait := incrementalWait(5*time.Second, 5*time.Second)
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw1, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.CreateVpnRouteEntry(request)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"VpnGateway.Configuring"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw1, request.RpcRequest, request)
		raw = raw1
		return nil
	})

	if err != nil {
		errmsg := ""
		if raw != nil {
			response, ok := raw.(*vpc.CreateVpnRouteEntryResponse)
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	response, _ := raw.(*vpc.CreateVpnRouteEntryResponse)
	id := response.VpnInstanceId + ":" + response.NextHop + ":" + response.RouteDest
	d.SetId(id)

	if err := vpnRouteEntryService.WaitForVpnRouteEntry(d.Id(), Active, 2*DefaultTimeoutMedium); err != nil {
		return errmsgs.WrapError(err)
	}
	return nil
}

func resourceAlibabacloudStackVpnRouteEntryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpnRouteEntryService := VpnGatewayService{client}

	object, err := vpnRouteEntryService.DescribeVpnRouteEntry(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("weight", object.Weight)
	d.Set("next_hop", object.NextHop)
	d.Set("route_dest", object.RouteDest)
	//d.Set("old_weight", object.Weight)
	d.Set("vpn_gateway_id", object.VpnInstanceId)

	if object.State == "published" {
		d.Set("publish_vpc", true)
	} else {
		d.Set("publish_vpc", false)
	}

	return nil
}

func resourceAlibabacloudStackVpnRouteEntryUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	d.Partial(true)

	if d.HasChange("publish_vpc") {
		request := vpc.CreatePublishVpnRouteEntryRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.VpnGatewayId = d.Get("vpn_gateway_id").(string)
		request.RouteDest = d.Get("route_dest").(string)
		request.NextHop = d.Get("next_hop").(string)
		request.RouteType = "dbr"
		request.PublishVpc = requests.NewBoolean(d.Get("publish_vpc").(bool))

		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.PublishVpnRouteEntry(request)
		})
		if err != nil {
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*vpc.PublishVpnRouteEntryResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		//d.SetPartial("public_vpc")
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	if d.HasChange("weight") {
		request := vpc.CreateModifyVpnRouteEntryWeightRequest()
		client.InitRpcRequest(*request.RpcRequest)
		oldWeight, newWeight := d.GetChange("weight")
		request.VpnGatewayId = d.Get("vpn_gateway_id").(string)
		request.RouteDest = d.Get("route_dest").(string)
		request.NextHop = d.Get("next_hop").(string)
		request.Weight = requests.NewInteger(oldWeight.(int))
		request.NewWeight = requests.NewInteger(newWeight.(int))

		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ModifyVpnRouteEntryWeight(request)
		})
		if err != nil {
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*vpc.ModifyVpnRouteEntryWeightResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		//d.SetPartial("weight")
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	d.Partial(false)
	return nil
}

func resourceAlibabacloudStackVpnRouteEntryDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpnRouteEntryService := VpnGatewayService{client}

	request := vpc.CreateDeleteVpnRouteEntryRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.VpnGatewayId = d.Get("vpn_gateway_id").(string)
	request.RouteDest = d.Get("route_dest").(string)
	request.NextHop = d.Get("next_hop").(string)
	request.Weight = requests.NewInteger(d.Get("weight").(int))
	request.ClientToken = buildClientToken(request.GetActionName())

	wait := incrementalWait(5*time.Second, 5*time.Second)
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteVpnRouteEntry(request)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"VpnGateway.Configuring"}) {
				wait()
				return resource.RetryableError(err)
			}
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*vpc.DeleteVpnRouteEntryResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		return err
	}
	return errmsgs.WrapError(vpnRouteEntryService.WaitForVpnRouteEntry(d.Id(), Deleted, DefaultTimeoutMedium))
}