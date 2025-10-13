---
subcategory: "PolarDBX"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardbx_account_database_binding" 
sidebar_current: "docs-Alibabacloudstack-polardbx-account-database-binding" 
description: |-
    Provides a PolarDBX Account Database Binding resource.
---

# alibabacloudstack_polardbx_account_database_binding
Provides a PolarDBX Account Database Binding resource.

## Example Usage
```

variable "name" {
  default = "YourTestName"
}

data "alibabacloudstack_zones" default {
  available_resource_creation = "VSwitch"
  enable_details = true
}

variable "creation" {
	default = "PolarDB"
}

resource "alibabacloudstack_polardb_dbinstance" "instance" {
	engine            = "MySQL"
	engine_version    = "5.7"
	instance_name = "${var.name}"
	db_instance_storage_type= "local_ssd"
	db_instance_storage = 5
	db_instance_class = "rds.mysql.t1.small"
	zone_id= "${data.alibabacloudstack_zones.default.zones.0.id}"
}

resource "alibabacloudstack_polardb_database" "default0" {
    data_base_instance_id  = "${alibabacloudstack_polardb_dbinstance.instance.id}"
	data_base_name = "${var.name}0"
	character_set_name = "utf8mb4"
}

resource "alibabacloudstack_polardb_database" "default1" {
    data_base_instance_id  = "${alibabacloudstack_polardb_dbinstance.instance.id}"
	data_base_name = "${var.name}1"
	character_set_name = "utf8mb4"
}

resource "alibabacloudstack_polardb_database" "default2" {
    data_base_instance_id  = "${alibabacloudstack_polardb_dbinstance.instance.id}"
	data_base_name = "${var.name}2"
	character_set_name = "utf8mb4"
}

resource "alibabacloudstack_polardbx_account_database_binding" "default" {
    db_instance_id  = "${alibabacloudstack_polardbx_instance.default.id}"
	account_name = "${alibabacloudstack_polardbx_account.default.account_name}"
	database_privileges {
        db_name = "${alibabacloudstack_polardb_database.default1.database_name}"
        privilege = "ReadOnly"
    }
    database_privileges {
        db_name = "${alibabacloudstack_polardb_database.default2.database_name}"
        privilege = "ReadWrite"
    }
}

```

## Argument Reference
The following arguments are supported:

* `account_name` - (Required) - The account name
* `db_instance_id` - (Required) - The ID of the PolarDBX instance.
* `database_privileges` - The Database permissions of the account.
    * `db_name` - The name of the database.
    * `privilege` - The permissions of the target account on the database.