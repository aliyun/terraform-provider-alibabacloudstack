package alibabacloudstack

import (
	"fmt"
	"log"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackDnsDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackDnsDomainCreate,
		Read:   resourceAlibabacloudStackDnsDomainRead,
		Update: resourceAlibabacloudStackDnsDomainUpdate,
		Delete: resourceAlibabacloudStackDnsDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"dns_servers": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"lang": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"remark": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlibabacloudStackDnsDomainCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	dnsService := DnsService{client}
	var requestInfo *ecs.Client
	DomainName := d.Get("domain_name").(string)
	check, err := dnsService.DescribeDnsDomain(DomainName)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_dns_domain", "domain alreadyExist", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if len(check.Data) == 0 {
		request := client.NewCommonRequest("POST", "CloudDns", "2021-06-24", "AddGlobalZone", "")
		request.QueryParams["Name"] = DomainName
		raw, err := client.WithEcsClient(func(dnsClient *ecs.Client) (interface{}, error) {
			return dnsClient.ProcessCommonRequest(request)
		})
		if err != nil {
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*responses.CommonResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_dns_domain", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw)
		bresponse, ok := raw.(*responses.CommonResponse)
		if !ok {
			return fmt.Errorf("Failed to cast raw response to CommonResponse")
		}
		if bresponse.GetHttpStatus() != 200 {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_dns_domain", "AddGlobalZone", errmsgs.AlibabacloudStackSdkGoERROR)
		}
		addDebug("AddGlobalZone", raw, requestInfo, bresponse.GetHttpContentString())
	}
	check, err = dnsService.DescribeDnsDomain(DomainName)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_dns_domain", "DescribeDnsDomain")
	}
	d.SetId(check.Data[0].Name + COLON_SEPARATED + fmt.Sprint(check.Data[0].Id))
	return resourceAlibabacloudStackDnsDomainUpdate(d, meta)
}

func resourceAlibabacloudStackDnsDomainRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
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
	dnsService := DnsService{client}
	remarkUpdate := false
	check, err := dnsService.DescribeDnsDomain(d.Id())
	did := strings.Split(d.Id(), COLON_SEPARATED)

	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "IsDomainExist", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	var desc string

	if d.HasChange("remark") {
		if v, ok := d.GetOk("remark"); ok {
			desc = v.(string)
		}
		check.Data[0].Remark = desc
		remarkUpdate = true
	} else {
		if v, ok := d.GetOk("remark"); ok {
			desc = v.(string)
		}
		check.Data[0].Remark = desc
	}

	if remarkUpdate {
		request := client.NewCommonRequest("POST", "CloudDns", "2021-06-24", "UpdateGlobalZoneRemark", "")
		request.QueryParams["Name"] = did[0]
		request.QueryParams["Id"] = did[1]
		request.QueryParams["Remark"] = desc
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		log.Printf(" response of raw UpdateGlobalZoneRemark : %s", raw)

		if err != nil {
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*responses.CommonResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_dns_domain", "UpdateGlobalZoneRemark", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request)
	}
	d.SetId(check.Data[0].Name + COLON_SEPARATED + fmt.Sprint(check.Data[0].Id))
	return resourceAlibabacloudStackDnsDomainRead(d, meta)
}

func resourceAlibabacloudStackDnsDomainDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	dnsService := DnsService{client}
	var requestInfo *ecs.Client
	check, err := dnsService.DescribeDnsDomain(d.Id())
	did := strings.Split(d.Id(), COLON_SEPARATED)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "IsDomainExist", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	addDebug("IsDomainExist", check, requestInfo, map[string]string{"Id": did[1]})

	if len(check.Data) != 0 {
		request := client.NewCommonRequest("POST", "CloudDns", "2021-06-24", "DeleteGlobalZone", "")
		request.QueryParams["Id"] = did[1]
		raw, err := client.WithEcsClient(func(csClient *ecs.Client) (interface{}, error) {
			return csClient.ProcessCommonRequest(request)
		})
		if err != nil {
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*responses.CommonResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_dns_domain", "DeleteGlobalZone", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
	}

	return nil
}
