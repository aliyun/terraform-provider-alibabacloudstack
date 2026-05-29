---
subcategory: "分布式关系型数据库"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_drds_account"
sidebar_current: "docs-Alibabacloudstack-drds-account"
description: |-
  Provides a drds Account resource.
---

# alibabacloudstack\_drds\_account

Provides a drds Account resource.

## 示例用法
```
variable "name" {
	default = "tf_acc_drds_db_11440"
}

variable "existed_drds_instance" {
	default = ""
}

locals {
	create_drds_instance_count = var.existed_drds_instance == "" ? 1: 0
}

variable "instance_series" {
	default = "drds.sn2.4c16g"
}

resource "random_password" "password" {
	count            = 2
	length           = 12
	special          = true
	override_special = "_"
	min_lower        = 1
	min_upper        = 1
	min_numeric      = 1
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

resource "alibabacloudstack_drds_instance" "default" {
	count                = local.create_drds_instance_count
	description          = "${var.name}"
	zone_id              = "${alibabacloudstack_vpc_vswitch.default.availability_zone}"
	instance_series      = "${var.instance_series}"
	instance_charge_type = "PostPaid"
	vswitch_id           = "${alibabacloudstack_vpc_vswitch.default.id}"
	specification        = "drds.sn2.4c16g.8C32G"
}

locals {
	drds_instance_id = var.existed_drds_instance == "" ? alibabacloudstack_drds_instance.default.0.id: var.existed_drds_instance
}

resource "alibabacloudstack_drds_rds_instance" "default" {
	zone_id             = data.alibabacloudstack_zones.default.zones.0.id
	db_instance_storage = "20"
	storage_type        = "local_ssd"
	category            = "HighAvailability"
	db_instance_class   = "rds.mysql.s1.small"
	drds_instance_id    = local.drds_instance_id
}

resource "alibabacloudstack_drds_database" "default" {
	count              = 2
	instance_id        = "${local.drds_instance_id}"
	drds_database_name = "${var.name}_${count.index}"
	password           = "${random_password.password.0.result}"
	rds_instance_ids   = [alibabacloudstack_drds_rds_instance.default.rds_instance_id,]
}

resource "alibabacloudstack_drds_account" "default" {
  db_privileges {
    db_name = "${alibabacloudstack_drds_database.default.0.drds_database_name}"
    privilege = "R"
  }
  
  instance_id = "${local.drds_instance_id}"
  drds_account_name = "${var.name}"
  password = "${random_password.password.0.result}"
  description = "tf_acc_drds_db_11440"
}
```

## 参数参考

支持以下参数：
  * `instance_id` - (必填, 强制新建) - 实例ID。
  * `drds_account_name` - (必填, 强制新建) - 账号名称。
  * `description` - (选填) - 账号备注。高级账号默认为**Created by DRDS**，普通账号无任何备注。备注信息可以在账号管理中自定义修改。
  * `password` - (必填) DRDS 账号的密码。
  * `db_privileges` - (必填) - 数据库权限信息。
    
    * `db_name` - (必填) - 数据库名称。
    
    * `privilege` - (必填) 数据库权限。取值：`R`（只读）、`RW`（读写）、`DDL`（数据定义语言）、`DML`（数据操作语言）。

## 属性参考

除了上述所有参数外，还导出了以下属性：
  * `host` - 可以访问数据库的IP地址。<note>**%**表示任何IP地址都能访问。</note>
  * `account_type` - 账号类型。**0** 表示高级账号。**1** 表示普通账号。

## 导入

DRDS 账号可以使用 instance_id 和 account_name 进行导入，用冒号连接，例如：

```
$ terraform import alibabacloudstack_drds_account.example <instance_id>:<account_name>@%
```
