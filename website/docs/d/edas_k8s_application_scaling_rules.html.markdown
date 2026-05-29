---
subcategory: "Enterprise Distributed Application Service"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_edas_k8s_application_scaling_rules"
sidebar_current: "docs-alibabacloudstack-datasource-edas-k8s-application-scaling-rules"
description: |-
  Provides a list of EDAS K8s Application Scaling Rules to be used by the caller.
---

# alibabacloudstack_edas_k8s_application_scaling_rules

This data source provides a list of EDAS K8s Application Scaling Rules that match the given criteria.

## Example Usage

```terraform
data "alibabacloudstack_edas_k8s_application_scaling_rules" "example" {
  app_id = "your-app-id"
}

output "first_scaling_rule_name" {
  value = data.alibabacloudstack_edas_k8s_application_scaling_rules.example.scaling_rules.0.scaling_rule_name
}
```

## Argument Reference

The following arguments are supported:

* `app_id` - (Required, ForceNew) The ID of the application.
* `ids` - (Optional, ForceNew) A list of scaling rule IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by scaling rule name.
* `scaling_rule_type` - (Optional) The type of scaling rule. Valid values: `metric` and `trigger`.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of scaling rule IDs.
* `names` - A list of scaling rule names.
* `scaling_rules` - A list of scaling rules. Each element contains the following attributes:
  * `id` - The ID of the scaling rule.
  * `app_id` - The ID of the application.
  * `scaling_rule_name` - The name of the scaling rule.
  * `scaling_rule_type` - The type of the scaling rule.
  * `max_replicas` - The maximum number of replicas.
  * `min_replicas` - The minimum number of replicas.
  * `metrics` - The metrics configuration.
    * `type` - The type of metric.
    * `utilization` - The target average utilization of the metric.
  * `enabled` - Whether the scaling rule is enabled.
  * `triggers` - (Optional) The triggers configuration.
  * `type` - (Optional) The type of trigger. Default: `cron`.
  * `name` - (Optional) The name of the trigger.
  * `period` - (Optional) The period of the trigger. Valid values: `daily`, `weekly` and `monthly`.
  * `timer_in_day` - (Optional) The timer configuration in a day.
    * `at_time` - (Required) The scheduled time in the day, e.g. "08:00".
    * `replicas` - (Required) The number of replicas. Valid values: 1 to 100.
    * `horizon_mode` - (Optional) Whether to enable horizon mode. Default: `false`.
  * `timer_in_week` - (Optional) The timer configuration in a week. 
  * `timer_in_month` - (Optional) The timer configuration in a month. 
  * `scale_up_stabilization_window_seconds` - (Optional) The stabilization window for scale up in seconds. Valid values: 0 to 3600.
  * `scale_up_select_policy` - (Optional) The select policy for scale up. Valid values: `Min`, `Max`, `Disabled`. 
  * `scale_up_policies` - (Optional) The policies for scale up.
    * `type` - (Required) The type of policy. Valid values: `Percent`, `Pods`.
    * `value` - (Required) The value of policy. Valid values: 1 to 100.
    * `period_seconds` - (Required) The period in seconds. Valid values: 0 to 3600.

  * `scale_down_stabilization_window_seconds` - (Optional) The stabilization window for scale down in seconds. Valid values: 0 to 3600. 
  * `scale_down_select_policy` - (Optional) The select policy for scale down. Valid values: `Min`, `Max`, `Disabled`. 
  * `scale_down_policies` - (Optional) The policies for scale down.
    * `type` - (Required) The type of policy. Valid values: `Percent`, `Pods`.
    * `value` - (Required) The value of policy. Valid values: 1 to 100.
    * `period_seconds` - (Required) The period in seconds. Valid values: 0 to 3600.