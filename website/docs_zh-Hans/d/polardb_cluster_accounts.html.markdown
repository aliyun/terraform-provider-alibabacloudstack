---
subcategory: "PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_cluster_accounts"
sidebar_current: "docs-Alibabacloudstack-datasource-polardb-cluster-accounts"
description: |-
  提供阿里云账户拥有的 PolarDB 集群账户列表。
---

# alibabacloudstack\_polardb\_cluster\_accounts

本数据源根据指定的过滤条件，提供阿里云账户中的 PolarDB 集群账户列表。

## 示例用法
```
variable "name" {
	default = "tfAccountsName"
}

variable "password" {
}

data  "alibabacloudstack_zones" "default" {
	available_resource_creation = "PolarDB"
}

resource "alibabacloudstack_polardb_cluster" "cluster" {
	engine            = "MySQL"
	engine_version    = "5.7"
	cluster_name      = "tfcluster"
	db_node_class     = "polar.mysql.x4.large"
	db_node_count     = 2
	db_node_storage   = 20
	zone_id           = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

resource "alibabacloudstack_polardb_cluster_account" "default" {
	db_cluster_id        = "${alibabacloudstack_polardb_cluster.cluster.id}"
	account_description  = "test"
	account_name         = "polardb_cluster_account"
	account_password     = "${var.password}"
	account_type         = "Normal"
}

data "alibabacloudstack_polardb_cluster_accounts" "default" {
  ids = [
          "${alibabacloudstack_polardb_cluster_account.default.id}"
        ]
  db_cluster_id = "${alibabacloudstack_polardb_cluster.cluster.id}"
}
```

## 参数参考

支持以下参数：
  * `ids` - (可选) - 用于过滤结果的账户 ID 列表。
  * `account_name` - (可选) - 账户名称
  * `db_cluster_id` - (必填) - 数据库集群 ID。
  * `name_regex` - (可选, 已弃用) - 用于通过账户名称过滤结果的正则表达式。字段 'name_regex' 已弃用，并将在未来版本中移除。请使用新字段 'description_regex' 替代。
  * `description_regex` - (可选) - 用于通过账户描述过滤结果的正则表达式。

## 属性参考

除了上述参数外，还导出以下属性：
  * `accounts` - 账户列表。每个元素包含以下属性：
    * `id` - 账户的 ID。
    * `account_description` - 账户备注
    * `account_name` - 账户名称
    * `account_type` - 账户类型
    * `account_lock_state` - 账户的锁定状态。
    * `status` - 资源的状态。
    * `database_privileges` - 目标账户的数据库权限。
      * `account_privilege` - 账户的权限。
      * `account_privilege_detail` - 账户的权限详情。
      * `data_base_name` - 数据库的名称。