package alibabacloudstack

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackAPIGatewayV2Service() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"gw_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"upstream_type": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"load_balance_type": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"real_service_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_group": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_group": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"service_source_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"dns", "ip"}, false),
			},
			"service_nodes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:     schema.TypeString,
							Required: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"weight": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"enable": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"health_check_struct": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"health_path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"http_statuses": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"timeout": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"health_interval": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"un_health_interval": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"http_successes": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"http_failures": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"service_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sql_input_parameters": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"original_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"target_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"isoptional": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"description": {
							Type:     schema.TypeString,
							Required: true,
						},
						"sample": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"sql_output_parameters": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"original_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"target_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"param_type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "java.lang.String",
						},
						"isoptional": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"description": {
							Type:     schema.TypeString,
							Required: true,
						},
						"sample": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackAPIGatewayV2ServiceCreate, resourceAlibabacloudStackAPIGatewayV2ServiceRead, resourceAlibabacloudStackAPIGatewayV2ServiceUpdate, resourceAlibabacloudStackAPIGatewayV2ServiceDelete)
	return resource
}

func resourceAlibabacloudStackAPIGatewayV2ServiceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := make(map[string]interface{})
	request["name"] = d.Get("name")
	request["description"] = d.Get("description")
	request["serviceType"] = 0
	request["gwInstanceId"] = d.Get("gw_instance_id")
	request["protocol"] = d.Get("protocol")
	request["healthCheckStruct"] = expandHealthCheckStruct(d.Get("health_check_struct"))
	request["upstreamType"] = d.Get("upstream_type")
	request["loadBalanceType"] = d.Get("load_balance_type")
	serviceStruct := make(map[string]interface{})
	serviceStruct["realServiceName"] = d.Get("real_service_name")
	serviceStruct["group"] = d.Get("service_group")
	serviceStruct["version"] = d.Get("service_version")
	serviceStruct["sourceId"] = d.Get("source_id")
	nodes := d.Get("service_nodes").([]interface{})
	nodeList := make([]map[string]interface{}, 0, len(nodes))
	for _, node := range nodes {
		n := node.(map[string]interface{})
		nodeMap := make(map[string]interface{})
		nodeMap["ip"] = n["ip"]
		nodeMap["port"] = n["port"]
		nodeMap["weight"] = n["weight"]
		nodeMap["enable"] = n["enable"]
		nodeList = append(nodeList, nodeMap)
	}
	serviceStruct["nodes"] = nodeList
	request["serviceStruct"] = serviceStruct
	healthCheckStruct := expandHealthCheckStruct(d.Get("health_check_struct"))
	if healthCheckStruct != nil {
		request["healthCheckStruct"] = healthCheckStruct
		request["isOpenHealthCheck"] = true
	} else {
		request["isOpenHealthCheck"] = false
	}
	if v, ok := d.GetOk("service_source_type"); ok {
		switch v.(string) {
		case "dns":
			request["serviceSourceType"] = 8
		case "ip":
			request["serviceSourceType"] = 7
		default:
			return errmsgs.WrapError(fmt.Errorf("invalid service_source_type: %s", v.(string)))
		}
	}
	resp, err := client.DoTeaRequest("POST", "csb2", "2023-02-06", "CreateService", "/microservice/createService", nil, nil, request)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_api_gateway_v2_service", "CreateService", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	serviceId, ok := resp["data"].(string)
	if !ok {
		return errmsgs.WrapError(fmt.Errorf("failed to retrieve serviceId from response data"))
	}

	gwInstanceId := d.Get("gw_instance_id").(string)
	d.SetId(fmt.Sprintf("%s:%s", gwInstanceId, serviceId))

	return nil
}

func expandHealthCheckStruct(i interface{}) map[string]interface{} {
	if i == nil || len(i.([]interface{})) == 0 {
		return nil
	}
	hc := i.([]interface{})[0].(map[string]interface{})
	result := make(map[string]interface{})
	result["type"] = hc["type"]
	result["healthPath"] = hc["health_path"]
	result["httpStatuses"] = hc["http_statuses"]
	result["timeout"] = hc["timeout"]
	result["healthInterval"] = hc["health_interval"]
	result["unHealthInterval"] = hc["un_health_interval"]
	result["httpSuccesses"] = hc["http_successes"]
	result["httpFailures"] = hc["http_failures"]
	return result
}

func resourceAlibabacloudStackAPIGatewayV2ServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	apiGatewayService := ApiGateWayV2Service{client}

	object, err := apiGatewayService.DescribeApiGatewayV2Service(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	gwInstanceId := parts[0]
	d.Set("gw_instance_id", gwInstanceId)
	d.Set("name", object["name"])
	d.Set("description", object["description"])
	d.Set("upstream_type", object["upstreamType"])
	d.Set("load_balance_type", object["loadBalanceType"])
	d.Set("protocol", object["protocol"])
	d.Set("service_id", object["serviceId"])
	switch object["serviceSourceType"].(int) {
	case 8:
		d.Set("service_source_type", "dns")
	case 7:
		d.Set("service_source_type", "ip")
	}

	if serviceStruct, ok := object["serviceStruct"]; ok && serviceStruct != nil {
		serviceStructMap := serviceStruct.(map[string]interface{})

		if nodes, ok := serviceStructMap["nodes"].([]interface{}); ok && len(nodes) > 0 {
			nodeList := make([]map[string]interface{}, 0, len(nodes))
			for _, node := range nodes {
				nodeMap := node.(map[string]interface{})
				nodeData := map[string]interface{}{
					"ip":     nodeMap["ip"],
					"port":   nodeMap["port"],
					"weight": nodeMap["weight"],
					"enable": nodeMap["enable"],
				}
				nodeList = append(nodeList, nodeData)
			}
			d.Set("service_nodes", nodeList)
		}
		if v, ok := serviceStructMap["realServiceName"]; ok {
			d.Set("real_service_name", v)
		}
		if v, ok := serviceStructMap["group"]; ok {
			d.Set("service_group", v)
		}
		if v, ok := serviceStructMap["version"]; ok {
			d.Set("service_version", v)
		}
		if v, ok := serviceStructMap["sourceId"]; ok {
			d.Set("source_id", v)
		}
		if v, ok := serviceStructMap["sourceGroup"]; ok {
			d.Set("source_group", v)
		}
	}

	if healthCheckStruct, ok := object["healthCheckStruct"]; ok && healthCheckStruct != nil {
		healthCheckMap := healthCheckStruct.(map[string]interface{})
		healthData := make([]map[string]interface{}, 0)
		healthDatum := map[string]interface{}{
			"type":               healthCheckMap["type"],
			"health_path":        healthCheckMap["healthPath"],
			"http_statuses":      healthCheckMap["httpStatuses"],
			"timeout":            healthCheckMap["timeout"],
			"health_interval":    healthCheckMap["healthInterval"],
			"un_health_interval": healthCheckMap["unHealthInterval"],
			"http_successes":     healthCheckMap["httpSuccesses"],
			"http_failures":      healthCheckMap["httpFailures"],
		}
		healthData = append(healthData, healthDatum)
		d.Set("health_check_struct", healthData)
	}

	return nil
}

func resourceAlibabacloudStackAPIGatewayV2ServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	gwInstanceId := parts[0]
	serviceId := parts[1]

	// If it is a new resource, do nothing
	if d.IsNewResource() {
		return nil
	}

	request := make(map[string]interface{})
	request["serviceId"] = serviceId
	request["gwInstanceId"] = gwInstanceId

	request["name"] = d.Get("name")
	request["description"] = d.Get("description")
	request["serviceType"] = 0
	request["protocol"] = d.Get("protocol")
	request["healthCheckStruct"] = expandHealthCheckStruct(d.Get("health_check_struct"))
	request["upstreamType"] = d.Get("upstream_type")
	request["loadBalanceType"] = d.Get("load_balance_type")
	serviceStruct := make(map[string]interface{})
	serviceStruct["realServiceName"] = d.Get("real_service_name")
	serviceStruct["group"] = d.Get("service_group")
	serviceStruct["version"] = d.Get("service_version")
	serviceStruct["sourceId"] = d.Get("source_id")
	nodes := d.Get("service_nodes").([]interface{})
	nodeList := make([]map[string]interface{}, 0, len(nodes))
	for _, node := range nodes {
		n := node.(map[string]interface{})
		nodeMap := make(map[string]interface{})
		nodeMap["ip"] = n["ip"]
		nodeMap["port"] = n["port"]
		nodeMap["weight"] = n["weight"]
		nodeMap["enable"] = n["enable"]
		nodeList = append(nodeList, nodeMap)
	}
	serviceStruct["nodes"] = nodeList
	request["serviceStruct"] = serviceStruct
	healthCheckStruct := expandHealthCheckStruct(d.Get("health_check_struct"))
	if healthCheckStruct != nil {
		request["healthCheckStruct"] = healthCheckStruct
		request["isOpenHealthCheck"] = true
	} else {
		request["isOpenHealthCheck"] = false
	}
	if v, ok := d.GetOk("service_source_type"); ok {
		switch v.(string) {
		case "dns":
			request["serviceSourceType"] = 8
		case "ip":
			request["serviceSourceType"] = 7
		default:
			return errmsgs.WrapError(fmt.Errorf("invalid service_source_type: %s", v.(string)))
		}
	}
	_, err = client.DoTeaRequest("POST", "csb2", "2023-02-06", "ModifyService", "/microservice/modifyService", nil, nil, request)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_api_gateway_v2_service", "ModifyService", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return nil
}

func resourceAlibabacloudStackAPIGatewayV2ServiceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	gwInstanceId := parts[0]
	serviceId := parts[1]

	reqQuery := map[string]interface{}{
		"serviceId":    serviceId,
		"gwInstanceId": gwInstanceId,
	}

	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		raw, err := client.DoTeaRequest("POST", "csb2", "2023-02-06", "DeleteService", "/microservice/deleteService", nil, nil, reqQuery)
		if err != nil {
			err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "DeleteService", errmsgs.AlibabacloudStackSdkGoERROR, "")
			return resource.RetryableError(err)
		}
		// No need to check response data as per documentation
		_ = raw
		return nil
	})
	if err != nil {
		return errmsgs.WrapError(err)
	}
	return nil
}
