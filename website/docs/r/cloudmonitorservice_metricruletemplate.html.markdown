---
subcategory: "CloudMonitorService"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cloudmonitorservice_metricruletemplate"
sidebar_current: "docs-Alibabacloudstack-cloudmonitorservice-metricruletemplate"
description: |- 
  Provides a cloudmonitorservice Metricruletemplate resource.
---

# alibabacloudstack_cloudmonitorservice_metricruletemplate
-> **NOTE:** Alias name has: `alibabacloudstack_cms_metric_rule_template`

Provides a cloudmonitorservice Metricruletemplate resource.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-testacccloudmonitorservicemetricruletemplate38509"
}

resource "alibabacloudstack_cms_monitor_group" "default" {
  monitor_group_name = var.name
}

resource "alibabacloudstack_cms_metric_rule_template" "default" {
  metric_rule_template_name = var.name
  description              = var.name
  enable                   = true
  notify_level             = 4
  silence_time             = 86400
  group_id                 = alibabacloudstack_cms_monitor_group.default.id

  alert_templates {
    category    = "ecs"
    metric_name = "cpu_total"
    namespace   = "acs_ecs_dashboard"
    rule_name   = "tf_testAcc_new"

    escalations {
      critical {
        comparison_operator = "GreaterThanThreshold"
        statistics          = "Average"
        threshold           = "90"
        times               = "3"
      }
    }
  }

  enable_start_time = "0"
  enable_end_time   = "23"
  webhook           = "https://www.example.com"
}
```

## Argument Reference

The following arguments are supported:

* `alert_templates` - (Optional) The details of alert rules that are generated based on the alert template. 
  * `category` - (Required) The abbreviation of the service name. Valid values include: `ecs`, `rds`, `ads`, `slb`, `vpc`, `apigateway`, `cdn`, `cs`, `dcdn`, `ddos`, `eip`, `elasticsearch`, `emr`, `ess`, `hbase`, `iot_edge`, `kvstore_sharding`, `kvstore_splitrw`, `kvstore_standard`, `memcache`, `mns`, `mongodb`, `mongodb_cluster`, `mongodb_sharding`, `mq_topic`, `ocs`, `opensearch`, `oss`, `polardb`, `petadata`, `scdn`, `sharebandwidthpackages`, `sls`, `vpn`.
  * `escalations` - (Optional) The information about the trigger condition based on the alert level.
    * `critical` - (Optional) The condition for triggering critical-level alerts.
      * `comparison_operator` - (Optional) The comparison operator of the threshold for critical-level alerts. Valid values: `GreaterThanOrEqualToThreshold`, `GreaterThanThreshold`, `LessThanOrEqualToThreshold`, `LessThanThreshold`, `NotEqualToThreshold`, `GreaterThanYesterday`, `LessThanYesterday`, `GreaterThanLastWeek`, `LessThanLastWeek`, `GreaterThanLastPeriod`, `LessThanLastPeriod`.
      * `statistics` - (Optional) The statistical aggregation method for critical-level alerts.
      * `threshold` - (Optional) The threshold for critical-level alerts.
      * `times` - (Optional) The consecutive number of times for which the metric value is measured before a critical-level alert is triggered.
    * `info` - (Optional) The condition for triggering info-level alerts.
      * `comparison_operator` - (Optional) The comparison operator of the threshold for info-level alerts. Valid values: `GreaterThanOrEqualToThreshold`, `GreaterThanThreshold`, `LessThanOrEqualToThreshold`, `LessThanThreshold`, `NotEqualToThreshold`, `GreaterThanYesterday`, `LessThanYesterday`, `GreaterThanLastWeek`, `LessThanLastWeek`, `GreaterThanLastPeriod`, `LessThanLastPeriod`.
      * `statistics` - (Optional) The statistical aggregation method for info-level alerts.
      * `threshold` - (Optional) The threshold for info-level alerts.
      * `times` - (Optional) The consecutive number of times for which the metric value is measured before an info-level alert is triggered.
    * `warn` - (Optional) The condition for triggering warn-level alerts.
      * `comparison_operator` - (Optional) The comparison operator of the threshold for warn-level alerts. Valid values: `GreaterThanOrEqualToThreshold`, `GreaterThanThreshold`, `LessThanOrEqualToThreshold`, `LessThanThreshold`, `NotEqualToThreshold`, `GreaterThanYesterday`, `LessThanYesterday`, `GreaterThanLastWeek`, `LessThanLastWeek`, `GreaterThanLastPeriod`, `LessThanLastPeriod`.
      * `statistics` - (Optional) The statistical aggregation method for warn-level alerts.
      * `threshold` - (Optional) The threshold for warn-level alerts.
      * `times` - (Optional) The consecutive number of times for which the metric value is measured before a warn-level alert is triggered.
  * `metric_name` - (Required) The name of the metric. For more information, see [DescribeMetricMetaList](https://www.alibabacloud.com/help/doc-detail/98846.htm) or [Appendix 1: Metrics](https://www.alibabacloud.com/help/doc-detail/28619.htm).
  * `namespace` - (Required) The namespace of the service. For more information, see [DescribeMetricMetaList](https://www.alibabacloud.com/help/doc-detail/98846.htm) or [Appendix 1: Metrics](https://www.alibabacloud.com/help/doc-detail/28619.htm).
  * `rule_name` - (Required) The name of the alert rule.
  * `webhook` - (Optional) The callback URL to which a POST request is sent when an alert is triggered based on the alert rule.
* `apply_mode` - (Optional) Template application method. Valid values:
  * `GROUP_INSTANCE_FIRST`: The application group instance takes precedence. When an alarm template is applied, the application group instance takes precedence. If the instance does not exist in the application Group, the rules in the template are ignored.
  * `ALARM_TEMPLATE_FIRST`: The alarm template instance takes precedence. When an alarm template is applied, an alarm rule is created regardless of whether the instance exists in the application Group.
* `description` - (Optional) The description of the alert template.
* `enable` - (Optional) Whether to apply the alert template to the application group. Valid values: `true` or `false`. Default value: `false`.
* `enable_end_time` - (Optional) The end time when the alarm takes effect. Value range: `00`~`23`, indicating `00:59` to `23:59`.
* `enable_start_time` - (Optional) The start time when the alarm takes effect. Value range: `00`~`23`, indicating `00:00` to `23:00`.
* `group_id` - (Optional) The ID of the application group.
* `metric_rule_template_name` - (Required, ForceNew) The name of the alert template.
* `notify_level` - (Optional) Alarm notification mode. Valid values:
  * `2`: Phone + SMS + email + Wangwang + DingTalk robot.
  * `3`: SMS + email + Wangwang + DingTalk robot.
  * `4`: Wangwang + DingTalk Robot.
* `rest_version` - (Optional) The version of the alert template.
* `silence_time` - (Optional) Channel Silence Cycle. Unit: seconds. Default value: `86400`. When the monitoring data continuously exceeds the alarm rule threshold, only one alarm notification is sent in each silent period.
* `webhook` - (Optional) When an alarm occurs, the specified URL address is called back and a POST request is sent.
* `overwrite` - (Optional) Whether to overwrite the alert template. Valid values: `true` or `false`. Default value: `true`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The resource ID in Terraform of Metric Rule Template.
* `rest_version` - The version of the alert template.