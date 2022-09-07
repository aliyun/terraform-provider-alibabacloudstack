---
subcategory: "Cloud Monitor"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cms_alarm"
sidebar_current: "docs-alibabacloudstack-resource-cms-alarm"
description: |-
  Provides a resource to build a alarm rule for cloud monitor.
---

# alibabacloudstack\_cms\_alarm

This resource provides a alarm rule resource and it can be used to monitor several cloud services according different metrics.

## Example Usage

Basic Usage

``` 
resource "alibabacloudstack_cms_alarm" "basic1" {
  name    = "tf-testAccCmsAlarm_basic"
  project = "acs_slb_dashboard"
  metric  = "ActiveConnection"
  dimensions = {
    instanceId = "lb-z0y7dpol"
  }
  escalations_critical {
    statistics = "Average"
    comparison_operator = "<="
    threshold = 35
    times = 2
  }
  period  =    300
  enabled =      true
  contact_groups     = ["test-group"]
  effective_interval = "0:00-2:00"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The alarm rule name.
* `project` - (Required, ForceNew) Monitor project name, such as "acs_ecs_dashboard" and "acs_rds_dashboard". For more information, see [Metrics Reference](https://www.alibabacloud.com/help/doc-detail/28619.htm).
* `metric` - (Required, ForceNew) Name of the monitoring metrics corresponding to a project, such as "CPUUtilization" and "networkin_rate". For more information, see [Metrics Reference](https://www.alibabacloud.com/help/doc-detail/28619.htm).
* `dimensions` - (Required, ForceNew) Map of the resources associated with the alarm rule, such as "instanceId", "device" and "port". Each key's value is a string and it uses comma to split multiple items. For more information, see [Metrics Reference](https://www.alibabacloud.com/help/doc-detail/28619.htm).
* `period` - Index query cycle, which must be consistent with that defined for metrics. 
* `escalations_critical` - (Optional) A configuration of critical alarm (documented below).
* `escalations_warn` - (Optional) A configuration of critical warn (documented below).
* `escalations_info` - (Optional) A configuration of critical info (documented below).
* `contact_groups` - (Required) List contact groups of the alarm rule, which must have been created on the console.
* `effective_interval` - (Available) The interval of effecting alarm rule. It foramt as "hh:mm-hh:mm", like "0:00-4:00". Default to "00:00-23:59".
* `silence_time` - Notification silence period in the alarm state, in seconds. Valid value range: [300, 86400]. Default to 86400
* `notify_type` - Notification type. Valid value [0, 1]. The value 0 indicates TradeManager+email, and the value 1 indicates that TradeManager+email+SMS
* `enabled` - Whether to enable alarm rule. Default to true.
* `webhook`- (Optional) The webhook that should be called when the alarm is triggered. Currently, only http protocol is supported. Default is empty string.

-> **NOTE:** Each resource supports the creation of one of the following three levels.

#### Block escalations critical alarm

The escalations_critical supports the following:

* `statistics` - Critical level alarm statistics method.. It must be consistent with that defined for metrics. Valid values: ["Average", "Minimum", "Maximum"]. Default to "Average".
* `comparison_operator` - Critical level alarm comparison operator. Valid values: ["<=", "<", ">", ">=", "==", "!="]. Default to "==".
* `threshold` - Critical level alarm threshold value, which must be a numeric value currently.
* `times` - Critical level alarm retry times. Default to 3.

#### Block escalations warn alarm

The escalations_warn supports the following:

* `statistics` - Critical level alarm statistics method.. It must be consistent with that defined for metrics. Valid values: ["Average", "Minimum", "Maximum"]. Default to "Average".
* `comparison_operator` - Critical level alarm comparison operator. Valid values: ["<=", "<", ">", ">=", "==", "!="]. Default to "==".
* `threshold` - Critical level alarm threshold value, which must be a numeric value currently.
* `times` - Critical level alarm retry times. Default to 3.

#### Block escalations info alarm

The escalations_info supports the following:

* `statistics` - Critical level alarm statistics method.. It must be consistent with that defined for metrics. Valid values: ["Average", "Minimum", "Maximum"]. Default to "Average".
* `comparison_operator` - Critical level alarm comparison operator. Valid values: ["<=", "<", ">", ">=", "==", "!="]. Default to "==".
* `threshold` - Critical level alarm threshold value, which must be a numeric value currently.
* `times` - Critical level alarm retry times. Default to 3.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the alarm rule.
* `status` - The current alarm rule status.

## Import

Alarm rule can be imported using the id, e.g.

```
$ terraform import alibabacloudstack_cms_alarm.alarm abc12345
```
