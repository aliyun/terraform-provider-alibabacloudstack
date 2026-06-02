---
subcategory: "Auto Scaling"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ess_lifecycle_hook"
sidebar_current: "docs-Alibabacloudstack-resource-ess-lifecycle-hook"
description: |-
  Provides a ESS lifecycle hook resource.
---

# alibabacloudstack_ess_lifecycle_hook

-> **NOTE:** This resource is unsupported on ApsaraStack and will be removed in version 3.21.0.

-> **NOTE:** Alias name has: `alibabacloudstack_autoscaling_lifecyclehook`

Provides a ESS lifecycle hook resource.

## Example Usage

```
data "alibabacloudstack_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vpc" "foo" {
  name       = "testAccEssScalingGroup_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vswitch" "foo" {
  vpc_id            = "${alibabacloudstack_vpc.foo.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

resource "alibabacloudstack_vswitch" "bar" {
  vpc_id            = "${alibabacloudstack_vpc.foo.id}"
  cidr_block        = "172.16.1.0/24"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

resource "alibabacloudstack_ess_scaling_group" "foo" {
  min_size           = 1
  max_size           = 1
  scaling_group_name = "testAccEssScaling_group"
  removal_policies   = ["OldestInstance", "NewestInstance"]
  vswitch_ids        = ["${alibabacloudstack_vswitch.foo.id}", "${alibabacloudstack_vswitch.bar.id}"]
}

resource "alibabacloudstack_ess_lifecycle_hook" "foo" {
  scaling_group_id      = "${alibabacloudstack_ess_scaling_group.foo.id}"
  lifecycle_hook_name   = "testAccEssLifecycle_hook"
  lifecycle_transition  = "SCALE_OUT"
  heartbeat_timeout     = 400
  notification_metadata = "helloworld"
}
```

## Argument Reference

The following arguments are supported:

* `scaling_group_id` - (Required, ForceNew) The ID of the Auto Scaling group to which you want to assign the lifecycle hook.

* `lifecycle_hook_name` - (Optional, ForceNew) The name of the lifecycle hook. The name must contain 2 to 128 characters. It must start with a letter and cannot start with http:// or https://. It can contain letters, digits, underscores (_), hyphens (-), and periods (.). This parameter conflicts with the `name` parameter.

* `lifecycle_transition` - (Required) The type of lifecycle hook. Valid values: `SCALE_IN`, `SCALE_OUT`.

* `heartbeat_timeout` - (Optional) The maximum waiting time before the lifecycle hook times out. Valid values: 30 to 21600. Unit: seconds. Default value: 600.

* `default_result` - (Optional) The action that Auto Scaling takes when the lifecycle hook times out. Valid values: `CONTINUE`, `ABANDON`. Default value: `CONTINUE`.

* `notification_arn` - (Optional) The ARN of the EventBridge event bus to which notifications are sent when the lifecycle hook is triggered.

* `notification_metadata` - (Optional) The additional information that you want to include when Auto Scaling sends a message to the notification target.

* `name` - (Optional, ForceNew, Deprecated) The name of the lifecycle hook. This field is deprecated. Use `lifecycle_hook_name` instead. If this parameter value is not specified, the default value is lifecycle hook id. This parameter conflicts with the `lifecycle_hook_name` parameter.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the lifecycle hook.

* `scaling_group_id` - The ID of the Auto Scaling group to which the lifecycle hook belongs.

* `lifecycle_hook_name` - The name of the lifecycle hook.

* `name` - (Deprecated) The name of the lifecycle hook.

* `lifecycle_transition` - The type of lifecycle hook.

* `heartbeat_timeout` - The maximum waiting time before the lifecycle hook times out.

* `default_result` - The action that Auto Scaling takes when the lifecycle hook times out.

* `notification_arn` - The ARN of the EventBridge event bus to which notifications are sent.

* `notification_metadata` - The additional information that is sent when the lifecycle hook is triggered.

## Import

ESS Lifecycle Hook can be imported using the lifecycle hook ID, e.g.

```
$ terraform import alibabacloudstack_ess_lifecycle_hook.example lch-bp1234567890abcdef
```
