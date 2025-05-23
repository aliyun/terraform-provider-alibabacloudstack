package alibabacloudstack

import (
	"encoding/json"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackSlb() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 80),
				Default:      resource.PrefixedUniqueId("tf-lb-"),
			},
			"address_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"internet", "intranet"}, false),
			},
			"network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "vpc",
				ValidateFunc: validation.StringInSlice([]string{"classic", "vpc"}, false),
			},

			"specification": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"slb.s1.small", "slb.s2.medium", "slb.s2.small", "slb.s3.large", "slb.s3.medium", "slb.s3.small", "slb.s4.large"}, false),
			},

			"vswitch_id": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: slbInternetDiffSuppressFunc,
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"address": {
				Type:         schema.TypeString,
				Computed:     true,
				ForceNew:     true,
				Optional:     true,
				ValidateFunc: validation.IsIPAddress,
			},
			"ip_version": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Default:  "ipv4",
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackSlbCreate, resourceAlibabacloudStackSlbRead, resourceAlibabacloudStackSlbUpdate, resourceAlibabacloudStackSlbDelete)
	return resource
}

func resourceAlibabacloudStackSlbCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	slbService := SlbService{client}
	request := client.NewCommonRequest("POST", "Slb", "2014-05-15", "CreateLoadBalancerPro", "")
	network_type := d.Get("network_type").(string)
	request.QueryParams["LoadBalancerName"] = d.Get("name").(string)
	request.QueryParams["AddressType"] = strings.ToLower(string(Intranet))
	request.QueryParams["AddressIPVersion"] = d.Get("ip_version").(string)
	request.QueryParams["NetworkType"] = network_type

	if network_type == "vpc" {
		if v, ok := d.GetOk("vswitch_id"); ok && v.(string) != "" {
			request.QueryParams["VSwitchId"] = d.Get("vswitch_id").(string)
		} else {
			return errmsgs.WrapError(errmsgs.Error("VSwitchId is required when network_type is vpc"))
		}
	}

	if v, ok := d.GetOk("specification"); ok && v.(string) != "" {
		request.QueryParams["LoadBalancerSpec"] = v.(string)
	}
	if v, ok := d.GetOk("address"); ok && v.(string) != "" {
		request.QueryParams["Address"] = v.(string)
	}

	bresponse, err := client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse != nil {
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			if errmsg != "" {
				return errmsgs.WrapError(errmsgs.Error(errmsg))
			}
		} else {
			return errmsgs.WrapError(err)
		}
	}
	response := slb.CreateLoadBalancerResponse{}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	d.SetId(response.LoadBalancerId)

	if err := slbService.WaitForSlb(response.LoadBalancerId, Active, DefaultTimeout); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}

func resourceAlibabacloudStackSlbRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	slbService := SlbService{client}
	object, err := slbService.DescribeSlb(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	d.Set("network_type", object.NetworkType)
	d.Set("name", object.LoadBalancerName)
	d.Set("address_type", object.AddressType)
	d.Set("vswitch_id", object.VSwitchId)
	d.Set("address", object.Address)
	d.Set("specification", object.LoadBalancerSpec)
	d.Set("ip_version", object.AddressIPVersion)

	tags, _ := slbService.DescribeTags(d.Id(), nil, TagResourceInstance)
	if len(tags) > 0 {
		if err := d.Set("tags", slbService.tagsToMap(tags)); err != nil {
			return errmsgs.WrapError(err)
		}
	}
	return nil
}

func resourceAlibabacloudStackSlbUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	slbService := SlbService{client}
	d.Partial(true)

	if err := slbService.setInstanceTags(d, TagResourceInstance); err != nil {
		return errmsgs.WrapError(err)
	}

	// Update Address

	if d.IsNewResource() {
		d.Partial(false)
		return nil
	}

	if d.HasChange("specification") {
		request := slb.CreateModifyLoadBalancerInstanceSpecRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.LoadBalancerId = d.Id()
		request.LoadBalancerSpec = d.Get("specification").(string)
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.ModifyLoadBalancerInstanceSpec(request)
		})
		if err != nil {
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*slb.ModifyLoadBalancerInstanceSpecResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	if d.HasChange("name") {
		request := slb.CreateSetLoadBalancerNameRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.LoadBalancerId = d.Id()
		request.LoadBalancerName = d.Get("name").(string)
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.SetLoadBalancerName(request)
		})
		if err != nil {
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*slb.SetLoadBalancerNameResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	update := false
	modifyLoadBalancerInternetSpecRequest := slb.CreateModifyLoadBalancerInternetSpecRequest()
	client.InitRpcRequest(*modifyLoadBalancerInternetSpecRequest.RpcRequest)
	modifyLoadBalancerInternetSpecRequest.LoadBalancerId = d.Id()
	if update {
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.ModifyLoadBalancerInternetSpec(modifyLoadBalancerInternetSpecRequest)
		})
		if err != nil {
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*slb.ModifyLoadBalancerInternetSpecResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), modifyLoadBalancerInternetSpecRequest.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(modifyLoadBalancerInternetSpecRequest.GetActionName(), raw, modifyLoadBalancerInternetSpecRequest.RpcRequest, modifyLoadBalancerInternetSpecRequest)
	}

	update = false
	modifyLoadBalancerPayTypeRequest := slb.CreateModifyLoadBalancerPayTypeRequest()
	client.InitRpcRequest(*modifyLoadBalancerPayTypeRequest.RpcRequest)
	modifyLoadBalancerPayTypeRequest.LoadBalancerId = d.Id()
	if update {
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.ModifyLoadBalancerPayType(modifyLoadBalancerPayTypeRequest)
		})
		if err != nil {
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*slb.ModifyLoadBalancerPayTypeResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), modifyLoadBalancerPayTypeRequest.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(modifyLoadBalancerPayTypeRequest.GetActionName(), raw, modifyLoadBalancerPayTypeRequest.RpcRequest, modifyLoadBalancerPayTypeRequest)
	}
	d.Partial(false)

	return nil
}

func resourceAlibabacloudStackSlbDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	slbService := SlbService{client}

	request := slb.CreateDeleteLoadBalancerRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.LoadBalancerId = d.Id()

	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DeleteLoadBalancer(request)
	})
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidLoadBalancerId.NotFound"}) {
			return nil
		}
		errmsg := ""
		if raw != nil {
			response, ok := raw.(*slb.DeleteLoadBalancerResponse)
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return errmsgs.WrapError(slbService.WaitForSlb(d.Id(), Deleted, DefaultTimeoutMedium))
}
