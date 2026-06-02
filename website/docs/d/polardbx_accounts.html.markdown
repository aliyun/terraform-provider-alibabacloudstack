---
subcategory: "Cloud-Native Distributed Database PolarDB-X 2.0"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardbx_accounts"
description: |-
  Provides a list of PolarDB-X accounts owned by an alibabacloudstack account.
---

# alibabacloudstack\_polardbx\_accounts

This data source provides a list of polardbx accounts in an alibabacloudstack account according to the specified filters.

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
	instance_id = alibabacloudstack_polardbx_instance.default.id
	account_name = var.name
	password = "${var.password}"
	description = var.name
}

data "alibabacloudstack_polardbx_accounts" "default" {
	instance_id = alibabacloudstack_polardbx_instance.default.id
}
```

## Argument Reference

The following arguments are supported:
  * `names` - (Optional) - A list of account names to filter results.
  * `instance_id` - (Required) - The ID of the PolarDB-X DB instance.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `ids` - A list of account IDs.
  * `names` - A list of account names.
  * `accounts` - A list of accounts. Each element contains the following attributes:
    * `id` - The ID of the account, formatted as `<instance_id>:<account_name>`.
    * `account_name` - The account name.
    * `description` - The account description.
    * `instance_id` - The ID of the PolarDB-X DB instance.
    * `db_privileges` - A list of database privileges granted to the account. Each element contains:
      * `db_name` - The name of the database.
      * `privilege` - The privilege level of the account on the database.
