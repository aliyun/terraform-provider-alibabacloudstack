---
subcategory: "Auto Scaling (ESS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_autoscaling_scheduledtasks"
sidebar_current: "docs-Alibabacloudstack-datasource-autoscaling-scheduledtasks"
description: |- 
  Provides a list of autoscaling scheduledtasks owned by an AlibabaCloudStack account.
---

# alibabacloudstack_autoscaling_scheduledtasks
-> **NOTE:** Alias name has: `alibabacloudstack_ess_scheduled_tasks`

This data source provides a list of autoscaling scheduled tasks in an AlibabaCloudStack account according to the specified filters.

## Example Usage

```hcl
data "alibabacloudstack_autoscaling_scheduledtasks" "example" {
  scheduled_task_id = "your-scheduled-task-id"
  name_regex        = "scheduled-task-name-.*"
  ids               = ["task-id-1", "task-id-2"]
  output_file       = "scheduled_tasks_output.txt"
}

output "first_scheduled_task_id" {
  value = data.alibabacloudstack_autoscaling_scheduledtasks.example.tasks.0.id
}
```

## Argument Reference

The following arguments are supported:

* `scheduled_task_id` - (Optional) The ID of the scheduled task.
* `scheduled_action` - (Optional) The action to be performed when the scheduled task is triggered.
* `name_regex` - (Optional) A regex string used to filter the resulting scheduled tasks by name.
* `ids` - (Optional) A list of scheduled task IDs used to filter the results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `ids` - A list of scheduled task IDs.
* `names` - A list of scheduled task names.
* `tasks` - A list of scheduled tasks. Each element contains the following attributes:
  * `id` - The ID of the scheduled task.
  * `name` - The name of the scheduled task.
  * `scheduled_action` - The action to be performed when the scheduled task is triggered.
  * `description` - Description of the scheduled task.
  * `launch_expiration_time` - After the scheduled task trigger operation fails, retry within this time. The unit is seconds, and the value range is 0~21600.
  * `launch_time` - The time at which the scheduled task is triggered.
  * `min_value` - The minimum number of instances in a scaling group when the scaling method of the scheduled task is to specify the number of instances in a scaling group.
  * `max_value` - The maximum number of instances in a scaling group when the scaling method of the scheduled task is to specify the number of instances in a scaling group.
  * `recurrence_type` - Specifies the recurrence type of the scheduled task.
  * `recurrence_value` - Specifies how often a scheduled task recurs.
  * `recurrence_end_time` - Specifies the end time after which the scheduled task is no longer repeated.
  * `task_enabled` - Specifies whether to start the scheduled task. Default to `true`.