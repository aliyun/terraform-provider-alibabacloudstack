package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackCmsAlarm() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"rule_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
			},
			"rule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'name' is deprecated and will be removed in a future release. Please use new field 'rule_name' instead.",
				ConflictsWith: []string{"rule_name"},
			},
			"namespace": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"project"},
			},
			"project": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				Deprecated:    "Field 'project' is deprecated and will be removed in a future release. Please use new field 'namespace' instead.",
				ConflictsWith: []string{"namespace"},
			},
			"metric_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"metric"},
			},
			"metric": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				Deprecated:    "Field 'metric' is deprecated and will be removed in a future release. Please use new field 'metric_name' instead.",
				ConflictsWith: []string{"metric_name"},
			},
			"dimensions": {
				Type:          schema.TypeMap,
				Optional:      true,
				ForceNew:      true,
				Elem:          schema.TypeString,
				Deprecated:    "Field 'dimensions' is deprecated and will be removed in a future release. Please use new field 'resources' instead.",
				ConflictsWith: []string{"resources"},
			},
			"resources": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					ValidateFunc: func(i interface{}, k string) ([]string, []error) {
						m := i.(map[string]interface{})
						if len(m) > 1 {
							return nil, []error{fmt.Errorf("too large map")}
						}
						return nil, nil
					},
				},
			},
			"period": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"escalations_critical": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"statistics": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      Average,
							ValidateFunc: validation.StringInSlice([]string{Average, Minimum, Maximum, ErrorCodeMaximum, Value}, false),
						},
						"comparison_operator": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  Equal,
							ValidateFunc: validation.StringInSlice([]string{
								MoreThan, MoreThanOrEqual, LessThan, LessThanOrEqual, NotEqual,
							}, false),
						},
						"threshold": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"times": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  3,
						},
					},
				},
				DiffSuppressFunc: cmsClientCriticalSuppressFunc,
				MaxItems:         1,
			},
			"escalations_warn": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"statistics": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      Average,
							ValidateFunc: validation.StringInSlice([]string{Average, Minimum, Maximum, ErrorCodeMaximum, Value}, false),
						},
						"comparison_operator": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  Equal,
							ValidateFunc: validation.StringInSlice([]string{
								MoreThan, MoreThanOrEqual, LessThan, LessThanOrEqual, NotEqual,
							}, false),
						},
						"threshold": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"times": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  3,
						},
					},
				},
				DiffSuppressFunc: cmsClientWarnSuppressFunc,
				MaxItems:         1,
			},
			"escalations_info": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"statistics": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      Average,
							ValidateFunc: validation.StringInSlice([]string{Average, Minimum, Maximum, ErrorCodeMaximum, Value}, false),
						},
						"comparison_operator": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  Equal,
							ValidateFunc: validation.StringInSlice([]string{
								MoreThan, MoreThanOrEqual, LessThan, LessThanOrEqual, NotEqual,
							}, false),
						},
						"threshold": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"times": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  3,
						},
					},
				},
				DiffSuppressFunc: cmsClientInfoSuppressFunc,
				MaxItems:         1,
			},
			"contact_groups": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"effective_interval": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "00:00-23:59",
			},
			"silence_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      86400,
				ValidateFunc: validation.IntBetween(300, 86400),
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"webhook": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackCmsAlarmCreate, resourceAlibabacloudStackCmsAlarmRead, nil, resourceAlibabacloudStackCmsAlarmDelete)
	return resource
}

func resourceAlibabacloudStackCmsAlarmCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	cmsService := CmsService{client}
	d.Partial(true)

	request := cms.CreatePutResourceMetricRuleRequest()
	client.InitRpcRequest(*request.RpcRequest)
	d.SetId(resource.UniqueId() + ":" + request.RuleName)
	request.RuleName = connectivity.GetResourceData(d, "rule_name", "name").(string)
	if err := errmsgs.CheckEmpty(request.RuleName, schema.TypeString, "rule_name", "name"); err != nil {
		return errmsgs.WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	request.RuleId = parts[0]

	request.Namespace = connectivity.GetResourceData(d, "namespace", "project").(string)
	if err := errmsgs.CheckEmpty(request.Namespace, schema.TypeString, "namespace", "project"); err != nil {
		return errmsgs.WrapError(err)
	}
	request.MetricName = connectivity.GetResourceData(d, "metric_name", "metric").(string)
	if err := errmsgs.CheckEmpty(request.MetricName, schema.TypeString, "metric_name", "metric"); err != nil {
		return errmsgs.WrapError(err)
	}
	request.Period = strconv.Itoa(d.Get("period").(int))

	request.ContactGroups = strings.Join(expandStringList(d.Get("contact_groups").([]interface{})), ",")
	if v, ok := d.GetOk("escalations_critical"); ok && len(v.([]interface{})) != 0 {
		for _, val := range v.([]interface{}) {
			val := val.(map[string]interface{})
			request.EscalationsCriticalStatistics = val["statistics"].(string)
			request.EscalationsCriticalComparisonOperator = convertOperator(val["comparison_operator"].(string))
			request.EscalationsCriticalThreshold = val["threshold"].(string)
			request.EscalationsCriticalTimes = requests.NewInteger(val["times"].(int))
		}
	}
	// Warn
	if v, ok := d.GetOk("escalations_warn"); ok && len(v.([]interface{})) != 0 {
		for _, val := range v.([]interface{}) {
			val := val.(map[string]interface{})
			request.EscalationsWarnStatistics = val["statistics"].(string)
			request.EscalationsWarnComparisonOperator = convertOperator(val["comparison_operator"].(string))
			request.EscalationsWarnThreshold = val["threshold"].(string)
			request.EscalationsWarnTimes = requests.NewInteger(val["times"].(int))
		}
	}
	// Info
	if v, ok := d.GetOk("escalations_info"); ok && len(v.([]interface{})) != 0 {
		for _, val := range v.([]interface{}) {
			val := val.(map[string]interface{})
			request.EscalationsInfoStatistics = val["statistics"].(string)
			request.EscalationsInfoComparisonOperator = convertOperator(val["comparison_operator"].(string))
			request.EscalationsInfoThreshold = val["threshold"].(string)
			request.EscalationsInfoTimes = requests.NewInteger(val["times"].(int))
		}
	}

	if v, ok := d.GetOk("effective_interval"); ok && v.(string) != "" {
		request.EffectiveInterval = v.(string)
	} else {
		start, startOk := d.GetOk("start_time")
		end, endOk := d.GetOk("end_time")
		if startOk && endOk && end.(int) > 0 {
			// The EffectiveInterval valid value between 00:00 and 23:59
			request.EffectiveInterval = fmt.Sprintf("%d:00-%d:59", start.(int), end.(int)-1)
		}
	}
	request.SilenceTime = requests.NewInteger(d.Get("silence_time").(int))

	var instanceId string
	var dimList []map[string]string

	if dimensions, ok := d.GetOk("dimensions"); ok {
		for k, v := range dimensions.(map[string]interface{}) {
			values := strings.Split(v.(string), COMMA_SEPARATED)
			if len(values) > 0 {
				instanceId = values[0]
				for _, vv := range values {
					dimList = append(dimList, map[string]string{k: Trim(vv)})
				}
			} else {
				dimList = append(dimList, map[string]string{k: Trim(v.(string))})
			}
		}
	} else if resources, ok := d.GetOk("resources"); ok {
		for _, item := range resources.([]interface{}) {
			for k, v := range item.(map[string]interface{}) {
				dimList = append(dimList, map[string]string{k: Trim(v.(string))})
			}
		}
	} else {
		return fmt.Errorf("dimensions and resources can not be empty at the same time")
	}

	if len(dimList) > 0 {
		if bytes, err := json.Marshal(dimList); err != nil {
			return fmt.Errorf("marshaling dimensions to json string got an error: %#v", err)
		} else {
			request.Resources = string(bytes[:])
		}
	}

	nrequest := client.NewCommonRequest("POST", "Cms", "2019-01-01", "PutResourceMetricRule", "")
	mergeMaps(nrequest.QueryParams, map[string]string{
		"RuleName":                       request.RuleName,
		"RuleId":                         request.RuleId,
		"Namespace":                      request.Namespace,
		"MetricName":                     request.MetricName,
		"Period":                         request.Period,
		"EffectiveInterval":              request.EffectiveInterval,
		"ContactGroups":                  request.ContactGroups,
		"InstanceID":                     instanceId,
		"Resources":                      request.Resources,
		"Escalations.Critical.Threshold": request.EscalationsCriticalThreshold,
		"Escalations.Critical.ComparisonOperator": request.EscalationsCriticalComparisonOperator,
		"Escalations.Critical.Statistics":         request.EscalationsCriticalStatistics,
		"Escalations.Critical.Times":              fmt.Sprint(request.EscalationsCriticalTimes),
		"Escalations.Warn.Threshold":              request.EscalationsWarnThreshold,
		"Escalations.Warn.ComparisonOperator":     request.EscalationsWarnComparisonOperator,
		"Escalations.Warn.Statistics":             request.EscalationsWarnStatistics,
		"Escalations.Warn.Times":                  fmt.Sprint(request.EscalationsWarnTimes),
		"Escalations.Info.Threshold":              request.EscalationsInfoThreshold,
		"Escalations.Info.ComparisonOperator":     request.EscalationsCriticalComparisonOperator,
		"Escalations.Info.Statistics":             request.EscalationsInfoStatistics,
		"Escalations.Info.Times":                  fmt.Sprint(request.EscalationsInfoTimes),
		"SilenceTime":                             fmt.Sprint(request.SilenceTime),
		"Format":                                  "JSON",
		"SignatureVersion":                        "1.0",
	})

	response, err := client.ProcessCommonRequest(nrequest)
	if err != nil {
		if response == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_cms", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	if d.Get("enabled").(bool) {
		request := cms.CreateEnableMetricRulesRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.RuleId = &[]string{d.Id()}

		wait := incrementalWait(1*time.Second, 2*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
				return cmsClient.EnableMetricRules(request)
			})
			response, ok := raw.(*responses.CommonResponse)
			if err != nil {
				errmsg := ""
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
				if errmsgs.IsExpectedErrors(err, []string{errmsgs.ThrottlingUser}) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_cms", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("Enabling alarm got an error: %#v", err)
		}
	} else if err != nil {
		return err
	} else {
		request := cms.CreateDisableMetricRulesRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.RuleId = &[]string{d.Id()}

		wait := incrementalWait(1*time.Second, 2*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
				return cmsClient.DisableMetricRules(request)
			})
			response, ok := raw.(*responses.CommonResponse)
			if err != nil {
				errmsg := ""
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
				if errmsgs.IsExpectedErrors(err, []string{errmsgs.ThrottlingUser}) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_cms", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("Disabling alarm got an error: %#v", err)
		}
	}
	if err := cmsService.WaitForCmsAlarm(d.Id(), d.Get("enabled").(bool), 102); err != nil {
		return err
	}

	d.Partial(false)

	return nil
}

func resourceAlibabacloudStackCmsAlarmRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	cmsService := CmsService{client}

	alarm, err := cmsService.DescribeCmsAlarm(d.Id())

	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	connectivity.SetResourceData(d, alarm.RuleName, "rule_name", "name")
	connectivity.SetResourceData(d, alarm.Namespace, "namespace", "project")
	connectivity.SetResourceData(d, alarm.MetricName, "metric_name", "metric")
	if period, err := strconv.Atoi(alarm.Period); err != nil {
		return errmsgs.WrapError(err)
	} else {
		d.Set("period", period)
	}

	escalationsCritical := make([]map[string]interface{}, 1)
	if alarm.Escalations.Critical.Times != 0 {
		mapping := map[string]interface{}{
			"statistics":          alarm.Escalations.Critical.Statistics,
			"comparison_operator": convertOperator(alarm.Escalations.Critical.ComparisonOperator),
			"threshold":           alarm.Escalations.Critical.Threshold,
			"times":               alarm.Escalations.Critical.Times,
		}
		escalationsCritical[0] = mapping
		d.Set("escalations_critical", escalationsCritical)
	}

	escalationsWarn := make([]map[string]interface{}, 1)
	if alarm.Escalations.Warn.Times != "" {
		if count, err := strconv.Atoi(alarm.Escalations.Warn.Times); err != nil {
			return errmsgs.WrapError(err)
		} else {
			mappingWarn := map[string]interface{}{
				"statistics":          alarm.Escalations.Warn.Statistics,
				"comparison_operator": convertOperator(alarm.Escalations.Warn.ComparisonOperator),
				"threshold":           alarm.Escalations.Warn.Threshold,
				"times":               count,
			}
			escalationsWarn[0] = mappingWarn
			d.Set("escalations_warn", escalationsWarn)
		}
	}

	escalationsInfo := make([]map[string]interface{}, 1)
	if alarm.Escalations.Info.Times != 0 {
		mappingInfo := map[string]interface{}{
			"statistics":          alarm.Escalations.Info.Statistics,
			"comparison_operator": convertOperator(alarm.Escalations.Info.ComparisonOperator),
			"threshold":           alarm.Escalations.Info.Threshold,
			"times":               alarm.Escalations.Info.Times,
		}
		escalationsInfo[0] = mappingInfo
		d.Set("escalations_info", escalationsInfo)
	}

	d.Set("rule_id", alarm.RuleId)
	d.Set("effective_interval", alarm.EffectiveInterval)
	d.Set("silence_time", alarm.SilenceTime)
	d.Set("status", alarm.AlertState)
	d.Set("enabled", alarm.EnableState)
	d.Set("contact_groups", strings.Split(alarm.ContactGroups, ","))

	//var dims map[string]interface{}
	// TODO: 当前接口无法正常返回Dimensions
	// if alarm.Dimensions != "" {
	// 	if err := json.Unmarshal([]byte(alarm.Dimensions), &dims); err != nil {
	// 		return fmt.Errorf("Unmarshaling Dimensions got an error: %#v.", err)
	// 	}
	// 	d.Set("dimensions", dims)
	// }

	return nil
}

func resourceAlibabacloudStackCmsAlarmDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	cmsService := CmsService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	request := cms.CreateDeleteMetricRulesRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.Id = &[]string{parts[0]}

	wait := incrementalWait(1*time.Second, 2*time.Second)
	return resource.Retry(10*time.Minute, func() *resource.RetryError {
		raw, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
			return cmsClient.DeleteMetricRules(request)
		})
		response, ok := raw.(*responses.CommonResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			if errmsgs.IsExpectedErrors(err, []string{errmsgs.ThrottlingUser}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_cms", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}

		_, err = cmsService.DescribeCmsAlarm(d.Id())
		if err != nil {
			if errmsgs.NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_cms", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, ""))
		}

		return resource.RetryableError(fmt.Errorf("Deleting alarm rule got an error: %#v", err))
	})
}

func convertOperator(operator string) string {
	switch operator {
	case MoreThan:
		return "GreaterThanThreshold"
	case MoreThanOrEqual:
		return "GreaterThanOrEqualToThreshold"
	case LessThan:
		return "LessThanThreshold"
	case LessThanOrEqual:
		return "LessThanOrEqualToThreshold"
	case NotEqual:
		return "NotEqualToThreshold"
	case Equal:
		return "GreaterThanThreshold"
	case "GreaterThanThreshold":
		return MoreThan
	case "GreaterThanOrEqualToThreshold":
		return MoreThanOrEqual
	case "LessThanThreshold":
		return LessThan
	case "LessThanOrEqualToThreshold":
		return LessThanOrEqual
	case "NotEqualToThreshold":
		return NotEqual
	default:
		return ""
	}
}