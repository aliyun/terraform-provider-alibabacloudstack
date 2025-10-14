---
subcategory: "EDAS"
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
  trigger_type       = "cron"
  trigger_name       = "example-trigger"
  trigger_period     = "daily"
  trigger_dryrun     = false
  enabled            = true

  trigger_timer_in_day {
    at_time  = "09:00"
    replicas = 10
  }

  trigger_timer_in_day {
    at_time  = "18:00"
    replicas = 5
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

* `trigger_type` - (可选) 触发器类型。默认值: `cron`。
* `trigger_name` - (可选) 触发器名称。
* `trigger_period` - (可选) 触发器周期。可选值: `daily` 和 `weekly`。
* `trigger_dryrun` - (可选) 是否执行试运行。
* `trigger_timer_in_day` - (可选) 每日定时器配置。
  * `at_time` - (必选) 当天的计划时间，例如 "08:00"。
  * `replicas` - (必选) 副本数。取值范围: 1 到 100。
* `trigger_timer_in_week` - (可选) 每周定时器配置。

## 属性导出

以下属性会被导出:

* `id` - 资源ID，格式为 `{app_id}:{scaling_rule_name}`。

## 导入

EDAS K8s应用扩缩容规则可以通过ID导入，例如:

```bash
$ terraform import alibabacloudstack_edas_k8s_application_scaling_rule.example your-app-id:scaling-rule-name
```