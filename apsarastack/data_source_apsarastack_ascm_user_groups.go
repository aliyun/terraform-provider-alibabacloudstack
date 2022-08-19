package apsarastack

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"regexp"
)

func dataSourceApsaraStackAscmUserGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApsaraStackAscmUserGroupsRead,
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

func dataSourceApsaraStackAscmUserGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	request := requests.NewCommonRequest()
	if client.Config.Insecure {
		request.SetHTTPSInsecure(client.Config.Insecure)
	}
	request.Method = "GET"
	request.Product = "ascm"
	request.Version = "2019-05-10"
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	userGroupName := d.Get("name_regex").(string)

	request.RegionId = client.RegionId
	request.ApiName = "ListUserGroups"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	response := UserGroup{}
	request.QueryParams = map[string]string{
		"AccessKeyId":     client.AccessKey,
		"AccessKeySecret": client.SecretKey,
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"Product":         "ascm",
		"RegionId":        client.RegionId,
		"Action":          "ListUserGroups",
		"Version":         "2019-05-10",
		"userGroupName":   userGroupName,
	}

	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})

		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "apsarastack_ascm_users", request.GetActionName(), ApsaraStackSdkGoERROR)
		}

		addDebug("ListUserGroups", raw, request)
		bresponse, _ := raw.(*responses.CommonResponse)

		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
		if err != nil {
			return WrapError(err)
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
			"id":              fmt.Sprint(group.Id),
			"group_name":      group.GroupName,
			"organization_id": group.Organization.Id,
			"user_group_id":   group.AugId,

			"role_ids": roleids,
			"users":    users,
		}

		ids = append(ids, fmt.Sprint(group.Id))
		groups = append(groups, mapping)
	}
	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("groups", groups); err != nil {
		return WrapError(err)
	}
	if err := d.Set("role_ids", roleids); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), groups)
	}
	return nil
}
