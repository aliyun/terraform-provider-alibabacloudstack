---
subcategory: "PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_database"
sidebar_current: "docs-Alibabacloudstack-resource-polardb-database"
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


	resource "alibabacloudstack_polardb_dbinstance" "instance" {
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
  data_base_instance_id = "${alibabacloudstack_polardb_dbinstance.instance.id}"
}
```

## Argument Reference

The following arguments are supported:
  * `character_set_name` - (Required) - Character set. For more information, see [Character Set Table](~~ 99716 ~~).
  * `data_base_instance_id` - (Required/ForceNew) - The ID of the PolarDB instance to which the database will be associated. Modification of this parameter forces a new resource to be created.
  * `data_base_name` - (Required) - The name of the database. Case-insensitive.
  * `data_base_description` - (Optional) - The description of the database.
  * `engine` - (Optional) - The database engine type. Valid values: **MySQL**, **Oracle**, **PostgreSQL**.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `status` - The status of the database.

## Import

PolarDB Database can be imported using the `data_base_instance_id` and `data_base_name` separated by a colon, e.g.

```
$ terraform import alibabacloudstack_polardb_database.example <data_base_instance_id>:<data_base_name>
```