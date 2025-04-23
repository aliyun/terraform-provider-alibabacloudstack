package alibabacloudstack

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackEssScheduledTask() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"scheduled_action": {
				Type:     schema.TypeString,
				Required: true,
			},
			"launch_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"scheduled_task_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"launch_expiration_time": {
				Type:         schema.TypeInt,
				Default:      600,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 21600),
			},
			"recurrence_type": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Daily", "Weekly", "Monthly"}, false),
			},
			"recurrence_value": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"recurrence_end_time": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"task_enabled": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackEssScheduledTaskCreate,
		resourceAlibabacloudStackEssScheduledTaskRead, resourceAlibabacloudStackEssScheduledTaskUpdate, resourceAlibabacloudStackEssScheduledTaskDelete)
	return resource
}

func resourceAlibabacloudStackEssScheduledTaskCreate(d *schema.ResourceData, meta interface{}) error {
	request := buildAlibabacloudStackEssScheduledTaskArgs(d)
	client := meta.(*connectivity.AlibabacloudStackClient)
	client.InitRpcRequest(*request.RpcRequest)

	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.CreateScheduledTask(request)
	})
	bresponse, ok := raw.(*ess.CreateScheduledTaskResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ess_scheduled_task", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	d.SetId(bresponse.ScheduledTaskId)

	return nil
}

func resourceAlibabacloudStackEssScheduledTaskRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AlibabacloudStackClient)
	essService := EssService{client}

	object, err := essService.DescribeEssScheduledTask(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("scheduled_action", object.ScheduledAction)
	d.Set("launch_time", object.LaunchTime)
	d.Set("scheduled_task_name", object.ScheduledTaskName)
	d.Set("description", object.Description)
	d.Set("launch_expiration_time", object.LaunchExpirationTime)
	d.Set("recurrence_type", object.RecurrenceType)
	d.Set("recurrence_value", object.RecurrenceValue)
	d.Set("recurrence_end_time", object.RecurrenceEndTime)
	d.Set("task_enabled", object.TaskEnabled)

	return nil
}

func resourceAlibabacloudStackEssScheduledTaskUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := ess.CreateModifyScheduledTaskRequest()
	client.InitRpcRequest(*request.RpcRequest)

	request.ScheduledTaskId = d.Id()
	request.LaunchExpirationTime = requests.NewInteger(d.Get("launch_expiration_time").(int))

	if d.HasChange("scheduled_task_name") {
		request.ScheduledTaskName = d.Get("scheduled_task_name").(string)
	}

	if d.HasChange("description") {
		request.Description = d.Get("description").(string)
	}

	if d.HasChange("scheduled_action") {
		request.ScheduledAction = d.Get("scheduled_action").(string)
	}

	if d.HasChange("launch_time") {
		request.LaunchTime = d.Get("launch_time").(string)
	}

	if d.HasChanges("recurrence_type", "recurrence_value", "recurrence_end_time") {
		request.RecurrenceType = d.Get("recurrence_type").(string)
		request.RecurrenceValue = d.Get("recurrence_value").(string)
		request.RecurrenceEndTime = d.Get("recurrence_end_time").(string)
	}

	if d.HasChange("task_enabled") {
		request.TaskEnabled = requests.NewBoolean(d.Get("task_enabled").(bool))
	}

	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.ModifyScheduledTask(request)
	})
	bresponse, ok := raw.(*ess.ModifyScheduledTaskResponse)
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

func resourceAlibabacloudStackEssScheduledTaskDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	essService := EssService{client}

	request := ess.CreateDeleteScheduledTaskRequest()
	client.InitRpcRequest(*request.RpcRequest)

	request.ScheduledTaskId = d.Id()

	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DeleteScheduledTask(request)
	})
	bresponse, ok := raw.(*ess.DeleteScheduledTaskResponse)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidScheduledTaskId.NotFound"}) {
			return nil
		}
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return errmsgs.WrapError(essService.WaitForEssScheduledTask(d.Id(), Deleted, DefaultTimeout))
}

func buildAlibabacloudStackEssScheduledTaskArgs(d *schema.ResourceData) *ess.CreateScheduledTaskRequest {
	request := ess.CreateCreateScheduledTaskRequest()

	request.ScheduledAction = d.Get("scheduled_action").(string)
	request.LaunchTime = d.Get("launch_time").(string)

	if v, ok := d.GetOk("task_enabled"); ok {
		request.TaskEnabled = requests.NewBoolean(v.(bool))
	}

	if v, ok := d.GetOk("scheduled_task_name"); ok && v.(string) != "" {
		request.ScheduledTaskName = v.(string)
	}

	if v, ok := d.GetOk("description"); ok && v.(string) != "" {
		request.Description = v.(string)
	}

	if v, ok := d.GetOk("recurrence_type"); ok && v.(string) != "" {
		request.RecurrenceType = v.(string)
	}

	if v, ok := d.GetOk("recurrence_value"); ok && v.(string) != "" {
		request.RecurrenceValue = v.(string)
	}

	if v, ok := d.GetOk("recurrence_end_time"); ok && v.(string) != "" {
		request.RecurrenceEndTime = v.(string)
	}

	if v, ok := d.GetOk("launch_expiration_time"); ok && v.(int) != 0 {
		request.LaunchExpirationTime = requests.NewInteger(v.(int))
	}

	return request
}
