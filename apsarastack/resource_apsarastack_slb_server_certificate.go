package apsarastack

import (
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApsaraStackSlbServerCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackSlbServerCertificateCreate,
		Read:   resourceApsaraStackSlbServerCertificateRead,
		Update: resourceApsaraStackSlbServerCertificateUpdate,
		Delete: resourceApsaraStackSlbServerCertificateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
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

func resourceApsaraStackSlbServerCertificateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)

	request := slb.CreateUploadServerCertificateRequest()
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "slb", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

	if val, ok := d.GetOk("name"); ok && val != "" {
		request.ServerCertificateName = val.(string)
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
			return WrapError(Error("UploadServerCertificate got an error, as server_certificate should be not null when apsarastack_certificate_id is null."))
		}

		if val := strings.Trim(request.PrivateKey, " "); val == "" {
			return WrapError(Error("UploadServerCertificate got an error, as either private_key or private_file  should be not null when apsarastack_certificate_id is null."))
		}
	}
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "slb", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.UploadServerCertificate(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*slb.UploadServerCertificateResponse)
	d.SetId(response.ServerCertificateId)

	return resourceApsaraStackSlbServerCertificateUpdate(d, meta)
}

func resourceApsaraStackSlbServerCertificateRead(d *schema.ResourceData, meta interface{}) error {
	wiatSecondsIfWithTest(1)
	client := meta.(*connectivity.ApsaraStackClient)
	slbService := SlbService{client}
	serverCertificate, err := slbService.DescribeSlbServerCertificate(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if err := d.Set("name", serverCertificate.ServerCertificateName); err != nil {
		return WrapError(err)
	}
	return nil
}

func resourceApsaraStackSlbServerCertificateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	if d.IsNewResource() {
		d.Partial(false)
		return resourceApsaraStackSlbServerCertificateRead(d, meta)
	}
	if !d.IsNewResource() && d.HasChange("name") {
		request := slb.CreateSetServerCertificateNameRequest()
		request.RegionId = client.RegionId
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "slb", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
		request.ServerCertificateId = d.Id()
		request.ServerCertificateName = d.Get("name").(string)
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.SetServerCertificateName(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}
	return resourceApsaraStackSlbServerCertificateRead(d, meta)
}

func resourceApsaraStackSlbServerCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	slbService := SlbService{client}

	request := slb.CreateDeleteServerCertificateRequest()
	request.RegionId = client.RegionId
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "slb", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.ServerCertificateId = d.Id()
	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.DeleteServerCertificate(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"CertificateAndPrivateKeyIsRefered"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)

		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"ServerCertificateId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
	}

	return WrapError(slbService.WaitForSlbServerCertificate(d.Id(), Deleted, DefaultTimeoutMedium))

}
