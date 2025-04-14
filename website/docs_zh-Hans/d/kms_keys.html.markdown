---
subcategory: "KMS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_kms_keys"
sidebar_current: "docs-alibabacloudstack-datasource-kms-keys"
description: |-
    查询KMS密钥列表
---

# alibabacloudstack_kms_keys

根据指定过滤条件列出当前凭证权限可以访问的KMS密钥列表。

## 示例用法

```
# 声明数据源
data "alibabacloudstack_kms_keys" "kms_keys_ds" {
  description_regex = "Hello KMS"
}

output "first_key_id" {
  value = "${data.alibabacloudstack_kms_keys.kms_keys_ds.keys.0.id}"
}
```

## 参数说明

以下参数被支持：

* `ids` - (可选) 用于筛选结果的KMS密钥ID列表。
* `description_regex` - (可选) 用于通过KMS密钥描述筛选结果的正则表达式字符串。此参数可以帮助您精确匹配特定描述的密钥。
* `status` - (可选) 用于通过KMS密钥状态筛选结果。有效值包括：`Enabled`（已启用）、`Disabled`（已禁用）和`PendingDeletion`（待删除）。

## 属性说明

除了上述列出的参数外，还导出以下属性：

* `ids` - 匹配条件后返回的KMS密钥ID列表。
* `keys` - 匹配条件后返回的KMS密钥详细信息列表。每个元素包含以下属性：
  * `id` - 密钥的唯一标识符。
  * `arn` - 密钥的阿里云资源名称（ARN），用于唯一标识该密钥。
  * `description` - 密钥的描述信息，帮助用户识别密钥用途。
  * `status` - 密钥的状态。可能的值为：`Enabled`（已启用）、`Disabled`（已禁用）和`PendingDeletion`（待删除）。
  * `creation_date` - 密钥的创建时间，格式为标准时间戳。
  * `delete_date` - 如果密钥处于`PendingDeletion`状态，则此字段表示密钥的计划删除时间。
  * `creator` - 密钥的创建者或所有者信息。
  * `computed_property` - 密钥的计算属性，通常由系统生成，具体含义取决于密钥类型及其配置。