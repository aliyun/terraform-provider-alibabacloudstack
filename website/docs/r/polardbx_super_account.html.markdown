---
subcategory: "Cloud-Native Distributed Database PolarDB-X 2.0"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardbx_super_account"
description: |-
  Orchestrates super accounts and ternary authorization configurations for POLARDB-X instances

---

# alibabacloudstack_polardbx_super_account

Creates and manages super accounts with ternary authorization configuration in specified POLARDB-X instances using provider credentials.

## Example Usage

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

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) POLARDB-X instance ID. Format: `pxc-************`
* `admin_account_name` - (Required, ForceNew) Super account name. 1-16 characters (letters, digits, underscores, hyphens).
* `admin_account_password` - (Required) Super account password. 8-32 characters containing uppercase, lowercase, digits, and special characters.
* `admin_account_description` - (Optional) Super account description. 2-256 characters.
* `security_account_name` - (Optional) Security admin account name. Must coexist with `audit_account_name`.
* `security_account_password` - (Optional) Security admin password. Required when `security_account_name` exists.
* `security_account_description` - (Optional) Security admin account description.
* `audit_account_name` - (Optional) Auditor account name. Must coexist with `security_account_name`.
* `audit_account_password` - (Optional) Auditor password. Required when `audit_account_name` exists.
* `audit_account_description` - (Optional) Auditor account description.

## Attributes Reference

The following attributes are exported:

* `id` - Resource identifier (matches `instance_id`)
* `three_roles` - (Boolean) Ternary authorization status. `true` = enabled, `false` = disabled