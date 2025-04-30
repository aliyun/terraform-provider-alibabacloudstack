---
subcategory: "AnalyticDB for MySQL (ADB)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_adb_account"
sidebar_current: "docs-Alibabacloudstack-adb-account"
description: |- 
  编排adb数据用户
---

# alibabacloudstack_adb_account

使用Provider配置的凭证在指定的资源集下编排adb数据用户。

## 示例用法

```hcl
variable "name" {
    default = "tf-testaccadbaccount96136"
}

variable "password" {
    default = "TestPassword123"
}

resource "alibabacloudstack_adb_account" "default" {
  db_cluster_id       = "am-bp1j43v9c35ef2cvf"
  account_name        = "nametest123"
  account_password    = var.password
  account_type        = "Normal"
  account_description = var.name
}
```

## 参数参考

支持以下参数：

* `db_cluster_id` - (必填，变更时重建) ADB集群的ID，该账户所属的集群ID。一旦设置，无法更改。
* `account_name` - (必填，变更时重建) 账户名称。它必须以字母开头，可以包含小写字母、数字和下划线(_)。长度不应超过16个字符。
* `account_password` - (选填)账户密码。它必须由字母、数字或下划线组成，长度在6到32个字符之间。必须指定 `account_password` 或 `kms_encrypted_password`。
* `kms_encrypted_password` - (选填)用于创建或更新数据库账户的KMS加密密码。如果提供了 `account_password`，此字段将被忽略。
* `kms_encryption_context` - (选填)用于在创建或更新数据库账户之前解密 `kms_encrypted_password` 的KMS加密上下文。仅当设置了 `kms_encrypted_password` 时有效。
* `account_type` - (选填，变更时重建) 数据库账户的类型。默认值：`Normal`。有效值：
  * `Normal`：标准账户。每个集群最多可以创建256个标准账户。
  * `Super`：特权账户。每个集群只能创建一个特权账户。
* `account_description` - (选填)账户描述。它必须以中文字符或英文字母开头，并可以包括中文字符、英文字母、下划线(_)、连字符(-)和数字。长度应在2到256个字符之间。

## 属性参考

除了上述所有参数外，还导出了以下属性：

* `id` - 账户的唯一标识符。它由实例ID和账户名称组成，格式为 `<instance_id>:<account_name>`。