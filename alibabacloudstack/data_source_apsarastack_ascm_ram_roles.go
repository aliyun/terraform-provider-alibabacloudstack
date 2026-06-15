package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackAscmRoles() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackAscmRolesRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:       schema.TypeInt,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "In future versions, searching by `id` is not supported, Please use `ids` instead. and is scheduled for removal in version 3.21.0",
			},
			"ids": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
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
				Optional: true,
			},
			"role_type": {
				Type:     schema.TypeString,
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
						"assume_role_policy_document": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"organization_visibility": {
							Type:     schema.TypeString,
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
	pageSize := 100
	request.QueryParams["pageSize"] = strconv.Itoa(pageSize)
	currentPage := 1

	//request.QueryParams["roleType"] = roleType

	response := ListAscmRolesResponse{}

	data := []AscmRoleData{}

	for {
		request.QueryParams["currentPage"] = strconv.Itoa(currentPage)
		bresponse, err := client.ProcessCommonRequest(request)
		addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
		if err != nil {
			if bresponse == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_ram_role", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		data = append(data, response.Data...)
		if response.PageInfo.TotalPage <= currentPage || len(response.Data) < pageSize {
			break
		}
		currentPage += 1
	}

	var r *regexp.Regexp
	if nameRegex, ok := d.GetOk("name_regex"); ok && nameRegex.(string) != "" {
		r = regexp.MustCompile(nameRegex.(string))
	}
	idsMap := getIdsStringFilter(d)
	var ids []string
	var s []map[string]interface{}
	ascmservice := AscmService{client}
	for _, rg := range data {
		roleid := fmt.Sprintf("%s:%d", rg.RoleName, rg.ID)
		if r != nil && !r.MatchString(rg.RoleName) {
			continue
		}
		if id != 0 && rg.ID != id {
			continue
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[roleid]; !ok {
				continue
			}
		}
		if roleType != "" && rg.RoleType != roleType {
			continue
		}
		log.Printf("[DEBUG] alibabacloudstack_ascm_ram_role ------------------------------------------ role.assume_role_policy_document: %s", rg.AssumeRolePolicyDocument)
		mapping := map[string]interface{}{
			"id":                          roleid,
			"name":                        rg.RoleName,
			"owner_organization_id":       rg.OwnerOrganizationID,
			"description":                 rg.Description,
			"user_count":                  rg.UserCount,
			"role_level":                  rg.RoleLevel,
			"role_type":                   rg.RoleType,
			"role_range":                  rg.RoleRange,
			"ram_role":                    rg.RAMRole,
			"enable":                      rg.Enable,
			"active":                      rg.Active,
			"default":                     rg.Default,
			"code":                        rg.Code,
			"assume_role_policy_document": rg.AssumeRolePolicyDocument,
			"organization_visibility":     rg.OrganizationVisibility,
		}
		if rg.RoleType == "ROLETYPE_RAMROLEAUTHORIZATION" {
			ramRole, err := ascmservice.DescribeAscmRamRoleForRoleid(roleid)
			if err == nil {
				mapping["assume_role_policy_document"] = ramRole.AssumeRolePolicyDocument
			}
		}
		ids = append(ids, roleid)
		s = append(s, mapping)
	}
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("roles", s); err != nil {
		return errmsgs.WrapError(err)
	}

	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), s); err != nil {
			return err
		}
	}
	return nil
}
