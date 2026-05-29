---
subcategory: "云原生数据库 PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_cluster_account"
sidebar_current: "docs-Alibabacloudstack-resource-polardb-cluster-account"
description: |-
  提供 PolarDB 集群账户资源。
---

# alibabacloudstack_polardb_cluster_account

提供 PolarDB 集群账户资源。该资源允许您管理 PolarDB 集群中的账户。

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
```

## 参数参考

支持以下参数：

* `db_cluster_id` - (必填, ForceNew) PolarDB 集群的 ID。
* `account_name` - (必填, ForceNew) 账户名称。
* `account_password` - (可选, 敏感) 账户密码。
* `account_type` - (可选) 账户类型。
* `account_description` - (可选) 账户描述。
* `account_lock_state` - (可选) 账户锁定状态。有效值：`UnLock`、`Lock`。默认值：`UnLock`。

## 属性参考

导出以下属性：

* `id` - 账户的 ID，格式为 `{DBClusterId}:{AccountName}`。
* `account_status` - 账户状态。
* `account_password_valid_time` - 账户密码的有效期。

## 导入

可以使用 id 导入 PolarDB 集群账户，例如：

```bash
$ terraform import alibabacloudstack_polardb_cluster_account.example pc-12345678:test_account
```