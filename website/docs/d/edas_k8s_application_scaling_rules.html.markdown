---
subcategory: "EDAS"
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
  * `trigger_type` - The type of trigger.
  * `trigger_name` - The name of the trigger.
  * `trigger_period` - The period of the trigger.
  * `trigger_dryrun` - Whether to perform a dry run.
  * `trigger_timer_in_day` - The timer configuration in a day.
    * `at_time` - The scheduled time in the day.
    * `replicas` - The number of replicas.
  * `trigger_timer_in_week` - The timer configuration in a week.
  * `enabled` - Whether the scaling rule is enabled.