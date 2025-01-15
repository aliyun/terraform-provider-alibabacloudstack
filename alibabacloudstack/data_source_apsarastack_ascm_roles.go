package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackAscmRoles() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackAscmRolesRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"role_type": {
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
						"role_level": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"role_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ram_role": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"role_range": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"default": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"active": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"owner_organization_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"code": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackAscmRolesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	id := d.Get("id").(int)
	roleType := d.Get("role_type").(string)

	request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "ListRoles", "/ascm/auth/role/listRoles")
	request.QueryParams["pageSize"] = "100000"
	//request.QueryParams["roleType"] = roleType

	response := AscmRoles{}

	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		log.Printf(" response of raw ListRoles : %s", raw)

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
		if response.AsapiErrorCode == "" || len(response.Data) < 1 {
			break
		}
	}

	var r *regexp.Regexp
	if nameRegex, ok := d.GetOk("name_regex"); ok && nameRegex.(string) != "" {
		r = regexp.MustCompile(nameRegex.(string))
	}
	var ids []string
	var s []map[string]interface{}

	for _, rg := range response.Data {
		if r != nil && !r.MatchString(rg.RoleName) {
			continue
		}
		if id != 0 && rg.ID == id {
			mapping := map[string]interface{}{
				"id":                    rg.ID,
				"name":                  rg.RoleName,
				"owner_organization_id": rg.OwnerOrganizationID,
				"description":           rg.Description,
				"user_count":            rg.UserCount,
				"role_level":            rg.RoleLevel,
				"role_type":             rg.RoleType,
				"role_range":            rg.RoleRange,
				"ram_role":              rg.RAMRole,
				"enable":                rg.Enable,
				"active":                rg.Active,
				"default":               rg.Default,
				"code":                  rg.Code,
			}
			ids = append(ids, fmt.Sprint(rg.ID))
			s = append(s, mapping)
			break
		}
		if id == 0 && roleType != "" && rg.RoleType == roleType {
			mapping := map[string]interface{}{
				"id":                    rg.ID,
				"name":                  rg.RoleName,
				"owner_organization_id": rg.OwnerOrganizationID,
				"description":           rg.Description,
				"user_count":            rg.UserCount,
				"role_level":            rg.RoleLevel,
				"role_type":             rg.RoleType,
				"role_range":            rg.RoleRange,
				"ram_role":              rg.RAMRole,
				"enable":                rg.Enable,
				"active":                rg.Active,
				"default":               rg.Default,
				"code":                  rg.Code,
			}
			ids = append(ids, fmt.Sprint(rg.ID))
			s = append(s, mapping)
		}
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
