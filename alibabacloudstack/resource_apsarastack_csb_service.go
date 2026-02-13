package alibabacloudstack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackCsbService() *schema.Resource {
	resource := &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Update: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"csb_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The unique ID of the CSB service.",
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of the project to which the service belongs.",
			},
			"service_version": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The version of the service (e.g., '1.0.0').",
			},
			"alias": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The alias name of the service.",
			},
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the service.",
			},
			"model_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "2.0",
				Description: "The model version of the service.",
			},
			"skip_auth": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to skip authentication for this service.",
			},
			"all_visiable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the service is visible to all users.",
			},
			"scope": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "0",
				Description: "The scope of the service (e.g., '0' for private, '1' for public).",
			},
			"consume_types": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"Restful", "WebService"}, true),
				},
				Description:           "List of consume types, e.g., ['Restful'].",
				MinItems:              1,
				DiffSuppressFunc:      ignoreCaseDiffSuppressFunc,
				DiffSuppressOnRefresh: true,
			},
			"provide_type": {
				Type:                  schema.TypeString,
				Optional:              true,
				Default:               "Restful",
				ValidateFunc:          validation.StringInSlice([]string{"RESTful", "SpringCloud", "HSF", "WebService", "DUBBO", "JDBC"}, true),
				Description:           "The provide type of the service.",
				DiffSuppressFunc:      ignoreCaseDiffSuppressFunc,
				DiffSuppressOnRefresh: true,
			},
			"cas_serv_targets": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of CAS service targets.",
			},
			"qps": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     "0",
				Description: "Queries per second (read-only monitoring metric).",
			},
			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Service status (e.g., 0=inactive, 1=active).",
			},
			"modified_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Last modified timestamp in milliseconds.",
			},
			"principal_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The principal name associated with the service.",
			},
			"owner_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The owner ID of the service.",
			},
			"user_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user ID of the service creator.",
			},
			"interface_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The interface name (usually empty for RESTful services).",
			},
			"ssl": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether SSL is enabled for the service.",
			},
			"ott_flag": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "OTF flag status.",
			},
			"policy_handler": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The policy handler (e.g., 'accept').",
			},
			"ip_white_str": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_black_str": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"err_def_json": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsJSON,
			},
			"access_params_json": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsJSON,
			},
			"route_conf_json": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsJSON,
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
					if oldValue == newValue {
						return true
					}

					var oldObj, newObj map[string]interface{}
					if err := json.Unmarshal([]byte(oldValue), &oldObj); err != nil {
						return false
					}
					if err := json.Unmarshal([]byte(newValue), &newObj); err != nil {
						return false
					}

					// Helper to safely get accessEndpointJSON string from importConf
					getAccessJSON := func(obj map[string]interface{}) (string, bool) {
						if imp, ok := obj["importConf"].(map[string]interface{}); ok {
							if val, exists := imp["accessEndpointJSON"]; exists {
								if s, ok := val.(string); ok {
									return s, true
								}
							}
						}
						return "", false
					}

					oldAccessStr, hasOld := getAccessJSON(oldObj)
					newAccessStr, hasNew := getAccessJSON(newObj)

					// If both have accessEndpointJSON, compare as JSON
					if hasOld && hasNew {
						var oldAccess, newAccess interface{}
						if err := json.Unmarshal([]byte(oldAccessStr), &oldAccess); err != nil {
							return false
						}
						if err := json.Unmarshal([]byte(newAccessStr), &newAccess); err != nil {
							return false
						}
						b1, _ := json.Marshal(oldAccess)
						b2, _ := json.Marshal(newAccess)
						if !bytes.Equal(b1, b2) {
							return false
						}
					} else if hasOld != hasNew {
						// One has it, the other doesn't
						return false
					}

					// Remove accessEndpointJSON from both for outer comparison
					if imp, ok := oldObj["importConf"].(map[string]interface{}); ok {
						delete(imp, "accessEndpointJSON")
					}
					if imp, ok := newObj["importConf"].(map[string]interface{}); ok {
						delete(imp, "accessEndpointJSON")
					}

					// Compare the rest of the structure
					b1, err1 := json.Marshal(oldObj)
					b2, err2 := json.Marshal(newObj)
					if err1 != nil || err2 != nil {
						return false
					}
					return bytes.Equal(b1, b2)
				},
				DiffSuppressOnRefresh: true,
			},
			"service_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackCsbServiceCreate,
		resourceAlibabacloudStackCsbServiceRead, resourceAlibabacloudStackCsbServiceUpdate, resourceAlibabacloudStackCsbServiceDelete)
	return resource
}

func resourceAlibabacloudStackCsbServiceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	action := "CreateService"
	reqQuery := map[string]interface{}{
		"CsbId": d.Get("csb_id"),
	}

	data, err := convertDataString(d)
	if err != nil {
		return err
	}

	body := map[string]interface{}{
		"Data": data,
	}

	response, err := client.DoTeaRequest("POST", "CSB", "2017-11-18", action, "", nil, reqQuery, body)
	if err != nil {
		return err
	}

	if v, existed := response["Data"]; !existed {
		return errmsgs.Error("error CreateProject response %v", response)
	} else {
		data := v.(map[string]interface{})
		if v, existed := data["Id"]; !existed {
			return errmsgs.Error("error Id not in CreateProject response %v", data)
		} else {
			d.Set("service_id", v)
		}
	}

	d.SetId(fmt.Sprint(d.Get("csb_id"), ":", d.Get("service_id")))
	return nil
}

func resourceAlibabacloudStackCsbServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	csbService := CsbService{client}
	object, err := csbService.DescribeCsbServiceDetail(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_csb_project csbService.DescribeCsbProjectDetail Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("csb_id", object["CsbId"])
	d.Set("project_id", object["ProjectId"])
	d.Set("service_version", object["ServiceVersion"])
	d.Set("alias", object["Alias"])
	d.Set("service_name", object["ServiceName"])
	d.Set("service_id", object["Id"])

	d.Set("model_version", object["ModelVersion"])
	d.Set("interface_name", object["InterfaceName"])
	d.Set("principal_name", object["PrincipalName"])
	d.Set("owner_id", object["owner_id"])
	d.Set("user_id", object["UserId"])
	d.Set("policy_handler", object["PolicyHandler"])
	if v, existed := object["IpWhiteStr"]; existed && v.(string) != "null" {
		d.Set("ip_white_str", object["IpWhiteStr"])
	} else {
		d.Set("ip_white_str", nil)
	}
	if v, existed := object["IpBlackStr"]; existed && v.(string) != "null" {
		d.Set("ip_black_str", object["IpBlackStr"])
	} else {
		d.Set("ip_black_str", nil)
	}
	if v, existed := object["ErrDefJson"]; existed && v.(string) != "null" {
		d.Set("err_def_json", object["ErrDefJson"])
	} else {
		d.Set("err_def_json", nil)
	}
	d.Set("access_params_json", object["AccessParamsJson"])

	d.Set("skip_auth", object["SkipAuth"])
	d.Set("all_visiable", object["AllVisiable"])
	d.Set("ssl", object["SSL"])
	d.Set("ott_flag", object["OttFlag"])
	d.Set("consume_types", object["ConsumeTypes"])
	d.Set("cas_serv_targets", object["CasServTargets"])

	d.Set("scope", object["Scope"])
	d.Set("provide_type", object["ProvideType"])

	d.Set("qps", object["Qps"])
	d.Set("status", object["Status"])
	d.Set("modified_time", object["ModifiedTime"])

	if v, ok := object["RouteConfJson"]; ok {
		if routeStr, ok := v.(string); ok {
			unescaped := html.UnescapeString(routeStr)
			d.Set("route_conf_json", unescaped)
		}
	}

	return nil
}

func resourceAlibabacloudStackCsbServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.IsNewResource() {
		return nil
	}
	client := meta.(*connectivity.AlibabacloudStackClient)
	if d.HasChanges("service_version", "alias", "project_id", "config") {
		reqQuery := map[string]interface{}{
			"CsbId": d.Get("csb_id"),
		}

		data, err := convertDataString(d)
		if err != nil {
			return err
		}

		body := map[string]interface{}{
			"Data": data,
		}

		action := "UpdateService"
		_, err = client.DoTeaRequest("POST", "CSB", "2017-11-18", action, "", nil, reqQuery, body)
		if err != nil {
			return err
		}
	}
	return nil
}

func resourceAlibabacloudStackCsbServiceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	action := "DeleteService"
	request := map[string]interface{}{
		"ServiceName": d.Get("service_name"),
		"CsbId":       d.Get("csb_id"),
	}

	_, err := client.DoTeaRequest("POST", "CSB", "2017-11-18", action, "", nil, nil, request)
	if err != nil {
		return err
	}
	return nil
}

func convertDataString(d *schema.ResourceData) (string, error) {
	data := make(map[string]interface{})

	data["serviceVersion"] = d.Get("service_version").(string)

	if v, ok := d.GetOk("service_id"); ok && v.(string) != "" {
		if id, err := toInt(v); err == nil {
			data["id"] = id
		}
	}

	data["serviceName"] = d.Get("service_name").(string)
	data["projectId"] = d.Get("project_id").(string)

	if v, ok := d.GetOk("alias"); ok {
		data["alias"] = v.(string)
	}
	if v, ok := d.GetOk("model_version"); ok {
		data["modelVersion"] = v.(string)
	}
	if v, ok := d.GetOk("interface_name"); ok {
		data["interfaceName"] = v.(string)
	}
	if v, ok := d.GetOk("scope"); ok {
		data["scope"] = v.(string)
	}
	if v, ok := d.GetOk("provide_type"); ok {
		data["provideType"] = v.(string)
	}
	if v, ok := d.GetOk("ip_white_str"); ok {
		data["ipWhiteStr"] = v.(string)
	}
	if v, ok := d.GetOk("ip_black_str"); ok {
		data["ipBlackStr"] = v.(string)
	}
	if v, ok := d.GetOk("err_def_json"); ok {
		data["errDefJson"] = v.(string)
	}
	if v, ok := d.GetOk("access_params_json"); ok {
		data["accessParamsJson"] = v.(string)
	}

	data["skipAuth"] = d.Get("skip_auth").(bool)
	data["allVisiable"] = d.Get("all_visiable").(bool)
	data["ssl"] = d.Get("ssl").(bool)
	data["ottFlag"] = d.Get("ott_flag").(bool)

	if v, ok := d.GetOk("consume_types"); ok {
		list := v.([]interface{})
		strList := make([]string, len(list))
		for i, item := range list {
			strList[i] = item.(string)
		}
		data["consumeTypes"] = strList
	}

	if v, ok := d.GetOk("cas_serv_targets"); ok {
		list := v.([]interface{})
		strList := make([]string, len(list))
		for i, item := range list {
			strList[i] = item.(string)
		}
		data["casServTargets"] = strList
	}

	if v, ok := d.GetOk("route_conf_json"); ok {
		routeConfStr := v.(string)
		var tmp interface{}
		if err := json.Unmarshal([]byte(routeConfStr), &tmp); err != nil {
			return "", fmt.Errorf("invalid JSON in route_conf_json: %w", err)
		}
		data["routeConfJson"] = routeConfStr
	} else {
		return "", fmt.Errorf("route_conf_json is required")
	}

	content, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal service data: %w", err)
	}

	return string(content), nil
}
