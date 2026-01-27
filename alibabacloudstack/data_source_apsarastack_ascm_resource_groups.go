package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackAscmResourceGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackAscmResourceGroupsRead,
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
			"organization_id": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"output_file": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "The 'output_file' field has been deprecated and is scheduled for removal in version 3.19.0. To write content to a file, use the 'local_file' provider instead.",
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
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"organization_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"rs_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creator": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gmt_created": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackAscmResourceGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "ListResourceGroup", "/ascm/auth/resource_group/list_resource_group")
	if v, ok := d.GetOk("organization_id"); ok {
		request.QueryParams["OrganizationId"] = fmt.Sprint(v)
		request.QueryParams["Department"] = fmt.Sprint(v)
	} else {
		request.QueryParams["OrganizationId"] = ""
		request.QueryParams["Department"] = ""

	}
	response := ResourceGroup{}

	bresponse, err := client.ProcessCommonRequest(request)
	addDebug(" ================== ListResourceGroup ==================", bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_resource_group", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		nameRegex = regexp.MustCompile(v.(string))
	}
	idsMap := getIdsStringFilter(d)
	var ids []string
	var s []map[string]interface{}
	for _, rg := range response.Data {
		resourceId := fmt.Sprintf("%d:%d", rg.OrganizationID, rg.ID)
		if nameRegex != nil && !nameRegex.MatchString(rg.ResourceGroupName) {
			continue
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[resourceId]; !ok {
				continue
			}
		}
		mapping := map[string]interface{}{
			"id":                  resourceId,
			"name":                rg.ResourceGroupName,
			"organization_id":     rg.OrganizationID,
			"creator":             rg.Creator,
			"gmt_created":         time.Unix(rg.GmtCreated/1000, 0).Format("2006-01-02 03:04:05"),
			"rs_id":               rg.RsID,
			"resource_group_type": rg.ResourceGroupType,
		}
		ids = append(ids, resourceId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("groups", s); err != nil {
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
