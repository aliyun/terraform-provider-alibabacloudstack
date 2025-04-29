package alibabacloudstack

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/denverdino/aliyungo/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackVpnGateway() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validation.StringLenBetween(1, 128),
				Computed:      true,
				Deprecated:    "Field 'name' is deprecated and will be removed in a future release. Please use new field 'vpn_gateway_name' instead.",
				ConflictsWith: []string{"vpn_gateway_name"},
			},
			"vpn_gateway_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validation.StringLenBetween(1, 128),
				Computed:      true,
				ConflictsWith: []string{"name"},
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      PostPaid,
				ValidateFunc: validation.StringInSlice([]string{string(common.PrePaid), string(common.PostPaid)}, false),
			},

			"period": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          1,
				ValidateFunc:     validation.Any(validation.IntBetween(1, 9), validation.IntInSlice([]int{12, 24, 36})),
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
			},

			"bandwidth": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntInSlice([]int{5, 10, 20, 50, 100, 200, 500, 1000}),
			},

			"enable_ipsec": {
				Type:          schema.TypeBool,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'enable_ipsec' is deprecated and will be removed in a future release. Please use new field 'ipsec_vpn' instead.",
				ConflictsWith: []string{"ipsec_vpn"},
			},
			"ipsec_vpn": {
				Type:          schema.TypeBool,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"enable_ipsec"},
			},

			"enable_ssl": {
				Type:          schema.TypeBool,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'enable_ssl' is deprecated and will be removed in a future release. Please use new field 'ssl_vpn' instead.",
				ConflictsWith: []string{"ssl_vpn"},
			},
			"ssl_vpn": {
				Type:          schema.TypeBool,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"enable_ssl"},
			},

			"ssl_connections": {
				Type:             schema.TypeInt,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: vpnSslConnectionsDiffSuppressFunc,
				Deprecated:       "Field 'ssl_connections' is deprecated and will be removed in a future release. Please use new field 'ssl_max_connections' instead.",
				ConflictsWith:    []string{"ssl_max_connections"},
			},
			"ssl_max_connections": {
				Type:             schema.TypeInt,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: vpnSslConnectionsDiffSuppressFunc,
				ConflictsWith:    []string{"ssl_connections"},
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},

			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"internet_ip": {
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
			"tags": tagsSchema(),
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackVpnGatewayCreate, resourceAlibabacloudStackVpnGatewayRead, resourceAlibabacloudStackVpnGatewayUpdate, resourceAlibabacloudStackVpnGatewayDelete)
	return resource
}

func resourceAlibabacloudStackVpnGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpnGatewayService := VpnGatewayService{client}
	request := vpc.CreateCreateVpnGatewayRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.VpcId = d.Get("vpc_id").(string)

	if v, ok := connectivity.GetResourceDataOk(d, "vpn_gateway_name", "name"); ok && v.(string) != "" {
		request.Name = v.(string)
	}

	if v, ok := d.GetOk("vswitch_id"); ok && v.(string) != "" {
		request.VSwitchId = d.Get("vswitch_id").(string)
	}

	if v, ok := d.GetOk("instance_charge_type"); ok && v.(string) != "" {
		if v.(string) == string(PostPaid) {
			request.InstanceChargeType = string("POSTPAY")
		} else {
			request.InstanceChargeType = string("PREPAY")
		}
	}

	if v, ok := d.GetOk("period"); ok && v.(int) != 0 && request.InstanceChargeType == string("PREPAY") {
		request.Period = requests.NewInteger(v.(int))
	}

	request.Bandwidth = requests.NewInteger(d.Get("bandwidth").(int))

	if v, ok := connectivity.GetResourceDataOk(d, "ipsec_vpn", "enable_ipsec"); ok {
		request.EnableIpsec = requests.NewBoolean(v.(bool))
	} else {
		//Default must be nil if computed
		request.EnableIpsec = requests.NewBoolean(true)
	}

	if v, ok := connectivity.GetResourceDataOk(d, "ssl_vpn", "enable_ssl"); ok {
		request.EnableSsl = requests.NewBoolean(v.(bool))
	} else {
		//Default must be nil if computed
		request.EnableSsl = requests.NewBoolean(false)
	}

	if v, ok := connectivity.GetResourceDataOk(d, "ssl_max_connections", "ssl_connections"); ok {
		request.SslConnections = requests.NewInteger(v.(int))
	} else {
		//Default must be nil if computed
		request.SslConnections = requests.NewInteger(5)
	}

	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.CreateVpnGateway(request)
	})

	if err != nil {
		errmsg := ""
		if raw != nil {
			response, ok := raw.(*vpc.CreateVpnGatewayResponse)
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_vpn_gateway", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*vpc.CreateVpnGatewayResponse)
	d.SetId(response.VpnGatewayId)

	time.Sleep(10 * time.Second)
	if err := vpnGatewayService.WaitForVpnGateway(d.Id(), Active, 2*DefaultTimeout); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}

func resourceAlibabacloudStackVpnGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpnGatewayService := VpnGatewayService{client}
	vpcService := VpcService{client}

	object, err := vpnGatewayService.DescribeVpnGateway(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	connectivity.SetResourceData(d, object.Name, "vpn_gateway_name", "name")
	d.Set("description", object.Description)
	d.Set("vpc_id", object.VpcId)
	d.Set("internet_ip", object.InternetIp)
	d.Set("status", object.Status)
	d.Set("vswitch_id", object.VSwitchId)
	connectivity.SetResourceData(d, "enable" == strings.ToLower(object.IpsecVpn), "ipsec_vpn", "enable_ipsec")
	connectivity.SetResourceData(d, "enable" == strings.ToLower(object.SslVpn), "ssl_vpn", "enable_ssl")
	connectivity.SetResourceData(d, object.SslMaxConnections, "ssl_max_connections", "ssl_connections")
	d.Set("business_status", object.BusinessStatus)

	spec := strings.Split(object.Spec, "M")[0]
	bandwidth, err := strconv.Atoi(spec)

	if err == nil {
		d.Set("bandwidth", bandwidth)
	} else {
		return errmsgs.WrapError(err)
	}

	if string("PostpayByFlow") == object.ChargeType {
		d.Set("instance_charge_type", string(PostPaid))
	} else {
		d.Set("instance_charge_type", string(PrePaid))
	}
	listTagResourcesObject, err := vpcService.ListTagResources(d.Id(), "VpnGateWay")
	if err != nil {
		return errmsgs.WrapError(err)
	}
	d.Set("tags", tagsToMap(listTagResourcesObject))

	return nil
}

func resourceAlibabacloudStackVpnGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}

	request := vpc.CreateModifyVpnGatewayAttributeRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.VpnGatewayId = d.Id()
	update := false
	if d.HasChange("tags") {
		if err := vpcService.SetResourceTags(d, "VpnGateWay"); err != nil {
			return errmsgs.WrapError(err)
		}
		//d.SetPartial("tags")
	}
	d.Partial(true)
	if d.HasChanges("name", "vpn_gateway_name") {
		request.Name = connectivity.GetResourceData(d, "vpn_gateway_name", "name").(string)
		update = true
	}

	if d.HasChange("description") {
		request.Description = d.Get("description").(string)
		update = true
	}

	if update {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ModifyVpnGatewayAttribute(request)
		})
		if err != nil {
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*vpc.ModifyVpnGatewayAttributeResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		//d.SetPartial("name")
		//d.SetPartial("description")
	}

	if d.IsNewResource() {
		d.Partial(false)
		return nil
	}

	if d.HasChange("bandwidth") {
		return fmt.Errorf("Now Cann't Support modify vpn gateway bandwidth, try to modify on the web console")
	}

	if d.HasChanges("enable_ipsec", "enable_ssl", "ipsec_vpn", "ssl_vpn") {
		return fmt.Errorf("Now Cann't Support modify ipsec/ssl switch, try to modify on the web console")
	}
	if d.HasChanges("ssl_max_connections", "ssl_connections") {
		return fmt.Errorf("Now Cann't Support modify ssl_max_connections, try to modify on the web console")
	}

	d.Partial(false)
	return nil
}

func resourceAlibabacloudStackVpnGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpnGatewayService := VpnGatewayService{client}

	request := vpc.CreateDeleteVpnGatewayRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.VpnGatewayId = d.Id()
	request.ClientToken = buildClientToken(request.GetActionName())

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		args := *request
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteVpnGateway(&args)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"VpnGateway.Configuring"}) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(err)
			}
			/*Vpn known issue: while the vpn is configuring, it will return unknown error*/
			if errmsgs.IsExpectedErrors(err, []string{"UnknownError"}) {
				return resource.RetryableError(err)
			}
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*vpc.DeleteVpnGatewayResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidVpnGatewayInstanceId.NotFound"}) {
			return nil
		}
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return errmsgs.WrapError(vpnGatewayService.WaitForVpnGateway(d.Id(), Deleted, DefaultTimeoutMedium))
}
