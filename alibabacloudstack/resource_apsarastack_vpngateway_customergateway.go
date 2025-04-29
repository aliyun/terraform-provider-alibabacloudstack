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

func resourceAlibabacloudStackVpnCustomerGateway() *schema.Resource {
	resource := &schema.Resource{
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
				Computed:true,
				ValidateFunc: validation.StringLenBetween(2, 128),
				Deprecated:   "Field 'name' is deprecated and will be removed in a future release. Please use new field 'customer_gateway_name' instead.",
				ConflictsWith: []string{"customer_gateway_name"},
			},
			"customer_gateway_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:true,
				ValidateFunc: validation.StringLenBetween(2, 128),
				ConflictsWith: []string{"name"},
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackVpnCustomerGatewayCreate, resourceAlibabacloudStackVpnCustomerGatewayRead, resourceAlibabacloudStackVpnCustomerGatewayUpdate, resourceAlibabacloudStackVpnCustomerGatewayDelete)
	return resource
}

func resourceAlibabacloudStackVpnCustomerGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpnGatewayService := VpnGatewayService{client}
	request := vpc.CreateCreateCustomerGatewayRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.IpAddress = d.Get("ip_address").(string)
	if v := connectivity.GetResourceData(d, "customer_gateway_name", "name").(string); v != "" {
		request.Name = v
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
	return nil
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
	if d.HasChanges("customer_gateway_name", "name") {
		request.Name = connectivity.GetResourceData(d, "customer_gateway_name", "name").(string)
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

	return nil
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