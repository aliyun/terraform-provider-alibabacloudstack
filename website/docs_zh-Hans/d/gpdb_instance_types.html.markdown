---
subcategory: "Greenplum Database"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_gpdb_instance_types"
sidebar_current: "docs-Alibabacloudstack-datasource-gpdb-instance-types"
description: |-
  提供GPDB实例类型列表。
---

# alibabacloudstack_gpdb_instance_types

该数据源用于获取专有云中可用的GPDB（Greenplum Database）实例类型列表。

-> **注意:** 适用于专有云环境。

## 示例

```hcl
data "alibabacloudstack_gpdb_instance_types" "example" {
  cpu    = 4
  memory = 16
}

output "instance_types" {
  value = data.alibabacloudstack_gpdb_instance_types.example.instance_types
}
```

## 参数说明

以下参数支持配置：

* `ids` - (可选, ForceNew) 用于过滤结果的实例类型ID列表。
* `engine_version` - (可选) 用于过滤实例类型的引擎版本。
* `cpu` - (可选) 用于过滤实例类型的CPU核心数。
* `memory` - (可选) 用于过滤实例类型的内存大小（GB）。
* `status` - (可选) 用于过滤实例类型的状态。
* `sorted_by` - (可选, ForceNew) 排序字段。有效值：`CPU`、`Memory`。

## 属性参考

以下属性会被导出：

* `ids` - 实例类型ID列表。
* `instance_types` - GPDB实例类型列表。每个元素包含以下属性：
  * `id` - 实例类型的ID。
  * `cpu` - CPU核心数。
  * `memory` - 内存大小（GB）。
  * `engine` - 数据库引擎。
  * `engine_version` - 引擎版本。
  * `connections` - 最大连接数。
  * `storage_type` - 存储类型。
  * `storage_min` - 最小存储空间（GB）。
  * `storage_max` - 最大存储空间（GB）。
  * `specification` - 规格代码。
  * `specification_label` - 规格标签。
  * `db_instance_mode` - 数据库实例模式。
  * `db_instance_mode_label` - 数据库实例模式标签。
  * `node` - 节点配置。
  * `region_id` - 区域ID。
  * `status` - 实例类型的状态。
  * `product` - 产品名称。
  * `storage` - 存储配置。
  * `cpu_label` - CPU标签。
  * `memory_label` - 内存标签。
  * `storage_label` - 存储标签。
  * `engine_version_label` - 引擎版本标签。
  * `gmt_create` - 创建时间戳。
  * `gmt_modify` - 修改时间戳。
  * `spec_from` - 规格来源。
