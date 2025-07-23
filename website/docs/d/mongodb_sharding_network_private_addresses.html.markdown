---
subcategory: "MongoDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_mongodb_shardingnetworkprivateaddresses"
sidebar_current: "docs-Alibabacloudstack-datasource-mongodb-shardingnetworkprivateaddresses"
description: |-
  Provides a list of mongodb shardingnetworkprivateaddresses owned by an alibabacloudstack account.
---

# alibabacloudstack\_mongodb\_shardingnetworkprivateaddresses

This data source provides a list of mongodb shardingnetworkprivateaddresses in an alibabacloudstack account according to the specified filters.

## Example Usage
```
variable "name" {
	default = "tf-testAlibabacloudstackMongodbShardingNetworkPrivateAddresses15244"
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
	db_instance_id="${alibabacloudstack_mongodb_sharding_instance.default.id}"
	account_name="terraform"
	account_password = "Y8Q_a@^MD$oM"
	node_id="${alibabacloudstack_mongodb_sharding_instance.default.mongo_list.0.node_id}"
}

data "alibabacloudstack_mongodb_sharding_network_private_addresses" "default" {
	ids = ["${alibabacloudstack_mongodb_sharding_network_private_address.default.id}"]
   db_instance_id = "${alibabacloudstack_mongodb_sharding_network_private_address.default.db_instance_id}"
}
```

## Argument Reference

The following arguments are supported:
  * `ids` - (Optional) - the IDs of network addresses
  * `node_id` - (Optional) - The ID of the Mongos node, Shard node, or ConfigServer node in the sharded cluster instance.> You can call the [DescribeDBInstanceAttribute](~~ 62010 ~~) interface to query the Mongos, Shard, and ConfigServer node ID.
  * `db_instance_id` - (Required) - the ID of the mongodb instance

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `network_addresses` - A list of connection addresses of an instance of the MongoDB protocol type.
