---
subcategory: "Alibaba Parallel File System (APFS)"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_apfs_zones"
sidebar_current: "docs-alibabacloudstack-datasource-apfs-zones"
description: |-
  Provides a list of Apfs Zones to the user.
---

# alibabacloudstack\_apfs\_zones

该数据源提供当前阿里云用户可用的APFS可用区列表。

## Example Usage

### 基础用法

```terraform
data "alibabacloudstack_apfs_zones" "example" {
}

output "first_apfs_zone_id" {
  value = data.alibabacloudstack_apfs_zones.example.zones.0.id
}
```

### x按可用区ID过滤

```terraform
data "alibabacloudstack_apfs_zones" "example" {
  zone_id = "cn-wulan-env119-amtest38002-b"
}

output "filtered_zone_id" {
  value = data.alibabacloudstack_apfs_zones.example.zones.0.id
}
```

## Argument Reference

以下参数可用于过滤数据：

* `zone_id` - (可选) 按可用区ID过滤。
* `cluster_id` - (可选) 按集群ID过滤。

## Attributes Reference

以下属性会被导出：

* `zones` - APFS可用区列表。每个元素包含以下属性：
  * `id` - 可用区ID。
  * `zone_id` - 可用区ID。
  * `clusters` - 可用区内的集群列表。每个元素包含：
    * `cluster_id` - 集群ID。
    * `available_capacity` - 集群可用容量。
    * `storage_type` - 实例的存储类型。