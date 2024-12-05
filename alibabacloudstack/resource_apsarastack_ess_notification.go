package alibabacloudstack

import (
	"fmt"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackEssNotification() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackEssNotificationCreate,
		Read:   resourceAlibabacloudStackEssNotificationRead,
		Update: resourceAlibabacloudStackEssNotificationUpdate,
		Delete: resourceAlibabacloudStackEssNotificationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
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
}

func resourceAlibabacloudStackEssNotificationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := ess.CreateCreateNotificationConfigurationRequest()
	request.RegionId = client.RegionId

	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{ "Product": "ess", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

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
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_ess_notification", request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	d.SetId(fmt.Sprintf("%s:%s", request.ScalingGroupId, request.NotificationArn))
	return resourceAlibabacloudStackEssNotificationRead(d, meta)
}

func resourceAlibabacloudStackEssNotificationRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	essService := EssService{client}
	object, err := essService.DescribeEssNotification(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("scaling_group_id", object.ScalingGroupId)
	d.Set("notification_arn", object.NotificationArn)
	d.Set("notification_types", object.NotificationTypes.NotificationType)
	return nil
}

func resourceAlibabacloudStackEssNotificationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := ess.CreateModifyNotificationConfigurationRequest()
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{ "Product": "ess", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

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
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return resourceAlibabacloudStackEssNotificationRead(d, meta)
}

func resourceAlibabacloudStackEssNotificationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	essService := EssService{client}
	request := ess.CreateDeleteNotificationConfigurationRequest()
	request.RegionId = client.RegionId
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{ "Product": "ess", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

	parts := strings.SplitN(d.Id(), ":", 2)

	request.ScalingGroupId = parts[0]
	request.NotificationArn = parts[1]

	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DeleteNotificationConfiguration(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"NotificationConfigurationNotExist", "InvalidScalingGroupId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return WrapError(essService.WaitForEssNotification(d.Id(), Deleted, DefaultTimeout))
}
