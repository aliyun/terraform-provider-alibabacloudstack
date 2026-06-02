---
subcategory: "Managed Service for Prometheus"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_prometheus_v2_alert"
description: |-
  Manage Prometheus v2 alert rules
---

# alibabacloudstack_prometheus_v2_alert

Manage Prometheus v2 alert rules.

## Example Usage

### Basic Usage

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

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the alert rule. The name must be 1-128 characters in length and cannot start with `http://` or `https://`.
* `notification` - (Required) The notification message template when an alert is triggered. Supports variable replacement, such as `${condition}` represents the trigger condition, and `${alert_result}` represents the hit record.
* `recover_notification` - (Required) The notification message template when an alert recovers. Must be set when `notify_recovered` is true.
* `trigger_period` - (Required) The check period, in the format of a number plus a time unit, such as "5m" means checking every 5 minutes.
* `trigger_promql` - (Required) The PromQL query statement used to define the alert trigger condition.
* `trigger_severity` - (Required) The alert severity level. Valid values are "warning", "serious", and "fatal".
* `is_check_all` - (Optional) Whether to check all clusters. Default is false. When set to true, the `trigger_clusters` parameter will be ignored.
* `notify_group_ids` - (Optional) A list of notification group IDs used to specify contact groups that receive alert notifications.
* `notify_interval` - (Optional) The notification interval, in the format of a number plus a time unit, such as "10m" means sending a notification every 10 minutes.
* `notify_recovered` - (Optional) Whether to send a notification when an alert recovers. Default is false. When set to true, the `recover_notification` parameter must be set.
* `notify_types` - (Optional) A list of notification types. Valid values include "EMAIL", "SMS", etc.
* `tag_set` - (Optional) A set of tags used to categorize and filter alert rules.
* `trigger_clusters` - (Optional) A list of cluster IDs that trigger the rule. Must be set when `is_check_all` is false.
* `trigger_cron` - (Optional) A Cron expression used to define the time rule for alert triggering. When this parameter is set, `trigger_period` will be ignored.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the resource.
* `alert_type` - The type of the alert, fixed as "PROMETHEUS".
* `is_check_all` - The actual value of whether to check all clusters.
* `tag_set` - The actual value of the tag set.
* `trigger_clusters` - The actual value of the list of cluster IDs that trigger the rule.
* `trigger_cron` - The actual value of the Cron expression.
* `trigger_period` - The actual value of the check period.
* `trigger_promql` - The actual value of the PromQL query statement.
* `trigger_severity` - The actual value of the alert severity level.