package alibabacloudstack

import (
	"encoding/json"
	"regexp"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackDnsPrivateDomains() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackDnsPrivateDomainsRead,

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
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"domains": {
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
						"remark": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"caller_uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"record_count": {
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
						"gtm_instance_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"region_and_vpcs": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"vpcs": {
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
											},
										},
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

func dataSourceAlibabacloudStackDnsPrivateDomainsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	idsMap := getIdsStringFilter(d)

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return errmsgs.WrapError(err)
		}
		nameRegex = r
	}

	request := make(map[string]interface{})
	if v, ok := d.GetOk("id"); ok {
		request["Id"] = v.(string)
	}
	if v, ok := d.GetOk("name"); ok {
		request["Name"] = v.(string)
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		request["VpcId"] = v.(string)
	}
	request["RegionId"] = client.RegionId
	request["PageNumber"] = 1
	request["PageSize"] = 100

	action := "DescribePrivateZones"
	response, err := client.DoTeaRequest("POST", "CloudDns", "2021-06-24", action, "", nil, nil, request)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_dns_private_domains", action, errmsgs.AlibabacloudStackSdkGoERROR)
	}

	if success, ok := response["success"].(bool); !ok || !success {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_dns_private_domains", action, errmsgs.AlibabacloudStackSdkGoERROR)
	}

	data, err := json.Marshal(response["Data"])
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_dns_private_domains", action, errmsgs.AlibabacloudStackSdkGoERROR)
	}

	var items []interface{}
	if err := json.Unmarshal(data, &items); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_dns_private_domains", action, errmsgs.AlibabacloudStackSdkGoERROR)
	}

	var filteredItems []interface{}
	for _, item := range items {
		itemMap := item.(map[string]interface{})
		id := itemMap["Id"].(string)
		name := itemMap["Name"].(string)

		// Filter by ids
		if len(idsMap) > 0 {
			if _, ok := idsMap[id]; !ok {
				continue
			}
		}

		// Filter by name_regex
		if nameRegex != nil && !nameRegex.MatchString(name) {
			continue
		}

		filteredItems = append(filteredItems, item)
	}

	// Prepare results
	var ids []string
	var domains []map[string]interface{}

	for _, item := range filteredItems {
		itemMap := item.(map[string]interface{})
		id := itemMap["Id"].(string)
		name := itemMap["Name"].(string)
		ids = append(ids, id)

		domain := map[string]interface{}{
			"id":                 id,
			"name":               name,
			"remark":             itemMap["Remark"],
			"caller_uid":         itemMap["CallerUid"],
			"record_count":       itemMap["RecordCount"],
			"create_timestamp":   itemMap["CreateTimestamp"],
			"update_timestamp":   itemMap["UpdateTimestamp"],
			"gtm_instance_count": itemMap["GtmInstanceCount"],
		}

		// Process region_and_vpcs
		var regionAndVpcs []map[string]interface{}
		if ravList, ok := itemMap["RegionAndVpcs"].([]interface{}); ok {
			for _, rav := range ravList {
				ravMap := rav.(map[string]interface{})
				regionId := ravMap["RegionId"].(string)

				var vpcs []map[string]interface{}
				if vpcList, ok := ravMap["Vpcs"].([]interface{}); ok {
					for _, vpc := range vpcList {
						vpcMap := vpc.(map[string]interface{})
						vpcs = append(vpcs, map[string]interface{}{
							"id":   vpcMap["Id"],
							"name": vpcMap["Name"],
						})
					}
				}

				regionAndVpcs = append(regionAndVpcs, map[string]interface{}{
					"region_id": regionId,
					"vpcs":      vpcs,
				})
			}
		}
		domain["region_and_vpcs"] = regionAndVpcs

		domains = append(domains, domain)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("domains", domains); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
