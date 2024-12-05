package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"

	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackAscmRamServiceRoles() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackAscmRamServiceRolesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
				MinItems: 1,
			},
			"product": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			// 后续未消费参数
			// 			"names": {
			// 				Type:     schema.TypeList,
			// 				Computed: true,
			// 				Elem:     &schema.Schema{Type: schema.TypeString},
			// 			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"roles": {
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
						"role_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"product": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"organization_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"aliyun_user_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackAscmRamServiceRolesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "ListRAMServiceRoles", "")
	request.QueryParams["roleType"] = "ROLETYPE_RAM"
	response := RamRole{}

	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		log.Printf(" response of raw ListRAMServiceRoles : %s", raw)

		bresponse, ok := raw.(*responses.CommonResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_roles", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
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
	if nameRegex, ok := d.GetOk("product"); ok && nameRegex.(string) != "" {
		r = regexp.MustCompile(strings.ToUpper(nameRegex.(string)))
	}
	var ids []string
	var s []map[string]interface{}
	for _, rg := range response.Data {
		if r != nil && !r.MatchString(rg.Product) {
			continue
		}
		mapping := map[string]interface{}{
			"id":                fmt.Sprint(rg.ID),
			"name":              rg.RoleName,
			"description":       rg.Description,
			"role_type":         rg.RoleType,
			"product":           rg.Product,
			"organization_name": rg.OrganizationName,
			"aliyun_user_id":    rg.AliyunUserID,
		}

		ids = append(ids, fmt.Sprint(rg.ID))
		s = append(s, mapping)
	}
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("roles", s); err != nil {
		return errmsgs.WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
