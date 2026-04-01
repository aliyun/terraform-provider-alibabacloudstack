package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
			"output_file": {
				Type:     schema.TypeString,
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

	response := Organization{}

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		bresponse, err := client.ProcessCommonRequest(request)
		addDebug(request.GetActionName(), bresponse, request, request.QueryParams)

		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{errmsgs.Throttling}) {
				time.Sleep(time.Duration(3) * time.Second)
				return resource.RetryableError(err)
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_organizations", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}

		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
		if err != nil {
			return resource.NonRetryableError(errmsgs.WrapError(err))
		}
		if response.Code == "200" || len(response.Data) < 1 {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("response code is not 200 or data is empty"))
	})
	if err != nil {
		return err
	}

	var r *regexp.Regexp
	if nameRegex, ok := d.GetOk("name_regex"); ok && nameRegex.(string) != "" {
		r = regexp.MustCompile(nameRegex.(string))
	}

	//parent_id
	var ids []string
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			idsMap[vv.(string)] = vv.(string)
		}
	}
	var s []map[string]interface{}
	for _, rg := range response.Data {
		if r != nil && !r.MatchString(rg.Name) {
			continue
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(rg.ID)]; !ok {
				continue
			}
		}
		mapping := map[string]interface{}{
			"id":          fmt.Sprint(rg.ID),
			"name":        rg.Name,
			"parent_id":   rg.ParentID,
			"muser_id":    rg.MuserID,
			"cuser_id":    rg.CuserID,
			"alias":       rg.Alias,
			"internal":    rg.Internal,
			"primary_key": rg.UUID,
		}
		ids = append(ids, fmt.Sprint(rg.ID))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("organizations", s); err != nil {
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
	if s == nil {
		d.SetId(parentId)
	}
	return nil
}
