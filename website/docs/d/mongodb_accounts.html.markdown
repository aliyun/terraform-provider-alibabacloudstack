---
subcategory: "MongoDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_mongodb_accounts"
sidebar_current: "docs-Alibabacloudstack-datasource-mongodb-accounts"
description: |-
  Provides a list of mongodb accounts owned by an alibabacloudstack account.
---

# alibabacloudstack\_mongodb\_accounts

This data source provides a list of mongodb accounts in an alibabacloudstack account according to the specified filters.

## Example Usage
```
variable "name" {
	default = "tf-testAlibabacloudstackMongodbAccounts83697"
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


resource "alibabacloudstack_mongodb_instance" "default" {
	vswitch_id          = alibabacloudstack_vpc_vswitch.default.id
	engine_version      = "3.0"
	db_instance_class   = "dds.mongo.mid"
	db_instance_storage = "10"
	name                = "${var.name}"
	storage_engine      = "WiredTiger"
	instance_charge_type = "PostPaid"
	replication_factor = "3"
  }

resource "alibabacloudstack_mongodb_account" "default" {
	account_name = "testaccountv1"
	account_password = "Y1cuo&67A!hU"
	instance_id = "${alibabacloudstack_mongodb_instance.default.id}"
}

data "alibabacloudstack_mongodb_accounts" "default" {
	ids = ["${alibabacloudstack_mongodb_account.default.id}"]
   instance_id = "${alibabacloudstack_mongodb_instance.default.id}"
}
```

## Argument Reference

The following arguments are supported:
  * `ids` - (Optional) -  A list of account name IDs to filter the results.
  * `account_name_regex` - (Optional) A regex pattern to filter account by name.
  * `instance_id` - (Required) - Instance Id

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `accounts` - the list of accounts. Each element contains the following attributes:
    * `id` - the unique id of the account.
    * `account_name` - Account Name
    * `character_type` - Account Of The Role Type
    * `status` - Account Status
    * `account_type` - the type of the account.
