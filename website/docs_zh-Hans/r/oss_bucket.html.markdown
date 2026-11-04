---
subcategory: "对象存储 OSS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_oss_bucket"
sidebar_current: "docs-Alibabacloudstack-resource-oss-bucket"
description: |-
  编排对象存储服务（OSS）存储桶
---

# alibabacloudstack_oss_bucket

使用Provider配置的凭证在指定的资源集编排对象存储服务（OSS）存储桶并设置其属性的资源。

-> **注意:** OSS 系统中的存储桶命名空间是所有用户共享的。请尽量将存储桶名称设置为唯一值。

-> **注意：** 支持的凭证：目标组织（主账号）的 AK/SK，或目标组织下 ASCM 侧的用户并使用配置正确的自定义角色。ASCM 侧的用户有如下要求：

* 裸 AK/SK 不行：OSS 权限是通过角色授予的，而不是直接挂在用户上，须使用配置正确的自定义角色。
* 跨组织账号不行（包括环境管理员账号 admin）：OSS 请求会被拒绝，报 `AccessDenied`（"Put bucket request is not resource owner."，EC 0003-00000202）。
* RAM 角色不行：ASCM 无法识别 RAM 角色会话的身份。
* 所需权限点：`oss:PutBucket`、`oss:GetBucketAcl`、`oss:PutBucketAcl`、`oss:GetBucketTags`、`oss:PutBucketTags`、`oss:DeleteBucketTags`、`oss:GetBucketLogging`、`oss:PutBucketLogging`、`oss:DeleteBucketLogging`、`oss:GetBucketEncryption`、`oss:PutBucketEncryption`、`oss:DeleteBucketEncryption`、`oss:GetBucketSync`、`oss:PutBucketSync`、`oss:DeleteBucketSync`、`oss:GetBucketStorageCapacity`、`oss:SetBucketStorageCapacity`、`oss:DeleteBucket`，以及 `ascm:UpdateInstanceBelong`（存储桶创建完成后，Provider 会通过 ASCM 管理面将其移入资源集；该调用失败时，刚创建的存储桶会被删除并返回错误）。


## 示例用法

私有存储桶

```
resource "alibabacloudstack_oss_bucket" "demo" {
  bucket = "sample_bucket"
  acl    = "public-read"
}
```

## 参数说明

以下是支持的参数：

* `bucket` - (可选，变更时重建) 存储桶的名称。如果省略，Terraform 将分配一个随机且唯一的名称。
* `acl` - (可选) 可以为 "private", "public-read" 和 "public-read-write"。默认为 "private"。
* `logging` - (可选) 日志记录对象支持以下内容：
  * `target_bucket` - (必填) 将接收日志对象的存储桶名称。
  * `target_prefix` - (可选) 为日志对象指定键前缀。
* `storage_class` - (可选，变更时重建) 对象存储类型。可能的值：`Standard`（标准存储）、`IA`（低频访问存储）和 `Archive`（归档存储）。
* `vpclist` - (可选) 可访问 VPC 的列表。
* `storage_capacity` - (可选) 为存储桶设置容量限制。如果达到容量限制，则写入操作将被拒绝。单位为 GB。
* `sse_algorithm` - (可选) 加密上传到 OSS 的文件。它可以是 "", "AES256", "SM4" 和 "KMS"。默认为 ""。
* `kms_key_id` - (可选) 当 `sse_algorithm` 为 KMS 时需要设置此参数，用于指定加密所使用的 KMS 密钥 ID。
* `bucket_sync` - (可选) 启用或禁用存储桶同步功能。默认值为 `true`。
* `tags` - (可选) 存储桶的标签映射。标签 `ascm:tag/rg` 由 provider 管理，不允许在 `tags` 中设置。

## 属性说明

导出以下属性：

* `id` - 存储桶的名称。
* `acl` - 存储桶的访问控制列表（ACL）。
* `creation_date` - 存储桶的创建日期。
* `extranet_endpoint` - 存储桶的公网访问端点。
* `intranet_endpoint` - 存储桶的内网访问端点。
* `location` - 存储桶所在的区域位置。
* `owner` - 存储桶的所有者信息。
* `storage_class` - 存储桶的存储类别，例如 `Standard`、`IA` 或 `Archive`。
* `vpclist` - 存储桶允许访问的 VPC 列表。

## 导入

OSS 存储桶可以使用以下导入命令导入：

```
$ terraform import alibabacloudstack_oss_bucket.example <oss_cluster>:<bucket>
```

其中 `oss_cluster` 是 OSS 集群的 ID，`bucket` 是存储桶的名称。