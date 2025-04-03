---
subcategory: "OOS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_oos_executions"
sidebar_current: "docs-Alibabacloudstack-datasource-oos-executions"
description: |- 
  Provides a list of oos executions owned by an alibabacloudstack account.
---

# alibabacloudstack_oos_executions

This data source provides a list of OOS Executions in an Alibaba Cloud account according to the specified filters.

## Example Usage

```hcl
# Declare the data source

data "alibabacloudstack_oos_executions" "example" {
  ids = ["execution_id"]
  template_name = "name"
  status = "Success"
}

output "first_execution_id" {
  value = "${data.alibabacloudstack_oos_executions.example.executions.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `category` - (Optional, ForceNew) The category of the template. Valid values: `AlarmTrigger`, `EventTrigger`, `Other`, and `TimerTrigger`.
* `end_date` - (Optional, ForceNew) The time when the execution was ended.
* `end_date_after` - (Optional, ForceNew) Execution whose end time is less than or equal to the specified time.
* `executed_by` - (Optional, ForceNew) The user who executed the template.
* `ids` - (Optional, ForceNew) A list of OOS Execution IDs.
* `include_child_execution` - (Optional, ForceNew) Whether to include sub-execution.
* `mode` - (Optional, ForceNew) The mode of OOS Execution. Valid values: `Automatic`, `Debug`.
* `parent_execution_id` - (Optional, ForceNew) The ID of the parent OOS Execution.
* `ram_role` - (Optional, ForceNew) The role that executes the current template.
* `sort_field` - (Optional, ForceNew) The sort field.
* `sort_order` - (Optional, ForceNew) The sort order.
* `start_date_after` - (Optional, ForceNew) The execution whose start time is greater than or equal to the specified time.
* `start_date_before` - (Optional, ForceNew) The execution with start time less than or equal to the specified time.
* `status` - (Optional, ForceNew) The Status of OOS Execution. Valid values: `Cancelled`, `Failed`, `Queued`, `Running`, `Started`, `Success`, `Waiting`.
* `template_name` - (Optional, ForceNew) The name of the execution template.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of OOS Execution IDs.
* `executions` - A list of OOS Executions. Each element contains the following attributes:
  * `id` - ID of the OOS Execution.
  * `parent_execution_id` - The ID of the parent OOS Execution.
  * `category` - The category of the template. Valid values: `AlarmTrigger`, `EventTrigger`, `Other`, and `TimerTrigger`.
  * `counters` - The counters of OOS Execution.
  * `create_date` - The time when the execution was created.
  * `end_date` - The time when the execution was ended.
  * `executed_by` - The user who executed the template.
  * `execution_id` - ID of the OOS Execution.
  * `is_parent` - Whether it includes subtasks.
  * `outputs` - The outputs of OOS Execution.
  * `parameters` - The parameters required by the template.
  * `mode` - The mode of OOS Execution. Valid values: `Automatic`, `Debug`.
  * `ram_role` - The role that executes the current template.
  * `start_date` - The time when the template was started.
  * `status_message` - The message of status.
  * `status_reason` - The reason of status.
  * `template_id` - The ID of the execution template.
  * `template_name` - The name of the execution template.
  * `template_version` - The version of the execution template.
  * `update_date` - The time when the template was updated.
  * `status` - The Status of OOS Execution. Valid values: `Cancelled`, `Failed`, `Queued`, `Running`, `Started`, `Success`, `Waiting`.