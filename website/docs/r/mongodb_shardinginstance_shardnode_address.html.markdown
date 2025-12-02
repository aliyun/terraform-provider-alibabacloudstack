---
subcategory: "MongoDB" 
layout: "alibabacloudstack" 
page_title: "Alibabacloudstack: alibabacloudstack_mongodb_shardinginstance_shardnode_address" 
sidebar_current: "docs-Alibabacloudstack-shardinginstance-shardnode-address" description: |- 
Manage Mongodb Sharding Instance Shard Node Address.
---

# alibabacloudstack_mongodb_shardinginstance_shardnode_address
Manage Mongodb Sharding Instance Shard Node Address

## Example Usage
```hcl
variable "name" {
  default = "tfmongodb_shardNode31768"
}

variable "existed_db_instance_id" {
  default = "dds-rw340c26dfb226b4"
}

data "alibabacloudstack_mongodb_instance_types" "mongos" {
  db_instnace_type = "sharding"
  node_type        = "mongos"
  sorted_by        = "CPU"
  engine_version   = "4.0"
}

data "alibabacloudstack_mongodb_instance_types" "configserver" {
  db_instnace_type = "sharding"
  node_type        = "configserver"
  sorted_by        = "CPU"
  engine_version   = "4.0"
}

data "alibabacloudstack_mongodb_instance_types" "shard" {
  db_instnace_type = "sharding"
  node_type        = "shard"
  sorted_by        = "CPU"
  engine_version   = "4.0"
}

resource "random_password" "password" {
  count            = 2
  length           = 12
  special          = true
  override_special = "!@#$^&*()_"
  min_lower        = 1
  min_upper        = 1
  min_numeric      = 1
}


data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details              = true
}

data "alibabacloudstack_mongodb_instances" "default" {
  ids           = var.existed_db_instance_id == "" ? [] : ["${var.existed_db_instance_id}", ]
  instance_type = "sharding"
}

resource "alibabacloudstack_mongodb_sharding_instance" "default" {
  count               = length(data.alibabacloudstack_mongodb_instances.default.instances) == 0 ? 1 : 0
  db_account_password = random_password.password.0.result
  zone_id             = data.alibabacloudstack_zones.default.zones.0.id
  engine_version      = "3.4"
  shard_list {
    node_storage = data.alibabacloudstack_mongodb_instance_types.shard.instance_types.0.storage_min
    description  = "shard1"
    node_class   = data.alibabacloudstack_mongodb_instance_types.shard.instance_types.0.id
  }
  shard_list {
    description  = "shard2"
    node_class   = data.alibabacloudstack_mongodb_instance_types.shard.instance_types.0.id
    node_storage = data.alibabacloudstack_mongodb_instance_types.shard.instance_types.0.storage_min
  }

  mongo_list {
    node_class  = data.alibabacloudstack_mongodb_instance_types.mongos.instance_types.0.id
    description = "mongo1"
  }
  mongo_list {
    description = "mongo2"
    node_class  = data.alibabacloudstack_mongodb_instance_types.mongos.instance_types.0.id
  }

  configserver_list {
    node_storage = data.alibabacloudstack_mongodb_instance_types.configserver.instance_types.0.storage_min
    description  = "cs1"
    node_class   = data.alibabacloudstack_mongodb_instance_types.configserver.instance_types.0.id
  }

  db_account_name = "tf_testacc"

}

locals {
  shard_instance = length(data.alibabacloudstack_mongodb_instances.default.instances) == 0 ? alibabacloudstack_mongodb_sharding_instance.default.0 : data.alibabacloudstack_mongodb_instances.default.instances.0
}



resource "alibabacloudstack_mongodb_shardinginstance_shardnode_address" "default" {
  enable_public_connection  = true
  enable_private_connection = true
  db_instance_id            = local.shard_instance.id
  node_id                   = [for shard in local.shard_instance.shard_list : shard if shard.description == "shard1"][0].node_id
}
```

## Argument Reference

The following arguments are supported:

* `db_instance_id` - (Required, ForceNew) - The ID of the database instance.
* `node_id` - (Required, ForceNew) - The ID of the node.
* `enable_public_connection` - (Optional) - Whether to enable public network connectivity.
* `enable_private_connection` - (Optional) - Whether to enable private network connectivity.


## Attributes Reference

The following attributes are exported:

* `public_connect_string` - Public connection string.
* `private_connect_string` - Private connection string.
* `public_connect_port` - Public connection port.
* `private_connect_port` - Private connection port.