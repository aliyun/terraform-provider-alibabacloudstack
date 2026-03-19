---
subcategory: "Enterprise Distributed Application Service (EDAS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_edas_k8s_application_scaling_rule"
sidebar_current: "docs-alibabacloudstack-resource-edas-k8s-application-scaling-rule"
description: |-
  提供一个EDAS K8s应用扩缩容规则资源
---

# alibabacloudstack_edas_k8s_application_scaling_rule

提供EDAS K8s应用扩缩容规则资源。该资源允许创建、修改和删除用于应用自动扩缩容的规则。

## 示例用法

### 基于指标的扩缩容规则

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

### 基于触发器的扩缩容规则

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

### 扩缩容行为配置

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

  # 扩容行为配置
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

  # 缩容行为配置
  scale_down_stabilization_window_seconds = 300
  scale_down_select_policy                = "Max"
  
  scale_down_policies {
    type           = "Pods"
    value          = 1
    period_seconds = 120
  }
}
```

## 参数说明

以下参数可用于配置资源:

* `app_id` - (必选, 不可变) 应用ID。更改此参数将创建新资源。
* `scaling_rule_name` - (必选, 不可变) 扩缩容规则名称。更改此参数将创建新资源。
* `scaling_rule_type` - (必选, 不可变) 扩缩容规则类型。可选值: `metric` 和 `trigger`。更改此参数将创建新资源。
* `max_replicas` - (可选) 最大副本数。取值范围: 2 到 100。
* `min_replicas` - (可选) 最小副本数。取值范围: 1 到 100。
* `enabled` - (可选) 是否启用扩缩容规则。默认值: `false`。

### 基于指标的扩缩容规则参数

当 `scaling_rule_type` 设置为 `metric` 时，支持以下参数:

* `metrics` - (可选) 指标配置。
  * `type` - (必选) 指标类型。可选值: `CPU` 和 `MEMORY`。
  * `utilization` - (必选) 指标的平均目标利用率。取值范围: 1 到 100。

### 基于触发器的扩缩容规则参数

当 `scaling_rule_type` 设置为 `trigger` 时，支持以下参数:

* `triggers` - (可选) 触发器配置。
  * `type` - (可选) 触发器类型。默认值: `cron`。
  * `name` - (可选) 触发器名称。
  * `period` - (可选) 触发器周期。可选值: `daily`、`weekly` 和 `monthly`。
  * `timer_in_day` - (可选) 每日定时器配置。
    * `at_time` - (必选) 当天的计划时间，例如 "08:00"。
    * `replicas` - (必选) 副本数。取值范围: 1 到 100。
    * `horizon_mode` - (可选) 是否启用水平模式。默认值: `false`。
  * `timer_in_week` - (可选) 每周定时器配置。未设置时为计算值。
  * `timer_in_month` - (可选) 每月定时器配置。未设置时为计算值。

### 扩缩容行为参数

以下参数用于配置扩缩容行为：

* `scale_up_stabilization_window_seconds` - (可选) 扩容稳定窗口时间（秒）。取值范围: 0 到 3600。未设置时为计算值。
* `scale_up_select_policy` - (可选) 扩容选择策略。可选值: `Min`、`Max`、`Disabled`。未设置时为计算值。
* `scale_up_policies` - (可选) 扩容策略。
  * `type` - (必选) 策略类型。可选值: `Percent`、`Pods`。
  * `value` - (必选) 策略值。取值范围: 1 到 100。
  * `period_seconds` - (必选) 周期时间（秒）。取值范围: 0 到 3600。

* `scale_down_stabilization_window_seconds` - (可选) 缩容稳定窗口时间（秒）。取值范围: 0 到 3600。未设置时为计算值。
* `scale_down_select_policy` - (可选) 缩容选择策略。可选值: `Min`、`Max`、`Disabled`。未设置时为计算值。
* `scale_down_policies` - (可选) 缩容策略。
  * `type` - (必选) 策略类型。可选值: `Percent`、`Pods`。
  * `value` - (必选) 策略值。取值范围: 1 到 100。
  * `period_seconds` - (必选) 周期时间（秒）。取值范围: 0 到 3600。

-> **注意：** 在设置扩缩容行为参数时（`scale_up_*` 或 `scale_down_*`），必须同时设置相应组中的所有参数。例如，如果指定了 `scale_up_stabilization_window_seconds`，则还必须指定 `scale_up_select_policy` 和 `scale_up_policies`。

## 属性导出

以下属性会被导出:

* `id` - 资源ID，格式为 `{app_id}:{scaling_rule_name}`。

## 导入

EDAS K8s应用扩缩容规则可以通过ID导入，例如:

```bash
$ terraform import alibabacloudstack_edas_k8s_application_scaling_rule.example your-app-id:scaling-rule-name
```