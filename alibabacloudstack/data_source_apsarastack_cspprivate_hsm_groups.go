package alibabacloudstack

import (
	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackCspprivateHsmGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackCspprivateHsmGroupsRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A list of group names to filter results by group name.",
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
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"unique_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"security_level_tag": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hsm_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_ids": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackCspprivateHsmGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// Build request parameters
	request := map[string]interface{}{
		"PageSize":   100,
		"PageNumber": 1,
	}

	// Initialize filters
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var finalObjects []interface{}
	pageNumber := 1

	// Paginate through all results
	for {
		request["PageNumber"] = pageNumber

		response, err := client.DoTeaRequest("POST", "cspprivate", "2022-02-17", "DescribeHsmGroups", "", nil, request, nil)
		if err != nil {
			return errmsgs.WrapError(err)
		}

		// Parse the response
		hsmGroups, err := jsonpath.Get("$.HsmGroups.HsmGroup", response)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		groups := hsmGroups.([]interface{})
		if len(groups) == 0 {
			break // No more data, exit loop normally
		}

		// Filter results based on ids and name_regex
		for _, item := range hsmGroups.([]interface{}) {
			hsmGroup := item.(map[string]interface{})

			groupName := hsmGroup["GroupName"].(string)

			// Check if id is in the ids filter
			if len(idsMap) > 0 {
				if _, ok := idsMap[groupName]; !ok {
					continue
				}
			}

			finalObjects = append(finalObjects, hsmGroup)
		}

		// Check if we have more pages
		totalCount, ok := response["TotalCount"].(float64)
		if !ok {
			break
		}

		pageSize := 100.0
		if pageSizeValue, ok := response["PageSize"].(float64); ok {
			pageSize = pageSizeValue
		}

		if float64(pageNumber)*pageSize >= totalCount {
			break
		}

		pageNumber++
	}

	// Prepare the result
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)

	for _, item := range finalObjects {
		hsmGroup := item.(map[string]interface{})

		mapping := map[string]interface{}{}

		mapping["id"] = hsmGroup["GroupName"]
		mapping["group_name"] = hsmGroup["GroupName"]
		mapping["status"] = hsmGroup["Status"]
		if uniqueId, ok := hsmGroup["UniqueId"].(float64); ok {
			mapping["unique_id"] = int(uniqueId)
		} else {
			mapping["unique_id"] = hsmGroup["UniqueId"]
		}
		mapping["security_level_tag"] = hsmGroup["SecurityLevelTag"]
		mapping["create_time"] = hsmGroup["CreateTime"]
		if hsmCount, ok := hsmGroup["HsmCount"].(float64); ok {
			mapping["hsm_count"] = int(hsmCount)
		} else {
			mapping["hsm_count"] = hsmGroup["HsmCount"]
		}
		mapping["vpc_id"] = hsmGroup["VpcId"]
		mapping["update_time"] = hsmGroup["UpdateTime"]
		mapping["zone_ids"] = hsmGroup["ZoneIds"]

		groupName := hsmGroup["GroupName"].(string)
		ids = append(ids, groupName)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("groups", s); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
