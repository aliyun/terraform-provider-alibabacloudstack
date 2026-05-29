---
subcategory: "分布式任务调度 SchedulerX"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_schedulerx2_job"
sidebar_current: "docs-Alibabacloudstack-schedulerx2-job"
description: |-
  管理Schedulerx2任务
---
# alibabacloudstack_schedulerx2_job

管理阿里云Schedulerx2任务，用于创建、读取、更新和删除定时任务。

## 示例用法

### 基础用法

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

## 参数说明

支持以下参数：

* `group_id` - (必填, 变更时重建) 任务组ID。创建任务时必须指定任务所属的任务组。
* `job_type` - (必填, 变更时重建) 任务类型。支持的类型包括：python、java、shell等。
* `name` - (必填) 任务名称。任务的显示名称。
* `execute_mode` - (必填) 执行模式。可选值：standalone（单机运行）、broadcast（广播模式）、parallel（并行模式）、grid（网格模式）、batch（批量模式）、sharding（分片模式）。
* `time_type` - (必填) 时间类型。1表示Cron表达式，其他值请参考API文档。
* `time_expression` - (必填) 时间表达式。当time_type为1时，为Cron表达式，如"8 59 15 */1 * ?"。
* `content` - (必填) 任务内容。根据job_type的不同，可以是Python脚本、Java类名、Shell命令等。
* `description` - (可选) 任务描述。对任务的简要描述。
* `max_attempt` - (可选) 最大重试次数。任务执行失败后的最大重试次数，默认值为3。
* `max_concurrency` - (可选) 最大并发数。任务的最大并发执行实例数，默认值为1。
* `attempt_interval` - (可选) 重试间隔。任务执行失败后，重试前的等待时间（秒），默认值为30。
* `monitor_fail_enable` - (可选) 是否启用失败监控。如果任务执行失败，是否触发失败监控。
* `monitor_miss_worker_enable` - (可选) 是否启用Worker缺失监控。如果任务执行时没有可用的Worker，是否触发缺失监控。
* `monitor_timeout` - (可选) 超时时间。任务执行的超时时间（秒）。
* `monitor_timeout_enable` - (可选) 是否启用超时监控。如果任务执行时间超过设定的超时时间，则触发超时监控。
* `monitor_timeout_kill_enable` - (可选) 是否启用超时终止。如果任务执行时间超过设定的超时时间，是否终止任务。
* `namespace` - (可选) 命名空间。默认值为"system_namespace"。
* `parameters` - (可选) 任务参数。以key=value格式传递给任务的参数，多个参数用空格分隔。
* `priority` - (可选) 任务优先级。数值越大优先级越高，默认值为5。

## 属性说明

以下属性会从API响应中导出：

* `id` - 任务ID。
* `execute_mode` - 任务执行模式（小写格式）。从API返回的执行模式会转换为小写格式。