---
subcategory: "Cloud-Native Distributed Database PolarDB-X 2.0"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardbx_databases"
description: |-
  Provides a list of polardbx databases owned by an alibabacloudstack account.
---

# alibabacloudstack\_polardbx\_databases

This data source provides a list of polardbx databases in an alibabacloudstack account according to the specified filters.

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

resource "alibabacloudstack_polardbx_database" "default0" {
  instance_id  = "${alibabacloudstack_polardbx_instance.default.id}"
	database_name = "${var.name}0"
	encode = "utf8mb4"
	mode = "auto"
}
```

## Argument Reference

The following arguments are supported:
  * `database_names` - (Optional) - A list of database names to filter results.
  * `instance_id` - (Required) - The ID of the PolarDBX DBinstance. 
  * `database_name` - (Optional) - The name of the database.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `databases` - A list of databases.
    * `encode` - Character set. For more information, see [Character Set Table](~~ 99716 ~~).
    * `description` - The description of the database.
    * `instance_id` - The ID of the PolarDBX DBinstance.
    * `database_name` - The name of the database.
    * `accounts` - The account information of the database.
      * `account_name` - The name of the account.
      * `privilege` - The permissions of the account.