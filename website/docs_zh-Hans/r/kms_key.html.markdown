---
subcategory: "密钥管理服务"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_kms_key"
sidebar_current: "docs-Alibabacloudstack-resource-kms-key"
description: |-
  编排KMS密钥
---

# alibabacloudstack_kms_key

> **注意：** 此资源也可以使用以下别名进行引用：
> - `apsarastack_kms_key`

使用Provider配置的凭证在指定的资源集编排KMS密钥。
KMS密钥可以帮助用户在传输过程中保护数据安全。

## 示例用法

### 基础用法

```
resource "alibabacloudstack_kms_key" "key" {
  description             = "Hello KMS"
  pending_window_in_days  = "7"
  key_state               = "Enabled"
}
```
## 参数说明

支持以下参数：

* `key_usage` - (可选，变更时重建) 指定CMK的用途。取值：`ENCRYPT/DECRYPT`、`SIGN/VERIFY`。默认值：`ENCRYPT/DECRYPT`。
* `origin` - (可选，变更时重建) CMK的密钥材料来源。取值：`Aliyun_KMS`、`EXTERNAL`。默认值：`Aliyun_KMS`。
* `protection_level` - (可选，变更时重建) CMK的保护级别。取值：`SOFTWARE`、`HSM`。默认值：`SOFTWARE`。
* `description` - (可选) 在Alibabacloudstack控制台中查看的密钥描述。
* `automatic_rotation` - (可选) 指定是否启用自动密钥轮换。取值：`Enabled`、`Disabled`。默认值：`Disabled`。
* `key_state` - (可选) CMK的状态。取值：`Enabled`、`Disabled`、`PendingDeletion`。默认值：`Enabled`。
* `pending_window_in_days` - (可选) 销毁资源后密钥删除前的等待天数，必须在7到30天之间。默认值：`7`。
* `rotation_interval` - (可选) 自动密钥轮换的周期。格式为数字后跟时间单位（d表示天，h表示小时，m表示分钟，s表示秒），例如 `7d` 或 `2678400s`。仅当 `automatic_rotation` 设置为 `Enabled` 时有效。
* `is_enabled` - (可选，已弃用) 字段 `is_enabled` 已从 provider 版本 1.85.0 开始弃用。请使用 `key_state` 替代。
* `deletion_window_in_days` - (可选，已弃用) 字段 `deletion_window_in_days` 已从 provider 版本 1.85.0 开始弃用。请使用 `pending_window_in_days` 替代。

-> **注意：** 当预删除等待期结束后，密钥将被永久删除且无法恢复。


## 属性说明

支持以下属性：

* `id` - 密钥的ID。
* `arn` - 密钥的Alibabacloudstack资源名称(ARN)。
* `creation_date` - CMK创建的日期和时间。时间以UTC显示。
* `creator` - CMK的创建者。
* `delete_date` - CMK计划删除的日期。时间以UTC显示。只有当KeyState值为PendingDeletion时，此值才返回。
* `key_state` - CMK的状态。
* `key_usage` - CMK的用途。
* `origin` - CMK的密钥材料来源。
* `protection_level` - CMK的保护级别。
* `primary_key_version` - 对称CMK当前主密钥版本的ID。
* `automatic_rotation` - 表示是否为密钥启用了自动轮换。
* `description` - 在Alibabacloudstack控制台中查看的密钥描述。
* `last_rotation_date` - 上次轮换执行的日期和时间。时间以UTC显示。
* `material_expire_time` - CMK的密钥材料过期的时间和日期。时间以UTC显示。如果值为空，则CMK的密钥材料不会过期。
* `next_rotation_date` - 下次轮换计划执行的时间。
* `rotation_interval` - 自动密钥轮换的周期。

## Import

KMS Key 可以使用 key ID 进行导入，例如：

```
$ terraform import alibabacloudstack_kms_key.example key-12345678
```