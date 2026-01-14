package alibabacloudstack

import (
	"regexp"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackDnsLines() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackDnsLinesRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A list of Universal DNS line IDs.",
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  "A regex string to filter results by the Universal DNS line name.",
			},
			"lines": {
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
						"priority": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"v4_addresses": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"v6_addresses": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
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

func dataSourceAlibabacloudStackDnsLinesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// Build request parameters
	request := make(map[string]interface{})
	request["PageNumber"] = 1
	request["PageSize"] = 100

	// Call API to get the list of lines
	response, err := client.DoTeaRequest("POST", "CloudDns", "2021-06-24", "DescribeGlobalLines", "", nil, nil, request)
	if err != nil {
		return err
	}

	// Check if the response contains data
	data, ok := response["Data"].([]interface{})
	if !ok {
		// Return empty list if no data found
		d.SetId("")
		return nil
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

	// Filter and process the returned data
	lines := make([]map[string]interface{}, 0)
	ids := make([]string, 0)

	for _, item := range data {
		line, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		id, ok := line["Id"].(string)
		if !ok {
			continue
		}

		name, ok := line["Name"].(string)
		if !ok {
			continue
		}

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
		// Append processed line to result
		lines = append(lines, map[string]interface{}{
			"id":               id,
			"name":             name,
			"priority":         line["Priority"],
			"v4_addresses":     line["V4Addresses"],
			"v6_addresses":     line["V6Addresses"],
			"create_timestamp": line["CreateTimestamp"],
			"update_timestamp": line["UpdateTimestamp"],
		})
		ids = append(ids, id)
	}

	// Set computed fields
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("lines", lines); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
