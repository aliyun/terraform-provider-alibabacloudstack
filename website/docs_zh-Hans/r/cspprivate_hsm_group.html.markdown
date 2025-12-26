---
subcategory: "Cspprivate HSM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cspprivate_hsm_group"
sidebar_current: "docs-Alibabacloudstack-cspprivate-hsm_group"
description: |-
  管理阿里云密码机组资源
---

# alibabacloudstack_cspprivate_hsm_group

管理阿里云密码机组（HSM Group），用于创建和管理硬件安全模块集群。

## 示例用法

### 基础用法

```hcl

variable "name" {
  default = "tf_hsm_group1476"
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

resource "alibabacloudstack_vpc" "vpc1" {
  provider   = alibabacloudstack-common
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16" # VPC CIDR block
}

resource "alibabacloudstack_vswitch" "vsw1" {
  provider          = alibabacloudstack-common
  vpc_id            = alibabacloudstack_vpc.vpc1.id
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

resource "alibabacloudstack_cspprivate_hsm_instance" "default1" {
  product_code   = data.alibabacloudstack_cspprivate_hsm_vendors.default.vendors.0.products.0.code
  vendor_code    = data.alibabacloudstack_cspprivate_hsm_vendors.default.vendors.0.code
  vsm_type       = "gvsm"
  zone_id        = data.alibabacloudstack_zones.default.zones.0.id
  alias_name     = "${var.name}1"
  vpc_id         = alibabacloudstack_vpc.vpc1.id
  vpc_cidr_block = alibabacloudstack_vpc.vpc1.cidr_block
  vswitch_id     = alibabacloudstack_vswitch.vsw1.id
  ip             = "192.168.0.101"
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
  zone_ids = [
    "${alibabacloudstack_cspprivate_hsm_instance.default.zone_id}"
  ]
  password  = random_password.password.0.result
  hsm_count = "3"
  hsm_list = [
    "${alibabacloudstack_cspprivate_hsm_instance.default.id}",
    "${alibabacloudstack_cspprivate_hsm_instance.default1.id}"
  ]
}
```

## 参数说明

支持以下参数：

* `password` - (必填, 变更时重建) 密码机组初始化密码。用于初始化密码机组的安全凭证，长度需符合密码策略要求。
* `hsm_count` - (可选, 变更时重建) HSM实例数量。指定密码机组中HSM实例的个数，最小值为1。变更时需重建整个密码机组。
* `zone_ids` - (可选, 变更时重建) 可用区ID列表。指定HSM实例部署的可用区ID集合，格式为JSON数组字符串（如`["cn-wulan-env17e-amtest17001-a"]`）。变更时需重建整个密码机组。
* `hsm_list` - (可选) HSM实例ID列表。指定要加入密码机组的现有HSM实例ID集合，格式为JSON数组字符串（如`["hsm-xxx"]`）。支持动态更新，无需重建资源。

## 属性说明

以下属性导出为资源属性：

* `id` - 密码机组的ID（与`group_name`相同）。
* `create_time` - 密码机组的创建时间，格式为ISO 8601标准时间字符串。
* `group_name` - 密码机组的名称，由系统自动生成。
* `security_level_tag` - 安全级别标签，表示密码机组的安全等级（如"公开"、"秘密"等）。
* `status` - 密码机组的当前状态（如"uninitialized"表示未初始化）。
* `update_time` - 密码机组的最后更新时间，格式为ISO 8601标准时间字符串。
* `vpc_id` - 密码机组关联的VPC ID，表示HSM实例部署的网络环境。