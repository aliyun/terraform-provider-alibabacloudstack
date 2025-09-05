---
subcategory: "PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_databases"
sidebar_current: "docs-Alibabacloudstack-datasource-polardb-databases"
description: |-
  Provides a list of polardb databases owned by an alibabacloudstack account.
---

# alibabacloudstack\_polardb\_databases

This data source provides a list of polardb databases in an alibabacloudstack account according to the specified filters.

## Example Usage
```
variable "name" {
	default = "tf-testAccPolardbBackups14307"
}

data "alibabacloudstack_zones" default {
  available_resource_creation = "VSwitch"
  enable_details = true
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
resource "alibabacloudstack_polardb_database" "default" {
	data_base_instance_id = "${alibabacloudstack_polardb_dbinstance.instance.id}"
	data_base_description = "Acc Test"
	data_base_name        = "tftest"
	character_set_name = "utf8"
}

data "alibabacloudstack_polardb_databases" "default" {
  ids = [
          "${alibabacloudstack_polardb_database.default.id}"
        ]
  data_base_instance_id = "${alibabacloudstack_polardb_dbinstance.instance.id}"
}
```

## Argument Reference

The following arguments are supported:
  * `ids` - (Optional) - A list of database IDs to filter results.
  * `data_base_instance_id` - (Required) - The ID of the PolarDB DBinstance. 
  * `data_base_name` - (Optional) - The name of the database.
  * `name_regex` - (Optional) - A name Regex of Database.
  * `description_regex` - (Optional) - A description Regex of Database.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `databases` - A list of databases.
    * `id` - The ID of the database.
    * `accounts` - Database account information details.> when the cluster is a PolarDB MySQL engine, it does not include a highly privileged account.
    * `character_set_name` - Character set. For more information, see [Character Set Table](~~ 99716 ~~).
    * `data_base_description` - The description of the database.
    * `data_base_instance_id` - The ID of the PolarDB DBinstance.
    * `data_base_name` - The name of the database.
    * `engine` - The database engine type. The value range is as follows:* **MySQL*** **Oracle*** **PostgreSQL**
    * `page_number` - Page number.
    * `page_size` - The number of records per page. The value range is as follows:* **30 * ** **50 * ** **100 * *The default value is **30 * *.
    * `status` - The status of the resource
