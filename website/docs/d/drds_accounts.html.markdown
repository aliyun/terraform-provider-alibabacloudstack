---
subcategory: "DRDS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_drds_accounts"
sidebar_current: "docs-Alibabacloudstack-datasource-drds-accounts"
description: |-
  Provides a list of drds accounts owned by an alibabacloudstack account.
---

# alibabacloudstack\_drds\_accounts

This data source provides a list of drds accounts in an alibabacloudstack account according to the specified filters.

## Example Usage
```
variable "name" {
	default = "tf_acc_drds_db_11062"
}

variable "instance_series" {
	default = "drds.sn2.4c16g"
}

resource "random_password" "password" {
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
	description          = "${var.name}"
	zone_id              = "${alibabacloudstack_vpc_vswitch.default.availability_zone}"
	instance_series      = "${var.instance_series}"
	instance_charge_type = "PostPaid"
	vswitch_id           = "${alibabacloudstack_vpc_vswitch.default.id}"
	specification        = "drds.sn2.4c16g.8C32G"
}

resource "alibabacloudstack_drds_rds_instance" "default" {
	zone_id             = data.alibabacloudstack_zones.default.zones.0.id
	db_instance_storage = "20"
	storage_type        = "local_ssd"
	category            = "HighAvailability"
	db_instance_class   = "rds.mysql.s1.small"
	drds_instance_id    = alibabacloudstack_drds_instance.default.id
}

resource "alibabacloudstack_drds_database" "default" {
	instance_id        = "${alibabacloudstack_drds_instance.default.id}"
	drds_database_name = "${var.name}_db"
	password           = random_password.password.result
	rds_instance_ids   = [alibabacloudstack_drds_rds_instance.default.rds_instance_id,]
}

resource "alibabacloudstack_drds_account" "default" {
	instance_id       = alibabacloudstack_drds_instance.default.id
	drds_account_name = var.name
	password          = random_password.password.result
	description       = var.name
	db_privileges {
		db_name   = alibabacloudstack_drds_database.default.drds_database_name
		privilege = "R"
	}
}
	
data "alibabacloudstack_drds_databases" "default" {
  instance_id = "${alibabacloudstack_drds_database.default.instance_id}"
}
```

## Argument Reference

The following arguments are supported:
  * `names` - (Optional) - A list of database account names.
  * `account_type` - (Optional) - Account type.-**0** indicates an advanced account.-**1** indicates a common account.
  * `instance_id` - (Required) - The ID of the instance.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `accounts` - The list of accounts.
    * `id` - The ID of the account.
    * `account_type` - Account type.-**0** indicates an advanced account.-**1** indicates a common account.
    * `db_privileges` - Database permission information.
    * `description` - Account remarks. The default value of the advanced account is **Created by DRDS**, and the normal account does not have any comments. Remarks can be customized in account management.
    * `drds_account_name` - The name of the account.
    * `host` - You can access the IP address of the database. <note>**%** indicates that any IP address can be accessed. </note>
    * `instance_id` - The ID of the instance.
