package alibabacloudstack

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackEssLifecycleHook() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"scaling_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Deprecated:   "Field 'name' is deprecated and will be removed in a future release. Please use new field 'lifecycle_hook_name' instead.",
				ConflictsWith: []string{"lifecycle_hook_name"},
			},
			"lifecycle_hook_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
				ConflictsWith: []string{"name"},
			},
			"lifecycle_transition": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"SCALE_IN", "SCALE_OUT"}, false),
			},
			"heartbeat_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      600,
				ValidateFunc: validation.IntBetween(30, 21600),
			},
			"default_result": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "CONTINUE",
				ValidateFunc: validation.StringInSlice([]string{"CONTINUE", "ABANDON"}, false),
			},
			"notification_arn": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"notification_metadata": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackEssLifeCycleHookCreate, resourceAlibabacloudStackEssLifeCycleHookRead, resourceAlibabacloudStackEssLifeCycleHookUpdate, resourceAlibabacloudStackEssLifeCycleHookDelete)
	return resource
}

func resourceAlibabacloudStackEssLifeCycleHookCreate(d *schema.ResourceData, meta interface{}) error {
	request := buildAlibabacloudStackEssLifeCycleHookArgs(d)
	client := meta.(*connectivity.AlibabacloudStackClient)
	client.InitRpcRequest(*request.RpcRequest)

	if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.CreateLifecycleHook(request)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{errmsgs.Throttling}) {
				return resource.RetryableError(err)
			}
			errmsg := ""
			if bresponse, ok := raw.(*ess.CreateLifecycleHookResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ess_lifecyclehook", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*ess.CreateLifecycleHookResponse)
		d.SetId(response.LifecycleHookId)
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func resourceAlibabacloudStackEssLifeCycleHookRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AlibabacloudStackClient)
	essService := EssService{client}

	object, err := essService.DescribeEssLifecycleHook(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("scaling_group_id", object.ScalingGroupId)
	connectivity.SetResourceData(d, object.LifecycleHookName, "lifecycle_hook_name", "name")
	d.Set("lifecycle_transition", object.LifecycleTransition)
	d.Set("heartbeat_timeout", object.HeartbeatTimeout)
	d.Set("default_result", object.DefaultResult)
	d.Set("notification_arn", object.NotificationArn)
	d.Set("notification_metadata", object.NotificationMetadata)

	return nil
}

func resourceAlibabacloudStackEssLifeCycleHookUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := ess.CreateModifyLifecycleHookRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.LifecycleHookId = d.Id()

	if d.HasChange("lifecycle_transition") {
		request.LifecycleTransition = d.Get("lifecycle_transition").(string)
	}

	if d.HasChange("heartbeat_timeout") {
		request.HeartbeatTimeout = requests.NewInteger(d.Get("heartbeat_timeout").(int))
	}

	if d.HasChange("default_result") {
		request.DefaultResult = d.Get("default_result").(string)
	}

	if d.HasChange("notification_arn") {
		request.NotificationArn = d.Get("notification_arn").(string)
	}

	if d.HasChange("notification_metadata") {
		request.NotificationMetadata = d.Get("notification_metadata").(string)
	}

	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.ModifyLifecycleHook(request)
	})
	if err != nil {
		errmsg := ""
		if bresponse, ok := raw.(*ess.ModifyLifecycleHookResponse); ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}

func resourceAlibabacloudStackEssLifeCycleHookDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	essService := EssService{client}
	request := ess.CreateDeleteLifecycleHookRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.LifecycleHookId = d.Id()

	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DeleteLifecycleHook(request)
	})
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidLifecycleHookId.NotExist"}) {
			return nil
		}
		errmsg := ""
		if bresponse, ok := raw.(*ess.DeleteLifecycleHookResponse); ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return errmsgs.WrapError(essService.WaitForEssLifecycleHook(d.Id(), Deleted, DefaultTimeout))
}

func buildAlibabacloudStackEssLifeCycleHookArgs(d *schema.ResourceData) *ess.CreateLifecycleHookRequest {
	request := ess.CreateCreateLifecycleHookRequest()

	request.ScalingGroupId = d.Get("scaling_group_id").(string)

	if v, ok := connectivity.GetResourceDataOk(d, "lifecycle_hook_name", "name"); ok && v.(string) != "" {
		request.LifecycleHookName = v.(string)
	}

	if transition := d.Get("lifecycle_transition").(string); transition != "" {
		request.LifecycleTransition = transition
	}

	if timeout, ok := d.GetOk("heartbeat_timeout"); ok && timeout.(int) > 0 {
		request.HeartbeatTimeout = requests.NewInteger(timeout.(int))
	}

	if v, ok := d.GetOk("default_result"); ok && v.(string) != "" {
		request.DefaultResult = v.(string)
	}

	if v, ok := d.GetOk("notification_arn"); ok && v.(string) != "" {
		request.NotificationArn = v.(string)
	}

	if v, ok := d.GetOk("notification_metadata"); ok && v.(string) != "" {
		request.NotificationMetadata = v.(string)
	}

	return request
}