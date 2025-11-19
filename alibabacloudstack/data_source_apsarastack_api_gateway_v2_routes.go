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

func dataSourceAlibabacloudStackAPIGatewayV2Routes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackAPIGatewayV2RoutesRead,

		Schema: map[string]*schema.Schema{
			"gw_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is_source_route": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"routes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"route_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"route_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"route_path": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"service_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enable_status": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// ... existing code ...

func dataSourceAlibabacloudStackAPIGatewayV2RoutesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// Get gwInstanceId from schema
	gwInstanceId := d.Get("gw_instance_id").(string)
	is_source_route := d.Get("is_source_route").(bool)
	// Prepare request body for ListRoutes API
	reqBody := map[string]interface{}{
		"gwInstanceId": gwInstanceId,
	}

	action := "ListRoutes"
	pattern := "/route/listRoutes"
	if is_source_route {
		action = "ListSourceRoutes"
		pattern = "/sourceRoute/listSourceRoutes"
	}

	resp, err := client.DoTeaRequest("POST", "csb2", "2023-02-06", action, pattern, nil, nil, reqBody)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	// Parse the response data into a list of routes
	routesData, err := jsonpath.Get("$.data.records", resp)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	// Build idsMap if ids are provided
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	// Compile nameRegex if provided
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		nameRegex = regexp.MustCompile(v.(string))
	}

	// Process each route and build result
	var filteredRoutes []interface{}
	var routeIds []string

	for _, routeItem := range routesData.([]interface{}) {
		route := routeItem.(map[string]interface{})

		routeId := route["routeId"].(string)
		// Construct full ID as gwInstanceId:routeId
		resourceId := fmt.Sprintf("%s:%s", gwInstanceId, routeId)

		// Filter by ids if specified
		if len(idsMap) > 0 {
			if _, exists := idsMap[resourceId]; !exists {
				continue
			}
		}

		// Filter by name_regex if specified
		if nameRegex != nil {
			routeName, _ := route["routeName"].(string)
			if !nameRegex.MatchString(routeName) {
				continue
			}
		}

		// Add route to filtered list
		filteredRoutes = append(filteredRoutes, map[string]interface{}{
			"id":            resourceId,
			"route_id":      routeId,
			"route_name":    route["routeName"],
			"group_id":      route["groupId"],
			"group_name":    route["groupName"],
			"service_name":  route["serviceName"],
			"route_path":    route["routePath"],
			"service_id":    route["serviceId"],
			"enable_status": route["enableStatus"],
			"create_time":   route["createTime"],
		})

		routeIds = append(routeIds, resourceId)
	}

	// Set the ID for the data source
	d.SetId(dataResourceIdHash(routeIds))
	// Set the computed routes attribute
	if err := d.Set("routes", filteredRoutes); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("ids", routeIds); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
