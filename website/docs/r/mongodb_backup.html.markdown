---
subcategory: "MongoDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_mongodb_backup"
sidebar_current: "docs-Alibabacloudstack-resource-mongodb-backup"
description: |-
  Provides a mongodb Backup resource.
---

# alibabacloudstack\_mongodb\_backup

Provides a mongodb Backup resource.

## Example Usage
```
variable "name" {
		default = "tfmongodb_backup24637"
	}


data "alibabacloudstack_zones" default {
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
		vswitch_id          = alibabacloudstack_vpc_vswitch.default.id
		engine_version      = "3.0"
		db_instance_class   = "dds.mongo.mid"
		db_instance_storage = "10"
		name                = "${var.name}"
		storage_engine      = "WiredTiger"
		instance_charge_type = "PostPaid"
		replication_factor = "3"
	  }



resource "alibabacloudstack_mongodb_backup" "default" {
  backup_method = "Physical"
  db_instance_id = "${alibabacloudstack_mongodb_instance.default.id}"
}
```

## Argument Reference

The following arguments are supported:
  * `backup_method` - (Required, ForceNew) - Backup Method
  * `db_instance_id` - (Required, ForceNew) - db instance id

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `backup_db_names` - backup db names
  * `backup_download_url` - backup download url
  * `backup_id` - Backup Id
  * `backup_intranet_download_url` - backup intranet download url
  * `backup_mode` - Backup Mode
  * `backup_size` - Backup Size
  * `backup_type` - Backup Type
  * `end_time` - Backup End Time
  * `start_time` - Backup Start Time
  * `status` - The status of the resource
