---
subcategory: "云原生数据库 PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_accounts"
sidebar_current: "docs-Alibabacloudstack-datasource-polardb-accounts"
description: |-
  提供阿里云账号下拥有的polardb accounts列表。
---

# alibabacloudstack\_polardb\_accounts

此数据源提供根据指定过滤条件列出的阿里云账号下的polardb accounts资源列表。

## 示例用法
```
variable "name" {
	default = "tf-testAccPolardbAccounts19559"
}

variable "password" {
}

data  "alibabacloudstack_zones" "default" {
	available_resource_creation = "PolarDB"
}
resource "alibabacloudstack_polardb_dbinstance" "instance" {
	engine            = "MySQL"
	engine_version    = "5.7"
	instance_name = "tfinstance"
	db_instance_storage_type= "local_ssd"
	db_instance_storage = 5
	db_instance_class = "rds.mysql.t1.small"
	zone_id= "${data.alibabacloudstack_zones.default.zones.0.id}"
}
resource "alibabacloudstack_polardb_account" "default" {
	data_base_instance_id = "${alibabacloudstack_polardb_dbinstance.instance.id}"
	account_description = "test"
	account_name        = "polardb_account"
	account_password = "${var.password}"
	account_type ="Normal"
}

data "alibabacloudstack_polardb_accounts" "default" {
  ids = [
          "${alibabacloudstack_polardb_account.default.id}"
        ]
  db_instance_id = "${alibabacloudstack_polardb_dbinstance.instance.id}"
}
```

## 参数参考
以下参数是支持的：
  * `ids` - (选填) - 用于过滤结果的Account ID列表。
  * `account_name` - (选填) - 账号名称，需符合如下要求：* 以小写字母开头，以字母或数字结尾。* 由小写字母、数字或下划线组成。* 长度为2~16个字符。* 不能使用某些预留的用户名，如root、admin等。
  * `db_instance_id` - (必填) - 数据库实例ID
  * `name_regex` - (选填) - 用于过滤结果的账号名称，支持正则表达式。
  * `description_regex` - (选填) - 用于过滤结果的账号备注说明，支持正则表达式。

## Attributes Reference
除了上述参数外，还导出以下属性：
  * `accounts` - 账号列表。
    * `id` - 账号ID
    * `account_description` - 账号备注说明，需满足如下要求：- 不能以`http://`或`https://`开头。- 长度为2~256个字符。
    * `account_name` - 账号名称，需符合如下要求：* 以小写字母开头，以字母或数字结尾。* 由小写字母、数字或下划线组成。* 长度为2~16个字符。* 不能使用某些预留的用户名，如root、admin等。
    * `account_type` - 账号类型，取值范围如下：- **Normal**：普通账号。 - **Super**：高权限账号。 > * 若该参数留空，则默认创建**Super**账号。* 当集群为PolarDB O引擎或PolarDB PostgreSQL引擎时，每个集群允许创建多个高权限账号，高权限账号相比普通账号拥有更多权限，创建数据库账号详情参见[创建数据库账号](~~68508~~)。* 当集群为PolarDB MySQL引擎时，每个集群最多只允许创建1个高权限账号，高权限账号相比普通账号拥有更多权限，创建数据库账号详情参见[创建数据库账号](~~68508~~)。
    * `db_instance_id` - 数据库实例ID
    * `database_privileges` - 目标账号拥有的数据库权限详情。
    * `priv_exceeded` - 权限超限
    * `status` - 代表资源状态的资源属性字段