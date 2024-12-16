package alibabacloudstack

import (
	"fmt"
	"log"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
				Type:     schema.TypeString,
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
			"line_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Default:  []string{"default"},
			},
			"remark": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlibabacloudStackDnsRecordCreate(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := make(map[string]interface{})
	ZoneId := d.Get("zone_id").(string)
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

	action := "AddGlobalZoneRecord"
	request["Type"] = Type
	request["Ttl"] = TTL
	request["ZoneId"] = ZoneId
	request["LbaStrategy"] = LbaStrategy
	request["Name"] = Name
	request["LineIds"] = d.Get("line_ids").([]string)
	request["ClientToken"] = buildClientToken(action)

	_, err = client.DoTeaRequest("POST", "CloudDns", "2021-06-24", action, "", nil, request)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprint(ZoneId))

	return resourceAlibabacloudStackDnsRecordRead(d, meta)
}

func resourceAlibabacloudStackDnsRecordRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)

	dnsService := &DnsService{client: client}
	object, err := dnsService.DescribeDnsRecord(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
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
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "IsRecordExist", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	attributeUpdate := false

	var desc string

	if d.HasChange("remark") {
		if v, ok := d.GetOk("remark"); ok {
			desc = v.(string)
		}
		check.Data[0].Remark = desc
		request := client.NewCommonRequest("POST", "CloudDns", "2021-06-24", "UpdateGlobalZoneRecordRemark", "")
		request.QueryParams["Id"] = fmt.Sprint(ID)
		request.QueryParams["Remark"] = desc
		raw, err := client.WithDnsClient(func(dnsClient *alidns.Client) (interface{}, error) {
			return dnsClient.ProcessCommonRequest(request)
		})
		log.Printf(" response of raw UpdateGlobalZoneRecordRemark : %s", raw)
		if err != nil {
			return err
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
		request["Type"] = Type
		request["Ttl"] = Ttl
		request["Id"] = ID
		request["ZoneId"] = ZoneId
		request["LbaStrategy"] = LbaStrategy
		request["Name"] = Name
		request["Remark"] = check.Data[0].Remark
		request["ClientToken"] = buildClientToken(action)

		_, err := client.DoTeaRequest("POST", "CloudDns", "2021-06-24", action, "", nil, request)
		if err != nil {
			return err
		}
	}

	return resourceAlibabacloudStackDnsRecordRead(d, meta)
}

func resourceAlibabacloudStackDnsRecordDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ID := d.Get("record_id").(int)
	ZoneId := d.Get("zone_id").(int)

	request := client.NewCommonRequest("POST", "CloudDns", "2021-06-24", "DeleteGlobalZoneRecord", "")
	request.QueryParams["Id"] = fmt.Sprint(ID)
	request.QueryParams["ZoneId"] = fmt.Sprint(ZoneId)
	raw, err := client.WithDnsClient(func(dnsClient *alidns.Client) (interface{}, error) {
		return dnsClient.ProcessCommonRequest(request)
	})
	addDebug(request.GetActionName(), raw)

	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, []string{"DomainRecordNotBelongToUser"}) {
			return nil
		}
		if errmsgs.IsExpectedErrors(err, []string{"RecordForbidden.DNSChange", "InternalError"}) {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	return nil
}
