package apsarastack

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strconv"
	"strings"
)

func resourceApsaraStackDnsDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackDnsDomainCreate,
		Read:   resourceApsaraStackDnsDomainRead,
		Update: resourceApsaraStackDnsDomainUpdate,
		Delete: resourceApsaraStackDnsDomainDelete,
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

func resourceApsaraStackDnsDomainCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	dnsService := DnsService{client}
	var requestInfo *ecs.Client
	DomainName := d.Get("domain_name").(string)
	check, err := dnsService.DescribeDnsDomain(DomainName)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_dns_domain", "domain alreadyExist", ApsaraStackSdkGoERROR)
	}
	if len(check.ZoneList) == 0 {

		request := requests.NewCommonRequest()
		request.Method = "POST"        // Set request method
		request.Product = "GenesisDns" // Specify product
		request.Domain = client.Domain // Location Service will not be enabled if the host is specified. For example, service with a Certification type-Bearer Token should be specified
		request.Version = "2018-07-20" // Specify product version
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		request.ApiName = "AddGlobalAuthZone"
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.QueryParams = map[string]string{
			"AccessKeySecret": client.SecretKey,
			"AccessKeyId":     client.AccessKey,
			"Product":         "GenesisDns",
			"RegionId":        client.RegionId,
			"Action":          "AddGlobalAuthZone",
			"Version":         "2018-07-20",
			"DomainName":      DomainName,
		}
		raw, err := client.WithEcsClient(func(dnsClient *ecs.Client) (interface{}, error) {
			return dnsClient.ProcessCommonRequest(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_dns_domain", request.GetActionName(), ApsaraStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
		bresponse, _ := raw.(*responses.CommonResponse)
		if bresponse.GetHttpStatus() != 200 {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_dns_domain", "AddGlobalAuthZone", ApsaraStackSdkGoERROR)
		}
		addDebug("AddGlobalAuthZone", raw, requestInfo, bresponse.GetHttpContentString())
	}
	//err = resource.Retry(5*time.Minute, func() *resource.RetryError {
	check, err = dnsService.DescribeDnsDomain(DomainName)
	//if err != nil {
	//	return resource.NonRetryableError(err)
	//}
	//return resource.RetryableError(err)
	//})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_dns_domain", "DescribeDnsDomain")
	}
	//id := strconv.Itoa(dnsresp.ID)
	//d.SetId(id)
	d.SetId(check.ZoneList[0].DomainName + COLON_SEPARATED + fmt.Sprint(check.ZoneList[0].DomainID))

	return resourceApsaraStackDnsDomainUpdate(d, meta)
}
func resourceApsaraStackDnsDomainRead(d *schema.ResourceData, meta interface{}) error {
	wiatSecondsIfWithTest(1)
	client := meta.(*connectivity.ApsaraStackClient)
	dnsService := DnsService{client}
	object, err := dnsService.DescribeDnsDomain(d.Id())
	did := strings.Split(d.Id(), COLON_SEPARATED)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
		}
		return WrapError(err)
	}

	d.Set("domain_name", did[0])
	d.Set("domain_id", strconv.Itoa(object.ZoneList[0].DomainID))
	d.Set("remark", object.ZoneList[0].Remark)
	return nil
}
func resourceApsaraStackDnsDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	dnsService := DnsService{client}
	remarkUpdate := false
	check, err := dnsService.DescribeDnsDomain(d.Id())
	did := strings.Split(d.Id(), COLON_SEPARATED)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "IsDomainExist", ApsaraStackSdkGoERROR)
	}

	var desc string

	if d.HasChange("remark") {
		if v, ok := d.GetOk("remark"); ok {
			desc = v.(string)
		}
		check.ZoneList[0].Remark = desc
		remarkUpdate = true
	} else {
		if v, ok := d.GetOk("remark"); ok {
			desc = v.(string)
		}
		check.ZoneList[0].Remark = desc
	}
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Product = "GenesisDns"
	request.Domain = client.Domain
	request.Version = "2018-07-20"
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "RemarkGlobalAuthZone"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.RegionId = client.RegionId

	request.QueryParams = map[string]string{
		"AccessKeySecret": client.SecretKey,
		"AccessKeyId":     client.AccessKey,
		"Product":         "GenesisDns",
		"RegionId":        client.RegionId,
		"Action":          "RemarkGlobalAuthZone",
		"Version":         "2018-07-20",
		"Id":              did[1],
		"Remark":          desc,
	}
	if remarkUpdate {

		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		log.Printf(" response of raw RemarkGlobalAuthZone : %s", raw)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_dns_domain", "RemarkGlobalAuthZone", raw)
		}
		addDebug(request.GetActionName(), raw, request)
	}
	d.SetId(check.ZoneList[0].DomainName + COLON_SEPARATED + fmt.Sprint(check.ZoneList[0].DomainID))
	return resourceApsaraStackDnsDomainRead(d, meta)
}
func resourceApsaraStackDnsDomainDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	dnsService := DnsService{client}
	var requestInfo *ecs.Client
	check, err := dnsService.DescribeDnsDomain(d.Id())
	did := strings.Split(d.Id(), COLON_SEPARATED)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "IsDomainExist", ApsaraStackSdkGoERROR)
	}
	addDebug("IsDomainExist", check, requestInfo, map[string]string{"Id": did[1]})

	if len(check.ZoneList) != 0 {
		request := requests.NewCommonRequest()
		request.Method = "POST"        // Set request method
		request.Product = "GenesisDns" // Specify product
		request.Domain = client.Domain // Location Service will not be enabled if the host is specified. For example, service with a Certification type-Bearer Token should be specified
		request.Version = "2018-07-20" // Specify product version
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		request.ApiName = "DeleteGlobalZone"
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.QueryParams = map[string]string{
			"AccessKeySecret": client.SecretKey,
			"AccessKeyId":     client.AccessKey,
			"Product":         "GenesisDns",
			"RegionId":        client.RegionId,
			"Action":          "DeleteGlobalZone",
			"Version":         "2018-07-20",
			"Id":              did[1],
		}
		_, err := client.WithEcsClient(func(csClient *ecs.Client) (interface{}, error) {
			return csClient.ProcessCommonRequest(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_dns_domain", "DeleteGlobalZone", ApsaraStackSdkGoERROR)
		}
	}

	return nil
}
