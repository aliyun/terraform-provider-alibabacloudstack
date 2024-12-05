package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"regexp"
)

func dataSourceAlibabacloudStackAscmUserGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackAscmUserGroupsRead,
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
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"organization_id": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"role_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"organization_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"user_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"users": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"role_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackAscmUserGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := client.NewCommonRequest("GET", "ascm", "2019-05-10", "ListUserGroups", "")
	userGroupName := d.Get("name_regex").(string)
	request.QueryParams["userGroupName"] = userGroupName

	response := UserGroup{}
	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})

		bresponse, ok := raw.(*responses.CommonResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_users", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		addDebug("ListUserGroups", raw, request)

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
	var roleids []int
	var users []string
	var groups []map[string]interface{}
	for _, group := range response.Data {
		if reg != nil && !reg.MatchString(group.GroupName) {
			continue
		}

		for _, rid := range response.Data[0].Roles {
			roleids = append(roleids, rid.Id)
		}

		for _, user := range response.Data[0].Users {
			users = append(users, user.Username)
		}

		mapping := map[string]interface{}{
			"id":             fmt.Sprint(group.Id),
			"group_name":     group.GroupName,
			"organization_id": group.Organization.Id,
			"user_group_id":  group.AugId,
			"role_ids":       roleids,
			"users":          users,
		}

		ids = append(ids, fmt.Sprint(group.Id))
		groups = append(groups, mapping)
	}
	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("groups", groups); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("role_ids", roleids); err != nil {
		return errmsgs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), groups)
	}
	return nil
}
