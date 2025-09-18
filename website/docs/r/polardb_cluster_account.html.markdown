---
subcategory: "PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_cluster_account"
sidebar_current: "docs-alibabacloudstack-resource-polardb-cluster-account"
description: |-
  Provides a PolarDB cluster account resource.
---

# alibabacloudstack_polardb_cluster_account

Provides a PolarDB cluster account resource. This resource allows you to manage accounts in PolarDB clusters.

## Example Usage

```hcl
variable "password" {
  type = string
}

variable "name" {
  default = "tf-polar-test"
}
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
}
resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name       = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}
resource "alibabacloudstack_vpc_vswitch" "default" {
  vpc_id            = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  zone_id           = "${data.alibabacloudstack_zones.default.zones.0.id}"
  vswitch_name      = "${var.name}_vsw"
}
resource "alibabacloudstack_polardb_cluster_instance" "default" {
  db_cluster_description 	=  "${var.name}"
  zone_id 				= "${data.alibabacloudstack_zones.default.zones.0.id}"
  db_type 				= "${var.db_type}"
  db_version 			= "${var.db_version}"
  storage_space 		= "20"
  vpc_id 				= "${alibabacloudstack_vpc_vpc.default.id}"
  vswitch_id			= "${alibabacloudstack_vpc_vswitch.default.id}"
  db_node_class 		= "${data.alibabacloudstack_polardb_cluster_instance_types.default.instance_types.0.id}"
  sub_category 			= "${data.alibabacloudstack_polardb_cluster_instance_types.default.instance_types.0.sub_category}"
  storage_type 			= "ESSDPL1"
}

resource "alibabacloudstack_polardb_cluster_account" "example" {
  db_cluster_id        = "${alibabacloudstack_polardb_cluster_instance_types.default.id}"
  account_name         = "test_account"
  account_password     = var.password
  account_description  = "Test account for PolarDB cluster"
  account_type         = "Normal"
  account_lock_state   = "UnLock"
}
```

## Argument Reference

The following arguments are supported:

* `db_cluster_id` - (Required, ForceNew) The ID of the PolarDB cluster.
* `account_name` - (Required, ForceNew) The name of the account.
* `account_password` - (Optional, Sensitive) The password of the account.
* `account_type` - (Optional) The type of the account.
* `account_description` - (Optional) The description of the account. 
* `account_lock_state` - (Optional) The lock state of the account. Valid values: `UnLock`, `Lock`. Default value: `UnLock`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the account, formatted as `{DBClusterId}:{AccountName}`.
* `account_status` - The status of the account.
* `account_password_valid_time` - The validity period of the account password.

## Import

PolarDB cluster account can be imported using the id, e.g.

```bash
$ terraform import alibabacloudstack_polardb_cluster_account.example pc-12345678:test_account
```