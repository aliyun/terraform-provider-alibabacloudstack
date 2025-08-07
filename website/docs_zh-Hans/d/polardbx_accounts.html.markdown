---
subcategory: "PolarDBX"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardbx_accounts"
sidebar_current: "docs-Alibabacloudstack-datasource-polardbx-accounts"
description: |-
  提供阿里云账号下拥有的polardbx accounts列表。
---

# alibabacloudstack\_polardbx\_accounts

此数据源提供根据指定过滤条件列出的阿里云账号下的polardbx accounts资源列表。

## 示例用法
```
variable "name" {
	default = "tf-testAccPolardbAccounts19559"
}

data  "alibabacloudstack_zones" "default" {
	available_resource_creation = "PolarDBX"
}
resource "alibabacloudstack_polardbx_dbinstance" "instance" {
	engine            = "MySQL"
	engine_version    = "5.7"
	instance_name = "tfinstance"
	db_instance_storage_type= "local_ssd"
	db_instance_storage = 5
	db_instance_class = "rds.mysql.t1.small"
	zone_id= "${data.alibabacloudstack_zones.default.zones.0.id}"
}
resource "alibabacloudstack_polardbx_account" "default" {
	data_base_instance_id = "${alibabacloudstack_polardbx_dbinstance.instance.id}"
	account_description = "test"
	account_name        = "polardbx_account"
	account_password = "NyCc0x6b!rH^"
	account_type ="Normal"
}

data "alibabacloudstack_polardbx_accounts" "default" {
  ids = [
          "${alibabacloudstack_polardbx_account.default.id}"
        ]
  db_instance_id = "${alibabacloudstack_polardbx_dbinstance.instance.id}"
}
```

## 参数参考

支持以下参数：
  * `ids` - (可选) - 用于过滤结果的账户 ID 列表。
  * `names` - (可选) - 用于过滤结果的账户名称列表。
  * `account_name` - (可选) - 账户名称。
  * `instance_id` - (必填) - 数据库实例 ID。
  * `account_type` - (可选) - 账户类型。取值范围如下：-**Normal**：普通账户。-**Super**：高权限账户。

## 属性参考

除上述参数外，还导出以下属性：
  * `accounts` - 账户列表。每个元素包含以下属性：
    * `id` - 账户 ID。
    * `description` - 账户描述。
    * `account_name` - 账户名称。
    * `account_type` - 账户类型。取值范围如下：-**Normal**：普通账户。-**Super**：高权限账户。
    * `instance_id` - PolarDBX 数据库实例的 ID。
    * `db_privileges` - 当前账户的数据库权限。
      * `db_name` - 数据库名称。
      * `privilege` - 数据库权限。