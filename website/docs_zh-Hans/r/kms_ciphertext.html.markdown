---
subcategory: "KMS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_kms_ciphertext"
sidebar_current: "docs-alibabacloudstack-datasource-kms-ciphertext"
description: |-
    使用KMS加密给定的明文
---

# alibabacloudstack_kms_ciphertext

使用Provider配置的凭证在指定的资源集下使用KMS加密给定的明文。

> **注意**：使用此数据源可以让你在资源定义中隐藏秘密数据，但并不能保证在所有 Terraform 日志和状态输出中保护这些数据。请确保在 Terraform 配置之外对秘密数据进行安全保护。

## 示例用法

```
resource "alibabacloudstack_kms_key" "key" {
  description             = "example key"
  is_enabled              = true
}

data "alibabacloudstack_kms_ciphertext" "encrypted" {
  key_id    = alibabacloudstack_kms_key.key.id
  plaintext = "example"
}

output "alibabacloudstack_kms_ciphertext" {
  value = "${data.alibabacloudstack_kms_ciphertext.encrypted}"
}
```

## 参数说明

支持以下参数：

* `plaintext` - (必选) 要加密的明文，必须以 Base64 编码。
* `key_id` - (必选) CMK 的全局唯一 ID。
* `encryption_context` - (可选) 加密上下文。如果你在此处指定此参数，则在调用 Decrypt API 操作时也需要提供它。
* `sensitive` - (必选) 表示该属性是否为敏感信息。

## 属性说明

除了上述参数外，还导出以下属性：

* `ciphertext_blob` - 使用主 CMK 版本加密的数据密钥的密文。
* `sensitive` - 自动计算的属性，表示该属性是否为敏感信息。