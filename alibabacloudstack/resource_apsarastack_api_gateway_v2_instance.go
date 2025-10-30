package alibabacloudstack

import (
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackAPIGateWayV2Instance() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"deploy_cluster_namespace": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"broker_engine_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"k8s_service_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"access_mode": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tid": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"k8s_cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"broker_engine_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"HIGRESS", "SCG"}, false),
			},
			"deploy_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"deploy_cluster_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"node_number": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"instance_class": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"edas_app_infos": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"edas_namespace": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"app_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"k8s_cluster_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"k8s_namespace": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "default",
						},
					},
				},
			},
			"broker_latest_engine_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"shared_instance": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"deploy_cluster_code": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"edas_namespace_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"prometheus_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"sls_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"edas_app_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackAPIGateWayV2InstanceCreate, resourceAlibabacloudStackAPIGateWayV2InstanceRead, resourceAlibabacloudStackAPIGateWayV2InstanceUpdate, resourceAlibabacloudStackAPIGateWayV2InstanceDelete)
	return resource
}

func resourceAlibabacloudStackAPIGateWayV2InstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	broker_engine_type := d.Get("broker_engine_type").(string)
	request := map[string]interface{}{
		"gwInstanceName":   d.Get("instance_name"),
		"instanceNumber":   d.Get("node_number"),
		"instanceClass":    d.Get("instance_class"),
		"brokerEngineType": broker_engine_type,
		"deployMode":       d.Get("deploy_mode"),
	}
	if broker_engine_type == "SCG" {
		edasAppInfos := make([]map[string]interface{}, 0)
		for _, v := range d.Get("edas_app_infos").(*schema.Set).List() {
			edasAppInfo := v.(map[string]interface{})
			edasAppInfos = append(edasAppInfos, map[string]interface{}{
				"edasNamespaceId": edasAppInfo["edas_namespace_id"],
				"k8sClusterId":    edasAppInfo["edas_k8s_id"],
				"k8sNamespace":    edasAppInfo["k8s_namespace"],
			})
		}
		if len(edasAppInfos) > 0 {
			request["edasAppInfos"] = edasAppInfos
		}
	} else {
		request["clusterCode"] = d.Get("deploy_cluster_code")
		request["ingressClassName"] = d.Get("deploy_cluster_namespace")
		request["namespace"] = d.Get("k8s_service_name")
		request["slsEnabled"] = d.Get("sls_enabled")
		request["prometheusEnabled"] = d.Get("prometheus_enabled")
		request["o11y"] = map[string]interface{}{
			"slsEnabled":        d.Get("sls_enabled"),
			"prometheusEnabled": d.Get("prometheus_enabled"),
		}
	}
	response, err := client.DoTeaRequest("POST", "csb2", "2023-02-06", "AscmCreateInstance", "/gatewayInstance/ascmCreateInstance", nil, nil, request)
	if err != nil {
		return err
	}
	id, ok := response["data"]
	if !ok {
		return errmsgs.Error("CreateInstance Failed! %v", response)
	}
	d.SetId(id.(string))
	return nil
}

func resourceAlibabacloudStackAPIGateWayV2InstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	apigatewayv2Service := ApiGateWayV2Service{client}
	instance, err := apigatewayv2Service.DescribeApiGatewayV2Instace(d.Id())
	if err != nil {
		return errmsgs.WrapError(err)
	}
	d.Set("instance_name", instance["gwInstanceName"])
	d.Set("deploy_cluster_namespace", instance["deployClusterNamespace"])
	d.Set("broker_engine_version", instance["brokerEngineVersion"])
	d.Set("k8s_service_name", instance["k8sServiceName"])
	d.Set("access_mode", instance["accessMode"])
	d.Set("tid", instance["tid"])
	d.Set("k8s_cluster_id", instance["k8sClusterId"])
	d.Set("broker_engine_type", instance["brokerEngineType"])
	d.Set("deploy_mode", instance["deployMode"])
	d.Set("deploy_cluster_name", instance["deployClusterName"])
	d.Set("node_number", instance["nodeNumber"])
	d.Set("instance_class", instance["instanceClass"])
	d.Set("broker_latest_engine_version", instance["brokerLatestEngineVersion"])
	d.Set("shared_instance", instance["sharedInstance"])
	d.Set("deploy_cluster_code", instance["deployClusterCode"])
	d.Set("create_time", instance["createTime"])
	d.Set("edas_namespace_id", instance["edasNamespaceId"])
	d.Set("edas_app_id", instance["edasAppId"])
	d.Set("gw_instance_id", instance["gwInstanceId"])
	d.Set("status", instance["status"])

	// Handle edas_app_infos
	if edasAppInfos, ok := instance["edasAppInfos"].([]interface{}); ok {
		edasAppInfosSet := make([]map[string]interface{}, 0)
		for _, item := range edasAppInfos {
			if appInfo, ok := item.(map[string]interface{}); ok {
				edasAppInfo := make(map[string]interface{})
				if v, ok := appInfo["edasNamespaceId"]; ok {
					edasAppInfo["edas_namespace"] = v
				}
				if v, ok := appInfo["appId"]; ok {
					edasAppInfo["app_id"] = v
				}
				if v, ok := appInfo["k8sClusterId"]; ok {
					edasAppInfo["k8s_cluster_id"] = v
				}
				if v, ok := appInfo["k8sNamespace"]; ok {
					edasAppInfo["k8s_namespace"] = v
				}
				edasAppInfosSet = append(edasAppInfosSet, edasAppInfo)
			}
		}
		d.Set("edas_app_infos", edasAppInfosSet)
	}

	if o11yInfo, ok := instance["o11y"].(map[string]interface{}); ok {
		if v, ok := o11yInfo["prometheusEnabled"]; ok {
			d.Set("prometheus_enabled", v)
		}
		if v, ok := o11yInfo["slsEnabled"]; ok {
			d.Set("sls_enabled", v)
		}
	}
	return nil
}

func resourceAlibabacloudStackAPIGateWayV2InstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	// client := meta.(*connectivity.AlibabacloudStackClient)
	return nil
}

func resourceAlibabacloudStackAPIGateWayV2InstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := map[string]interface{}{
		"gwInstanceId": d.Id(),
		"forcedDelete": true,
	}
	_, err := client.DoTeaRequest("POST", "csb2", "2023-02-06", "DeleteInstance", "/gatewayInstance/deleteInstance", nil, request, nil)
	if err != nil {
		return err
	}
	return nil
}
