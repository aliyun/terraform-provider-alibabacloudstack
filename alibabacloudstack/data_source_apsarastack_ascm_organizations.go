package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackAscmOrganizations() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackAscmOrganizationsRead,
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
			"parent_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"organizations": {
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
						"cuser_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"muser_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"alias": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parent_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"internal": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackAscmOrganizationsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	var parentId string
	if v, ok := d.GetOk("parent_id"); ok {
		parentId = fmt.Sprint(v.(int))
	} else {
		parentId = client.Department
	}

	request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "GetOrganizationList", "/ascm/auth/organization/queryList")
	request.QueryParams["id"] = parentId

	response := OrganizationListResponse{}

	for {
		bresponse, err := client.ProcessCommonRequest(request)
		if err != nil {
			if bresponse == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_acm_configuration", "DescribeInstanceTypeFamilies", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
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
	if nameRegex, ok := d.GetOk("name_regex"); ok && nameRegex.(string) != "" {
		r = regexp.MustCompile(nameRegex.(string))
	}

	//parent_id
	var ids []string
	var s []map[string]interface{}
	for _, rg := range response.Data {
		if r != nil && !r.MatchString(rg.Name) {
			continue
		}
		mapping := map[string]interface{}{
			"id":        fmt.Sprint(rg.ID),
			"name":      rg.Name,
			"parent_id": rg.ParentID,
			"muser_id":  rg.MuserID,
			"cuser_id":  rg.CuserID,
			"alias":     rg.Alias,
			"internal":  rg.Internal,
		}
		ids = append(ids, fmt.Sprint(rg.ID))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("organizations", s); err != nil {
		return errmsgs.WrapError(err)
	}
	if s == nil {
		d.SetId(parentId)
	}
	return nil
}
