package alibabacloudstack

import (
	"fmt"
	"log"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackAPIGatewayV2Mcpserver() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"gw_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"DIRECT_ROUTE", "DATABASE", "OPEN_API"}, false),
			},
			"service": {
				Type:     schema.TypeString,
				Required: true,
			},
			"domains": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"consumer_auth": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"direct_route_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "When the type is `DIRECT_ROUTE`, this field is required.",
			},
			"direct_route_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "When the type is `DIRECT_ROUTE`, this field is required.",
				ValidateFunc: validation.StringInSlice([]string{"streamable", "sse"}, false),
			},
			"services": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"weight": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"db_host": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "When the type is `DATABASE`, this field is required.",
			},
			"db_port": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "When the type is `DATABASE`, this field is required.",
			},
			"db_username": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "When the type is `DATABASE`, this field is required.",
			},
			"db_password": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "When the type is `DATABASE`, this field is required.",
			},
			"db_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "When the type is `DATABASE`, this field is required.",
			},
			"other_params": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "When the type is `DATABASE`, this field is required.",
			},
			"db_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "When the type is `DATABASE`, this field is required.",
				ValidateFunc: validation.StringInSlice([]string{"MYSQL", "POSTGRESQL", "CLICKHOUSE"}, false),
			},
			"raw_configurations": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackAPIGatewayV2McpserverCreate, resourceAlibabacloudStackAPIGatewayV2McpserverRead, resourceAlibabacloudStackAPIGatewayV2McpserverUpdate, resourceAlibabacloudStackAPIGatewayV2McpserverDelete)
	return resource
}

func resourceAlibabacloudStackAPIGatewayV2McpserverCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := make(map[string]interface{})

	request["name"] = d.Get("name")
	request["type"] = d.Get("type")
	request["service"] = d.Get("service")
	request["gwInstanceId"] = d.Get("gw_instance_id")
	request["description"] = d.Get("description")
	request["domains"] = d.Get("domains").(*schema.Set).List()
	request["consumerAuth"] = d.Get("consumer_auth")

	if d.Get("type").(string) == "DIRECT_ROUTE" {
		directRouteConfig := make(map[string]string)
		directRouteConfig["path"] = d.Get("direct_route_path").(string)
		directRouteConfig["type"] = d.Get("direct_route_type").(string)
		request["directRouteConfig"] = directRouteConfig
		request["serviceProtocol"] = "MCP"
	}
	apiGatewayV2Service := ApiGateWayV2Service{client}
	services, err := apiGatewayV2Service.DescribeService(d.Get("service").(string), d.Get("gw_instance_id").(string))
	if err != nil {
		return errmsgs.WrapError(err)
	}
	nodes, err := jsonpath.Get("$.serviceStruct.nodes", services)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	for _, v := range nodes.([]interface{}) {
		node := v.(map[string]interface{})
		if node["enable"].(bool) {
			service := make(map[string]interface{})
			service["name"] = d.Get("service")
			service["port"] = node["port"]
			service["version"] = "1.0"
			service["weight"] = node["weight"]
			request["services"] = []interface{}{service}
			break
		}

	}
	if d.Get("consumer_auth").(bool) {
		request["consumerAuth"] = true
		request["consumerAuthInfo"] = map[string]interface{}{
			"enable":           true,
			"type":             "API_KEY",
			"allowedConsumers": []string{},
		}
	} else {
		request["consumerAuth"] = false
	}

	if d.Get("type").(string) == "DATABASE" {
		dbConfig := make(map[string]interface{})
		dbConfig["host"] = d.Get("db_host")
		dbConfig["port"] = d.Get("db_port")
		dbConfig["username"] = d.Get("db_username")
		dbConfig["password"] = d.Get("db_password")
		dbConfig["dbname"] = d.Get("db_name")
		dbConfig["otherParams"] = d.Get("other_params")
		request["dbConfig"] = dbConfig
		request["dbType"] = d.Get("db_type")
	}

	resp, err := client.DoTeaRequest("POST", "csb2", "2023-02-06", "CreateMcpServer", "/mcpServer/createMcpServer", nil, nil, request)
	if err != nil {
		return err
	}

	// Get the name from response data
	name, ok := resp["data"].(string)
	if !ok {
		return fmt.Errorf("failed to get name from response data")
	}

	gwInstanceId := d.Get("gw_instance_id").(string)
	// Set ID: {gwInstanceId}:{name}
	d.SetId(fmt.Sprintf("%s:%s", gwInstanceId, name))

	return nil
}

func resourceAlibabacloudStackAPIGatewayV2McpserverRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	apiGatewayV2Service := ApiGateWayV2Service{client}
	object, err := apiGatewayV2Service.DescribeMcpserver(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_api_gateway_v2_mcpserver apiGatewayV2Service.DescribeMcpserver Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("name", object["name"])
	d.Set("description", object["description"])
	d.Set("type", object["type"])
	d.Set("domains", object["domains"])
	d.Set("gw_instance_id", object["gwInstanceId"])
	d.Set("raw_configurations", object["rawConfigurations"])

	// Set consumer_auth from consumerAuthInfo.enable
	if consumerAuthInfo, ok := object["consumerAuthInfo"].(map[string]interface{}); ok {
		if enable, ok := consumerAuthInfo["enable"]; ok {
			d.Set("consumer_auth", enable)
		}
	}

	// Set direct_route_config fields for DIRECT_ROUTE type
	if directRouteConfig, ok := object["directRouteConfig"].(map[string]interface{}); ok {
		if path, ok := directRouteConfig["path"]; ok {
			d.Set("direct_route_path", path)
		}
		if transportType, ok := directRouteConfig["transportType"]; ok {
			d.Set("direct_route_type", transportType)
		}
	}

	// Set services
	if services, ok := object["services"].([]interface{}); ok {
		serviceList := make([]map[string]interface{}, 0)
		for _, s := range services {
			if serviceMap, ok := s.(map[string]interface{}); ok {
				service := make(map[string]interface{})
				if name, ok := serviceMap["name"]; ok {
					service["name"] = name
					d.Set("service", name)
				}
				if port, ok := serviceMap["port"]; ok {
					service["port"] = port
				}
				if version, ok := serviceMap["version"]; ok {
					service["version"] = version
				}
				if weight, ok := serviceMap["weight"]; ok {
					service["weight"] = weight
				}
				serviceList = append(serviceList, service)
			}
		}
		d.Set("services", serviceList)
	}

	// Set database config fields for DATABASE type
	if dbConfig, ok := object["dbConfig"].(map[string]interface{}); ok {
		if dbHost, ok := dbConfig["host"]; ok {
			d.Set("db_host", dbHost)
		}
		if port, ok := dbConfig["port"]; ok {
			d.Set("db_port", port)
		}
		if username, ok := dbConfig["username"]; ok {
			d.Set("db_username", username)
		}
		if password, ok := dbConfig["password"]; ok {
			d.Set("db_password", password)
		}
		if dbname, ok := dbConfig["dbname"]; ok {
			d.Set("db_name", dbname)
		}
		if otherParams, ok := dbConfig["otherParams"]; ok {
			d.Set("other_params", otherParams)
		}
	}

	if dbType, ok := object["dbType"]; ok {
		d.Set("db_type", dbType)
	}

	return nil
}

func resourceAlibabacloudStackAPIGatewayV2McpserverUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// If it is a new resource, do nothing
	if d.IsNewResource() {
		return nil
	}

	request := make(map[string]interface{})

	// Always required fields
	request["name"] = d.Get("name")
	request["type"] = d.Get("type")
	request["gwInstanceId"] = d.Get("gw_instance_id")

	if d.HasChanges("description", "domains", "direct_route_path", "direct_route_type", "db_host", "db_port", "db_username",
		"db_password", "db_name", "other_params", "db_type", "service", "consumer_auth") {
		request["domains"] = d.Get("domains").(*schema.Set).List()
		request["consumerAuth"] = d.Get("consumer")
		request["service"] = d.Get("service")
		request["description"] = d.Get("description")
		apiGatewayV2Service := ApiGateWayV2Service{client}
		services, err := apiGatewayV2Service.DescribeService(d.Get("service").(string), d.Get("gw_instance_id").(string))
		if err != nil {
			return errmsgs.WrapError(err)
		}
		nodes, err := jsonpath.Get("$.serviceStruct.nodes", services)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		for _, v := range nodes.([]interface{}) {
			node := v.(map[string]interface{})
			if node["enable"].(bool) {
				service := make(map[string]interface{})
				service["name"] = d.Get("service")
				service["port"] = node["port"]
				service["version"] = "1.0"
				service["weight"] = node["weight"]
				request["services"] = []interface{}{service}
				break
			}

		}

		if d.Get("type").(string) == "DIRECT_ROUTE" {
			directRouteConfig := make(map[string]string)
			directRouteConfig["path"] = d.Get("direct_route_path").(string)
			directRouteConfig["type"] = d.Get("direct_route_type").(string)
			request["directRouteConfig"] = directRouteConfig
			request["serviceProtocol"] = "MCP"
		}
		if d.Get("type").(string) == "DATABASE" {
			dbConfig := make(map[string]interface{})
			dbConfig["host"] = d.Get("db_host")
			dbConfig["port"] = d.Get("db_port")
			dbConfig["username"] = d.Get("db_username")
			dbConfig["password"] = d.Get("db_password")
			dbConfig["dbname"] = d.Get("db_name")
			dbConfig["otherParams"] = d.Get("other_params")
			request["dbConfig"] = dbConfig
			request["dbType"] = d.Get("db_type")
		}
		_, err = client.DoTeaRequest("POST", "csb2", "2023-02-06", "UpdateMcpServer", "/mcpServer/updateMcpServer", nil, nil, request)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_api_gateway_v2_mcpserver", "UpdateMcpServer", errmsgs.AlibabacloudStackSdkGoERROR)
		}

	}
	return nil
}

func resourceAlibabacloudStackAPIGatewayV2McpserverDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	gwInstanceId := d.Get("gw_instance_id").(string)
	name := d.Get("name").(string)

	reqBody := map[string]interface{}{
		"gwInstanceId":  gwInstanceId,
		"mcpServerName": name,
	}

	_, err := client.DoTeaRequest("POST", "csb2", "2023-02-06", "DeleteMcpServer", "/mcpServer/deleteMcpServer", nil, nil, reqBody)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "DeleteMcpServer", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return nil
}
