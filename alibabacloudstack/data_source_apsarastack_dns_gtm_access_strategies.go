package alibabacloudstack

import (
	"fmt"
	"regexp"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackDnsGtmAccessStrategies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackDnsGtmAccessStrategiesRead,

		Schema: map[string]*schema.Schema{
			"gtm_instance_id": {
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
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"strategies": {
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
						"gtm_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"default_gtm_address_pool_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"default_gtm_address_pool_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"default_gtm_address_pool_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"default_min_available_addr_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"default_available_addr_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"failover_gtm_address_pool_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"failover_gtm_address_pool_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"failover_gtm_address_pool_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"failover_min_available_addr_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"failover_available_addr_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"specified_gtm_address_pool": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"switch_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"in_use_gtm_address_pool_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"in_use_gtm_address_pool_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"line_ids": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackDnsGtmAccessStrategiesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// Build request parameters
	request := make(map[string]interface{})
	request["GtmInstanceId"] = d.Get("gtm_instance_id")
	request["PageNumber"] = 1
	request["PageSize"] = 100 // Set a large enough page size to get all data at once

	action := "DescribeDnsGtmAccessStrategies"
	response, err := client.DoTeaRequest("POST", "CloudDns", "2021-06-24", action, "", nil, nil, request)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	// Parse response data
	data, ok := response["Data"].([]interface{})
	if !ok {
		data = []interface{}{}
	}

	// Filter by ids if provided
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
		nameRegex = regexp.MustCompile(v.(string))
	}

	// Process and filter results
	var filteredResults []interface{}
	for _, item := range data {
		obj, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		id := fmt.Sprintf("%s:%s", obj["GtmInstanceId"].(string), obj["Id"].(string))
		name := obj["Name"].(string)

		// Apply id filter
		if len(idsMap) > 0 {
			if _, exists := idsMap[id]; !exists {
				continue
			}
		}

		// Apply name regex filter
		if nameRegex != nil && !nameRegex.MatchString(name) {
			continue
		}

		filteredResults = append(filteredResults, obj)
	}

	// Prepare final result structures
	strategyMaps := make([]map[string]interface{}, 0)
	ids := make([]string, 0)

	for _, item := range filteredResults {
		obj := item.(map[string]interface{})

		id := fmt.Sprintf("%s:%s", obj["GtmInstanceId"].(string), obj["Id"].(string))
		lineIdsInterface, _ := obj["LineIds"].([]interface{})
		lineIds := make([]string, 0)
		for _, lineId := range lineIdsInterface {
			if lid, ok := lineId.(string); ok {
				lineIds = append(lineIds, lid)
			}
		}
		mapping := map[string]interface{}{
			"id":                              id,
			"name":                            obj["Name"],
			"gtm_instance_id":                 obj["GtmInstanceId"],
			"default_gtm_address_pool_type":   obj["DefaultGtmAddressPoolType"],
			"default_gtm_address_pool_id":     obj["DefaultGtmAddressPoolId"],
			"default_gtm_address_pool_name":   obj["DefaultGtmAddressPoolName"],
			"default_min_available_addr_num":  obj["DefaultMinAvailableAddrNum"],
			"default_available_addr_num":      obj["DefaultAvailableAddrNum"],
			"failover_gtm_address_pool_id":    obj["FailoverGtmAddressPoolId"],
			"failover_gtm_address_pool_name":  obj["FailoverGtmAddressPoolName"],
			"failover_gtm_address_pool_type":  obj["FailoverGtmAddressPoolType"],
			"failover_min_available_addr_num": obj["FailoverMinAvailableAddrNum"],
			"failover_available_addr_num":     obj["FailoverAvailableAddrNum"],
			"specified_gtm_address_pool":      obj["SpecifiedGtmAddressPool"],
			"switch_mode":                     obj["SwitchMode"],
			"in_use_gtm_address_pool_id":      obj["InUseGtmAddressPoolId"],
			"in_use_gtm_address_pool_name":    obj["InUseGtmAddressPoolName"],
			"line_ids":                        obj["LineIds"],
		}
		strategyMaps = append(strategyMaps, mapping)
		ids = append(ids, id)
	}

	// Set computed fields
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("strategies", strategyMaps); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
