---
subcategory: "云解析 DNS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dns_gtm_access_strategies"
sidebar_current: "docs-Alibabacloudstack-datasource-dns-gtm-access-strategies"
description: |-
  提供DNS GTM访问策略列表。
---

# alibabacloudstack_dns_gtm_access_strategies

该数据源用于获取专有云中可用的DNS GTM访问策略列表。

-> **注意:** 适用于专有云环境。

## 示例

```hcl
data "alibabacloudstack_dns_gtm_access_strategies" "example" {
  gtm_instance_id = "gtm-xxx"
  name_regex      = "^test-.*"
}

output "strategies" {
  value = data.alibabacloudstack_dns_gtm_access_strategies.example.strategies
}
```

## 参数说明

以下参数支持配置：

* `gtm_instance_id` - (必选) GTM实例的ID。
* `name_regex` - (可选) 用于按名称过滤策略的正则表达式字符串。
* `ids` - (可选, Computed) 用于过滤结果的策略ID列表。

## 属性参考

以下属性会被导出：

* `ids` - 策略ID列表。
* `strategies` - DNS GTM访问策略列表。每个元素包含以下属性：
  * `id` - 策略的ID，格式为 `{gtm_instance_id}:{strategy_id}`。
  * `name` - 策略的名称。
  * `gtm_instance_id` - GTM实例的ID。
  * `default_gtm_address_pool_type` - 默认GTM地址池的类型。
  * `default_gtm_address_pool_id` - 默认GTM地址池的ID。
  * `default_gtm_address_pool_name` - 默认GTM地址池的名称。
  * `default_min_available_addr_num` - 默认地址池的最小可用地址数。
  * `default_available_addr_num` - 默认地址池的可用地址数。
  * `failover_gtm_address_pool_id` - 故障转移GTM地址池的ID。
  * `failover_gtm_address_pool_name` - 故障转移GTM地址池的名称。
  * `failover_gtm_address_pool_type` - 故障转移GTM地址池的类型。
  * `failover_min_available_addr_num` - 故障转移地址池的最小可用地址数。
  * `failover_available_addr_num` - 故障转移地址池的可用地址数。
  * `specified_gtm_address_pool` - 指定的GTM地址池。
  * `switch_mode` - 切换模式。
  * `in_use_gtm_address_pool_id` - 当前使用的GTM地址池的ID。
  * `in_use_gtm_address_pool_name` - 当前使用的GTM地址池的名称。
  * `line_ids` - 线路ID列表。
