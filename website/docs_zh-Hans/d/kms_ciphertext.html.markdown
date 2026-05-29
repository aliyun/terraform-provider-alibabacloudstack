---
subcategory: "Key Management Service"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_kms_ciphertexts"
sidebar_current: "docs-alibabacloudstack-datasource-kms-ciphertext"
description: |-
    查询KMS加密数据。
---

# alibabacloudstack_kms_ciphertexts

使用指定KMS加密给定的明文。 

~> **注意**: 使用此数据源可以让你在资源定义中隐藏秘密数据，但这并不会保护所有 Terraform 日志和状态输出中的数据。请确保在 Terraform 配置之外妥善保护你的秘密数据。

## 示例用法

```
resource "alibabacloudstack_kms_key" "key" {
  description             = "example key"
  is_enabled              = true
}

data "alibabacloudstack_kms_ciphertexts" "encrypted" {
  key_id    = alibabacloudstack_kms_key.key.id
  plaintext = "example"
}

output "alibabacloudstack_kms_ciphertexts" {
  value = "${data.alibabacloudstack_kms_ciphertexts.encrypted}"
}
```

## 参数说明

支持以下参数：

* `plaintext` - (必填) 要加密的明文数据。
* `key_id` - (必填) CMK 的全局唯一 ID。
* `encryption_context` - (可选) 加密上下文。如果你在此处指定此参数，则在调用 Decrypt API 操作时也需要提供该参数。

## 属性说明

除了上述参数外，还导出以下属性：

* `ciphertext_blob` - 使用主 CMK 版本加密的数据密钥的密文。