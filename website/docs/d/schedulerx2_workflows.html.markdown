---
subcategory: "SchedulerX2"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_schedulerx2_workflow"
sidebar_current: "docs-Alibabacloudstack-datasource-schedulerx2-workflow"
description: |-
  Query the list of SchedulerX2 workflows
---

# alibabacloudstack_schedulerx2_workflow

Query the list of SchedulerX2 workflows to obtain information about created workflows.

## Example Usage

```hcl

variable "name" {
  default = "tf-test-1"
}

resource "alibabacloudstack_schedulerx2_app_group" "example" {
  group_id        = "${var.name}.terra"
  app_name        = var.name
  description     = var.name
  max_jobs        = 30
  max_concurrency = 10
  monitor_config {
    send_channel = "mail,ding"
    alarm_type   = "CustomContacts"
  }
  contacts {
    username    = "test"
    user_email  = "123@123.com"
    dingding_ak = "testakkkkk"
  }
  metrics_threshold {
    load5       = 10
    heap5_usage = 100
    disk_usage  = 100
  }
}

resource "alibabacloudstack_schedulerx2_workflow" "example" {
  group_id        = alibabacloudstack_schedulerx2_app_group.example.group_id
  name            = var.name
  description     = "Initial description"
  time_type       = "cron"
  time_zone       = "PRC"
  time_expression = "34 14 14 */1 * ?"
  max_concurrency = 1
}

data "alibabacloudstack_schedulerx2_workflows" "default" {
  ids = ["${alibabacloudstack_schedulerx2_workflow.example.id}"]
}

```

## Argument Reference

The following arguments support filtering the workflow list:

* `group_id` (Optional): The application group ID used to filter workflows under a specific application group.
* `ids` (Optional): A list of workflow IDs used to filter workflows with specified IDs.
* `name_regex` (Optional): A regular expression for names used to filter workflows with matching names.
* `namespace` (Optional): The namespace used to isolate resources in different environments. Default value is "system_namespace".

## Attributes Reference

The following attributes are exported:

* `id` (String): The unique identifier of the workflow, formatted as the workflow ID.
* `app_group_id` (Integer): The application group ID that identifies the application group to which the workflow belongs.
* `creator` (String): The account ID of the workflow creator.
* `description` (String): The description information of the workflow.
* `group_id` (String): The application group ID, corresponding to the group_id parameter.
* `max_concurrency` (Integer): The maximum number of concurrent executions for the workflow.
* `name` (String): The name of the workflow.
* `time_expression` (String): The time expression, such as a Cron expression.
* `time_type` (String): The time type, indicating the workflow trigger method.
* `updater` (String): The account ID of the last updater of the workflow.
* `workflow_id` (Integer): The workflow ID, a system-generated unique identifier.