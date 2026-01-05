package alibabacloudstack

import (
	"encoding/json"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackAscmRamPoliciesForUser() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackAscmRamPoliciesForUserRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
				Computed: true,
				ForceNew: true,
				MinItems: 1,
			},
			"login_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"attach_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"default_version": {
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

func dataSourceAlibabacloudStackAscmRamPoliciesForUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	lname := d.Get("login_name").(string)

	request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "ListRAMPoliciesForUser", "/ascm/auth/dataPrivilege/listRAMPoliciesForUser")
	request.QueryParams["loginName"] = lname

	response := RamPolicyUser{}

	bresponse, err := client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_ram_policies_for_user", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	var names []string
	var s []map[string]interface{}
	for _, rp := range response.Data.DataItem {
		mapping := map[string]interface{}{
			"policy_name":     rp.PolicyName,
			"policy_type":     rp.PolicyType,
			"description":     rp.Description,
			"default_version": rp.DefaultVersion,
			"attach_date":     time.Unix(rp.AttachDate/1000, 0).Format("2006-01-02 03:04:05"), // 使用转换后的时间戳
			"policy_document": rp.PolicyDocument,
		}
		names = append(names, rp.PolicyName)
		s = append(s, mapping)
	}
	d.SetId(dataResourceIdHash(names))
	if err := d.Set("policies", s); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
