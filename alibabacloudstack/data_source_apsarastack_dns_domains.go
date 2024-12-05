package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"regexp"
	"strings"
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
// 后续未消费使用
// 			"group_name_regex": {
// 				Type:     schema.TypeString,
// 				Optional: true,
// 				ForceNew: true,
// 			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
// 			"ali_domain": {
// 				Type:     schema.TypeBool,
// 				Optional: true,
// 				ForceNew: true,
// 			},
// 			"instance_id": {
// 				Type:     schema.TypeString,
// 				Optional: true,
// 				ForceNew: true,
// 			},
// 			"version_code": {
// 				Type:     schema.TypeString,
// 				Optional: true,
// 				ForceNew: true,
// 			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
// 			"resource_group_id": {
// 				Type:     schema.TypeString,
// 				Optional: true,
// 				ForceNew: true,
// 			},
			// Computed values
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
// 						"ali_domain": {
// 							Type:     schema.TypeBool,
// 							Computed: true,
// 						},
// 						"group_id": {
// 							Type:     schema.TypeString,
// 							Computed: true,
// 						},
// 						"group_name": {
// 							Type:     schema.TypeString,
// 							Computed: true,
// 						},
// 						"instance_id": {
// 							Type:     schema.TypeString,
// 							Computed: true,
// 						},
// 						"version_code": {
// 							Type:     schema.TypeString,
// 							Computed: true,
// 						},
// 						"puny_code": {
// 							Type:     schema.TypeString,
// 							Computed: true,
// 						},
						"dns_servers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
// 						"tags": tagsSchema(),
					},
				},
			},
		},
	}
}
func dataSourceAlibabacloudStackDnsDomainsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Product = "CloudDns"
	request.Domain = client.Domain
	request.Version = "2021-06-24"
	name := d.Get("domain_name").(string)
	request.PageNumber = requests.NewInteger(2)
	request.PageSize = requests.NewInteger(PageSizeLarge)
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "DescribeGlobalZones"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{
		
		
		"Product":           "CloudDns",
		"RegionId":          client.RegionId,
		"Action":            "DescribeGlobalZones",
		"Version":           "2021-06-24",
		"PageNumber":        fmt.Sprint(1),
		"PageSize":          fmt.Sprint(PageSizeLarge),
		"Name":              name,
		"Forwardedregionid": client.RegionId,
		"SignatureVersion":  "2.1",
	}

	var addDomains = DnsDomains{}
	for {
		raw, err := client.WithEcsClient(func(alidnsClient *ecs.Client) (interface{}, error) {
			return alidnsClient.ProcessCommonRequest(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alibabacloudstack_dns_domains", request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request)
		response, _ := raw.(*responses.CommonResponse)
		err = json.Unmarshal(response.GetHttpContentBytes(), &addDomains)
		if err != nil {
			return WrapError(err)
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
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
