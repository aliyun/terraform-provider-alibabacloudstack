package alibabacloudstack

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackSlbCACertificate() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Deprecated:   "Field 'name' is deprecated and will be removed in a future release. Please use new field 'ca_certificate_name' instead.",
				ConflictsWith: []string{"ca_certificate_name"},
			},
			"ca_certificate_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ConflictsWith: []string{"name"},
			},
			"ca_certificate": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackSlbCACertificateCreate,
		resourceAlibabacloudStackSlbCACertificateRead, resourceAlibabacloudStackSlbCACertificateUpdate, resourceAlibabacloudStackSlbCACertificateDelete)
	return resource
}

func resourceAlibabacloudStackSlbCACertificateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	slbService := SlbService{client}

	request := slb.CreateUploadCACertificateRequest()
	client.InitRpcRequest(*request.RpcRequest)
	
	if val, ok := connectivity.GetResourceDataOk(d, "ca_certificate_name", "name"); ok && val.(string) != "" {
		request.CACertificateName = val.(string)
	}

	if val, ok := d.GetOk("ca_certificate"); ok && val.(string) != "" {
		request.CACertificate = val.(string)
	} else {
		return errmsgs.WrapError(errmsgs.Error("UploadCACertificate got an error, ca_certificate should be not null"))
	}

	raw, err := slbService.client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.UploadCACertificate(request)
	})
	response, ok := raw.(*slb.UploadCACertificateResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_slb_ca_certificate", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	d.SetId(response.CACertificateId)

	return nil
}

func resourceAlibabacloudStackSlbCACertificateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	slbService := SlbService{client}

	object, err := slbService.DescribeSlbCACertificate(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			//return nil
		}
		return errmsgs.WrapError(err)
	}

	err = connectivity.SetResourceData(d, object.CACertificateName, "ca_certificate_name", "name")
	if err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}

func resourceAlibabacloudStackSlbCACertificateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	if d.HasChanges("name", "ca_certificate_name") {
		request := slb.CreateSetCACertificateNameRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.CACertificateId = d.Id()
		request.CACertificateName = connectivity.GetResourceData(d, "ca_certificate_name", "name").(string)

		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.SetCACertificateName(request)
		})
		response, ok := raw.(*slb.SetCACertificateNameResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	return nil
}

func resourceAlibabacloudStackSlbCACertificateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	slbService := SlbService{client}

	request := slb.CreateDeleteCACertificateRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.CACertificateId = d.Id()

	err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.DeleteCACertificate(request)
		})
		response, ok := raw.(*slb.DeleteCACertificateResponse)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, errmsgs.SlbIsBusy) {
				return resource.RetryableError(err)
			}
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"CACertificateId.NotFound"}) {
			return nil
		}
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return errmsgs.WrapError(slbService.WaitForSlbCACertificate(d.Id(), Deleted, DefaultTimeoutMedium))
}
