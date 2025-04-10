---
subcategory: "Zone"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack"
sidebar_current: "docs-alibabacloudstack-datasource-zones"
description: |-
  查询可用区
---

# alibabacloudstack_zones

根据指定过滤条件列出当前凭证权限可以访问的可用区


-> **注意:** 如果某个可用区已售罄，则不会导出该可用区。

## 示例用法

```
# 声明数据源
data "alibabacloudstack_zones" "zones_ds" {
  available_instance_type = "ecs.n4.large"
  available_disk_category = "cloud_ssd"
}

output "zones" {
  value = data.alibabacloudstack_zones.zones_ds.zones.*
}
```

## 参数说明

支持以下参数：

* `available_instance_type` - (可选) 通过特定实例类型过滤结果。
* `available_resource_creation` - (可选) 通过特定资源类型过滤结果。有效值：`Instance`, `Disk`, `VSwitch`, `Rds`, `KVStore`, `Slb`。
* `available_disk_category` - (可选) 通过特定磁盘类别过滤结果。可以是 `cloud`, `cloud_efficiency`, `cloud_ssd`, `ephemeral_ssd`。
* `multi` - (可选，类型：bool) 指示这些可用区是否可以在多AZ配置中使用。默认为 `false`。多AZ通常用于启动RDS实例。
* `instance_charge_type` - (可选) 通过特定ECS实例计费类型过滤结果。有效值：`PrePaid` 和 `PostPaid`。默认为 `PostPaid`。
* `network_type` - (可选) 通过特定网络类型过滤结果。有效值：`Classic` 和 `Vpc`。
* `spot_strategy` - (可选) 通过特定ECS竞价实例类型过滤结果。有效值：`NoSpot`, `SpotWithPriceLimit` 和 `SpotAsPriceGo`。默认为 `NoSpot`。
* `enable_details` - (可选) 默认为 false，仅在 `zones` 块中输出 `id`。将其设置为 true 可以输出更多详细信息。
* `available_slb_address_type` - (可选) 通过负载均衡实例地址类型过滤结果。可以是 `Vpc`, `classic_internet` 或 `classic_intranet`。
* `available_slb_address_ip_version` - (可选) 通过负载均衡实例地址版本过滤结果。可以是 `ipv4` 或 `ipv6`。

-> **注意:** 磁盘类别 `cloud` 已过时，只能由非I/O优化型ECS实例使用。许多可用区不再支持它。建议使用 `cloud_efficiency` 或 `cloud_ssd`。

## 属性说明

除了上述参数外，还导出以下属性：

* `ids` - 可用区ID列表。
* `zones` - 可用区列表。每个元素包含以下属性：
  * `id` - 可用区的ID。
  * `local_name` - 本地语言中的可用区名称。
  * `available_instance_types` - 允许的实例类型集合。
  * `available_resource_creation` - 可以创建的资源类型集合。可能的值包括：`Instance`, `Disk`, `VSwitch`, `Rds`, `KVStore`, `Slb`。
  * `available_disk_categories` - 支持的磁盘类别集合。可能的值包括：`cloud`, `cloud_efficiency`, `cloud_ssd`, `ephemeral_ssd`。
  * `multi_zone_ids` - 多可用区配置中使用的可用区ID列表。
  * `slb_slave_zone_ids` - 负载均衡主可用区中对应的从可用区ID列表。