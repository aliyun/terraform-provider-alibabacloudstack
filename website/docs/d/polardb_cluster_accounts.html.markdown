---
subcategory: "PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_cluster_accounts"
sidebar_current: "docs-Alibabacloudstack-datasource-polardb-cluster-accounts"
description: |-
  Provides a list of polardb cluster accounts owned by an alibabacloudstack account.
---

# alibabacloudstack\_polardb\_cluster\_accounts

This data source provides a list of polardb cluster accounts in an alibabacloudstack account according to the specified filters.

## Example Usage
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

## Argument Reference

The following arguments are supported:
  * `ids` - (Optional) - A list of account IDs to filter results.
  * `account_name` - (Optional) - The account name must meet the following requirements:* Start with a lowercase letter and end with a letter or number.* Consists of lowercase letters, numbers, or underscores.* The length is 2 to 16 characters.* You cannot use some reserved usernames, such as root and admin.
  * `db_cluster_id` - (Required) - db cluster id.
  * `name_regex` - (Optional, Deprecated) - A regex string to filter results by account name. Field 'name_regex' is deprecated and will be removed in a future release. Please use new field 'description_regex' instead.
  * `description_regex` - (Optional) - A regex string to filter results by account description.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `accounts` - A list of accounts. Each element contains the following attributes:
    * `id` - The ID of the account.
    * `account_description` - The account number Notes shall meet the following requirements:-Cannot start with' http:// 'or' https.-2 to 256 characters in length.
    * `account_name` - The account name must meet the following requirements:* Start with a lowercase letter and end with a letter or number.* Consists of lowercase letters, numbers, or underscores.* The length is 2 to 16 characters.* You cannot use some reserved usernames, such as root and admin.
    * `account_type` - Account type. The value range is as follows:-**Normal**: Normal account.-**Super**: a highly privileged account.> * If this parameter is left blank, the **Super** account is created by default.* When the cluster is PolarDB O engine or PolarDB PostgreSQL engine, each cluster can create multiple high-permission accounts. High-permission accounts have more permissions than normal accounts. For more information about creating database accounts, see [create database accounts](~~ 68508 ~~).* When the cluster is the PolarDB MySQL engine, each cluster can only create one high-permission account at most. High-permission accounts have more permissions than normal accounts. For more information about creating database accounts, see [create database accounts](~~ 68508 ~~).
    * `account_lock_state` - The lock state of the account.
    * `status` - The status of the resource.
    * `database_privileges` - The Database permissions of the target account.
      * `account_privilege` - The privilege of the account.
      * `account_privilege_detail` - The privilege detail of the account.
      * `data_base_name` - The name of the database.
