---
subcategory: "OSS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_oss_bucket_kms"
sidebar_current: "docs-Alibabacloudstack-oss-bucket-kms"
description: |-
  编排对象存储服务（OSS）加密配置
---

# alibabacloudstack_oss_bucket_kms

使用Provider配置的凭证在指定的资源集编排对象存储服务（OSS）加密配置。
**注意：** 此资源已弃用。`oss_bucket` 资源已经包含用于管理使用 KMS 的服务器端加密的相应功能。

## 示例用法

```hcl
resource "alibabacloudstack_oss_bucket" "default" {
  bucket = "my-tf-test-bucket"
}

resource "alibabacloudstack_oss_bucket_kms" "default" {
  bucket          = alibabacloudstack_oss_bucket.default.bucket
  sse_algorithm   = "KMS"
  kms_data_encryption = "AES256"
  kms_master_key_id = "your_kms_master_key_id"
}
```

## 参数参考

支持以下参数：

* `bucket` - (必填，变更时重建) OSS 存储桶的名称。
* `sse_algorithm` - (必填，变更时重建) 服务器端加密算法。有效值：KMS。
* `kms_data_encryption` - (可选，变更时重建) KMS 使用的数据加密算法。有效值：AES256。
* `kms_master_key_id` - (可选，变更时重建) 用于加密的 KMS 主密钥 ID。

## 属性参考

除了上述参数列表中的参数外，还导出以下属性：

* `bucket` - OSS 存储桶的名称。
* `sse_algorithm` - 服务器端加密算法。
* `kms_data_encryption` - KMS 使用的数据加密算法。
* `kms_master_key_id` - 用于加密的 KMS 主密钥 ID。