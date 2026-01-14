package alibabacloudstack

import (
	"regexp"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackAPIGatewayV2Instances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackAPIGatewayV2InstancesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"broker_engine_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"HIGRESS", "SCG"}, false),
			},
			"deploy_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"k8s", "edas", "custom"}, false),
			},
			"name_regex": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validation.StringIsValidRegExp,
				Deprecated:    "Field 'name_regex' is deprecated and will be removed in a future release. Please use new field 'description_regex' instead.",
				ConflictsWith: []string{"description_regex"},
			},
			"description_regex": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validation.StringIsValidRegExp,
				ConflictsWith: []string{"name_regex"},
			},

			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"broker_engine_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"deploy_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_class": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"k8s_cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"k8s_namespace": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"edas_app_infos": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"edas_namespace": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"app_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"k8s_cluster_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"k8s_namespace": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"edas_app_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"edas_namespace_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"shared_instance": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"node_number": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"custom_deploy_config": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackAPIGatewayV2InstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	apigatewayv2Service := ApiGateWayV2Service{client}
	request := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		request["gwInstanceId"] = v.(string)
	}
	if v, ok := d.GetOk("deploy_mode"); ok {
		request["deployMode"] = v.(string)
	}
	if v, ok := d.GetOk("broker_engine_type"); ok {
		request["brokerEngineType"] = v.(string)
	}
	request["current"]=1
	request["size"]=1000
	response, err := client.DoTeaRequest("POST", "csb2", "2023-02-06", "ListInstances", "/gatewayInstance/listInstances", nil, request, request)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	v, err := jsonpath.Get("$.data.records", response)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	idsMap := getIdsStringFilter(d)

	var ids []string
	datas := make([]interface{}, 0)
	for _, item := range v.([]interface{}) {
		data := item.(map[string]interface{})
		if description_regex, ok := connectivity.GetResourceDataOk(d, "description_regex", "name_regex"); ok {
			r := regexp.MustCompile(description_regex.(string))
			if !r.MatchString(data["instanceName"].(string)) {
				continue
			}
		}

		if len(idsMap) > 0 {
			if _, exist := idsMap[data["gwInstanceId"].(string)]; !exist {
				continue
			}
		}
		edasAppInfosSet := make([]map[string]interface{}, 0)
		if edasAppInfos, ok := data["edasAppInfos"].([]interface{}); ok {
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
		}
		customDeployConfig := make(map[string]interface{})
		if data["deployMode"].(string) == "custom" {
			customDeployConfig, err = apigatewayv2Service.GetCustomDeployConfig(data["gwInstanceId"].(string))
			if err != nil {
				return errmsgs.WrapError(err)
			}
		}
		datas = append(datas, map[string]interface{}{
			"id":                   data["gwInstanceId"],
			"instance_id":          data["gwInstanceId"],
			"instance_name":        data["instanceName"],
			"broker_engine_type":   data["brokerEngineType"],
			"deploy_mode":          data["deployMode"],
			"instance_class":       data["instanceClass"],
			"status":               data["status"],
			"k8s_cluster_id":       data["k8sClusterId"],
			"k8s_namespace":        data["k8sNamespace"],
			"edas_app_id":          data["edasAppId"],
			"edas_namespace_id":    data["edasNamespaceId"],
			"shared_instance":      data["sharedInstance"],
			"node_number":          data["nodeNumber"],
			"edas_app_infos":       edasAppInfosSet,
			"custom_deploy_config": customDeployConfig,
		})
		ids = append(ids, data["gwInstanceId"].(string))
	}
	d.Set("ids", ids)
	d.Set("instances", datas)
	d.SetId(dataResourceIdHash(ids))

	return nil
}
