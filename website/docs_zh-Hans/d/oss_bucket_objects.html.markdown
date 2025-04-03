---
subcategory: "OSS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_oss_bucket_objects"
sidebar_current: "docs-alibabacloudstack-datasource-oss-bucket-objects"
description: |-
    查询块存储（OSS）存储桶中的对象
---

# alibabacloudstack_oss_bucket_objects

根据指定过滤条件列出当前凭证权限可以访问的块存储（OSS）存储桶中的对象列表。

## 示例用法

```
data "alibabacloudstack_oss_bucket_objects" "bucket_objects_ds" {
  bucket_name = "sample_bucket"
  key_regex   = "sample/sample_object.txt"
}

output "first_object_key" {
  value = "${data.alibabacloudstack_oss_bucket_objects.bucket_objects_ds.objects.0.key}"
}
```

## 参数参考

支持以下参数：

* `bucket_name` - 包含要查找的对象的存储桶名称。
* `key_regex` - (可选) 用于通过键过滤结果的正则表达式字符串。
* `key_prefix` - (可选) 通过给定的键前缀过滤结果(例如“path/to/folder/logs-”)。

## 属性参考

除了上述参数外，还导出以下属性：

* `objects` - 存储桶对象的列表。每个元素包含以下属性：
  * `key` - 对象键。
  * `acl` - 对象访问控制列表。可能的值：`default`、`private`、`public-read` 和 `public-read-write`。
  * `content_type` - 描述对象数据格式的标准MIME类型，例如“application/octet-stream”。
  * `cache_control` - 沿请求/响应链的缓存行为。阅读 [RFC2616 Cache-Control](https://www.ietf.org/rfc/rfc2616.txt) 获取更多详细信息。
  * `content_disposition` - 对象的展示信息。阅读 [RFC2616 Content-Disposition](https://www.ietf.org/rfc/rfc2616.txt) 获取更多详细信息。
  * `content_encoding` - 已应用于对象的内容编码，因此必须应用哪些解码机制才能获得由Content-Type头字段引用的媒体类型。阅读 [RFC2616 Content-Encoding](https://www.ietf.org/rfc/rfc2616.txt) 获取更多详细信息。
  * `content_md5` - 内容的MD5值。阅读 [MD5](https://www.alibabacloud.com/help/doc-detail/31978.htm) 获取计算方法。
  * `expires` - 请求/响应的过期日期。阅读 [RFC2616 Expires](https://www.ietf.org/rfc/rfc2616.txt) 获取更多详细信息。
  * `server_side_encryption` - OSS中的对象服务器端加密。它可以为空或`AES256`。
  * `sse_kms_key_id` - 如果存在，则指定用于该对象的密钥管理服务(KMS)主加密密钥ID。
  * `storage_class` - 对象存储类型。可能的值：`Standard`、`IA` 和 `Archive`。
  * `last_modification_time` - 对象的最后修改时间。
  * `computed_attribute` - 计算属性的描述。