---
subcategory: "DNS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dns_gtm_access_strategy"
sidebar_current: "docs-Alibabacloudstack-dns-dns_gtm_access_strategy"
description: |-
  Access strategy for Global Traffic Manager (GTM) instance in Alibaba Cloud DNS
---

# alibabacloudstack_dns_gtm_access_strategy

Creates an access strategy for Global Traffic Manager (GTM) instance in Alibaba Cloud DNS using the credentials configured in the provider.

## Example Usage

### Basic Usage

```hcl
variable "name" {
  default = "tfacc35194"
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

resource "alibabacloudstack_dns_gtm_addresspool" "update" {
  name         = "${var.name}-update"
  type         = "A"
  lba_strategy = "ALL_RR"
  addrs {
    value = "3.3.3.3"
    mode  = "SMART"
  }
}


resource "alibabacloudstack_dns_gtm_access_strategy" "default" {
  default_gtm_address_pool_id    = alibabacloudstack_dns_gtm_addresspool.default.id
  failover_gtm_address_pool_id   = alibabacloudstack_dns_gtm_addresspool.failover.id
  failover_gtm_address_pool_type = "IPV4"
  default_min_available_addr_num = "1"
  line_ids = [
    "${alibabacloudstack_dns_private_line.default.id}"
  ]
  gtm_instance_id                 = alibabacloudstack_dns_gtm_instance.default.id
  default_gtm_address_pool_type   = "IPV4"
  failover_min_available_addr_num = "1"
  name                            = var.name
  switch_mode                     = "BY_PROBE_RESULT"
}
```

## Argument Reference

The following arguments are supported:

* `default_gtm_address_pool_id` - (Required) The ID of the primary address pool.
* `default_gtm_address_pool_type` - (Required) The type of the primary address pool. Valid values: DOMAIN (domain name), IPV6 (IPv6), IPV4 (IPv4).
* `default_min_available_addr_num` - (Required) The minimum number of available addresses in the primary address pool.
* `line_ids` - (Required) A list of line IDs to specify the lines to which the access strategy applies.
* `name` - (Required) The name of the access strategy.
* `switch_mode` - (Required) The address pool switching strategy. Valid values: BY_HAND (manual), BY_PROBE_RESULT (automatic).
* `gtm_instance_id` - (Required, ForceNew) The ID of the GTM instance. Changing this value will recreate the resource.
* `failover_gtm_address_pool_id` - (Optional) The ID of the failover address pool.
* `failover_gtm_address_pool_type` - (Optional) The type of the failover address pool. Valid values: DOMAIN (domain name), IPV6 (IPv6), IPV4 (IPv4).
* `failover_min_available_addr_num` - (Optional) The minimum number of available addresses in the failover address pool.
* `specified_gtm_address_pool` - (Optional) Manually specify the address pool currently in use. Valid values: DEFAULT (primary address pool), FAILOVER (failover address pool).

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the access strategy.
* `default_available_addr_num` - The current number of available addresses in the primary address pool.
* `default_gtm_address_pool_name` - The name of the primary address pool.
* `failover_available_addr_num` - The current number of available addresses in the failover address pool.
* `failover_gtm_address_pool_name` - The name of the failover address pool.
* `in_use_gtm_address_pool_id` - The ID of the address pool currently in use.
* `in_use_gtm_address_pool_name` - The name of the address pool currently in use.