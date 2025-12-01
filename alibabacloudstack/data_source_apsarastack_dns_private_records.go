package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackDnsPrivateRecords() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackDnsPrivateRecordsRead,

		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"records": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ttl": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"lba_strategy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remark": {
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
						"rdatas": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"lba_weight": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"line_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackDnsPrivateRecordsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// Parse filters
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		nameRegex = regexp.MustCompile(v.(string))
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

	zoneId := d.Get("zone_id").(string)
	pageNumber := 1
	pageSize := 100
	var allRecords []map[string]interface{}

	for {
		request := map[string]interface{}{
			"ZoneId":     zoneId,
			"PageNumber": pageNumber,
			"PageSize":   pageSize,
		}

		response, err := client.DoTeaRequest("POST", "CloudDns", "2021-06-24", "DescribePrivateZoneRecords", "", nil, nil, request)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_dns_private_records", "DescribePrivateZoneRecords", errmsgs.AlibabacloudStackSdkGoERROR)
		}

		records, ok := response["Data"]
		if !ok || records == nil {
			break
		}
		for _, v := range records.([]interface{}) {
			record := v.(map[string]interface{})
			id := record["Id"].(string)
			name := record["Name"].(string)
			resourceId := fmt.Sprintf("%s:%s", zoneId, id)
			// Apply name regex filter
			if nameRegex != nil && !nameRegex.MatchString(name) {
				continue
			}

			// Apply ids filter
			if len(idsMap) > 0 {
				if _, exists := idsMap[resourceId]; !exists {
					continue
				}
			}

			allRecords = append(allRecords, record)
		}

		totalItems, _ := response["TotalItems"].(json.Number).Int64()
		currentEnd := pageNumber * pageSize
		if int64(currentEnd) >= totalItems {
			break
		}
		pageNumber++
	}

	// Build final result
	ids := make([]string, 0)
	recordsList := make([]map[string]interface{}, 0)

	for _, record := range allRecords {
		id := record["Id"].(string)
		resourceId := fmt.Sprintf("%s:%s", zoneId, id)

		mapping := map[string]interface{}{
			"id":               id,
			"zone_id":          record["ZoneId"],
			"name":             record["Name"],
			"type":             record["Type"],
			"ttl":              record["Ttl"],
			"lba_strategy":     record["LbaStrategy"],
			"remark":           record["Remark"],
			"create_timestamp": record["CreateTimestamp"],
			"update_timestamp": record["UpdateTimestamp"],
			"line_ids":         record["LineIds"],
		}

		// Process RDatas
		rdatas := make([]map[string]interface{}, 0)
		if rdataList, ok := record["RDatas"].([]interface{}); ok {
			for _, rdataItem := range rdataList {
				if rdata, ok := rdataItem.(map[string]interface{}); ok {
					rdataMap := map[string]interface{}{
						"value": rdata["Value"],
					}
					lbaWeight, ok := rdata["LbaWeight"]
					if ok {
						rdataMap["lba_weight"] = lbaWeight
					}
					rdatas = append(rdatas, rdataMap)
				}
			}
		}
		mapping["rdatas"] = rdatas

		recordsList = append(recordsList, mapping)
		ids = append(ids, resourceId)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("records", recordsList); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
