---
subcategory: "云原生分布式数据库PolarDB-X 2.0"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardbx_databases"
sidebar_current: "docs-Alibabacloudstack-datasource-polardbx-databases"
description: |-
  提供由阿里云账户拥有的 polardbx 数据库列表。
---

# alibabacloudstack\_polardbx\_databases

此数据源根据指定的过滤条件提供阿里云账户中的 polardbx 数据库列表

## 示例用法
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

resource "alibabacloudstack_polardbx_database" "default0" {
    instance_id  = "pxc-unrhxvdzmkvtxg"
	database_name = "${var.name}0"
	encode = "utf8mb4"
	mode = "auto"
}

data "alibabacloudstack_polardbx_databases" "default" {
}
```

## 参数参考

支持以下参数：
  * `database_names` - (可选) - 用于过滤结果的数据库名称列表。
  * `instance_id` - (必填) - polardbx 数据库实例的 ID。
  * `database_name` - (可选) - 数据库的名称。

## 属性参考

除上述参数外，还导出以下属性：
  * `databases` - 数据库列表。
    * `encode` - 字符集。更多信息请参见[字符集表](~~ 99716 ~~)。
    * `description` - 数据库的描述。
    * `instance_id` - polardbx 数据库实例的 ID。
    * `database_name` - 数据库的名称。
    * `accounts` - 数据库的账户信息。
      * `account_name` - 账户名称。
      * `privilege` - 账户的权限。