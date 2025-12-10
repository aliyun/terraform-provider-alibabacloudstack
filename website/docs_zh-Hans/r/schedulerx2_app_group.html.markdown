---
subcategory: "SchedulerX2"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_schedulerx2_app_group"
sidebar_current: "docs-Alibabacloudstack-schedulerx2-app_group"
description: |-
  管理SchedulerX2应用组
---

# alibabacloudstack_schedulerx2_app_group

管理SchedulerX2应用组，用于分布式任务调度系统中的应用管理。

## 示例用法

### 基础用法

```hcl

variable "name" {
  default = "tf-testAccschedulerx2AppGroup11340"
}


resource "alibabacloudstack_schedulerx2_app_group" "example" {
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

  group_id        = "${var.name}.terra"
  app_name        = var.name
  description     = var.name
  max_jobs        = "30"
  max_concurrency = "10"
  monitor_config {
    send_channel = "mail,ding"
    alarm_type   = "CustomContacts"
  }

}
```

## 参数说明

支持以下参数：

* `app_name` - (必填) 应用名称。长度为1-128个字符，不能以`http://`或`https://`开头。
* `group_id` - (必填, 变更时重建) 应用组唯一标识。格式为字符串，通常为应用的唯一命名空间标识。
* `description` - (可选) 应用描述信息。长度为1-256个字符，不能以`http://`或`https://`开头。
* `max_concurrency` - (可选) 最大并发数。表示该应用组能同时执行的最大任务数。
* `max_jobs` - (可选) 最大任务数。表示该应用组能同时运行的最大任务数量。
* `metrics_threshold` - (可选) 指标阈值配置，用于监控应用运行状态。
  * `disk_usage` - (可选) 磁盘使用率阈值，超过此值将触发告警。默认值为`95`。
  * `heap5_usage` - (可选) 堆内存使用率阈值，超过此值将触发告警。默认值为`90`。
  * `load5` - (可选) 5分钟系统负载阈值，超过此值将触发告警。默认值为`0`。
* `monitor_config` - (可选) 监控告警配置。
  * `alarm_type` - (可选) 告警类型。默认值为`CustomContacts`。
  * `send_channel` - (可选) 告警通知渠道。有效值：`ding`（钉钉）、`mail`（邮件）、`mail,ding`（邮件和钉钉）、空字符串（不通知）。
* `namespace` - (可选) 命名空间。默认值为`system_namespace`。
* `contacts` - (可选) 告警联系人列表。
  * `dingding_ak` - (可选) 钉钉AK，用于钉钉告警通知。
  * `user_email` - (可选) 用户邮箱，用于邮件告警通知。
  * `username` - (可选) 用户名称，用于标识联系人。

## 属性说明

以下属性会从API响应中导出：

* `id` - 应用组ID。
* `app_key` - 应用密钥，用于应用认证。