---
subcategory: "Cloud-Native Distributed Database PolarDB-X 2.0"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardbx_account_database_binding"
description: |-
  Provides a PolarDBX Account Database Binding resource.
---

# alibabacloudstack_polardbx_account_database_binding

Provides a PolarDBX Account Database Binding resource.

## Example Usage
```
variable "name" {
  default = "accdbbind90366"
}

variable "password" {
  default = ""
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

resource "alibabacloudstack_polardbx_instance" "default" {
  description = "testtf1111"
	zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
	engine_version = "5.7"
	storage = "50"
	network_type = "vpc"
	vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
	vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
	cn_node_class = "polarx.x4.medium.2e"
	cn_node_count = "2"
	dn_node_class = "mysql.n4.medium.25"
	dn_node_count = "2"
}

resource "alibabacloudstack_polardbx_account" "default" {
  instance_id  = "${alibabacloudstack_polardbx_instance.default.id}"
	account_name = "${var.name}"
	account_type = "Normal"
	password     = "${var.password}"
	description  = "Normal user"
}

resource "alibabacloudstack_polardbx_database" "default0" {
    instance_id  = "${alibabacloudstack_polardbx_instance.default.id}"
	database_name = "${var.name}0"
	encode = "utf8mb4"
	mode = "auto"
}

resource "alibabacloudstack_polardbx_database" "default1" {
    instance_id  = "${alibabacloudstack_polardbx_instance.default.id}"
	database_name = "${var.name}1"
	encode = "utf8mb4"
	mode = "auto"
}

resource "alibabacloudstack_polardbx_account_database_binding" "default" {
    instance_id  = "${alibabacloudstack_polardbx_instance.default.id}"
	account_name = "${alibabacloudstack_polardbx_account.default.account_name}"
	db_privileges {
        db_name = "${alibabacloudstack_polardbx_database.default0.database_name}"
        privilege = "ReadOnly"
    }
    db_privileges {
        db_name = "${alibabacloudstack_polardbx_database.default1.database_name}"
        privilege = "ReadWrite"
    }
}

```

## Argument Reference

The following arguments are supported:
  * `account_name` - (Required) - The account name
  * `instance_id` - (Required) -  The ID of the PolarDBX instance.
  * `db_privileges` - (Required) - The Database permissions of the account.
    * `db_name` - (Required) - The name of the database.
    * `privilege` - (Required) - The permissions of the target account on the database. Valid values: ReadOnly, ReadWrite, DMLOnly, DDLOnly.

## Import

PolarDBX Account Database Binding can be imported using the instance_id and account_name, e.g.

```
$ terraform import alibabacloudstack_polardbx_account_database_binding.example pxc-12345678:myaccount
```

