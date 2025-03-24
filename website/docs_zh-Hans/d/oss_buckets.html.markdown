---
subcategory: "OSS"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_oss_buckets"
sidebar_current: "docs-alibabacloudstack-datasource-oss-buckets"
description: |-
    查询块存储（OSS）存储桶
---

# alibabacloudstack_oss_buckets

根据指定过滤条件列出当前凭证权限可以访问的块存储（OSS）存储桶列表。

## 示例用法

```
data "alibabacloudstack_oss_buckets" "oss_buckets_ds" {
  name_regex = "sample_oss_bucket"
}

output "first_oss_bucket_name" {
  value = "${data.alibabacloudstack_oss_buckets.oss_buckets_ds.buckets.0.name}"
}
```

## 参数参考

支持以下参数：

* `name_regex` - (可选) 用于通过存储桶名称过滤结果的正则表达式字符串。
* `ids` - (可选) 存储桶 ID 列表。

## 属性参考

除了上述参数外，还导出以下属性：

* `names` - 存储桶名称列表。
* `buckets` - 存储桶列表。每个元素包含以下属性：
  * `name` - 存储桶名称。
  * `acl` - 存储桶访问控制列表。可能的值：`private`、`public-read` 和 `public-read-write`。
  * `extranet_endpoint` - 从外部访问存储桶的互联网域名。
  * `intranet_endpoint` - 从同一区域内的 ECS 实例访问存储桶的内网域名。
  * `location` - 存储桶所在的数据中心区域。
  * `owner` - 存储桶所有者。
  * `storage_class` - 对象存储类型。可能的值：`Standard`、`IA` 和 `Archive`。
  * `creation_date` - 存储桶创建日期。
  * `id` - 存储桶 ID。