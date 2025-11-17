package alibabacloudstack

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackApiGatewayV2Route() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"gw_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "DEFAULT",
			},
			"route_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"path": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"match_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"match_value": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"case_sensitive": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"route_path": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"open_strip_prefix": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"strip_prefix": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"order": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"methods": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"header": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"cookie": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"query_param": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"domain_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"enable_status": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"service_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_ids": {
				Type:          schema.TypeSet,
				Optional:      true,
				ConflictsWith: []string{"service_id"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"weight": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 100),
						},
					},
				},
			},
			"service_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"service_ids"},
			},
			"route_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

	resource.Create = resourceAlibabacloudStackApiGatewayV2RouteCreate
	resource.Read = resourceAlibabacloudStackApiGatewayV2RouteRead
	resource.Update = resourceAlibabacloudStackApiGatewayV2RouteUpdate
	resource.Delete = resourceAlibabacloudStackApiGatewayV2RouteDelete

	return resource
}

func resourceAlibabacloudStackApiGatewayV2RouteCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// Prepare request parameters
	reqBody := make(map[string]interface{})

	if v, ok := d.GetOk("gw_instance_id"); ok {
		reqBody["gwInstanceId"] = v.(string)
	}
	if v, ok := d.GetOk("group_id"); ok {
		reqBody["groupId"] = v.(string)
	}
	if v, ok := d.GetOk("route_name"); ok {
		reqBody["routeName"] = v.(string)
	}

	if v, ok := d.GetOk("path"); ok && len(v.([]interface{})) > 0 {
		pathMap := v.([]interface{})[0].(map[string]interface{})
		path := make(map[string]interface{})
		if val, ok := pathMap["match_type"]; ok {
			path["matchType"] = val.(string)
		}
		if val, ok := pathMap["match_value"]; ok {
			path["matchValue"] = val.(string)
		}
		if val, ok := pathMap["case_sensitive"]; ok {
			path["caseSensitive"] = val.(bool)
		}
		reqBody["path"] = path
	}

	if v, ok := d.GetOk("route_path"); ok {
		routePaths := make([]string, 0)
		for _, item := range v.([]interface{}) {
			routePaths = append(routePaths, item.(string))
		}
		reqBody["routePath"] = routePaths
	}

	if v, ok := d.GetOk("open_strip_prefix"); ok {
		reqBody["openStripPrefix"] = v.(bool)
	}
	if v, ok := d.GetOk("strip_prefix"); ok {
		reqBody["stripPrefix"] = v.(int)
	}
	if v, ok := d.GetOk("order"); ok {
		reqBody["order"] = v.(int)
	}

	if v, ok := d.GetOk("methods"); ok {
		methods := make([]string, 0)
		for _, item := range v.(*schema.Set).List() {
			methods = append(methods, item.(string))
		}
		reqBody["methods"] = methods
	}

	if v, ok := d.GetOk("header"); ok {
		headers := v.(*schema.Set).List()
		for _, item := range v.(*schema.Set).List() {
			header := item.(map[string]interface{})
			h := make(map[string]interface{})
			h["key"] = header["key"]
			h["value"] = header["value"]
			headers = append(headers, h)
		}
		reqBody["header"] = headers
	}

	if v, ok := d.GetOk("cookie"); ok {
		cookies := make([]map[string]interface{}, 0)
		for _, item := range v.(*schema.Set).List() {
			cookie := item.(map[string]interface{})
			c := make(map[string]interface{})
			c["key"] = cookie["key"]
			c["value"] = cookie["value"]
			cookies = append(cookies, c)
		}
		reqBody["cookie"] = cookies
	}

	if v, ok := d.GetOk("query_param"); ok {
		queryParams := make([]map[string]interface{}, 0)
		for _, item := range v.(*schema.Set).List() {
			param := item.(map[string]interface{})
			q := make(map[string]interface{})
			q["key"] = param["key"]
			q["value"] = param["value"]
			queryParams = append(queryParams, q)
		}
		reqBody["queryParam"] = queryParams
	}

	if v, ok := d.GetOk("domain_ids"); ok {
		domainIds := make([]string, 0)
		for _, item := range v.(*schema.Set).List() {
			domainIds = append(domainIds, item.(string))
		}
		reqBody["domainIds"] = domainIds
	}

	if v, ok := d.GetOk("enable_status"); ok {
		reqBody["enableStatus"] = v.(bool)
	}
	if v, ok := d.GetOk("service_type"); ok {
		reqBody["serviceType"] = v.(string)
	}

	if v, ok := d.GetOk("service_ids"); ok {
		serviceIds := make([]map[string]interface{}, 0)
		for _, item := range v.(*schema.Set).List() {
			svc := item.(map[string]interface{})
			s := make(map[string]interface{})
			if val, ok := svc["service_id"]; ok {
				s["serviceId"] = val
			}
			if val, ok := svc["weight"]; ok {
				s["weight"] = val
			}
			serviceIds = append(serviceIds, s)
		}
		reqBody["serviceIds"] = serviceIds
	}

	if v, ok := d.GetOk("service_id"); ok {
		reqBody["serviceId"] = v.(string)
	}

	// Call the API to create the route group
	resp, err := client.DoTeaRequest("POST", "csb2", "2023-02-06", "CreateRoute", "", nil, nil, reqBody)
	if err != nil {
		return err
	}

	// Extract routeId from response data
	routeIdRaw, exists := resp["data"]
	if !exists {
		return fmt.Errorf("failed to get routeId from CreateRoute response")
	}
	routeId := routeIdRaw.(string)

	// Construct resource ID using gwInstanceId and routeId
	gwInstanceId := reqBody["gwInstanceId"].(string)
	resourceId := fmt.Sprintf("%s:%s", gwInstanceId, routeId)
	d.SetId(resourceId)

	// Set computed field route_id
	d.Set("route_id", routeId)

	return nil
}

func resourceAlibabacloudStackApiGatewayV2RouteRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	resourceId := d.Id()
	parts := strings.Split(resourceId, ":")
	if len(parts) != 2 {
		return errmsgs.GetNotFoundErrorFromString("Invalid resource id, should be gwInstanceId:routeId")
	}
	gwInstanceId := parts[0]
	routeId := parts[1]

	query := make(map[string]interface{})
	query["routeId"] = routeId
	query["gwInstanceId"] = gwInstanceId

	resp, err := client.DoTeaRequest("POST", "csb2", "2023-02-06", "GetRoute", "", nil, query, nil)
	if err != nil {
		return err
	}

	if !resp["asapiSuccess"].(bool) || resp["data"] == nil {
		return errmsgs.GetNotFoundErrorFromString("Resource not found")
	}

	data := resp["data"].(map[string]interface{})

	d.Set("gw_instance_id", gwInstanceId)
	d.Set("group_id", data["groupId"])
	d.Set("route_name", data["routeName"])
	d.Set("open_strip_prefix", data["openStripPrefix"])
	d.Set("strip_prefix", data["stripPrefix"])
	d.Set("order", data["order"])
	d.Set("enable_status", data["enableStatus"])
	d.Set("service_type", data["serviceType"])

	if v, ok := data["routePath"].([]interface{}); ok {
		d.Set("route_path", v)
	}

	if v, ok := data["methods"].([]interface{}); ok {
		d.Set("methods", v)
	}

	if headerList, ok := data["header"].([]interface{}); ok {
		var headers []map[string]interface{}
		for _, item := range headerList {
			if m, ok := item.(map[string]interface{}); ok {
				headers = append(headers, map[string]interface{}{
					"value": m["value"],
					"key":   m["key"],
				})
			}
		}
		d.Set("header", headers)
	}

	if cookieList, ok := data["cookie"].([]interface{}); ok {
		var cookies []map[string]interface{}
		for _, item := range cookieList {
			if m, ok := item.(map[string]interface{}); ok {
				cookies = append(cookies, map[string]interface{}{
					"value": m["value"],
					"key":   m["key"],
				})
			}
		}
		d.Set("cookie", cookies)
	}

	if queryParams, ok := data["queryParam"].([]interface{}); ok {
		var params []map[string]interface{}
		for _, item := range queryParams {
			if m, ok := item.(map[string]interface{}); ok {
				params = append(params, map[string]interface{}{
					"value": m["value"],
					"key":   m["key"],
				})
			}
		}
		d.Set("query_param", params)
	}

	if domainVOs, ok := data["domainVO"].([]interface{}); ok {
		var domainIds []string
		for _, item := range domainVOs {
			if m, ok := item.(map[string]interface{}); ok {
				if domainId, ok := m["domainId"].(string); ok {
					domainIds = append(domainIds, domainId)
				}
			}
		}
		d.Set("domain_ids", domainIds)
	}

	if serviceIds, ok := data["serviceIds"].([]interface{}); ok {
		var services []map[string]interface{}
		for _, item := range serviceIds {
			if m, ok := item.(map[string]interface{}); ok {
				services = append(services, map[string]interface{}{
					"service_id": m["serviceId"],
					"weight":     m["weight"],
				})
			}
		}
		d.Set("service_ids", services)
	}

	d.Set("service_id", data["serviceId"])

	return nil
}

func resourceAlibabacloudStackApiGatewayV2RouteUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// Parse resource ID to get routeId and gwInstanceId
	resourceIdParts := strings.Split(d.Id(), ":")
	if len(resourceIdParts) != 2 {
		return errmsgs.WrapErrorf(fmt.Errorf("invalid resource id: %s", d.Id()), errmsgs.DefaultErrorMsg,
			"alibabacloudstack_api_gateway_v2_route_group", "ParseResourceId", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	routeId := resourceIdParts[1]
	gwInstanceId := resourceIdParts[0]

	if d.IsNewResource() {
		// For new resources, no update is needed as all fields are already set in Create
		return nil
	}

	// Prepare request parameters for ModifyRoute
	request := make(map[string]interface{})
	request["routeId"] = routeId
	request["gwInstanceId"] = gwInstanceId

	// Only include fields that can be updated (all fields in this case)
	if v, ok := d.GetOk("group_id"); ok {
		request["groupId"] = v.(string)
	}

	if v, ok := d.GetOk("route_name"); ok {
		request["routeName"] = v.(string)
	}

	if v, ok := d.GetOk("path"); ok && len(v.([]interface{})) > 0 {
		pathMap := v.([]interface{})[0].(map[string]interface{})
		pathObj := make(map[string]interface{})
		if val, ok := pathMap["match_type"]; ok {
			pathObj["matchType"] = val.(string)
		}
		if val, ok := pathMap["match_value"]; ok {
			pathObj["matchValue"] = val.(string)
		}
		if val, ok := pathMap["case_sensitive"]; ok {
			pathObj["caseSensitive"] = val.(bool)
		}
		request["path"] = pathObj
	}

	if v, ok := d.GetOk("route_path"); ok {
		routePaths := make([]string, 0)
		for _, item := range v.([]interface{}) {
			routePaths = append(routePaths, item.(string))
		}
		request["routePath"] = routePaths
	}

	if v, ok := d.GetOk("open_strip_prefix"); ok {
		request["openStripPrefix"] = v.(bool)
	}

	if v, ok := d.GetOk("strip_prefix"); ok {
		request["stripPrefix"] = v.(int)
	}

	if v, ok := d.GetOk("order"); ok {
		request["order"] = v.(int)
	}

	if v, ok := d.GetOk("methods"); ok {
		methods := make([]string, 0)
		for _, item := range v.([]interface{}) {
			methods = append(methods, item.(string))
		}
		request["methods"] = methods
	}

	if v, ok := d.GetOk("header"); ok {
		headers := make([]map[string]interface{}, 0)
		for _, item := range v.([]interface{}) {
			headerItem := item.(map[string]interface{})
			headerObj := make(map[string]interface{})
			if val, exists := headerItem["name"]; exists {
				headerObj["name"] = val.(string)
			}
			if val, exists := headerItem["value"]; exists {
				headerObj["value"] = val.(string)
			}
			if val, exists := headerItem["required"]; exists {
				headerObj["required"] = val.(bool)
			}
			if val, exists := headerItem["description"]; exists {
				headerObj["description"] = val.(string)
			}
			if val, exists := headerItem["item_key"]; exists {
				headerObj["itemKey"] = val.(int)
			}
			if val, exists := headerItem["key"]; exists {
				headerObj["key"] = val.(string)
			}
			headers = append(headers, headerObj)
		}
		request["header"] = headers
	}

	if v, ok := d.GetOk("cookie"); ok {
		cookies := make([]map[string]interface{}, 0)
		for _, item := range v.([]interface{}) {
			cookieItem := item.(map[string]interface{})
			cookieObj := make(map[string]interface{})
			if val, exists := cookieItem["name"]; exists {
				cookieObj["name"] = val.(string)
			}
			if val, exists := cookieItem["value"]; exists {
				cookieObj["value"] = val.(string)
			}
			if val, exists := cookieItem["required"]; exists {
				cookieObj["required"] = val.(bool)
			}
			if val, exists := cookieItem["description"]; exists {
				cookieObj["description"] = val.(string)
			}
			if val, exists := cookieItem["item_key"]; exists {
				cookieObj["itemKey"] = val.(int)
			}
			if val, exists := cookieItem["key"]; exists {
				cookieObj["key"] = val.(string)
			}
			cookies = append(cookies, cookieObj)
		}
		request["cookie"] = cookies
	}

	if v, ok := d.GetOk("query_param"); ok {
		queryParams := make([]map[string]interface{}, 0)
		for _, item := range v.([]interface{}) {
			queryItem := item.(map[string]interface{})
			queryObj := make(map[string]interface{})
			if val, exists := queryItem["name"]; exists {
				queryObj["name"] = val.(string)
			}
			if val, exists := queryItem["value"]; exists {
				queryObj["value"] = val.(string)
			}
			if val, exists := queryItem["required"]; exists {
				queryObj["required"] = val.(bool)
			}
			if val, exists := queryItem["description"]; exists {
				queryObj["description"] = val.(string)
			}
			if val, exists := queryItem["item_key"]; exists {
				queryObj["itemKey"] = val.(int)
			}
			if val, exists := queryItem["key"]; exists {
				queryObj["key"] = val.(string)
			}
			queryParams = append(queryParams, queryObj)
		}
		request["queryParam"] = queryParams
	}

	if v, ok := d.GetOk("domain_ids"); ok {
		domainIds := make([]string, 0)
		for _, item := range v.([]interface{}) {
			domainIds = append(domainIds, item.(string))
		}
		request["domainIds"] = domainIds
	}

	if v, ok := d.GetOk("enable_status"); ok {
		request["enableStatus"] = v.(bool)
	}

	if v, ok := d.GetOk("service_type"); ok {
		request["serviceType"] = v.(string)
	}

	if v, ok := d.GetOk("service_ids"); ok {
		serviceIds := make([]map[string]interface{}, 0)
		for _, item := range v.([]interface{}) {
			serviceItem := item.(map[string]interface{})
			serviceObj := make(map[string]interface{})
			if val, exists := serviceItem["service_id"]; exists {
				serviceObj["serviceId"] = val.(string)
			}
			if val, exists := serviceItem["weight"]; exists {
				serviceObj["weight"] = val.(int)
			}
			serviceIds = append(serviceIds, serviceObj)
		}
		request["serviceIds"] = serviceIds
	}

	if v, ok := d.GetOk("service_id"); ok {
		request["serviceId"] = v.(string)
	}

	// Call ModifyRoute API
	_, err := client.DoTeaRequest("POST", "csb2", "2023-02-06", "ModifyRoute", "", nil, nil, request)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg,
			"alibabacloudstack_api_gateway_v2_route_group", "ModifyRoute", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return nil
}

func resourceAlibabacloudStackApiGatewayV2RouteDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	parts := strings.Split(d.Id(), ":")
	if len(parts) != 2 {
		return errmsgs.WrapError(fmt.Errorf("invalid resource id: %s", d.Id()))
	}
	gwInstanceId := parts[0]
	routeId := parts[1]

	reqQuery := map[string]interface{}{
		"routeId":      routeId,
		"gwInstanceId": gwInstanceId,
	}

	err := resource.Retry(10*time.Minute, func() *resource.RetryError {
		_, err := client.DoTeaRequest("POST", "csb2", "2023-02-06", "DeleteRoute", "", nil, reqQuery, nil)
		if err != nil {
			err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "DeleteRoute", errmsgs.AlibabacloudStackSdkGoERROR, "")
			return resource.RetryableError(err)
		}
		return nil
	})

	if err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
