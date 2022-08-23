package apsarastack

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alibabacloudstack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApsaraStackSlbDomainExtension() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackSlbDomainExtensionCreate,
		Read:   resourceApsaraStackSlbDomainExtensionRead,
		Update: resourceApsaraStackSlbDomainExtensionUpdate,
		Delete: resourceApsaraStackSlbDomainExtensionDelete,

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
				Required:     true,
				ForceNew:     true,
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
			"delete_protection_validation": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceApsaraStackSlbDomainExtensionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)

	request := slb.CreateCreateDomainExtensionRequest()
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "slb", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.LoadBalancerId = d.Get("load_balancer_id").(string)
	request.ListenerPort = requests.NewInteger(d.Get("frontend_port").(int))
	request.Domain = d.Get("domain").(string)
	request.ServerCertificateId = d.Get("server_certificate_id").(string)

	var response *slb.CreateDomainExtensionResponse
	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.CreateDomainExtension(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"DomainExtensionProcessing"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response = raw.(*slb.CreateDomainExtensionResponse)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_slb_domain_extension", request.GetActionName(), ApsaraStackSdkGoERROR)
	}

	d.SetId(response.DomainExtensionId)
	return resourceApsaraStackSlbDomainExtensionRead(d, meta)
}

func resourceApsaraStackSlbDomainExtensionRead(d *schema.ResourceData, meta interface{}) error {
	wiatSecondsIfWithTest(1)
	client := meta.(*connectivity.ApsaraStackClient)
	slbService := SlbService{client}

	domainExtension, err := slbService.DescribeDomainExtensionAttribute(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("id", domainExtension.DomainExtensionId)
	d.Set("load_balancer_id", domainExtension.LoadBalancerId)
	d.Set("domain", domainExtension.Domain)
	d.Set("server_certificate_id", domainExtension.ServerCertificateId)
	d.Set("frontend_port", domainExtension.ListenerPort)

	return nil
}

func resourceApsaraStackSlbDomainExtensionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	if d.HasChange("server_certificate_id") {
		request := slb.CreateSetDomainExtensionAttributeRequest()
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "slb", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
		request.DomainExtensionId = d.Id()
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		request.ServerCertificateId = d.Get("server_certificate_id").(string)
		client := meta.(*connectivity.ApsaraStackClient)
		err := resource.Retry(3*time.Minute, func() *resource.RetryError {
			raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
				return slbClient.SetDomainExtensionAttribute(request)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{"BackendServer.configuring", "DomainExtensionProcessing"}) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_slb_domain_extension", request.GetActionName(), ApsaraStackSdkGoERROR)
		}
		//d.SetPartial("server_certificate_id")
	}
	return resourceApsaraStackSlbDomainExtensionRead(d, meta)
}

func resourceApsaraStackSlbDomainExtensionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	slbService := SlbService{client}

	if d.Get("delete_protection_validation").(bool) {
		lbId := d.Get("load_balancer_id").(string)
		lbInstance, err := slbService.DescribeSlb(lbId)
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return WrapError(err)
		}
		if lbInstance.DeleteProtection == "on" {
			return WrapError(fmt.Errorf("Current domain extension's SLB Instance %s has enabled DeleteProtection. Please set delete_protection_validation to false to delete the resource.", lbId))
		}
	}

	request := slb.CreateDeleteDomainExtensionRequest()
	request.DomainExtensionId = d.Id()
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "slb", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.DeleteDomainExtension(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"DomainExtensionProcessing", "InternalError"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidParameter.DomainExtensionId"}) {
			return nil
		}
		return WrapErrorf(err, DataDefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)

	}
	return WrapError(slbService.WaitForSlbDomainExtension(d.Id(), Deleted, DefaultTimeoutMedium))
}
