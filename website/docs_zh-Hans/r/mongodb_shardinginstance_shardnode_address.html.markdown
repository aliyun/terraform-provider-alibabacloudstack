---
subcategory: "MongoDB" 
layout: "alibabacloudstack" 
page_title: "Alibabacloudstack: alibabacloudstack_mongodb_shardinginstance_shardnode_address" 
sidebar_current: "docs-Alibabacloudstack-shardinginstance-shardnode-address" description: |- 
管理Mongodb分片实例shard节点的网络链接地址。
---

# alibabacloudstack_mongodb_shardinginstance_shardnode_address
管理Mongodb分片实例shard节点的网络链接地址。

## 示例用法
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

## 参数说明
以下参数被支持：

* `db_instance_id` - (必填, 强制新资源) 实例ID
* `node_id` - (必填, 强制新资源) 节点ID
* `enable_public_connection` - (选填) 是否开启公网连接
* `enable_private_connection` - (选填) 是否开启私网连接


## 属性说明
除了上述列出的参数外，还导出了以下属性：

* `public_connect_string` - 公网链接地址
* `private_connect_string` - 私网链接地址
* `public_connect_port` - 公网链接地址端口
* `private_connect_port` - 私网链接端口