---
subcategory: "Alibaba Cloud DNS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dns_gtm_access_strategies"
sidebar_current: "docs-Alibabacloudstack-datasource-dns-gtm-access-strategies"
description: |-
  Provides a list of DNS GTM Access Strategies.
---

# alibabacloudstack_dns_gtm_access_strategies

This data source provides the DNS GTM Access Strategies available in ApsaraStack.

-> **NOTE:** Available in ApsaraStack.

## Example Usage

```hcl
data "alibabacloudstack_dns_gtm_access_strategies" "example" {
  gtm_instance_id = "gtm-xxx"
  name_regex      = "^test-.*"
}

output "strategies" {
  value = data.alibabacloudstack_dns_gtm_access_strategies.example.strategies
}
```

## Argument Reference

The following arguments are supported:

* `gtm_instance_id` - (Required) The ID of the GTM instance.
* `name_regex` - (Optional) A regex string to filter strategies by name.
* `ids` - (Optional, Computed) A list of strategy IDs to filter results.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of strategy IDs.
* `strategies` - A list of DNS GTM Access Strategies. Each element contains the following attributes:
  * `id` - The ID of the strategy, in the format `{gtm_instance_id}:{strategy_id}`.
  * `name` - The name of the strategy.
  * `gtm_instance_id` - The ID of the GTM instance.
  * `default_gtm_address_pool_type` - The type of the default GTM address pool.
  * `default_gtm_address_pool_id` - The ID of the default GTM address pool.
  * `default_gtm_address_pool_name` - The name of the default GTM address pool.
  * `default_min_available_addr_num` - The minimum number of available addresses for the default pool.
  * `default_available_addr_num` - The number of available addresses for the default pool.
  * `failover_gtm_address_pool_id` - The ID of the failover GTM address pool.
  * `failover_gtm_address_pool_name` - The name of the failover GTM address pool.
  * `failover_gtm_address_pool_type` - The type of the failover GTM address pool.
  * `failover_min_available_addr_num` - The minimum number of available addresses for the failover pool.
  * `failover_available_addr_num` - The number of available addresses for the failover pool.
  * `specified_gtm_address_pool` - The specified GTM address pool.
  * `switch_mode` - The switch mode.
  * `in_use_gtm_address_pool_id` - The ID of the currently used GTM address pool.
  * `in_use_gtm_address_pool_name` - The name of the currently used GTM address pool.
  * `line_ids` - A list of line IDs.
