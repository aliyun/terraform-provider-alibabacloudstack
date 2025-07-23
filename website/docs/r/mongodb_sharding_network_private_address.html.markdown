---
subcategory: "MongoDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_mongodb_shardingnetworkprivateaddress"
sidebar_current: "docs-Alibabacloudstack-mongodb-shardingnetworkprivateaddress"
description: |-
  Provides a mongodb Shardingnetworkprivateaddress resource.
---

# alibabacloudstack\_mongodb\_shardingnetworkprivateaddress

Provides a mongodb Shardingnetworkprivateaddress resource.

## Example Usage
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
  account_password = "123@123qwe"
  db_instance_id = "${alibabacloudstack_mongodb_sharding_instance.default.id}"
  node_id = "${alibabacloudstack_mongodb_sharding_instance.default.shard_list.0.node_id}"
}
```

## Argument Reference

The following arguments are supported:
  * `account_name` - (Optional) - Account name.>-starts with a lowercase letter, has 4 to 16 digits in length, and consists of lowercase letters, numbers, or underscores.-Only when you apply for the Shard/ConfigServer address for the first time, you need to set the account name and password. That is, all Shard nodes and ConfigServer nodes will use the account and password set when applying for the address for the first time.-The permissions of this account are fixed to read-only.
  * `account_password` - (Optional) - Account password.-The password consists of at least three of uppercase letters, lowercase letters, numbers, and special characters. The special character is '! #$%^& *()_+-='-The password length is 8-32 bits.
  * `network_address` - (Optional) - the network address of the instance.
  * `port` - (Optional) - the port of the instance.
  * `db_instance_id` - (Required) - the ID of the instance.
  * `node_id` - (Required) - The ID of the Mongos node, Shard node, or ConfigServer node in the sharded cluster instance.> You can call the [DescribeDBInstanceAttribute](~~ 62010 ~~) interface to query the Mongos, Shard, and ConfigServer node ID.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `network_addresses` - A list of connection addresses of an instance of the MongoDB protocol type.
    * `ip_address` - the IP address of the connection address.
    * `network_address` - Connection address (string).
    * `network_type` - Network type.-**VPC**: VPC.-**Classic**: Classic network.-**Public**: Public network.
    * `node_id` - The ID of the Mongos node.
    * `node_type` - Node type, return value is-**mongos**:mongos node.-**shard**:shard node.-**configserver**:configserver node.
    * `port` - Connection port.
    * `role` - Node role, return value:-Primary: Primary node.-Secondary: Slave node.
