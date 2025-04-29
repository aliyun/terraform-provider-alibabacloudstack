---
subcategory: "Auto Scaling (ESS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_autoscaling_notifications"
sidebar_current: "docs-Alibabacloudstack-datasource-autoscaling-notifications"
description: |- 
  Provides a list of autoscaling notifications owned by an alibabacloudstack account.
---

# alibabacloudstack_autoscaling_notifications
-> **NOTE:** Alias name has: `alibabacloudstack_ess_notifications`

This data source provides a list of autoscaling notifications in an Alibabacloudstack account according to the specified filters.

## Example Usage

```hcl
data "alibabacloudstack_autoscaling_notifications" "example" {
  scaling_group_id = "sg-1234567890abcdef"
  ids              = ["notification-id-1", "notification-id-2"]
  output_file      = "notifications_output.txt"
}

output "first_notification_arn" {
  value = data.alibabacloudstack_autoscaling_notifications.example.notifications[0].notification_arn
}
```

## Argument Reference

The following arguments are supported:

* `scaling_group_id` - (Required) The ID of the scaling group to which the notifications belong.
* `ids` - (Optional) A list of notification IDs. If specified, the data source will only return notifications that match these IDs.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of notification IDs.
* `notifications` - A list of autoscaling notifications. Each element contains the following attributes:
  * `id` - The unique identifier of the notification.
  * `scaling_group_id` - The ID of the scaling group associated with the notification.
  * `notification_arn` - The Alibaba Cloud Resource Name (ARN) for the notification object.
  * `notification_types` - One or more types of auto-scaling events and resource change notifications. These can include lifecycle events, scaling activities, etc.