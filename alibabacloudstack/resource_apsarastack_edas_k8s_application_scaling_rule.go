package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackEdasK8sApplicationScalingRule() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"scaling_rule_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"scaling_rule_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"metric", "trigger"}, false),
			},
			"max_replicas": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(2, 100),
			},
			"min_replicas": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 100),
			},
			"metrics": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"CPU", "MEMORY"}, false),
						},
						"utilization": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 100),
						},
					},
				},
			},
			"triggers": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "cron",
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"period": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"daily", "weekly", "monthly"}, false),
						},
						"timer_in_day": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"at_time": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "like '08:00'",
									},
									"replicas": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(1, 100),
									},
									"horizon_mode": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
								},
							},
						},
						"timer_in_week": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"timer_in_month": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"scale_up_stabilization_window_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 3600),
			},
			"scale_up_select_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Min", "Max", "Disabled"}, false),
			},

			"scale_up_policies": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"Percent", "Pods"}, false),
						},
						"value": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 100),
						},
						"period_seconds": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 3600),
						},
					},
				},
			},
			"scale_down_stabilization_window_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 3600),
			},
			"scale_down_select_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Min", "Max", "Disabled"}, false),
			},
			"scale_down_policies": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"Percent", "Pods"}, false),
						},
						"value": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 100),
						},
						"period_seconds": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 3600),
						},
					},
				},
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
	setResourceFunc(resource,
		resourceAlibabacloudStackEdasK8sApplicationScalingRuleCreate,
		resourceAlibabacloudStackEdasK8sApplicationScalingRuleRead,
		resourceAlibabacloudStackEdasK8sApplicationScalingRuleUpdate,
		resourceAlibabacloudStackEdasK8sApplicationScalingRuleDelete)
	return resource
}

func resourceAlibabacloudStackEdasK8sApplicationScalingRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	scaling_rule_type := d.Get("scaling_rule_type").(string)
	behaviour, err := BuildScalingBehaviour(d)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	trigger, err := BuildScalingRuleTrigger(d)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	metric, err := BuildscalingRuleMetric(d)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	request := map[string]interface{}{
		"AppId":           d.Get("app_id"),
		"ScalingRuleName": d.Get("scaling_rule_name"),
		"ScalingRuleType": scaling_rule_type,
	}
	if scaling_rule_type == "trigger" {
		request["ScalingRuleTrigger"] = trigger
	} else {
		request["ScalingRuleMetric"] = metric
	}
	if behaviour != "" && behaviour != "{}" {
		request["ScalingBehaviour"] = behaviour
	}
	response, err := client.DoTeaRequest("POST", "Edas", "2017-08-01", "CreateApplicationScalingRule", "/pop/v1/eam/scale/application_scaling_rule", nil, request, nil)
	if err != nil {
		return err
	}
	if fmt.Sprint(response["Code"]) != "200" {
		return errmsgs.Error("Create k8s application scaling rule failed for " + response["Message"].(string))
	}
	resourceId := fmt.Sprintf("%s:%s", d.Get("app_id").(string), d.Get("scaling_rule_name").(string))
	d.SetId(resourceId)
	return nil
}

func resourceAlibabacloudStackEdasK8sApplicationScalingRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	if d.HasChange("enabled") {
		params := strings.Split(d.Id(), ":")
		req := map[string]interface{}{
			"AppId":           params[0],
			"ScalingRuleName": params[1],
		}
		enabled := d.Get("enabled").(bool)
		todo := true
		var action, pattern string
		if enabled {
			action = "EnableApplicationScalingRule"
			pattern = "/pop/v1/eam/scale/enable_application_scaling_rule"
		} else {
			action = "DisableApplicationScalingRule"
			pattern = "/pop/v1/eam/scale/disable_application_scaling_rule"
			if d.IsNewResource() {
				todo = false
			}

		}
		if todo {
			_, err := client.DoTeaRequest("PUT", "Edas", "2017-08-01", action, pattern, nil, req, nil)
			if err != nil {
				return err
			}
		}
	}
	if !d.IsNewResource() && d.HasChanges("metrics", "max_replicas", "min_replicas",
		"trigger_type", "trigger_name", "trigger_period", "trigger_dryrun",
		"trigger_timer_in_day", "trigger_timer_in_week", "scale_up_stabilization_window_seconds",
		"scale_up_select_policy", "scale_up_policies", "scale_down_stabilization_window_seconds",
		"scale_down_select_policy", "scale_down_policies") {
		scaling_rule_type := d.Get("scaling_rule_type").(string)
		behaviour, err := BuildScalingBehaviour(d)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		trigger, err := BuildScalingRuleTrigger(d)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		metric, err := BuildscalingRuleMetric(d)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		request := map[string]interface{}{
			"AppId":           d.Get("app_id"),
			"ScalingRuleName": d.Get("scaling_rule_name"),
			"ScalingRuleType": scaling_rule_type,
		}
		if scaling_rule_type == "trigger" {
			request["ScalingRuleTrigger"] = trigger
		} else {
			request["ScalingRuleMetric"] = metric
		}
		if behaviour != "" && behaviour != "{}" {
			request["ScalingBehaviour"] = behaviour
		}
		response, err := client.DoTeaRequest("PUT", "Edas", "2017-08-01", "UpdateApplicationScalingRule", "/pop/v1/eam/scale/application_scaling_rule", nil, request, nil)
		if err != nil {
			return err
		}
		if fmt.Sprint(response["Code"]) != "200" {
			return errmsgs.Error("Update k8s application scaling rule failed for " + response["Message"].(string))
		}
	}

	return nil
}

func resourceAlibabacloudStackEdasK8sApplicationScalingRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	edasService := EdasService{client}

	object, err := edasService.DescribeEdasK8sApplicationScalingRule(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	d.Set("app_id", object["appId"])
	d.Set("scaling_rule_name", object["scaleRuleName"])
	d.Set("scaling_rule_type", object["scaleRuleType"])
	d.Set("max_replicas", object["maxReplicas"])
	d.Set("min_replicas", object["minReplicas"])
	scalingRuleMetrics, ok := object["metric"].(map[string]interface{})
	if ok {
		metrics := make([]map[string]interface{}, 0)
		for _, metric := range scalingRuleMetrics["metrics"].([]interface{}) {
			metrics = append(metrics, map[string]interface{}{
				"type":        metric.(map[string]interface{})["metricType"],
				"utilization": metric.(map[string]interface{})["metricTargetAverageUtilization"],
			})
		}
		d.Set("metrics", metrics)
	}
	scalingRuleTrigger, ok := object["trigger"]
	if ok {
		triggerDatas := scalingRuleTrigger.(map[string]interface{})["triggers"].([]interface{})
		triggers := make([]map[string]interface{}, 0)
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
		d.Set("triggers", triggers)
	}
	behaviour, ok := object["behaviour"]
	if ok {
		behaviourData := behaviour.(map[string]interface{})
		scale_up, ok := behaviourData["scaleUp"]
		if ok {
			scaleUpData := scale_up.(map[string]interface{})
			d.Set("scale_up_stabilization_window_seconds", scaleUpData["stabilizationWindowSeconds"])
			d.Set("scale_up_select_policy", scaleUpData["selectPolicy"])
			policies := make([]map[string]interface{}, 0)
			for _, v := range scaleUpData["policies"].([]interface{}) {
				policy := v.(map[string]interface{})
				policies = append(policies, map[string]interface{}{
					"period_seconds": policy["periodSeconds"],
					"value":          policy["value"],
					"type":           policy["type"],
				})
			}
			d.Set("scale_up_policies", policies)
		}
		scale_down, ok := behaviourData["scaleDown"]
		if ok {
			scaleDownData := scale_down.(map[string]interface{})
			d.Set("scale_down_stabilization_window_seconds", scaleDownData["stabilizationWindowSeconds"])
			d.Set("scale_down_select_policy", scaleDownData["selectPolicy"])
			policies := make([]map[string]interface{}, 0)
			for _, v := range scaleDownData["policies"].([]interface{}) {
				policy := v.(map[string]interface{})
				policies = append(policies, map[string]interface{}{
					"period_seconds": policy["periodSeconds"],
					"value":          policy["value"],
					"type":           policy["type"],
				})
			}
			d.Set("scale_down_policies", policies)
		}
	}
	d.Set("enabled", object["scaleRuleEnabled"])

	return nil
}

func resourceAlibabacloudStackEdasK8sApplicationScalingRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := client.NewCommonRequest("DELETE", "Edas", "2017-08-01", "DeleteApplicationScalingRule", "/pop/v1/eam/scale/application_scaling_rule")
	params := strings.Split(d.Id(), ":")
	request.QueryParams["AppId"] = params[0]
	request.QueryParams["ScalingRuleName"] = params[1]
	// request.Headers["x-acs-content-type"] = "application/x-www-form-urlencoded"

	bresponse, err := client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request)
	if err != nil {
		if bresponse == nil {
			return err
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_edas_k8s_application_scaling_rule", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	return nil
}

func BuildScalingBehaviour(d *schema.ResourceData) (string, error) {
	scalingBehaviour := make(map[string]interface{})
	scale_up_stabilization_window_seconds := d.Get("scale_up_stabilization_window_seconds").(int)
	scale_up_select_policy := d.Get("scale_up_select_policy").(string)
	scale_up_policies := d.Get("scale_up_policies").(*schema.Set).List()
	scale_down_stabilization_window_seconds := d.Get("scale_down_stabilization_window_seconds").(int)
	scale_down_select_policy := d.Get("scale_down_select_policy").(string)
	scale_down_policies := d.Get("scale_down_policies").(*schema.Set).List()

	// Verify the scale_up attribute: if any one is set, all three must be set
	scaleUpSet := scale_up_stabilization_window_seconds > 0 || scale_up_select_policy != "" || len(scale_up_policies) > 0
	scaleDownSet := scale_down_stabilization_window_seconds > 0 || scale_down_select_policy != "" || len(scale_down_policies) > 0
	if scaleUpSet || scaleDownSet {
		if scale_up_stabilization_window_seconds <= 0 || scale_up_select_policy == "" || len(scale_up_policies) == 0 || scale_down_stabilization_window_seconds <= 0 || scale_down_select_policy == "" || len(scale_down_policies) == 0 {
			return "", fmt.Errorf("`scale_up_stabilization_window_seconds`, `scale_up_select_policy`, `scale_up_policies`, `scale_down_stabilization_window_seconds`, `scale_down_select_policy`, `scale_down_policies`  must be set together")
		}

		scaleUp := make(map[string]interface{})
		scaleUp["stabilizationWindowSeconds"] = scale_up_stabilization_window_seconds
		scaleUp["selectPolicy"] = scale_up_select_policy
		if len(scale_up_policies) > 0 {
			policies := make([]map[string]interface{}, 0)
			for _, p := range scale_up_policies {
				policy := p.(map[string]interface{})
				policies = append(policies, map[string]interface{}{
					"type":          policy["type"],
					"value":         policy["value"],
					"periodSeconds": policy["period_seconds"],
				})
			}
			scaleUp["policies"] = policies
		}
		scalingBehaviour["scaleUp"] = scaleUp

		scaleDown := make(map[string]interface{})
		scaleDown["stabilizationWindowSeconds"] = scale_down_stabilization_window_seconds
		scaleDown["selectPolicy"] = scale_down_select_policy
		if len(scale_down_policies) > 0 {
			policies := make([]map[string]interface{}, 0)
			for _, p := range scale_down_policies {
				policy := p.(map[string]interface{})
				policies = append(policies, map[string]interface{}{
					"type":          policy["type"],
					"value":         policy["value"],
					"periodSeconds": policy["period_seconds"],
				})
			}
			scaleDown["policies"] = policies
		}
		scalingBehaviour["scaleDown"] = scaleDown
	}

	behaviour, _ := json.Marshal(scalingBehaviour)
	return string(behaviour), nil
}

func BuildScalingRuleTrigger(d *schema.ResourceData) (string, error) {
	ruleTriggers := make([]map[string]interface{}, 0)
	triggers := d.Get("triggers").(*schema.Set).List()
	for _, t := range triggers {
		triggerData := t.(map[string]interface{})
		triggerInDay := make([]map[string]interface{}, 0)
		timer_in_day := triggerData["timer_in_day"].(*schema.Set).List()
		timer_in_week := triggerData["timer_in_week"].(*schema.Set).List()
		timer_in_month := triggerData["timer_in_month"].(*schema.Set).List()
		timer_type := triggerData["type"].(string)
		timer_name := triggerData["name"].(string)
		timer_period := triggerData["period"].(string)
		for _, v := range timer_in_day {
			timer := v.(map[string]interface{})
			triggerInDay = append(triggerInDay, map[string]interface{}{
				"atTime":         timer["at_time"],
				"targetReplicas": timer["replicas"],
				"horizonMode":    fmt.Sprint(timer["horizon_mode"]),
			})
		}
		metadata := map[string]interface{}{
			"period":     timer_period,
			"timerInDay": triggerInDay,
		}
		if timer_period == "monthly" {
			metadata["timerInMonth"] = timer_in_month
		} else if timer_period == "weekly" {
			metadata["timerInWeek"] = timer_in_week
		}
		trigger := map[string]interface{}{
			"type":     timer_type,
			"name":     timer_name,
			"metadata": metadata,
		}
		ruleTriggers = append(ruleTriggers, trigger)
	}
	data := map[string]interface{}{
		"maxReplicas": d.Get("max_replicas").(int),
		"minReplicas": d.Get("min_replicas").(int),
		"triggers":    ruleTriggers,
	}
	triggerStr, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(triggerStr), nil
}

func BuildscalingRuleMetric(d *schema.ResourceData) (string, error) {
	scalingRuleMetric := make(map[string]interface{})
	scalingRuleMetric["maxReplicas"] = d.Get("max_replicas")
	scalingRuleMetric["minReplicas"] = d.Get("min_replicas")
	metrics := d.Get("metrics").(*schema.Set).List()
	metricData := make([]map[string]interface{}, 0)
	for _, m := range metrics {
		metric := m.(map[string]interface{})
		metricData = append(metricData, map[string]interface{}{
			"metricType":                     metric["type"],
			"metricTargetAverageUtilization": metric["utilization"],
		})
	}
	scalingRuleMetric["metrics"] = metricData
	metricStr, err := json.Marshal(scalingRuleMetric)
	if err != nil {
		return "", err
	}
	return string(metricStr), nil
}
