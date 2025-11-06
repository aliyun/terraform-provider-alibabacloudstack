package alibabacloudstack

import (
	"fmt"
	"regexp"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackAPIGatewayV2RouteGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackApiGatewayV2RouteGroupsRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"route_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"base_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"editable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"domains": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"create_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"domain": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"domain_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackApiGatewayV2RouteGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// Build request parameters
	request := make(map[string]interface{})
	gwInstanceId := d.Get("instance_id").(string)
	request["gwInstanceId"] = gwInstanceId
	request["size"] = 100
	request["current"] = 1

	// Call ListGroups API
	resp, err := client.DoTeaRequest("POST", "csb2", "2023-02-06", "ListGroups", "/group/listGroups", nil, nil, request)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	groups, err := jsonpath.Get("$.data.records", resp)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		nameRegex = regexp.MustCompile(v.(string))
	}
	var result []map[string]interface{}
	var ids []string
	for _, item := range groups.([]interface{}) {
		group, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		groupId := fmt.Sprintf("%s:%s", gwInstanceId, group["groupId"].(string))
		name, _ := group["name"].(string)

		if len(idsMap) > 0 {
			if _, exists := idsMap[groupId]; !exists {
				continue
			}
		}

		if nameRegex != nil && !nameRegex.MatchString(name) {
			continue
		}
		mapping := map[string]interface{}{
			"id":          groupId,
			"instance_id": group["gwInstanceId"],
			"group_id":    group["groupId"],
			"name":        group["name"],
			"base_path":   group["basePath"],
			"description": group["description"],
			"create_time": group["createTime"],
			"update_time": group["updateTime"],
			"editable":    group["editable"],
		}

		if domainList, ok := group["domains"].([]interface{}); ok {
			domains := make([]map[string]interface{}, 0)
			for _, domainItem := range domainList {
				if domainMap, ok := domainItem.(map[string]interface{}); ok {
					domainEntry := map[string]interface{}{
						"protocol":    domainMap["protocol"],
						"create_time": domainMap["createTime"],
						"domain":      domainMap["domain"],
						"domain_id":   domainMap["domainId"],
					}
					domains = append(domains, domainEntry)
				}
			}
			mapping["domains"] = domains
		} else {
			mapping["domains"] = []interface{}{}
		}

		result = append(result, mapping)
		ids = append(ids, groupId)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("route_groups", result); err != nil {
		return err
	}

	return nil
}
