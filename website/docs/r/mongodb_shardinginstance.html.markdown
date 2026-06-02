---
subcategory: "ApsaraDB for MongoDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_mongodb_shardinginstance"
description: |- 
  Provides a mongodb Shardinginstance resource.
---

# alibabacloudstack_mongodb_shardinginstance
-> **NOTE:** Alias name has: `alibabacloudstack_mongodb_sharding_instance`

Provides a mongodb Shardinginstance resource.

## Example Usage

### Create a MongoDB Sharding Instance with VPC Configuration

```hcl
variable "name" {
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


resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name   = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  name       = "${var.name}_vsw"
  vpc_id     = alibabacloudstack_vpc_vpc.default.id
  cidr_block = "172.16.1.0/24"
  zone_id    = data.alibabacloudstack_zones.default.zones.0.id
}

resource "alibabacloudstack_mongodb_sharding_instance" "default" {
  db_account_password = random_password.password.0.result
  zone_id             = data.alibabacloudstack_zones.default.zones.0.id
  engine_version      = "4.0"
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

  vswitch_id      = alibabacloudstack_vpc_vswitch.default.id
  db_account_name = "tf_testacc"
}
```

## Argument Reference

The following arguments are supported:

* `engine_version` - (Required, ForceNew) Database version. Value options can refer to the latest docs [CreateDBInstance](https://www.alibabacloud.com/help/zh/doc-detail/61884.htm) `EngineVersion`.
* `storage_engine` - (Optional, ForceNew) Storage engine type of the instance. Valid values: `WiredTiger`, `RocksDB`. Default value: `WiredTiger`.
* `zone_id` - (Optional, ForceNew) The Zone to launch the DB instance. MongoDB sharding instance does not support multiple zones. If it is a multi-zone and `vswitch_id` is specified, the vswitch must be in one of them.
* `vswitch_id` - (Optional, ForceNew) The virtual switch ID to launch DB instances in one VPC.
* `name` - (Optional) The name of the DB instance. It is a string of 2 to 256 characters.
* `db_instance_description` - (Optional) A description of the DB instance. It is a string of 2 to 256 characters.
* `cs_root_account_password` - (Optional, Sensitive) Password of the cs node root account. It is a string of 6 to 32 characters and is composed of letters, numbers, and underlines.
* `db_account_name` - (Required) Account name of the shard node.
* `db_account_password` - (Required, Sensitive) Password of the shard node account. It is a string of 6 to 32 characters and is composed of letters, numbers, and underlines.
* `kms_encrypted_password` - (Optional) An KMS encrypts password used to create an instance. If the `account_password` is filled in, this field will be ignored.
* `kms_encryption_context` - (Optional) An KMS encryption context used to decrypt `kms_encrypted_password` before creating or updating an instance with `kms_encrypted_password`. See [Encryption Context](https://www.alibabacloud.com/help/doc-detail/42975.htm). It is valid when `kms_encrypted_password` is set.
* `tde_status` - (Optional, ForceNew) The TDE (Transparent Data Encryption) status. Valid values: `Enabled`, `Disabled`.
* `backup_time` - (Optional) MongoDB instance backup time. In the format of HH:mmZ- HH:mmZ. Time setting interval is one hour. If not set, the system will return a default, like "23:00Z-24:00Z".
* `preferred_backup_time` - (Optional) Backup time in the format of HH:mmZ-HH:mmZ (UTC time).
* `shard_list` - (Required) The list of shard nodes. Each shard node has the following properties:
  * `node_class` - (Required) Node specification. See [Instance specifications](https://www.alibabacloud.com/help/doc-detail/57141.htm).
  * `node_storage` - (Required) Custom storage space; value range: [10, 1,000] in 10-GB increments. Unit: GB.
  * `description` - (Required) Node name.
* `mongo_list` - (Required) The list of mongo nodes. Each mongo node has the following properties:
  * `node_class` - (Required) Node specification. See [Instance specifications](https://www.alibabacloud.com/help/doc-detail/57141.htm).
  * `description` - (Required) Node name.
* `configserver_list` - (Required) The list of config server nodes. Each shard node has the following properties:
  * `node_class` - (Required) Node specification. See [Instance specifications](https://www.alibabacloud.com/help/doc-detail/57141.htm).
  * `node_storage` - (Required) Custom storage space; value range: [10, 1,000] in 10-GB increments. Unit: GB.
  * `description` - (Required) Node name.
* `security_ip_list` - (Optional) - List of IP addresses allowed to access all databases of an instance. The list contains up to 1,000 IP addresses, separated by commas. Supported formats include 0.0.0.0/0, 10.23.12.24 (IP), and 10.23.12.24/24 (CIDR mode). CIDR group mode can authorize a continuous segment of IP addresses, including private network types authorized by IP segments.

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

* `storage_engine` - (Computed) - Storage engine type of the instance.
* `instance_charge_type` - (Computed) - Charging type of the instance.
* `period` - (Computed) - Duration of the subscription for prepaid instances.
* `zone_id` - (Computed) - Zone where the instance resides.
* `vswitch_id` - (Computed) - Virtual switch ID within a VPC.
* `name` - (Computed) - Name of the DB instance.
* `db_instance_description` - (Computed) - Description of the DB instance.
* `security_group_id` - (Computed) - Security Group ID associated with the instance.
* `tde_status` - (Computed) - Transparent Data Encryption status of the instance.
* `preferred_backup_time` - (Computed) - Preferred backup time window.

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