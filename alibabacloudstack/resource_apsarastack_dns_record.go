package alibabacloudstack

import (
	"fmt"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"strings"
	"time"
)

func resourceAlibabacloudStackDnsRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackDnsRecordCreate,
		Read:   resourceAlibabacloudStackDnsRecordRead,
		Update: resourceAlibabacloudStackDnsRecordUpdate,
		Delete: resourceAlibabacloudStackDnsRecordDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"record_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"A", "AAAA", "CNAME", "MX", "TXT", "PTR", "SRV", "NAPRT", "CAA", "NS"}, false),
			},
			"lba_strategy": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"ALL_RR", "RATIO"}, false),
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
			"remark": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlibabacloudStackDnsRecordCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	//var requestInfo *ecs.Client
	request := make(map[string]interface{})
	ZoneId := d.Get("zone_id").(int)
	LbaStrategy := d.Get("lba_strategy").(string)
	Type := d.Get("type").(string)
	Name := d.Get("name").(string)
	TTL := d.Get("ttl").(int)
	var rrsets []string
	if v, ok := d.GetOk("rr_set"); ok {
		rrsets = expandStringList(v.(*schema.Set).List())
		for i, key := range rrsets {
			request[fmt.Sprintf("RDatas.%d.Value", i+1)] = key

		}
	}
	var response map[string]interface{}

	action := "AddGlobalZoneRecord"
	request["Product"] = "CloudDns"
	request["product"] = "CloudDns"
	request["OrganizationId"] = client.Department
	request["RegionId"] = client.RegionId

	request["Type"] = Type
	request["Ttl"] = TTL
	request["ZoneId"] = ZoneId
	request["LbaStrategy"] = LbaStrategy
	request["Name"] = Name
	conn, err := client.NewCloudApiClient()
	if err != nil {
		return WrapError(err)
	}
	request["ClientToken"] = buildClientToken("AddGlobalZoneRecord")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequesttowpoint1(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-06-24"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	addDebug("AddGlobalZoneRecord", response, request)
	d.SetId(fmt.Sprint(ZoneId))

	return resourceAlibabacloudStackDnsRecordRead(d, meta)
}

func resourceAlibabacloudStackDnsRecordRead(d *schema.ResourceData, meta interface{}) error {
	wiatSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)

	dnsService := &DnsService{client: client}
	object, err := dnsService.DescribeDnsRecord(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("ttl", object.Data[0].TTL)
	d.Set("record_id", object.Data[0].Id)
	d.Set("name", object.Data[0].Name)
	d.Set("type", object.Data[0].Type)
	d.Set("remark", object.Data[0].Remark)
	d.Set("zone_id", object.Data[0].ZoneId)
	d.Set("lba_strategy", object.Data[0].LbaStrategy)

	return nil
}
func resourceAlibabacloudStackDnsRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	dnsService := DnsService{client}
	ID := d.Get("record_id").(int)
	ZoneId := d.Get("zone_id").(int)
	Name := d.Get("name").(string)
	LbaStrategy := d.Get("lba_strategy").(string)
	check, err := dnsService.DescribeDnsRecord(d.Id())
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "IsRecordExist", AlibabacloudStackSdkGoERROR)
	}
	attributeUpdate := false

	var desc string

	if d.HasChange("remark") {
		if v, ok := d.GetOk("remark"); ok {
			desc = v.(string)
		}
		check.Data[0].Remark = desc
		request := requests.NewCommonRequest()
		request.Method = "POST"
		request.Product = "CloudDns"
		request.Domain = client.Domain
		request.Version = "2021-06-24"
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		request.ApiName = "UpdateGlobalZoneRecordRemark"
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.RegionId = client.RegionId

		request.QueryParams = map[string]string{
			"AccessKeySecret": client.SecretKey,
			"AccessKeyId":     client.AccessKey,
			"Product":         "CloudDns",
			"RegionId":        client.RegionId,
			"Action":          "UpdateGlobalZoneRecordRemark",
			"Version":         "2021-06-24",
			"Id":              fmt.Sprint(ID),
			"Remark":          desc,
		}
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		log.Printf(" response of raw UpdateGlobalZoneRecordRemark : %s", raw)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_dns_record", "UpdateGlobalZoneRecordRemark", raw)
		}
		addDebug(request.GetActionName(), raw, request)
	} else {
		if v, ok := d.GetOk("remark"); ok {
			desc = v.(string)
		}
		check.Data[0].Remark = desc
	}

	var Type string
	var Ttl int

	if d.HasChange("type") {
		if v, ok := d.GetOk("type"); ok {
			Type = v.(string)
		}
		check.Data[0].Type = Type
		attributeUpdate = true
	} else {
		if v, ok := d.GetOk("type"); ok {
			Type = v.(string)
		}
		check.Data[0].Type = Type
	}
	if d.HasChange("ttl") {
		if v, ok := d.GetOk("ttl"); ok {
			Ttl = v.(int)
		}
		check.Data[0].TTL = Ttl
		attributeUpdate = true
	} else {
		if v, ok := d.GetOk("ttl"); ok {
			Ttl = v.(int)
		}
		check.Data[0].TTL = Ttl
	}

	//var rrset string

	//if v, ok := d.GetOk("rr_set"); ok {
	//	rrsets = expandStringList(v.(*schema.Set).List())
	//	for i, key := range rrsets {
	//		request[fmt.Sprintf("RDatas.%d.Value", i+1)] = key
	//
	//	}
	//}

	if d.HasChange("rr_set") {
		attributeUpdate = true
	}

	if attributeUpdate {
		request := make(map[string]interface{})
		var rrsets []string
		if v, ok := d.GetOk("rr_set"); ok {
			rrsets = expandStringList(v.(*schema.Set).List())
			for i, key := range rrsets {
				request[fmt.Sprintf("RDatas.%d.Value", i+1)] = key

			}
		}
		action := "UpdateGlobalZoneRecord"
		request["Product"] = "CloudDns"
		request["product"] = "CloudDns"
		request["OrganizationId"] = client.Department
		request["RegionId"] = client.RegionId
		request["Type"] = Type
		request["Ttl"] = Ttl
		request["Id"] = ID
		request["ZoneId"] = ZoneId
		request["LbaStrategy"] = LbaStrategy
		request["Name"] = Name
		request["Remark"] = check.Data[0].Remark
		conn, err := client.NewCloudApiClient()
		if err != nil {
			return WrapError(err)
		}
		var response map[string]interface{}
		request["ClientToken"] = buildClientToken("UpdateGlobalZoneRecord")
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			response, err = conn.DoRequesttowpoint1(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-06-24"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		addDebug("UpdateGlobalZoneRecord", response, request)

	}

	return resourceAlibabacloudStackDnsRecordRead(d, meta)
}

func resourceAlibabacloudStackDnsRecordDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ID := d.Get("record_id").(int)
	ZoneId := d.Get("zone_id").(int)
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Product = "CloudDns"
	request.Domain = client.Domain
	request.Version = "2021-06-24"
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "DeleteGlobalZoneRecord"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{
		"AccessKeySecret": client.SecretKey,
		"AccessKeyId":     client.AccessKey,
		"Product":         "CloudDns",
		"RegionId":        client.RegionId,
		"Action":          "DeleteGlobalZoneRecord",
		"Version":         "2021-06-24",
		"Id":              fmt.Sprint(ID),
		"ZoneId":          fmt.Sprint(ZoneId),
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
	}
	return nil
}
