package alibabacloudstack

import (
	"regexp"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackMqttInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackMqttInstancesRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
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
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"max_conn": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_sub": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_up_tps": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_down_tps": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"independent_naming": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"store_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"store_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"endpoints": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"instance_status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"instance_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"namespace_rules_type": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"order_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sp_instance_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_tps": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackMqttInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// Build request query parameters
	requestQuery := map[string]interface{}{
		"Platform":       "onsConsole",
		"OnsRegionId":    client.RegionId,
		"Dauth_url_hash": "mqtt%2Fconsole%2Finstances",
		"PreventCache":   time.Now().UnixNano() / 1e6,
		"CurrentPage":    1,
		"PageSize":       100, // Assuming maximum page size to get all instances
	}

	// Call API to retrieve MQTT instances
	resp, err := client.DoTeaRequest("POST", "Ons-inner", "2018-02-05", "ConsoleMqttInstanceSearchByKeyword", "", nil, requestQuery, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_mqtt_instances", "ConsoleMqttInstanceSearchByKeyword", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	// Parse response data
	data, ok := resp["Data"]
	if !ok || data == nil {
		d.SetId("")
		return nil
	}

	var instancesData []map[string]interface{}
	if dataList, ok := data.(map[string]interface{})["data"].([]interface{}); ok {
		for _, item := range dataList {
			if instanceMap, ok := item.(map[string]interface{}); ok {
				instancesData = append(instancesData, instanceMap)
			}
		}
	} else {
		// If single instance returned (from detail API), wrap it in slice
		if instanceMap, ok := data.(map[string]interface{}); ok {
			instancesData = append(instancesData, instanceMap)
		}
	}

	// Filter by IDs if provided
	idsMap := getIdsStringFilter(d)

	// Compile name regex if provided
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		nameRegex = regexp.MustCompile(v.(string))
	}

	var filteredInstances []map[string]interface{}
	for _, instance := range instancesData {
		id := instance["instanceId"].(string)
		name := instance["instanceName"].(string)

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

	// Prepare results
	var ids []string
	var names []string
	var resultInstances []map[string]interface{}

	for _, instance := range filteredInstances {
		id := instance["instanceId"].(string)
		name := instance["instanceName"].(string)

		ids = append(ids, id)
		names = append(names, name)

		endpoints := make(map[string]string)
		if ep, ok := instance["endpoints"].(map[string]interface{}); ok {
			for k, v := range ep {
				if vs, ok := v.(string); ok {
					endpoints[k] = vs
				}
			}
		}

		resultInstance := map[string]interface{}{
			"instance_id":          id,
			"instance_name":        name,
			"max_conn":             instance["maxConn"],
			"max_sub":              instance["maxSub"],
			"max_up_tps":           instance["maxUpTps"],
			"max_down_tps":         instance["maxDownTps"],
			"independent_naming":   instance["independentNaming"],
			"store_instance_id":    instance["storeInstanceId"],
			"store_type":           instance["storeType"],
			"endpoints":            endpoints,
			"create_time":          instance["createTime"],
			"instance_status":      instance["instanceStatus"],
			"instance_type":        instance["instanceType"],
			"namespace_rules_type": instance["namespaceRulesType"],
			"order_id":             instance["orderId"],
			"sp_instance_type":     instance["spInstanceType"],
			"max_tps":              instance["maxTps"],
		}

		resultInstances = append(resultInstances, resultInstance)
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("ids", ids); err != nil {
		return err
	}
	if err := d.Set("names", names); err != nil {
		return err
	}
	if err := d.Set("instances", resultInstances); err != nil {
		return err
	}

	return nil
}
