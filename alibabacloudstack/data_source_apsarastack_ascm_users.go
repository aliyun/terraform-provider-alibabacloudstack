package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackAscmUsers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackAscmUsersRead,
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
			//后续未消费
// 			"names": {
// 				Type:     schema.TypeList,
// 				Computed: true,
// 				Elem:     &schema.Schema{Type: schema.TypeString},
// 			},
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
			"users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"login_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"organization_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cell_phone_number": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"email": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mobile_nation_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
						"default_role_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"login_policy_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						// 后续未消费字段
// 						"acc": {
// 							Type:     schema.TypeInt,
// 							Computed: true,
// 						},
						"primary_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackAscmUsersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
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
	loginName := d.Get("name_regex").(string)
	orgId := d.Get("organization_id")

	request.RegionId = client.RegionId
	request.ApiName = "ListUsers"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	response := User{}
	request.QueryParams = map[string]string{
		
		
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"Product":         "ascm",
		"RegionId":        client.RegionId,
		"Action":          "ListUsers",
		"Version":         "2019-05-10",
		"loginName":       loginName,
		"organizationId":  fmt.Sprint(orgId),
	}

	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		log.Printf(" response of raw ListUsers : %s", raw)

		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alibabacloudstack_ascm_users", request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}

		bresponse, _ := raw.(*responses.CommonResponse)

		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
		if err != nil {
			return WrapError(err)
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
	var roleids []int
	var s []map[string]interface{}
	for _, u := range response.Data {
		if r != nil && !r.MatchString(u.LoginName) {
			continue
		}

		for _, rid := range response.Data[0].Roles {
			roleids = append(roleids, rid.ID)
		}
		//for _, r := range response.Data[0].Roles {

		mapping := map[string]interface{}{
			"id":                 fmt.Sprint(u.ID),
			"login_name":         u.LoginName,
			"organization_id":    u.Organization.ID,
			"cell_phone_number":  u.CellphoneNum,
			"display_name":       u.DisplayName,
			"email":              u.Email,
			"primary_key":        u.PrimaryKey,
			"mobile_nation_code": u.MobileNationCode,
			"default_role_id":    u.DefaultRole.ID,
			"role_ids":           roleids,
			"login_policy_id":    u.LoginPolicy.ID,
		}

		ids = append(ids, fmt.Sprint(u.ID))
		//roleids = append(roleids, r.ID)
		s = append(s, mapping)
	}
	//}
	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("users", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("role_ids", roleids); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
