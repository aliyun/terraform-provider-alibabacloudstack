package alibabacloudstack

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackVpnConnection() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"customer_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"vpn_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:true,
				ValidateFunc: validation.StringLenBetween(2, 128),
				Deprecated:  "Field 'name' is deprecated and will be removed in a future release. Please use new field 'vpn_connection_name' instead.",
				ConflictsWith: []string{"vpn_connection_name"},
			},

			"vpn_connection_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:true,
				ValidateFunc: validation.StringLenBetween(2, 128),
				ConflictsWith: []string{"name"},
			},

			"local_subnet": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateCIDRNetworkAddress,
				},
				MinItems: 1,
				MaxItems: 10,
			},

			"remote_subnet": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateCIDRNetworkAddress,
				},
				MinItems: 1,
				MaxItems: 10,
			},

			"effect_immediately": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"ike_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"psk": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(1, 100),
						},
						"ike_version": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      IKE_VERSION_1,
							ValidateFunc: validation.StringInSlice([]string{IKE_VERSION_1, IKE_VERSION_2}, false),
						},
						"ike_mode": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      IKE_MODE_MAIN,
							ValidateFunc: validation.StringInSlice([]string{IKE_MODE_MAIN, IKE_MODE_AGGRESSIVE}, false),
						},
						"ike_enc_alg": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      VPN_ENC_AES,
							ValidateFunc: validation.StringInSlice([]string{VPN_ENC_AES, VPN_ENC_AES_3DES, VPN_ENC_AES_192, VPN_ENC_AES_256, VPN_ENC_AES_DES}, false),
						},
						"ike_auth_alg": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      VPN_AUTH_SHA,
							ValidateFunc: validation.StringInSlice([]string{VPN_AUTH_SHA, VPN_AUTH_MD5, VPN_AUTH_SHA256, VPN_AUTH_SHA386, VPN_AUTH_SHA512}, false),
						},
						"ike_pfs": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      VPN_PFS_G2,
							ValidateFunc: validation.StringInSlice([]string{VPN_PFS_G1, VPN_PFS_G2, VPN_PFS_G5, VPN_PFS_G14, VPN_PFS_G24}, false),
						},
						"ike_lifetime": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      86400,
							ValidateFunc: validation.IntBetween(0, 86400),
						},
						"ike_local_id": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(1, 100),
						},
						"ike_remote_id": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(1, 100),
						},
					},
				},
			},

			"ipsec_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipsec_enc_alg": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      VPN_ENC_AES,
							ValidateFunc: validation.StringInSlice([]string{VPN_ENC_AES, VPN_ENC_AES_3DES, VPN_ENC_AES_192, VPN_ENC_AES_256, VPN_ENC_AES_DES}, false),
						},
						"ipsec_auth_alg": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      VPN_AUTH_SHA,
							ValidateFunc: validation.StringInSlice([]string{VPN_AUTH_SHA, VPN_AUTH_MD5, VPN_AUTH_SHA256, VPN_AUTH_SHA386, VPN_AUTH_SHA512}, false),
						},
						"ipsec_pfs": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      VPN_PFS_G2,
							ValidateFunc: validation.StringInSlice([]string{VPN_PFS_G1, VPN_PFS_G2, VPN_PFS_G5, VPN_PFS_G14, VPN_PFS_G24, VPN_PFS_DISABLED}, false),
						},
						"ipsec_lifetime": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 86400),
						},
					},
				},
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackVpnConnectionCreate, resourceAlibabacloudStackVpnConnectionRead, resourceAlibabacloudStackVpnConnectionUpdate, resourceAlibabacloudStackVpnConnectionDelete)
	return resource
}

func resourceAlibabacloudStackVpnConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpnGatewayService := VpnGatewayService{client}
	request, err := buildAlibabacloudStackVpnConnectionArgs(d, meta)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	client.InitRpcRequest(*request.RpcRequest)
	var response *vpc.CreateVpnConnectionResponse
	var ok bool
	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		args := *request
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.CreateVpnConnection(&args)
		})
		response, ok = raw.(*vpc.CreateVpnConnectionResponse)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"VpnGateway.Configuring"}) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(err)
			}
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_vpn_connection", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_vpn_connection", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
	}

	d.SetId(response.VpnConnectionId)

	if err := vpnGatewayService.WaitForVpnConnection(d.Id(), Null, DefaultTimeoutMedium); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}

func resourceAlibabacloudStackVpnConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpnGatewayService := VpnGatewayService{client}
	response, err := vpnGatewayService.DescribeVpnConnection(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	d.Set("customer_gateway_id", response.CustomerGatewayId)
	d.Set("vpn_gateway_id", response.VpnGatewayId)
	connectivity.SetResourceData(d, response.Name, "vpn_connection_name", "name")

	localSubnet := strings.Split(response.LocalSubnet, ",")
	d.Set("local_subnet", localSubnet)

	remoteSubnet := strings.Split(response.RemoteSubnet, ",")
	d.Set("remote_subnet", remoteSubnet)

	d.Set("effect_immediately", response.EffectImmediately)
	d.Set("status", response.Status)

	if err := d.Set("ike_config", vpnGatewayService.ParseIkeConfig(response.IkeConfig)); err != nil {
		return errmsgs.WrapError(err)
	}

	if err := d.Set("ipsec_config", vpnGatewayService.ParseIpsecConfig(response.IpsecConfig)); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}

func resourceAlibabacloudStackVpnConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpnGatewayService := VpnGatewayService{client}
	request := vpc.CreateModifyVpnConnectionAttributeRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.VpnConnectionId = d.Id()

	if d.HasChanges("name", "vpn_connection_name") {
		request.Name = connectivity.GetResourceData(d, "vpn_connection_name", "name").(string)
	}

	request.LocalSubnet = vpnGatewayService.AssembleNetworkSubnetToString(d.Get("local_subnet").(*schema.Set).List())
	request.RemoteSubnet = vpnGatewayService.AssembleNetworkSubnetToString(d.Get("remote_subnet").(*schema.Set).List())

	/* If not set effect_immediately value, VPN connection will automatically set the value to false*/
	if v, ok := d.GetOk("effect_immediately"); ok {
		request.EffectImmediately = requests.NewBoolean(v.(bool))
	}

	if d.HasChange("ike_config") {
		ike_config, err := vpnGatewayService.AssembleIkeConfig(d.Get("ike_config").([]interface{}))
		if err != nil {
			return errmsgs.WrapError(err)
		}
		request.IkeConfig = ike_config
	}

	if d.HasChange("ipsec_config") {
		ipsec_config, err := vpnGatewayService.AssembleIpsecConfig(d.Get("ipsec_config").([]interface{}))
		if err != nil {
			return errmsgs.WrapError(err)
		}
		request.IpsecConfig = ipsec_config
	}

	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.ModifyVpnConnectionAttribute(request)
	})
	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return nil
}

func resourceAlibabacloudStackVpnConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpnGatewayService := VpnGatewayService{client}
	request := vpc.CreateDeleteVpnConnectionRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.VpnConnectionId = d.Id()

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		args := *request
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteVpnConnection(&args)
		})

		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"VpnGateway.Configuring"}) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(err)
			}
			errmsg := ""
			bresponse, ok := raw.(*responses.CommonResponse)
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidVpnConnectionInstanceId.NotFound"}) {
			return nil
		}
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return errmsgs.WrapError(vpnGatewayService.WaitForVpnConnection(d.Id(), Deleted, DefaultTimeout))
}

func buildAlibabacloudStackVpnConnectionArgs(d *schema.ResourceData, meta interface{}) (*vpc.CreateVpnConnectionRequest, error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpnGatewayService := VpnGatewayService{client}

	request := vpc.CreateCreateVpnConnectionRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.CustomerGatewayId = d.Get("customer_gateway_id").(string)
	request.VpnGatewayId = d.Get("vpn_gateway_id").(string)
	request.LocalSubnet = vpnGatewayService.AssembleNetworkSubnetToString(d.Get("local_subnet").(*schema.Set).List())
	request.RemoteSubnet = vpnGatewayService.AssembleNetworkSubnetToString(d.Get("remote_subnet").(*schema.Set).List())

	if v := connectivity.GetResourceData(d, "vpn_connection_name", "name"); v != ""{
		request.Name = v.(string)
	}

	if v, ok := d.GetOk("effect_immediately"); ok {
		request.EffectImmediately = requests.NewBoolean(v.(bool))
	}

	if v, ok := d.GetOk("ike_config"); ok {
		ikeConfig, err := vpnGatewayService.AssembleIkeConfig(v.([]interface{}))
		if err != nil {
			return nil, errmsgs.WrapError(err)
		}
		request.IkeConfig = ikeConfig
	}

	if v, ok := d.GetOk("ipsec_config"); ok {
		ipsecConfig, err := vpnGatewayService.AssembleIpsecConfig(v.([]interface{}))
		if err != nil {
			return nil, fmt.Errorf("wrong ipsec_config: %#v", err)
		}
		request.IpsecConfig = ipsecConfig
	}
	request.ClientToken = buildClientToken(request.GetActionName())

	return request, nil
}