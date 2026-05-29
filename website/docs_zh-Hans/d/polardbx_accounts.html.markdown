---
subcategory: "云原生分布式数据库PolarDB-X 2.0"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardbx_accounts"
sidebar_current: "docs-Alibabacloudstack-datasource-polardbx-accounts"
description: |-
  提供阿里云账号下拥有的 PolarDB-X 账户列表。
---

# alibabacloudstack\_polardbx\_accounts

此数据源提供根据指定过滤条件列出的阿里云账号下的polardbx accounts资源列表。

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

resource "alibabacloudstack_polardbx_account" "default" {
	instance_id = alibabacloudstack_polardbx_instance.default.id
	account_name = var.name
	password = "${var.password}"
	description = var.name
}

data "alibabacloudstack_polardbx_accounts" "default" {
	instance_id = alibabacloudstack_polardbx_instance.default.id
}
```

## 参数参考

支持以下参数：
  * `names` - (可选) - 用于过滤结果的账户名称列表。
  * `instance_id` - (必填) - PolarDB-X 数据库实例 ID。

## 属性参考

除上述参数外，还导出以下属性：
  * `ids` - 账户 ID 列表。
  * `names` - 账户名称列表。
  * `accounts` - 账户列表。每个元素包含以下属性：
    * `id` - 账户 ID，格式为 `<instance_id>:<account_name>`。
    * `account_name` - 账户名称。
    * `description` - 账户描述。
    * `instance_id` - PolarDB-X 数据库实例的 ID。
    * `db_privileges` - 账户的数据库权限列表。每个元素包含：
      * `db_name` - 数据库名称。
      * `privilege` - 账户在该数据库上的权限级别。