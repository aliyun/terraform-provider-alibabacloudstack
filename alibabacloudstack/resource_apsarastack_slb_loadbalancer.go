package alibabacloudstack

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackSlb() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackSlbCreate,
		Read:   resourceAlibabacloudStackSlbRead,
		Update: resourceAlibabacloudStackSlbUpdate,
		Delete: resourceAlibabacloudStackSlbDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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
		},
	}
}

func resourceAlibabacloudStackSlbCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	slbService := SlbService{client}
	request := slb.CreateCreateLoadBalancerRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.LoadBalancerName = d.Get("name").(string)
	request.AddressType = strings.ToLower(string(Intranet))
	request.ClientToken = buildClientToken(request.GetActionName())

	if v, ok := d.GetOk("name"); ok && v.(string) != "" {
		request.LoadBalancerName = strings.ToLower(v.(string))
	}
	if v, ok := d.GetOk("address_type"); ok && v.(string) != "" {
		request.AddressType = strings.ToLower(v.(string))
	}

	if v, ok := d.GetOk("vswitch_id"); ok && v.(string) != "" {
		request.VSwitchId = v.(string)
	}

	if v, ok := d.GetOk("specification"); ok && v.(string) != "" {
		request.LoadBalancerSpec = v.(string)
	}

	var raw interface{}

	invoker := Invoker{}
	invoker.AddCatcher(Catcher{"OperationFailed.TokenIsProcessing", 10, 5})
	log.Printf("[DEBUG] slb request %v", request)
	if err := invoker.Run(func() error {
		resp, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.CreateLoadBalancer(request)
		})
		raw = resp
		return err
	}); err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"OrderFailed"}) {
			return errmsgs.WrapError(err)
		}
		errmsg := ""
		if raw != nil {
			response, ok := raw.(*slb.CreateLoadBalancerResponse)
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_slb", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*slb.CreateLoadBalancerResponse)
	d.SetId(response.LoadBalancerId)

	if err := slbService.WaitForSlb(response.LoadBalancerId, Active, DefaultTimeout); err != nil {
		return errmsgs.WrapError(err)
	}

	return resourceAlibabacloudStackSlbUpdate(d, meta)
}

func resourceAlibabacloudStackSlbRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
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

	d.Set("name", object.LoadBalancerName)
	d.Set("address_type", object.AddressType)
	d.Set("vswitch_id", object.VSwitchId)
	d.Set("address", object.Address)
	d.Set("specification", object.LoadBalancerSpec)

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

	if d.IsNewResource() {
		d.Partial(false)
		return resourceAlibabacloudStackSlbRead(d, meta)
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

	return resourceAlibabacloudStackSlbRead(d, meta)
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
