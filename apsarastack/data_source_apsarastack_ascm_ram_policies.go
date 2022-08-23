package apsarastack

import (
	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"regexp"
	"strings"
	"time"
)

func dataSourceApsaraStackAscmRamPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApsaraStackAscmRamPoliciesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
				Computed: true,
				ForceNew: true,
				MinItems: 1,
			},
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
							Type:     schema.TypeInt,
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
						"ctime": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cuser_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_document": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceApsaraStackAscmRamPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	request := requests.NewCommonRequest()
	request.Product = "ascm"
	request.Version = "2019-05-10"
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.RegionId = client.RegionId
	request.ApiName = "ListRamPolicies"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{
		"AccessKeyId":     client.AccessKey,
		"AccessKeySecret": client.SecretKey,
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"Product":         "ascm",
		"RegionId":        client.RegionId,
		"Action":          "ListRamPolicies",
		"Version":         "2019-05-10",
		"pageSize":        "1000",
		//"policyName":name,
	}
	response := RamPolicies{}

	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		log.Printf(" response of raw ListRamPolicies : %s", raw)

		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "apsarastack_ascm_ram_policies", request.GetActionName(), ApsaraStackSdkGoERROR)
		}

		bresponse, _ := raw.(*responses.CommonResponse)

		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
		if err != nil {
			return WrapError(err)
		}
		if response.Code == "200" || len(response.Data) < 1 {
			break
		}

	}

	var r *regexp.Regexp
	if nameRegex, ok := d.GetOk("name_regex"); ok && nameRegex.(string) != "" {
		r = regexp.MustCompile(nameRegex.(string))
	}
	var ids []string
	var s []map[string]interface{}
	for _, rp := range response.Data {
		if r != nil && !r.MatchString(rp.PolicyName) {
			continue
		}
		mapping := map[string]interface{}{
			"id":              rp.ID,
			"name":            rp.PolicyName,
			"description":     rp.Description,
			"region":          rp.Region,
			"cuser_id":        rp.CuserID,
			"ctime":           time.Unix(rp.Ctime/1000, 0).Format("2006-01-02 03:04:05"),
			"policy_document": rp.PolicyDocument,
		}
		ids = append(ids, string(rune(rp.ID)))
		s = append(s, mapping)
	}
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("policies", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
