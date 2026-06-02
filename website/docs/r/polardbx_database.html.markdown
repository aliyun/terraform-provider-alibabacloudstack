---
subcategory: "Cloud-Native Distributed Database PolarDB-X 2.0"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardbx_database"
description: |-
  Provide PolarDBX database resources.
---

# alibabacloudstack_polardbx_database

Provide PolarDBX database resources.

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
    instance_id  = "pxc-unrhxvdzmkvtxg"
	database_name = "${var.name}0"
	encode = "utf8mb4"
	mode = "auto"
}

```
## Argument Reference

The following arguments are supported:

  * `encode` - (Required) - Character set. For more information, refer to the [Character Set Table](~~ 99716 ~~).
  * `description` - (Optional) - Description of the database.
  * `instance_id` - (Required) - ID of the PolarDBX instance to which the database will be associated.
  * `database_name` - (Required) - Name of the database.
  * `mode` - (Optional) - Mode of the database. Valid values: `auto` and `drds`.

