---
subcategory: "AutoScaling"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_autoscaling_notification"
sidebar_current: "docs-Alibabacloudstack-autoscaling-notification"
description: |- 
  Provides a autoscaling Notification resource.
---

# alibabacloudstack_autoscaling_notification
-> **NOTE:** Alias name has: `alibabacloudstack_ess_notification`

Provides a autoscaling Notification resource.

## Example Usage

```hcl
variable "name" {
    default = "tf-testAccAutoscalingNotification-%d"
}

data "alibabacloudstack_regions" "default" {
    current = true
}

data "alibabacloudstack_account" "default" {
}

data "alibabacloudstack_zones" "default" {
    available_disk_category     = "cloud_efficiency"
    available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vpc" "default" {
    name       = "${var.name}"
    cidr_block = "172.16.0.0/16"
}
    
resource "alibabacloudstack_vswitch" "default" {
    vpc_id            = "${alibabacloudstack_vpc.default.id}"
    cidr_block        = "172.16.0.0/24"
    availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
    name              = "${var.name}"
}

resource "alibabacloudstack_autoscaling_scaling_group" "default" {
    min_size = 1
    max_size = 1
    scaling_group_name = "${var.name}"
    removal_policies = ["OldestInstance", "NewestInstance"]
    vswitch_ids = ["${alibabacloudstack_vswitch.default.id}"]
}

resource "alibabacloudstack_mns_queue" "default"{
    name="${var.name}"
}

resource "alibabacloudstack_autoscaling_notification" "default" {
    scaling_group_id = "${alibabacloudstack_autoscaling_scaling_group.default.id}"
    notification_types = ["AUTOSCALING:SCALE_OUT_SUCCESS","AUTOSCALING:SCALE_OUT_ERROR"]
    notification_arn = "acs:ess:${data.alibabacloudstack_regions.default.regions.0.id}:${data.alibabacloudstack_account.default.id}:queue/${alibabacloudstack_mns_queue.default.name}"
}
```

## Argument Reference

The following arguments are supported:

* `scaling_group_id` - (Required, ForceNew) The ID of the Auto Scaling group. Changing this will force the creation of a new resource.
* `notification_arn` - (Required, ForceNew) The AlibabaCloudStack Cloud Resource Name (ARN) for the notification object. The format of `notification_arn` is `acs:ess:{region}:{account-id}:{resource-relative-id}`. Valid values for `resource-relative-id` include:
  * `cloudmonitor`: For CloudMonitor notifications.
  * `queue/{queue-name}`: For Message Queue (MNS) queue-based notifications.
  * `topic/{topic-name}`: For Message Queue (MNS) topic-based notifications.
* `notification_types` - (Required) One or more types of auto-scaling events and resource change notifications. Supported notification types:
  * `AUTOSCALING:SCALE_OUT_SUCCESS`
  * `AUTOSCALING:SCALE_IN_SUCCESS`
  * `AUTOSCALING:SCALE_OUT_ERROR`
  * `AUTOSCALING:SCALE_IN_ERROR`
  * `AUTOSCALING:SCALE_REJECT`
  * `AUTOSCALING:SCALE_OUT_START`
  * `AUTOSCALING:SCALE_IN_START`
  * `AUTOSCALING:SCHEDULE_TASK_EXPIRING`

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the notification resource, which is composed of `scaling_group_id` and `notification_arn` in the format of `<scaling_group_id>:<notification_arn>`.