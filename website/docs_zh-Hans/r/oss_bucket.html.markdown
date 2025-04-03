---
subcategory: "OSS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_oss_bucket"
sidebar_current: "docs-alibabacloudstack-resource-oss-bucket"
description: |-
  编排对象存储服务（OSS）存储桶
---

# alibabacloudstack_oss_bucket

使用Provider配置的凭证在指定的资源集编排对象存储服务（OSS）存储桶并设置其属性的资源。

-> **注意:** OSS 系统中的存储桶命名空间是所有用户共享的。请尽量将存储桶名称设置为唯一值。


## 示例用法

私有存储桶

```
resource "alibabacloudstack_oss_bucket" "demo" {
  bucket = "sample_bucket"
  acl    = "public-read"
}
```

## 参数参考

以下是支持的参数：

* `bucket` - (可选，变更时重建) 存储桶的名称。如果省略，Terraform 将分配一个随机且唯一的名称。
* `acl` - (可选) 可以为 "private", "public-read" 和 "public-read-write"。默认为 "private"。
* `logging` - (可选) 日志记录对象支持以下内容：
  * `target_bucket` - (必填)  将接收日志对象的存储桶名称。
  * `target_prefix` - (可选) 为日志对象指定键前缀。 
* `storage_class` - (可选，变更时重建) 对象存储类型。可能的值：`Standard`, `IA` 和 `Archive`。
* `vpclist` - (可选) 可访问 VPC 的列表
* `storage_capacity` - (可选) 为存储桶设置容量限制。如果达到容量限制，则写入操作将被拒绝。单位 GB
* `sse_algorithm` - (可选) 加密上传到 OSS 的文件。它可以是 "", "AES256", "SM4" 和 "KMS"。默认为 ""。
* `kms_key_id` - (可选) 当 sse_algorithm 为 KMS 时需要设置。
* `bucket_sync` - (可选) 启用或禁用存储桶同步。默认为 `true`。

## 属性参考

导出以下属性：

* `id` - 存储桶的名称。
* `acl` - 存储桶的 ACL。
* `creation_date` - 存储桶的创建日期。
* `extranet_endpoint` - 存储桶的公网访问端点。
* `intranet_endpoint` - 存储桶的内网访问端点。
* `location` - 存储桶的位置。
* `owner` - 存储桶的所有者。
* `storage_class` - 存储桶的存储类别。
* `vpclist` - 存储桶的可访问 VPC 列表。