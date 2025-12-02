package alibabacloudstack

import (
	"regexp"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackDnsGtmAddressPools() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackDnsGtmAddressPoolsRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A list of address pool IDs.",
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  "A regex string to filter results by address pool name.",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Type of the address pool. Valid values: A, AAAA, CNAME.",
			},
			"lba_strategy": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Load balancing strategy of the address pool. Valid values: ALL_RR, RATIO.",
			},
			"address_pools": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the address pool.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the address pool.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the address pool. Valid values: A, AAAA, CNAME.",
						},
						"lba_strategy": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The load balancing strategy of the address pool. Valid values: ALL_RR, RATIO.",
						},
						"create_timestamp": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The timestamp when the address pool was created (in seconds).",
						},
						"update_timestamp": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The timestamp when the address pool was last updated (in seconds).",
						},
						"addrs": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The value of the address.",
									},
									"mode": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The mode of the address. Valid values: SMART, ONLINE, OFFLINE.",
									},
									"lba_weight": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The weight of the address for load balancing.",
									},
									"available": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Indicates whether the address is available.",
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

func dataSourceAlibabacloudStackDnsGtmAddressPoolsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// Build request parameters
	request := map[string]interface{}{
		"PageNumber": 1,
		"PageSize":   100,
	}
	// Handle ids filter
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	// Handle name_regex filter
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return errmsgs.WrapError(err)
		}
		nameRegex = r
	}
	resp, err := client.DoTeaRequest("POST", "CloudDns", "2021-06-24", "DescribeDnsGtmAddressPools", "", nil, nil, request)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	if success, ok := resp["success"].(bool); !ok || !success {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_dns_gtm_addresspools", "DescribeDnsGtmAddressPools", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	// Filter results based on ids and name_regex
	var filteredAddressPools []interface{}
	for _, item := range resp["Data"].([]interface{}) {
		pool := item.(map[string]interface{})

		// Apply ids filter
		if len(idsMap) > 0 {
			id := pool["Id"].(string)
			if _, exists := idsMap[id]; !exists {
				continue
			}
		}

		// Apply name_regex filter
		if nameRegex != nil {
			name := pool["Name"].(string)
			if !nameRegex.MatchString(name) {
				continue
			}
		}

		if v, ok := d.GetOk("type"); ok {
			if v != pool["Type"] {
				continue
			}
		}

		if v, ok := d.GetOk("lba_strategy"); ok {
			if v != pool["LbaStrategy"] {
				continue
			}
		}

		filteredAddressPools = append(filteredAddressPools, pool)
	}

	// Prepare final result
	addressPools := make([]map[string]interface{}, 0)
	ids := make([]string, 0)

	for _, item := range filteredAddressPools {
		pool := item.(map[string]interface{})
		// Process Addrs
		addrsList := make([]map[string]interface{}, 0)
		if addrsRaw, ok := pool["Addrs"].([]interface{}); ok {
			for _, addrRaw := range addrsRaw {
				addr := addrRaw.(map[string]interface{})
				addrMap := make(map[string]interface{})
				addrMap["value"] = addr["Value"]
				addrMap["mode"] = addr["Mode"]
				if lbaWeight, ok := addr["LbaWeight"].(float64); ok {
					addrMap["lba_weight"] = int(lbaWeight)
				}
				if available, ok := addr["Available"].(bool); ok {
					addrMap["available"] = available
				}
				addrsList = append(addrsList, addrMap)
			}
		}
		id := pool["Id"].(string)

		mapping := map[string]interface{}{
			"id":               id,
			"name":             pool["Name"],
			"type":             pool["Type"],
			"lba_strategy":     pool["LbaStrategy"],
			"create_timestamp": pool["CreateTimestamp"],
			"update_timestamp": pool["UpdateTimestamp"],
			"addrs":            addrsList,
		}

		addressPools = append(addressPools, mapping)
		ids = append(ids, id)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("address_pools", addressPools); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}
	return nil
}
