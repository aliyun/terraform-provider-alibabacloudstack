---
subcategory: "MongoDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_mongodb_shardinginstance"
sidebar_current: "docs-Alibabacloudstack-mongodb-shardinginstance"
description: |- 
  编排mongodb共享实例
---

# alibabacloudstack_mongodb_shardinginstance
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_mongodb_sharding_instance`

使用Provider配置的凭证在指定的资源集编排mongodb共享实例。

## 示例用法

### 创建具有 VPC 配置的 MongoDB 分片实例

```hcl
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "MongoDB"
}

variable "password" {}

resource "alibabacloudstack_vpc" "example" {
  name       = "tf-example-vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vswitch" "example" {
  vpc_id     = alibabacloudstack_vpc.example.id
  cidr_block = "172.16.0.0/24"
  zone_id    = data.alibabacloudstack_zones.default.zones[0].id
  name       = "tf-example-vswitch"
}

resource "alibabacloudstack_mongodb_sharding_instance" "default" {
  zone_id        = data.alibabacloudstack_zones.default.zones[0].id
  vswitch_id     = alibabacloudstack_vswitch.example.id
  engine_version = "3.4"
  storage_engine = "WiredTiger"
  name           = "tf-example-instance"

  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
  }

  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }

  mongo_list {
    node_class = "dds.mongos.mid"
  }

  mongo_list {
    node_class = "dds.mongos.large"
  }

  account_password = var.password
}
```

## 参数参考

支持以下参数：

* `engine_version` - (必填，变更时重建) 数据库版本。值选项可以参考最新文档 [CreateDBInstance](https://www.alibabacloud.com/help/zh/doc-detail/61884.htm) 的 `EngineVersion`。
* `storage_engine` - (选填，变更时重建) 实例的存储引擎类型。有效值：`WiredTiger`、`RocksDB`。默认值：`WiredTiger`。
* `instance_charge_type` - (选填，变更时重建) 有效的值为 `PrePaid` 和 `PostPaid`。系统默认值为 `PostPaid`。**注意**：从 v1.141.0 版本开始，可以从 `PostPaid` 修改为 `PrePaid`。
* `period` - (选填)购买 DB 实例的时长(以月为单位)。当 `instance_charge_type` 为 `PrePaid` 时有效。有效值：[1~9]、12、24、36。系统默认值为 1。
* `zone_id` - (选填，变更时重建) 启动 DB 实例的可用区。MongoDB 分片实例不支持多可用区。如果它是多可用区并且指定了 `vswitch_id`，交换机必须在其中一个可用区内。
* `vswitch_id` - (选填，变更时重建) 用于在 VPC 中启动 DB 实例的虚拟交换机 ID。
* `name` - (选填)DB 实例的名称。它是一个长度为 2 到 256 个字符的字符串。
* `db_instance_description` - (选填)DB 实例的描述。它是一个长度为 2 到 256 个字符的字符串。
* `security_group_id` - (选填)ECS 的安全组 ID。
* `account_password` - (选填，敏感)root 账户的密码。它是一个长度为 6 到 32 个字符的字符串，由字母、数字和下划线组成。
* `kms_encrypted_password` - (选填)用于创建实例的 KMS 加密密码。如果填写了 `account_password`，此字段将被忽略。
* `kms_encryption_context` - (选填)用于在使用 `kms_encrypted_password` 创建或更新实例之前解密 `kms_encrypted_password` 的 KMS 加密上下文。参见 [加密上下文](https://www.alibabacloud.com/help/doc-detail/42975.htm)。当设置了 `kms_encrypted_password` 时有效。
* `tde_status` - (选填，变更时重建) 透明数据加密 (TDE) 状态。有效值：`Enabled`、`Disabled`。
* `backup_time` - (选填)MongoDB 实例备份时间。格式为 HH:mmZ- HH:mmZ。时间设置间隔为一小时。如果不设置，系统将返回默认值，例如 "23:00Z-24:00Z"。
* `preferred_backup_time` - (选填)备份时间，格式为 HH:mmZ-HH:mmZ(UTC 时间)。
* `shard_list` - (必填) 分片节点列表。每个分片节点具有以下属性：
  * `node_class` - (必填) 节点规格。参见 [实例规格](https://www.alibabacloud.com/help/doc-detail/57141.htm)。
  * `node_storage` - (必填) 自定义存储空间；范围：[10, 1,000]，以 10 GB 为增量。单位：GB。
  * `readonly_replicas` - (选填)分片节点中的只读节点数量。有效值：0 到 5。默认值：0。
* `mongo_list` - (必填) Mongo 节点列表。每个 Mongo 节点具有以下属性：
  * `node_class` - (必填) 节点规格。参见 [实例规格](https://www.alibabacloud.com/help/doc-detail/57141.htm)。

## 属性参考

除了上述所有参数外，还导出了以下属性：

* `id` - MongoDB 实例的 ID。
* `mongo_list` - Mongo 节点列表。每个 Mongo 节点包含以下属性：
  * `node_id` - Mongo 节点的 ID。
  * `connect_string` - Mongo 节点连接字符串。
  * `port` - Mongo 节点端口。
* `shard_list` - 分片节点列表。每个分片节点包含以下属性：
  * `node_id` - 分片节点的 ID。
* `retention_period` - 实例日志备份保留天数。
* `config_server_list` - Config Server 节点信息列表。每个 Config Server 节点包含以下属性：
  * `max_iops` - Config Server 节点的最大 IOPS。
  * `connect_string` - Config Server 节点的连接地址。
  * `node_class` - Config Server 节点的节点类。
  * `max_connections` - Config Server 节点的最大连接数。
  * `port` - Config Server 节点的连接端口。
  * `node_description` - Config Server 节点的描述。
  * `node_id` - Config Server 节点的 ID。
  * `node_storage` - Config Server 节点的存储容量。