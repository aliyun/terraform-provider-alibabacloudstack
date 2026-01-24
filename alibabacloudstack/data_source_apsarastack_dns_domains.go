package alibabacloudstack

import (
	"encoding/json"
	"fmt"

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
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "The 'output_file' field has been deprecated and is scheduled for removal in version 3.19.0. To write content to a file, use the 'local_file' provider instead.",
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
		// FIXME: pagination logic error
		response, err := client.ProcessCommonRequest(request)
		if err != nil {
			if response == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_dns_domains", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), response, request)
		err = json.Unmarshal(response.GetHttpContentBytes(), &addDomains)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		if response.IsSuccess() == true || len(addDomains.Data) < 1 {
			break
		}
	}

	var domain_name string
	if v, ok := d.GetOk("domain_name"); ok {
		domain_name = v.(string)
	}
	idsMap := getIdsStringFilter(d)
	namesMap := getStringListFilters(d, "names")
	var ids []string
	var names []string
	var s []map[string]interface{}
	for _, rg := range addDomains.Data {
		if domain_name != "" && rg.Name != domain_name {
			continue
		}
		id := (rg.Id)
		if _, existed := idsMap[id]; len(idsMap) > 0 && !existed {
			continue
		}
		if _, existed := namesMap[rg.Name]; len(namesMap) > 0 && !existed {
			continue
		}
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
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return errmsgs.WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), s); err != nil {
			return err
		}
	}

	return nil
}
