---
subcategory: "RDS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_rds_account"
sidebar_current: "docs-Alibabacloudstack-rds-account"
description: |- 
  编排RDS数据库帐号
---

# alibabacloudstack_rds_account
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_db_account`

使用Provider配置的凭证在指定的资源集编排RDS数据库帐号。

## 示例用法

```hcl
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

variable "creation" {
  default = "Rds"
}

variable "name" {
  default = "tf-testAccdbaccount-52957"
}

variable "password" {
  default = "YourSecurePassword123"
}

resource "alibabacloudstack_db_instance" "instance" {
  engine               = "MySQL"
  engine_version       = "5.6"
  instance_type        = "rds.mysql.s2.large"
  instance_storage     = "30"
  vswitch_id           = "${alibabacloudstack_vswitch.default.id}"
  instance_name        = "${var.name}"
  storage_type         = "local_ssd"
}

resource "alibabacloudstack_rds_account" "default" {
  data_base_instance_id = "${alibabacloudstack_db_instance.instance.id}"
  account_name          = "tftestnormal"
  password              = var.password
  account_type          = "Normal"
  account_description   = "This is a test account."
}
```

## 参数说明

支持以下参数：

* `data_base_instance_id` - (必填，变更时重建) 实例 ID。您可以通过调用 `DescribeDBInstances` 操作查询实例 ID。
* `account_name` - (必填，变更时重建) 数据库账号名称。由小写字母、数字或下划线组成。对于 MySQL，还支持大写字母。它必须以字母开头，并以字母或数字结尾。
  - **长度**：
    - MySQL 8.0 和 5.7：2 到 32 个字符。
    - MySQL 5.6：2 到 16 个字符。
    - SQL Server：2 到 64 个字符。
    - PostgreSQL 云盘版：2 到 63 个字符。
    - PostgreSQL 本地盘版：2 到 16 个字符。
    - MariaDB：2 到 16 个字符。
  - **注意**：普通账号名和高权限账号名不能相似。例如，如果高权限账号名为 `Test1`，那么普通账号名不能为 `test1`。
* `password` - (必填) 操作密码。可以包含字母、数字或下划线，长度为 6 到 32 个字符。
* `kms_encrypted_password` - (可选) 使用 KMS 加密的数据库账号密码。如果设置了 `password`，此字段将被忽略。
* `kms_encryption_context` - (可选) 用于在创建或更新数据库账号时解密 `kms_encrypted_password` 的 KMS 加密上下文。详见 [加密上下文](https://www.alibabacloud.com/help/doc-detail/42975.htm)。当设置了 `kms_encrypted_password` 时有效。
* `account_type` - (可选，变更时重建) 账号类型。有效值：
  - `Normal`：普通账号(默认)。
  - `Super`：高权限账号。
  - `Sysadmin`：具有 SA 权限的数据库账号(仅支持 RDS SQL Server 实例)。
  - **注意**：在创建具有 SA 权限的数据库账号之前，请检查实例是否满足先决条件。更多信息，请参阅 [创建具有 SA 权限的数据库账号](https://www.alibabacloud.com/help/doc-detail/122334.htm)。
* `account_description` - (可选) 账号描述。长度可以为 2 到 256 个字符。它以中文或英文字母开头，可以包含数字、中文、英文、下划线 (`_`) 和连字符 (`-`)。
  - **注意**：不能以 `http://` 或 `https://` 开头。
* `instance_id` - (可选，变更时重建) 已废弃字段，建议使用 `data_base_instance_id` 替代。
* `name` - (可选，变更时重建) 已废弃字段，建议使用 `account_name` 替代。
* `type` - (可选，变更时重建) 已废弃字段，建议使用 `account_type` 替代。
* `description` - (可选) 已废弃字段，建议使用 `account_description` 替代。

## 属性说明

除了上述所有参数外，还导出了以下属性：

* `id` - 当前账号资源 ID。由实例 ID 和账号名称组成，格式为 `<instance_id>:<account_name>`。
* `data_base_instance_id` - 实例 ID。您可以通过调用 `DescribeDBInstances` 操作查询实例 ID。
* `account_name` - 数据库账号名称。由小写字母、数字或下划线组成。对于 MySQL，还支持大写字母。它必须以字母开头，并以字母或数字结尾。
* `account_type` - 账号类型。有效值：`Normal`、`Super` 或 `Sysadmin`。
* `account_description` - 账号描述。长度可以为 2 到 256 个字符。它以中文或英文字母开头，可以包含数字、中文、英文、下划线 (`_`) 和连字符 (`-`)。
* `instance_id` - 已废弃字段，建议使用 `data_base_instance_id` 替代。
* `name` - 已废弃字段，建议使用 `account_name` 替代。
* `type` - 已废弃字段，建议使用 `account_type` 替代。
* `description` - 已废弃字段，建议使用 `account_description` 替代。