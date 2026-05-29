---
subcategory: "PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_cluster_databases"
sidebar_current: "docs-Alibabacloudstack-datasource-polardb-cluster-databases"
description: |-
  Provides a list of PolarDB cluster databases.
---

# alibabacloudstack_polardb_cluster_databases

This data source provides a list of PolarDB cluster databases in an Alibaba Cloud Stack environment.

## Example Usage

```hcl
variable "name" {
	default = "tfAccountsName"
}

variable "password" {
}

data  "alibabacloudstack_zones" "default" {
	available_resource_creation = "PolarDB"
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name       = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}
resource "alibabacloudstack_vpc_vswitch" "default" {
  vpc_id            = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  zone_id           = "${data.alibabacloudstack_zones.default.zones.0.id}"
  vswitch_name      = "${var.name}_vsw"
}

resource "alibabacloudstack_polardb_cluster_instance" "default" {
  db_cluster_description 	=  "${var.name}"
  zone_id 				= "${data.alibabacloudstack_zones.default.zones.0.id}"
  db_type 				= "${var.db_type}"
  db_version 			= "${var.db_version}"
  storage_space 		= "20"
  vpc_id 				= "${alibabacloudstack_vpc_vpc.default.id}"
  vswitch_id			= "${alibabacloudstack_vpc_vswitch.default.id}"
  db_node_class 		= "${data.alibabacloudstack_polardb_cluster_instance_types.default.instance_types.0.id}"
  sub_category 			= "${data.alibabacloudstack_polardb_cluster_instance_types.default.instance_types.0.sub_category}"
  storage_type 			= "ESSDPL1"
}

resource "alibabacloudstack_polardb_cluster_database" "default" {
    db_cluster_id		= "${alibabacloudstack_polardb_cluster_instance.default.id}"
    db_name 			= "${var.name}"
    character_set_name 	= "utf8"
}

data "alibabacloudstack_polardb_cluster_databases" "example" {
  db_cluster_id = "pc-12345678"
  db_name       = "test_database"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of database IDs.
* `db_name` - (Optional) The name of the database.
* `db_cluster_id` - (Required) The ID of the PolarDB cluster.
* `name_regex` - (Optional, Deprecated) A regex string to filter results by database name. **Field 'name_regex' is deprecated and will be removed in a future release. Please use new field 'description_regex' instead.**
* `description_regex` - (Optional) A regex string to filter results by database name.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of database IDs.
* `databases` - A list of PolarDB cluster databases. Each element contains the following attributes:
  * `id` - The ID of the database.
  * `db_name` - The name of the database.
  * `character_set_name` - The character set of the database.
  * `db_description` - The description of the database.
  * `engine` - The database engine.
  * `db_status` - The status of the database.