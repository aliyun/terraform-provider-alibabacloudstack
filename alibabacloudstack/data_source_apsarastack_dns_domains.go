package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackDnsDomains() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackDnsDomainsRead,

		Schema: map[string]*schema.Schema{
			"domain_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"domains": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dns_servers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackDnsDomainsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := client.NewCommonRequest("POST", "CloudDns", "2021-06-24", "DescribeGlobalZones", "")
	request.QueryParams["PageNumber"] = fmt.Sprint(1)
	request.QueryParams["PageSize"] = fmt.Sprint(PageSizeLarge)
	request.QueryParams["Name"] = d.Get("domain_name").(string)
	request.QueryParams["Forwardedregionid"] = client.RegionId
	request.QueryParams["SignatureVersion"] = "2.1"

	var addDomains = DnsDomains{}
	for {
		raw, err := client.WithDnsClient(func(dnsClient *alidns.Client) (interface{}, error) {
			return dnsClient.ProcessCommonRequest(request)
		})
		response, ok := raw.(*responses.CommonResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_dns_domains", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request)
		err = json.Unmarshal(response.GetHttpContentBytes(), &addDomains)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		if response.IsSuccess() == true || len(addDomains.Data) < 1 {
			break
		}
	}

	var r *regexp.Regexp
	if nameRegex, ok := d.GetOk("domain_name"); ok && nameRegex.(string) != "" {
		r = regexp.MustCompile(nameRegex.(string))
	}
	var ids []string
	var names []string
	var s []map[string]interface{}
	for _, rg := range addDomains.Data {
		if r != nil && !r.MatchString(rg.Name) {
			continue
		}
		id := (rg.Id)
		mapping := map[string]interface{}{
			"domain_id":   id,
			"domain_name": rg.Name,
		}

		names = append(names, rg.Name)
		ids = append(ids, id)
		s = append(s, mapping)
	}
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("domains", s); err != nil {
		return errmsgs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), s); err != nil {
			return err
		}
	}

	return nil
}
