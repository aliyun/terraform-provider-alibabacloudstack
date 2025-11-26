package alibabacloudstack

import (
	"regexp"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackUniversalDnsDomains() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackUniversalDnsDomainsRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A list of domain IDs that the data source will be filtered by.",
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  "A regex string to filter results by domain name.",
			},
			"names": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A list of domain names corresponding to the domains found.",
			},
			"domains": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of universal DNS domains matching the criteria.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the universal DNS domain.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the universal DNS domain.",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The remark or description of the domain.",
						},
						"record_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The total number of DNS records in the domain.",
						},
						"create_timestamp": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The timestamp (in seconds) when the domain was created.",
						},
						"update_timestamp": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The timestamp (in seconds) when the domain was last updated.",
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackUniversalDnsDomainsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// Build ids filter map
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	// Compile name regex if provided
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return errmsgs.WrapError(err)
		}
		nameRegex = r
	}

	// Call API to get all domains
	request := make(map[string]interface{})
	request["PageNumber"] = 1
	request["PageSize"] = 100 // Max page size to reduce pagination needs

	var allDomains []map[string]interface{}
	for {
		response, err := client.DoTeaRequest("POST", "UniversalDns", "2021-06-24", "DescribeUniversalZones", "", nil, nil, request)
		if err != nil {
			return errmsgs.WrapError(err)
		}

		data, ok := response["Data"].([]interface{})
		if !ok {
			break
		}

		totalItems, _ := response["TotalItems"].(float64)
		pageSize, _ := response["PageSize"].(float64)
		pageNumber, _ := response["PageNumber"].(float64)

		for _, item := range data {
			domain, ok := item.(map[string]interface{})
			if !ok {
				continue
			}

			// Apply id filtering
			if len(idsMap) > 0 {
				id, ok := domain["Id"].(string)
				if !ok || idsMap[id] == "" {
					continue
				}
			}

			// Apply name regex filtering
			if nameRegex != nil {
				name, ok := domain["Name"].(string)
				if !ok || !nameRegex.MatchString(name) {
					continue
				}
			}

			allDomains = append(allDomains, domain)
		}

		// Check for more pages
		if int(pageNumber)*int(pageSize) >= int(totalItems) {
			break
		}
		request["PageNumber"] = int(pageNumber) + 1
	}

	// Prepare output structures
	domains := make([]map[string]interface{}, 0)
	names := make([]string, 0)
	ids := make([]string, 0)

	for _, domain := range allDomains {
		mapping := map[string]interface{}{
			"id":               domain["Id"],
			"name":             domain["Name"],
			"remark":           domain["Remark"],
			"record_count":     domain["RecordCount"],
			"create_timestamp": domain["CreateTimestamp"],
			"update_timestamp": domain["UpdateTimestamp"],
		}

		domains = append(domains, mapping)
		if name, ok := domain["Name"].(string); ok {
			names = append(names, name)
		}
		if id, ok := domain["Id"].(string); ok {
			ids = append(ids, id)
		}
	}

	// Set computed fields
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("domains", domains); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
