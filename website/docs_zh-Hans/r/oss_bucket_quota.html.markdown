---
subcategory: "OSS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_oss_bucket_quota"
sidebar_current: "docs-Alibabacloudstack-oss-bucket-quota"
description: |-
  编排对象存储服务（OSS）配额资源
---

# alibabacloudstack_oss_bucket_quota

提供一个 OSS Bucket 配额资源。
使用Provider配置的凭证在指定的资源集编排对象存储服务（OSS）配额资源。
**注意：** 此资源已弃用。`oss_bucket` 资源已经包含用于管理配额资源的相应功能。

## 示例用法

```hcl
resource "alibabacloudstack_oss_bucket" "default" {
  bucket = "my-tf-test-bucket"
}

resource "alibabacloudstack_oss_bucket_quota" "default" {
  bucket = alibabacloudstack_oss_bucket.default.bucket
  quota  = 10240  # 配额以 MB 为单位
}
```

## 参数说明

以下参数被支持：

* `bucket` - (必填，变更时重建) - 指定需要设置配额的 OSS 存储桶名称。此参数必须与已创建的 OSS 存储桶名称一致。
* `quota` - (必填，变更时重建) - 设置 OSS 存储桶的存储配额，单位为兆字节（MB）。通过此参数可以限制存储桶的最大存储容量。

## 属性说明

除了上述列出的参数外，还导出以下属性：

* `bucket` - 返回所设置配额的 OSS 存储桶名称。
* `quota` - 返回 OSS 存储桶的存储配额，单位为兆字节（MB）。此值与设置时的值一致，表示当前存储桶的存储限制。