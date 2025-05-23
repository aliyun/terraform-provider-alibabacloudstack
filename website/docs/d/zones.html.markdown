---
subcategory: "Zone"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack"
sidebar_current: "docs-alibabacloudstack-datasource-zones"
description: |-
  Provides a list of availability zones that can be used by an Alibabacloudstack Cloud account.

---

# alibabacloudstack_zones

This data source provides availability zones that can be accessed by an Alibabacloudstack Cloud account within the region configured in the provider.


-> **NOTE:** If one zone is sold out, it will not be exported.

## Example Usage

```
# Declare the data source
data "alibabacloudstack_zones" "zones_ds" {
  available_instance_type = "ecs.n4.large"
  available_disk_category = "cloud_ssd"
}

output "zones" {
  value = data.alibabacloudstack_zones.zones_ds.zones.*
}
```

## Argument Reference

The following arguments are supported:

* `available_instance_type` - (Optional) Filter the results by a specific instance type.
* `available_resource_creation` - (Optional) Filter the results by a specific resource type.
Valid values: `Instance`, `Disk`, `VSwitch`, `Rds`, `KVStore`, `Slb`.
* `available_disk_category` - (Optional) Filter the results by a specific disk category. Can be either `cloud`, `cloud_efficiency`, `cloud_ssd`, `ephemeral_ssd`.
* `multi` - (Optional, type: bool) Indicate whether the zones can be used in a multi AZ configuration. Default to `false`. Multi AZ is usually used to launch RDS instances.
* `instance_charge_type` - (Optional) Filter the results by a specific ECS instance charge type. Valid values: `PrePaid` and `PostPaid`. Default to `PostPaid`.
* `network_type` - (Optional) Filter the results by a specific network type. Valid values: `Classic` and `Vpc`.
* `spot_strategy` - - (Optional) Filter the results by a specific ECS spot type. Valid values: `NoSpot`, `SpotWithPriceLimit` and `SpotAsPriceGo`. Default to `NoSpot`.
* `enable_details` - (Optional) Default to false and only output `id` in the `zones` block. Set it to true can output more details.
* `available_slb_address_type` - (Optional) Filter the results by a slb instance address type. Can be either `Vpc`, `classic_internet` or `classic_intranet`.
* `available_slb_address_ip_version` - (Optional) Filter the results by a slb instance address version. Can be either `ipv4`, or `ipv6`.

-> **NOTE:** The disk category `cloud` has been outdated and can only be used by non-I/O Optimized ECS instances. Many availability zones don't support it. It is recommended to use `cloud_efficiency` or `cloud_ssd`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of zone IDs.
* `zones` - A list of availability zones. Each element contains the following attributes:
  * `id` - ID of the zone.
  * `local_name` - Name of the zone in the local language.
  * `available_instance_types` - Allowed instance types.
  * `available_resource_creation` - Type of resources that can be created.
  * `available_disk_categories` - Set of supported disk categories.
  * `multi_zone_ids` - A list of zone ids in which the multi zone.
  * `slb_slave_zone_ids` - A list of slb slave zone ids in which the slb master zone.