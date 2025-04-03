---
subcategory: "Redis And Memcache (KVStore)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_kvstore_zones"
sidebar_current: "docs-alibabacloudstack-datasource-kvstore-zones"
description: |-
    Provides a list of availability zones for KVStore that can be used by an Apsara Stack Cloud account.
---

# alibabacloudstack_kvstore_zones

This data source provides availability zones for KVStore that can be accessed by an Apsara Stack Cloud account within the region configured in the provider.


## Example Usage

```
# Declare the data source
data "alibabacloudstack_kvstore_zones" "zones_ids" {}

output "kvstore_zones" {
  value = "${data.alibabacloudstack_kvstore_zones.zones_ids.zones}"
}
```

## Argument Reference

The following arguments are supported:

* `multi` - (Optional) Indicate whether the zones can be used in a multi AZ configuration. Default to `false`. Multi AZ is usually used to launch KVStore instances.
* `instance_charge_type` - (Optional) Filter the results by a specific instance charge type. Valid values: `PrePaid` and `PostPaid`. Default to `PostPaid`.
* `ids` - (Optional) A list of zone IDs.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of zone IDs.
* `zones` - A list of availability zones. Each element contains the following attributes:
  * `id` - ID of the zone.
  * `multi_zone_ids` - A list of zone ids in which the multi zone.
  * `zones` - A list of availability zones.