---
subcategory: "KMS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_kms_key"
sidebar_current: "docs-alibabacloudstack-resource-kms-key"
description: |-
  编排KMS密钥
---

# alibabacloudstack_kms_key

使用Provider配置的凭证在指定的资源集编排KMS密钥对。
KMS密钥对可以帮助用户在传输过程中保护数据安全。

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

* `description` - (可选) 在Alibabacloudstack控制台中查看的密钥描述。
* `key_usage` - (可选，变更时重建) 指定CMK的用途。目前，默认为'ENCRYPT/DECRYPT'，表示CMK用于加密和解密。
* `automatic_rotation` - (可选) 指定是否启用自动密钥轮换。默认："Disabled"。
* `is_enabled` - (可选，已弃用) 字段 'is_enabled' 已被弃用。使用新的字段 'key_state' 替代。
* `key_state` - (可选) CMK的状态。默认为启用状态（Enabled）。
* `origin` - (可选，变更时重建) CMK的密钥材料来源。默认为 "Aliyun_KMS"。
* `deletion_window_in_days` - (可选，已弃用) 字段 'deletion_window_in_days' 已被弃用。使用新的字段 'pending_window_in_days' 替代。
* `pending_window_in_days` - (可选) 销毁资源后密钥删除前的天数，必须在7到30天之间。默认为30天。
* `protection_level` - (可选，变更时重建) CMK的保护级别。默认为 "SOFTWARE"。
* `rotation_interval` - (可选) 自动密钥轮换的周期。单位：秒。

-> **注意：** 当预删除天数到期后，密钥将被永久删除且无法恢复。


## 属性说明

* `id` - 密钥的ID。
* `arn` - 密钥的Alibabacloudstack资源名称(ARN)。
* `creation_date` - CMK创建的日期和时间。时间以UTC显示。
* `creator` - CMK的创建者。
* `delete_date` - CMK计划删除的日期。时间以UTC显示。只有当KeyState值为PendingDeletion时，此值才返回。
* `last_rotation_date` - 上次轮换执行的日期和时间。时间以UTC显示。
* `material_expire_time` - CMK的密钥材料过期的时间和日期。时间以UTC显示。如果值为空，则CMK的密钥材料不会过期。
* `next_rotation_date` - 下次轮换计划执行的时间。
* `primary_key_version` - 对称CMK当前主密钥版本的ID。
* `automatic_rotation` - 表示是否为密钥启用了自动轮换的属性。
* `description` - 在Alibabacloudstack控制台中查看的密钥描述。