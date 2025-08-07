---
subcategory: "PolarDBX"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardbx_account_database_binding"
sidebar_current: "docs-Alibabacloudstack-PolarDBX-account-database-binding"
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
```

## Argument Reference

The following arguments are supported:
  * `account_name` - (Required) - The account name must meet the following requirements:* Start with a lowercase letter and end with a letter or number.* Consists of lowercase letters, numbers, or underscores.* The length is 2 to 16 characters.* You cannot use some reserved usernames, such as root and admin.
  * `account_type` - (Optional) - Account type. The value range is as follows:-**Normal**: Normal account.-**Super**: a highly privileged account.
  * `instance_id` - (Required) -  The ID of the PolarDBX instance.
  * `db_privileges` - The Database permissions of the account.
    * `db_name` - The name of the database.
    * `privilege` - The permissions of the target account on the database.
