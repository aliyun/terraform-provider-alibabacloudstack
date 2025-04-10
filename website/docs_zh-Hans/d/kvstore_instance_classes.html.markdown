---
subcategory: "Redis 和 Memcache (KVStore)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_kvstore_instance_classes"
sidebar_current: "docs-alibabacloudstack-datasource-kvstore-instance-classes"
description: |-
   查询KVStore实例类信息
---

# alibabacloudstack_kvstore_instance_classes

根据指定过滤条件列出当前凭证权限可以访问的KVStore实例类列表。

## 示例用法

```tf
data "alibabacloudstack_zones" "resources" {
  available_resource_creation = "KVStore"
}

data "alibabacloudstack_kvstore_instance_classes" "resources" {
  zone_id              = "${data.alibabacloudstack_zones.resources.zones.0.id}"
  instance_charge_type = "PrePaid"
  engine               = "Redis"
  engine_version       = "5.0"
  output_file          = "./classes.txt"
}

output "first_kvstore_instance_class" {
  value = "${data.alibabacloudstack_kvstore_instance_classes.resources.instance_classes}"
}
```

## 参数说明

以下是支持的参数：

* `zone_id` - (必填) 启动 KVStore 实例所在的可用区。
* `instance_charge_type` - (可选) 通过计费类型筛选结果。有效值：`PrePaid`（预付费）和 `PostPaid`（后付费）。默认为 `PrePaid`。
* `sorted_by` - (可选，强制更新) 排序模式，有效值：`Price`（按价格排序）。
* `engine` - (可选) 数据库类型。选项为 `Redis` 和 `Memcache`。默认为 `Redis`。
* `engine_version` - (可选) 用户所需的数据库版本。Redis 的可选项可以参考最新文档 [详细信息](https://www.alibabacloud.com/help/doc-detail/60873.htm) `EngineVersion`。对于 Memcache，该值应为空。
* `architecture` - (可选) 用户所需的 KVStore 实例系统架构。有效值：`standard`（标准版）、`cluster`（集群版）和 `rwsplit`（读写分离版）。
* `node_type` - (可选) 用户所需的 KVStore 实例节点类型。有效值：`double`（双副本）、`single`（单副本）、`readone`（一个只读副本）、`readthree`（三个只读副本）和 `readfive`（五个只读副本）。
* `edition_type` - (可选) 用户所需的 KVStore 实例版本类型。有效值：`Community`（社区版）和 `Enterprise`（企业版）。
* `series_type` - (可选) 用户所需的 KVStore 实例系列类型。有效值：`enhanced_performance_type`（性能增强型）和 `hybrid_storage`（混合存储型）。
* `shard_number` - (可选) 分片数量。有效值：`1`、`2`、`4`、`8`、`16`、`32`、`64`、`128` 和 `256`。

## 属性说明

除了上述列出的参数外，还导出以下属性：

* `instance_classes` - KVStore 可用实例类的列表。
* `classes` - 当 `sorted_by` 为 "Price" 时，KVStore 可用实例类的列表。包括：
  * `price` - 实例类型的单价。
  * `instance_class` - KVStore 可用实例类。