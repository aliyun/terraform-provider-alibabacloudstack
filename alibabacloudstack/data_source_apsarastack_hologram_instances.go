package alibabacloudstack

import (
	"regexp"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackHologramInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackHologramInstancesRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
						"compute_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"node": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enable_hive_access": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu_arch": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu_brand": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu_brand_i18n": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"apsara_ascm_cpu_brand": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"support_replica": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"ascm_create_user": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"commodity_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"endpoints": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"endpoint": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"vpc_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"vswitch_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"vpc_instance_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackHologramInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	response, err := client.DoTeaRequest("GET", "Hologram", "2022-06-01", "ListInstances", "/api/v1/instances", nil, nil, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_hologram_instances", "ListInstances", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	instances, err := jsonpath.Get("$.InstanceList", response)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_hologram_instances", "ListInstances", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	idsMap := getIdsStringFilter(d)

	var ids []string
	var s []map[string]interface{}

	for _, item := range instances.([]interface{}) {
		instance := item.(map[string]interface{})
		if name_regex, ok := d.GetOk("name_regex"); ok && name_regex.(string) != "" {
			r := regexp.MustCompile(name_regex.(string))
			if !r.MatchString(instance["InstanceName"].(string)) {
				continue
			}
		}

		if len(idsMap) > 0 {
			if _, exist := idsMap[instance["InstanceId"].(string)]; !exist {
				continue
			}
		}

		endpointList := make([]map[string]interface{}, 0)
		if endpoints, ok := instance["Endpoints"].([]interface{}); ok {
			for _, endpointItem := range endpoints {
				if endpointMap, ok := endpointItem.(map[string]interface{}); ok {
					endpoint := make(map[string]interface{})
					if v, ok := endpointMap["Type"]; ok {
						endpoint["type"] = v
					}
					if v, ok := endpointMap["Endpoint"]; ok {
						endpoint["endpoint"] = v
					}
					if v, ok := endpointMap["Enabled"]; ok {
						endpoint["enabled"] = v
					}
					if v, ok := endpointMap["VpcId"]; ok {
						endpoint["vpc_id"] = v
					}
					if v, ok := endpointMap["VSwitchId"]; ok {
						endpoint["vswitch_id"] = v
					}
					if v, ok := endpointMap["VpcInstanceId"]; ok {
						endpoint["vpc_instance_id"] = v
					}
					endpointList = append(endpointList, endpoint)
				}
			}
		}

		mapping := map[string]interface{}{
			"id":                    instance["InstanceId"],
			"instance_name":         instance["InstanceName"],
			"instance_status":       instance["InstanceStatus"],
			"compute_type":          instance["InstanceType"],
			"creation_time":         instance["CreationTime"],
			"version":               instance["Version"],
			"enable_hive_access":    instance["EnableHiveAccess"],
			"cluster":               instance["Cluster"],
			"cpu":                   instance["Cpu"],
			"instance_type":         instance["InstanceType"],
			"instance_charge_type":  instance["InstanceChargeType"],
			"cpu_arch":              instance["CpuArch"],
			"cpu_brand":             instance["CpuBrand"],
			"cpu_brand_i18n":        instance["CpuBrandI18n"],
			"apsara_ascm_cpu_brand": instance["ApsaraAscmCpuBrand"],
			"support_replica":       instance["SupportReplica"],
			"ascm_create_user":      instance["AscmCreateUser"],
			"commodity_code":        instance["CommodityCode"],
			"endpoints":             endpointList,
		}

		ids = append(ids, instance["InstanceId"].(string))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("instances", s); err != nil {
		return errmsgs.WrapError(err)
	}

	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
