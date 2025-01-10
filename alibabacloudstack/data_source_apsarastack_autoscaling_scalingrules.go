package alibabacloudstack

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackEssScalingRules() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceAlibabacloudStackEssScalingRulesRead,
		Schema: map[string]*schema.Schema{
			"scaling_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				ForceNew:     true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scaling_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cooldown": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"adjustment_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"adjustment_value": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"scaling_rule_ari": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackEssScalingRulesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := ess.CreateDescribeScalingRulesRequest()
	client.InitRpcRequest(*request.RpcRequest)

	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)

	if scalingGroupId, ok := d.GetOk("scaling_group_id"); ok && scalingGroupId.(string) != "" {
		request.ScalingGroupId = scalingGroupId.(string)
	}

	if ruleType, ok := d.GetOk("type"); ok && ruleType.(string) != "" {
		request.ScalingRuleType = ruleType.(string)
	}

	var allScalingRules []ess.ScalingRule

	for {
		raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.DescribeScalingRules(request)
		})
		response, ok := raw.(*ess.DescribeScalingRulesResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ess_scalingrules", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		if len(response.ScalingRules.ScalingRule) < 1 {
			break
		}

		allScalingRules = append(allScalingRules, response.ScalingRules.ScalingRule...)

		if len(response.ScalingRules.ScalingRule) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return err
		} else {
			request.PageNumber = page
		}
	}

	var filteredScalingRulesTemp = make([]ess.ScalingRule, 0)

	nameRegex, okNameRegex := d.GetOk("name_regex")
	idsMap := make(map[string]string)
	ids, okIds := d.GetOk("ids")
	if okIds {
		for _, i := range ids.([]interface{}) {
			idsMap[i.(string)] = i.(string)
		}
	}

	if okNameRegex || okIds {
		for _, rule := range allScalingRules {
			if okNameRegex && nameRegex != "" {
				var r = regexp.MustCompile(nameRegex.(string))
				if r != nil && !r.MatchString(rule.ScalingRuleName) {
					continue
				}
			}
			if okIds && len(idsMap) > 0 {
				if _, ok := idsMap[rule.ScalingRuleId]; !ok {
					continue
				}
			}
			filteredScalingRulesTemp = append(filteredScalingRulesTemp, rule)
		}
	} else {
		filteredScalingRulesTemp = allScalingRules
	}
	return scalingRulesDescriptionAttribute(d, filteredScalingRulesTemp, meta)
}

func scalingRulesDescriptionAttribute(d *schema.ResourceData, scalingRules []ess.ScalingRule, meta interface{}) error {
	var ids []string
	var names []string
	var s = make([]map[string]interface{}, 0)
	for _, scalingRule := range scalingRules {
		mapping := map[string]interface{}{
			"id":                  scalingRule.ScalingRuleId,
			"scaling_group_id":    scalingRule.ScalingGroupId,
			"name":                scalingRule.ScalingRuleName,
			"type":                scalingRule.ScalingRuleType,
			"cooldown":            scalingRule.Cooldown,
			"adjustment_type":     scalingRule.AdjustmentType,
			"adjustment_value":    scalingRule.AdjustmentValue,
			"scaling_rule_ari":    scalingRule.ScalingRuleAri,
		}
		ids = append(ids, scalingRule.ScalingRuleId)
		names = append(names, scalingRule.ScalingRuleName)
		s = append(s, mapping)
	}
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("rules", s); err != nil {
		return errmsgs.WrapError(err)
	}

	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return errmsgs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), s); err != nil {
			return err
		}
	}
	return nil
}
