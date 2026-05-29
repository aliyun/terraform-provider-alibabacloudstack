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
  * `account_name` - (Optional) - The account name
  * `db_cluster_id` - (Required) - db cluster id.
  * `name_regex` - (Optional, Deprecated) - A regex string to filter results by account name. Field 'name_regex' is deprecated and will be removed in a future release. Please use new field 'description_regex' instead.
  * `description_regex` - (Optional) - A regex string to filter results by account description.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `accounts` - A list of accounts. Each element contains the following attributes:
    * `id` - The ID of the account.
    * `account_description` - The account number Notes shall meet the following requirements:-Cannot start with' http:// 'or' https.-2 to 256 characters in length.
    * `account_name` - The account name
    * `account_type` - Account type.
    * `account_lock_state` - The lock state of the account.
    * `status` - The status of the resource.
    * `database_privileges` - The Database permissions of the target account.
      * `privilege` - The privilege of the account.
      * `db_name` - The name of the database.
