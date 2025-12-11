---
subcategory: "SchedulerX2"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_schedulerx2_workflow"
sidebar_current: "docs-Alibabacloudstack-datasource-schedulerx2-workflow"
description: |-
  查询SchedulerX2工作流列表
---

# alibabacloudstack_schedulerx2_workflow

查询SchedulerX2工作流列表，用于获取已创建的工作流信息。

## 示例用法

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

## 参数说明
以下参数支持过滤工作流列表：

* `group_id` (可选)：应用分组ID，用于过滤特定应用分组下的工作流。
* `ids` (可选)：工作流ID列表，用于过滤指定ID的工作流。
* `name_regex` (可选)：名称正则表达式，用于过滤名称匹配的工作流。
* `namespace` (可选)：命名空间，用于隔离不同环境的资源，默认值为"system_namespace"。

## 属性说明
以下属性被导出：

* `id` (字符串)：工作流的唯一标识符，格式为工作流ID。
* `app_group_id` (整数)：应用分组ID，标识工作流所属的应用分组。
* `creator` (字符串)：工作流创建者账号ID。
* `description` (字符串)：工作流描述信息。
* `group_id` (字符串)：应用分组ID，与参数中的group_id对应。
* `max_concurrency` (整数)：工作流最大并发执行数。
* `name` (字符串)：工作流名称。
* `time_expression` (字符串)：时间表达式，如Cron表达式。
* `time_type` (字符串)：时间类型，表示工作流触发方式。
* `updater` (字符串)：工作流最后更新者账号ID。
* `workflow_id` (整数)：工作流ID，系统生成的唯一标识。