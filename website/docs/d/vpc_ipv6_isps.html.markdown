---
subcategory: "Virtual Private Cloud (VPC)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpc_ipv6_isps"
description: |- 
  Provides a list of vpc ipv6 isps in an alibabacloudstack account.
---

The alibabacloudstack_vpc_ipv6_isps data source provides a list of vpc ipv6 isps in an alibabacloudstack account.

## Example Usage

```
data "alibabacloudstack_vpc_ipv6_isps" "example" {
  service_provider = "BGP"
  lock_status      = "unlocked"
}
```

## Argument Reference
The following arguments are supported:

* `ids` - (Optional) A list of specific IPv6 ISP pool IDs to filter results.
* `lock_status` - (Optional) Filter results by the lock status (locked or unlocked).
* `service_provider` - (Optional) Filter results by the service provider name, such as China Mobile, China Telecom, etc.


## Attributes Reference
The following attributes are exported:

* `ipv6_isps` - A list of IPv6 ISP entries. Each entry contains:
  * `id` - The ID of the IPv6 ISP pool (PoolId).
  * `service_provider` - The service provider name.
  * `zone_id` - The availability zone ID where the ISP is available.
  * `type` - The type of the IPv6 ISP.
  * `cidr_block` - The CIDR block assigned to this ISP.
  * `available_count` - Number of available IPv6 addresses.
  * `in_use_count` - In-use IPv6 addresses count.
  * `lock_status` - The current lock status of the pool.
  * `need_declare` - Indicates if declaration is required for this ISP.
  * `pool_id` - The unique identifier of the IPv6 address pool.
  * `ula` - Unique Local Address (ULA) associated with this ISP.
