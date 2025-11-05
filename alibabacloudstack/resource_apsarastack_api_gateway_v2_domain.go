package alibabacloudstack

import (
	"fmt"
	"strings"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackAPIGateWayV2Domain() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"HTTPS", "HTTP"}, false),
			},
			"certificate_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ca_certificate_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"client_auth": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "0",
				ValidateFunc: validation.StringInSlice([]string{"0", "1"}, false),
			},
			"subject_dn": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"issuer_dn": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackAPIGateWayV2DomainCreate, resourceAlibabacloudStackAPIGateWayV2DomainRead, resourceAlibabacloudStackAPIGateWayV2DomainUpdate, resourceAlibabacloudStackAPIGateWayV2DomainDelete)
	return resource
}

func resourceAlibabacloudStackAPIGateWayV2DomainCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	protocol := d.Get("protocol").(string)
	client_auth := d.Get("client_auth").(string)
	subject_dn := d.Get("subject_dn").(string)
	issuer_dn := d.Get("issuer_dn").(string)
	certificate_id := d.Get("certificate_id").(string)
	ca_certificate_id := d.Get("ca_certificate_id").(string)
	if protocol == "https" && certificate_id == "" {
		return errmsgs.WrapError(fmt.Errorf("[ERROR] The certificate_id is required when the protocol is https"))
	}
	if protocol == "https" && client_auth == "1" && ca_certificate_id == "" {
		return errmsgs.WrapError(fmt.Errorf("[ERROR] The ca_certificate_id is required when the protocol is https and client_auth is 1"))
	}

	request := map[string]interface{}{
		"gwInstanceId":    d.Get("instance_id"),
		"domain":          d.Get("domain"),
		"protocol":        protocol,
		"clientAuth":      client_auth,
		"isSubjectDn":     subject_dn != "",
		"isIssuerDn":      issuer_dn != "",
		"certificateId":   certificate_id,
		"caCertificateId": ca_certificate_id,
		"subjectDn":       subject_dn,
		"issuerDn":        issuer_dn,
	}
	response, err := client.DoTeaRequest("POST", "csb2", "2023-02-06", "CreateDomain", "/domain/createDomain", nil, nil, request)
	if err != nil {
		return err
	}
	if fmt.Sprint(response["code"]) != "200" {
		return errmsgs.Error("CreateDomain Failed! %v", response)
	}
	if domainId, ok := response["data"]; ok {
		id := fmt.Sprintf("%s:%s", d.Get("instance_id"), domainId.(string))
		d.SetId(id)
	} else {
		return errmsgs.Error("CreateDomain Failed! %v", response)
	}
	return nil
}

func resourceAlibabacloudStackAPIGateWayV2DomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	apigatewayv2Service := ApiGateWayV2Service{client}
	domain, err := apigatewayv2Service.DescribeApiGatewayV2Domain(d.Id())
	if err != nil {
		return errmsgs.WrapError(err)
	}
	params := strings.Split(d.Id(), ":")
	instanceId := params[0]
	d.Set("instance_id", instanceId)
	d.Set("domain", domain["domain"])
	d.Set("domain_id", domain["domainId"])
	d.Set("protocol", domain["protocol"])
	d.Set("certificate_id", domain["certificateId"])
	d.Set("client_auth", domain["clientAuth"])
	d.Set("ca_certificate_id", domain["caCertificateId"])
	d.Set("subject_dn", domain["subjectDn"])
	d.Set("issuer_dn", domain["issuerDn"])

	return nil
}

func resourceAlibabacloudStackAPIGateWayV2DomainUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.IsNewResource() {
		return nil
	}
	if d.HasChanges("protocol", "certificate_id", "client_auth", "ca_certificate_id", "subject_dn", "issuer_dn") {
		client := meta.(*connectivity.AlibabacloudStackClient)
		params := strings.Split(d.Id(), ":")
		protocol := d.Get("protocol").(string)
		client_auth := d.Get("client_auth").(string)
		subject_dn := d.Get("subject_dn").(string)
		issuer_dn := d.Get("issuer_dn").(string)
		certificate_id := d.Get("certificate_id").(string)
		ca_certificate_id := d.Get("ca_certificate_id").(string)
		if protocol == "https" && certificate_id == "" {
			return errmsgs.WrapError(fmt.Errorf("[ERROR] The certificate_id is required when the protocol is https"))
		}
		if protocol == "https" && client_auth == "1" && ca_certificate_id == "" {
			return errmsgs.WrapError(fmt.Errorf("[ERROR] The ca_certificate_id is required when the protocol is https and client_auth is 1"))
		}
		request := map[string]interface{}{
			"gwInstanceId":    params[0],
			"domainId":        params[1],
			"domain":          d.Get("domain").(string),
			"protocol":        protocol,
			"clientAuth":      client_auth,
			"isSubjectDn":     subject_dn != "",
			"isIssuerDn":      issuer_dn != "",
			"certificateId":   certificate_id,
			"caCertificateId": ca_certificate_id,
			"subjectDn":       subject_dn,
			"issuerDn":        issuer_dn,
		}
		response, err := client.DoTeaRequest("POST", "csb2", "2023-02-06", "ModifyDomain", "/domain/modifyDomain", nil, nil, request)
		if err != nil {
			return err
		}
		if fmt.Sprint(response["code"]) != "200" {
			return errmsgs.Error("ModifyDomain Failed! %v", response)
		}
	}

	return nil
}

func resourceAlibabacloudStackAPIGateWayV2DomainDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	params := strings.Split(d.Id(), ":")
	request := map[string]interface{}{
		"gwInstanceId": params[0],
		"domainId":     params[1],
	}

	_, err := client.DoTeaRequest("POST", "csb2", "2023-02-06", "DeleteDomain", "/domain/deleteDomain", nil, nil, request)
	if err != nil {
		return err
	}

	return nil
}
