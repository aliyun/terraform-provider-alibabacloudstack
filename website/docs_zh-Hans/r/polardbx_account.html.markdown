---
subcategory: "PolarDBX"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardbx_account"
sidebar_current: "docs-Alibabacloudstack-PolarDBX-account"
description: |-
  提供 PolarDBX 账户资源。
---

# alibabacloudstack_polardb_account

提供 PolarDBX 账户资源。

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
  instance_id  = "${alibabacloudstack_polardbx_instance.default.id}"
	account_name = "${var.name}"
	account_type = "Normal"
	password     = "${var.password}"
	description  = "Normal user"
}
```


## 参数参考

支持以下参数：
  * `description` - (可选) - 账户备注需满足以下要求：-不能以' http:'或' https'开头。-长度为2到256个字符。
  * `account_name` - (必填) - 账户名称必须满足以下要求：* 以小写字母开头，以字母或数字结尾。* 由小写字母、数字或下划线组成。* 长度为2到16个字符。* 不能使用一些保留用户名，如 root 和 admin。
  * `password` - (必填) - 账户密码
  * `account_type` - (可选) - 账户类型。取值范围如下：-**Normal**：普通账户。-**Super**：高权限账户。
  * `instance_id` - (必填) - 账户将关联的 PolarDBX 实例的 ID。


## 属性参考

除上述参数外，还导出以下属性：
  * `description` - 账户备注需满足以下要求：-不能以' http:'或' https'开头。-长度为2到256个字符。
  * `account_type` - 账户类型。取值范围如下：-**Normal**：普通账户。-**Super**：高权限账户。