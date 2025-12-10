---
subcategory: "Schedulerx2"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_schedulerx2_job"
sidebar_current: "docs-Alibabacloudstack-schedulerx2-job"
description: |-
  Manages Schedulerx2 jobs
---
# alibabacloudstack_schedulerx2_job

Manages Alibaba Cloud Schedulerx2 jobs for creating, reading, updating, and deleting scheduled tasks.

## Example Usage

### Basic Usage

```hcl

variable "name" {
  default = "tf-test45165"
}


resource "alibabacloudstack_schedulerx2_app_group" "example" {
  group_id        = "${var.name}.terra"
  app_name        = var.name
  description     = var.name
  max_jobs        = 30
  max_concurrency = 10
  metrics_threshold {
    load5       = 10
    heap5_usage = 100
    disk_usage  = 100
  }
}


resource "alibabacloudstack_schedulerx2_job" "example" {
  max_concurrency             = "2"
  monitor_timeout_kill_enable = "true"
  name                        = var.name
  max_attempt                 = "3"
  job_type                    = "java"
  description                 = "ddddddd"
  monitor_timeout             = "3600"
  attempt_interval            = "30"
  time_type                   = "1"
  content                     = "{\"className\":\"Create\"}"
  monitor_timeout_enable      = "true"
  parameters                  = "testargs=1"
  execute_mode                = "standalone"
  group_id                    = alibabacloudstack_schedulerx2_app_group.example.group_id
  monitor_fail_enable         = "true"
  time_expression             = "8 59 15 */1 * ?"
  monitor_miss_worker_enable  = "true"
  priority                    = "5"
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required, Forces new resource) The ID of the job group. When creating a job, you must specify the job group to which the job belongs.
* `job_type` - (Required, Forces new resource) The type of the job. Supported types include: python, java, shell, etc.
* `name` - (Required) The name of the job. The display name of the job.
* `execute_mode` - (Required) The execution mode. Valid values: standalone (single machine), broadcast (broadcast mode), parallel (parallel mode), grid (grid mode), batch (batch mode), sharding (sharding mode).
* `time_type` - (Required) The time type. 1 indicates a Cron expression. For other values, see the API documentation.
* `time_expression` - (Required) The time expression. When `time_type` is 1, it is a Cron expression, such as "8 59 15 */1 * ?".
* `content` - (Required) The content of the job. Depending on the `job_type`, it can be a Python script, Java class name, Shell command, etc.
* `description` - (Optional) The description of the job. A brief description of the job.
* `max_attempt` - (Optional) The maximum number of retries. The maximum number of retries after a job execution failure. Default value: 3.
* `max_concurrency` - (Optional) The maximum concurrency. The maximum number of concurrent execution instances for the job. Default value: 1.
* `attempt_interval` - (Optional) The retry interval. The waiting time (in seconds) before retrying after a job execution failure. Default value: 30.
* `monitor_fail_enable` - (Optional) Whether to enable failure monitoring. If the job execution fails, failure monitoring is triggered.
* `monitor_miss_worker_enable` - (Optional) Whether to enable worker missing monitoring. If no available worker exists when the job is executed, missing monitoring is triggered.
* `monitor_timeout` - (Optional) The timeout period. The timeout period (in seconds) for job execution.
* `monitor_timeout_enable` - (Optional) Whether to enable timeout monitoring. If the job execution time exceeds the set timeout period, timeout monitoring is triggered.
* `monitor_timeout_kill_enable` - (Optional) Whether to enable timeout termination. If the job execution time exceeds the set timeout period, whether to terminate the job.
* `namespace` - (Optional) The namespace. Default value: "system_namespace".
* `parameters` - (Optional) The job parameters. Parameters passed to the job in key=value format, with multiple parameters separated by spaces.
* `priority` - (Optional) The job priority. A higher value indicates a higher priority. Default value: 5.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the job.
* `execute_mode` - The execution mode of the job (in lowercase). The execution mode returned by the API is converted to lowercase.