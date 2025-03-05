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

func dataSourceAlibabacloudStackBastionhostInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackBastionhostInstancesRead,

		Schema: map[string]*schema.Schema{
			"description_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
			},
			"descriptions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_domain": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"instance_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"license_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_network_access": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"security_group_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"tags": tagsSchema(),
					},
				},
			},
			"tags": tagsSchema(),
		},
	}
}

func dataSourceAlibabacloudStackBastionhostInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	action := "QueryInstance"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1

	var objects []map[string]interface{}

	// get name Regex
	var descriptionRegex *regexp.Regexp
	if v, ok := d.GetOk("description_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return errmsgs.WrapError(err)
		}
		descriptionRegex = r
	}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		var idsStr []string
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
			idsStr = append(idsStr, vv.(string))
		}
		request["InstanceId"] = idsStr
	}

	if v, ok := d.GetOk("tags"); ok {
		tags := make([]map[string]interface{}, 0)
		for key, value := range v.(map[string]interface{}) {
			tags = append(tags, map[string]interface{}{
				"Key":   key,
				"Value": value.(string),
			})
		}
		request["Tag.*"] = tags
	}
	for {
		response, err := client.DoTeaRequest("POST", "Bastionhostprivate", "2023-03-23", action, "", nil, nil, request)
		if err != nil {
			return err
		}
		addDebug(action, response, request)
		// addDebug(action, response, request)
		resp, err := jsonpath.Get("$.Instances", response)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, action, "$.Instances", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if descriptionRegex != nil {
				if !descriptionRegex.MatchString(fmt.Sprint(item["Description"])) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["InstanceId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":                    object["InstanceId"],
			"description":           object["Description"],
			"user_vswitch_id":       object["VswitchId"],
			"private_domain":        object["IntranetEndpoint"],
			"public_domain":         object["InternetEndpoint"],
			"instance_status":       object["InstanceStatus"],
			"license_code":          object["LicenseCode"],
			"public_network_access": object["PublicNetworkAccess"],
		}

		id := fmt.Sprint(object["InstanceId"])
		bastionhostService := YundunBastionhostService{client}

		getResp, err := bastionhostService.DescribeBastionhostInstance(id)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		mapping["security_group_ids"] = getResp["SecurityGroupIds"]

		// getResp2, err := bastionhostService.ListTagResources(id, "instance")
		// if err != nil {
		// 	return errmsgs.WrapError(err)
		// }
		// mapping["tags"] = tagsToMap(getResp2)

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["Description"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}

	if err := d.Set("descriptions", names); err != nil {
		return errmsgs.WrapError(err)
	}

	if err := d.Set("instances", s); err != nil {
		return errmsgs.WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
