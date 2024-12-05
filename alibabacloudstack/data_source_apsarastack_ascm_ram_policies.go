package alibabacloudstack

import (
	"encoding/json"
	"log"
	"regexp"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackAscmRamPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackAscmRamPoliciesRead,
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

func dataSourceAlibabacloudStackAscmRamPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := client.NewCommonRequest("GET", "ascm", "2019-05-10", "ListRamPolicies", "")
	request.QueryParams["pageSize"] = "1000"
	//request.QueryParams["policyName"]= name

	response := RamPolicies{}

	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		log.Printf(" response of raw ListRamPolicies : %s", raw)

		bresponse, ok := raw.(*responses.CommonResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_ram_policies", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
		if err != nil {
			return errmsgs.WrapError(err)
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
		return errmsgs.WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
