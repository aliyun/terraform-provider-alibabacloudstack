package alibabacloudstack

import (
	"strings"
	"time"
	"reflect"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackSlbServerCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackSlbServerCertificateCreate,
		Read:   resourceAlibabacloudStackSlbServerCertificateRead,
		Update: resourceAlibabacloudStackSlbServerCertificateUpdate,
		Delete: resourceAlibabacloudStackSlbServerCertificateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Deprecated:   "Field 'name' is deprecated and will be removed in a future release. Please use 'server_certificate_name' instead.",
			},
			"server_certificate_name": {
				Type:         schema.TypeString,
				Optional:     true,
			},
			"server_certificate": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"private_key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlibabacloudStackSlbServerCertificateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := slb.CreateUploadServerCertificateRequest()
	client.InitRpcRequest(*request.RpcRequest)

	if v, err := connectivity.GetResourceData(d, reflect.TypeOf(""), "server_certificate_name", "name"); err == nil {
		request.ServerCertificateName = v.(string)
	} else {
		return err
	}

	if val, ok := d.GetOk("server_certificate"); ok && val != "" {
		request.ServerCertificate = val.(string)
	}

	if val, ok := d.GetOk("private_key"); ok && val != "" {
		request.PrivateKey = val.(string)
	}
	//check server_certificate and private_key
	if request.AliCloudCertificateId == "" {
		if val := strings.Trim(request.ServerCertificate, " "); val == "" {
			return errmsgs.WrapError(errmsgs.Error("UploadServerCertificate got an error, as server_certificate should be not null when alibabacloudstack_certificate_id is null."))
		}

		if val := strings.Trim(request.PrivateKey, " "); val == "" {
			return errmsgs.WrapError(errmsgs.Error("UploadServerCertificate got an error, as either private_key or private_file  should be not null when alibabacloudstack_certificate_id is null."))
		}
	}

	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.UploadServerCertificate(request)
	})
	if err != nil {
		errmsg := ""
		if response, ok := raw.(*slb.UploadServerCertificateResponse); ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*slb.UploadServerCertificateResponse)
	d.SetId(response.ServerCertificateId)

	return resourceAlibabacloudStackSlbServerCertificateUpdate(d, meta)
}

func resourceAlibabacloudStackSlbServerCertificateRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	slbService := SlbService{client}
	serverCertificate, err := slbService.DescribeSlbServerCertificate(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	if err := connectivity.SetResourceData(d, serverCertificate.ServerCertificateName, "server_certificate_name", "name"); err != nil {
		return errmsgs.WrapError(err)
	}
	return nil
}

func resourceAlibabacloudStackSlbServerCertificateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	if d.IsNewResource() {
		d.Partial(false)
		return resourceAlibabacloudStackSlbServerCertificateRead(d, meta)
	}
	if !d.IsNewResource() && d.HasChange("name") {
		request := slb.CreateSetServerCertificateNameRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.ServerCertificateId = d.Id()

		if v, err := connectivity.GetResourceData(d, reflect.TypeOf(""), "server_certificate_name", "name"); err == nil {
			request.ServerCertificateName = v.(string)
		} else {
			return err
		}

		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.SetServerCertificateName(request)
		})
		if err != nil {
			errmsg := ""
			if response, ok := raw.(*slb.SetServerCertificateNameResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}
	return resourceAlibabacloudStackSlbServerCertificateRead(d, meta)
}

func resourceAlibabacloudStackSlbServerCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	slbService := SlbService{client}

	request := slb.CreateDeleteServerCertificateRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.ServerCertificateId = d.Id()

	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.DeleteServerCertificate(request)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"CertificateAndPrivateKeyIsRefered"}) {
				return resource.RetryableError(err)
			}
			errmsg := ""
			if response, ok := raw.(*slb.DeleteServerCertificateResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"ServerCertificateId.NotFound"}) {
			return nil
		}
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return errmsgs.WrapError(slbService.WaitForSlbServerCertificate(d.Id(), Deleted, DefaultTimeoutMedium))
}
