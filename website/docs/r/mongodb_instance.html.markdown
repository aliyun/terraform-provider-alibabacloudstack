---
subcategory: "MongoDB"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_mongodb_instance"
sidebar_current: "docs-apsarastack-resource-mongodb-instance"
description: |-
  Provides a MongoDB instance resource supports replica set instances only. the MongoDB provides stable, reliable, and automatic scalable database services. It offers a full range of database solutions, such as disaster recovery, backup, recovery, monitoring, and alarms.
---

# apsarastack\_mongodb\_instance

Provides a MongoDB instance resource supports replica set instances only. the MongoDB provides stable, reliable, and automatic scalable database services. It offers a full range of database solutions, such as disaster recovery, backup, recovery, monitoring, and alarms.


## Example Usage

Basic usage

```

data "apsarastack_zones" "default" {
  available_resource_creation = "MongoDB"
}

resource "apsarastack_mongodb_instance" "example" {
  engine_version      = "3.0"
  db_instance_class   = "dds.mongo.mid"
  db_instance_storage = 10
  zone_id=data.apsarastack_zones.default.zones[0].id
  backup_period=["Monday","Wednesday","Thursday","Friday","Saturday","Sunday"]
  backup_time="20:00Z-21:00Z"
  name="testMongoDB"
  audit_policy={
    enable_audit_policy=true
    storage_period=20
  }
  security_ip_list    = ["10.168.1.12", "100.69.7.112"] 
  ssl_action="Open"
  tde_status="Enabled"
}
```

## Argument Reference

The following arguments are supported:

* `engine_version` - (Required, ForceNew) The database engine version of the instance.
* `db_instance_class` - (Required) The instance type.
* `db_instance_storage` - (Required) The storage capacity of the instance. Valid values: 10 to 3000. The value must be a multiple of 10. Unit: GB.
* `zone_id` - (Optional, ForceNew) The zone ID of the instance.
* `backup_period` - (Optional) MongoDB Instance backup period. It is required when backup_time was existed. Valid values: [Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday]. Default to [Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday].
* `backup_time` - (Required, ForceNew) MongoDB instance backup time. It is required when backup_period was existed. In the format of HH:mmZ- HH:mmZ. Time setting interval is one hour. If not set, the system will return a default, like "23:00Z-24:00Z".
* `name` - (Optional) The name of DB instance. It's a string of 2 to 256 characters.
* `audit_policy` - (Optional) Enables or disables the audit log feature or set the log retention period for MongoDB instance.
    * `enable_audit_policy` - (Required) Specifies whether the audit log feature is enabled.
    * `storage_period` - (Optional) The number of days that audit logs are stored. Valid values: 1 to 365 days. Default value: 30 days.
* `security_ip_list` - (Optional) Description of the nat gateway, This description can have a string of 2 to 256 characters, It cannot begin with http:// or https://. Defaults to null.
* `ssl_action` - (Optional) A list of bandwidth packages for the nat gateway. Only support nat gateway created before 00:00 on November 4, 2017.
* `tde_status` - (Optional, ForceNew) The TDE(Transparent Data Encryption) status.
* `new_connection_string` -(Optional) The new connection string.
* `connection_string` - (Optional) The current connection string, which is to be modified.
* `replication_factor` -(Optional) Number of replica set nodes. Valid values: [1, 3, 5, 7]
* `storage_engine` - (Optional,  ForceNew) Storage engine: WiredTiger or RocksDB. System Default value: WiredTiger.
* `instance_charge_type` - (Optional) Valid values are PrePaid, PostPaid, System default to PostPaid.
* `period` - (Optional) The duration that you will buy DB instance (in month). It is valid when instance_charge_type is PrePaid. Valid values: [1~9], 12, 24, 36. System default to 1.
* `vswitch_id` - (Optional, ForceNew) The virtual switch ID to launch DB instances in one VPC.
* `security_group_id` - (Optional) The Security Group ID of ECS.
* `account_password` - (Optional) Password of the root account. It is a string of 6 to 32 characters and is composed of letters, numbers, and underlines.
* `kms_encrypted_password` - (Optional)  An KMS encrypts password used to a instance. If the account_password is filled in, this field will be ignored.
* `kms_encryption_context` - (Optional) An KMS encryption context used to decrypt kms_encrypted_password before creating or updating instance with kms_encrypted_password
* `ssl_action` - (Optional) Actions performed on SSL functions, Valid values: Open: turn on SSL encryption; Close: turn off SSL encryption; Update: update SSL certificate.
* `maintain_start_time` - (Optional) The start time of the maintenance window. Specify the time in the HH:mmZ format. The time must be in UTC.
* `maintain_end_time` - (Optional) The end time of the maintenance window. Specify the time in the HH:mmZ format. The time must be in UTC.

## Attributes Reference

The following attributes are exported:

* `retention_period` - Instance log backup retention days.
* `replica_set_name` - The name of the mongo replica set.
* `maintain_start_time` - The start time of the maintenance window.
* `maintain_end_time` - The end time of the maintenance window.
* `ssl_status` - Status of the SSL feature. Open: SSL is turned on; Closed: SSL is turned off.
* `connection_string` - The current connection string.
* `zone_id` - The zone ID of the instance.
* `vswitch_id` - The virtual switch ID to launch DB instances in one VPC.
* `security_ip_list` - Description of the nat gateway, This description can have a string of 2 to 256 characters, It cannot begin with http:// or https://. Defaults to null.
* `security_group_id` - The Security Group ID of ECS.
* `backup_period` - MongoDB Instance backup period.


