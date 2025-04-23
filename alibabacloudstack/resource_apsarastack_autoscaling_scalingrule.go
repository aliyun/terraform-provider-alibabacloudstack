package alibabacloudstack

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackEssScalingRule() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"scaling_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"adjustment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"QuantityChangeInCapacity", "PercentChangeInCapacity", "TotalCapacity"}, false),
			},
			"adjustment_value": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"scaling_rule_name": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 40),
			},
			"ari": {
				Type:         schema.TypeString,
				Computed:     true,
				Deprecated:   "Field 'ari' is deprecated and will be removed in a future release. Please use new field 'scaling_rule_aris' instead.",
			},
			"scaling_rule_aris": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cooldown": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 86400),
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackEssScalingRuleCreate, resourceAlibabacloudStackEssScalingRuleRead, resourceAlibabacloudStackEssScalingRuleUpdate, resourceAlibabacloudStackEssScalingRuleDelete)
	return resource
}

func resourceAlibabacloudStackEssScalingRuleCreate(d *schema.ResourceData, meta interface{}) error {
	request, err := buildAlibabacloudStackEssScalingRuleArgs(d, meta)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	client := meta.(*connectivity.AlibabacloudStackClient)
	client.InitRpcRequest(*request.RpcRequest)

	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.CreateScalingRule(request)
	})
	bresponse, ok := raw.(*ess.CreateScalingRuleResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ess_scalingrule", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	d.SetId(bresponse.ScalingRuleId)

	return nil
}

func resourceAlibabacloudStackEssScalingRuleRead(d *schema.ResourceData, meta interface{}) error {

	//Compatible with older versions id
	if strings.Contains(d.Id(), COLON_SEPARATED) {
		parts, _ := ParseResourceId(d.Id(), 2)
		d.SetId(parts[1])
	}

	client := meta.(*connectivity.AlibabacloudStackClient)
	essService := EssService{client}

	object, err := essService.DescribeEssScalingRule(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("scaling_group_id", object.ScalingGroupId)
	connectivity.SetResourceData(d, object.ScalingRuleAri, "scaling_rule_aris", "ari")
	d.Set("adjustment_type", object.AdjustmentType)
	d.Set("adjustment_value", object.AdjustmentValue)
	d.Set("scaling_rule_name", object.ScalingRuleName)
	d.Set("cooldown", object.Cooldown)

	return nil
}

func resourceAlibabacloudStackEssScalingRuleDelete(d *schema.ResourceData, meta interface{}) error {
	//Compatible with older versions id
	if strings.Contains(d.Id(), COLON_SEPARATED) {
		parts, _ := ParseResourceId(d.Id(), 2)
		d.SetId(parts[1])
	}

	client := meta.(*connectivity.AlibabacloudStackClient)
	essService := EssService{client}
	request := ess.CreateDeleteScalingRuleRequest()
	request.ScalingRuleId = d.Id()
	client.InitRpcRequest(*request.RpcRequest)

	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DeleteScalingRule(request)
	})
	bresponse, ok := raw.(*ess.DeleteScalingRuleResponse)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidScalingRuleId.NotFound"}) {
			return nil
		}
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return errmsgs.WrapError(essService.WaitForEssScalingRule(d.Id(), Deleted, DefaultTimeout))
}

func resourceAlibabacloudStackEssScalingRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	//Compatible with older versions id
	if strings.Contains(d.Id(), COLON_SEPARATED) {
		parts, _ := ParseResourceId(d.Id(), 2)
		d.SetId(parts[1])
	}

	client := meta.(*connectivity.AlibabacloudStackClient)
	request := ess.CreateModifyScalingRuleRequest()
	client.InitRpcRequest(*request.RpcRequest)

	if d.HasChange("scaling_rule_name") {
		request.ScalingRuleName = d.Get("scaling_rule_name").(string)
	}
	if d.HasChange("adjustment_type") {
		request.AdjustmentType = d.Get("adjustment_type").(string)
	}
	if d.HasChange("adjustment_value") {
		request.AdjustmentValue = requests.NewInteger(d.Get("adjustment_value").(int))
	}
	if d.HasChange("cooldown") {
		request.Cooldown = requests.NewInteger(d.Get("cooldown").(int))
	}

	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.ModifyScalingRule(request)
	})
	bresponse, ok := raw.(*ess.ModifyScalingRuleResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}

func buildAlibabacloudStackEssScalingRuleArgs(d *schema.ResourceData, meta interface{}) (*ess.CreateScalingRuleRequest, error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := ess.CreateCreateScalingRuleRequest()
	client.InitRpcRequest(*request.RpcRequest)

	// common params
	request.ScalingGroupId = d.Get("scaling_group_id").(string)

	if v, ok := d.GetOk("scaling_rule_name"); ok && v.(string) != "" {
		request.ScalingRuleName = v.(string)
	}
	if v, ok := d.GetOk("adjustment_type"); ok && v.(string) != "" {
		request.AdjustmentType = v.(string)
	}
	if v, ok := d.GetOkExists("adjustment_value"); ok {
		request.AdjustmentValue = requests.NewInteger(v.(int))
	}
	if v, ok := d.GetOk("cooldown"); ok {
		request.Cooldown = requests.NewInteger(v.(int))
	}

	return request, nil
}