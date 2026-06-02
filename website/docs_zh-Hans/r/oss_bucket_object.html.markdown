---
subcategory: "对象存储 OSS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_oss_bucket_object"
sidebar_current: "docs-Alibabacloudstack-resource-oss-bucket-object"
description: |-
  将对象(内容或文件)上传到OSS存储桶
---

# alibabacloudstack_oss_bucket_object

使用Provider配置的凭证将对象(内容或文件)上传到OSS存储桶。

## 示例用法

### 调整provider.tf中的配置

```hcl
provider "alibabacloudstack" {
  popgw_domain            = "inter.env205.shuguang.com"
  access_key              = "xxx"
  secret_key              = "xxx"
  region                  = "cn-center-a"
  insecure                = "true"
  resource_group_set_name = "ResourceSet(xxx)"
}
```

### 首先创建一个新的存储桶

```hcl
resource "alibabacloudstack_oss_bucket" "example" {
  bucket = "your_bucket_name"
  acl    = "public-read"
}
```

### 将文件上传到存储桶

```hcl
resource "alibabacloudstack_oss_bucket_object" "object-source" {
  bucket  = alibabacloudstack_oss_bucket.example.bucket
  key     = "new_object_key"
  source  = "path/to/file"
}
```

### 将内容上传到存储桶

```hcl
resource "alibabacloudstack_oss_bucket_object" "object-content" {
  bucket  = alibabacloudstack_oss_bucket.example.bucket
  key     = "new_object_key"
  content = "the content that you want to upload."
}
```

### 上传带有服务端加密的文件

```hcl
resource "alibabacloudstack_oss_bucket_object" "encrypted-object" {
  bucket                 = alibabacloudstack_oss_bucket.example.bucket
  key                    = "encrypted_object_key"
  source                 = "path/to/file"
  server_side_encryption = "AES256"
}
```

## 参数说明

-> **注意:** 如果您指定了`content_encoding`，则需要负责正确编码主体(即`source`和`content`都期望已经编码/压缩的字节)

以下是支持的参数：

* `bucket` - (必填，变更后重建) 要上传文件的目标存储桶名称。修改此参数会强制重新创建资源。
* `key` - (必填，变更后重建) 对象在存储桶中的名称。修改此参数会强制重新创建资源。
* `oss_cluster` - (可选，变更后重建) OSS集群标识符。如果未指定，将使用默认集群。修改此参数会强制重新创建资源。
* `source` - (可选) 要上传到存储桶的源文件路径。与`content`互斥。必须提供`source`或`content`之一。
* `content` - (可选) 要上传到存储桶的字面内容。与`source`互斥。必须提供`source`或`content`之一。
* `acl` - (可选) 要应用的[标准ACL](https://www.alibabacloud.com/help/doc-detail/52284.htm)。有效值：`private`、`public-read`、`public-read-write`。默认值为`private`。
* `content_type` - (可选，Computed) 描述对象数据格式的标准MIME类型，例如`application/octet-stream`。所有有效的MIME类型都可以作为此输入。
* `cache_control` - (可选) 指定请求/响应链中的缓存行为。阅读[RFC2616 Cache-Control](https://www.ietf.org/rfc/rfc2616.txt)获取更多详细信息。
* `content_disposition` - (可选) 指定对象的呈现信息。阅读[RFC2616 Content-Disposition](https://www.ietf.org/rfc/rfc2616.txt)获取更多详细信息。
* `content_encoding` - (可选) 指定已应用于对象的内容编码，因此必须应用哪些解码机制才能获得Content-Type头字段引用的媒体类型。阅读[RFC2616 Content-Encoding](https://www.ietf.org/rfc/rfc2616.txt)获取更多详细信息。
* `content_md5` - (可选，Computed) 内容的MD5值。阅读[MD5](https://www.alibabacloud.com/help/doc-detail/31978.htm)获取计算方法。
* `expires` - (可选) 指定请求/响应的过期日期。必须使用RFC1123格式。阅读[RFC2616 Expires](https://www.ietf.org/rfc/rfc2616.txt)获取更多详细信息。
* `server_side_encryption` - (可选) 指定OSS中的对象服务器端加密。有效值：`AES256`、`KMS`。
* `kms_key_id` - (可选，Computed) 指定由KMS管理的主要密钥。当`server_side_encryption`的值设置为`KMS`时，此参数有效。

## 属性说明

以下属性被导出：

* `id` - 资源的ID，格式为`{oss_cluster}:{bucket}:{key}`。
* `version_id` - 如果启用了存储桶版本控制，则为对象提供唯一的版本ID值。
* `content_type` - 对象数据的MIME类型。
* `content_md5` - 对象内容的MD5哈希值。
* `kms_key_id` - 用于服务端加密的KMS密钥ID(当`server_side_encryption`为`KMS`时)。

## 导入

OSS存储桶对象可以使用复合ID导入，格式为`{oss_cluster}:{bucket}:{key}`，例如：

```shell
$ terraform import alibabacloudstack_oss_bucket_object.example default:my-bucket:my-object-key
```