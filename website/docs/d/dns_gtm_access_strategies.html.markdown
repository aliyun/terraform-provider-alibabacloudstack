---
subcategory: "Universal DNS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dns_gtm_access_strategy"
sidebar_current: "docs-alibabacloudstack-datasource-dns_gtm_access_strategy"
description: |-
  Provides access strategies for Global Traffic Management (GTM) instances in Alibaba Cloud DNS.
---

# alibabacloudstack_dns_gtm_access_strategy

This data source queries access strategies for Global Traffic Management (GTM) instances in Alibaba Cloud DNS.

## Example Usage

```hcl
variable "name" {
  default = "tfacc3910"
}

resource "alibabacloudstack_dns_private_domain" "default" {
  name = "${var.name}.local."
}

resource "alibabacloudstack_dns_gtm_instance" "default" {
  name    = var.name
  prefix  = var.name
  zone_id = alibabacloudstack_dns_private_domain.default.id
  ttl     = 300
}

resource "alibabacloudstack_dns_private_line" "default" {
  name         = var.name
  v4_addresses = ["192.168.0.1"]
  v6_addresses = ["2020:148:2:28::", "2020:148:3:28::"]
}

resource "alibabacloudstack_dns_gtm_addresspool" "default" {
  name         = "${var.name}-default"
  type         = "A"
  lba_strategy = "ALL_RR"
  addrs {
    value = "1.1.1.1"
    mode  = "SMART"
  }
}

resource "alibabacloudstack_dns_gtm_addresspool" "failover" {
  name         = "${var.name}-failover"
  type         = "A"
  lba_strategy = "ALL_RR"
  addrs {
    value = "2.2.2.2"
    mode  = "SMART"
  }
}

resource "alibabacloudstack_dns_gtm_access_strategy" "default" {
  name                            = var.name
  gtm_instance_id                 = alibabacloudstack_dns_gtm_instance.default.id
  switch_mode                     = "BY_PROBE_RESULT"
  default_min_available_addr_num  = 1
  failover_min_available_addr_num = 1
  line_ids                        = [alibabacloudstack_dns_private_line.default.id]
  default_gtm_address_pool_id     = alibabacloudstack_dns_gtm_addresspool.default.id
  default_gtm_address_pool_type   = "IPV4"
  failover_gtm_address_pool_id    = alibabacloudstack_dns_gtm_addresspool.failover.id
  failover_gtm_address_pool_type  = "IPV4"
}

data "alibabacloudstack_dns_gtm_access_strategies" "default" {
  gtm_instance_id = alibabacloudstack_dns_gtm_instance.default.id
  name_regex      = alibabacloudstack_dns_gtm_access_strategy.default.name
}
```

## Argument Reference

The following arguments are supported:

* `gtm_instance_id` (String, Required): The ID of the GTM instance.
* `name_regex` (String, Optional): A regular expression to filter access strategy names.
* `ids` (List, Optional): A list of access strategy IDs to filter.

## Attributes Reference

The following attributes are exported:

* `id` (String): The unique identifier of the access strategy, in the format "GtmInstanceId:AccessStrategyId".
* `default_available_addr_num` (Integer): The current number of available addresses in the default address pool.
* `default_gtm_address_pool_id` (String): The ID of the default address pool.
* `default_gtm_address_pool_name` (String): The name of the default address pool.
* `default_gtm_address_pool_type` (String): The type of the default address pool. Valid values: DOMAIN (domain), IPV6 (IPv6), IPV4 (IPv4).
* `default_min_available_addr_num` (Integer): The minimum number of available addresses in the default address pool.
* `failover_available_addr_num` (Integer): The current number of available addresses in the failover address pool.
* `failover_gtm_address_pool_id` (String): The ID of the failover address pool.
* `failover_gtm_address_pool_name` (String): The name of the failover address pool.
* `failover_gtm_address_pool_type` (String): The type of the failover address pool. Valid values: DOMAIN (domain), IPV6 (IPv6), IPV4 (IPv4).
* `failover_min_available_addr_num` (Integer): The minimum number of available addresses in the failover address pool.
* `gtm_instance_id` (String): The ID of the GTM instance.
* `in_use_gtm_address_pool_id` (String): The ID of the currently used address pool.
* `in_use_gtm_address_pool_name` (String): The name of the currently used address pool.
* `line_ids` (Set): A list of line IDs.
* `name` (String): The name of the access strategy.
* `specified_gtm_address_pool` (String): The manually specified address pool currently in use. Valid values: DEFAULT (default address pool), FAILOVER (failover address pool).
* `switch_mode` (String): The address pool switching strategy. Valid values: BY_HAND (manual), BY_PROBE_RESULT (automatic).