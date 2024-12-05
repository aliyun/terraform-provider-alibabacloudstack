package alibabacloudstack

import (
	"time"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackVpnCustomerGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackVpnCustomerGatewayCreate,
		Read:   resourceAlibabacloudStackVpnCustomerGatewayRead,
		Update: resourceAlibabacloudStackVpnCustomerGatewayUpdate,
		Delete: resourceAlibabacloudStackVpnCustomerGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"ip_address": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsIPAddress,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
				Deprecated:   "Field 'name' is deprecated and will be removed in a future release. Please use 'customer_gateway_name' instead.",
			},
			"customer_gateway_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},
		},
	}
}

func resourceAlibabacloudStackVpnCustomerGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpnGatewayService := VpnGatewayService{client}
	request := vpc.CreateCreateCustomerGatewayRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.IpAddress = d.Get("ip_address").(string)
	if v, err := connectivity.GetResourceData(d, reflect.TypeOf(""), "customer_gateway_name", "name"); err == nil && v.(string) != "" {
		request.Name = v.(string)
	} else if err != nil {
		return err
	}

	if v := d.Get("description").(string); v != "" {
		request.Description = v
	}
	request.ClientToken = buildClientToken(request.GetActionName())

	request.Headers["x-acs-organizationId"] = client.Department
	wait := incrementalWait(3*time.Second, 5*time.Second)
	var raw interface{}
	var err error
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		args := *request
		raw, err = client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.CreateCustomerGateway(&args)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{errmsgs.Throttling, "OperationConflict"}) {
				wait()
				return resource.RetryableError(err)
			}

			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		errmsg := ""
		if raw != nil {
			response, ok := raw.(*vpc.CreateCustomerGatewayResponse)
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_vpn_customer_gateway", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	response, _ := raw.(*vpc.CreateCustomerGatewayResponse)

	d.SetId(response.CustomerGatewayId)

	err = vpnGatewayService.WaitForVpnCustomerGateway(d.Id(), Null, 60)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	return resourceAlibabacloudStackVpnCustomerGatewayRead(d, meta)
}

func resourceAlibabacloudStackVpnCustomerGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpnGatewayService := VpnGatewayService{client}

	object, err := vpnGatewayService.DescribeVpnCustomerGateway(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("ip_address", object.IpAddress)
	if err := connectivity.SetResourceData(d, object.Name, "customer_gateway_name", "name"); err != nil {
		return err
	}
	d.Set("description", object.Description)

	return nil
}

func resourceAlibabacloudStackVpnCustomerGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := vpc.CreateModifyCustomerGatewayAttributeRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.CustomerGatewayId = d.Id()
	if d.HasChange("customer_gateway_name") || d.HasChange("name") {
		if v, err := connectivity.GetResourceData(d, reflect.TypeOf(""), "customer_gateway_name", "name"); err == nil {
			request.Name = v.(string)
		} else {
			return err
		}
	}

	if d.HasChange("description") {
		request.Description = d.Get("description").(string)
	}

	request.Headers["x-acs-organizationId"] = client.Department
	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.ModifyCustomerGatewayAttribute(request)
	})
	if err != nil {
		errmsg := ""
		if raw != nil {
			response, ok := raw.(*vpc.ModifyCustomerGatewayAttributeResponse)
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return resourceAlibabacloudStackVpnCustomerGatewayRead(d, meta)
}

func resourceAlibabacloudStackVpnCustomerGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpnGatewayService := VpnGatewayService{client}
	request := vpc.CreateDeleteCustomerGatewayRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.CustomerGatewayId = d.Id()
	request.ClientToken = buildClientToken(request.GetActionName())

	request.Headers["x-acs-organizationId"] = client.Department
	err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		args := *request
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteCustomerGateway(&args)
		})

		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"VpnGateway.Configuring"}) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(err)
			}
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*vpc.DeleteCustomerGatewayResponse)
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
		if errmsgs.IsExpectedErrors(err, []string{"InvalidCustomerGatewayInstanceId.NotFound"}) {
			return nil
		}
		return err
	}
	return errmsgs.WrapError(vpnGatewayService.WaitForVpnCustomerGateway(d.Id(), Deleted, DefaultTimeout))
}
