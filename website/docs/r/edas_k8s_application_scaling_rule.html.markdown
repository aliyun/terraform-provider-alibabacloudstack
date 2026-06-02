---
subcategory: "Enterprise Distributed Application Service (EDAS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_edas_k8s_application_scaling_rule"
sidebar_current: "docs-Alibabacloudstack-resource-edas-k8s-application-scaling-rule"
description: |-
  Provides an EDAS K8s Application Scaling Rule resource.
---

# alibabacloudstack_edas_k8s_application_scaling_rule

Provides an EDAS K8s application scaling rule resource. This allows scaling rules to be created, modified and deleted for automatic application scaling.

## Example Usage

### Metric-based Scaling Rule

```terraform
resource "alibabacloudstack_edas_k8s_application_scaling_rule" "example" {
  app_id             = "your-app-id"
  scaling_rule_name  = "example-scaling-rule"
  scaling_rule_type  = "metric"
  max_replicas       = 10
  min_replicas       = 2
  enabled            = true

  metrics {
    type        = "CPU"
    utilization = 80
  }

  metrics {
    type        = "MEMORY"
    utilization = 80
  }
}
```

### Trigger-based Scaling Rule

```terraform
resource "alibabacloudstack_edas_k8s_application_scaling_rule" "example_cron" {
  app_id             = "your-app-id"
  scaling_rule_name  = "example-cron-scaling-rule"
  scaling_rule_type  = "trigger"
  max_replicas       = 20
  min_replicas       = 5
  triggers {
    type   = "cron"
    name   = "example-trigger"
    period = "daily"
    
    timer_in_day {
      at_time  = "09:00"
      replicas = 10
    }

    timer_in_day {
      at_time  = "18:00"
      replicas = 5
    }
  }
  
  enabled = true
}
```

### Scaling Behavior Configuration

```terraform
resource "alibabacloudstack_edas_k8s_application_scaling_rule" "example_behavior" {
  app_id             = "your-app-id"
  scaling_rule_name  = "example-behavior-scaling-rule"
  scaling_rule_type  = "metric"
  max_replicas       = 10
  min_replicas       = 2
  enabled            = true

  metrics {
    type        = "CPU"
    utilization = 80
  }

  # Scale up behavior
  scale_up_stabilization_window_seconds = 300
  scale_up_select_policy                = "Max"
  
  scale_up_policies {
    type           = "Pods"
    value          = 2
    period_seconds = 60
  }

  scale_up_policies {
    type           = "Percent"
    value          = 10
    period_seconds = 60
  }

  # Scale down behavior
  scale_down_stabilization_window_seconds = 300
  scale_down_select_policy                = "Max"
  
  scale_down_policies {
    type           = "Pods"
    value          = 1
    period_seconds = 120
  }
}
```

## Argument Reference

The following arguments are supported:

* `app_id` - (Required, ForceNew) The ID of the application. Changing this forces a new resource to be created.
* `scaling_rule_name` - (Required, ForceNew) The name of the scaling rule. Changing this forces a new resource to be created.
* `scaling_rule_type` - (Required, ForceNew) The type of the scaling rule. Valid values: `metric` and `trigger`. Changing this forces a new resource to be created.
* `max_replicas` - (Optional) The maximum number of replicas. Valid values: 2 to 100.
* `min_replicas` - (Optional) The minimum number of replicas. Valid values: 1 to 100.
* `enabled` - (Optional) Whether to enable the scaling rule. Default: `false`.

### Metric-based Scaling Rule Arguments

When `scaling_rule_type` is set to `metric`, the following arguments are supported:

* `metrics` - (Optional) The metrics configuration.
  * `type` - (Required) The type of metric. Valid values: `CPU` and `MEMORY`.
  * `utilization` - (Required) The target average utilization of the metric. Valid values: 1 to 100.

### Trigger-based Scaling Rule Arguments

When `scaling_rule_type` is set to `trigger`, the following arguments are supported:

* `triggers` - (Optional) The triggers configuration.
  * `type` - (Optional) The type of trigger. Default: `cron`.
  * `name` - (Optional) The name of the trigger.
  * `period` - (Optional) The period of the trigger. Valid values: `daily`, `weekly` and `monthly`.
  * `timer_in_day` - (Optional) The timer configuration in a day.
    * `at_time` - (Required) The scheduled time in the day, e.g. "08:00".
    * `replicas` - (Required) The number of replicas. Valid values: 1 to 100.
    * `horizon_mode` - (Optional) Whether to enable horizon mode. Default: `false`.
  * `timer_in_week` - (Optional) The timer configuration in a week. Effective when `period` is set to `weekly`. Elements are day abbreviations like `Mon`, `Tue`, etc. Computed when not set.
  * `timer_in_month` - (Optional) The timer configuration in a month. Effective when `period` is set to `monthly`. Elements are day numbers like `1`, `2`, etc. Computed when not set.

### Scaling Behavior Arguments

The following arguments configure scaling behaviors:

* `scale_up_stabilization_window_seconds` - (Optional) The stabilization window for scale up in seconds. Valid values: 0 to 3600. Computed when not set.
* `scale_up_select_policy` - (Optional) The select policy for scale up. Valid values: `Min`, `Max`, `Disabled`. Computed when not set.
* `scale_up_policies` - (Optional) The policies for scale up.
  * `type` - (Required) The type of policy. Valid values: `Percent`, `Pods`.
  * `value` - (Required) The value of policy. Valid values: 1 to 100.
  * `period_seconds` - (Required) The period in seconds. Valid values: 0 to 3600.

* `scale_down_stabilization_window_seconds` - (Optional) The stabilization window for scale down in seconds. Valid values: 0 to 3600. Computed when not set.
* `scale_down_select_policy` - (Optional) The select policy for scale down. Valid values: `Min`, `Max`, `Disabled`. Computed when not set.
* `scale_down_policies` - (Optional) The policies for scale down.
  * `type` - (Required) The type of policy. Valid values: `Percent`, `Pods`.
  * `value` - (Required) The value of policy. Valid values: 1 to 100.
  * `period_seconds` - (Required) The period in seconds. Valid values: 0 to 3600.

-> **NOTE:** When setting scaling behavior arguments (`scale_up_*` or `scale_down_*`), all arguments in the respective group must be set together. For example, if you specify `scale_up_stabilization_window_seconds`, you must also specify `scale_up_select_policy` and `scale_up_policies`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID, formatted as `{app_id}:{scaling_rule_name}`.

## Import

EDAS K8s application scaling rule can be imported using the ID, e.g.

```bash
$ terraform import alibabacloudstack_edas_k8s_application_scaling_rule.example your-app-id:scaling-rule-name
```