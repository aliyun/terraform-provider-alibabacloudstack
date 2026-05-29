---
subcategory: "PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_cluster_account_database_binding"
sidebar_current: "docs-Alibabacloudstack-resource-polardb-cluster-account-database-binding"
description: |-
  Provides a PolarDB cluster account database binding resource.
---

# alibabacloudstack_polardb_cluster_account_database_binding

Provides a PolarDB cluster account database binding resource. This resource allows you to manage database privileges for accounts in PolarDB clusters.

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
resource "alibabacloudstack_polardb_cluster_database" "default0" {
  db_cluster_id      = "${alibabacloudstack_polardb_cluster_instance.default.id}"
  db_name            = "${var.name}"
  character_set_name = "utf8"
  db_description     = "Test database for PolarDB cluster"
}
resource "alibabacloudstack_polardb_cluster_database" "default1" {
  db_cluster_id      = "${alibabacloudstack_polardb_cluster_instance.default.id}"
  db_name            = "${var.name}"
  character_set_name = "utf8"
  db_description     = "Test database for PolarDB cluster"
}

resource "alibabacloudstack_polardb_cluster_account_database_binding" "example" {
  db_cluster_id = "${alibabacloudstack_polardb_cluster_instance.default.id}"
  account_name  = "${alibabacloudstack_polardb_cluster_account.example.account_name}"

  database_privileges {
    db_name    = "${alibabacloudstack_polardb_cluster_database.default0.db_name}"
    privilege  = "ReadWrite"
  }

  database_privileges {
    db_name    = "${alibabacloudstack_polardb_cluster_database.default1.db_name}"
    privilege  = "ReadOnly"
  }
}
```

## Argument Reference

The following arguments are supported:

* `account_name` - (Required) The name of the account.
* `db_cluster_id` - (Required) The ID of the PolarDB cluster.
* `database_privileges` - (Optional) The database privileges for the account. Each entry supports the following:
  * `db_name` - (Optional) The name of the database.
  * `privilege` - (Optional) The privilege level for the database. Valid values: `ReadWrite`, `ReadOnly`, `DDLOnly`, `DMLOnly`, `ReadIndex`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the binding, formatted as `{DBClusterId}:{AccountName}`.

## Import

PolarDB cluster account database binding can be imported using the id, e.g.

```bash
$ terraform import alibabacloudstack_polardb_cluster_account_database_binding.example pc-12345678:test_account
```