package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			"output_file": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "The 'output_file' field has been deprecated and is scheduled for removal in version 3.19.0. To write content to a file, use the 'local_file' provider instead.",
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

	request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "ListLoginPolicies", "/ascm/auth/loginPolicy/listLoginPolicies")
	request.QueryParams["name"] = name

	response := LoginPolicy{}

	page := 1

	for {
		request.QueryParams["currentPage"] = fmt.Sprintf("%d", page)
		bresponse, err := client.ProcessCommonRequest(request)
		if err != nil {
			if bresponse == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_acm_configuration", "ListLoginPolicies", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
		if err != nil {
			return errmsgs.WrapError(err)
		}

		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		if response.Code == "200" {
			break
		}
		page += 1
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
		return errmsgs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), t); err != nil {
			return err
		}
	}
	return nil
}
