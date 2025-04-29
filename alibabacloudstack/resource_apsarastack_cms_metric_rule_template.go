package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudCmsMetricRuleTemplate() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"alert_templates": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"category": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"ecs", "rds", "ads", "slb", "vpc", "apigateway", "cdn", "cs", "dcdn", "ddos", "eip", "elasticsearch", "emr", "ess", "hbase", "iot_edge", "kvstore_sharding", "kvstore_splitrw", "kvstore_standard", "memcache", "mns", "mongodb", "mongodb_cluster", "mongodb_sharding", "mq_topic", "ocs", "opensearch", "oss", "polardb", "petadata", "scdn", "sharebandwidthpackages", "sls", "vpn"}, false),
						},
						"escalations": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"critical": {
										Type:     schema.TypeSet,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"comparison_operator": {
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: validation.StringInSlice([]string{"GreaterThanOrEqualToThreshold", "GreaterThanThreshold", "LessThanOrEqualToThreshold", "LessThanThreshold", "NotEqualToThreshold", "GreaterThanYesterday", "LessThanYesterday", "GreaterThanLastWeek", "LessThanLastWeek", "GreaterThanLastPeriod", "LessThanLastPeriod"}, false),
												},
												"statistics": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"threshold": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"times": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									"info": {
										Type:     schema.TypeSet,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"comparison_operator": {
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: validation.StringInSlice([]string{"GreaterThanOrEqualToThreshold", "GreaterThanThreshold", "LessThanOrEqualToThreshold", "LessThanThreshold", "NotEqualToThreshold", "GreaterThanYesterday", "LessThanYesterday", "GreaterThanLastWeek", "LessThanLastWeek", "GreaterThanLastPeriod", "LessThanLastPeriod"}, false),
												},
												"statistics": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"threshold": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"times": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									"warn": {
										Type:     schema.TypeSet,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"comparison_operator": {
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: validation.StringInSlice([]string{"GreaterThanOrEqualToThreshold", "GreaterThanThreshold", "LessThanOrEqualToThreshold", "LessThanThreshold", "NotEqualToThreshold", "GreaterThanYesterday", "LessThanYesterday", "GreaterThanLastWeek", "LessThanLastWeek", "GreaterThanLastPeriod", "LessThanLastPeriod"}, false),
												},
												"statistics": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"threshold": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"times": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
								},
							},
						},
						"metric_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"namespace": {
							Type:     schema.TypeString,
							Required: true,
						},
						"rule_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"webhook": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"apply_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enable_end_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_start_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"metric_rule_template_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"notify_level": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rest_version": {
				Optional: true,
				Type:     schema.TypeString,
				Computed: true,
			},
			"silence_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 86400),
			},
			"webhook": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"overwrite": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudCmsMetricRuleTemplateCreate,
		resourceAlibabacloudCmsMetricRuleTemplateRead, resourceAlibabacloudCmsMetricRuleTemplateUpdate,
		resourceAlibabacloudCmsMetricRuleTemplateDelete)
	return resource
}

func resourceAlibabacloudCmsMetricRuleTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := cms.CreateCreateMetricRuleTemplateRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.Name = d.Get("metric_rule_template_name").(string)
	request.Description = d.Get("description").(string)

	if v, ok := d.GetOk("alert_templates"); ok {
		alertTemplatesMaps := make([]cms.CreateMetricRuleTemplateAlertTemplates, 0)
		for _, alertTemplates := range v.(*schema.Set).List() {
			alertTemplatesArg := alertTemplates.(map[string]interface{})
			if escalationsMaps, ok := alertTemplatesArg["escalations"]; ok {
				for _, escalationsArg := range escalationsMaps.(*schema.Set).List() {
					alertTemplate := cms.CreateMetricRuleTemplateAlertTemplates{}
					alertTemplate.Category = alertTemplatesArg["category"].(string)
					alertTemplate.MetricName = alertTemplatesArg["metric_name"].(string)
					alertTemplate.Namespace = alertTemplatesArg["namespace"].(string)
					alertTemplate.RuleName = alertTemplatesArg["rule_name"].(string)
					alertTemplate.Webhook = alertTemplatesArg["webhook"].(string)
					if criticalMaps, ok := escalationsArg.(map[string]interface{})["critical"]; ok {
						for _, criticalMap := range criticalMaps.(*schema.Set).List() {
							criticalArg := criticalMap.(map[string]interface{})
							alertTemplate.EscalationsCriticalComparisonOperator = criticalArg["comparison_operator"].(string)
							alertTemplate.EscalationsCriticalStatistics = criticalArg["statistics"].(string)
							alertTemplate.EscalationsCriticalThreshold = criticalArg["threshold"].(string)
							alertTemplate.EscalationsCriticalTimes = criticalArg["times"].(string)
						}
					}
					if infoMaps, ok := escalationsArg.(map[string]interface{})["info"]; ok {
						for _, infoMap := range infoMaps.(*schema.Set).List() {
							infoArg := infoMap.(map[string]interface{})
							alertTemplate.EscalationsInfoComparisonOperator = infoArg["comparison_operator"].(string)
							alertTemplate.EscalationsInfoStatistics = infoArg["statistics"].(string)
							alertTemplate.EscalationsInfoThreshold = infoArg["threshold"].(string)
							alertTemplate.EscalationsInfoTimes = infoArg["times"].(string)
						}
					}
					if warnMaps, ok := escalationsArg.(map[string]interface{})["warn"]; ok {
						for _, warnMap := range warnMaps.(*schema.Set).List() {
							warnArg := warnMap.(map[string]interface{})
							alertTemplate.EscalationsWarnComparisonOperator = warnArg["comparison_operator"].(string)
							alertTemplate.EscalationsWarnStatistics = warnArg["statistics"].(string)
							alertTemplate.EscalationsWarnThreshold = warnArg["threshold"].(string)
							alertTemplate.EscalationsWarnTimes = warnArg["times"].(string)
						}
					}
					alertTemplatesMaps = append(alertTemplatesMaps, alertTemplate)
				}
			}
		}
		request.AlertTemplates = &alertTemplatesMaps
	}
	raw, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
		return cmsClient.CreateMetricRuleTemplate(request)
	})
	addDebug(request.GetActionName(), raw, request, request.QueryParams)
	if err != nil {
		errmsg := ""
		if response, ok := raw.(*responses.BaseResponse); ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "cms_metric_rule_templates", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	response, _ := raw.(*cms.CreateMetricRuleTemplateResponse)
	d.SetId(fmt.Sprint(response.Id))

	return nil
}

func resourceAlibabacloudCmsMetricRuleTemplateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	cmsService := CmsService{client}
	templateattr, err := cmsService.DescribeMetricRuleTemplateAttribute(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloud_cms_metric_rule_template cmsService.DescribeCmsMetricRuleTemplate Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	resource := templateattr.Resource
	if len(resource.AlertTemplates.AlertTemplate) > 0 {
		alertTemplatesMaps := make([]map[string]interface{}, 0)
		for _, alertTemplate := range resource.AlertTemplates.AlertTemplate {
			alertTempArg := make(map[string]interface{}, 0)
			alertTempArg["category"] = alertTemplate.Category
			alertTempArg["metric_name"] = alertTemplate.MetricName
			alertTempArg["namespace"] = alertTemplate.Namespace
			alertTempArg["rule_name"] = alertTemplate.RuleName
			alertTempArg["webhook"] = alertTemplate.Webhook
			escalationsMaps := make([]map[string]interface{}, 0)
			escalationsMap := map[string]interface{}{}

			criticalMap := alertTemplate.Escalations.Critical
			criticalMaps := make([]map[string]interface{}, 0)
			criticalArg := map[string]interface{}{}
			criticalArg["comparison_operator"] = criticalMap.ComparisonOperator
			criticalArg["statistics"] = criticalMap.Statistics
			criticalArg["threshold"] = criticalMap.Threshold
			criticalArg["times"] = criticalMap.Times
			criticalMaps = append(criticalMaps, criticalArg)
			escalationsMap["critical"] = criticalMaps

			infoMap := alertTemplate.Escalations.Info
			infoMaps := make([]map[string]interface{}, 0)
			infoArg := map[string]interface{}{}
			infoArg["comparison_operator"] = infoMap.ComparisonOperator
			infoArg["statistics"] = infoMap.Statistics
			infoArg["threshold"] = infoMap.Threshold
			infoArg["times"] = infoMap.Times
			infoMaps = append(infoMaps, infoArg)
			escalationsMap["info"] = infoMaps

			warnMap := alertTemplate.Escalations.Warn
			warnMaps := make([]map[string]interface{}, 0)
			warnArg := map[string]interface{}{}
			warnArg["comparison_operator"] = warnMap.ComparisonOperator
			warnArg["statistics"] = warnMap.Statistics
			warnArg["threshold"] = warnMap.Threshold
			warnArg["times"] = warnMap.Times
			warnMaps = append(warnMaps, warnArg)
			escalationsMap["warn"] = warnMaps

			escalationsMaps = append(escalationsMaps, escalationsMap)

			alertTempArg["escalations"] = escalationsMaps
			alertTemplatesMaps = append(alertTemplatesMaps, alertTempArg)
			d.Set("alert_templates", alertTemplatesMaps)
		}
	}
	d.Set("description", resource.Description)
	d.Set("metric_rule_template_name", resource.Name)
	d.Set("rest_version", resource.RestVersion)
	return nil
}

func resourceAlibabacloudCmsMetricRuleTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	d.Partial(true)
	update := false

	if v, ok := d.GetOk("enable"); ok && v.(bool) {
		request := client.NewCommonRequest("POST", "cms", "2019-01-01", "ApplyMetricRuleTemplate", "")
		request.QueryParams["TemplateId"] = "[]"
		request.QueryParams["TemplateIds"] = d.Id()
		request.QueryParams["GroupId"] = client.ResourceGroup
		request.QueryParams["Overwrite"] = fmt.Sprintf("%t", d.Get("overwrite").(bool))
		if v, ok := d.GetOk("apply_mode"); ok {
			request.QueryParams["ApplyMode"] = v.(string)
		}
		if v, ok := d.GetOk("enable_end_time"); ok {
			request.QueryParams["EnableEndTime"] = v.(string)
		}
		if v, ok := d.GetOk("enable_start_time"); ok {
			request.QueryParams["EnableStartTime"] = v.(string)
		}
		if v, ok := d.GetOk("notify_level"); ok {
			request.QueryParams["NotifyLevel"] = v.(string)
		}
		if v, ok := d.GetOk("silence_time"); ok {
			request.QueryParams["SilenceTime"] = fmt.Sprint(v)
		}
		if v, ok := d.GetOk("webhook"); ok {
			request.QueryParams["Webhook"] = v.(string)
		}

		raw, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
			return cmsClient.ProcessCommonRequest(request)
		})
		addDebug(request.GetActionName(), raw, request, request.QueryParams)
		if err != nil {
			errmsg := ""
			if response, ok := raw.(*responses.BaseResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "ApplyMetricRuleTemplate", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		bresponse, _ := raw.(*responses.CommonResponse)
		resource := make(map[string]interface{})
		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &resource)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DataDefaultErrorMsg, "ApplyMetricRuleTemplate", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
		}
		if resource["Code"].(float64) != 200 {
			return errmsgs.WrapError(fmt.Errorf("ApplyMetricRuleTemplate Error: %v", resource))
		}
		d.Set("group_id", client.ResourceGroup)
	}

	update = false
	modifyMetricRuleTemplateReq := cms.CreateModifyMetricRuleTemplateRequest()
	client.InitRpcRequest(*modifyMetricRuleTemplateReq.RpcRequest)

	if v, ok := d.GetOk("rest_version"); ok {
		rest_version, _ := strconv.Atoi(v.(string))
		modifyMetricRuleTemplateReq.RestVersion = requests.NewInteger(rest_version)
	}
	if !d.IsNewResource() && d.HasChange("alert_templates") {
		update = true
		if v, ok := d.GetOk("alert_templates"); ok {
			alertTemplatesMaps := make([]cms.ModifyMetricRuleTemplateAlertTemplates, 0)
			for _, alertTemplates := range v.(*schema.Set).List() {
				alertTemplatesArg := alertTemplates.(map[string]interface{})
				if escalationsMaps, ok := alertTemplatesArg["escalations"]; ok {
					for _, escalationsArg := range escalationsMaps.(*schema.Set).List() {
						alertTemplate := cms.ModifyMetricRuleTemplateAlertTemplates{}
						alertTemplate.Category = alertTemplatesArg["category"].(string)
						alertTemplate.MetricName = alertTemplatesArg["metric_name"].(string)
						alertTemplate.Namespace = alertTemplatesArg["namespace"].(string)
						alertTemplate.RuleName = alertTemplatesArg["rule_name"].(string)
						alertTemplate.Webhook = alertTemplatesArg["webhook"].(string)
						if criticalMaps, ok := escalationsArg.(map[string]interface{})["critical"]; ok {
							for _, criticalMap := range criticalMaps.(*schema.Set).List() {
								criticalArg := criticalMap.(map[string]interface{})
								alertTemplate.EscalationsCriticalComparisonOperator = criticalArg["comparison_operator"].(string)
								alertTemplate.EscalationsCriticalStatistics = criticalArg["statistics"].(string)
								alertTemplate.EscalationsCriticalThreshold = criticalArg["threshold"].(string)
								alertTemplate.EscalationsCriticalTimes = criticalArg["times"].(string)
							}
						}
						if infoMaps, ok := escalationsArg.(map[string]interface{})["info"]; ok {
							for _, infoMap := range infoMaps.(*schema.Set).List() {
								infoArg := infoMap.(map[string]interface{})
								alertTemplate.EscalationsInfoComparisonOperator = infoArg["comparison_operator"].(string)
								alertTemplate.EscalationsInfoStatistics = infoArg["statistics"].(string)
								alertTemplate.EscalationsInfoThreshold = infoArg["threshold"].(string)
								alertTemplate.EscalationsInfoTimes = infoArg["times"].(string)
							}
						}
						if warnMaps, ok := escalationsArg.(map[string]interface{})["warn"]; ok {
							for _, warnMap := range warnMaps.(*schema.Set).List() {
								warnArg := warnMap.(map[string]interface{})
								alertTemplate.EscalationsWarnComparisonOperator = warnArg["comparison_operator"].(string)
								alertTemplate.EscalationsWarnStatistics = warnArg["statistics"].(string)
								alertTemplate.EscalationsWarnThreshold = warnArg["threshold"].(string)
								alertTemplate.EscalationsWarnTimes = warnArg["times"].(string)
							}
						}
						alertTemplatesMaps = append(alertTemplatesMaps, alertTemplate)
					}
				}
			}
			modifyMetricRuleTemplateReq.AlertTemplates = &alertTemplatesMaps
		}
	}
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			modifyMetricRuleTemplateReq.Description = v.(string)
		}
	}
	if !d.IsNewResource() && d.HasChange("metric_rule_template_name") {
		update = true
		modifyMetricRuleTemplateReq.Name = d.Get("metric_rule_template_name").(string)
	}
	if update {
		raw, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
			return cmsClient.ModifyMetricRuleTemplate(modifyMetricRuleTemplateReq)
		})
		addDebug(modifyMetricRuleTemplateReq.GetActionName(), raw, modifyMetricRuleTemplateReq, modifyMetricRuleTemplateReq.QueryParams)
		if err != nil {
			errmsg := ""
			if response, ok := raw.(*responses.BaseResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "ApplyMetricRuleTemplate", modifyMetricRuleTemplateReq.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		response, _ := raw.(*cms.ModifyMetricRuleTemplateResponse)
		if response.Code != 200 {
			return errmsgs.WrapError(fmt.Errorf("%s", response.Message))
		}
	}
	d.Partial(false)
	return nil
}

func resourceAlibabacloudCmsMetricRuleTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := cms.CreateDeleteMetricRuleTemplateRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.TemplateId = d.Id()

	raw, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
		return cmsClient.DeleteMetricRuleTemplate(request)
	})
	addDebug(request.GetActionName(), raw, request, request.QueryParams)
	if err != nil {
		errmsg := ""
		if response, ok := raw.(*responses.BaseResponse); ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "DeleteMetricRuleTemplate", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	response, _ := raw.(*cms.DeleteMetricRuleTemplateResponse)
	if response.Code != 200 {
		return errmsgs.WrapError(fmt.Errorf("%s", response.Message))
	}
	return nil
}
