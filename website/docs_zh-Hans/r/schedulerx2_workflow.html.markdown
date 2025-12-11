---
subcategory: "Schedulerx2"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_schedulerx2_workflow"
sidebar_current: "docs-Alibabacloudstack-schedulerx2-workflow"
description: |-
  管理Schedulerx2工作流资源
---

# alibabacloudstack_schedulerx2_workflow

管理Schedulerx2工作流资源，用于创建、配置和调度定时任务流程。

## 示例用法

### 基础用法

```hcl

variable "name" {
  default = "tf-test97924"
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
  monitor_config {
    send_channel = "mail,ding"
    alarm_type   = "CustomContacts"
  }
}



resource "alibabacloudstack_schedulerx2_workflow" "example" {
  time_expression = "34 14 14 */1 * ?"
  time_zone       = "PRC"
  max_concurrency = "1"
  group_id        = alibabacloudstack_schedulerx2_app_group.example.group_id
  name            = var.name
  description     = "Initial description"
  time_type       = "cron"
}
```

## 参数说明

支持以下参数：

* `description` - (必填) 工作流的描述信息。长度1-256字符，不能以`http://`或`https://`开头。
* `name` - (必填) 工作流的名称。长度1-128字符，不能以`http://`或`https://`开头。
* `time_type` - (必填) 时间触发类型。取值：`cron`（定时触发）或`api`（API触发）。
* `time_zone` - (必填) 时区配置。取值范围：`PRC`, `Hongkong`, `Japan`, `Singapore`, `GTM`, `GTM-0`至`GTM-14`（含正负时区）。
* `group_id` - (必填, 变更时重建) 工作流所属的组ID。创建后不可修改，变更时将重建资源。
* `enabled` - (可选) 是否启用工作流。设置为`true`时启用，`false`时禁用。默认值：`false`。
* `max_concurrency` - (可选) 最大并发执行任务数。取值范围：1-1000。
* `namespace` - (可选) 命名空间。默认值：`system_namespace`。
* `time_expression` - (可选) 时间表达式。当`time_type`为`cron`时生效，需符合Quartz cron表达式规范（如`"34 14 14 */1 * ?"`）。

## 属性说明

以下属性导出：

* `id` - 工作流ID（与`workflow_id`相同）。
* `workflow_id` - 工作流ID（系统自动生成的唯一标识）。