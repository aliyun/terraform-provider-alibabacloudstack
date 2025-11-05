package alibabacloudstack

import (
	"fmt"
	"strings"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackAPIGatewayV2Certificate() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cert_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"0", "1"}, false),
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"certificates": {
				Type:     schema.TypeString,
				Required: true,
			},
			"private_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"certificate_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"certificate_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expire_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackAPIGatewayV2CertificateCreate, resourceAlibabacloudStackAPIGatewayV2CertificateRead, resourceAlibabacloudStackAPIGatewayV2CertificateUpdate, resourceAlibabacloudStackAPIGatewayV2CertificateDelete)
	return resource
}

func resourceAlibabacloudStackAPIGatewayV2CertificateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := map[string]interface{}{
		"certType":        d.Get("cert_type"),
		"certificates":    d.Get("certificates"),
		"certificateName": d.Get("certificate_name"),
		"gwInstanceId":    d.Get("instance_id"),
		"organizationId":  client.Department,
		"resourceGroupId": client.ResourceGroup,
		"regionId":        client.RegionId,
	}
	if v, ok := d.GetOk("private_key"); ok {
		request["privateKey"] = v
	}
	response, err := client.DoTeaRequest("POST", "csb2", "2023-02-06", "CreateCertificate", "/certificate/createCertificate", nil, nil, request)
	if err != nil {
		return err
	}
	if fmt.Sprint(response["code"]) != "200" {
		return errmsgs.Error("ModifyCertificate Failed! %v", response)
	}
	certificateId, ok := response["data"]
	if !ok {
		return errmsgs.Error("CreateCertificate Failed! %v", response)
	}
	id := fmt.Sprintf("%s:%s", d.Get("instance_id").(string), certificateId)
	d.SetId(id)
	return nil
}

func resourceAlibabacloudStackAPIGatewayV2CertificateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	apigatewayv2Service := ApiGateWayV2Service{client}
	certificate, err := apigatewayv2Service.DescribeApiGatewayV2Certificate(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	params := strings.Split(d.Id(), ":")
	d.Set("certificate_id", certificate["certificateId"])
	d.Set("instance_id", params[0])
	d.Set("certificate_name", certificate["certificateName"])
	d.Set("cert_type", certificate["certType"])
	d.Set("expire_time", certificate["expireTime"])
	d.Set("create_time", certificate["createTime"])
	d.Set("update_time", certificate["updateTime"])
	return nil
}

func resourceAlibabacloudStackAPIGatewayV2CertificateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	if d.IsNewResource() {
		return nil
	}
	if d.HasChanges("certificates", "private_key") {
		params := strings.Split(d.Id(), ":")
		request := map[string]interface{}{
			"gwInstanceId":  params[0],
			"certificateId": params[1],
			"certType":      d.Get("cert_type"),
			"certificates":  d.Get("certificates"),
			"privateKey":    d.Get("private_key"),
		}
		response, err := client.DoTeaRequest("POST", "csb2", "2023-02-06", "ModifyCertificate", "/certificate/modifyCertificate", nil, nil, request)
		if err != nil {
			return err
		}
		if fmt.Sprint(response["code"]) != "200" {
			return errmsgs.Error("ModifyCertificate Failed! %v", response)
		}
	}
	return nil
}

func resourceAlibabacloudStackAPIGatewayV2CertificateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	params := strings.Split(d.Id(), ":")
	request := map[string]interface{}{
		"gwInstanceId":  params[0],
		"certificateId": params[1],
	}
	_, err := client.DoTeaRequest("POST", "csb2", "2023-02-06", "DeleteCertificate", "/certificate/deleteCertificate", nil, nil, request)
	if err != nil {
		return err
	}
	return nil
}
