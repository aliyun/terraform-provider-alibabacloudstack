---
subcategory: "MongoDB"  
layout: "alibabacloudstack"  
page_title: "Alibabacloudstack: alibabacloudstack_mongodb_instance"  
sidebar_current: "docs-alibabacloudstack-resource-mongodb-instance"  
description: |-  
  Provides a MongoDB instance resource supports replica set instances only. the MongoDB provides stable, reliable, and automatic scalable database services. It offers a full range of database solutions, such as disaster recovery, backup, recovery, monitoring, and alarms.  
---

# alibabacloudstack_mongodb_instance

Provides a MongoDB instance resource that supports replica set instances only. The MongoDB service provides stable, reliable, and automatically scalable database services. It offers a full range of database solutions, such as disaster recovery, backup, recovery, monitoring, and alarms.

## Example Usage

Basic usage:

```hcl
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "MongoDB"
}

variable "password" {
}

resource "alibabacloudstack_mongodb_instance" "example" {
  engine_version      = "3.4"
  db_instance_class   = "dds.mongo.mid"
  db_instance_storage = 10
  zone_id             = data.alibabacloudstack_zones.default.zones[0].id
  backup_period       = ["Monday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"]
  preferred_backup_time = "20:00Z-21:00Z"
  name                = "testMongoDB"
  security_ip_list    = ["10.168.1.12", "100.69.7.112"] 
  ssl_action          = "Open"
  tde_status          = "Enabled"
  replication_factor  = 3
  storage_engine      = "WiredTiger"
  instance_charge_type= "PostPaid"
  vswitch_id          = "vsw-abc123"
  security_group_id   = "sg-abc123"
  account_password    = var.password
}
```

## Argument Reference

The following arguments are supported:

* `engine_version` - (Required, ForceNew) The database engine version of the instance. Valid values include: `3.4`, `4.0`, etc.
* `db_instance_class` - (Required) The instance type. For example: `dds.mongo.s.small`, `dds.mongo.mid`.
* `db_instance_storage` - (Required) The storage capacity of the instance. Valid values: 10 to 3000. The value must be a multiple of 10. Unit: GB.
* `zone_id` - (Optional, ForceNew) The zone ID of the instance. If not specified, the system will select one by default.
* `backup_period` - (Optional) The backup period for the MongoDB instance. It is required when `preferred_backup_time` is set. Valid values: `[Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday]`. Default: `[Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday]`.
* `preferred_backup_time` - (Optional) The backup time window for the MongoDB instance. In the format of `HH:mmZ-HH:mmZ`. Time setting interval is one hour. If not set, the system will return a default, like `23:00Z-24:00Z`.
* `name` - (Optional) The name of the DB instance. It's a string of 2 to 256 characters.
* `security_ip_list` - (Optional) A list of IP addresses that are allowed to access the MongoDB instance. Each IP address can have up to 256 characters. Defaults to an empty list.
* `ssl_action` - (Optional) Actions performed on SSL functions. Valid values: `Open`: turn on SSL encryption; `Close`: turn off SSL encryption; `Update`: update SSL certificate.
* `tde_status` - (Optional, ForceNew) The Transparent Data Encryption (TDE) status. Valid values: `Enabled`, `Disabled`.
* `replication_factor` - (Optional) Number of replica set nodes. Valid values: `1`, `3`, `5`, `7`. Default: `3`.
* `storage_engine` - (Optional, ForceNew) Storage engine for the instance. Valid values: `WiredTiger`, `RocksDB`. System default: `WiredTiger`.
* `instance_charge_type` - (Optional) The charge type of the instance. Valid values: `PrePaid`, `PostPaid`. Default: `PostPaid`.
* `period` - (Optional) The duration that you will buy the DB instance (in months). It is valid when `instance_charge_type` is `PrePaid`. Valid values: `[1~9], 12, 24, 36`. Default: `1`.
* `vswitch_id` - (Optional, ForceNew) The virtual switch ID to launch DB instances in one VPC.
* `security_group_id` - (Optional) The Security Group ID of ECS. One instance can bind up to 10 ECS security groups.
* `account_password` - (Optional) Password of the root account. It is a string of 6 to 32 characters and is composed of letters, numbers, and underlines.
* `kms_encrypted_password` - (Optional) An KMS encrypted password used to create or update the instance. If `account_password` is provided, this field will be ignored.
* `kms_encryption_context` - (Optional) An KMS encryption context used to decrypt `kms_encrypted_password` before creating or updating the instance.
* `maintain_start_time` - (Optional) The start time of the maintenance window. Specify the time in the `HH:mmZ` format. The time must be in UTC.
* `maintain_end_time` - (Optional) The end time of the maintenance window. Specify the time in the `HH:mmZ` format. The time must be in UTC.
* `tags` - (Optional, Map) A mapping of tags to assign to the resource.
* `db_instance_description` - (Optional) The description of the DB instance.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `retention_period` - Instance log backup retention days.
* `replica_set_name` - The name of the mongo replica set.
* `maintain_start_time` - The start time of the maintenance window.
* `maintain_end_time` - The end time of the maintenance window.
* `ssl_status` - Status of the SSL feature. `Open`: SSL is turned on; `Closed`: SSL is turned off.
* `zone_id` - The zone ID of the instance.
* `vswitch_id` - The virtual switch ID to launch DB instances in one VPC.
* `security_ip_list` - A list of IP addresses that are allowed to access the MongoDB instance.
* `security_group_id` - The Security Group ID of ECS.
* `backup_period` - MongoDB Instance backup period.