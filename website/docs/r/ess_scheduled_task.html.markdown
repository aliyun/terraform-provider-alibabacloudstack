---
subcategory: "Auto Scaling(ESS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ess_schedule"
sidebar_current: "docs-alibabacloudstack-resource-ess-schedule"
description: |-
  Provides a ESS schedule resource.
---

# alibabacloudstack\_ess\_scheduled\_task

Provides a ESS schedule resource.

## Example Usage

```
variable "name" {
  default = "essscheduleconfig"
}

data "alibabacloudstack_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alibabacloudstack_instance_types" "default" {
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  cpu_core_count    = 2
  memory_size       = 4
}

data "alibabacloudstack_images" "default" {
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
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

resource "alibabacloudstack_security_group" "default" {
  name   = "${var.name}"
  vpc_id = "${alibabacloudstack_vpc.default.id}"
}

resource "alibabacloudstack_security_group_rule" "default" {
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "22/22"
  priority          = 1
  security_group_id = "${alibabacloudstack_security_group.default.id}"
  cidr_ip           = "172.16.0.0/24"
}

resource "alibabacloudstack_ess_scaling_group" "default" {
  min_size           = 1
  max_size           = 1
  scaling_group_name = "${var.name}"
  vswitch_ids        = ["${alibabacloudstack_vswitch.default.id}"]
  removal_policies   = ["OldestInstance", "NewestInstance"]
}

resource "alibabacloudstack_ess_scaling_configuration" "default" {
  scaling_group_id  = "${alibabacloudstack_ess_scaling_group.default.id}"
  image_id          = "${data.alibabacloudstack_images.default.images.0.id}"
  instance_type     = "${data.alibabacloudstack_instance_types.default.instance_types.0.id}"
  security_group_id = "${alibabacloudstack_security_group.default.id}"
  force_delete      = "true"
}

resource "alibabacloudstack_ess_scaling_rule" "default" {
  scaling_group_id = "${alibabacloudstack_ess_scaling_group.default.id}"
  adjustment_type  = "TotalCapacity"
  adjustment_value = 2
  cooldown         = 60
}

resource "alibabacloudstack_ess_scheduled_task" "default" {
  scheduled_action    = "${alibabacloudstack_ess_scaling_rule.default.ari}"
  launch_time         = "2019-05-21T11:37Z"
  scheduled_task_name = "${var.name}"
}
```


## Argument Reference

The following arguments are supported:

* `scheduled_action` - (Required) The operation to be performed when a scheduled task is triggered. Enter the unique identifier of a scaling rule.
* `scheduled_task_name` - (Optional) Display name of the scheduled task, which must be 2-40 characters (English or Chinese) long.
* `description` - (Optional) Description of the scheduled task, which is 2-200 characters (English or Chinese) long.
* `launch_time` - (Required) The time at which the scheduled task is triggered. Specify the time in the ISO 8601 standard in the YYYY-MM-DDThh:mm:ssZ format. 
The time must be in UTC. You cannot enter a time point later than 90 days from the date of scheduled task creation. 
If the `recurrence_type` parameter is specified, the task is executed repeatedly at the time specified by LaunchTime. 
Otherwise, the task is only executed once at the date and time specified by LaunchTime.
* `launch_expiration_time` - (Optional) The time period during which a failed scheduled task is retried. Unit: seconds. Valid values: 0 to 21600. Default value: 600
* `recurrence_type` - (Optional) Specifies the recurrence type of the scheduled task. 
If set, both `recurrence_value` and `recurrence_end_time` must be set. Valid values:
    - Daily: The scheduled task is executed once every specified number of days.
    - Weekly: The scheduled task is executed on each specified day of a week.
    - Monthly: The scheduled task is executed on each specified day of a month.
    - Cron:  The scheduled task is executed based on the specified cron expression.
* `recurrence_value` - (Optional) Specifies how often a scheduled task recurs. The valid value depends on `recurrence_type`
    - Daily: You can enter one value. Valid values: 1 to 31.
    - Weekly: You can enter multiple values and separate them with commas (,). For example, the values 0 to 6 correspond to the days of the week in sequence from Sunday to Saturday.
    - Monthly: You can enter two values in A-B format. Valid values of A and B: 1 to 31. The value of B must be greater than or equal to the value of A.
    - Cron: You can enter a cron expression which is written in UTC and consists of five fields: minute, hour, day of month (date), month, and day of week. The expression can contain wildcard characters including commas (,), question marks (?), hyphens (-), asterisks (*), number signs (#), forward slashes (/), and the L and W letters.
* `recurrence_end_time` - (Optional) Specifies the end time after which the scheduled task is no longer repeated. 
Specify the time in the ISO 8601 standard in the YYYY-MM-DDThh:mm:ssZ format. 
The time must be in UTC. You cannot enter a time point later than 365 days from the date of scheduled task creation.                                
* `task_enabled` - (Optional) Specifies whether to start the scheduled task. Default to true.
                                  
                                 
## Attributes Reference

The following attributes are exported:

* `id` - The schedule task ID.
