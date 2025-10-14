package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
			"trigger_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"trigger_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"trigger_period": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"daily", "weekly"}, false),
			},
			"trigger_dryrun": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"trigger_timer_in_day": {
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
					},
				},
			},
			"trigger_timer_in_week": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
	scalingRuleMetric := make(map[string]interface{})
	scalingRuleTrigger := make(map[string]interface{})
	scaling_rule_type := d.Get("scaling_rule_type").(string)
	if scaling_rule_type == "metric" {
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
	} else {
		scalingRuleTrigger["maxReplicas"] = d.Get("max_replicas")
		scalingRuleTrigger["minReplicas"] = d.Get("min_replicas")
		trigger_type, ok := d.GetOk("trigger_type")
		if !ok {
			trigger_type = "cron"
		}
		triggers := map[string]interface{}{
			"type": trigger_type,
			"name": d.Get("trigger_name"),
		}
		metadata := map[string]interface{}{
			"period": d.Get("trigger_period"),
			"dryRun": fmt.Sprint(d.Get("trigger_dryrun")),
		}
		timerInDay := make([]map[string]interface{}, 0)
		for _, v := range d.Get("trigger_timer_in_day").(*schema.Set).List() {
			timer := v.(map[string]interface{})
			timerInDay = append(timerInDay, map[string]interface{}{
				"atTime":         timer["at_time"],
				"targetReplicas": timer["replicas"],
			})
		}
		metadata["timerInDay"] = timerInDay
		if d.Get("trigger_period").(string) == "weekly" {
			metadata["timeInWeek"] = d.Get("trigger_timer_in_week").(*schema.Set).List()
		}
		triggers["metadata"] = metadata
		scalingRuleTrigger["triggers"] = triggers
	}
	metric, _ := json.Marshal(scalingRuleMetric)
	trigger, err := json.Marshal(scalingRuleTrigger)
	if err != nil {
		return fmt.Errorf("scalingRuleTrigger data to marshal JSON failed: %w \n%v", err, scalingRuleTrigger)
	}
	request := map[string]interface{}{
		"AppId":              d.Get("app_id"),
		"ScalingRuleName":    d.Get("scaling_rule_name"),
		"ScalingRuleType":    scaling_rule_type,
		"ScalingRuleMetric":  string(metric),
		"ScalingRuleTrigger": string(trigger),
	}
	_, err = client.DoTeaRequest("POST", "Edas", "2017-08-01", "CreateApplicationScalingRule", "/pop/v1/eam/scale/application_scaling_rule", nil, request, nil)
	if err != nil {
		return err
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
		"trigger_timer_in_day", "trigger_timer_in_week") {

		scalingRuleMetric := make(map[string]interface{})
		scalingRuleTrigger := make(map[string]interface{})
		scaling_rule_type := d.Get("scaling_rule_type").(string)
		if scaling_rule_type == "metric" {
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
		} else {
			scalingRuleTrigger["maxReplicas"] = d.Get("max_replicas")
			scalingRuleTrigger["minReplicas"] = d.Get("min_replicas")
			triggers := map[string]interface{}{
				"type": d.Get("trigger_type"),
				"name": d.Get("trigger_name"),
			}
			metadata := map[string]interface{}{
				"period": d.Get("trigger_period"),
				"dryRun": fmt.Sprint(d.Get("trigger_dryrun")),
			}
			timerInDay := make([]map[string]interface{}, 0)
			for _, v := range d.Get("trigger_timer_in_day").(*schema.Set).List() {
				timer := v.(map[string]interface{})
				timerInDay = append(timerInDay, map[string]interface{}{
					"atTime":         timer["at_time"],
					"targetReplicas": timer["replicas"],
				})
			}
			metadata["timerInDay"] = timerInDay
			if d.Get("trigger_period").(string) == "weekly" {
				metadata["timeInWeek"] = d.Get("trigger_timer_in_week").(*schema.Set).List()
			}
			triggers["metadata"] = metadata
			scalingRuleTrigger["triggers"] = triggers
		}
		metric, _ := json.Marshal(scalingRuleMetric)
		trigger, _ := json.Marshal(scalingRuleTrigger)
		request := map[string]interface{}{
			"AppId":              d.Get("app_id"),
			"ScalingRuleName":    d.Get("scaling_rule_name"),
			"ScalingRuleType":    scaling_rule_type,
			"ScalingRuleMetric":  string(metric),
			"scalingRuleTrigger": string(trigger),
		}
		_, err := client.DoTeaRequest("PUT", "Edas", "2017-08-01", "UpdateApplicationScalingRule", "/pop/v1/eam/scale/application_scaling_rule", nil, request, nil)
		if err != nil {
			return err
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
		triggers := scalingRuleTrigger.(map[string]interface{})["triggers"].([]interface{})
		trigger := triggers[0].(map[string]interface{})
		d.Set("trigger_type", trigger["type"])
		d.Set("trigger_name", trigger["name"])
		matedata := make(map[string]interface{})
		_ = json.Unmarshal([]byte(trigger["metadata"].(string)), &matedata)
		d.Set("trigger_period", matedata["period"])
		d.Set("trigger_dryrun", matedata["dryRun"] == "true")
		timerInWeek, ok := matedata["timeInWeek"]
		if ok {
			d.Set("trigger_timer_in_week", timerInWeek)
		}
		timerInDay, ok := matedata["timerInDay"]
		timers := make([]map[string]interface{}, 0)
		if ok {
			for _, v := range timerInDay.([]interface{}) {
				timer := v.(map[string]interface{})
				timers = append(timers, map[string]interface{}{
					"at_time":  timer["atTime"],
					"replicas": timer["targetReplicas"],
				})
			}
			d.Set("trigger_timer_in_day", timers)
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
	wait := incrementalWait(1*time.Second, 2*time.Second)
	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		bresponse, err := client.ProcessCommonRequest(request)
		addDebug(request.GetActionName(), bresponse, request)
		if err != nil {
			if bresponse == nil {
				return resource.RetryableError(err)
			}
			if errmsgs.IsExpectedErrors(err, []string{errmsgs.ThrottlingUser}) {
				wait()
				return resource.RetryableError(err)
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_edas_k8s_application_scaling_rule", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		return nil
	})
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return nil
}
