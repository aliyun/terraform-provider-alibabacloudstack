---
subcategory: "Alibaba Parallel File System (APFS)"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_apfs_zones"
sidebar_current: "docs-alibabacloudstack-datasource-apfs-zones"
description: |-
  Provides a list of Apfs Zones to the user.
---

# alibabacloudstack\_apfs\_zones

This data source provides the APFS Zones of the current Alibaba Cloud user.

## Example Usage

Basic Usage

```terraform
data "alibabacloudstack_apfs_zones" "example" {
}

output "first_apfs_zone_id" {
  value = data.alibabacloudstack_apfs_zones.example.zones.0.id
}
```

Filter by zone ID

```terraform
data "alibabacloudstack_apfs_zones" "example" {
  zone_id = "cn-wulan-env119-amtest38002-b"
}

output "filtered_zone_id" {
  value = data.alibabacloudstack_apfs_zones.example.zones.0.id
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Optional) Filter by zone ID.
* `cluster_id` - (Optional) Filter by cluster ID.

## Attributes Reference

The following attributes are exported:

* `zones` - A list of APFS Zones. Each element contains the following attributes:
  * `id` - Zone ID.
  * `zone_id` - Zone ID.
  * `clusters` - List of clusters in the zone. Each element contains:
    * `cluster_id` - Cluster ID.
    * `available_capacity` - Available capacity of the cluster.
    * `storage_type` - Storage type of the instance.
