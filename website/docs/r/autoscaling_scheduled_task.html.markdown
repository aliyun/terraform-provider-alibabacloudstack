---
subcategory: "Auto Scaling"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_autoscaling_scheduled_task"
description: |-
  Provides an Auto Scaling Scheduled Task resource.
---

# alibabacloudstack_autoscaling_scheduled_task

Provides an Auto Scaling Scheduled Task resource that allows you to create, modify, and delete scheduled tasks for scaling groups.

## Example Usage

### Basic Usage

```hcl
resource "alibabacloudstack_autoscaling_scheduled_task" "default" {
  scheduled_action       = "acs:ess:::scalingGroup/${alibabacloudstack_autoscaling_scaling_group.default.id}:changeInCapacity:1"
  launch_time            = "2026-06-01T15:04Z"
  scheduled_task_name    = "tf-testAccScheduledTask"
  description            = "Test scheduled task"
  launch_expiration_time = 600
  recurrence_type        = "Daily"
  recurrence_value       = "1"
  task_enabled           = true
  scaling_group_id       = alibabacloudstack_autoscaling_scaling_group.default.id
}
```

## Argument Reference

The following arguments are supported:

* `scheduled_action` - (Required) The operation to be performed when the scheduled task is triggered. Format: `acs:ess:::scalingGroup/<scalingGroupId>:changeInCapacity:<value>`.
* `launch_time` - (Required) The time when the scheduled task is triggered. Format: `YYYY-MM-DDThh:mmZ` (UTC time).
* `scheduled_task_name` - (Optional) The name of the scheduled task. Must be 2-64 characters in length, starting with a letter or digit, and can contain letters, digits, underscores (_), hyphens (-), and periods (.).
* `description` - (Optional, Computed) The description of the scheduled task. Must be 2-200 characters in length.
* `launch_expiration_time` - (Optional) The validity period of the scheduled task after it is triggered. Unit: seconds. Valid values: 0 to 21600. Default value: 600.
* `recurrence_type` - (Optional, Computed) The type of recurrence for the scheduled task. Valid values: `Daily`, `Weekly`, `Monthly`.
* `recurrence_value` - (Optional, Computed) The value of recurrence. The format varies by recurrence type:
  - Daily: A positive integer indicating the interval between executions.
  - Weekly: Comma-separated days of the week (0-6, where 0 is Sunday).
  - Monthly: Comma-separated days of the month (1-31).
* `recurrence_end_time` - (Optional, Computed) The end time for recurring scheduled tasks. Format: `YYYY-MM-DDThh:mmZ` (UTC time).
* `task_enabled` - (Optional) Whether the scheduled task is enabled. Valid values: `true`, `false`. Default value: `true`.
* `scaling_group_id` - (Optional, Computed) The ID of the scaling group to which the scheduled task belongs.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the scheduled task.

## Import

Auto Scaling Scheduled Task can be imported using the scheduled task ID, e.g.

```
$ terraform import alibabacloudstack_autoscaling_scheduled_task.example st-12345678
```
