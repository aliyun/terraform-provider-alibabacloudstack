---
subcategory: "云密码机"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cspprivate_hsm_instance"
sidebar_current: "docs-alibabacloudstack-resource-cspprivate-hsm-instance"
description: |-
  提供阿里云密码服务(HSM)实例资源
---

# alibabacloudstack_cspprivate_hsm_instance

提供密码服务(HSM)实例资源。

## 典型用法

### 基本用法

```hcl
variable "name" {
  default = "tf-example"
}

data "alibabacloudstack_zones" "default" {
  enable_details = true
}

data "alibabacloudstack_cspprivate_hsm_vendors" "default" {
}

resource "alibabacloudstack_vpc" "vpc" {
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alibabacloudstack_vswitch" "vsw" {
  vpc_id       = alibabacloudstack_vpc.vpc.id
  cidr_block   = "192.168.0.0/24"
  zone_id      = data.alibabacloudstack_zones.default.zones.0.id
}

resource "alibabacloudstack_cspprivate_hsm_instance" "default" {
  product_code   = data.alibabacloudstack_cspprivate_hsm_vendors.default.vendors.0.products.0.code
  vendor_code    = data.alibabacloudstack_cspprivate_hsm_vendors.default.vendors.0.code
  vsm_type       = "gvsm"
  zone_id        = data.alibabacloudstack_zones.default.zones.0.id
  alias_name     = var.name
  vpc_id         = alibabacloudstack_vpc.vpc.id
  vpc_cidr_block = alibabacloudstack_vpc.vpc.cidr_block
  vswitch_id     = alibabacloudstack_vswitch.vsw.id
  ip             = "192.168.0.100"
}
```

### 使用指定设备

```hcl
variable "name" {
  default = "tf-example"
}

data "alibabacloudstack_zones" "default" {
  enable_details = true
}

data "alibabacloudstack_cspprivate_hsm_vendors" "default" {
}

data "alibabacloudstack_cspprivate_hsms" "default" {
  vendor_code  = data.alibabacloudstack_cspprivate_hsm_vendors.default.vendors.0.code
  product_code = data.alibabacloudstack_cspprivate_hsm_vendors.default.vendors.0.products.0.code
  zone_id      = data.alibabacloudstack_zones.default.zones.0.id
}

resource "alibabacloudstack_cspprivate_hsm_instance" "default" {
  product_code = data.alibabacloudstack_cspprivate_hsm_vendors.default.vendors.0.products.0.code
  vendor_code  = data.alibabacloudstack_cspprivate_hsm_vendors.default.vendors.0.code
  vsm_type     = "gvsm"
  zone_id      = data.alibabacloudstack_zones.default.zones.0.id
  alias_name   = var.name
  device_id    = data.alibabacloudstack_cspprivate_hsms.default.hsms.0.id
}
```

## 参数说明

以下参数可用于配置HSM实例：

* `product_code` - (必填, 强制新建) HSM设备的产品型号代码。
* `vendor_code` - (必填, 强制新建) HSM设备的供应商代码。
* `vsm_type` - (必填, 强制新建) HSM实例的类型。可选值: `evsm`, `gvsm`, `svsm`。
* `zone_id` - (必填, 强制新建) HSM实例所在的可用区ID。
* `device_id` - (可选, 自动计算) HSM设备的ID。
* `vpc_id` - (可选, 自动计算) HSM实例所在的VPC ID。
* `vpc_cidr_block` - (可选, 自动计算) VPC的CIDR块。
* `vswitch_id` - (可选, 自动计算) HSM实例所在的VSwitch ID。
* `ip` - (可选, 自动计算) HSM实例的IP地址。
* `alias_name` - (可选) HSM实例的别名。

## 属性说明

以下属性将会被导出：

* `id` - HSM实例的ID。
* `port` - HSM实例的端口。
* `status` - HSM实例的状态。

## 导入方式

HSM实例可以通过ID进行导入，例如：

```shell
$ terraform import alibabacloudstack_cspprivate_hsm_instance.example <id>
```