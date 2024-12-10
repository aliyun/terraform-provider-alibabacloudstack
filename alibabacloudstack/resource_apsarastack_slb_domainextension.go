package alibabacloudstack

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackSlbDomainExtension() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackSlbDomainExtensionCreate,
		Read:   resourceAlibabacloudStackSlbDomainExtensionRead,
		Update: resourceAlibabacloudStackSlbDomainExtensionUpdate,
		Delete: resourceAlibabacloudStackSlbDomainExtensionDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"load_balancer_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"frontend_port": {
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntBetween(1, 65535),
				Optional:true,
				Computed:true,
				ForceNew:     true,
				Deprecated:   "Field 'frontend_port' is deprecated and will be removed in a future release. Please use new field 'listener_port' instead.",
				ConflictsWith: []string{"listener_port"},
			},
			"listener_port": {
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntBetween(1, 65535),
				Optional:true,
				Computed:true,
				ForceNew:     true,
				ConflictsWith: []string{"frontend_port"},
			},
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"server_certificate_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"id": {
				Type:         schema.TypeString,
				Computed:     true,
				Deprecated:   "Field 'id' is deprecated and will be removed in a future release. Please use new field 'domain_extension_id' instead.",
			},
			"domain_extension_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"delete_protection_validation": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceAlibabacloudStackSlbDomainExtensionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := slb.CreateCreateDomainExtensionRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.LoadBalancerId = d.Get("load_balancer_id").(string)
	request.ListenerPort = requests.NewInteger(connectivity.GetResourceData(d, "listener_port", "frontend_port").(int))
	if err := errmsgs.CheckEmpty(request.ListenerPort, schema.TypeString, "listener_port", "frontend_port"); err != nil {
		return errmsgs.WrapError(err)
	}
	request.Domain = d.Get("domain").(string)
	request.ServerCertificateId = d.Get("server_certificate_id").(string)

	var response *slb.CreateDomainExtensionResponse
	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.CreateDomainExtension(request)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"DomainExtensionProcessing"}) {
				return resource.RetryableError(err)
			}
			bresponse, ok := raw.(*slb.CreateDomainExtensionResponse)
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_slb_domain_extension", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response = raw.(*slb.CreateDomainExtensionResponse)
		return nil
	})
	if err != nil {
		return err
	}

	d.SetId(response.DomainExtensionId)
	return resourceAlibabacloudStackSlbDomainExtensionRead(d, meta)
}

func resourceAlibabacloudStackSlbDomainExtensionRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	slbService := SlbService{client}

	domainExtension, err := slbService.DescribeDomainExtensionAttribute(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	connectivity.SetResourceData(d, domainExtension.DomainExtensionId, "domain_extension_id", "id")
	d.Set("load_balancer_id", domainExtension.LoadBalancerId)
	d.Set("domain", domainExtension.Domain)
	d.Set("server_certificate_id", domainExtension.ServerCertificateId)
	connectivity.SetResourceData(d, domainExtension.ListenerPort, "listener_port", "frontend_port")

	return nil
}

func resourceAlibabacloudStackSlbDomainExtensionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	if d.HasChange("server_certificate_id") {
		request := slb.CreateSetDomainExtensionAttributeRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.DomainExtensionId = d.Id()
		request.ServerCertificateId = d.Get("server_certificate_id").(string)

		err := resource.Retry(3*time.Minute, func() *resource.RetryError {
			raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
				return slbClient.SetDomainExtensionAttribute(request)
			})
			if err != nil {
				if errmsgs.IsExpectedErrors(err, []string{"BackendServer.configuring", "DomainExtensionProcessing"}) {
					return resource.RetryableError(err)
				}
				bresponse, ok := raw.(*slb.SetDomainExtensionAttributeResponse)
				errmsg := ""
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
				}
				return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_slb_domain_extension", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		})
		if err != nil {
			return err
		}
		//d.SetPartial("server_certificate_id")
	}
	return resourceAlibabacloudStackSlbDomainExtensionRead(d, meta)
}

func resourceAlibabacloudStackSlbDomainExtensionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	slbService := SlbService{client}

	if d.Get("delete_protection_validation").(bool) {
		lbId := d.Get("load_balancer_id").(string)
		lbInstance, err := slbService.DescribeSlb(lbId)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				return nil
			}
			return errmsgs.WrapError(err)
		}
		if lbInstance.DeleteProtection == "on" {
			return errmsgs.WrapError(fmt.Errorf("Current domain extension's SLB Instance %s has enabled DeleteProtection. Please set delete_protection_validation to false to delete the resource.", lbId))
		}
	}

	request := slb.CreateDeleteDomainExtensionRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.DomainExtensionId = d.Id()

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.DeleteDomainExtension(request)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"DomainExtensionProcessing", "InternalError"}) {
				return resource.RetryableError(err)
			}
			bresponse, ok := raw.(*slb.DeleteDomainExtensionResponse)
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_slb_domain_extension", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidParameter.DomainExtensionId"}) {
			return nil
		}
		return err
	}
	return errmsgs.WrapError(slbService.WaitForSlbDomainExtension(d.Id(), Deleted, DefaultTimeoutMedium))
}
