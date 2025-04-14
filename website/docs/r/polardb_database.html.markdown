---
subcategory: "PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_database"
sidebar_current: "docs-Alibabacloudstack-polardb-database"
description: |-
  Provides a polardb Database resource.
---

# alibabacloudstack_polardb_database

Provides a polardb Database resource.

## Example Usage
```
variable "name" {
		default = "tf-testaccdbdatabase_basic"
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
  cidr_block = "172.16.0.0/24"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
}


	resource "alibabacloudstack_polardb_instance" "instance" {
		engine            = "MySQL"
		engine_version    = "5.7"
		instance_name = "${var.name}"
		db_instance_storage_type= "local_ssd"
		db_instance_storage = 5
		db_instance_class = "rds.mysql.t1.small"
		zone_id= "${data.alibabacloudstack_zones.default.zones.0.id}"
		vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
	}

resource "alibabacloudstack_polardb_database" "default" {
  data_base_name = "tf-testaccdbdatabase_basic"
  character_set_name = "utf8"
  data_base_instance_id = "${alibabacloudstack_polardb_instance.instance.id}"
}
```

## Argument Reference

The following arguments are supported:
  * `accounts` - (Optional) - Database account information details.> when the cluster is a PolarDB MySQL engine, it does not include a highly privileged account.
    
    * `account` - (Optional) - The name of the database account.
    
    * `account_privilege` - (Optional) - Account permissions. The value range is as follows:* **ReadWrite**: read and write* **ReadOnly**: Read-only* **DMLOnly**: only DML is allowed.* **DDLOnly**: only DDL is allowed* **ReadIndex**: Read-only + index
    
    * `account_privilege_detail` - (Optional) - Detailed information about the account privileges.
  * `character_set_name` - (Required) - Character set. For more information, see [Character Set Table](~~ 99716 ~~).
  * `data_base_description` - (Optional) - The description of the database.
  * `data_base_instance_id` - (Required) -  The ID of the PolarDB instance to which the database will be associated.
  * `data_base_name` - (Required) - The name of the database.
  * `engine` - (Optional) - The database engine type. The value range is as follows:* **MySQL*** **Oracle*** **PostgreSQL**
  * `status` - (Optional) - The status of the resource.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `accounts` - Database account information details.> when the cluster is a PolarDB MySQL engine, it does not include a highly privileged account.
    * `account` - The name of the database account.
    * `account_privilege` - Account permissions. The value range is as follows:* **ReadWrite**: read and write* **ReadOnly**: Read-only* **DMLOnly**: only DML is allowed.* **DDLOnly**: only DDL is allowed* **ReadIndex**: Read-only + index
    * `account_privilege_detail` - Detailed information about the account privileges.
  * `data_base_description` - The description of the database.
  * `engine` - The database engine type. The value range is as follows:* **MySQL*** **Oracle*** **PostgreSQL**
  * `status` - The status of the resource
  * `data_base_instance_id` - The ID of the PolarDB instance to which the database will be associated.