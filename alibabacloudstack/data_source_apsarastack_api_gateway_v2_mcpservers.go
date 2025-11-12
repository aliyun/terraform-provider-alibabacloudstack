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

func dataSourceAlibabacloudStackAPIGatewayV2McpServers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackAPIGatewayV2McpServersRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A list of MCP server IDs to filter results by.",
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  "A regex string to filter results by the MCP server name.",
			},
			"gw_instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the API Gateway instance.",
			},
			"mcp_servers": {
				Type:     schema.TypeList,
				Computed: true,

				Description: "A list of MCP servers satisfying the criteria.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier of the MCP server. Format: {gwInstanceId}:{name}.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the MCP server.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the MCP server.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the MCP server (e.g., OPEN_API, DATABASE, DIRECT_ROUTE).",
						},
						"service": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The backend service address for the MCP server.",
						},
						"domains": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "The list of domain associated with the MCP server.",
						},
						"services": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The service name in the upstream.",
									},
									"port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The service port in the upstream.",
									},
									"version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The version of the service.",
									},
									"weight": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The weight of the service used for load balancing.",
									},
								},
							},
							Description: "The list of backend services configured for this MCP server.",
						},
						"consumer_auth_info": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enable": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether consumer authentication is enabled.",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type of consumer authentication (e.g., API_KEY).",
									},
									"allowed_consumers": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "List of allowed consumers.",
									},
								},
							},
							Description: "Information about consumer authentication settings.",
						},
						"raw_configurations": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Raw configuration content of the MCP server in YAML format.",
						},
						"db_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of database (e.g., MYSQL), only applicable when type is DATABASE.",
						},
						"dsn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Data Source Name used for connecting to the database.",
						},
						"upstream_path_prefix": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Prefix path for upstream requests.",
						},
						"direct_route_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"path": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The route path for direct routing.",
									},
									"transport_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The transport type for direct routing (e.g., sse).",
									},
								},
							},
							Description: "Configuration for direct route type MCP servers.",
						},
						"gw_instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the gateway instance to which the MCP server belongs.",
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackAPIGatewayV2McpServersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// Parse filters
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
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return fmt.Errorf("invalid name_regex pattern: %v", err)
		}
		nameRegex = r
	}

	gwInstanceId := d.Get("gw_instance_id").(string)

	// Prepare request body
	requestBody := map[string]interface{}{
		"gwInstanceId": gwInstanceId,
		"current":      1,
		"size":         100,
	}

	resp, err := client.DoTeaRequest("POST", "csb2", "2023-02-06", "ListMcpServers", "/mcpServer/listMcpServers", nil, nil, requestBody)
	if err != nil {
		return fmt.Errorf("failed to list MCP servers: %v", err)
	}

	records, err := jsonpath.Get("$.data.records", resp)
	if err != nil {
		return fmt.Errorf("failed to marshal response data: %v", err)
	}
	var mcpServers []map[string]interface{}
	ids := make([]string, 0)

	for _, v := range records.([]interface{}) {
		record := v.(map[string]interface{})
		name := record["name"].(string)
		id := fmt.Sprintf("%s:%s", gwInstanceId, name)
		if len(idsMap) > 0 {
			if _, exists := idsMap[id]; !exists {
				continue
			}
		}

		// Apply name regex filter
		if nameRegex != nil && !nameRegex.MatchString(name) {
			continue
		}

		server := map[string]interface{}{
			"id":                   id,
			"name":                 name,
			"description":          record["description"],
			"type":                 record["type"],
			"service":              record["service"],
			"domains":              record["domains"],
			"gw_instance_id":       gwInstanceId,
			"raw_configurations":   record["rawConfigurations"],
			"db_type":              record["dbType"],
			"dsn":                  record["dsn"],
			"upstream_path_prefix": record["upstreamPathPrefix"],
		}

		// Handle services
		if services, ok := record["services"].([]interface{}); ok {
			serviceList := make([]map[string]interface{}, 0)
			for _, svc := range services {
				if svcMap, ok := svc.(map[string]interface{}); ok {
					service := map[string]interface{}{
						"name":    svcMap["name"],
						"port":    svcMap["port"],
						"version": svcMap["version"],
						"weight":  svcMap["weight"],
					}
					serviceList = append(serviceList, service)
				}
			}
			server["services"] = serviceList
		}

		// Handle consumerAuthInfo
		if consumerAuthInfo, ok := record["consumerAuthInfo"].(map[string]interface{}); ok {
			cai := map[string]interface{}{
				"enable":            consumerAuthInfo["enable"],
				"type":              consumerAuthInfo["type"],
				"allowed_consumers": consumerAuthInfo["allowedConsumers"],
			}
			server["consumer_auth_info"] = []map[string]interface{}{cai}
		}

		// Handle directRouteConfig
		if directRouteConfig, ok := record["directRouteConfig"].(map[string]interface{}); ok {
			drc := map[string]interface{}{
				"path":           directRouteConfig["path"],
				"transport_type": directRouteConfig["transportType"],
			}
			server["direct_route_config"] = []map[string]interface{}{drc}
		}

		mcpServers = append(mcpServers, server)
		ids = append(ids, id)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("mcp_servers", mcpServers); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}
	return nil
}
