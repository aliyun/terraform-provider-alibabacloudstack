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
						"trigger_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"trigger_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"trigger_period": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"trigger_dryrun": {
							Type:     schema.TypeBool,
							Computed: true,
						},

						"trigger_timer_in_day": {
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
								},
							},
						},
						"trigger_timer_in_week": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
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

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
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
		timers := make([]map[string]interface{}, 0)
		var triggerType, triggerName, triggerPeriod string
		var dryRun bool
		timerInWeek := make([]interface{}, 0)
		if ok {
			triggers := scalingRuleTrigger.(map[string]interface{})["triggers"].([]interface{})
			trigger := triggers[0].(map[string]interface{})
			triggerType = trigger["type"].(string)
			triggerName = trigger["name"].(string)
			triggerPeriod = trigger["period"].(string)
			dryRun = trigger["dryRun"] == "true"
			matedata := make(map[string]interface{})
			_ = json.Unmarshal([]byte(trigger["metadata"].(string)), &matedata)
			v, ok := trigger["timerInWeek"]
			if ok {
				timerInWeek = v.([]interface{})
			}
			timerInDay, ok := trigger["timerInDay"]
			if ok {
				for _, v := range timerInDay.([]interface{}) {
					timer := v.(map[string]interface{})
					timers = append(timers, map[string]interface{}{
						"at_time":  timer["atTime"],
						"replicas": timer["targetReplicas"],
					})
				}
			}
		}
		mapping := map[string]interface{}{
			"id":                    key,
			"app_id":                object["appId"],
			"scaling_rule_name":     object["scaleRuleName"],
			"scaling_rule_type":     object["scaleRuleType"],
			"max_replicas":          object["maxReplicas"],
			"min_replicas":          object["minReplicas"],
			"metrics":               metrics,
			"trigger_timer_in_day":  timers,
			"trigger_type":          triggerType,
			"trigger_name":          triggerName,
			"trigger_period":        triggerPeriod,
			"trigger_dryrun":        dryRun,
			"trigger_timer_in_week": timerInWeek,
			"enabled":               object["scaleRuleEnabled"],
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
