package alibabacloudstack

import (
	"regexp"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackDnsForwardDomains() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackDnsForwardDomainsRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"forward_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"FORWARD_FIRST", "FORWARD_ONLY"}, false),
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"forward_domains": {
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
						"forward_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"forwarders": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"remark": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"caller_uid": {
							Type:     schema.TypeString,
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

func dataSourceAlibabacloudStackDnsForwardDomainsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// Build request parameters
	request := make(map[string]interface{})
	if v, ok := d.GetOk("name"); ok {
		request["Name"] = v.(string)
	}
	request["PageNumber"] = 1
	request["PageSize"] = 50

	// Call API to get the list of forward domains
	resp, err := client.DoTeaRequest("POST", "CloudDns", "2021-06-24", "DescribeGlobalForwardZones", "", nil, nil, request)
	if err != nil {
		return err
	}

	// Extract data from response
	data, ok := resp["Data"].([]interface{})
	if !ok {
		data = []interface{}{}
	}

	// Prepare filters
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

	forwardModeFilter := ""
	if v, ok := d.GetOk("forward_mode"); ok {
		forwardModeFilter = v.(string)
	}

	// Filter results based on conditions
	var filteredResults []interface{}
	for _, item := range data {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		id, ok := itemMap["Id"].(string)
		if !ok || (len(idsMap) > 0 && idsMap[id] == "") {
			continue
		}

		name, ok := itemMap["Name"].(string)
		if !ok || (nameRegex != nil && !nameRegex.MatchString(name)) {
			continue
		}

		if forwardModeFilter != "" {
			forwardMode, ok := itemMap["ForwardMode"].(string)
			if !ok || forwardMode != forwardModeFilter {
				continue
			}
		}

		filteredResults = append(filteredResults, itemMap)
	}

	// Prepare final result structure
	ids := make([]string, 0)
	result := make([]map[string]interface{}, 0)

	for _, item := range filteredResults {
		itemMap := item.(map[string]interface{})
		id := itemMap["Id"].(string)
		name := itemMap["Name"].(string)

		mapping := map[string]interface{}{
			"id":               id,
			"name":             name,
			"forward_mode":     itemMap["ForwardMode"],
			"forwarders":       itemMap["Forwarders"],
			"remark":           itemMap["Remark"],
			"caller_uid":       itemMap["CallerUid"],
			"create_timestamp": itemMap["CreateTimestamp"],
			"update_timestamp": itemMap["UpdateTimestamp"],
		}

		ids = append(ids, id)
		result = append(result, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return err
	}
	if err := d.Set("forward_domains", result); err != nil {
		return err
	}
	return nil
}
