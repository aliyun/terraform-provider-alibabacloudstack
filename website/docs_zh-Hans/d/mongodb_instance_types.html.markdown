---
subcategory: "云数据库 MongoDB 版"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_mongodb_instance_types"
sidebar_current: "docs-Alibabacloudstack-datasource-mongodb-instance-types"
description: |-
  查询Mongodb实例规格
---

# alibabacloudstack_mongodb_instance_types

根据指定过滤条件列出当前凭证权限可以访问的Mongodb实例规格列表。

## 示例用法

```
data "alibabacloudstack_mongodb_instance_types" "default" {
  engine_version = "4.0"
  db_instnace_type = "replicate"
  sorted_by = "CPU"
}

data "alibabacloudstack_mongodb_instance_types" "shard" {
  engine_version = "4.0"
  db_instnace_type = "sharding"
  node_type = "shard"
  sorted_by = "CPU"
}

```

## 参数说明

支持以下参数：

* `ids` - (可选，变更时重建) 指定实例规格ID范围。如果未指定，则返回所有可用区的实例规格。
* `db_instnace_type` - (必填, 变更时重建) Mongodb实例类型，有效值：`sharding`, `replicate`。
* `node_type` - (可选, 变更时重建) 分片存储实例的节点类型。有效值： `configserver`, `shard`, `mongos`。
* `engine_version` - (可选，变更时重建) 筛选特定例引擎版本的结果。
* `cpu` - (可选) 筛选特定数量CPU核心的结果。
* `cpu_type` - (可选) 筛选特定数量CPU类型的结果。
* `memory` - (可选) 筛选特定内存大小(GB)的结果。
* `sorted_by` - (可选，强制更新) 排序模式，有效值：`CPU`, `Memory`。

## 属性说明

除了上述列出的参数外，还导出以下属性：

* `ids` - 包含所有匹配条件的实例规格ID的列表。
* `names` - 包含所有匹配条件的实例规格名称的列表。
* `instance_types` - 实例规格的详细信息列表。每个元素包含以下属性：
  * `id` - 实例规格的唯一标识符。
  * `name` - 实例规格的名称。
  * `cpu` - CPU核心数量。
  * `memory` - 内存大小，以GB为单位。
  * `engine_version` - 引擎版本。
  * `cpu_type` - CPU类型。
  * `series` - 隶属的实例规格系列的ID。
  * `connections` - 最大链接数。
  * `storage_min` - 最小存储规格。
  * `storage_max` - 最大存储规格。
