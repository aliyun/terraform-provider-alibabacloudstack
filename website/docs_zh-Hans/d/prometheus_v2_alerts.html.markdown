---
subcategory: "Prometheus"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_prometheus_v2_alerts"
sidebar_current: "docs-alibabacloudstack-datasource-prometheus-v2-alerts"
description: |-
  提供 Prometheus V2 告警列表数据源
---

# alibabacloudstack\_prometheus\_v2\_alerts

该数据源提供可用的 Prometheus V2 告警列表。

## 典型用法

```terraform
variable "name" {
  default = "tfacc"
}

resource "alibabacloudstack_prometheus_v2_instance" "default" {
  cluster_name = "${var.name}"
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
  name                 = "${var.name}"
  notify_recovered 	   = true
  is_check_all 		   = false
  tag_set 			   = ["aaa", "ccc"]
  trigger_clusters     = ["${alibabacloudstack_prometheus_v2_instance.default.id}"]
  trigger_promql 	   = "select testfield from testtable where testfield >= 0"
  trigger_severity     = "warning"
  trigger_cron 		   = "0 /5 * * * ?"
  trigger_period       = "5m"
  recover_notification = "触发条件$${alert_source} \命中记录$${alert_time}"
  notification         = "触发条件: {condition}\n命中记录 :{alert_result}"
  notify_group_ids     = ["${alibabacloudstack_prometheus_v2_notify_group.default.id}"]
  notify_types         = ["EMAIL", "SMS"]
  notify_interval      = "10m"
}

data "alibabacloudstack_prometheus_v2_alerts" "default" {
  ids = ["${alibabacloudstack_prometheus_v2_alert.default.id}"]
}

output "first_alert_id" {
  value = data.alibabacloudstack_prometheus_v2_alerts.default.alerts.0.id
}
```

## 参数参考

以下参数可用于过滤告警：

* `ids` - (可选) 告警 ID 列表。
* `name_regex` - (可选) 用于按名称过滤告警的正则表达式。

## 属性参考

以下属性将会被导出：

* `alerts` - Prometheus V2 告警列表。每个元素包含以下属性：
  * `id` - 告警的 ID。
  * `name` - 告警的名称。
  * `notify_recovered` - 是否在告警恢复时发送通知。
  * `is_check_all` - 是否检查所有集群。
  * `trigger_clusters` - 触发告警的集群 ID 列表。
  * `trigger_promql` - 触发告警的 PromQL 查询语句。
  * `trigger_period` - 告警条件评估的时间段。
  * `trigger_severity` - 告警的严重级别。
  * `trigger_cron` - 调度告警的 cron 表达式。
  * `tag_set` - 与告警关联的标签集合。
  * `notify_types` - 发送的通知类型 (例如 EMAIL, SMS)。
  * `notify_group_ids` - 通知组的 ID。
  * `notify_interval` - 通知之间的间隔。
  * `notification` - 通知消息模板。
  * `recover_notification` - 恢复通知消息模板。
```