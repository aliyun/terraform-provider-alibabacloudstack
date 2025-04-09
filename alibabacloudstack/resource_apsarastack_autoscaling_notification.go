package alibabacloudstack

import (
	"fmt"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackEssNotification() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"notification_arn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"notification_types": {
				Required: true,
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				MinItems: 1,
			},
			"scaling_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackEssNotificationCreate,
		resourceAlibabacloudStackEssNotificationRead, resourceAlibabacloudStackEssNotificationUpdate, resourceAlibabacloudStackEssNotificationDelete)
	return resource
}

func resourceAlibabacloudStackEssNotificationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := ess.CreateCreateNotificationConfigurationRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.ScalingGroupId = d.Get("scaling_group_id").(string)
	request.NotificationArn = d.Get("notification_arn").(string)
	if v, ok := d.GetOk("notification_types"); ok {
		notificationTypes := make([]string, 0)
		notificationTypeList := v.(*schema.Set).List()
		if len(notificationTypeList) > 0 {
			for _, n := range notificationTypeList {
				notificationTypes = append(notificationTypes, n.(string))
			}
		}
		if len(notificationTypes) > 0 {
			request.NotificationType = &notificationTypes
		}
	}

	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.CreateNotificationConfiguration(request)
	})
	response, ok := raw.(*ess.CreateNotificationConfigurationResponse)
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ess_notification", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	d.SetId(fmt.Sprintf("%s:%s", request.ScalingGroupId, request.NotificationArn))
	return nil
}

func resourceAlibabacloudStackEssNotificationRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	essService := EssService{client}
	object, err := essService.DescribeEssNotification(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	d.Set("scaling_group_id", object.ScalingGroupId)
	d.Set("notification_arn", object.NotificationArn)
	d.Set("notification_types", object.NotificationTypes.NotificationType)
	return nil
}

func resourceAlibabacloudStackEssNotificationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := ess.CreateModifyNotificationConfigurationRequest()
	client.InitRpcRequest(*request.RpcRequest)
	parts := strings.SplitN(d.Id(), ":", 2)
	request.ScalingGroupId = parts[0]
	request.NotificationArn = parts[1]
	if d.HasChange("notification_types") {
		v := d.Get("notification_types")
		notificationTypes := make([]string, 0)
		notificationTypeList := v.(*schema.Set).List()
		if len(notificationTypeList) > 0 {
			for _, n := range notificationTypeList {
				notificationTypes = append(notificationTypes, n.(string))
			}
		}
		if len(notificationTypes) > 0 {
			request.NotificationType = &notificationTypes
		}
	}
	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.ModifyNotificationConfiguration(request)
	})
	response, ok := raw.(*ess.ModifyNotificationConfigurationResponse)
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	return nil
}

func resourceAlibabacloudStackEssNotificationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	essService := EssService{client}
	request := ess.CreateDeleteNotificationConfigurationRequest()
	client.InitRpcRequest(*request.RpcRequest)
	parts := strings.SplitN(d.Id(), ":", 2)
	request.ScalingGroupId = parts[0]
	request.NotificationArn = parts[1]

	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DeleteNotificationConfiguration(request)
	})
	response, ok := raw.(*ess.DeleteNotificationConfigurationResponse)
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"NotificationConfigurationNotExist", "InvalidScalingGroupId.NotFound"}) {
			return nil
		}
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	return errmsgs.WrapError(essService.WaitForEssNotification(d.Id(), Deleted, DefaultTimeout))
}
