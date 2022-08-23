package apsarastack

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/aliyun/terraform-provider-alibabacloudstack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApsaraStackCmsAlarm() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackCmsAlarmCreate,
		Read:   resourceApsaraStackCmsAlarmRead,
		Update: resourceApsaraStackCmsAlarmUpdate,
		Delete: resourceApsaraStackCmsAlarmDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"metric": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dimensions": {
				Type:     schema.TypeMap,
				Required: true,
				ForceNew: true,
				Elem:     schema.TypeString,
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
							ValidateFunc: validation.StringInSlice([]string{Average, Minimum, Maximum, ErrorCodeMaximum}, false),
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
							ValidateFunc: validation.StringInSlice([]string{Average, Minimum, Maximum, ErrorCodeMaximum}, false),
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
							ValidateFunc: validation.StringInSlice([]string{Average, Minimum, Maximum, ErrorCodeMaximum}, false),
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
}

func resourceApsaraStackCmsAlarmCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.ApsaraStackClient)
	cmsService := CmsService{client}
	d.Partial(true)

	request := cms.CreatePutResourceMetricRuleRequest()
	request.RuleName = d.Get("name").(string)
	d.SetId(resource.UniqueId() + ":" + request.RuleName)
	parts, err := ParseResourceId(d.Id(), 2)
	request.RuleId = parts[0]

	request.Namespace = d.Get("project").(string)
	request.MetricName = d.Get("metric").(string)
	request.Period = strconv.Itoa(d.Get("period").(int))
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "cms", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
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
	}
	if len(dimList) > 0 {
		if bytes, err := json.Marshal(dimList); err != nil {
			return fmt.Errorf("marshaling dimensions to json string got an error: %#v", err)
		} else {
			request.Resources = string(bytes[:])
		}
	}

	// make a request
	nrequest := requests.NewCommonRequest()
	nrequest.RegionId = client.RegionId
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	nrequest.Method = "POST"
	nrequest.Product = "cms"
	nrequest.Domain = client.Domain
	nrequest.Version = "2019-01-01"
	nrequest.ApiName = "PutResourceMetricRule"

	nrequest.Headers = map[string]string{"RegionId": client.RegionId}
	nrequest.QueryParams = map[string]string{
		"AccessKeySecret":                client.SecretKey,
		"Product":                        "cms",
		"Department":                     client.Department,
		"ResourceGroup":                  client.ResourceGroup,
		"RegionId":                       client.RegionId,
		"Action":                         "PutResourceMetricRule",
		"Version":                        "2019-01-01",
		"Namespace":                      request.Namespace,
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
		"RuleName":                                request.RuleName,
		"RuleId":                                  request.RuleId,
		"MetricName":                              request.MetricName,
		"Format":                                  "JSON",
		"SilenceTime":                             fmt.Sprint(request.SilenceTime),
		"SignatureVersion":                        "1.0",
		"Period":                                  request.Period,
	}

	raw, err := client.WithEcsClient(func(cmsClient *ecs.Client) (interface{}, error) {
		return cmsClient.ProcessCommonRequest(nrequest)
	})

	log.Printf("testing cms %v", raw)
	if err != nil {
		return fmt.Errorf("Putting alarm got an error: %#v", err)
	}
	//d.SetPartial("name")
	//d.SetPartial("period")
	//d.SetPartial("statistics")
	//d.SetPartial("operator")
	//d.SetPartial("threshold")
	//d.SetPartial("triggered_count")
	//d.SetPartial("contact_groups")
	//d.SetPartial("effective_interval")
	//d.SetPartial("start_time")
	//d.SetPartial("end_time")
	//d.SetPartial("silence_time")
	//d.SetPartial("notify_type")

	if d.Get("enabled").(bool) {
		request := cms.CreateEnableMetricRulesRequest()
		request.RuleId = &[]string{d.Id()}
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "cms", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

		wait := incrementalWait(1*time.Second, 2*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			_, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
				return cmsClient.EnableMetricRules(request)
			})

			if err != nil {
				if IsExpectedErrors(err, []string{ThrottlingUser}) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(fmt.Errorf("Enabling alarm got an error: %#v", err))
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("Enabling alarm got an error: %#v", err)
		}
	} else {
		request := cms.CreateDisableMetricRulesRequest()
		request.RuleId = &[]string{d.Id()}
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "cms", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

		wait := incrementalWait(1*time.Second, 2*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			_, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
				return cmsClient.DisableMetricRules(request)
			})

			if err != nil {
				if IsExpectedErrors(err, []string{ThrottlingUser}) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(fmt.Errorf("Disableing alarm got an error: %#v", err))
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("Disableing alarm got an error: %#v", err)
		}
	}
	if err := cmsService.WaitForCmsAlarm(d.Id(), d.Get("enabled").(bool), 102); err != nil {
		return err
	}

	d.Partial(false)

	return resourceApsaraStackCmsAlarmUpdate(d, meta)
}

func resourceApsaraStackCmsAlarmRead(d *schema.ResourceData, meta interface{}) error {
	wiatSecondsIfWithTest(1)
	client := meta.(*connectivity.ApsaraStackClient)
	cmsService := CmsService{client}

	alarm, err := cmsService.DescribeCmsAlarm(d.Id())

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", alarm.RuleName)
	d.Set("project", alarm.Namespace)
	d.Set("metric", alarm.MetricName)
	if period, err := strconv.Atoi(alarm.Period); err != nil {
		return WrapError(err)
	} else {
		d.Set("period", period)
	}

	escalationsCritical := make([]map[string]interface{}, 1)
	if alarm.Escalations.Critical.Times != "" {
		if count, err := strconv.Atoi(alarm.Escalations.Critical.Times); err != nil {
			return WrapError(err)
		} else {
			mapping := map[string]interface{}{
				"statistics":          alarm.Escalations.Critical.Statistics,
				"comparison_operator": convertOperator(alarm.Escalations.Critical.ComparisonOperator),
				"threshold":           alarm.Escalations.Critical.Threshold,
				"times":               count,
			}
			escalationsCritical[0] = mapping
			d.Set("escalations_critical", escalationsCritical)
		}
	}

	escalationsWarn := make([]map[string]interface{}, 1)
	if alarm.Escalations.Warn.Times != "" {
		if count, err := strconv.Atoi(alarm.Escalations.Warn.Times); err != nil {
			return WrapError(err)
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
	if alarm.Escalations.Info.Times != "" {
		if count, err := strconv.Atoi(alarm.Escalations.Info.Times); err != nil {
			return WrapError(err)
		} else {
			mappingInfo := map[string]interface{}{
				"statistics":          alarm.Escalations.Info.Statistics,
				"comparison_operator": convertOperator(alarm.Escalations.Info.ComparisonOperator),
				"threshold":           alarm.Escalations.Info.Threshold,
				"times":               count,
			}
			escalationsInfo[0] = mappingInfo
			d.Set("escalations_info", escalationsInfo)
		}
	}

	d.Set("effective_interval", alarm.EffectiveInterval)
	if silence, err := strconv.Atoi(alarm.SilenceTime); err != nil {
		return fmt.Errorf("Atoi SilenceTime got an error: %#v.", err)
	} else {
		d.Set("silence_time", silence)
	}
	d.Set("status", alarm.AlertState)
	d.Set("enabled", alarm.EnableState)
	d.Set("contact_groups", strings.Split(alarm.ContactGroups, ","))

	var dims map[string]interface{}
	if alarm.Dimensions != "" {
		if err := json.Unmarshal([]byte(alarm.Dimensions), &dims); err != nil {
			return fmt.Errorf("Unmarshaling Dimensions got an error: %#v.", err)
		}
	}
	d.Set("dimensions", dims)

	return nil
}

func resourceApsaraStackCmsAlarmUpdate(d *schema.ResourceData, meta interface{}) error {

	return resourceApsaraStackCmsAlarmRead(d, meta)
}

func resourceApsaraStackCmsAlarmDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	cmsService := CmsService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request := cms.CreateDeleteMetricRulesRequest()
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "cms", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

	request.Id = &[]string{parts[0]}

	wait := incrementalWait(1*time.Second, 2*time.Second)
	return resource.Retry(10*time.Minute, func() *resource.RetryError {
		_, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
			return cmsClient.DeleteMetricRules(request)
		})

		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(fmt.Errorf("Deleting alarm rule got an error: %#v", err))
		}

		_, err = cmsService.DescribeCmsAlarm(d.Id())
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Describe alarm rule got an error: %#v", err))
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
