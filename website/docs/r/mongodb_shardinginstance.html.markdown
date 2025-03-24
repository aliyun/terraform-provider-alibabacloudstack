---
subcategory: "MongoDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_mongodb_shardinginstance"
sidebar_current: "docs-Alibabacloudstack-mongodb-shardinginstance"
description: |- 
  Provides a mongodb Shardinginstance resource.
---

# alibabacloudstack_mongodb_shardinginstance
-> **NOTE:** Alias name has: `alibabacloudstack_mongodb_sharding_instance`

Provides a mongodb Shardinginstance resource.

## Example Usage

### Create a MongoDB Sharding Instance with VPC Configuration

```hcl
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "MongoDB"
}

variable "password" {
}

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

## Argument Reference

The following arguments are supported:

* `engine_version` - (Required, ForceNew) Database version. Value options can refer to the latest docs [CreateDBInstance](https://www.alibabacloud.com/help/zh/doc-detail/61884.htm) `EngineVersion`.
* `storage_engine` - (Optional, ForceNew) Storage engine type of the instance. Valid values: `WiredTiger`, `RocksDB`. Default value: `WiredTiger`.
* `instance_charge_type` - (Optional, ForceNew) Valid values are `PrePaid` and `PostPaid`. System default is `PostPaid`. **NOTE:** It can be modified from `PostPaid` to `PrePaid` after version v1.141.0.
* `period` - (Optional) The duration that you will buy DB instance (in months). It is valid when `instance_charge_type` is `PrePaid`. Valid values: [1~9], 12, 24, 36. System default is 1.
* `zone_id` - (Optional, ForceNew) The Zone to launch the DB instance. MongoDB sharding instance does not support multiple zones. If it is a multi-zone and `vswitch_id` is specified, the vswitch must be in one of them.
* `vswitch_id` - (Optional, ForceNew) The virtual switch ID to launch DB instances in one VPC.
* `name` - (Optional) The name of the DB instance. It is a string of 2 to 256 characters.
* `db_instance_description` - (Optional) A description of the DB instance. It is a string of 2 to 256 characters.
* `security_group_id` - (Optional) The Security Group ID of ECS.
* `account_password` - (Optional, Sensitive) Password of the root account. It is a string of 6 to 32 characters and is composed of letters, numbers, and underlines.
* `kms_encrypted_password` - (Optional) An KMS encrypts password used to create an instance. If the `account_password` is filled in, this field will be ignored.
* `kms_encryption_context` - (Optional) An KMS encryption context used to decrypt `kms_encrypted_password` before creating or updating an instance with `kms_encrypted_password`. See [Encryption Context](https://www.alibabacloud.com/help/doc-detail/42975.htm). It is valid when `kms_encrypted_password` is set.
* `tde_status` - (Optional, ForceNew) The TDE (Transparent Data Encryption) status. Valid values: `Enabled`, `Disabled`.
* `backup_time` - (Optional) MongoDB instance backup time. In the format of HH:mmZ- HH:mmZ. Time setting interval is one hour. If not set, the system will return a default, like "23:00Z-24:00Z".
* `preferred_backup_time` - (Optional) Backup time in the format of HH:mmZ-HH:mmZ (UTC time).
* `shard_list` - (Required) The list of shard nodes. Each shard node has the following properties:
  * `node_class` - (Required) Node specification. See [Instance specifications](https://www.alibabacloud.com/help/doc-detail/57141.htm).
  * `node_storage` - (Required) Custom storage space; value range: [10, 1,000] in 10-GB increments. Unit: GB.
  * `readonly_replicas` - (Optional) The number of read-only nodes in shard node. Valid values: 0 to 5. Default value: 0.
* `mongo_list` - (Required) The list of mongo nodes. Each mongo node has the following properties:
  * `node_class` - (Required) Node specification. See [Instance specifications](https://www.alibabacloud.com/help/doc-detail/57141.htm).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the MongoDB instance.
* `mongo_list` - The list of mongo nodes. Each mongo node contains the following properties:
  * `node_id` - The ID of the mongo node.
  * `connect_string` - Mongo node connection string.
  * `port` - Mongo node port.
* `shard_list` - The list of shard nodes. Each shard node contains the following properties:
  * `node_id` - The ID of the shard node.
* `retention_period` - Instance log backup retention days.
* `config_server_list` - The node information list of config server. Each config server node contains the following properties:
  * `max_iops` - The maximum IOPS of the Config Server node.
  * `connect_string` - The connection address of the Config Server node.
  * `node_class` - The node class of the Config Server node.
  * `max_connections` - The max connections of the Config Server node.
  * `port` - The connection port of the Config Server node.
  * `node_description` - The description of the Config Server node.
  * `node_id` - The ID of the Config Server node.
  * `node_storage` - The node storage of the Config Server node.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 mins) Used when creating the MongoDB instance (until it reaches the initial `Running` status).
* `update` - (Defaults to 30 mins) Used when updating the MongoDB instance (until it reaches the initial `Running` status).
* `delete` - (Defaults to 30 mins) Used when terminating the MongoDB instance.

## Import

MongoDB can be imported using the id, e.g.

```bash
$ terraform import alibabacloudstack_mongodb_sharding_instance.example dds-bp1291daeda44195
```