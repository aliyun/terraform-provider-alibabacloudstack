---
subcategory: "云数据库 MongoDB" 
layout: "alibabacloudstack" 
page_title: "Alibabacloudstack: alibabacloudstack_mongodb_accounts" 
sidebar_current: "docs-Alibabacloudstack-datasource-mongodb-accounts"
description: |- 
提供一个由 Alibabacloudstack 账户拥有的 MongoDB 账号列表。
---
# alibabacloudstack\_mongodb\_accounts
该数据源根据指定的过滤条件提供 Alibabacloudstack 账户下的 MongoDB 账号列表。

## 示例用法
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
	account_password = "xxxxxx"
	instance_id = "${alibabacloudstack_mongodb_instance.default.id}"
}

data "alibabacloudstack_mongodb_accounts" "default" {
	ids = ["${alibabacloudstack_mongodb_account.default.id}"]
   instance_id = "${alibabacloudstack_mongodb_instance.default.id}"
}
```
## 参数说明
支持以下参数：

* `ids` -（可选）- 用于过滤结果的账号名称 ID 列表。
* `account_name_regex` -（可选）- 通过名称进行过滤的正则表达式模式。
* `instance_id` -（必填）- 实例 ID。
导出属性
除了上述列出的参数外，还导出了以下属性：

* `accounts` - 账号列表。每个元素包含以下属性：
* `id` - 账号的唯一标识。
* `account_name` - 账号名称。
* `character_type` - 账户的角色类型。
* `status` - 账号状态。
* `account_type` - 账号类型。