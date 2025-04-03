---
subcategory: "CloudMonitorService"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cloudmonitorservice_metricalarmrules"
sidebar_current: "docs-Alibabacloudstack-datasource-cloudmonitorservice-metricalarmrules"
description: |- 
  Provides a list of cloudmonitorservice metricalarmrules owned by an alibabacloudstack account.
---

# alibabacloudstack_cloudmonitorservice_metricalarmrules
-> **NOTE:** Alias name has: `alibabacloudstack_cms_alarms`

This data source provides a list of cloudmonitorservice metricalarmrules in an alibabacloudstack account according to the specified filters.

## Example Usage

```hcl
data "alibabacloudstack_cloudmonitorservice_metricalarmrules" "example" {
  name_regex = "TEST123456"
}

output "metricalarmrules_list" {
  value = data.alibabacloudstack_cloudmonitorservice_metricalarmrules.example.alarms
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to filter CloudMonitorService Metric Alarm Rules by Rule name.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `alarms` - The list of the Metric Alarm Rules. Each element contains the following attributes:
  * `group_name` - The name of the application group.
  * `metric_name` - The name of the metric associated with the alarm rule.
  * `no_effective_interval` - The time period during which the alert rule is ineffective.
  * `silence_time` - Notification silence period in the alarm state, in seconds. Valid value range: [300, 86400]. Default to 86400.
  * `contact_groups` - List of contact groups of the alarm rule, which must have been created on the console.
  * `mail_subject` - The subject of the alert notification email.
  * `source_type` - The type of the alert rule. Valid values:
    - `METRIC`: the alert rule for time series metrics.
    - `EVENT`: the alert rule for event-type metrics. This type was used in earlier versions and has been discarded.
  * `rule_id` - The ID of the alarm rule.
  * `period` - Index query cycle, which must be consistent with that defined for metrics. Default to 300, in seconds.
  * `dimensions` - Map of the resources associated with the alarm rule, such as "instanceId", "device" and "port". Each key's value is a string and it uses comma to split multiple items. For more information, see [Metrics Reference](https://www.alibabacloud.com/help/doc-detail/28619.htm).
  * `effective_interval` - The interval of effecting alarm rule. It format as "hh:mm-hh:mm", like "0:00-4:00". Default to "00:00-23:59".
  * `namespace` - The namespace of the monitored service.
  * `enable_state` - Indicates whether the alert rule was enabled. Valid values:
    - `true`: indicates that the alert rule was enabled.
    - `false`: indicates that the alert rule was disabled. This parameter is not specified in the request by default. In this case, both enabled and disabled rules are returned.
  * `webhook` - The callback URL.
  * `resources` - The resources associated with the alert rule.
  * `rule_name` - The alarm rule name.
  * `escalations` - The conditions for triggering different levels of alerts.
    * `critical_comparison_operator` - Critical level alarm comparison operator for critical alarm. Valid values: ["<=", "<", ">", ">=", "==", "!="]. Default to "==".
    * `critical_times` - Critical level alarm retry times. Default to 3.
    * `critical_statistics` - Critical level alarm statistics method for critical alarm. It must be consistent with that defined for metrics. Valid values: ["Average", "Minimum", "Maximum"]. Default to "Average".
    * `critical_threshold` - Critical level alarm threshold value for critical alarm, which must be a numeric value currently.
    * `info_comparison_operator` - Info level alarm comparison operator for info alarm. Valid values: ["<=", "<", ">", ">=", "==", "!="]. Default to "==".
    * `info_times` - Info level alarm retry times for info alarm. Default to 3.
    * `info_statistics` - Info level alarm statistics method for info alarm. It must be consistent with that defined for metrics. Valid values: ["Average", "Minimum", "Maximum"]. Default to "Average".
    * `info_threshold` - Info level alarm threshold value for info alarm, which must be a numeric value currently.
    * `warn_comparison_operator` - Warn level alarm comparison operator for warn alarm. Valid values: ["<=", "<", ">", ">=", "==", "!="]. Default to "==".
    * `warn_times` - Warn level alarm retry times for warn alarm. Default to 3.
    * `warn_statistics` - Warn level alarm statistics method for warn alarm. It must be consistent with that defined for metrics. Valid values: ["Average", "Minimum", "Maximum"]. Default to "Average".
    * `warn_threshold` - Warn level alarm threshold value for warn alarm, which must be a numeric value currently.