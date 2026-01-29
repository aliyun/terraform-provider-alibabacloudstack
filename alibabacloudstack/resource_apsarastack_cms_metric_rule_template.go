package alibabacloudstack

import (
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

func resourceAlibabacloudStackCmsMetricRuleTemplate() *schema.Resource {
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
													Type:     schema.TypeInt,
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
													Type:     schema.TypeInt,
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
													Type:     schema.TypeInt,
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
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"metric_rule_template_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackCmsMetricRuleTemplateCreate,
		resourceAlibabacloudStackCmsMetricRuleTemplateRead, resourceAlibabacloudStackCmsMetricRuleTemplateUpdate,
		resourceAlibabacloudStackCmsMetricRuleTemplateDelete)
	return resource
}

func resourceAlibabacloudStackCmsMetricRuleTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	reqQuery := map[string]interface{}{
		"Name": d.Get("metric_rule_template_name"),
	}

	if v, ok := d.GetOk("description"); ok {
		reqQuery["Description"] = v
	}

	if v, ok := d.GetOk("alert_templates"); ok {
		alertTemplatesList := v.(*schema.Set).List()
		if len(alertTemplatesList) > 0 {
			alertTemplates := make([]interface{}, 0, len(alertTemplatesList))

			for _, alertTemplate := range alertTemplatesList {
				tmpl := alertTemplate.(map[string]interface{})
				templateMap := map[string]interface{}{
					"Category":   tmpl["category"],
					"MetricName": tmpl["metric_name"],
					"Namespace":  tmpl["namespace"],
					"RuleName":   tmpl["rule_name"],
					"Webhook":    tmpl["webhook"],
				}

				// Handle escalations
				if escRaw, ok := tmpl["escalations"]; ok {
					escList := escRaw.(*schema.Set).List()
					if len(escList) > 0 {
						escalation := escList[0].(map[string]interface{})
						escalationsMap := map[string]interface{}{}

						// Helper to extract level config
						extractLevel := func(levelRaw interface{}) map[string]interface{} {
							if levelRaw == nil {
								return nil
							}
							levelList := levelRaw.(*schema.Set).List()
							if len(levelList) == 0 {
								return nil
							}
							lm := levelList[0].(map[string]interface{})
							return map[string]interface{}{
								"ComparisonOperator": lm["comparison_operator"],
								"Statistics":         lm["statistics"],
								"Threshold":          lm["threshold"],
								"Times":              lm["times"],
							}
						}

						if critical := extractLevel(escalation["critical"]); critical != nil {
							escalationsMap["Critical"] = critical
						}
						if info := extractLevel(escalation["info"]); info != nil {
							escalationsMap["Info"] = info
						}
						if warn := extractLevel(escalation["warn"]); warn != nil {
							escalationsMap["Warn"] = warn
						}

						if len(escalationsMap) > 0 {
							templateMap["Escalations"] = escalationsMap
						}
					}
				}

				alertTemplates = append(alertTemplates, templateMap)
			}

			reqQuery["AlertTemplates"] = alertTemplates
		}
	}

	response, err := client.DoTeaRequest("POST", "Cms", "2019-01-01", "CreateMetricRuleTemplate", "", nil, reqQuery, nil)

	if err != nil {
		return err
	}

	if code, ok := response["Code"]; ok {
		if codeFloat, ok := code.(float64); ok && codeFloat != 200 {
			return errmsgs.WrapError(fmt.Errorf("CreateMetricRuleTemplate failed: %v", response["Message"]))
		} else if codeInt, ok := code.(int); ok && codeInt != 200 {
			return errmsgs.WrapError(fmt.Errorf("CreateMetricRuleTemplate failed: %v", response["Message"]))
		}
	} else {
		if _, ok := response["Id"]; !ok {
			return errmsgs.WrapError(fmt.Errorf("CreateMetricRuleTemplate unexpected response: %v", response))
		}
	}

	if id, exists := response["Id"]; exists {
		d.SetId(fmt.Sprintf("%v", id))
	} else {
		return errmsgs.WrapError(fmt.Errorf("CreateMetricRuleTemplate response missing Id: %v", response))
	}

	return nil
}

func resourceAlibabacloudStackCmsMetricRuleTemplateRead(d *schema.ResourceData, meta interface{}) error {
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
			if criticalMap.ComparisonOperator != "" {
				criticalMaps := make([]map[string]interface{}, 0)
				criticalArg := map[string]interface{}{}
				criticalArg["comparison_operator"] = criticalMap.ComparisonOperator
				criticalArg["statistics"] = criticalMap.Statistics
				criticalArg["threshold"] = criticalMap.Threshold
				criticalArg["times"] = criticalMap.Times
				criticalMaps = append(criticalMaps, criticalArg)
				escalationsMap["critical"] = criticalMaps
			}

			infoMap := alertTemplate.Escalations.Info
			if infoMap.ComparisonOperator != "" {
				infoMaps := make([]map[string]interface{}, 0)
				infoArg := map[string]interface{}{}
				infoArg["comparison_operator"] = infoMap.ComparisonOperator
				infoArg["statistics"] = infoMap.Statistics
				infoArg["threshold"] = infoMap.Threshold
				infoArg["times"] = infoMap.Times
				infoMaps = append(infoMaps, infoArg)
				escalationsMap["info"] = infoMaps
			}

			warnMap := alertTemplate.Escalations.Warn
			if warnMap.ComparisonOperator != "" {
				warnMaps := make([]map[string]interface{}, 0)
				warnArg := map[string]interface{}{}
				warnArg["comparison_operator"] = warnMap.ComparisonOperator
				warnArg["statistics"] = warnMap.Statistics
				warnArg["threshold"] = warnMap.Threshold
				warnArg["times"] = warnMap.Times
				warnMaps = append(warnMaps, warnArg)
				escalationsMap["warn"] = warnMaps
			}

			escalationsMaps = append(escalationsMaps, escalationsMap)

			alertTempArg["escalations"] = escalationsMaps
			alertTemplatesMaps = append(alertTemplatesMaps, alertTempArg)
			d.Set("alert_templates", alertTemplatesMaps)
		}
	}
	d.Set("description", resource.Description)
	d.Set("metric_rule_template_name", resource.Name)

	return nil
}

func resourceAlibabacloudStackCmsMetricRuleTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	if d.IsNewResource() {
		return nil
	}

	if d.HasChanges("alert_templates", "description", "metric_rule_template_name") {
		modifyMetricRuleTemplateReq := cms.CreateModifyMetricRuleTemplateRequest()
		client.InitRpcRequest(*modifyMetricRuleTemplateReq.RpcRequest)

		cmsService := CmsService{client}
		templateattr, err := cmsService.DescribeMetricRuleTemplateAttribute(d.Id())
		if err != nil {
			return errmsgs.WrapError(err)
		}
		modifyMetricRuleTemplateReq.RestVersion = requests.NewInteger(templateattr.Resource.RestVersion)
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
								alertTemplate.EscalationsCriticalTimes = strconv.Itoa(criticalArg["times"].(int))
							}
						}
						if infoMaps, ok := escalationsArg.(map[string]interface{})["info"]; ok {
							for _, infoMap := range infoMaps.(*schema.Set).List() {
								infoArg := infoMap.(map[string]interface{})
								alertTemplate.EscalationsInfoComparisonOperator = infoArg["comparison_operator"].(string)
								alertTemplate.EscalationsInfoStatistics = infoArg["statistics"].(string)
								alertTemplate.EscalationsInfoThreshold = infoArg["threshold"].(string)
								alertTemplate.EscalationsInfoTimes = strconv.Itoa(infoArg["times"].(int))
							}
						}
						if warnMaps, ok := escalationsArg.(map[string]interface{})["warn"]; ok {
							for _, warnMap := range warnMaps.(*schema.Set).List() {
								warnArg := warnMap.(map[string]interface{})
								alertTemplate.EscalationsWarnComparisonOperator = warnArg["comparison_operator"].(string)
								alertTemplate.EscalationsWarnStatistics = warnArg["statistics"].(string)
								alertTemplate.EscalationsWarnThreshold = warnArg["threshold"].(string)
								alertTemplate.EscalationsWarnTimes = strconv.Itoa(warnArg["times"].(int))
							}
						}
						alertTemplatesMaps = append(alertTemplatesMaps, alertTemplate)
					}
				}
			}
			modifyMetricRuleTemplateReq.AlertTemplates = &alertTemplatesMaps
		}
		if v, ok := d.GetOk("description"); ok {
			modifyMetricRuleTemplateReq.Description = v.(string)
		}
		modifyMetricRuleTemplateReq.Name = d.Get("metric_rule_template_name").(string)
		id, err := strconv.Atoi(d.Id())
		if err != nil {
			return err
		}
		modifyMetricRuleTemplateReq.TemplateId = requests.NewInteger(id)
		raw, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
			return cmsClient.ModifyMetricRuleTemplate(modifyMetricRuleTemplateReq)
		})
		addDebug(modifyMetricRuleTemplateReq.GetActionName(), raw, modifyMetricRuleTemplateReq, modifyMetricRuleTemplateReq.QueryParams)
		if err != nil {
			errmsg := ""
			if response, ok := raw.(*responses.BaseResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "ModifyMetricRuleTemplate", modifyMetricRuleTemplateReq.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		response, _ := raw.(*cms.ModifyMetricRuleTemplateResponse)
		if response.Code != 200 {
			return errmsgs.WrapError(fmt.Errorf("%s", response.Message))
		}
	}
	return nil
}

func resourceAlibabacloudStackCmsMetricRuleTemplateDelete(d *schema.ResourceData, meta interface{}) error {
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
