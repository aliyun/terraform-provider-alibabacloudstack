---
subcategory: "PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_account"
sidebar_current: "docs-Alibabacloudstack-polardb-account"
description: |- 
  编排polardb用户
---

# alibabacloudstack_polardb_account

使用Provider配置的凭证在指定的资源集编排polardb用户。

## 示例用法
```hcl
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "${var.creation}"
}

variable "password" {
  default = "SecurePassword123!"
}

variable "name" {
  default = "tf-testAccdbaccount-533060"
}

variable "creation" {
  default = "PolarDB"
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

resource "alibabacloudstack_polardb_dbinstance" "instance" {
  engine                  = "MySQL"
  engine_version          = "5.7"
  instance_name           = "${var.name}"
  db_instance_storage_type = "local_ssd"
  db_instance_storage     = 5
  db_instance_class       = "rds.mysql.t1.small"
  zone_id                = "${data.alibabacloudstack_zones.default.zones.0.id}"
  vswitch_id             = "${alibabacloudstack_vswitch.default.id}"
}

resource "alibabacloudstack_polardb_account" "default" {
  data_base_instance_id = "${alibabacloudstack_polardb_dbinstance.instance.id}"
  account_name          = "tftestnormal"
  account_password      = var.password
  account_description   = "This is a test account for PolarDB."
  account_type          = "Normal"

  database_privileges {
    data_base_name       = "test_db"
    account_privilege    = "ReadOnly"
    account_privilege_detail = "SELECT"
  }
}
```

## 参数说明

支持以下参数：
  * `account_description` - (可选) - 账号备注说明，需满足如下要求：不能以`http://`或`https://`开头；长度为2~256个字符。
  * `account_name` - (必填) - 账号名称，需符合如下要求：以小写字母开头，以字母或数字结尾；由小写字母、数字或下划线组成；长度为2~16个字符；不能使用某些预留的用户名，如root、admin等。
  * `account_password` - (必填) - 账号密码。
  * `account_type` - (可选) - 账号类型，取值范围如下：**Normal**：普通账号；**Super**：高权限账号。> 如果该参数留空，则默认创建**Super**账号。当集群为PolarDB O引擎或PolarDB PostgreSQL引擎时，每个集群允许创建多个高权限账号，高权限账号相比普通账号拥有更多权限，详情参见[创建数据库账号](~~68508~~)。当集群为PolarDB MySQL引擎时，每个集群最多只允许创建1个高权限账号，高权限账号相比普通账号拥有更多权限，详情参见[创建数据库账号](~~68508~~)。
  * `data_base_instance_id` - (必填) - 将要关联的 PolarDB 实例 ID。
  * `database_privileges` - (可选) - 目标账号拥有的数据库权限详情。
    
    * `data_base_name` - (必填) - 权限应用的数据库名称。
    * `account_privilege` - (必填) - 在指定数据库上对账号的权限级别，例如：`ReadOnly`, `ReadWrite` 等。
    * `account_privilege_detail` - (可选) - 账号权限的详细信息，例如具体的SQL操作权限(如`SELECT`, `INSERT` 等)。
  * `priv_exceeded` - (可选) - 表示权限是否超出允许的限制。如果设置为`true`，表示权限超出限制。
  * `status` - (可选) - 代表资源状态的资源属性字段。

## 属性说明

除了上述所有参数外，还导出了以下属性：
  * `account_description` - 账号备注说明，需满足如下要求：不能以`http://`或`https://`开头；长度为2~256个字符。
  * `account_type` - 账号类型，取值范围如下：**Normal**：普通账号；**Super**：高权限账号。> 如果该参数留空，则默认创建**Super**账号。当集群为PolarDB O引擎或PolarDB PostgreSQL引擎时，每个集群允许创建多个高权限账号，高权限账号相比普通账号拥有更多权限，详情参见[创建数据库账号](~~68508~~)。当集群为PolarDB MySQL引擎时，每个集群最多只允许创建1个高权限账号，高权限账号相比普通账号拥有更多权限，详情参见[创建数据库账号](~~68508~~)。
  * `priv_exceeded` - 表示权限是否超出允许的限制。如果设置为`true`，表示权限超出限制。
  * `status` - 代表资源状态的资源属性字段。