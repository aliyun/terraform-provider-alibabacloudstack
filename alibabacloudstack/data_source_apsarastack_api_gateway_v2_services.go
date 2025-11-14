package alibabacloudstack

import (
	"fmt"
	"regexp"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackAPIGatewayV2Services() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackAPIGatewayV2ServicesRead,

		Schema: map[string]*schema.Schema{
			"gw_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"names": {
				Type:     schema.TypeList,
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
			"services": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"upstream_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"load_balance_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"real_service_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_group": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_nodes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"weight": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"enable": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"health_check_struct": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"health_path": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"http_statuses": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"timeout": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"health_interval": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"un_health_interval": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"http_successes": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"http_failures": {
										Type:     schema.TypeInt,
										Computed: true,
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

func dataSourceAlibabacloudStackAPIGatewayV2ServicesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	apiGatewayService := ApiGateWayV2Service{client}
	gwInstanceId := d.Get("gw_instance_id").(string)
	request := map[string]interface{}{
		"gwInstanceId": gwInstanceId,
		"current":      1,
		"size":         100, // Assuming a reasonable page size; adjust if needed
	}

	resp, err := client.DoTeaRequest("POST", "csb2", "2023-02-06", "ListServices", "/microservice/listServices", nil, nil, request)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_api_gateway_v2_services", "ListServices", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	// Parse response
	records, err := jsonpath.Get("$.data.records", resp)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_api_gateway_v2_services", "ListServices", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	// Prepare filtering maps
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

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	services := make([]map[string]interface{}, 0)
	for _, v := range records.([]interface{}) {
		record := v.(map[string]interface{})
		serviceId := record["serviceId"].(string)
		resourceId := fmt.Sprintf("%s:%s", gwInstanceId, serviceId)

		obj, err := apiGatewayService.DescribeApiGatewayV2Service(resourceId)
		if err != nil {
			return err
		}
		if len(idsMap) > 0 {
			if _, exists := idsMap[resourceId]; !exists {
				continue
			}
		}

		if nameRegex != nil && !nameRegex.MatchString(obj["name"].(string)) {
			continue
		}

		mapping := map[string]interface{}{
			"id":                resourceId,
			"service_id":        serviceId,
			"name":              obj["name"],
			"description":       obj["description"],
			"upstream_type":     obj["upstreamType"],
			"load_balance_type": obj["loadBalanceType"],
			"protocol":          obj["protocol"],
			"create_time":       obj["createTime"],
		}

		if serviceStructRaw, ok := obj["serviceStruct"]; ok && serviceStructRaw != nil {
			serviceStruct := serviceStructRaw.(map[string]interface{})
			nodes := make([]map[string]interface{}, 0)
			if nodeList, ok := serviceStruct["nodes"].([]interface{}); ok {
				for _, nodeRaw := range nodeList {
					node := nodeRaw.(map[string]interface{})
					nodes = append(nodes, map[string]interface{}{
						"ip":     node["ip"],
						"port":   node["port"],
						"weight": node["weight"],
						"enable": node["enable"],
					})
				}
			}
			mapping["service_nodes"] = nodes
			mapping["real_service_name"] = serviceStruct["realServiceName"]
			mapping["service_group"] = serviceStruct["group"]
			mapping["service_version"] = serviceStruct["version"]
			mapping["source_id"] = serviceStruct["sourceId"]
		}

		// Handle health_check_struct
		if healthCheckRaw, ok := obj["healthCheckStruct"]; ok && healthCheckRaw != nil {
			healthCheck := healthCheckRaw.(map[string]interface{})
			mapping["health_check_struct"] = []map[string]interface{}{
				{
					"type":               healthCheck["type"],
					"health_path":        healthCheck["healthPath"],
					"http_statuses":      healthCheck["httpStatuses"],
					"timeout":            healthCheck["timeout"],
					"health_interval":    healthCheck["healthInterval"],
					"un_health_interval": healthCheck["unHealthInterval"],
					"http_successes":     healthCheck["httpSuccesses"],
					"http_failures":      healthCheck["httpFailures"],
				},
			}
		}

		ids = append(ids, resourceId)
		names = append(names, obj["name"])
		services = append(services, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("services", services); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
