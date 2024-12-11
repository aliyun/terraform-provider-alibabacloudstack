package alibabacloudstack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackSnapshotPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackSnapshotPolicyCreate,
		Read:   resourceAlibabacloudStackSnapshotPolicyRead,
		Update: resourceAlibabacloudStackSnapshotPolicyUpdate,
		Delete: resourceAlibabacloudStackSnapshotPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:true,
				Computed:true,
				ValidateFunc: validation.StringLenBetween(2, 128),
				Deprecated:   "Field 'name' is deprecated and will be removed in a future release. Please use new field 'auto_snapshot_policy_name' instead.",
				ConflictsWith: []string{"auto_snapshot_policy_name"},
			},
			"auto_snapshot_policy_name": {
				Type:         schema.TypeString,
				Optional:true,
				Computed:true,
				ValidateFunc: validation.StringLenBetween(2, 128),
				ConflictsWith: []string{"name"},
			},
			"repeat_weekdays": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"retention_days": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"tags": tagsSchema(),
			"time_points": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAlibabacloudStackSnapshotPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := ecs.CreateCreateAutoSnapshotPolicyRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.AutoSnapshotPolicyName = connectivity.GetResourceData(d, "auto_snapshot_policy_name", "name").(string)
	if err := errmsgs.CheckEmpty(request.AutoSnapshotPolicyName, schema.TypeString, "auto_snapshot_policy_name", "name"); err != nil {
		return errmsgs.WrapError(err)
	}
	request.RepeatWeekdays = convertListToJsonString(d.Get("repeat_weekdays").(*schema.Set).List())
	request.RetentionDays = requests.NewInteger(d.Get("retention_days").(int))
	request.TimePoints = convertListToJsonString(d.Get("time_points").(*schema.Set).List())

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.CreateAutoSnapshotPolicy(request)
	})
	bresponse, ok := raw.(*ecs.CreateAutoSnapshotPolicyResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_snapshot_policy", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	d.SetId(bresponse.AutoSnapshotPolicyId)

	ecsService := EcsService{client}
	if err := ecsService.WaitForSnapshotPolicy(d.Id(), SnapshotPolicyNormal, DefaultTimeout); err != nil {
		return errmsgs.WrapError(err)
	}

	return resourceAlibabacloudStackSnapshotPolicyRead(d, meta)
}

func resourceAlibabacloudStackSnapshotPolicyRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}
	object, err := ecsService.DescribeSnapshotPolicy(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	connectivity.SetResourceData(d, object.AutoSnapshotPolicyName, "auto_snapshot_policy_name", "name")
	weekdays, err := convertJsonStringToList(object.RepeatWeekdays)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	d.Set("repeat_weekdays", weekdays)
	d.Set("retention_days", object.RetentionDays)
	timePoints, err := convertJsonStringToList(object.TimePoints)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	d.Set("tags", ecsService.tagsToMap(object.Tags.Tag))
	d.Set("time_points", timePoints)

	return nil
}

func resourceAlibabacloudStackSnapshotPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	ecsService := EcsService{client}
	if d.HasChange("tags") {
		if err := ecsService.SetResourceTagsNew(d, "auto_snapshot_policy"); err != nil {
			return errmsgs.WrapError(err)
		}
	}

	request := ecs.CreateModifyAutoSnapshotPolicyExRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.AutoSnapshotPolicyId = d.Id()
	if d.HasChanges("auto_snapshot_policy_name", "name") {
		request.AutoSnapshotPolicyName = connectivity.GetResourceData(d, "auto_snapshot_policy_name", "name").(string)
		if err := errmsgs.CheckEmpty(request.AutoSnapshotPolicyName, schema.TypeString, "auto_snapshot_policy_name", "name"); err != nil {
			return errmsgs.WrapError(err)
		}
	}
	if d.HasChange("repeat_weekdays") {
		request.RepeatWeekdays = convertListToJsonString(d.Get("repeat_weekdays").(*schema.Set).List())
	}
	if d.HasChange("retention_days") {
		request.RetentionDays = requests.NewInteger(d.Get("retention_days").(int))
	}
	if d.HasChange("time_points") {
		request.TimePoints = convertListToJsonString(d.Get("time_points").(*schema.Set).List())
	}

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ModifyAutoSnapshotPolicyEx(request)
	})
	bresponse, ok := raw.(*ecs.ModifyAutoSnapshotPolicyExResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return resourceAlibabacloudStackSnapshotPolicyRead(d, meta)
}

func resourceAlibabacloudStackSnapshotPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}

	request := ecs.CreateDeleteAutoSnapshotPolicyRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.AutoSnapshotPolicyId = d.Id()

	err := resource.Retry(DefaultTimeout*time.Second, func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DeleteAutoSnapshotPolicy(request)
		})
		bresponse, ok := raw.(*ecs.DeleteAutoSnapshotPolicyResponse)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, errmsgs.SnapshotPolicyInvalidOperations) {
				return resource.RetryableError(err)
			}
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return err
	}

	return errmsgs.WrapError(ecsService.WaitForSnapshotPolicy(d.Id(), Deleted, DefaultTimeout))
}
