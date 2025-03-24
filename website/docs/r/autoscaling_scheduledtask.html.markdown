---
subcategory: "AutoScaling"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_autoscaling_scheduledtask"
sidebar_current: "docs-Alibabacloudstack-autoscaling-scheduledtask"
description: |- 
  Provides a autoscaling Scheduledtask resource.
---

# alibabacloudstack_autoscaling_scheduledtask
-> **NOTE:** Alias name has: `alibabacloudstack_ess_scheduled_task`

Provides a autoscaling Scheduledtask resource.

## Example Usage

```hcl
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  name = "${var.name}_vsw"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block = "172.16.0.0/24"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

resource "alibabacloudstack_ecs_securitygroup" "default" {
  name   = "${var.name}_sg"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
}

resource "alibabacloudstack_security_group_rule" "default" {
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "22/22"
  priority          = 1
  security_group_id = "${alibabacloudstack_ecs_securitygroup.default.id}"
  cidr_ip           = "172.16.0.0/24"
}

data "alibabacloudstack_images" "default" {
  name_regex  = "^ubuntu_"
  most_recent = true
  owners      = "system"
}

data "alibabacloudstack_instance_types" "all" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
}

data "alibabacloudstack_instance_types" "any_n4" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  instance_type_family = "ecs.n4"
  sorted_by            = "Memory"
}

data "alibabacloudstack_instance_types" "default" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  cpu_core_count       = 1
  memory_size          = 1
  instance_type_family = "ecs.n4"
  sorted_by            = "Memory"
}

locals {
  default_instance_type_id = try(element(sort(length(data.alibabacloudstack_instance_types.default.instance_types) > 0 ? data.alibabacloudstack_instance_types.default.ids : data.alibabacloudstack_instance_types.any_n4.ids), 0), sort(data.alibabacloudstack_instance_types.all.ids)[0])
}

resource "alibabacloudstack_ecs_instance" "default" {
  image_id             = "${data.alibabacloudstack_images.default.images.0.id}"
  instance_type        = "${local.default_instance_type_id}"
  system_disk_category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
  system_disk_size     = 20
  system_disk_name     = "test_sys_disk"
  security_groups      = [alibabacloudstack_ecs_securitygroup.default.id]
  instance_name        = "${var.name}_ecs"
  vswitch_id           = alibabacloudstack_vpc_vswitch.default.id
  zone_id             = data.alibabacloudstack_zones.default.zones.0.id
  is_outdated          = false
  lifecycle {
    ignore_changes = [
      instance_type
    ]
  }
}

variable "name" {
  default = "tf-testAccEssScheduleConfig-227425"
}

resource "alibabacloudstack_ess_scaling_group" "default" {
  min_size           = 0
  max_size           = 2
  default_cooldown   = 20
  removal_policies   = ["OldestInstance", "NewestInstance"]
  scaling_group_name = "${var.name}"
  vswitch_ids        = ["${alibabacloudstack_vpc_vswitch.default.id}"]
}

resource "alibabacloudstack_ecs_deployment_set" "default" {
  strategy            = "Availability"
  domain              = "Default"
  granularity         = "Host"
  deployment_set_name = "example_value"
  description         = "example_value"
}

resource "alibabacloudstack_ess_scaling_configuration" "default" {
  scaling_group_id = "${alibabacloudstack_ess_scaling_group.default.id}"
  image_id         = "${data.alibabacloudstack_images.default.images.0.id}"
  instance_type    = "ecs.e4.small"
  security_group_ids = [alibabacloudstack_ecs_securitygroup.default.id]
  force_delete     = true
  active          = true
  enable          = true
  deployment_set_id = alibabacloudstack_ecs_deployment_set.default.id
}

resource "alibabacloudstack_ess_scaling_rule" "default" {
  scaling_group_id = "${alibabacloudstack_ess_scaling_group.default.id}"
  adjustment_type  = "TotalCapacity"
  adjustment_value = "1"
  cooldown        = 0
}

resource "alibabacloudstack_ess_scheduled_task" "default" {
  scheduled_action     = "${alibabacloudstack_ess_scaling_rule.default.ari}"
  launch_time          = "2025-03-08T17:50Z"
  scheduled_task_name  = "${var.name}"
  description         = "This is an example of a scheduled task."
  launch_expiration_time = 3600
  recurrence_type      = "Daily"
  recurrence_value     = "1"
  recurrence_end_time  = "2025-04-08T17:50Z"
  task_enabled        = true
}
```

## Argument Reference

The following arguments are supported:

* `scheduled_action` - (Required) The action to be performed when the scheduled task is triggered. It must be the unique identifier of a scaling rule.
* `launch_time` - (Required) The time at which the scheduled task is triggered. Specify the time in the ISO 8601 standard in the `YYYY-MM-DDThh:mm:ssZ` format. The time must be in UTC. You cannot enter a time point later than 90 days from the date of scheduled task creation. If the `recurrence_type` parameter is specified, the task is executed repeatedly at the time specified by `launch_time`. Otherwise, the task is only executed once at the date and time specified by `launch_time`.
* `scheduled_task_name` - (Optional) Display name of the scheduled task, which must be 2-40 characters (English or Chinese) long.
* `description` - (Optional) Description of the scheduled task, which is 2-200 characters (English or Chinese) long.
* `launch_expiration_time` - (Optional) After the scheduled task trigger operation fails, retry within this time. The unit is seconds, and the value range is 0~21600. Default value: 600.
* `recurrence_type` - (Optional) Specifies the recurrence type of the scheduled task. If set, both `recurrence_value` and `recurrence_end_time` must be set. Valid values:
  * `Daily`: The scheduled task is executed once every specified number of days.
  * `Weekly`: The scheduled task is executed on each specified day of a week.
  * `Monthly`: The scheduled task is executed on each specified day of a month.
  * `Cron`: The scheduled task is executed based on the specified cron expression.
* `recurrence_value` - (Optional) Specifies how often a scheduled task recurs. The valid value depends on `recurrence_type`:
  * `Daily`: You can enter one value. Valid values: 1 to 31.
  * `Weekly`: You can enter multiple values and separate them with commas (,). For example, the values 0 to 6 correspond to the days of the week in sequence from Sunday to Saturday.
  * `Monthly`: You can enter two values in A-B format. Valid values of A and B: 1 to 31. The value of B must be greater than or equal to the value of A.
  * `Cron`: You can enter a cron expression which is written in UTC and consists of five fields: minute, hour, day of month (date), month, and day of week. The expression can contain wildcard characters including commas (,), question marks (?), hyphens (-), asterisks (*), number signs (#), forward slashes (/), and the L and W letters.
* `recurrence_end_time` - (Optional) Specifies the end time after which the scheduled task is no longer repeated. Specify the time in the ISO 8601 standard in the `YYYY-MM-DDThh:mm:ssZ` format. The time must be in UTC. You cannot enter a time point later than 365 days from the date of scheduled task creation.
* `task_enabled` - (Optional) Specifies whether to start the scheduled task. Default value: `true`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the scheduled task.
* `description` - Description of the scheduled task.
* `recurrence_type` - The type of the scheduled task that is repeated.
* `recurrence_value` - The value of the repeated execution of the scheduled task.
* `recurrence_end_time` - The end time of the recurring scheduled task.