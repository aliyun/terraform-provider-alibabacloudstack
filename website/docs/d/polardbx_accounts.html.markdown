---
subcategory: "PolarDBX"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardbx_accounts"
sidebar_current: "docs-Alibabacloudstack-datasource-polardbx-accounts"
description: |-
  Provides a list of polardbx accounts owned by an alibabacloudstack account.
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

resource "alibabacloudstack_polardbxx_instance" "default" {
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

resource "alibabacloudstack_polardbxx_account" "default" {
	instance_id = alibabacloudstack_polardbxx_instance.default.id
	account_name = var.name
	password = "${var.password}"
	description = var.name
}

data "alibabacloudstack_polardbxx_accounts" "default" {
}
```

## Argument Reference

The following arguments are supported:
  * `ids` - (Optional) - A list of account IDs to filter results.
  * `names` - (Optional) - A list of account names to filter results
  * `account_name` - (Optional) - The account name.
  * `instance_id` - (Required) - db instance id.
  * `account_type` - (Optional) - Account type. The value range is as follows:-**Normal**: Normal account.-**Super**: a highly privileged account.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `accounts` - A list of accounts. Each element contains the following attributes:
    * `id` - The ID of the account.
    * `description` - The account description.
    * `account_name` - The account name
    * `account_type` - Account type. The value range is as follows:-**Normal**: Normal account.-**Super**: a highly privileged account.
    * `instance_id` - The ID of the PolarDBX Db instance.
    * `db_privileges` - The Database permissions of the account.
      * `db_name` - The name of the database.
      * `privilege` - The permissions of the target account on the database.
