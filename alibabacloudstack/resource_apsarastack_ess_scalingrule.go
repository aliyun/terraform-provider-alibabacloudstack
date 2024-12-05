package alibabacloudstack

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackEssScalingRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackEssScalingRuleCreate,
		Read:   resourceAlibabacloudStackEssScalingRuleRead,
		Update: resourceAlibabacloudStackEssScalingRuleUpdate,
		Delete: resourceAlibabacloudStackEssScalingRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"scaling_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"adjustment_type": {
				Type: schema.TypeString,
				//Optional:     true,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"QuantityChangeInCapacity", "PercentChangeInCapacity", "TotalCapacity"}, false),
			},
			"adjustment_value": {
				Type: schema.TypeInt,
				//Optional: true,
				Required: true,
			},
			"scaling_rule_name": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 40),
			},
			"ari": {
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
}

func resourceAlibabacloudStackEssScalingRuleCreate(d *schema.ResourceData, meta interface{}) error {

	request, err := buildAlibabacloudStackEssScalingRuleArgs(d, meta)
	if err != nil {
		return WrapError(err)
	}

	client := meta.(*connectivity.AlibabacloudStackClient)
	request.RegionId = client.RegionId
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{ "Product": "ess", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.CreateScalingRule(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_ess_scalingrule", request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ess.CreateScalingRuleResponse)
	d.SetId(response.ScalingRuleId)

	return resourceAlibabacloudStackEssScalingRuleRead(d, meta)
}

func resourceAlibabacloudStackEssScalingRuleRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)

	//Compatible with older versions id
	if strings.Contains(d.Id(), COLON_SEPARATED) {
		parts, _ := ParseResourceId(d.Id(), 2)
		d.SetId(parts[1])
	}

	client := meta.(*connectivity.AlibabacloudStackClient)
	essService := EssService{client}

	object, err := essService.DescribeEssScalingRule(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("scaling_group_id", object.ScalingGroupId)
	d.Set("ari", object.ScalingRuleAri)
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
	request.RegionId = client.RegionId
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{ "Product": "ess", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DeleteScalingRule(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidScalingRuleId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return WrapError(essService.WaitForEssScalingRule(d.Id(), Deleted, DefaultTimeout))
}

func resourceAlibabacloudStackEssScalingRuleUpdate(d *schema.ResourceData, meta interface{}) error {

	//Compatible with older versions id
	if strings.Contains(d.Id(), COLON_SEPARATED) {
		parts, _ := ParseResourceId(d.Id(), 2)
		d.SetId(parts[1])
	}

	client := meta.(*connectivity.AlibabacloudStackClient)
	request := ess.CreateModifyScalingRuleRequest()
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ScalingRuleId = d.Id()
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{ "Product": "ess", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

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
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return resourceAlibabacloudStackEssScalingRuleRead(d, meta)
}

func buildAlibabacloudStackEssScalingRuleArgs(d *schema.ResourceData, meta interface{}) (*ess.CreateScalingRuleRequest, error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := ess.CreateCreateScalingRuleRequest()
	request.RegionId = client.RegionId
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{ "Product": "ess", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

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
