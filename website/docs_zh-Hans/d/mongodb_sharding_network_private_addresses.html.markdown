---
subcategory: "MongoDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_mongodb_shardingnetworkprivateaddress"
sidebar_current: "docs-Alibabacloudstack-mongodb-shardingnetworkprivateaddress"
description: |-
提供一个 MongoDB 分片网络私有地址资源。
---

# alibabacloudstack\_mongodb\_shardingnetworkprivateaddresses
提供一个 MongoDB 分片网络私有地址资源。

## 示例用法
```
variable "name" {
	default = "tfmongodb_sharding_network_privateaddrdss32446"
}

data "alibabacloudstack_zones" "default" {
	available_resource_creation = "MongoDB"
  }

resource "alibabacloudstack_mongodb_sharding_instance" "default" {
	zone_id        = "${data.alibabacloudstack_zones.default.zones.0.id}"
	engine_version = "3.4"
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
	  node_class = "dds.mongos.mid"
	}
}



resource "alibabacloudstack_mongodb_sharding_network_private_address" "default" {
  account_name = "terraform"
  account_password = "xxxx"
  db_instance_id = "${alibabacloudstack_mongodb_sharding_instance.default.id}"
  node_id = "${alibabacloudstack_mongodb_sharding_instance.default.shard_list.0.node_id}"
}
```
## 参数说明
支持以下参数：

* `account_name` - (可选) - 账户名称。>- 以小写字母开头，长度为4到16位，由小写字母、数字或下划线组成。- 只有在第一次申请分片/配置服务器地址时，需要设置账户名和密码。即所有分片节点和配置服务器节点将使用第一次申请地址时设置的账户和密码。- 该账户的权限固定为只读。
* `account_password` - (可选) - 账户密码。- 密码需包含大小写字母、数字、特殊字符中的至少三项。特殊字符包括：'! #$%^& *()_+-='。- 密码长度为8到32位。
* `network_address` - (可选) - 实例的网络地址。
* `port` - (可选) - 实例的端口。
* `db_instance_id` - (必填) - 实例的 ID。
* `node_id` - (必填) - 分片集群实例中的 Mongos 节点、Shard 节点或 ConfigServer 节点的 ID。> 可以调用 [DescribeDBInstanceAttribute](~~ 62010 ~~) 接口查询 Mongos、Shard 和 ConfigServer 节点的 ID。
属性输出
除了上述参数外，还导出以下属性：

* `network_addresses` - MongoDB 协议类型实例的连接地址列表。
* `ip_address` - 连接地址的 IP。
* `network_address` - 连接地址（字符串）。
* `network_type` - 网络类型。- VPC：VPC 网络。- Classic：经典网络。- Public：公网。
* `node_id` - Mongos 节点的 ID。
* `node_type` - 节点类型，返回值为：- mongos：mongos 节点。- shard：shard 节点。- configserver：* configserver 节点。
* `port` - 连接端口。
* `role` - 节点角色，返回值为：- Primary：主节点。- Secondary：从节点。