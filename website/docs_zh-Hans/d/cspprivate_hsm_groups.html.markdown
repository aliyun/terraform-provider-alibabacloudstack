---
subcategory: "Cspprivate HSM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cspprivate_hsm_group"
sidebar_current: "docs-alibabacloudstack-datasource-cspprivate-hsm-group"
description: |-
    Provides a data source for CSP Private HSM Group to query HSM group information in Alibaba Cloud CSP Private HSM service.
---

# alibabacloudstack_cspprivate_hsm_group

> 密码机组数据源，用于查询阿里云密码机服务中的密码机组信息

## 示例用法

```hcl

variable "name" {
  default = "tf_hsm_group34787"
}

data "alibabacloudstack_zones" "default" {
  provider                    = alibabacloudstack-common
  available_resource_creation = "VSwitch"
}

data "alibabacloudstack_cspprivate_hsm_vendors" "default" {
}

resource "alibabacloudstack_vpc" "vpc" {
  provider   = alibabacloudstack-common
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16" # VPC CIDR block
}

resource "alibabacloudstack_vswitch" "vsw" {
  provider          = alibabacloudstack-common
  vpc_id            = alibabacloudstack_vpc.vpc.id
  cidr_block        = "192.168.0.0/24" # VSwitch CIDR block
  availability_zone = data.alibabacloudstack_zones.default.zones.0.id
}


resource "alibabacloudstack_cspprivate_hsm_instance" "default" {
  product_code   = data.alibabacloudstack_cspprivate_hsm_vendors.default.vendors.0.products.0.code
  vendor_code    = data.alibabacloudstack_cspprivate_hsm_vendors.default.vendors.0.code
  vsm_type       = "gvsm"
  zone_id        = data.alibabacloudstack_zones.default.zones.0.id
  alias_name     = "${var.name}0"
  vpc_id         = alibabacloudstack_vpc.vpc.id
  vpc_cidr_block = alibabacloudstack_vpc.vpc.cidr_block
  vswitch_id     = alibabacloudstack_vswitch.vsw.id
  ip             = "192.168.0.100"
}


resource "random_password" "password" {
  count            = 1
  length           = 12
  special          = true
  override_special = "!@#$^&*()_"
  min_lower        = 1
  min_upper        = 1
  min_numeric      = 1
}

resource "alibabacloudstack_cspprivate_hsm_group" "default" {
  hsm_list  = ["${alibabacloudstack_cspprivate_hsm_instance.default.id}"]
  zone_ids  = ["${data.alibabacloudstack_zones.default.zones.0.id}"]
  password  = random_password.password.0.result
  hsm_count = 1
}

data "alibabacloudstack_cspprivate_hsm_groups" "default" {
  ids = ["${alibabacloudstack_cspprivate_hsm_group.default.id}"]
}

```

## 参数说明
以下参数用于过滤和查询密码机组：

* `ids` (列表, 可选)：用于按组名过滤结果的密码机组名称列表。如果指定，将只返回名称匹配的密码机组。

## 属性说明
以下属性被导出：

* `id` (字符串)：密码机组的唯一标识符，等同于组名。
* `create_time` (字符串)：密码机组的创建时间，格式为ISO 8601标准时间。
* `group_name` (字符串)：密码机组的名称。
* `hsm_count` (整数)：密码机组中HSM（硬件安全模块）的数量。
* `security_level_tag` (字符串)：密码机组的安全级别标签，如"公开"。
* `status` (字符串)：密码机组的当前状态，如"uninitialized"表示未初始化。
* `unique_id` (整数)：密码机组的唯一ID标识。
* `update_time` (字符串)：密码机组的最后更新时间，格式为ISO 8601标准时间。
* `vpc_id` (字符串)：密码机组所属VPC的ID。
* `zone_ids` (字符串)：密码机组所在可用区ID列表，以JSON数组格式表示。