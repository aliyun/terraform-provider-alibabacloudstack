---
subcategory: "PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_account"
sidebar_current: "docs-Alibabacloudstack-polardb-account"
description: |-
  Provides a polardb Account resource.
---

# alibabacloudstack_polardb_account

Provides a polardb Account resource.

## Example Usage
```
data  "alibabacloudstack_zones" "default" {
  available_resource_creation = "${var.creation}"
}
variable "password" {
}
resource "alibabacloudstack_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = "${alibabacloudstack_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  name              = "${var.name}"
}

	variable "name" {
		default = "tf-testAccdbaccount-533060"
	}
	variable "creation" {
		default = "PolarDB"
	}
	resource "alibabacloudstack_polardb_dbinstance" "instance" {
		engine            = "MySQL"
		engine_version    = "5.7"
		instance_name = "${var.name}"
		db_instance_storage_type= "local_ssd"
		db_instance_storage = 5
		db_instance_class = "rds.mysql.t1.small"
		zone_id= "${data.alibabacloudstack_zones.default.zones.0.id}"
		vswitch_id = "${alibabacloudstack_vswitch.default.id}"
	}
	

resource "alibabacloudstack_polardb_account" "default" {
  data_base_instance_id = "${alibabacloudstack_polardb_dbinstance.instance.id}"
  account_name = "tftestnormal"
  account_password = var.password
}
```

## Argument Reference

The following arguments are supported:
  * `account_description` - (Optional) - The account number Notes shall meet the following requirements:-Cannot start with' http:// 'or' https.-2 to 256 characters in length.
  * `account_name` - (Required) - The account name must meet the following requirements:* Start with a lowercase letter and end with a letter or number.* Consists of lowercase letters, numbers, or underscores.* The length is 2 to 16 characters.* You cannot use some reserved usernames, such as root and admin.
  * `account_password` - (Required) - Account password
  * `account_type` - (Optional) - Account type. The value range is as follows:-**Normal**: Normal account.-**Super**: a highly privileged account.> * If this parameter is left blank, the **Super** account is created by default.* When the cluster is PolarDB O engine or PolarDB PostgreSQL engine, each cluster can create multiple high-permission accounts. High-permission accounts have more permissions than normal accounts. For more information about creating database accounts, see [create database accounts](~~ 68508 ~~).* When the cluster is the PolarDB MySQL engine, each cluster can only create one high-permission account at most. High-permission accounts have more permissions than normal accounts. For more information about creating database accounts, see [create database accounts](~~ 68508 ~~).
  * `data_base_instance_id` - (Required) -  The ID of the PolarDB instance to which the account will be associated.
  * `database_privileges` - (Optional) - The Database permissions of the target account.
    
    * `account_privilege` - (Optional) - The privilege level for the account on the specified database.
    
    * `account_privilege_detail` - (Optional) - Detailed information about the account privileges.
    
    * `data_base_name` - (Optional) - The name of the database to which the privileges are applied.
  * `priv_exceeded` - (Optional) - Indicates whether the privileges exceed the allowed limits.
  * `status` - (Optional) - The status of the resource

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `account_description` - The account number Notes shall meet the following requirements:-Cannot start with' http:// 'or' https.-2 to 256 characters in length.
  * `account_type` - Account type. The value range is as follows:-**Normal**: Normal account.-**Super**: a highly privileged account.> * If this parameter is left blank, the **Super** account is created by default.* When the cluster is PolarDB O engine or PolarDB PostgreSQL engine, each cluster can create multiple high-permission accounts. High-permission accounts have more permissions than normal accounts. For more information about creating database accounts, see [create database accounts](~~ 68508 ~~).* When the cluster is the PolarDB MySQL engine, each cluster can only create one high-permission account at most. High-permission accounts have more permissions than normal accounts. For more information about creating database accounts, see [create database accounts](~~ 68508 ~~).
  * `priv_exceeded` -  Indicates whether the privileges exceed the allowed limits.
  * `status` - The status of the resource
