package alibabacloudstack

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
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
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "slb", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
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
		if IsExpectedErrors(err, []string{"OrderFailed"}) {
			return WrapError(err)
		}
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_slb", request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*slb.CreateLoadBalancerResponse)
	d.SetId(response.LoadBalancerId)

	if err := slbService.WaitForSlb(response.LoadBalancerId, Active, DefaultTimeout); err != nil {
		return WrapError(err)
	}

	return resourceAlibabacloudStackSlbUpdate(d, meta)
}

func resourceAlibabacloudStackSlbRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	slbService := SlbService{client}
	object, err := slbService.DescribeSlb(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", object.LoadBalancerName)
	d.Set("address_type", object.AddressType)
	d.Set("vswitch_id", object.VSwitchId)
	d.Set("address", object.Address)
	d.Set("specification", object.LoadBalancerSpec)

	tags, _ := slbService.DescribeTags(d.Id(), nil, TagResourceInstance)
	if len(tags) > 0 {
		if err := d.Set("tags", slbService.tagsToMap(tags)); err != nil {
			return WrapError(err)
		}
	}
	return nil
}

func resourceAlibabacloudStackSlbUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AlibabacloudStackClient)
	slbService := SlbService{client}
	d.Partial(true)

	// set instance tags

	if err := slbService.setInstanceTags(d, TagResourceInstance); err != nil {
		return WrapError(err)
	}

	if d.IsNewResource() {
		d.Partial(false)
		return resourceAlibabacloudStackSlbRead(d, meta)
	}

	if d.HasChange("specification") {
		request := slb.CreateModifyLoadBalancerInstanceSpecRequest()
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		request.RegionId = client.RegionId
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "slb", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
		request.LoadBalancerId = d.Id()
		request.LoadBalancerSpec = d.Get("specification").(string)
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.ModifyLoadBalancerInstanceSpec(request)
		})
		if err != nil {
			WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	if d.HasChange("name") {
		request := slb.CreateSetLoadBalancerNameRequest()
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		request.RegionId = client.RegionId
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "slb", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
		request.LoadBalancerId = d.Id()
		request.LoadBalancerName = d.Get("name").(string)
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.SetLoadBalancerName(request)
		})
		if err != nil {
			WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		//d.SetPartial("name")
	}
	update := false
	modifyLoadBalancerInternetSpecRequest := slb.CreateModifyLoadBalancerInternetSpecRequest()
	if strings.ToLower(client.Config.Protocol) == "https" {
		modifyLoadBalancerInternetSpecRequest.Scheme = "https"
	} else {
		modifyLoadBalancerInternetSpecRequest.Scheme = "http"
	}
	modifyLoadBalancerInternetSpecRequest.Headers = map[string]string{"RegionId": client.RegionId}
	modifyLoadBalancerInternetSpecRequest.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "slb", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	modifyLoadBalancerInternetSpecRequest.RegionId = client.RegionId
	modifyLoadBalancerInternetSpecRequest.LoadBalancerId = d.Id()
	if update {
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.ModifyLoadBalancerInternetSpec(modifyLoadBalancerInternetSpecRequest)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), modifyLoadBalancerInternetSpecRequest.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
		addDebug(modifyLoadBalancerInternetSpecRequest.GetActionName(), raw, modifyLoadBalancerInternetSpecRequest.RpcRequest, modifyLoadBalancerInternetSpecRequest)
	}

	update = false
	modifyLoadBalancerPayTypeRequest := slb.CreateModifyLoadBalancerPayTypeRequest()
	modifyLoadBalancerPayTypeRequest.RegionId = client.RegionId
	if strings.ToLower(client.Config.Protocol) == "https" {
		modifyLoadBalancerPayTypeRequest.Scheme = "https"
	} else {
		modifyLoadBalancerPayTypeRequest.Scheme = "http"
	}
	modifyLoadBalancerInternetSpecRequest.Headers = map[string]string{"RegionId": client.RegionId}
	modifyLoadBalancerInternetSpecRequest.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "slb", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	modifyLoadBalancerPayTypeRequest.LoadBalancerId = d.Id()

	if update {
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.ModifyLoadBalancerPayType(modifyLoadBalancerPayTypeRequest)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), modifyLoadBalancerPayTypeRequest.GetActionName(), AlibabacloudStackSdkGoERROR)
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
	request.RegionId = client.RegionId
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "slb", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.LoadBalancerId = d.Id()

	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DeleteLoadBalancer(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidLoadBalancerId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return WrapError(slbService.WaitForSlb(d.Id(), Deleted, DefaultTimeoutMedium))
}
