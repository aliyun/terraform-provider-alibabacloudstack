package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

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
			// Subsequent unconsumed parameter
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
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "The 'output_file' field has been deprecated and is scheduled for removal in version 3.19.0. To write content to a file, use the 'local_file' provider instead.",
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
	request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "ListRAMServiceRoles", "/ascm/auth/role/listRAMServiceRoles")
	request.QueryParams["pageSize"] = "10"
	pageNumber := 1
	response := RamRole{}

	var ids []string
	var s []map[string]interface{}
	for {
		request.QueryParams["currentPage"] = strconv.Itoa(pageNumber)
		bresponse, err := client.ProcessCommonRequest(request)
		addDebug("ListRAMServiceRoles", bresponse, request, request.QueryParams)
		if err != nil {
			errmsg := ""
			if bresponse != nil {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_roles", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		if response.Code != "200" || len(response.Data) < 1 {
			break
		}

		var r *regexp.Regexp
		if nameRegex, ok := d.GetOk("product"); ok && nameRegex.(string) != "" {
			r = regexp.MustCompile(strings.ToUpper(nameRegex.(string)))
		}

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
		if len(response.Data) < 10 {
			break
		}
		pageNumber += 1
	}
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("roles", s); err != nil {
		return errmsgs.WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), s); err != nil {
			return err
		}
	}
	return nil
}
