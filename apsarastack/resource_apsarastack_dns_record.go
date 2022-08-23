package apsarastack

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"strings"
)

func resourceApsaraStackDnsRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackDnsRecordCreate,
		Read:   resourceApsaraStackDnsRecordRead,
		Update: resourceApsaraStackDnsRecordUpdate,
		Delete: resourceApsaraStackDnsRecordDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"record_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"host_record": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateRR,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"A", "NS", "MX", "TXT", "CNAME", "SRV", "AAAA", "CAA", "REDIRECT_URL", "FORWORD_URL"}, false),
			},
			"rr_set": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				MinItems: 1,
			},
			"ttl": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  300,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceApsaraStackDnsRecordCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	var requestInfo *ecs.Client
	DomainID := d.Get("domain_id").(string)
	RR := d.Get("host_record").(string)
	Type := d.Get("type").(string)
	var rrset string
	var rrsets []string
	if v, ok := d.GetOk("rr_set"); ok {
		rrsets = expandStringList(v.(*schema.Set).List())
		for i, k := range rrsets {
			if i != 0 {
				rrset = fmt.Sprintf("%s\",\"%s", rrset, k)
			} else {
				rrset = k
			}
		}
	}
	TTL := d.Get("ttl").(int)
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
	request.ApiName = "AddGlobalRrSet"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.RegionId = client.RegionId

	request.QueryParams = map[string]string{
		"AccessKeySecret": client.SecretKey,
		"AccessKeyId":     client.AccessKey,
		"Product":         "GenesisDns",
		"RegionId":        client.RegionId,
		"Action":          "AddGlobalRrSet",
		"Version":         "2018-07-20",
		"Type":            Type,
		"Ttl":             fmt.Sprint(TTL),
		"RrSet":           fmt.Sprintf("[\"%s\"]", rrset),
		"ZoneId":          DomainID,
		"Rr":              RR,
	}

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_dns_record", "AddGlobalRrSet", raw)
	}
	addDebug("AddGlobalRrSet", raw, requestInfo, request)

	bresponse, _ := raw.(*responses.CommonResponse)
	if bresponse.GetHttpStatus() != 200 {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_dns_record", "AddGlobalRrSet", ApsaraStackSdkGoERROR)
	}
	addDebug("AddGlobalRrSet", raw, requestInfo, bresponse.GetHttpContentString())

	d.SetId(RR + COLON_SEPARATED + DomainID)

	return resourceApsaraStackDnsRecordRead(d, meta)
}

func resourceApsaraStackDnsRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	dnsService := DnsService{client}
	RecordID := d.Get("record_id").(int)
	Rr := d.Get("host_record").(string)
	check, err := dnsService.DescribeDnsRecord(d.Id())
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "IsRecordExist", ApsaraStackSdkGoERROR)
	}
	attributeUpdate := false

	var desc string

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			desc = v.(string)
		}
		check.Records[0].Remark = desc
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
		request.ApiName = "RemarkGlobalRrSet"
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.RegionId = client.RegionId

		request.QueryParams = map[string]string{
			"AccessKeySecret": client.SecretKey,
			"AccessKeyId":     client.AccessKey,
			"Product":         "GenesisDns",
			"RegionId":        client.RegionId,
			"Action":          "RemarkGlobalRrSet",
			"Version":         "2018-07-20",
			"Id":              fmt.Sprint(RecordID),
			"Remark":          desc,
		}
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		log.Printf(" response of raw RemarkGlobalRrSet : %s", raw)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_dns_record", "RemarkGlobalRrSet", raw)
		}
		addDebug(request.GetActionName(), raw, request)
	} else {
		if v, ok := d.GetOk("description"); ok {
			desc = v.(string)
		}
		check.Records[0].Remark = desc
	}

	var Type string
	var Ttl int

	if d.HasChange("type") {
		if v, ok := d.GetOk("type"); ok {
			Type = v.(string)
		}
		check.Records[0].Type = Type
		attributeUpdate = true
	} else {
		if v, ok := d.GetOk("type"); ok {
			Type = v.(string)
		}
		check.Records[0].Type = Type
	}
	if d.HasChange("ttl") {
		if v, ok := d.GetOk("ttl"); ok {
			Ttl = v.(int)
		}
		check.Records[0].TTL = Ttl
		attributeUpdate = true
	} else {
		if v, ok := d.GetOk("ttl"); ok {
			Ttl = v.(int)
		}
		check.Records[0].TTL = Ttl
	}

	var rrset string
	var rrsets []string

	if d.HasChange("rr_set") {
		if v, ok := d.GetOk("rr_set"); ok {
			rrsets = expandStringList(v.(*schema.Set).List())

			for i, k := range rrsets {
				if i != 0 {
					rrset = fmt.Sprintf("%s\",\"%s", rrset, k)
				} else {
					rrset = k
				}
			}
			check.Records[0].RrSet = rrsets
		}
		attributeUpdate = true
	} else {
		if v, ok := d.GetOk("rr_set"); ok {
			rrsets = expandStringList(v.(*schema.Set).List())
			for i, k := range rrsets {
				if i != 0 {
					rrset = fmt.Sprintf("%s\",\"%s", rrset, k)
				} else {
					rrset = k
				}
			}
			check.Records[0].RrSet = rrsets
		}
	}

	if attributeUpdate {

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
		request.ApiName = "UpdateGlobalRrSet"
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.RegionId = client.RegionId

		request.QueryParams = map[string]string{
			"AccessKeySecret": client.SecretKey,
			"AccessKeyId":     client.AccessKey,
			"Product":         "GenesisDns",
			"RegionId":        client.RegionId,
			"Action":          "UpdateGlobalRrSet",
			"Version":         "2018-07-20",
			"RrsetId":         fmt.Sprint(RecordID),
			"RrSet":           fmt.Sprintf("[\"%s\"]", rrset),
			"Ttl":             fmt.Sprint(Ttl),
			"Type":            Type,
			"Rr":              Rr,
		}
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		log.Printf(" response of raw UpdateGlobalRrSet : %s", raw)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_dns_record", "UpdateGlobalRrSet", raw)
		}
		addDebug(request.GetActionName(), raw, request)

	}

	return resourceApsaraStackDnsRecordRead(d, meta)
}

func resourceApsaraStackDnsRecordRead(d *schema.ResourceData, meta interface{}) error {
	wiatSecondsIfWithTest(1)
	client := meta.(*connectivity.ApsaraStackClient)

	dnsService := &DnsService{client: client}
	object, err := dnsService.DescribeDnsRecord(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("ttl", object.Records[0].TTL)
	d.Set("record_id", object.Records[0].RecordID)
	d.Set("host_record", object.Records[0].Rr)
	d.Set("type", object.Records[0].Type)
	d.Set("description", object.Records[0].Remark)

	return nil
}

func resourceApsaraStackDnsRecordDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	RecordID := d.Get("record_id").(int)
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
	request.ApiName = "DeleteGlobalRrSet"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{
		"AccessKeySecret": client.SecretKey,
		"AccessKeyId":     client.AccessKey,
		"Product":         "GenesisDns",
		"RegionId":        client.RegionId,
		"Action":          "DeleteGlobalRrSet",
		"Version":         "2018-07-20",
		"Id":              fmt.Sprint(RecordID),
	}
	raw, err := client.WithEcsClient(func(dnsClient *ecs.Client) (interface{}, error) {
		return dnsClient.ProcessCommonRequest(request)
	})
	addDebug(request.GetActionName(), raw)

	if err != nil {
		if IsExpectedErrors(err, []string{"DomainRecordNotBelongToUser"}) {
			return nil
		}
		if IsExpectedErrors(err, []string{"RecordForbidden.DNSChange", "InternalError"}) {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
		}
	}
	return nil
}
