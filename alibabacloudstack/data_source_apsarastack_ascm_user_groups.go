package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackAscmUserGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackAscmUserGroupsRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id":              {Type: schema.TypeString, Computed: true},
						"group_name":      {Type: schema.TypeString, Computed: true},
						"organization_id": {Type: schema.TypeString, Computed: true},
						"user_group_id":   {Type: schema.TypeString, Computed: true},
						"role_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"users": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackAscmUserGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "ListUserGroups", "/ascm/auth/user/listUserGroups")
	userGroupName := d.Get("name_regex").(string)
	request.QueryParams["userGroupName"] = userGroupName

	response := UserGroup{}
	for {
		bresponse, err := client.ProcessCommonRequest(request)
		if err != nil {
			errmsg := ""
			if bresponse != nil {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_users", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		addDebug("ListUserGroups", bresponse, request)

		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		if response.Code == "200" || len(response.Data) < 1 {
			break
		}
	}

	var reg *regexp.Regexp
	if nameRegex, ok := d.GetOk("name_regex"); ok && nameRegex.(string) != "" {
		reg = regexp.MustCompile(nameRegex.(string))
	}

	var ids []string
	var names []string
	var groups []map[string]interface{}

	for _, group := range response.Data {
		if reg != nil && !reg.MatchString(group.GroupName) {
			continue
		}

		// ✅ 提取当前组的 role_ids
		var roleIds []string
		for _, role := range group.Roles {
			if role.Id != 0 {
				roleIds = append(roleIds, fmt.Sprintf("%d", role.Id))
			}
		}

		// ✅ 提取当前组的 users
		var users []string
		for _, user := range group.Users {
			if user.Username != "" {
				users = append(users, user.Username)
			}
		}

		mapping := map[string]interface{}{
			"id":              fmt.Sprint(group.Id),
			"group_name":      group.GroupName,
			"organization_id": strconv.Itoa(group.Organization.Id),
			"user_group_id":   group.AugId,
			"role_ids":        roleIds,
			"users":           users,
		}

		ids = append(ids, fmt.Sprint(group.Id))
		names = append(names, group.GroupName)
		groups = append(groups, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("groups", groups); err != nil {
		return errmsgs.WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), groups); err != nil {
			return err
		}
	}
	return nil
}
