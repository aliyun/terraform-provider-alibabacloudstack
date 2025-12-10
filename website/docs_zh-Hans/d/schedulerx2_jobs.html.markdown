---
subcategory: "SchedulerX2"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_schedulerx2_job"
sidebar_current: "docs-Alibabacloudstack-resource-schedulerx2-job"
description: |-
  管理SchedulerX2任务
---

# alibabacloudstack_schedulerx2_job

> schedulerx2 任务管理

## 示例用法

```hcl

variable "name" {
  default = "tf-test96911"
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

resource "alibabacloudstack_schedulerx2_job" "example" {
  namespace                   = "system_namespace"
  group_id                    = alibabacloudstack_schedulerx2_app_group.example.group_id
  name                        = var.name
  job_type                    = "python"
  execute_mode                = "standalone"
  description                 = var.name
  priority                    = 5
  parameters                  = "testargs=1"
  max_attempt                 = 3
  attempt_interval            = 30
  max_concurrency             = 2
  time_type                   = 1
  time_expression             = "8 59 15 */1 * ?"
  content                     = "python test"
  monitor_timeout_enable      = true
  monitor_timeout_kill_enable = true
  monitor_fail_enable         = true
  monitor_miss_worker_enable  = true
  monitor_timeout             = 3600
}

data "alibabacloudstack_schedulerx2_jobs" "default" {
  name_regex = alibabacloudstack_schedulerx2_job.example.name
}

```

## 参数说明
以下参数支持创建和管理SchedulerX2任务：

* `content` (字符串, 必填)：任务执行内容，例如"python test"或Java类路径。
* `execute_mode` (字符串, 必填)：任务执行模式，可选值包括"standalone"(单机运行)、"broadcast"(广播模式)等。
* `group_id` (字符串, 必填, 变更时重建)：任务所属的组ID，用于任务分组管理。
* `job_type` (字符串, 必填)：任务类型，可选值包括"python"、"java"、"shell"等。
* `name` (字符串, 必填, 变更时重建)：任务名称，用于标识任务。
* `priority` (整数, 必填)：任务优先级，数值越大优先级越高。
* `time_expression` (字符串, 必填)：时间表达式，根据time_type类型不同格式不同，如cron表达式"8 59 15 */1 * ?"。
* `time_type` (字符串, 必填)：时间类型，1表示固定频率，2表示cron表达式，3表示固定时间。
* `attempt_interval` (整数, 可选)：任务重试间隔时间(秒)，默认值30。
* `description` (字符串, 可选)：任务描述信息。
* `max_attempt` (整数, 可选)：最大重试次数，当任务执行失败时的重试次数。
* `max_concurrency` (整数, 可选)：最大并发数，控制任务同时执行的实例数量。
* `parameters` (字符串, 可选)：任务参数，格式为"key1=value1 key2=value2"。
* `xattrs` (字符串, 可选)：扩展属性，JSON格式字符串，用于存储额外配置。
* `monitor_config` (字符串, 可选)：监控配置，JSON格式字符串，包含超时、失败等监控设置。

## 属性说明
以下属性导出为资源属性：

* `id` (字符串)：任务ID，由系统生成的唯一标识符。
* `app_group_id` (整数)：应用组ID，关联的应用组标识。
* `calendar` (字符串)：日历配置，用于特殊日期调度。
* `clean_mode` (字符串)：清理模式，任务执行完成后的资源清理策略。
* `contact` (字符串)：联系人信息，任务异常时的通知联系人。
* `content_type` (整数)：内容类型，标识任务内容的格式。
* `creator` (字符串)：创建者ID，任务创建者的用户标识。
* `data_offset` (整数)：数据偏移量，用于数据分片处理。
* `gmt_create` (字符串)：任务创建时间，格式为ISO8601。
* `gmt_modified` (字符串)：任务最后修改时间，格式为ISO8601。
* `monitor_fail_enable` (布尔)：是否启用失败监控，当任务执行失败时触发告警。
* `monitor_miss_worker_enable` (布尔)：是否启用工作节点丢失监控，当工作节点离线时触发告警。
* `monitor_timeout` (整数)：超时时间(秒)，任务执行超过该时间视为超时。
* `monitor_timeout_enable` (布尔)：是否启用超时监控，当任务执行超时时触发告警。
* `monitor_timeout_kill_enable` (布尔)：是否启用超时终止，当任务执行超时时自动终止任务。
* `resource` (字符串)：资源标识，任务关联的资源信息。
* `template` (字符串)：模板信息，任务使用的模板配置。
* `timezone` (字符串)：时区设置，任务调度使用的时区。
* `updater` (字符串)：更新者ID，最后修改任务的用户标识。
* `version` (整数)：任务版本号，每次修改任务时递增。
* `workflow_id` (字符串)：工作流ID，关联的工作流标识。