---
subcategory: "Auto Scaling(ESS)"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_ess_notification"
sidebar_current: "docs-apsarastack-resource-ess-notification"
description: |-
  Provides a ESS notification resource.
---

# apsarastack\_ess\_notification

Provides a ESS notification resource.

## Example Usage
```
variable "name" {
    default = "tf-testAccEssNotification-%d"
}

data "apsarastack_regions" "default" {
    current = true
}

data "apsarastack_account" "default" {
}

data "apsarastack_zones" "default" {
    available_disk_category     = "cloud_efficiency"
    available_resource_creation = "VSwitch"
}

resource "apsarastack_vpc" "default" {
    name       = "${var.name}"
    cidr_block = "172.16.0.0/16"
}
    
resource "apsarastack_vswitch" "default" {
    vpc_id            = "${apsarastack_vpc.default.id}"
    cidr_block        = "172.16.0.0/24"
    availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
    name              = "${var.name}"
}

resource "apsarastack_ess_scaling_group" "default" {
    min_size = 1
    max_size = 1
    scaling_group_name = "${var.name}"
    removal_policies = ["OldestInstance", "NewestInstance"]
    vswitch_ids = ["${apsarastack_vswitch.default.id}"]
}

resource "apsarastack_mns_queue" "default"{
    name="${var.name}"
}

resource "apsarastack_ess_notification" "default" {
    scaling_group_id = "${apsarastack_ess_scaling_group.default.id}"
    notification_types = ["AUTOSCALING:SCALE_OUT_SUCCESS","AUTOSCALING:SCALE_OUT_ERROR"]
    notification_arn = "acs:ess:${data.apsarastack_regions.default.regions.0.id}:${data.apsarastack_account.default.id}:queue/${apsarastack_mns_queue.default.name}"
}

```

## Argument Reference

The following arguments are supported:

* `scaling_group_id` - (Required, ForceNew) The ID of the Auto Scaling group.
* `notification_arn` - (Required, ForceNew) The Apsarastack Cloud Resource Name (ARN) for the notification object. The format of `notification_arn` is acs:ess:{region}:{account-id}:{resource-relative-id}. Valid values for `resource-relative-id`: 'cloudmonitor', 'queue/', 'topic/'.
* `notification_types` - (Required) The notification types of Auto Scaling events and resource changes. Supported notification types: 'AUTOSCALING:SCALE_OUT_SUCCESS', 'AUTOSCALING:SCALE_IN_SUCCESS', 'AUTOSCALING:SCALE_OUT_ERROR', 'AUTOSCALING:SCALE_IN_ERROR', 'AUTOSCALING:SCALE_REJECT', 'AUTOSCALING:SCALE_OUT_START', 'AUTOSCALING:SCALE_IN_START', 'AUTOSCALING:SCHEDULE_TASK_EXPIRING'.

## Attribute Reference

The following attributes are exported:

* `id` - The ID of notification resource, which is composed of 'scaling_group_id' and 'notification_arn' in the format of '<scaling_group_id>:<notification_arn>'.
