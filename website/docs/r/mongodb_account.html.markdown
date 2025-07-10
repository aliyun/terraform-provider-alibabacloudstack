---
subcategory: "MongoDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_mongodb_account"
sidebar_current: "docs-Alibabacloudstack-mongodb-account"
description: |-
  Provides a mongodb Account resource.
---

# alibabacloudstack\_mongodb\_account

Provides a mongodb Account resource.

## Example Usage
```
variable "name" {
	default = "tfaccount44135"
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
  account_password = "YqvqrI2@wGwc"
  instance_id = "${alibabacloudstack_mongodb_instance.default.id}"
  account_name = "tfaccount44135"
}
```

## Argument Reference

The following arguments are supported:
  * `account_name` - (Required, ForceNew) - Account Name
  * `account_password` - (Required) - Account Password
  * `instance_id` - (Required, ForceNew) - Instance Id

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `character_type` - Account Of The Role Type
  * `status` - Account Status
  * `account_type` - the Account Type
