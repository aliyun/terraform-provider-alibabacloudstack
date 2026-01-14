package alibabacloudstack

import (
	"regexp"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackDnsGtmInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackDnsGtmInstancesRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
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
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"prefix": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ttl": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"create_timestamp": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"update_timestamp": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackDnsGtmInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// Build request parameters
	request := make(map[string]interface{})
	request["PageNumber"] = 1
	request["PageSize"] = 100

	response, err := client.DoTeaRequest("POST", "CloudDns", "2021-06-24", "DescribeDnsGtmInstances", "", nil, nil, request)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	if _, ok := response["Data"]; !ok {
		return errmsgs.GetNotFoundErrorFromString("DNS GTM instance list not found")
	}

	// Prepare filters
	idsMap := getIdsStringFilter(d)

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		nameRegex, err = regexp.Compile(v.(string))
		if err != nil {
			return errmsgs.WrapError(err)
		}
	}

	// Filter results
	var filteredInstances []map[string]interface{}
	for _, item := range response["Data"].([]interface{}) {
		instance, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		id, _ := instance["Id"].(string)
		name, _ := instance["Name"].(string)

		// Apply ID filter
		if len(idsMap) > 0 {
			if _, exists := idsMap[id]; !exists {
				continue
			}
		}

		// Apply name regex filter
		if nameRegex != nil && !nameRegex.MatchString(name) {
			continue
		}

		filteredInstances = append(filteredInstances, instance)
	}

	// Prepare output
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)

	for _, object := range filteredInstances {
		mapping := map[string]interface{}{}
		mapping["id"] = object["Id"]
		mapping["name"] = object["Name"]
		mapping["zone_id"] = object["ZoneId"]
		mapping["zone_name"] = object["ZoneName"]
		mapping["prefix"] = object["Prefix"]
		mapping["ttl"] = object["Ttl"]
		mapping["create_timestamp"] = object["CreateTimestamp"]
		mapping["update_timestamp"] = object["UpdateTimestamp"]

		ids = append(ids, object["Id"].(string))
		names = append(names, object["Name"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("instances", s); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
