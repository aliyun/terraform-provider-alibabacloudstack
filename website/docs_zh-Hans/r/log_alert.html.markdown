---
subcategory: "Log Service (SLS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_log_alert"
sidebar_current: "docs-alibabacloudstack-resource-log-alert"
description: |-
  编排日志告警
---

# alibabacloudstack_log_alert

使用Provider配置的凭证在指定的资源集编排日志告警。  
日志告警是日志服务的一个单元，用于监控和告警用户logstore的状态信息。  
日志服务允许您基于仪表盘中的图表配置告警，实时监控服务状态。

## 示例用法

### 基础用法

```
resource "alibabacloudstack_log_project" "example" {
  name        = "test-tf"
  description = "create by terraform"
}

resource "alibabacloudstack_log_store" "example" {
  project               = alibabacloudstack_log_project.example.name
  name                  = "tf-test-logstore"
  retention_period      = 3650
  shard_count           = 3
  auto_split            = true
  max_split_shard_count = 60
  append_meta           = true
}

resource "alibabacloudstack_log_alert" "example" {
  project_name      = alibabacloudstack_log_project.example.name
  alert_name        = "tf-test-alert"
  alert_displayname = "tf-test-alert-displayname"
  condition         = "count> 100"
  dashboard         = "tf-test-dashboard"
  query_list {
    logstore    = "tf-test-logstore"
    chart_title = "chart_title"
    start       = "-60s"
    end         = "20s"
    query       = "* AND aliyun"
  }
  notification_list {
    type        = "SMS"
    mobile_list = ["12345678", "87654321"]
    content     = "alert content"
  }
  notification_list {
    type       = "Email"
    email_list = ["aliyun@alibaba-inc.com", "tf-test@123.com"]
    content    = "alert content"
  }
  notification_list {
    type        = "DingTalk"
    service_uri = "www.aliyun.com"
    content     = "alert content"
  }
}
```

## 参数说明

以下是支持的参数：

* `project_name` - (必填，变更时重建) 项目名称。
* `alert_name` - (必填，变更时重建) 配置告警服务的日志存储名称。
* `alert_displayname` - (必填) 告警显示名称。
* `alert_description` - (可选) 告警描述。
* `condition` - (必填) 条件表达式，例如：`count > 100`。
* `dashboard` - (必填) 与告警关联的仪表盘名称。如果不存在这样的仪表盘，Terraform 将帮助您创建一个空的仪表盘。
* `mute_until` - (可选) 时间戳，在此之前关闭通知。
* `throttling` - (可选) 通知间隔，默认为无间隔。支持数字+单位类型，例如 `60s`、`1h`。
* `notify_threshold` - (可选) 通知阈值，达到触发次数后才进行通知，默认为 `1`。
* `query_list` - (必填) 配置告警查询的多个条件。
  * `chart_title` - (必填) 图表标题。
  * `logstore` - (必填) 查询日志存储。
  * `query` - (必填) 与图表对应的查询。示例：`* AND aliyun`。
  * `start` - (必填) 开始时间。示例：`-60s`。
  * `end` - (必填) 结束时间。示例：`20s`。
  * `time_span_type` - (可选) 默认为 `Custom`。无需配置此参数。
* `notification_list` - (必填) 告警通知列表。
  * `type` - (必填) 通知类型。支持 `Email`、`SMS`、`DingTalk`、`MessageCenter`。
  * `content` - (必填) 告警通知内容。
  * `service_uri` - (可选) 请求地址。
  * `mobile_list` - (可选) 短信发送的手机号码。
  * `email_list` - (可选) 电子邮件地址列表。
* `schedule_interval` - (可选) 执行间隔。最小为 `60` 秒，例如 `60s`、`1h`。
* `schedule_type` - (可选) 默认为 `FixedRate`。无需配置此参数。
* `mute_until` - (可选) 时间戳，在此之前关闭通知。此属性允许大于或等于 `0` 的值。

## 属性说明

以下属性将被导出：

* `id` - 日志告警的ID。格式为 `<project>:<alert_name>`。
* `mute_until` - 时间戳，在此之前关闭通知。

## 导入

可以使用ID导入日志告警，例如：

```bash
$ terraform import alibabacloudstack_log_alert.example tf-log:tf-log-alert
```