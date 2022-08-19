package apsarastack

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
)

func dataSourceApsaraStackEssNotifications() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApsaraStackEssNotificationsRead,
		Schema: map[string]*schema.Schema{
			"scaling_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"notifications": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"notification_arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"notification_types": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"scaling_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceApsaraStackEssNotificationsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	request := ess.CreateDescribeNotificationConfigurationsRequest()
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ess", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.RegionId = client.RegionId
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	if scalingGroupId, ok := d.GetOk("scaling_group_id"); ok && scalingGroupId.(string) != "" {
		request.ScalingGroupId = scalingGroupId.(string)
	}
	var allNotifications []ess.NotificationConfigurationModel
	for {
		raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.DescribeNotificationConfigurations(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "apsarastack_ess_notifications", request.GetActionName(), ApsaraStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response := raw.(*ess.DescribeNotificationConfigurationsResponse)
		if len(response.NotificationConfigurationModels.NotificationConfigurationModel) < 1 {
			break
		}
		allNotifications = append(allNotifications, response.NotificationConfigurationModels.NotificationConfigurationModel...)
		if len(response.NotificationConfigurationModels.NotificationConfigurationModel) < PageSizeLarge {
			break
		} else {
			continue
		}
	}
	var filteredNotifications = make([]ess.NotificationConfigurationModel, 0)
	idsMap := make(map[string]string)
	if ids, okIds := d.GetOk("ids"); okIds {
		for _, i := range ids.([]interface{}) {
			idsMap[i.(string)] = i.(string)
		}
		for _, n := range allNotifications {
			if _, ok := idsMap[n.NotificationArn]; !ok {
				continue
			}
			filteredNotifications = append(filteredNotifications, n)
		}

	} else {
		filteredNotifications = allNotifications
	}

	return notificationsDescriptionAttribute(d, filteredNotifications, meta)
}

func notificationsDescriptionAttribute(d *schema.ResourceData, notifications []ess.NotificationConfigurationModel, meta interface{}) error {
	var ids []string
	var s = make([]map[string]interface{}, 0)
	for _, n := range notifications {
		mapping := map[string]interface{}{
			"notification_arn":   n.NotificationArn,
			"notification_types": n.NotificationTypes.NotificationType,
			"scaling_group_id":   n.ScalingGroupId,
		}
		ids = append(ids, n.NotificationArn)
		s = append(s, mapping)
	}
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("notifications", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
