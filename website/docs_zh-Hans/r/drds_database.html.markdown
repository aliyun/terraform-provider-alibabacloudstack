---
subcategory: "DRDS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_drds_database"
sidebar_current: "docs-Alibabacloudstack-drds-database"
description: |-
  Provides a drds Database resource.
---

# alibabacloudstack\_drds\_database

Provides a drds Database resource.

## 示例用法
```
variable "name" {
	default = "tf_acc_drds_db_10040"
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
	count               = 3
	zone_id             = data.alibabacloudstack_zones.default.zones.0.id
	db_instance_storage = "20"
	storage_type        = "local_ssd"
	category            = "HighAvailability"
	db_instance_class   = "rds.mysql.s1.small"
	drds_instance_id    = local.drds_instance_id
}

resource "alibabacloudstack_drds_database" "default" {
  instance_id = "${local.drds_instance_id}"
  drds_database_name = "tf_acc_drds_db_10040"
  password = "${random_password.password.0.result}"
  rds_instance_ids = [
                       "${alibabacloudstack_drds_rds_instance.default.0.rds_instance_id}",
                       "${alibabacloudstack_drds_rds_instance.default.1.rds_instance_id}"
                     ]
  ip_white_list = {
                    test1 = "127.0.0.1,192.168.1.1"
                  }
}
```

## 参数参考

支持以下参数：
  * `instance_id` - (必填) - 实例id。
  * `drds_database_name` - (必填) - 数据库名称。
  * `split_mode` - (选填) - 数据库拆分模式。水平拆分: HORIZONTAL或垂直拆分: VERTICAL。
  * `encode` - (选填) - 数据库编码。
  * `password` - (必填) - 密码
  * `rds_instance_ids` - (必填) - 数据库实例id列表。
  * `ip_white_list` - (选填) - 数据库白名单列表。
  * `storage_type` - (选填) - 数据库存储模式。

## 属性参考

除了上述所有参数外，还导出了以下属性：
  * `create_time` - 数据库创建时间。
  * `status` - 数据库状态。
