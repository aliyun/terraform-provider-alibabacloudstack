---
subcategory: "MongoDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_mongodb_instances"
sidebar_current: "docs-Alibabacloudstack-datasource-mongodb-instances"
description: |- 
  查询mongodb instances
---

# alibabacloudstack_mongodb_instances

根据指定过滤条件列出当前凭证权限可以访问的mongodb instances列表。

## 示例用法

```hcl
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "MongoDB"
}

variable "name" {
  default = "tf-testAccMongoDBInstance_datasource_15722"
}

resource "alibabacloudstack_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = "${alibabacloudstack_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  name              = "${var.name}"
}

resource "alibabacloudstack_mongodb_instance" "default" {
  vswitch_id          = alibabacloudstack_vswitch.default.id
  engine_version      = "3.4"
  db_instance_class   = "dds.mongo.mid"
  db_instance_storage = "10"
  name                = "${var.name}"
  storage_engine      = "WiredTiger"
  instance_charge_type = "PostPaid"
  replication_factor = "3"
}

data "alibabacloudstack_mongodb_instances" "default" {
  name_regex        = "${alibabacloudstack_mongodb_instance.default.name}"
  instance_type     = "replicate"
  instance_class    = "dds.mongo.mid"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"

  tags = {
    Environment = "Test"
    CreatedBy   = "Terraform"
  }
}

output "mongodb_instance_ids" {
  value = data.alibabacloudstack_mongodb_instances.default.ids
}

output "mongodb_instance_names" {
  value = data.alibabacloudstack_mongodb_instances.default.names
}
```

## 参数参考

以下参数是支持的：

* `name_regex` - (选填) 应用于实例名称的正则表达式字符串。这允许基于实例名称使用正则表达式进行过滤。
* `ids` - (选填, 可用版本 v1.53.0+) MongoDB 实例 ID 列表。使用此参数通过特定实例 ID 过滤结果。
* `instance_type` - (选填) 要查询的实例类型。如果设置为 `sharding`，将列出分片集群实例；如果设置为 `replicate`，将列出副本集实例。默认值是 `replicate`。
* `instance_class` - (选填) 要查询的实例规格。这对应于 MongoDB 实例的性能级别。
* `availability_zone` - (选填) 实例可用区。使用此参数通过特定可用区过滤结果。
* `tags` - (选填, 可用版本 v1.66.0+) 分配给资源的标签映射。使用此参数通过标签过滤结果。

## 属性参考

除了上述参数外，还导出以下属性：

* `names` - 所有匹配的 MongoDB 实例的名称列表。
* `instances` - MongoDB 实例列表。每个元素包含以下属性：
  * `id` - MongoDB 实例的 ID。
  * `name` - MongoDB 实例的名称。
  * `charge_type` - 计费方式。选项包括 `PostPaid`(按量付费)和 `PrePaid`(包年包月订阅)。
  * `instance_type` - 实例类型。选项包括 `sharding`(分片集群)或 `replicate`(副本集)。
  * `region_id` - 实例所属的区域 ID。
  * `creation_time` - 实例创建时间，格式为 RFC3339。
  * `expiration_time` - 实例到期时间，格式为 RFC3339。按量付费实例不会过期。
  * `status` - 实例状态。
  * `replication` - 副本因子，对应节点数量。选项包括 `1`(单节点)和 `3`(三节点副本集)。
  * `engine` - 数据库引擎类型。支持的选项是 `MongoDB`。
  * `engine_version` - 数据库引擎版本。
  * `network_type` - 网络类型。选项包括经典网络或 VPC。
  * `lock_mode` - 实例锁定状态。可能的值包括：`Unlock`(正常)、`ManualLock`(手动触发锁定)、`LockByExpiration`(实例过期自动锁定)、`LockByRestoration`(实例回滚前自动锁定)、`LockByDiskQuota`(实例空间满自动锁定)、`Released`(实例已释放)。
  * `instance_class` - MongoDB 实例规格。
  * `storage` - 存储大小，单位为 GB。
  * `mongos` - Mongos 节点组成的数组。每个元素包含：
    * `node_id` - Mongos 实例 ID。
    * `description` - Mongos 实例描述。
    * `class` - Mongos 实例规格。
  * `shards` - 分片组成的数组。每个元素包含：
    * `node_id` - 分片实例 ID。
    * `description` - 分片实例描述。
    * `class` - 分片实例规格。
    * `storage` - 分片磁盘大小，单位为 GB。
  * `availability_zone` - 实例可用区。
  * `tags` - 分配给资源的标签映射。