---
subcategory: "PolarDB-X"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardbx_super_account"
sidebar_current: "docs-Alibabacloudstack-polardbx-super-account"
description: |-
  编排PolarDB-X实例的高权限账号及三权分立配置
---

# alibabacloudstack_polardbx_super_account

使用Provider配置的凭证在指定PolarDB-X实例中创建高权限账号并管理三权分立功能。

## 示例用法

### 基础用法

```hcl
variable "name" {
}


resource "random_password" "password" {
  count            = 2
  length           = 12
  special          = true
  override_special = "!@#$^&*()_"
  min_lower        = 1
  min_upper        = 1
  min_numeric      = 1
}


data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details              = true
}


resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name   = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  name       = "${var.name}_vsw"
  vpc_id     = alibabacloudstack_vpc_vpc.default.id
  cidr_block = "172.16.1.0/24"
  zone_id    = data.alibabacloudstack_zones.default.zones.0.id
}




variable "existed_polardbx_id" {
  type    = string
}

data "alibabacloudstack_polardbx_instance_types" "cn" {
  sorted_by = "CPU"
  spec_type = "CN"
}

data "alibabacloudstack_polardbx_instance_types" "dn" {
  sorted_by = "CPU"
  spec_type = "DN"
}

data "alibabacloudstack_polardbx_instances" "default" {
  ids = var.existed_polardbx_id == "" ? [" ", ] : ["${var.existed_polardbx_id}", ]
}

resource "alibabacloudstack_polardbx_instance" "default" {
  count          = length(data.alibabacloudstack_polardbx_instances.default.polardbx_instances) == 0 ? 1 : 0
  zone_id        = data.alibabacloudstack_zones.default.zones.0.id
  engine_version = "5.7"
  storage        = 50
  vswitch_id     = alibabacloudstack_vpc_vswitch.default.id
  cn_node_class  = data.alibabacloudstack_polardbx_instance_types.cn.instance_types.0.id
  cn_node_count  = "2"
  dn_node_class  = data.alibabacloudstack_polardbx_instance_types.dn.instance_types.0.id
  dn_node_count  = "2"
}
locals {
  polardbx_instance = length(data.alibabacloudstack_polardbx_instances.default.polardbx_instances) == 0 ? alibabacloudstack_polardbx_instance.default.0 : data.alibabacloudstack_polardbx_instances.default.polardbx_instances.0
}

resource "alibabacloudstack_polardbx_super_account" "default" {
  security_account_password    = random_password.password.0.result
  audit_account_password       = random_password.password.0.result
  instance_id                  = local.polardbx_instance.id
  audit_account_name           = "audit_user"
  audit_account_description    = "audit user"
  admin_account_description    = "system user"
  security_account_name        = "security_user"
  security_account_description = "security user"
  admin_account_name           = "admin_user"
  admin_account_password       = random_password.password.0.result
}
```

## 参数说明

支持以下参数：

* `instance_id` - (必填, 变更时重建) PolarDB-X实例ID。格式为pxc-************
* `admin_account_name` - (必填, 变更时重建) 高权限账号名称。长度1-16个字符，支持字母、数字、下划线和连字符
* `admin_account_password` - (必填) 高权限账号密码。长度8-32个字符，需包含大小写字母、数字和特殊字符
* `admin_account_description` - (可选) 高权限账号描述信息。长度2-256个字符
* `security_account_name` - (可选) 安全管理员账号名称。与audit_account_name必须同时存在
* `security_account_password` - (可选) 安全管理员账号密码。与security_account_name必须同时存在
* `security_account_description` - (可选) 安全管理员账号描述信息
* `audit_account_name` - (可选) 审计员账号名称。与security_account_name必须同时存在
* `audit_account_password` - (可选) 审计员账号密码。与audit_account_name必须同时存在
* `audit_account_description` - (可选) 审计员账号描述信息

## 属性说明

以下属性会导出：

* `id` - 资源唯一标识符，与instance_id保持一致
* `three_roles` - (布尔值) 是否启用三权分立功能。true表示已启用，false表示未启用
