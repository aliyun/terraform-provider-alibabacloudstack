package alibabacloudstack

import (
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackEssAlarm() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackEssAlarmCreate,
		Read:   resourceAlibabacloudStackEssAlarmRead,
		Update: resourceAlibabacloudStackEssAlarmUpdate,
		Delete: resourceAlibabacloudStackEssAlarmDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Deprecated:   "Field 'name' is deprecated and will be removed in a future release. Please use new field 'alarm_task_name' instead.",
				ConflictsWith: []string{"alarm_task_name"},
			},
			"alarm_task_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ConflictsWith: []string{"name"},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable": {
				Type:         schema.TypeBool,
				Optional:     true,
				Default:      true,
				Deprecated:   "Field 'enable' is deprecated and will be removed in a future release. Please use new field 'status' instead.",
				ConflictsWith: []string{"status"},
			},
			"status": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ConflictsWith: []string{"enable"},
			},
			"alarm_actions": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
				MaxItems: 5,
				MinItems: 1,
			},
			"scaling_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"metric_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "system",
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"system", "custom"}, false),
			},
			"metric_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"period": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      300,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{60, 120, 300, 900}),
			},
			"statistics": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  Average,
				ValidateFunc: validation.StringInSlice([]string{
					string(Average),
					string(Minimum),
					string(Maximum),
				}, false),
			},
			"threshold": {
				Type:     schema.TypeString,
				Required: true,
			},
			"comparison_operator": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      ">=",
				ValidateFunc: validation.StringInSlice([]string{">", ">=", "<", "<="}, false),
			},
			"evaluation_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      3,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"cloud_monitor_group_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"dimensions": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},
			"state": {
				Type:         schema.TypeString,
				Computed:     true,
				Deprecated:   "Field 'state' is deprecated and will be removed in a future release. Please use new field 'alarm_trigger_state' instead.",
				ConflictsWith: []string{"alarm_trigger_state"},
			},
			"alarm_trigger_state": {
				Type:     schema.TypeString,
				Computed: true,
				ConflictsWith: []string{"state"},
			},
		},
	}
}

func resourceAlibabacloudStackEssAlarmCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request, err := buildAlibabacloudStackEssAlarmArgs(d)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	client.InitRpcRequest(*request.RpcRequest)

	var raw interface{}
	if err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.CreateAlarm(request)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{errmsgs.Throttling}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	}); err != nil {
		errmsg := ""
		if raw != nil {
			response, ok := raw.(*ess.CreateAlarmResponse)
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ess_alarm", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	response, _ := raw.(*ess.CreateAlarmResponse)
	d.SetId(response.AlarmTaskId)

	// enable or disable alarm
	enable := connectivity.GetResourceData(d, "status", "enable")
	if !enable.(bool) {
		disableAlarmRequest := ess.CreateDisableAlarmRequest()
		client.InitRpcRequest(*disableAlarmRequest.RpcRequest)
		disableAlarmRequest.AlarmTaskId = response.AlarmTaskId
		raw, err = client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.DisableAlarm(disableAlarmRequest)
		})
		if err != nil {
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*ess.DisableAlarmResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), disableAlarmRequest.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(disableAlarmRequest.GetActionName(), raw, disableAlarmRequest.RpcRequest, disableAlarmRequest)
	}
	return resourceAlibabacloudStackEssAlarmRead(d, meta)
}

func resourceAlibabacloudStackEssAlarmRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	essService := EssService{client}

	object, err := essService.DescribeEssAlarm(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	connectivity.SetResourceData(d, object.Name, "alarm_task_name", "name")
	d.Set("description", object.Description)
	d.Set("alarm_actions", object.AlarmActions.AlarmAction)
	d.Set("scaling_group_id", object.ScalingGroupId)
	d.Set("metric_type", object.MetricType)
	d.Set("metric_name", object.MetricName)
	d.Set("period", object.Period)
	d.Set("statistics", object.Statistics)
	d.Set("threshold", strconv.FormatFloat(object.Threshold, 'f', -1, 32))
	d.Set("comparison_operator", object.ComparisonOperator)
	d.Set("evaluation_count", object.EvaluationCount)
	connectivity.SetResourceData(d, object.Enable, "status", "enable")
	connectivity.SetResourceData(d, object.State, "alarm_trigger_state", "state")

	dims := make([]ess.Dimension, 0, len(object.Dimensions.Dimension))
	for _, dimension := range object.Dimensions.Dimension {
		if dimension.DimensionKey == GroupId {
			d.Set("cloud_monitor_group_id", dimension.DimensionValue)
		} else {
			dims = append(dims, dimension)
		}
	}

	if err := d.Set("dimensions", essService.flattenDimensionsToMap(dims)); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}

func resourceAlibabacloudStackEssAlarmUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := ess.CreateModifyAlarmRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.AlarmTaskId = d.Id()

	d.Partial(true)
	if metricType, ok := d.GetOk("metric_type"); ok && metricType.(string) != "" {
		request.MetricType = metricType.(string)
	}
	if d.HasChanges("alarm_task_name","name")  {
		request.Name = connectivity.GetResourceData(d, "alarm_task_name", "name").(string)
	}

	if d.HasChange("description") {
		request.Description = d.Get("description").(string)
	}

	if d.HasChange("alarm_actions") {
		if v, ok := d.GetOk("alarm_actions"); ok {
			alarmActions := expandStringList(v.(*schema.Set).List())
			if len(alarmActions) > 0 {
				request.AlarmAction = &alarmActions
			}
		}
	}
	if d.HasChange("metric_name") {
		request.MetricName = d.Get("metric_name").(string)
	}
	if d.HasChange("statistics") {
		request.Statistics = d.Get("statistics").(string)
	}
	if d.HasChange("threshold") {
		request.Threshold = requests.Float(d.Get("threshold").(string))
	}
	if d.HasChange("comparison_operator") {
		request.ComparisonOperator = d.Get("comparison_operator").(string)
	}
	if d.HasChange("evaluation_count") {
		request.EvaluationCount = requests.NewInteger(d.Get("evaluation_count").(int))
	}
	if v, ok := d.GetOk("cloud_monitor_group_id"); ok {
		request.GroupId = requests.NewInteger(v.(int))
	}

	dimensions := d.Get("dimensions").(map[string]interface{})
	createAlarmDimensions := make([]ess.ModifyAlarmDimension, 0, len(dimensions))
	for k, v := range dimensions {
		if k == UserId || k == ScalingGroup {
			return errmsgs.WrapError(errmsgs.Error("Invalide dimension keys, %s", k))
		}
		if k != "" {
			dimension := ess.ModifyAlarmDimension{
				DimensionKey:   k,
				DimensionValue: v.(string),
			}
			createAlarmDimensions = append(createAlarmDimensions, dimension)
		}
	}
	request.Dimension = &createAlarmDimensions

	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.ModifyAlarm(request)
	})
	if err != nil {
		errmsg := ""
		if raw != nil {
			response, ok := raw.(*ess.ModifyAlarmResponse)
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if d.HasChanges("status", "enable"){
		enable := connectivity.GetResourceData(d, "status", "enable")
		if enable.(bool) {
			enableAlarmRequest := ess.CreateEnableAlarmRequest()
			client.InitRpcRequest(*enableAlarmRequest.RpcRequest)
			enableAlarmRequest.AlarmTaskId = d.Id()
			raw, err = client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
				return essClient.EnableAlarm(enableAlarmRequest)
			})
			if err != nil {
				errmsg := ""
				if raw != nil {
					response, ok := raw.(*ess.EnableAlarmResponse)
					if ok {
						errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
					}
				}
				return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), enableAlarmRequest.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
			}
			addDebug(enableAlarmRequest.GetActionName(), raw)
		} else {
			disableAlarmRequest := ess.CreateDisableAlarmRequest()
			client.InitRpcRequest(*disableAlarmRequest.RpcRequest)
			disableAlarmRequest.AlarmTaskId = d.Id()
			raw, err = client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
				return essClient.DisableAlarm(disableAlarmRequest)
			})
			if err != nil {
				errmsg := ""
				if raw != nil {
					response, ok := raw.(*ess.DisableAlarmResponse)
					if ok {
						errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
					}
				}
				return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), disableAlarmRequest.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
			}
			addDebug(disableAlarmRequest.GetActionName(), raw)
		}
	}
	d.Partial(false)
	return resourceAlibabacloudStackEssAlarmRead(d, meta)
}

func resourceAlibabacloudStackEssAlarmDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	essService := EssService{client}

	request := ess.CreateDeleteAlarmRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.AlarmTaskId = d.Id()

	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DeleteAlarm(request)
	})
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"404"}) {
			return nil
		}
		errmsg := ""
		if raw != nil {
			response, ok := raw.(*ess.DeleteAlarmResponse)
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return errmsgs.WrapError(essService.WaitForEssAlarm(d.Id(), Deleted, DefaultTimeout))
}

func buildAlibabacloudStackEssAlarmArgs(d *schema.ResourceData) (*ess.CreateAlarmRequest, error) {
	request := ess.CreateCreateAlarmRequest()

	if name, ok := connectivity.GetResourceDataOk(d, "alarm_task_name", "name"); ok && name.(string) != "" {
		request.Name = name.(string)
	}

	if description, ok := d.GetOk("description"); ok && description.(string) != "" {
		request.Description = description.(string)
	}

	if v, ok := d.GetOk("alarm_actions"); ok {
		alarmActions := expandStringList(v.(*schema.Set).List())
		request.AlarmAction = &alarmActions
	}

	if scalingGroupId := d.Get("scaling_group_id").(string); scalingGroupId != "" {
		request.ScalingGroupId = scalingGroupId
	}

	if metricType, ok := d.GetOk("metric_type"); ok && metricType.(string) != "" {
		request.MetricType = metricType.(string)
	}

	if metricName := d.Get("metric_name").(string); metricName != "" {
		request.MetricName = metricName
	}

	if period, ok := d.GetOk("period"); ok && period.(int) > 0 {
		request.Period = requests.NewInteger(period.(int))
	}

	if statistics, ok := d.GetOk("statistics"); ok && statistics.(string) != "" {
		request.Statistics = statistics.(string)
	}

	if v, ok := d.GetOk("threshold"); ok {
		threshold, err := strconv.ParseFloat(v.(string), 32)
		if err != nil {
			return nil, errmsgs.WrapError(err)
		}
		request.Threshold = requests.NewFloat(threshold)
	}

	if comparisonOperator, ok := d.GetOk("comparison_operator"); ok && comparisonOperator.(string) != "" {
		request.ComparisonOperator = comparisonOperator.(string)
	}

	if evaluationCount, ok := d.GetOk("evaluation_count"); ok && evaluationCount.(int) > 0 {
		request.EvaluationCount = requests.NewInteger(evaluationCount.(int))
	}

	if groupId, ok := d.GetOk("cloud_monitor_group_id"); ok {
		request.GroupId = requests.NewInteger(groupId.(int))
	}

	if v, ok := d.GetOk("dimensions"); ok {
		dimensions := v.(map[string]interface{})
		createAlarmDimensions := make([]ess.CreateAlarmDimension, 0, len(dimensions))
		for k, v := range dimensions {
			if k == UserId || k == ScalingGroup {
				return nil, errmsgs.WrapError(errmsgs.Error("Invalide dimension keys, %s", k))
			}
			if k != "" {
				dimension := ess.CreateAlarmDimension{
					DimensionKey:   k,
					DimensionValue: v.(string),
				}
				createAlarmDimensions = append(createAlarmDimensions, dimension)
			}
		}
		request.Dimension = &createAlarmDimensions
	}

	return request, nil
}
