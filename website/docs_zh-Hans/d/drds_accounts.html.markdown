---
subcategory: "DRDS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_drds_accounts"
sidebar_current: "docs-Alibabacloudstack-datasource-drds-accounts"
description: |-
  提供阿里云账号下拥有的drds accounts列表。
---

# alibabacloudstack\_drds\_accounts

此数据源提供根据指定过滤条件列出的阿里云账号下的drds accounts资源列表。

## 示例用法
```
variable "name" {
	default = "tf_acc_drds_db_11062"
}

variable "instance_series" {
	default = "drds.sn2.4c16g"
}

resource "random_password" "password" {
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
	description          = "${var.name}"
	zone_id              = "${alibabacloudstack_vpc_vswitch.default.availability_zone}"
	instance_series      = "${var.instance_series}"
	instance_charge_type = "PostPaid"
	vswitch_id           = "${alibabacloudstack_vpc_vswitch.default.id}"
	specification        = "drds.sn2.4c16g.8C32G"
}

resource "alibabacloudstack_drds_rds_instance" "default" {
	zone_id             = data.alibabacloudstack_zones.default.zones.0.id
	db_instance_storage = "20"
	storage_type        = "local_ssd"
	category            = "HighAvailability"
	db_instance_class   = "rds.mysql.s1.small"
	drds_instance_id    = alibabacloudstack_drds_instance.default.id
}

resource "alibabacloudstack_drds_database" "default" {
	instance_id        = "${alibabacloudstack_drds_instance.default.id}"
	drds_database_name = "${var.name}_db"
	password           = random_password.password.result
	rds_instance_ids   = [alibabacloudstack_drds_rds_instance.default.rds_instance_id,]
}

resource "alibabacloudstack_drds_account" "default" {
	instance_id       = alibabacloudstack_drds_instance.default.id
	drds_account_name = var.name
	password          = random_password.password.result
	description       = var.name
	db_privileges {
		db_name   = alibabacloudstack_drds_database.default.drds_database_name
		privilege = "R"
	}
}
	
data "alibabacloudstack_drds_databases" "default" {
  instance_id = "${alibabacloudstack_drds_database.default.instance_id}"
}

```

## 参数参考
以下参数是支持的：
  * `names` - (选填) - 账号名称列表。
  * `account_type` - (选填) - 账号类型。- **0**表示高级账号。- **1**表示普通账号。
  * `instance_id` - (必填) - 实例ID。

## Attributes Reference
除了上述参数外，还导出以下属性：
  * `accounts` - 账号列表。
    * `id` - 账号ID。
    * `account_type` - 账号类型。- **0**表示高级账号。- **1**表示普通账号。
    * `db_privileges` - 数据库权限信息。
    * `description` - 账号备注。高级账号默认为**Created by DRDS**，普通账号无任何备注。备注信息可以在账号管理中自定义修改。
    * `drds_account_name` - 账号名称。
    * `host` - 可以访问数据库的IP地址。<note>**%**表示任何IP地址都能访问。</note>
    * `instance_id` - 实例ID。
