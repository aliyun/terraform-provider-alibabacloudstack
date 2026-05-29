---
subcategory: "Distributed Relational Database Service(DRDS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_drds_account"
sidebar_current: "docs-Alibabacloudstack-drds-account"
description: |-
  Provides a drds Account resource.
---

# alibabacloudstack\_drds\_account

Provides a drds Account resource.

## Example Usage
```
variable "name" {
	default = "tf_acc_drds_db_11440"
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
	zone_id             = data.alibabacloudstack_zones.default.zones.0.id
	db_instance_storage = "20"
	storage_type        = "local_ssd"
	category            = "HighAvailability"
	db_instance_class   = "rds.mysql.s1.small"
	drds_instance_id    = local.drds_instance_id
}

resource "alibabacloudstack_drds_database" "default" {
	count              = 2
	instance_id        = "${local.drds_instance_id}"
	drds_database_name = "${var.name}_${count.index}"
	password           = "${random_password.password.0.result}"
	rds_instance_ids   = [alibabacloudstack_drds_rds_instance.default.rds_instance_id,]
}

resource "alibabacloudstack_drds_account" "default" {
  db_privileges {
    db_name = "${alibabacloudstack_drds_database.default.0.drds_database_name}"
    privilege = "R"
  }
  
  instance_id = "${local.drds_instance_id}"
  drds_account_name = "${var.name}"
  password = "${random_password.password.0.result}"
  description = "tf_acc_drds_db_11440"
}
```

## Argument Reference

The following arguments are supported:
  * `instance_id` - (Required, ForceNew) - The ID of the instance.
  * `drds_account_name` - (Required, ForceNew) - The name of the account.
  * `description` - (Optional) - Account remarks. The default value of the advanced account is **Created by DRDS**, and the normal account does not have any comments. Remarks can be customized in account management.
  * `password` - (Required) The password of the DRDS account.
  * `db_privileges` - (Required) - Database permission information.
    
    * `db_name` - (Required) - The name of the database.
    
    * `privilege` - (Required) The permission of the database. Valid values: `R`, `RW`, `DDL`, `DML`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `host` - You can access the IP address of the database. <note>**%** indicates that any IP address can be accessed. </note>
  * `account_type` - Account type. **0** indicates an advanced account. **1** indicates a common account.

## Import

DRDS Account can be imported using the instance ID and account name joined with a colon, e.g.

```
$ terraform import alibabacloudstack_drds_account.example <instance_id>:<account_name>@%
```
