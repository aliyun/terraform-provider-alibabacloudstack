---
subcategory: "云数据库 MongoDB 版" 
layout: "alibabacloudstack" 
page_title: "Alibabacloudstack: 
alibabacloudstack_mongodb_account" 
sidebar_current: "docs-Alibabacloudstack-mongodb-account"
description: |- 
提供一个 MongoDB 账号资源。
---

# alibabacloudstack_mongodb_account
提供一个 MongoDB 账号资源。

### 示例用法
```hcl
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
  account_password = "xxxxxxxx"
  instance_id = "${alibabacloudstack_mongodb_instance.default.id}"
  account_name = "tfaccount44135"
}
```
## 参数说明
以下参数被支持：

* `account_name` - (必填, 强制新资源) - 账号名称
* `account_password` - (必填) - 账号密码
* `instance_id` - (必填, 强制新资源) - 实例 ID
导出属性
除了上述列出的参数外，还导出了以下属性：

* `character_type` - 账户的角色类型
* `status` - 账户状态
* `account_type` - 账户类型