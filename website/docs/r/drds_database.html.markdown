---
subcategory: "Distributed Relational Database Service(DRDS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_drds_database"
sidebar_current: "docs-Alibabacloudstack-resource-drds-database"
description: |-
  Provides a drds Database resource.
---

# alibabacloudstack\_drds\_database

Provides a drds Database resource.

## Example Usage
```
variable "name" {
	default = "tf_acc_drds_db_10040"
}

variable "existed_drds_instance" {
	default = ""
}

locals {
	create_drds_instance_count = var.existed_drds_instance == "" ? 1: 0
}

variable "instance_series" {
	default = "drds.sn2.4c16g"
}

resource "random_password" "password" {
	count            = 2
	length           = 12
	special          = true
	override_special = "_"
	min_lower        = 1
	min_upper        = 1
	min_numeric      = 1
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

resource "alibabacloudstack_drds_instance" "default" {
	count                = local.create_drds_instance_count
	description          = "${var.name}"
	zone_id              = "${alibabacloudstack_vpc_vswitch.default.availability_zone}"
	instance_series      = "${var.instance_series}"
	instance_charge_type = "PostPaid"
	vswitch_id           = "${alibabacloudstack_vpc_vswitch.default.id}"
	specification        = "drds.sn2.4c16g.8C32G"
}

locals {
	drds_instance_id = var.existed_drds_instance == "" ? alibabacloudstack_drds_instance.default.0.id: var.existed_drds_instance
}

resource "alibabacloudstack_drds_rds_instance" "default" {
	count               = 3
	zone_id             = data.alibabacloudstack_zones.default.zones.0.id
	db_instance_storage = "20"
	storage_type        = "local_ssd"
	category            = "HighAvailability"
	db_instance_class   = "rds.mysql.s1.small"
	drds_instance_id    = local.drds_instance_id
}

resource "alibabacloudstack_drds_database" "default" {
  instance_id = "${local.drds_instance_id}"
  drds_database_name = "tf_acc_drds_db_10040"
  password = "${random_password.password.0.result}"
  rds_instance_ids = [
                       "${alibabacloudstack_drds_rds_instance.default.0.rds_instance_id}",
                       "${alibabacloudstack_drds_rds_instance.default.1.rds_instance_id}"
                     ]
  ip_white_list = {
                    test1 = "127.0.0.1,192.168.1.1"
                  }
}
```

## Argument Reference

The following arguments are supported:

**Required Arguments:**
  * `instance_id` - (Required) The ID of the DRDS instance.
  * `drds_database_name` - (Required) The name of the DRDS database. The name must be 1 to 24 characters in length, can contain lowercase letters, digits, and underscores (_), and must start with a letter.
  * `password` - (Required) The password of the DRDS database. The password must be 8 to 30 characters in length. This attribute is sensitive.
  * `rds_instance_ids` - (Required) The list of RDS instance IDs. At least one RDS instance must be specified.

**Optional Arguments:**
  * `encode` - (Optional) The character set encoding of the database. Default value: `utf8`.
  * `ip_white_list` - (Optional) The IP whitelist of the database. The key is the group name and the value is a comma-separated list of IP addresses.
  * `split_mode` - (Optional) The database split mode. Valid values: `HORIZONTAL` (horizontal split) or `VERTICAL` (vertical split). Default value: `HORIZONTAL`.
  * `storage_type` - (Optional) The storage type of the database. Valid values: `RDS` or other storage types. Default value: `RDS`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `id` - The unique identifier of the resource, formatted as `<instance_id>:<drds_database_name>`.
  * `create_time` - The creation timestamp of the database in ISO 8601 format.
  * `status` - The status of the DRDS database.

## Import

DRDS Database can be imported using the combination of `instance_id` and `drds_database_name`, formatted as `<instance_id>:<drds_database_name>`, e.g.

```
$ terraform import alibabacloudstack_drds_database.example drds-abc123:my_database
```
