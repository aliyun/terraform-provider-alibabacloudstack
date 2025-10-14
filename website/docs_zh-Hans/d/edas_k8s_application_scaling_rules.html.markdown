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
* `scaling_rule_type` - (可选) 扩缩容规则类型。可选值: `metric` 和 `trigger`。

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
  * `trigger_type` - 触发器类型。
  * `trigger_name` - 触发器名称。
  * `trigger_period` - 触发器周期。
  * `trigger_dryrun` - 是否执行试运行。
  * `trigger_timer_in_day` - 每日定时器配置。
    * `at_time` - 当天的计划时间。
    * `replicas` - 副本数。
  * `trigger_timer_in_week` - 每周定时器配置。
  * `enabled` - 扩缩容规则是否启用。