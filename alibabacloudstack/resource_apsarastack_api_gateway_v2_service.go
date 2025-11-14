package alibabacloudstack

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			"max_limit": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"timeout": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"is_scan_protect": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"original_sql": {
				Type:     schema.TypeString,
				Optional: true,
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
			"service_type": {
				Type:     schema.TypeInt,
				Required: true,
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
	service_type := d.Get("service_type").(int)
	request["name"] = d.Get("name")
	request["description"] = d.Get("description")
	request["serviceType"] = service_type
	request["gwInstanceId"] = d.Get("gw_instance_id")
	request["protocol"] = d.Get("protocol")
	switch service_type {
	case 0:
		request["healthCheckStruct"] = expandHealthCheckStruct(d.Get("health_check_struct"))
		request["upstreamType"] = d.Get("upstream_type")
		request["loadBalanceType"] = d.Get("load_balance_type")
		serviceStruct := make(map[string]interface{})
		serviceStruct["realServiceName"] = d.Get("real_service_name")
		serviceStruct["group"] = d.Get("service_group")
		serviceStruct["version"] = d.Get("service_version")
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
	case 1:
		request["sourceId"] = d.Get("source_id")
		request["maxLimit"] = d.Get("max_limit")
		request["timeout"] = d.Get("timeout")
		request["isScanProtect"] = d.Get("is_scan_protect")
		request["originalSql"] = d.Get("original_sql")

		inputParams := d.Get("sql_input_parameters").(*schema.Set).List()
		request["sqlInputParameters"] = buildSqlInputParams(inputParams)
		outputParams := d.Get("sql_output_parameters").(*schema.Set).List()
		request["sqlInputParameters"] = buildSqlOutputParams(outputParams)
	default:

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
	d.Set("service_type", object["serviceType"])
	d.Set("service_id", object["serviceId"])

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
		if v, ok := serviceStructMap["maxLimit"]; ok {
			d.Set("max_limit", v)
		}
		if v, ok := serviceStructMap["timeout"]; ok {
			d.Set("timeout", v)
		}
		if v, ok := serviceStructMap["isScanProtect"]; ok {
			d.Set("is_scan_protect", v)
		}
		if v, ok := serviceStructMap["originalSql"]; ok {
			d.Set("original_sql", v)
		}
		if v, ok := serviceStructMap["inputParameterList"]; ok {
			inputParams := readSqlInputParams(v.([]interface{}))
			d.Set("sql_input_parameters", inputParams)
		}
		if v, ok := serviceStructMap["outputParameterList"]; ok {
			outputParams := readSqlOutputParams(v.([]interface{}))
			d.Set("sql_output_parameters", outputParams)
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

	// Set all fields that can be updated
	service_type := d.Get("service_type").(int)
	request["name"] = d.Get("name")
	request["description"] = d.Get("description")
	request["serviceType"] = service_type
	request["protocol"] = d.Get("protocol")
	switch service_type {
	case 0:
		request["healthCheckStruct"] = expandHealthCheckStruct(d.Get("health_check_struct"))
		request["upstreamType"] = d.Get("upstream_type")
		request["loadBalanceType"] = d.Get("load_balance_type")
		serviceStruct := make(map[string]interface{})
		serviceStruct["realServiceName"] = d.Get("real_service_name")
		serviceStruct["group"] = d.Get("service_group")
		serviceStruct["version"] = d.Get("service_version")
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
	case 1:
		request["sourceId"] = d.Get("source_id")
		request["maxLimit"] = d.Get("max_limit")
		request["timeout"] = d.Get("timeout")
		request["isScanProtect"] = d.Get("is_scan_protect")
		request["originalSql"] = d.Get("original_sql")

		inputParams := d.Get("sql_input_parameters").(*schema.Set).List()
		request["sqlInputParameters"] = buildSqlInputParams(inputParams)
		outputParams := d.Get("sql_output_parameters").(*schema.Set).List()
		request["sqlInputParameters"] = buildSqlOutputParams(outputParams)
	default:

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

func buildSqlInputParams(input_parameters []interface{}) []map[string]interface{} {
	input_fields := make([]map[string]interface{}, 0)
	for _, v := range input_parameters {
		params := v.(map[string]interface{})
		input_fields = append(input_fields, map[string]interface{}{
			"originalName": params["original_name"],
			"targetName":   params["target_name"],
			"optional":     params["isoptional"],
			"sample":       params["sample"],
			"description":  params["description"],
			"fields":       []string{},
		})
	}
	input := map[string]interface{}{
		"originalName": "queryValues",
		"targetName":   "queryValues",
		"paramType":    "QueryValues",
		"optional":     false,
		"sample":       "",
		"description":  "查询条件字段集",
		"fields":       input_fields,
	}
	result := []map[string]interface{}{
		{
			"originalName": "requireTotalCount",
			"targetName":   "requireTotalCount",
			"paramType":    "java.lang.Boolean",
			"optional":     true,
		},
		{
			"originalName": "offset",
			"targetName":   "offset",
			"paramType":    "java.lang.Integer",
			"optional":     true,
		},
		{
			"originalName": "limit",
			"targetName":   "limit",
			"paramType":    "java.lang.Integer",
			"optional":     true,
		},
		input,
	}
	return result
}

func buildSqlOutputParams(output_parameters []interface{}) []map[string]interface{} {
	output_fields := []map[string]interface{}{
		{
			"originalName": "__csbRecordError",
			"targetName":   "errorMsg",
			"paramType":    "java.lang.String",
			"optional":     true,
		},
	}
	for _, v := range output_parameters {
		params := v.(map[string]interface{})
		output_fields = append(output_fields, map[string]interface{}{
			"originalName": params["original_name"],
			"targetName":   params["target_name"],
			"optional":     params["isoptional"],
			"sample":       params["sample"],
			"description":  params["description"],
			"fields":       []string{},
		})
	}
	oupput := map[string]interface{}{
		"originalName": "queryValues",
		"targetName":   "queryValues",
		"paramType":    "QueryValues",
		"optional":     false,
		"fields":       output_fields,
	}
	result := []map[string]interface{}{
		{
			"originalName": "resultCode",
			"targetName":   "resultCode",
			"paramType":    "java.lang.Integer",
			"optional":     false,
			"sample":       "",
		},
		{
			"originalName": "resultMsg",
			"targetName":   "resultMsg",
			"paramType":    "java.lang.String",
			"optional":     false,
			"sample":       "",
		},
		{
			"originalName": "resultCount",
			"targetName":   "resultCount",
			"paramType":    "java.lang.Integer",
			"optional":     true,
			"sample":       "",
		},
		{
			"originalName": "count",
			"targetName":   "count",
			"paramType":    "java.lang.Integer",
			"optional":     false,
			"sample":       "",
		},
		oupput,
	}
	return result
}

func readSqlInputParams(inputParams []interface{}) []map[string]interface{} {
	input := make([]map[string]interface{}, 0)
	for _, v := range inputParams {
		params := v.(map[string]interface{})
		if params["originalName"] == "queryValues" {
			fields := params["fields"].([]interface{})
			for _, field := range fields {
				fieldParams := field.(map[string]interface{})
				input = append(input, map[string]interface{}{
					"original_name": fieldParams["originalName"],
					"target_name":   fieldParams["targetName"],
					"isoptional":    fieldParams["optional"],
					"sample":        fieldParams["sample"],
					"description":   fieldParams["description"],
				})
			}
		}
	}
	return input
}

func readSqlOutputParams(outputParams []interface{}) []map[string]interface{} {
	output := make([]map[string]interface{}, 0)
	for _, v := range outputParams {
		params := v.(map[string]interface{})
		if params["originalName"] == "result" {
			fields := params["fields"].([]interface{})
			for _, v := range fields {
				fieldParams := v.(map[string]interface{})
				if fieldParams["originalName"] == "__csbRecordError" {
					continue
				}
				output = append(output, map[string]interface{}{
					"original_name": fieldParams["originalName"],
					"target_name":   fieldParams["targetName"],
					"isoptional":    fieldParams["optional"],
					"sample":        fieldParams["sample"],
					"description":   fieldParams["description"],
				})
			}
		}
	}
	return output
}
