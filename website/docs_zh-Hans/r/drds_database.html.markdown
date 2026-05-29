---
subcategory: "云原生分布式数据库PolarDB-X 1.0"
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

**必填参数：**
  * `instance_id` - (必填) DRDS 实例 ID。
  * `drds_database_name` - (必填) 数据库名称。长度为 1-24 个字符，只能包含小写字母、数字和下划线 (_)，且必须以字母开头。
  * `password` - (必填) 数据库密码。长度为 8-30 个字符。该属性是敏感的。
  * `rds_instance_ids` - (必填) RDS 实例 ID 列表。至少需要指定 1 个 RDS 实例。

**可选参数：**
  * `encode` - (选填) 数据库字符集编码。默认为 `utf8`。
  * `ip_white_list` - (选填) 数据库 IP 白名单列表。键为分组名称，值为逗号分隔的 IP 地址。
  * `split_mode` - (选填) 数据库拆分模式。`HORIZONTAL`（水平拆分）或 `VERTICAL`（垂直拆分）。默认为 `HORIZONTAL`。
  * `storage_type` - (选填) 数据库存储类型。`RDS` 或其他存储类型。默认为 `RDS`。

## 属性参考

除了上述所有参数外，还导出了以下属性：
  * `id` - 资源的唯一标识，格式为 `<instance_id>:<drds_database_name>`。
  * `create_time` - 数据库创建时间（ISO 8601 格式）。
  * `status` - 数据库状态。

## Import

DRDS 数据库可以使用 `instance_id` 和 `drds_database_name` 的组合进行导入，格式为 `<instance_id>:<drds_database_name>`，例如：

```
$ terraform import alibabacloudstack_drds_database.example drds-abc123:my_database
```
