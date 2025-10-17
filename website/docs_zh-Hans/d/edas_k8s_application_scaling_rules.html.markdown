---
subcategory: "EDAS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_edas_k8s_application_scaling_rules"
sidebar_current: "docs-alibabacloudstack-datasource-edas-k8s-application-scaling-rules"
description: |-
  提供EDAS K8s应用扩缩容规则列表供调用方使用
---

# alibabacloudstack_edas_k8s_application_scaling_rules

该数据源提供符合给定条件的EDAS K8s应用扩缩容规则列表。

## 示例用法

```terraform
data "alibabacloudstack_edas_k8s_application_scaling_rules" "example" {
  app_id = "your-app-id"
}

output "first_scaling_rule_name" {
  value = data.alibabacloudstack_edas_k8s_application_scaling_rules.example.scaling_rules.0.scaling_rule_name
}
```

## 参数说明

以下参数可用于查询：

* `app_id` - (必选, 不可变) 应用ID。
* `ids` - (可选, 不可变) 扩缩容规则ID列表。
* `name_regex` - (可选, 不可变) 用于根据扩缩容规则名称过滤结果的正则表达式。
* `scaling_rule_type` - 扩缩容规则类型。可选值: `metric` 和 `trigger`。

## 属性导出

以下属性会被导出：

* `ids` - 扩缩容规则ID列表。
* `names` - 扩缩容规则名称列表。
* `scaling_rules` - 扩缩容规则列表。每个元素包含以下属性：
  * `id` - 扩缩容规则ID。
  * `app_id` - 应用ID。
  * `scaling_rule_name` - 扩缩容规则名称。
  * `scaling_rule_type` - 扩缩容规则类型。
  * `max_replicas` - 最大副本数。
  * `min_replicas` - 最小副本数。
  * `metrics` - 指标配置。
    * `type` - 指标类型。
    * `utilization` - 指标的平均目标利用率。
  * `triggers` - 触发器配置。
  * `type` - 触发器类型。默认值: `cron`。
  * `name` - 触发器名称。
  * `period` - 触发器周期。可选值: `daily`、`weekly` 和 `monthly`。
  * `timer_in_day` - 每日定时器配置。
    * `at_time` - 当天的计划时间，例如 "08:00"。
    * `replicas` - 副本数。取值范围: 1 到 100。
    * `horizon_mode` - 是否启用水平模式。默认值: `false`。
  * `timer_in_week` - 每周定时器配置。未设置时为计算值。
  * `timer_in_month` - 每月定时器配置。未设置时为计算值。
  * `enabled` - 扩缩容规则是否启用。
  * `scale_up_stabilization_window_seconds` - 扩容稳定窗口时间（秒）。取值范围: 0 到 3600。未设置时为计算值。
  * `scale_up_select_policy` - 扩容选择策略。可选值: `Min`、`Max`、`Disabled`。未设置时为计算值。
  * `scale_up_policies` - 扩容策略。
    * `type` - 策略类型。可选值: `Percent`、`Pods`。
    * `value` - 策略值。取值范围: 1 到 100。
    * `period_seconds` - 周期时间（秒）。取值范围: 0 到 3600。

  * `scale_down_stabilization_window_seconds` - 缩容稳定窗口时间（秒）。取值范围: 0 到 3600。未设置时为计算值。
  * `scale_down_select_policy` - 缩容选择策略。可选值: `Min`、`Max`、`Disabled`。未设置时为计算值。
  * `scale_down_policies` - 缩容策略。
    * `type` - 策略类型。可选值: `Percent`、`Pods`。
    * `value` - 策略值。取值范围: 1 到 100。
    * `period_seconds` - 周期时间（秒）。取值范围: 0 到 3600。