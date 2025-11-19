package alibabacloudstack

import (
	"fmt"
	"log"
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
				ForceNew: true,
			},
			"route_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
			"cascade_link_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"enable_status": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
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

	setResourceFunc(resource, resourceAlibabacloudStackApiGatewayV2RouteCreate, resourceAlibabacloudStackApiGatewayV2RouteRead, resourceAlibabacloudStackApiGatewayV2RouteUpdate, resourceAlibabacloudStackApiGatewayV2RouteDelete)

	return resource
}

func resourceAlibabacloudStackApiGatewayV2RouteCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// Prepare request parameters
	reqBody := make(map[string]interface{})

	reqBody["gwInstanceId"] = d.Get("gw_instance_id")
	reqBody["routeName"] = d.Get("route_name")
	reqBody["groupId"] = d.Get("group_id")
	cascade_link_ids := d.Get("cascade_link_ids").(*schema.Set).List()
	action := "CreateRoute"
	pattern := "/route/createRoute"
	idpre := "route"
	if len(cascade_link_ids) > 0 {
		cascadeLinkIds := make([]string, 0)
		for _, item := range cascade_link_ids {
			cascadeLinkIds = append(cascadeLinkIds, item.(string))
		}
		reqBody["cascadeLinkIds"] = cascadeLinkIds
		action = "CreateSourceRoute"
		pattern = "/sourceRoute/createSourceRoute"
		idpre = "sourceRoute"
	}

	if v, ok := d.GetOk("route_path"); ok {
		routePaths := make([]string, 0)
		for _, item := range v.([]interface{}) {
			routePaths = append(routePaths, item.(string))
		}
		reqBody["routePath"] = routePaths
	}

	strip_prefix := d.Get("strip_prefix").(int)
	if strip_prefix > 0 {
		reqBody["stripPrefix"] = strip_prefix
		reqBody["openStripPrefix"] = true
	} else {
		reqBody["openStripPrefix"] = false
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
		headers := make([]map[string]interface{}, 0)
		for _, item := range v.(*schema.Set).List() {
			header := item.(map[string]interface{})
			h := make(map[string]interface{})
			h["key"] = header["key"]
			h["value"] = header["value"]
			headers = append(headers, h)
		}
		if len(cascade_link_ids) > 0 {
			headers = append(headers, map[string]interface{}{
				"key":   "csb_cascade",
				"value": "true",
			})
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

	if v, ok := d.GetOk("enable_status"); ok {
		reqBody["enableStatus"] = v.(bool)
	}

	if v, ok := d.GetOk("service_ids"); ok {
		serviceIds := make([]map[string]interface{}, 0)
		for _, item := range v.(*schema.Set).List() {
			svc := item.(map[string]interface{})
			serviceIds = append(serviceIds, map[string]interface{}{
				"serviceId": svc["service_id"],
				"weight":    svc["weight"],
			})
		}
		reqBody["serviceIds"] = serviceIds
		reqBody["serviceType"] = "MULTI"
	}

	if v, ok := d.GetOk("service_id"); ok {
		reqBody["serviceId"] = v.(string)
		reqBody["serviceType"] = "SINGLE"
	}
	if v, ok := d.GetOk("domain_ids"); ok {
		domainIds := make([]string, 0)
		for _, item := range v.(*schema.Set).List() {
			domainIds = append(domainIds, item.(string))
		}
		reqBody["domainIds"] = domainIds
	}

	// Call the API to create the route group
	resp, err := client.DoTeaRequest("POST", "csb2", "2023-02-06", action, pattern, nil, nil, reqBody)
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
	resourceId := fmt.Sprintf("%s:%s:%s", idpre, gwInstanceId, routeId)
	d.SetId(resourceId)

	return nil
}

func resourceAlibabacloudStackApiGatewayV2RouteRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	apiGatewayV2Service := &ApiGateWayV2Service{client}

	object, err := apiGatewayV2Service.DescribeApigwV2Route(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_api_gateway_v2_route apiGatewayV2Service.DescribeApigwV2Route Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return err
	}
	gwInstanceId := parts[1]

	d.Set("gw_instance_id", gwInstanceId)
	if v, ok := object["groupId"]; ok && v != nil {
		d.Set("group_id", v)
	}
	d.Set("route_name", object["routeName"])
	d.Set("strip_prefix", object["stripPrefix"])
	d.Set("order", object["order"])
	d.Set("enable_status", object["enableStatus"])

	if v, ok := object["routePath"].([]interface{}); ok {
		d.Set("route_path", v)
	}

	if v, ok := object["methods"].([]interface{}); ok {
		d.Set("methods", v)
	}

	if headerList, ok := object["header"].([]interface{}); ok {
		var headers []map[string]interface{}
		for _, item := range headerList {
			if m, ok := item.(map[string]interface{}); ok {
				if m["key"] == "csb_cascade" {
					continue
				}
				headers = append(headers, map[string]interface{}{
					"value": m["value"],
					"key":   m["key"],
				})
			}
		}
		d.Set("header", headers)
	}

	if cookieList, ok := object["cookie"].([]interface{}); ok {
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

	if queryParams, ok := object["queryParam"].([]interface{}); ok {
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

	if domainVOs, ok := object["domainVO"].([]interface{}); ok {
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
	serviceType, ok := object["serviceType"]
	if ok && serviceType.(string) == "MULTI" {
		if serviceIds, ok := object["serviceIds"].([]interface{}); ok {
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
			d.Set("service_id", nil)
		}
	} else {
		d.Set("service_ids", nil)
		d.Set("service_id", object["serviceId"])

	}
	linkRouteRelations, ok := object["linkRouteRelations"]
	if ok && linkRouteRelations != nil {
		cls := make([]interface{}, 0)
		for _, item := range linkRouteRelations.([]interface{}) {
			data := item.(map[string]interface{})
			if cascadeInstanceId, ok := data["cascadeInstanceId"]; ok {
				cls = append(cls, cascadeInstanceId)
			}
		}
		d.Set("cascade_link_ids", cls)
	}
	d.Set("route_id", object["routeId"])

	return nil
}

func resourceAlibabacloudStackApiGatewayV2RouteUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return err
	}
	gwInstanceId := parts[1]
	routeId := parts[2]

	if d.IsNewResource() {
		// For new resources, no update is needed as all fields are already set in Create
		return nil
	}

	// Prepare request parameters for ModifyRoute
	reqBody := make(map[string]interface{})
	reqBody["routeId"] = routeId
	reqBody["gwInstanceId"] = gwInstanceId
	reqBody["routeName"] = d.Get("route_name")
	reqBody["groupId"] = d.Get("group_id")
	cascade_link_ids := d.Get("cascade_link_ids").(*schema.Set).List()
	action := "ModifyRoute"
	pattern := "/route/modifyRoute"
	if len(cascade_link_ids) > 0 {
		cascadeLinkIds := make([]string, 0)
		for _, item := range cascade_link_ids {
			cascadeLinkIds = append(cascadeLinkIds, item.(string))
		}
		reqBody["cascadeLinkIds"] = cascadeLinkIds
		action = "ModifySourceRoute"
		pattern = "/sourceRoute/modifySourceRoute"
	}

	// if v, ok := d.GetOk("path"); ok && len(v.([]interface{})) > 0 {
	// 	pathMap := v.([]interface{})[0].(map[string]interface{})
	// 	pathObj := make(map[string]interface{})
	// 	if val, ok := pathMap["match_type"]; ok {
	// 		pathObj["matchType"] = val.(string)
	// 	}
	// 	if val, ok := pathMap["match_value"]; ok {
	// 		pathObj["matchValue"] = val.(string)
	// 	}
	// 	if val, ok := pathMap["case_sensitive"]; ok {
	// 		pathObj["caseSensitive"] = val.(bool)
	// 	}
	// 	request["path"] = pathObj
	// }

	if v, ok := d.GetOk("route_path"); ok {
		routePaths := make([]string, 0)
		for _, item := range v.([]interface{}) {
			routePaths = append(routePaths, item.(string))
		}
		reqBody["routePath"] = routePaths
	}

	strip_prefix := d.Get("strip_prefix").(int)
	if strip_prefix > 0 {
		reqBody["stripPrefix"] = strip_prefix
		reqBody["openStripPrefix"] = true
	} else {
		reqBody["openStripPrefix"] = false
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
		headers := make([]map[string]interface{}, 0)
		for _, item := range v.(*schema.Set).List() {
			header := item.(map[string]interface{})
			h := make(map[string]interface{})
			h["key"] = header["key"]
			h["value"] = header["value"]
			headers = append(headers, h)
		}
		if len(cascade_link_ids) > 0 {
			headers = append(headers, map[string]interface{}{
				"key":   "csb_cascade",
				"value": "true",
			})
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
	if v, ok := d.GetOk("enable_status"); ok {
		reqBody["enableStatus"] = v.(bool)
	}

	if v, ok := d.GetOk("service_ids"); ok {
		serviceIds := make([]map[string]interface{}, 0)
		for _, item := range v.(*schema.Set).List() {
			svc := item.(map[string]interface{})
			serviceIds = append(serviceIds, map[string]interface{}{
				"serviceId": svc["service_id"],
				"weight":    svc["weight"],
			})
		}
		reqBody["serviceIds"] = serviceIds
		reqBody["serviceType"] = "MULTI"
	}

	if v, ok := d.GetOk("service_id"); ok {
		reqBody["serviceId"] = v.(string)
		reqBody["serviceType"] = "SINGLE"
	}
	if v, ok := d.GetOk("domain_ids"); ok {
		domainIds := make([]string, 0)
		for _, item := range v.(*schema.Set).List() {
			domainIds = append(domainIds, item.(string))
		}
		reqBody["domainIds"] = domainIds
	}

	_, err = client.DoTeaRequest("POST", "csb2", "2023-02-06", action, pattern, nil, nil, reqBody)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg,
			"alibabacloudstack_api_gateway_v2_route_group", action, errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return nil
}

func resourceAlibabacloudStackApiGatewayV2RouteDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	parts := strings.Split(d.Id(), ":")
	if len(parts) != 3 {
		return errmsgs.WrapError(fmt.Errorf("invalid resource id: %s", d.Id()))
	}
	idpre := parts[0]
	gwInstanceId := parts[1]
	routeId := parts[2]

	reqQuery := map[string]interface{}{
		"routeId":      routeId,
		"gwInstanceId": gwInstanceId,
	}
	action := "DeleteRoute"
	pattern := "/route/deleteRoute"
	if idpre == "sourceRoute" {
		action = "DeleteSourceRoute"
		pattern = "/sourceRoute/deleteSourceRoute"
	}

	err := resource.Retry(10*time.Minute, func() *resource.RetryError {
		_, err := client.DoTeaRequest("POST", "csb2", "2023-02-06", action, pattern, nil, nil, reqQuery)
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
