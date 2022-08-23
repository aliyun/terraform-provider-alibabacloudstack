package apsarastack

import (
	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"regexp"
	"strconv"
	"strings"
)

func dataSourceApsaraStackDnsRecords() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApsaraStackDnsRecordsRead,

		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"host_record_regex": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				// must be one of [A, NS, MX, TXT, CNAME, SRV, AAAA, CAA, REDIRECT_URL, FORWORD_URL]
				ValidateFunc: validation.StringInSlice([]string{"A", "NS", "MX", "TXT", "CNAME", "SRV", "AAAA", "CAA", "REDIRECT_URL", "FORWORD_URL"}, false),
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values
			"records": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"record_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"domain_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_record": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rr_set": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"ttl": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceApsaraStackDnsRecordsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	DomainId := d.Get("domain_id").(string)
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
	request.ApiName = "ObtainGlobalAuthRecordList"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{
		"AccessKeySecret": client.SecretKey,
		"AccessKeyId":     client.AccessKey,
		"Product":         "GenesisDns",
		"RegionId":        client.RegionId,
		"Action":          "ObtainGlobalAuthRecordList",
		"Version":         "2018-07-20",
		"Id":              DomainId,
	}

	response := DnsRecord{}

	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		log.Printf(" response of raw ObtainGlobalAuthRecordList : %s", raw)

		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "apsarastack_dns_records", request.GetActionName(), ApsaraStackSdkGoERROR)
		}

		bresponse, _ := raw.(*responses.CommonResponse)

		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
		if err != nil {
			return WrapError(err)
		}
		if response.Records != nil {
			break
		}

	}
	var r *regexp.Regexp
	if nameRegex, ok := d.GetOk("host_record_regex"); ok && nameRegex.(string) != "" {
		r = regexp.MustCompile(nameRegex.(string))
	}

	var ids []string
	var s []map[string]interface{}
	for _, record := range response.Records {
		if r != nil && !r.MatchString(record.Rr) {
			continue
		}
		mapping := map[string]interface{}{
			"record_id":   record.RecordID,
			"domain_id":   DomainId,
			"host_record": record.Rr,
			"type":        record.Type,
			"description": record.Remark,
			"rr_set":      record.RrSet,
			"ttl":         record.TTL,
		}
		ids = append(ids, strconv.Itoa(record.RecordID))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("records", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
