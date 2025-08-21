package alibabacloudstack

import (
	"fmt"
	"regexp"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackMaxcomputeUsers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackMaxcomputeUsersRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				MinItems: 1,
			},
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"organization_id": {
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
						"user_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_pk": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"organization_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"organization_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackMaxcomputeUsersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := make(map[string]interface{})
	if v, ok := d.GetOk("organization_id"); ok {
		request["Department"] = v.(string)
	}
	request["Region"] = client.RegionId
	request["Action"] = "GetOdpsUserList"
	request["AccessKeyId"] = client.AccessKey

	responseData, err := client.DoTeaRequest("GET", "ascm", "2019-05-10", "GetOdpsUserList", "", nil, request, nil)
	addDebug("GetOdpsUserList", responseData, request, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"Error OdpsUser Not Found"}) {
			return errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return err
	}
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			idsMap[Trim(vv.(string))] = Trim(vv.(string))
		}
	}
	users := make([]map[string]interface{}, 0)
	ids := make([]string, 0)
	datas, err := jsonpath.Get("$.data", responseData)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	for _, v := range datas.([]interface{}) {
		object := v.(map[string]interface{})

		if description_regex, ok := connectivity.GetResourceDataOk(d, "description_regex", "name_regex"); ok {
			r := regexp.MustCompile(description_regex.(string))
			if !r.MatchString(object["userName"].(string)) {
				continue
			}
		}

		if len(idsMap) > 0 {
			if _, exist := idsMap[fmt.Sprint(object["id"])]; !exist {
				continue
			}
		}
		user := map[string]interface{}{
			"id":                fmt.Sprint(object["id"]),
			"user_id":           object["userId"].(string),
			"user_name":         object["userName"].(string),
			"user_type":         object["userType"].(string),
			"organization_id":   fmt.Sprint(object["organizationId"]),
			"organization_name": object["organizationName"].(string),
			"description":       object["description"].(string),
			"user_pk":           object["aasPk"].(string),
		}
		users = append(users, user)
		ids = append(ids, fmt.Sprint(object["id"]))

	}
	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("users", users); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}
	return nil
}
