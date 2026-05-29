---
subcategory: "Universal DNS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dns_gtm_addresspool"
sidebar_current: "docs-alibabacloudstack-dns-dns_gtm_addresspool"
description: |-
  Manages a Cloud DNS Global Traffic Management Address Pool resource in Alibaba Cloud.
---

# alibabacloudstack_dns_gtm_addresspool

Creates a Cloud DNS Global Traffic Management Address Pool in the specified resource set using the credentials configured in the provider.

## Example Usage

### Basic Usage

```hcl
variable "name" {
  default = "tfacc16606"
}

resource "alibabacloudstack_dns_gtm_addresspool" "default" {
  name         = var.name
  type         = "A"
  lba_strategy = "RATIO"
  addrs {
    value      = "192.168.1.1"
    mode       = "SMART"
    lba_weight = 20
  }
  addrs {
    value      = "127.0.0.1"
    mode       = "SMART"
    lba_weight = 80
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the address pool.
* `type` - (Required) The type of the address pool. Valid values: `A` (IPv4 address), `AAAA` (IPv6 address), `CNAME` (domain name).
* `lba_strategy` - (Required) The load balancing strategy. Valid values: `ALL_RR` (return all addresses without weight), `RATIO` (return addresses by weight).
* `addrs` - (Required) The list of addresses. Each address block contains the following attributes:
  * `value` - (Required) The address value.
  * `mode` - (Required) The mode. Valid values: `SMART` (smart return), `ONLINE` (always online), `OFFLINE` (always offline).
  * `lba_weight` - (Optional) The weight. Valid values: 0 to 100. This parameter is valid only when `lba_strategy` is set to `RATIO`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the address pool.

## Import

DNS GTM Address Pool can be imported using the address pool ID, e.g.

```
$ terraform import alibabacloudstack_dns_gtm_addresspool.example <address_pool_id>
```