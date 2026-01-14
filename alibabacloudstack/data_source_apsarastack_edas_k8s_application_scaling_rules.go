package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackEdasScalingRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackEdasScalingRulesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				ForceNew:     true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"app_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"scaling_rule_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"metric", "trigger"}, false),
			},
			"scaling_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"app_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"scaling_rule_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scaling_rule_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"max_replicas": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"min_replicas": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"metrics": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"utilization": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"triggers": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"period": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"timer_in_day": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"at_time": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"replicas": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"horizon_mode": {
													Type:     schema.TypeBool,
													Computed: true,
												},
											},
										},
									},
									"timer_in_week": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"timer_in_month": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"scale_up_stabilization_window_seconds": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"scale_up_select_policy": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"scale_up_policies": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"period_seconds": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"scale_down_stabilization_window_seconds": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"scale_down_select_policy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scale_down_policies": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"period_seconds": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackEdasScalingRulesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := map[string]interface{}{
		"AppId": d.Get("app_id"),
	}
	resp, err := client.DoTeaRequest("GET", "Edas", "2017-08-01", "DescribeApplicationScalingRules", "/pop/v1/eam/scale/application_scaling_rules", nil, request, nil)
	if err != nil {
		return errmsgs.WrapError(err)

	}
	result, err := jsonpath.Get("$.Data.result", resp)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, "alibabacloudstack_edas_k8s_application_scaling_rules", "$.Data.result", resp)
	}

	idsMap := getIdsStringFilter(d)
	ids := make([]string, 0)
	scaling_rules := make([]map[string]interface{}, 0)
	for _, data := range result.([]interface{}) {
		object := data.(map[string]interface{})
		key := fmt.Sprintf("%s:%s", object["appId"].(string), object["scaleRuleName"].(string))
		if len(idsMap) > 0 {
			if _, ok := idsMap[key]; !ok {
				continue
			}
		}

		if nameRegex, ok := d.GetOk("name_regex"); ok {
			r := regexp.MustCompile(nameRegex.(string))
			if !r.MatchString(object["scaleRuleName"].(string)) {
				continue
			}
		}

		if scaleRuleType, ok := d.GetOk("scaling_rule_type"); ok && scaleRuleType.(string) != "" && scaleRuleType.(string) != object["scaleRuleType"].(string) {
			continue
		}
		scalingRuleMetrics, ok := object["metric"].(map[string]interface{})
		metrics := make([]map[string]interface{}, 0)
		if ok {
			for _, metric := range scalingRuleMetrics["metrics"].([]interface{}) {
				metrics = append(metrics, map[string]interface{}{
					"type":        metric.(map[string]interface{})["metricType"],
					"utilization": metric.(map[string]interface{})["metricTargetAverageUtilization"],
				})
			}
		}
		scalingRuleTrigger, ok := object["trigger"]
		triggers := make([]map[string]interface{}, 0)
		if ok {
			triggerDatas := scalingRuleTrigger.(map[string]interface{})["triggers"].([]interface{})
			for _, v := range triggerDatas {
				trigger := make(map[string]interface{})
				data := v.(map[string]interface{})
				metadata := make(map[string]interface{})
				mdate := data["metadata"].(string)
				err := json.Unmarshal([]byte(mdate), &metadata)
				if err != nil {
					return errmsgs.WrapErrorf(err, "Failed to unmarshal trigger metadata: %s", mdate)
				}
				trigger["name"] = data["name"]
				trigger["type"] = data["type"]
				trigger["period"] = metadata["period"]
				if v, ok := metadata["timerInWeek"]; ok {
					trigger["timer_in_week"] = v
				}
				if v, ok := metadata["timerInMonth"]; ok {
					trigger["timer_in_month"] = v
				}
				if v, ok := metadata["timerInDay"]; ok {
					timer_in_day := make([]map[string]interface{}, 0)
					days := v.([]interface{})
					for _, d := range days {
						day := d.(map[string]interface{})
						timer_in_day = append(timer_in_day, map[string]interface{}{
							"at_time":      day["atTime"],
							"horizon_mode": day["horizonMode"] == "true",
							"replicas":     day["targetReplicas"],
						})
					}
					trigger["timer_in_day"] = timer_in_day
				}
				triggers = append(triggers, trigger)

			}
		}
		mapping := map[string]interface{}{
			"id":                key,
			"app_id":            object["appId"],
			"scaling_rule_name": object["scaleRuleName"],
			"scaling_rule_type": object["scaleRuleType"],
			"max_replicas":      object["maxReplicas"],
			"min_replicas":      object["minReplicas"],
			"metrics":           metrics,
			"triggers":          triggers,
			"enabled":           object["scaleRuleEnabled"],
		}
		behaviour, ok := object["behaviour"]
		if ok {
			behaviourData := behaviour.(map[string]interface{})
			scale_up, ok := behaviourData["scaleUp"]
			if ok {
				scaleUpData := scale_up.(map[string]interface{})
				mapping["scale_up_stabilization_window_seconds"] = scaleUpData["stabilizationWindowSeconds"]
				mapping["scale_up_select_policy"] = scaleUpData["selectPolicy"]
				policies := make([]map[string]interface{}, 0)
				for _, v := range scaleUpData["policies"].([]interface{}) {
					policy := v.(map[string]interface{})
					policies = append(policies, map[string]interface{}{
						"period_seconds": policy["periodSeconds"],
						"value":          policy["value"],
						"type":           policy["type"],
					})
				}
				mapping["scale_up_policies"] = policies
			}
			scale_down, ok := behaviourData["scaleDown"]
			if ok {
				scaleDownData := scale_down.(map[string]interface{})
				d.Set("scale_down_stabilization_window_seconds", scaleDownData["stabilizationWindowSeconds"])
				d.Set("scale_down_select_policy", scaleDownData["selectPolicy"])
				mapping["scale_down_stabilization_window_seconds"] = scaleDownData["stabilizationWindowSeconds"]
				mapping["scale_down_select_policy"] = scaleDownData["selectPolicy"]
				policies := make([]map[string]interface{}, 0)
				for _, v := range scaleDownData["policies"].([]interface{}) {
					policy := v.(map[string]interface{})
					policies = append(policies, map[string]interface{}{
						"period_seconds": policy["periodSeconds"],
						"value":          policy["value"],
						"type":           policy["type"],
					})
				}
				mapping["scale_down_policies"] = policies
			}
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		scaling_rules = append(scaling_rules, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}

	if err := d.Set("scaling_rules", scaling_rules); err != nil {
		return errmsgs.WrapError(err)
	}
	return nil
}
