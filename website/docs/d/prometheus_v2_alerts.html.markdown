---
subcategory: "Managed Service for Prometheus"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_prometheus_v2_alerts"
description: |-
  Provides a list of Prometheus V2 Alerts to the user.
---

# alibabacloudstack\_prometheus\_v2\_alerts

This data source provides a list of Prometheus V2 Alerts available to the user.

## Example Usage

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
  recover_notification = "Trigger condition$${alert_source} \Hit record$${alert_time}"
  notification         = "Trigger condition: {condition}\nHit record :{alert_result}"
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

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of Alert IDs.
* `name_regex` - (Optional) A regex string to filter alerts by name.

## Attributes Reference

The following attributes are exported:

* `alerts` - A list of Prometheus V2 Alerts. Each element contains the following attributes:
  * `id` - The ID of the alert.
  * `name` - The name of the alert.
  * `notify_recovered` - Whether to send notification when the alert recovers.
  * `is_check_all` - Whether the alert checks all clusters.
  * `trigger_clusters` - List of cluster IDs that trigger the alert.
  * `trigger_promql` - The PromQL query that triggers the alert.
  * `trigger_period` - The period during which the alert condition is evaluated.
  * `trigger_severity` - The severity level of the alert.
  * `trigger_cron` - The cron expression for scheduling the alert.
  * `tag_set` - A set of tags associated with the alert.
  * `notify_types` - Types of notifications to send (e.g., EMAIL, SMS).
  * `notify_group_ids` - IDs of notification groups.
  * `notify_interval` - Interval between notifications.
  * `notification` - The notification message template.
  * `recover_notification` - The recovery notification message template.
```