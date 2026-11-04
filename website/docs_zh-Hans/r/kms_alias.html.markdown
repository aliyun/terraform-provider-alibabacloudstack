---
subcategory: "密钥管理服务"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_kms_alias"
sidebar_current: "docs-Alibabacloudstack-resource-kms-alias"
description: |-
  编排加密主密钥 (CMK) 的别名
---

# alibabacloudstack_kms_alias

使用Provider配置的凭证在指定的资源集下编排加密主密钥 (CMK) 的别名。

-> **注意：** 支持的凭证：目标组织（主账号）的 AK/SK，或目标组织下 ASCM 侧的用户并使用配置正确的自定义角色。ASCM 侧的用户有如下要求：

* 裸 AK/SK 不行：KMS 权限是通过角色授予的，而不是直接挂在用户上，须使用配置正确的自定义角色。
* 跨组织账号不行（包括环境管理员账号 admin）：KMS 请求会被拒绝，报 `Forbidden.NoPermission`。
* RAM 角色不行：ASCM 无法识别 RAM 角色会话的身份。
* 所需权限点：`kms:CreateAlias`、`kms:DescribeKey`、`kms:UpdateAlias`、`kms:DeleteAlias`。

## 示例用法

### 基础用法

```
resource "alibabacloudstack_kms_key" "key" {}

resource "alibabacloudstack_kms_alias" "alias" {
  alias_name = "alias/test_kms_alias"
  key_id     = alibabacloudstack_kms_key.key.id
}
```

## 参数说明

支持以下参数：

* `alias_name` - (必填，变更时重建) CMK 的别名。可以使用别名调用 `Encrypt`、`GenerateDataKey` 和 `DescribeKey`。字符长度(不包括前缀)：最小长度为 1 个字符，最大长度为 255 个字符。必须包含前缀 `alias/`。
* `key_id` - (必填) 密钥的 ID。

-> **注意：** 每个别名只能代表一个主密钥 (CMK)。

-> **注意：** 在同一区域的同一用户内，别名不可重复。

-> **注意：** 可以使用 `UpdateAlias` 更新别名与主密钥 (CMK) 的映射关系。


## 属性说明

* `id` - 别名的唯一标识符，用于标识该别名资源。

## Import

KMS Alias 可以使用 alias_name 导入，例如:

```
$ terraform import alibabacloudstack_kms_alias.example alias/test_kms_alias
```