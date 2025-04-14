---
subcategory: "OSS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_oss_bucket_object"
sidebar_current: "docs-alibabacloudstack-resource-oss-bucket-object"
description: |-
  将对象(内容或文件)上传到OSS存储桶
---

# alibabacloudstack_oss_bucket_object

使用Provider配置的凭证将对象(内容或文件)上传到OSS存储桶。

## 示例用法

### 调整provider.tf中的配置
> 由于对对象的操作通过OSS数据网关进行，
> 因此需要在provider.tf中的provider配置中设置单独的数据网关地址。

```
provider "alibabacloudstack" {
  #popgw_domain            = "inter.env205.shuguang.com"
  access_key              = "xxx"
  ...
  insecure                = "true"
  resource_group_set_name = "xxx"
  ossservice_domain = "<oss data endpoint>"
}
```

### 首先创建一个新的存储桶
```
resource "alibabacloudstack_oss_bucket" "example" {
  bucket = "your_bucket_name"
  acl    = "public-read"
}
```

### 授予OSS特定用户的权限
> 目前仅支持通过ascm进行。

### 将文件上传到存储桶

```
resource "alibabacloudstack_oss_bucket_object" "object-source" {
  bucket  = "${alibabacloudstack_oss_bucket.example.bucket}"
  key    = "new_object_key"
  source = "path/to/file"
}
```

### 将内容上传到存储桶

```
resource "alibabacloudstack_oss_bucket_object" "object-content" {
  bucket  = "${alibabacloudstack_oss_bucket.example.bucket}"
  key     = "new_object_key"
  content = "the content that you want to upload."
}
```

## 参数说明

-> **注意:** 如果您指定了`content_encoding`，则需要负责正确编码主体(即`source`和`content`都期望已经编码/压缩的字节)

以下是支持的参数：

* `bucket` - (必填) 要上传文件的目标存储桶名称。
* `key` - (必填) 对象在存储桶中的名称。
* `source` - (可选) 要上传到存储桶的源文件路径。
* `content` - (可选，除非提供了`source`) 要上传到存储桶的字面内容。
* `acl` - (可选) 要应用的[标准ACL](https://www.alibabacloud.com/help/doc-detail/52284.htm)。默认为“private”。
* `content_type` - (可选) 描述对象数据格式的标准MIME类型，例如application/octet-stream。所有有效的MIME类型都可以作为此输入。
* `cache_control` - (可选) 指定请求/响应链中的缓存行为。阅读[RFC2616 Cache-Control](https://www.ietf.org/rfc/rfc2616.txt)获取更多详细信息。
* `content_disposition` - (可选) 指定对象的呈现信息。阅读[RFC2616 Content-Disposition](https://www.ietf.org/rfc/rfc2616.txt)获取更多详细信息。
* `content_encoding` - (可选) 指定已应用于对象的内容编码，因此必须应用哪些解码机制才能获得Content-Type头字段引用的媒体类型。阅读[RFC2616 Content-Encoding](https://www.ietf.org/rfc/rfc2616.txt)获取更多详细信息。
* `content_md5` - (可选) 内容的MD5值。阅读[MD5](https://www.alibabacloud.com/help/doc-detail/31978.htm)获取计算方法。
* `expires` - (可选) 指定请求/响应的过期日期。阅读[RFC2616 Expires](https://www.ietf.org/rfc/rfc2616.txt)获取更多详细信息。
* `server_side_encryption` - (可选) 指定OSS中的对象服务器端加密。有效值为`AES256`、`KMS`。默认值为`AES256`。
* `kms_key_id` - (可选，自1.62.1版本起可用) 指定由KMS管理的主要密钥。当`server_side_encryption`的值设置为KMS时，此参数有效。
* `version_id` - (可选) 如果启用了存储桶版本控制，则为对象提供唯一的版本ID值。

必须提供`source`或`content`之一来指定存储桶内容。这两个参数是互斥的。

## 属性说明

以下属性被导出：

* `id` - 上述资源的`key`。
* `version_id` - 如果启用了存储桶版本控制，则为对象提供唯一的版本ID值。