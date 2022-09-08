package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"regexp"
)

func dataSourceAlibabacloudStackAscmLogonPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackAscmLogonPoliciesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
				MinItems: 1,
			},
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ip_ranges": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"rule": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rule": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_range": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"login_policy_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"start_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackAscmLogonPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	name := d.Get("name_regex").(string)
	request := requests.NewCommonRequest()
	if client.Config.Insecure {
		request.SetHTTPSInsecure(client.Config.Insecure)
	}
	request.Method = "GET"
	request.Product = "ascm"
	request.Version = "2019-05-10"
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}

	request.RegionId = client.RegionId
	request.ApiName = "ListLoginPolicies"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{
		"AccessKeyId":     client.AccessKey,
		"AccessKeySecret": client.SecretKey,
		"Product":         "ascm",
		"RegionId":        client.RegionId,
		"Action":          "ListLoginPolicies",
		"Version":         "2019-05-10",
		"name":            name,
	}
	response := LoginPolicy{}

	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		log.Printf(" rsponse of raw ListLoginPolicies : %s", raw)

		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alibabacloudstack_ascm_logon_policies", request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}

		bresponse, _ := raw.(*responses.CommonResponse)

		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
		if err != nil {
			return WrapError(err)
		}
		if response.Code == "200" {
			break
		}
	}

	var r *regexp.Regexp
	if nameRegex, ok := d.GetOk("name_regex"); ok && nameRegex.(string) != "" {
		r = regexp.MustCompile(nameRegex.(string))
	}
	var ids []string
	var t []map[string]interface{}
	for _, u := range response.Data {
		if r != nil && !r.MatchString(u.Name) {
			continue
		}
		for _, times := range response.Data {
			for _, k := range times.TimeRanges {

				var ipranges []string
				var iprange string
				for _, k := range u.IPRanges {
					ipranges = append(ipranges, k.IPRange)
					if len(ipranges) > 1 {
						iprange = iprange + "," + k.IPRange
					} else {
						iprange = k.IPRange
					}
				}
				allmapping := map[string]interface{}{
					"id":              fmt.Sprint(u.ID),
					"name":            u.Name,
					"rule":            u.Rule,
					"description":     u.Description,
					"ip_range":        iprange,
					"login_policy_id": u.LpID,
					"start_time":      k.StartTime,
					"end_time":        k.EndTime,
				}
				t = append(t, allmapping)
			}
		}

	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("policies", t); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), t)
	}
	return nil
}
