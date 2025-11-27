package alibabacloudstack

import (
	"encoding/json"
	"regexp"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackUniversalDnsLines() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackUniversalDnsLinesRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				ForceNew:     true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
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
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"v6_addresses": {
							Type:     schema.TypeList,
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

func dataSourceAlibabacloudStackUniversalDnsLinesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

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

	var allLines []interface{}
	var ids []string
	pageNumber := 1
	pageSize := 100

	for {
		request := map[string]interface{}{
			"PageNumber": pageNumber,
			"PageSize":   pageSize,
		}

		if v, ok := d.GetOk("name"); ok {
			request["Name"] = v.(string)
		}

		respRaw, err := client.DoTeaRequest("POST", "UniversalDns", "2021-06-24", "DescribeUniversalLines", "", nil, nil, request)
		if err != nil {
			return err
		}

		if !respRaw["success"].(bool) {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_universal_dns_lines", "DescribeUniversalLines", errmsgs.AlibabacloudStackSdkGoERROR)
		}

		dataBytes, err := json.Marshal(respRaw["Data"])
		if err != nil {
			return err
		}

		var lines []map[string]interface{}
		if err := json.Unmarshal(dataBytes, &lines); err != nil {
			return err
		}

		totalItems, _ := respRaw["TotalItems"].(json.Number).Int64()
		currentPageSize := len(lines)

		for _, line := range lines {
			id := line["Id"].(string)
			name := line["Name"].(string)

			// Filter by ids
			if len(idsMap) > 0 {
				if _, exists := idsMap[id]; !exists {
					continue
				}
			}

			// Filter by name_regex
			if nameRegex != nil && !nameRegex.MatchString(name) {
				continue
			}

			allLines = append(allLines, map[string]interface{}{
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

		// Check pagination
		if currentPageSize < pageSize || int64(pageNumber*pageSize) >= totalItems {
			break
		}
		pageNumber++
	}
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("lines", allLines); err != nil {
		return err
	}

	return nil
}
