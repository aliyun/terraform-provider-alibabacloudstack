---
subcategory: "KMS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_kms_secret"
sidebar_current: "docs-alibabacloudstack-resource-kms-secret"
description: |-
  编排KMS密钥
---

# alibabacloudstack_kms_secret

使用Provider配置的凭证在指定的资源集编排KMS密钥。



## 示例用法

### 基础用法

```
resource "alibabacloudstack_kms_secret" "default" {
  secret_name                   = "secret-foo"
  description                   = "from terraform"
  secret_data                   = "Secret data."
  version_id                    = "000000000001"
  force_delete_without_recovery = true
}
```

## 参数参考

以下参数受支持：

* `description` - (可选) 密钥的描述。
* `encryption_key_id` - (可选，变更时重建) 用于加密密钥值的 KMS CMK 的 ID。如果您不指定此参数，密钥管理器将自动创建一个加密密钥来加密密钥。
* `force_delete_without_recovery` - (可选) 指定是否强制删除密钥。如果此参数设置为 true，则密钥无法恢复。有效值：true, false。默认值为：false。
* `recovery_window_in_days` - (可选) 如果您不强制删除它，密钥的恢复期。默认值：30。当 `force_delete_without_recovery` 为 true 时将被忽略。
* `secret_data` - (必填) 您要创建的密钥值。密钥管理器会加密密钥值并将其存储在初始版本中。
* `secret_data_type` - (可选) 密钥值的类型。有效值：text, binary。默认为 "text"。
* `secret_name` - (必填，变更时重建) 密钥的名称。
* `version_id` - (必填) 初始版本的版本号。每个密钥对象中的版本号是唯一的。
* `version_stages` - (可选，字符串列表)标记新密钥版本的阶段标签。如果您不指定此参数，密钥管理器会将其标记为 "ACSCurrent"。
* `tags` - (可选) 分配给资源的标签映射。
* `encryption_key_id` - (可选，变更时重建) 用于加密密钥值的 KMS CMK 的 ID。如果您不指定此参数，密钥管理器将自动创建一个加密密钥来加密密钥。

## 属性参考

* `id` - 密钥的 ID。它与 `secret_name` 相同。
* `arn` - 密钥的 Alibabacloudstack 资源名称 (ARN)。
* `planned_delete_time` - 密钥计划被删除的时间。