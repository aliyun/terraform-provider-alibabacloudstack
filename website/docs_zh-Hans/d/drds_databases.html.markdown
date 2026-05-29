---
subcategory: "分布式关系型数据库"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_drds_databases"
sidebar_current: "docs-Alibabacloudstack-datasource-drds-databases"
description: |-
  提供阿里云账号下拥有的drds databases列表。
---

# alibabacloudstack\_drds\_databases

此数据源提供根据指定过滤条件列出的阿里云账号下的drds databases资源列表。

## 示例用法
```
variable "name" {
	default = "tf_acc_drds_db_19277"
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
	count            = 1
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
	description          = var.name
	zone_id              = alibabacloudstack_vpc_vswitch.default.availability_zone
	instance_series      = var.instance_series
	instance_charge_type = "PostPaid"
	vswitch_id           = alibabacloudstack_vpc_vswitch.default.id
	specification        = "drds.sn2.4c16g.8C32G"
}

locals {
	drds_instance_id = var.existed_drds_instance == "" ? alibabacloudstack_drds_instance.default.0.id : var.existed_drds_instance
}

resource "alibabacloudstack_drds_rds_instance" "default" {
	count               = 1
	zone_id             = data.alibabacloudstack_zones.default.zones.0.id
	db_instance_storage = "20"
	storage_type        = "local_ssd"
	category            = "HighAvailability"
	db_instance_class   = "rds.mysql.s1.small"
	drds_instance_id    = local.drds_instance_id
}

resource "alibabacloudstack_drds_database" "default" {
	instance_id        = local.drds_instance_id
	drds_database_name = var.name
	password = random_password.password.0.result
	rds_instance_ids = [
		"${alibabacloudstack_drds_rds_instance.default.0.rds_instance_id}",
	]
}
	
data "alibabacloudstack_drds_databases" "default" {
  instance_id = "${alibabacloudstack_drds_database.default.instance_id}"
}
```

## 参数参考
以下参数是支持的：
  * `instance_id` - (必填) - 实例id。
  * `drds_database_name` - (选填) - 数据库名称。

## Attributes Reference
除了上述参数外，还导出以下属性：
  * `drds_database_names` - 数据库名称列表。
  * `databases` - 数据库列表。
    * `create_time` - 数据库创建时间。
    * `drds_database_name` - 数据库名称。
    * `split_mode` - 数据库拆分模式。水平拆分: HORIZONTAL或垂直拆分: VERTICAL。
    * `status` - 数据库状态。
    * `storage_type` - 数据库存储模式。
