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
	setResourceFunc(resource, resourceAlibabacloudStackCmsMetricRuleTemplateCreate,
		resourceAlibabacloudStackCmsMetricRuleTemplateRead, resourceAlibabacloudStackCmsMetricRuleTemplateUpdate,
		resourceAlibabacloudStackCmsMetricRuleTemplateDelete)
	return resource
}

func resourceAlibabacloudStackCmsMetricRuleTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := client.NewCommonRequest("POST", "Cms", "2019-01-01", "CreateMetricRuleTemplate", "")

	request.QueryParams["Name"] = d.Get("metric_rule_template_name").(string)
	if v, ok := d.GetOk("description"); ok {
		request.QueryParams["Description"] = v.(string)
	}

	if v, ok := d.GetOk("alert_templates"); ok {
		alertTemplatesList := v.(*schema.Set).List()
		if len(alertTemplatesList) > 0 {

			for i, alertTemplate := range alertTemplatesList {
				alertTemplateMap := alertTemplate.(map[string]interface{})

				request.QueryParams[fmt.Sprintf("AlertTemplates.AlertTemplate.%d.Category", i+1)] = alertTemplateMap["category"].(string)
				request.QueryParams[fmt.Sprintf("AlertTemplates.AlertTemplate.%d.MetricName", i+1)] = alertTemplateMap["metric_name"].(string)
				request.QueryParams[fmt.Sprintf("AlertTemplates.AlertTemplate.%d.Namespace", i+1)] = alertTemplateMap["namespace"].(string)
				request.QueryParams[fmt.Sprintf("AlertTemplates.AlertTemplate.%d.RuleName", i+1)] = alertTemplateMap["rule_name"].(string)
				request.QueryParams[fmt.Sprintf("AlertTemplates.AlertTemplate.%d.Webhook", i+1)] = alertTemplateMap["webhook"].(string)

				if escalations, ok := alertTemplateMap["escalations"]; ok {
					escalationsList := escalations.(*schema.Set).List()
					if len(escalationsList) > 0 {
						escalationMap := escalationsList[0].(map[string]interface{}) // MaxItems 为 1，所以只取第一个

						if critical, ok := escalationMap["critical"]; ok {
							criticalList := critical.(*schema.Set).List()
							if len(criticalList) > 0 {
								criticalMap := criticalList[0].(map[string]interface{}) // MaxItems 为 1
								request.QueryParams[fmt.Sprintf("AlertTemplates.AlertTemplate.%d.Escalations.Critical.%d.ComparisonOperator", i+1, 1)] = criticalMap["comparison_operator"].(string)
								request.QueryParams[fmt.Sprintf("AlertTemplates.AlertTemplate.%d.Escalations.Critical.%d.Statistics", i+1, 1)] = criticalMap["statistics"].(string)
								request.QueryParams[fmt.Sprintf("AlertTemplates.AlertTemplate.%d.Escalations.Critical.%d.Threshold", i+1, 1)] = criticalMap["threshold"].(string)
								request.QueryParams[fmt.Sprintf("AlertTemplates.AlertTemplate.%d.Escalations.Critical.%d.Times", i+1, 1)] = criticalMap["times"].(string)
							}
						}

						if info, ok := escalationMap["info"]; ok {
							infoList := info.(*schema.Set).List()
							if len(infoList) > 0 {
								infoMap := infoList[0].(map[string]interface{}) // MaxItems 为 1
								request.QueryParams[fmt.Sprintf("AlertTemplates.AlertTemplate.%d.Escalations.Info.%d.ComparisonOperator", i+1, 1)] = infoMap["comparison_operator"].(string)
								request.QueryParams[fmt.Sprintf("AlertTemplates.AlertTemplate.%d.Escalations.Info.%d.Statistics", i+1, 1)] = infoMap["statistics"].(string)
								request.QueryParams[fmt.Sprintf("AlertTemplates.AlertTemplate.%d.Escalations.Info.%d.Threshold", i+1, 1)] = infoMap["threshold"].(string)
								request.QueryParams[fmt.Sprintf("AlertTemplates.AlertTemplate.%d.Escalations.Info.%d.Times", i+1, 1)] = infoMap["times"].(string)
							}
						}

						if warn, ok := escalationMap["warn"]; ok {
							warnList := warn.(*schema.Set).List()
							if len(warnList) > 0 {
								warnMap := warnList[0].(map[string]interface{}) // MaxItems 为 1
								request.QueryParams[fmt.Sprintf("AlertTemplates.AlertTemplate.%d.Escalations.Warn.%d.ComparisonOperator", i+1, 1)] = warnMap["comparison_operator"].(string)
								request.QueryParams[fmt.Sprintf("AlertTemplates.AlertTemplate.%d.Escalations.Warn.%d.Statistics", i+1, 1)] = warnMap["statistics"].(string)
								request.QueryParams[fmt.Sprintf("AlertTemplates.AlertTemplate.%d.Escalations.Warn.%d.Threshold", i+1, 1)] = warnMap["threshold"].(string)
								request.QueryParams[fmt.Sprintf("AlertTemplates.AlertTemplate.%d.Escalations.Warn.%d.Times", i+1, 1)] = warnMap["times"].(string)
							}
						}
					}
				}
			}
		}
	}

	bresponse, err := client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	log.Printf(" response of raw CreateMetricRuleTemplate : %s", bresponse)

	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "cms_metric_rule_templates", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	var response map[string]interface{}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DataDefaultErrorMsg, "CreateMetricRuleTemplate", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
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

func resourceAlibabacloudStackCmsMetricRuleTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	d.Partial(true)
	update := false

	if v, ok := d.GetOk("enable"); ok && v.(bool) {
		request := client.NewCommonRequest("POST", "Cms", "2019-01-01", "ApplyMetricRuleTemplate", "")
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

		bresponse, err := client.ProcessCommonRequest(request)
		addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
		log.Printf(" response of raw ApplyMetricRuleTemplate : %s", bresponse)
		if err != nil {
			if bresponse == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_cms_metric_rule_template", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
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
