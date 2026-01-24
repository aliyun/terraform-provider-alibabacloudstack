package alibabacloudstack

import (
	"fmt"
	"log"
	"strings"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackDnsDomain() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"remark": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackDnsDomainCreate, resourceAlibabacloudStackDnsDomainRead, resourceAlibabacloudStackDnsDomainUpdate, resourceAlibabacloudStackDnsDomainDelete)
	return resource
}

func resourceAlibabacloudStackDnsDomainCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	dnsService := DnsService{client}
	DomainName := d.Get("domain_name").(string)
	check, err := dnsService.DescribeDnsDomain(DomainName)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_dns_domain", "domain alreadyExist", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if check == nil || len(check.Data) == 0 {
		request := client.NewCommonRequest("POST", "CloudDns", "2021-06-24", "AddGlobalZone", "")
		request.QueryParams["Name"] = DomainName
		bresponse, err := client.ProcessCommonRequest(request)
		if err != nil {
			if bresponse == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_dns_domain", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), bresponse)
		if bresponse.GetHttpStatus() != 200 {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_dns_domain", "AddGlobalZone", errmsgs.AlibabacloudStackSdkGoERROR)
		}
	}
	check, err = dnsService.DescribeDnsDomain(DomainName)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_dns_domain", "DescribeDnsDomain")
	}
	d.SetId(check.Data[0].Name + COLON_SEPARATED + fmt.Sprint(check.Data[0].Id))
	return nil
}

func resourceAlibabacloudStackDnsDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	dnsService := DnsService{client}
	object, err := dnsService.DescribeDnsDomain(d.Id())
	did := strings.Split(d.Id(), COLON_SEPARATED)
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
		}
		return errmsgs.WrapError(err)
	}

	d.Set("domain_name", did[0])
	d.Set("domain_id", (object.Data[0].Id))
	d.Set("remark", object.Data[0].Remark)
	return nil
}

func resourceAlibabacloudStackDnsDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	did := strings.Split(d.Id(), COLON_SEPARATED)

	if d.HasChange("remark") {
		request := client.NewCommonRequest("POST", "CloudDns", "2021-06-24", "UpdateGlobalZoneRemark", "")
		request.QueryParams["Name"] = did[0]
		request.QueryParams["Id"] = did[1]
		request.QueryParams["Remark"] = d.Get("remark").(string)
		response, err := client.ProcessCommonRequest(request)
		log.Printf(" response of raw UpdateGlobalZoneRemark : %s", response)

		if err != nil {
			if response == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_dns_domain", "UpdateGlobalZoneRemark", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), response, request)
	}
	return nil
}

func resourceAlibabacloudStackDnsDomainDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	dnsService := DnsService{client}
	check, err := dnsService.DescribeDnsDomain(d.Id())
	did := strings.Split(d.Id(), COLON_SEPARATED)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "IsDomainExist", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	if len(check.Data) != 0 {
		request := client.NewCommonRequest("POST", "CloudDns", "2021-06-24", "DeleteGlobalZone", "")
		request.QueryParams["Id"] = did[1]
		response, err := client.ProcessCommonRequest(request)
		if err != nil {
			if response == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_dns_domain", "DeleteGlobalZone", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
	}

	return nil
}
