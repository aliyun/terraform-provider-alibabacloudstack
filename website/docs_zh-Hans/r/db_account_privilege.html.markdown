---
subcategory: "RDS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_db_account_privilege"
sidebar_current: "docs-alibabacloudstack-resource-db-account-privilege"
description: |-
  RDS实例授权
---

# alibabacloudstack_db_account_privilege
使用Provider配置的凭证在指定的资源集下为RDS实例授予多个数据库某些访问权限。一个数据库可以被多个账户授予权限。

## 示例用法

```
variable "creation" {
  default = "Rds"
}

variable "name" {
  default = "dbaccountprivilegebasic"
}

variable "password" {
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "${var.creation}"
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

resource "alibabacloudstack_db_instance" "instance" {
  engine           = "MySQL"
  engine_version   = "5.6"
  instance_type    = "rds.mysql.s1.small"
  instance_storage = "10"
  vswitch_id       = "${alibabacloudstack_vswitch.default.id}"
  instance_name    = "${var.name}"
}

resource "alibabacloudstack_db_database" "db" {
  count       = 2
  instance_id = "${alibabacloudstack_db_instance.instance.id}"
  name        = "tfaccountpri_${count.index}"
  description = "from terraform"
}

resource "alibabacloudstack_db_account" "account" {
  instance_id = "${alibabacloudstack_db_instance.instance.id}"
  name        = "tftestprivilege"
  password    = var.password
  description = "from terraform"
}

resource "alibabacloudstack_db_account_privilege" "privilege" {
  instance_id  = "${alibabacloudstack_db_instance.instance.id}"
  account_name = "${alibabacloudstack_db_account.account.name}"
  privilege    = "ReadOnly"
  db_names     = "${alibabacloudstack_db_database.db.*.name}"
}
```

## 参数说明

支持以下参数：

* `instance_id` - (必填，变更时重建) 账户所属实例的ID。
* `account_name` - (必填，变更时重建) 指定的账户名称。
* `privilege` - 一个账户访问数据库的权限。有效值为：
    - ReadOnly: 此值仅适用于 MySQL、MariaDB 和 SQL Server。
    - ReadWrite: 此值仅适用于 MySQL、MariaDB 和 SQL Server。
    
   默认值为 "ReadOnly"。
* `db_names` - (必填) 指定的数据库名称列表。
* `data_base_instance_id` - (可选，变更时重建) 数据库实例的ID。此参数与 `instance_id` 类似，但在某些场景下可能需要单独指定。

## 属性说明

导出以下属性：

* `id` - 当前账户资源ID。由实例ID、账户名称和权限组成，格式为 `<instance_id>:<name>:<privilege>`。
* `instance_id` - 账户所属实例的ID。
* `data_base_instance_id` - 数据库实例的ID。此属性与 `instance_id` 类似，但可能在某些情况下提供额外的信息。