package alibabacloudstack

import (
	"encoding/json"
	"regexp"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackNasNamespaceGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackNasNamespaceGroupsRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A list of NAS namespace group IDs.",
			},
			"mapped_path_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  "A regex string to filter results by namespace group name.",
			},
			"nas_namespace_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter results by NAS namespace ID.",
			},
			"mount_target_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter results by mount target domain.",
			},
			"network_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter results by network type (e.g., Vpc or Classic).",
			},
			"names": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A list of NAS namespace group names.",
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nas_namespace_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mapped_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mount_target_domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"member_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Description: "A list of NAS namespace groups. Each element contains attributes of a namespace group.",
			},
		},
	}
}

func dataSourceAlibabacloudStackNasNamespaceGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	action := "DescribeNamespaceGroup"

	// Create request parameters
	request := make(map[string]interface{})

	if v, ok := d.GetOk("nas_namespace_id"); ok && v.(string) != "" {
		request["NasNamespaceId"] = v.(string)
	}
	if v, ok := d.GetOk("mount_target_domain"); ok && v.(string) != "" {
		request["MountTargetDomain"] = v.(string)
	}
	request["PageSize"] = 10

	idsMap := getIdsStringFilter(d)

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("mapped_path_regex"); ok {
		nameRegex = regexp.MustCompile(v.(string))
	}
	pageNumber := 1
	var namespaceGroups []interface{}
	for {
		request["PageNumber"] = pageNumber
		response, err := client.DoTeaRequest("POST", "nas", "2017-06-26", action, "", nil, nil, request)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_nas_namespace_groups", action, errmsgs.AlibabacloudStackSdkGoERROR, "")
		}

		data, ok := response["NasNamespaces"].([]interface{})
		if ok {
			namespaceGroups = append(namespaceGroups, data...)
		}
		totalCount, _ := response["TotalCount"].(json.Number).Int64()
		if pageNumber*10 >= int(totalCount) {
			break
		}
		pageNumber++

	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	groups := make([]map[string]interface{}, 0)

	for _, item := range namespaceGroups {
		namespace := item.(map[string]interface{})

		mapping := map[string]interface{}{}
		mountTargetDomain := ""
		if domain, ok := namespace["MountTargetDomain"]; ok && domain != nil {
			mountTargetDomain = domain.(string)
			if len(mountTargetDomain) > 0 && mountTargetDomain[len(mountTargetDomain)-1] == '.' {
				mountTargetDomain = mountTargetDomain[:len(mountTargetDomain)-1]
			}
		}
		if len(idsMap) > 0 {
			if _, exists := idsMap[mountTargetDomain]; !exists {
				continue
			}
		}

		if nameRegex != nil {
			mappedPath := ""
			if desc, ok := namespace["MappedPath"]; ok && desc != nil {
				mappedPath = desc.(string)
			}
			if !nameRegex.MatchString(mappedPath) {
				continue
			}
		}
		if v, ok := d.GetOk("network_type"); ok {
			if v.(string) != namespace["NetworkType"].(string) {
				continue
			}
		}
		mapping["mount_target_domain"] = mountTargetDomain
		mapping["nas_namespace_id"] = namespace["NasNamespaceId"]
		mapping["mapped_path"] = namespace["MappedPath"]

		mapping["network_type"] = namespace["NetworkType"]
		mapping["member_id"] = namespace["MemberId"]
		mapping["status"] = namespace["Status"]
		mapping["create_time"] = namespace["CreateTime"]

		// Add to results
		if id, ok := namespace["NasNamespaceId"]; ok {
			ids = append(ids, id.(string))
			names = append(names, id.(string)) // Using NasNamespaceId as name
		}

		groups = append(groups, mapping)
	}

	// Set resource data
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

	return nil
}
