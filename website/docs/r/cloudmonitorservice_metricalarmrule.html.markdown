---
subcategory: "Cloud Monitor Service (CMS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cloudmonitorservice_metricalarmrule"
sidebar_current: "docs-Alibabacloudstack-cloudmonitorservice-metricalarmrule"
description: |- 
  Provides a cloudmonitorservice Metricalarmrule resource.
---

# alibabacloudstack_cloudmonitorservice_metricalarmrule
-> **NOTE:** Alias name has: `alibabacloudstack_cms_alarm`

Provides a cloudmonitorservice Metricalarmrule resource.

## Example Usage

Basic Usage

```hcl
variable "name" {
 default = "tf_testacc_cmsalarm4649479"
}

resource "alibabacloudstack_slb" "basic" {
 name          = "${var.name}"
}
resource "alibabacloudstack_cloudmonitorservice_metricalarmrule" "default" {
  rule_name    = "${var.name}"
  namespace = "acs_slb_dashboard"
  metric_name  = "ActiveConnection"
  dimensions = {
    instanceId = alibabacloudstack_slb.basic.id
  }
  escalations_critical {
    statistics = "Average"
    comparison_operator = "<="
    threshold = 35
    times = 2
  }
  enabled =      true
  contact_groups     = ["test-group"]
  effective_interval = "0:00-2:00"
}
```

## Argument Reference

The following arguments are supported:

* `rule_name` - (Required) The name of the alarm rule.
* `namespace` - (Required, ForceNew) The data namespace of the product is used to distinguish different products. For example, `acs_ecs_dashboard` for ECS and `acs_rds_dashboard` for RDS.
* `metric_name` - (Required, ForceNew) Monitoring item name. For example, `CPUUtilization` for ECS or `ActiveConnection` for SLB. Refer to [Metrics Reference](https://www.alibabacloud.com/help/doc-detail/28619.htm) for more details.
* `dimensions` - (Required, ForceNew) Map of the resources associated with the alarm rule, such as "instanceId", "device" and "port". Each key's value is a string and it uses comma to split multiple items. For more information, see [Metrics Reference](https://www.alibabacloud.com/help/doc-detail/28619.htm).
* `period` - (Optional) Statistical cycle in seconds. Valid values depend on the metrics. Default value is `300`.
* `escalations_critical` - (Optional) A configuration block defining the critical alarm settings (documented below). Only one level can be defined per alarm rule.
  * `statistics` - Critical level alarm statistics method. It must be consistent with that defined for metrics. Valid values: ["Average", "Minimum", "Maximum"]. Default to "Average".
  * `comparison_operator` - Critical level alarm comparison operator. Valid values: ["<=", "<", ">", ">=", "==", "!="]. Default to "==".
  * `threshold` - Critical level alarm threshold value, which must be a numeric value currently.
  * `times` - Critical level alarm retry times. Default to `3`.
* `escalations_warn` - (Optional) A configuration block defining the warning alarm settings (documented below). Only one level can be defined per alarm rule.
  * `statistics` - Warning level alarm statistics method. It must be consistent with that defined for metrics. Valid values: ["Average", "Minimum", "Maximum"]. Default to "Average".
  * `comparison_operator` - Warning level alarm comparison operator. Valid values: ["<=", "<", ">", ">=", "==", "!="]. Default to "==".
  * `threshold` - Warning level alarm threshold value, which must be a numeric value currently.
  * `times` - Warning level alarm retry times. Default to `3`.
* `escalations_info` - (Optional) A configuration block defining the informational alarm settings (documented below). Only one level can be defined per alarm rule.
  * `statistics` - Informational level alarm statistics method. It must be consistent with that defined for metrics. Valid values: ["Average", "Minimum", "Maximum"]. Default to "Average".
  * `comparison_operator` - Informational level alarm comparison operator. Valid values: ["<=", "<", ">", ">=", "==", "!="]. Default to "==".
  * `threshold` - Informational level alarm threshold value, which must be a numeric value currently.
  * `times` - Informational level alarm retry times. Default to `3`.
* `contact_groups` - (Required) List of contact groups of the alarm rule, which must have been created on the console. Alarm notifications are sent to the alarm contacts in the alarm contact group.
* `effective_interval` - (Optional) The interval during which the alarm rule is effective. It is formatted as "hh:mm-hh:mm", like "0:00-4:00". Default to "00:00-23:59".
* `silence_time` - (Optional) Notification silence period in the alarm state, in seconds. Valid value range: [300, 86400]. Default to `86400`.
* `enabled` - (Optional) Whether to enable the alarm rule. Default to `true`.
* `webhook` - (Optional) The webhook that should be called when the alarm is triggered. Currently, only HTTP protocol is supported. Default is an empty string.

-> **NOTE:** Each resource supports the creation of one of the following three levels: critical, warning, or informational.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `rule_id` - The unique ID of the alarm rule.
* `status` - The current status of the alarm rule (`true` means enabled, `false` means disabled).