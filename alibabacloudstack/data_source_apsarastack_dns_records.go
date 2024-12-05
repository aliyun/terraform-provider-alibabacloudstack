package alibabacloudstack

import (
	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"regexp"
	"strconv"
)

func dataSourceAlibabacloudStackDnsRecords() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackDnsRecordsRead,

		Schema: map[string]*schema.Schema{
			"zone_id": {
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
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remark": {
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

func dataSourceAlibabacloudStackDnsRecordsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ZoneId := d.Get("zone_id").(string)

	request := client.NewCommonRequest("POST", "CloudDns", "2021-06-24", "DescribeGlobalZoneRecords", "")
	request.QueryParams["ZoneId"] = ZoneId

	response := DnsRecord{}

	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		log.Printf(" response of raw ObtainGlobalAuthRecordList : %s", raw)

		bresponse, ok := raw.(*responses.CommonResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_dns_records", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		if response.Data != nil {
			break
		}
	}

	var r *regexp.Regexp
	if nameRegex, ok := d.GetOk("name"); ok && nameRegex.(string) != "" {
		r = regexp.MustCompile(nameRegex.(string))
	}

	var ids []string
	var s []map[string]interface{}
	for _, record := range response.Data {
		if r != nil && !r.MatchString(record.Name) {
			continue
		}
		mapping := map[string]interface{}{
			"record_id": record.Id,
			"zone_id":   ZoneId,
			"name":      record.Name,
			"type":      record.Type,
			"remark":    record.Remark,
			"ttl":       record.TTL,
		}
		ids = append(ids, strconv.Itoa(record.Id))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("records", s); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}
	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
