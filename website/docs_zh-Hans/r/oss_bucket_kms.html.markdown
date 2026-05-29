---
subcategory: "对象存储 OSS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_oss_bucket_kms"
sidebar_current: "docs-Alibabacloudstack-resource-oss-bucket-kms"
description: |-
  编排对象存储服务（OSS）加密配置
---

# alibabacloudstack_oss_bucket_kms

使用 Provider 配置的凭证在指定的资源集编排对象存储服务（OSS）加密配置。

**注意：** 此资源已弃用，计划在 3.21.0 版本中移除。`oss_bucket` 资源已包含用于管理使用 KMS 的服务器端加密的相应功能。

## 示例用法

```hcl
resource "alibabacloudstack_oss_bucket" "default" {
  bucket = "my-tf-test-bucket"
}

resource "alibabacloudstack_oss_bucket_kms" "default" {
  bucket          = alibabacloudstack_oss_bucket.default.bucket
  oss_cluster     = "your-oss-cluster-id"
  sse_algorithm   = "KMS"
  kms_master_key_id = "your_kms_master_key_id"
}
```

## 参数说明

支持以下参数：

* `bucket` - (必填，ForceNew) OSS 存储桶的名称。此参数用于指定要配置加密的 OSS 存储桶。
* `oss_cluster` - (可选，ForceNew) OSS 集群的 ID。修改此参数会强制重新创建资源。
* `sse_algorithm` - (必填，ForceNew) 服务器端加密算法。有效值为：`KMS`，表示使用 KMS 进行服务器端加密。修改此参数会强制重新创建资源。
* `kms_master_key_id` - (可选，ForceNew) 用于加密的 KMS 主密钥 ID。当 `sse_algorithm` 设置为 `KMS` 时需要指定此参数。修改此参数会强制重新创建资源。

## 属性说明

除了上述参数列表中的参数外，还导出以下属性：

* `bucket` - OSS 存储桶的名称。此属性返回配置加密的 OSS 存储桶名称。
* `oss_cluster` - OSS 集群的 ID。
* `sse_algorithm` - 服务器端加密算法。此属性返回当前存储桶使用的服务器端加密算法。
* `kms_master_key_id` - 用于加密的 KMS 主密钥 ID。此属性返回当前存储桶中用于加密的 KMS 主密钥 ID。

## Import

OSS Bucket KMS 可以通过 `oss_cluster:bucket_name` 格式导入，例如：

```
$ terraform import alibabacloudstack_oss_bucket_kms.example oss-cluster-id:my-tf-test-bucket
```