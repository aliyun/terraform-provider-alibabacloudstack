---
subcategory: "Prometheus"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_prometheus_v2_alert"
sidebar_current: "docs-Alibabacloudstack-prometheus-prometheus_v2_alert"
description: |-
  管理Prometheus v2警告规则
---

# alibabacloudstack_prometheus_v2_alert

管理Prometheus v2警告规则。

## 示例用法

### 基础用法

```hcl

variable "name" {
  default = "tfacc14757"
}

resource "alibabacloudstack_prometheus_v2_instance" "default" {
  cluster_name = var.name
  tags         = ["test1", "test2"]
}

resource "alibabacloudstack_prometheus_v2_notify_group" "default" {
  name        = "${var.name}_notify_group"
  type        = "WEBHOOK"
  description = "${var.name}_notify_group_description"
  webhook_url = "https://test.com"
  webhook_header_params {
    key   = "Content-Type"
    value = "application/json"
  }
}


resource "alibabacloudstack_prometheus_v2_alert" "default" {
  notify_types = [
    "EMAIL",
    "SMS"
  ]
  trigger_period = "5m"
  notify_group_ids = [
    "${alibabacloudstack_prometheus_v2_notify_group.default.id}"
  ]
  notify_interval  = "10m"
  notify_recovered = "true"
  is_check_all     = "false"
  tag_set = [
    "aaa",
    "ccc"
  ]
  trigger_clusters = [
    "${alibabacloudstack_prometheus_v2_instance.default.id}"
  ]
  trigger_severity     = "warning"
  trigger_cron         = "0 /5 * * * ?"
  recover_notification = "Trigger condition\\$${alert_source} \\Hit record\\$${alert_time}"
  notification         = "Trigger condition: {condition}\nHit record :{alert_result}"
  name                 = var.name
  trigger_promql       = "select testfield from testtable where testfield >= 0"
}
```

## 参数说明

支持以下参数：

* `name` - (必填) 告警规则的名称。名称长度为1-128个字符，不能以`http://`或`https://`开头。
* `notification` - (必填) 触发告警时的通知消息模板。支持变量替换，如`${condition}`表示触发条件，`${alert_result}`表示命中记录。
* `recover_notification` - (必填) 告警恢复时的通知消息模板。当`notify_recovered`为true时必须设置。
* `trigger_period` - (必填) 检查周期，格式为数字加时间单位，如"5m"表示每5分钟检查一次。
* `trigger_promql` - (必填) PromQL查询语句，用于定义告警触发条件。
* `trigger_severity` - (必填) 告警级别，可选值为"warning"、"serious"、"fatal"。
* `is_check_all` - (可选) 是否检查所有集群。默认为false，当设置为true时，`trigger_clusters`参数将被忽略。
* `notify_group_ids` - (可选) 通知组ID列表，用于指定接收告警通知的联系组。
* `notify_interval` - (可选) 通知间隔，格式为数字加时间单位，如"10m"表示每10分钟发送一次通知。
* `notify_recovered` - (可选) 是否在告警恢复时发送通知。默认为false，当设置为true时，`recover_notification`参数必须设置。
* `notify_types` - (可选) 通知类型列表，可选值包括"EMAIL"、"SMS"等。
* `tag_set` - (可选) 标签集合，用于对告警规则进行分类和过滤。
* `trigger_clusters` - (可选) 触发规则的集群ID列表。当`is_check_all`为false时必须设置。
* `trigger_cron` - (可选) Cron表达式，用于定义告警触发的时间规则。当设置此参数时，`trigger_period`将被忽略。

## 属性说明

除了上述所有参数外，还导出了以下属性：

* `id` - 资源的唯一标识符。
* `alert_type` - 告警类型，固定为"PROMETHEUS"。
* `is_check_all` - 是否检查所有集群的实际值。
* `tag_set` - 标签集合的实际值。
* `trigger_clusters` - 触发规则的集群ID列表的实际值。
* `trigger_cron` - Cron表达式的实际值。
* `trigger_period` - 检查周期的实际值。
* `trigger_promql` - PromQL查询语句的实际值。
* `trigger_severity` - 告警级别的实际值。