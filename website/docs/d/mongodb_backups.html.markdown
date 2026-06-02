---
subcategory: "ApsaraDB for MongoDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_mongodb_backups"
description: |-
  Provides a list of MongoDB backups owned by an Alibabacloudstack account.
---

# alibabacloudstack\_mongodb\_backups

This data source provides a list of MongoDB backups in an Alibabacloudstack account according to the specified filters.

## Example Usage

```hcl
variable "name" {
	default = "tf-testAlibabacloudstackMongodbBackups39113"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  name = "${var.name}_vsw"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block = "172.16.1.0/24"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

resource "alibabacloudstack_mongodb_instance" "default" {
	vswitch_id           = alibabacloudstack_vpc_vswitch.default.id
	engine_version       = "3.0"
	db_instance_class    = "dds.mongo.mid"
	db_instance_storage  = "10"
	name                 = "${var.name}"
	storage_engine       = "WiredTiger"
	instance_charge_type = "PostPaid"
	replication_factor   = "3"
}

resource "alibabacloudstack_mongodb_backup" "default" {
	backup_method  = "Physical"
	db_instance_id = alibabacloudstack_mongodb_instance.default.id
}

data "alibabacloudstack_mongodb_backups" "default" {
	ids            = ["${alibabacloudstack_mongodb_backup.default.id}"]
	db_instance_id = "${alibabacloudstack_mongodb_backup.default.db_instance_id}"
	start_time     = "${alibabacloudstack_mongodb_backup.default.start_time}"
	end_time       = "${alibabacloudstack_mongodb_backup.default.end_time}"
}
```

## Argument Reference

The following arguments are supported:

* `db_instance_id` - (Required) The ID of the MongoDB instance.
* `start_time` - (Required) The start time of the backup query, in the format `YYYY-MM-DDTHH:MMZ`.
* `end_time` - (Required) The end time of the backup query, in the format `YYYY-MM-DDTHH:MMZ`.
* `backup_id` - (Optional) The ID of the backup.
* `ids` - (Optional) A list of backup IDs.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of backup IDs.
* `backups` - A list of backups. Each element contains the following attributes:
  * `id` - The ID of the backup, formatted as `<backup_id>&<db_instance_id>&<start_time>`.
  * `backup_id` - The backup ID.
  * `backup_db_names` - The names of the backed-up databases.
  * `backup_download_url` - The backup download URL.
  * `backup_intranet_download_url` - The backup intranet download URL.
  * `backup_method` - The backup method.
  * `backup_mode` - The backup mode.
  * `backup_size` - The backup size.
  * `backup_type` - The backup type.
  * `status` - The status of the backup.
  * `db_instance_id` - The MongoDB instance ID.
  * `start_time` - The backup start time.
  * `end_time` - The backup end time.
