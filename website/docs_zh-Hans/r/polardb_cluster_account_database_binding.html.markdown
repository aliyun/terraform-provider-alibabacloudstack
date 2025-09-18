---
subcategory: "PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_cluster_account_database_binding"
sidebar_current: "docs-alibabacloudstack-resource-polardb-cluster-account-database-binding"
description: |-
  提供 PolarDB 集群账户数据库绑定资源。
---

# alibabacloudstack_polardb_cluster_account_database_binding

提供 PolarDB 集群账户数据库绑定资源。该资源允许您管理 PolarDB 集群中账户的数据库权限。

## 示例用法

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

## 参数参考

支持以下参数：

* `account_name` - (必填) 账户名称。
* `db_cluster_id` - (必填) PolarDB 集群的 ID。
* `database_privileges` - (可选) 账户的数据库权限。每个条目支持以下参数：
  * `db_name` - (可选) 数据库名称。
  * `privilege` - (可选) 数据库的权限级别。有效值：`ReadWrite`、`ReadOnly`、`DDLOnly`、`DMLOnly`、`ReadIndex`。

## 属性参考

导出以下属性：

* `id` - 绑定的 ID，格式为 `{DBClusterId}:{AccountName}`。

## 导入

可以使用 id 导入 PolarDB 集群账户数据库绑定，例如：

```bash
$ terraform import alibabacloudstack_polardb_cluster_account_database_binding.example pc-12345678:test_account
```