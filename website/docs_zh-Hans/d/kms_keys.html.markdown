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

## 参数参考

以下参数被支持：

* `ids` - (可选) KMS 密钥 ID 列表。
* `description_regex` - (可选) 用于通过 KMS 密钥描述筛选结果的正则表达式字符串。
* `status` - (可选) 通过 KMS 密钥状态筛选结果。有效值：`Enabled`, `Disabled`, `PendingDeletion`。

## 属性参考

除了上述列出的参数外，还导出以下属性：

* `ids` - KMS 密钥 ID 列表。
* `keys` - KMS 密钥列表。每个元素包含以下属性：
  * `id` - 密钥的 ID。
  * `arn` - 密钥的 Alibabacloudstack Cloud 资源名称 (ARN)。
  * `description` - 密钥的描述。
  * `status` - 密钥的状态。可能的值：`Enabled`, `Disabled` 和 `PendingDeletion`。
  * `creation_date` - 密钥的创建日期。
  * `delete_date` - 密钥的删除日期。
  * `creator` - 密钥的所有者。
  * `computed_property` - 表示密钥的计算属性。