---
subcategory: "Redis"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_redis_backup"
sidebar_current: "docs-Alibabacloudstack-redis-backup"
description: |-
  Provides a redis Backup resource.
---

# alibabacloudstack\_redis\_backup

Provides a redis Backup resource.

## Example Usage
```
variable "name" {
		default = "tfredis_backup77672"
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



resource "alibabacloudstack_db_instance" "default" {
	engine               = "MySQL"
	engine_version       = "5.6"
	instance_type        = "rds.mysql.s2.large"
	instance_storage     = "20"
	instance_name        = "${var.name}"
	vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
	storage_type         = "local_ssd"
  }



resource "alibabacloudstack_redis_backup" "default" {
  backup_method = "Physical"
  instance_id = "${alibabacloudstack_db_instance.default.id}"
}
```

## Argument Reference

The following arguments are supported:
  * `backup_download_url` - (Optional) - the backup download url.
  * `backup_method` - (Required) - Backup method.
  * `instance_id` - (Required) - InstanceId

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `backup_db_names` - the backup database names.
  * `backup_download_url` - the backup download url.
  * `backup_id` - Backup ID.
  * `backup_intranet_download_url` - the backup intranet download url.
  * `backup_mode` - Backup mode.
  * `backup_size` - Backup size.
  * `backup_type` - Backup type.
  * `end_time` - The query end time must be later than the query start time in the format <I> yyyy-MM-dd</I> T <I> HH:mm</I> Z(UTC time).
  * `start_time` - The query start time, in the format <I> yyyy-MM-dd</I> T <I> HH:mm</I> Z(UTC time).
  * `status` - Backup status.
