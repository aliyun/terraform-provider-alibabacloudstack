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

## 参数说明

支持以下参数：

* `name_regex` - (可选) 用于通过存储桶名称过滤结果的正则表达式字符串。例如，可以通过设置 `"^example"` 来筛选所有以 `example` 开头的存储桶。
* `ids` - (可选) 存储桶 ID 列表。可以通过此参数指定需要查询的具体存储桶 ID。

## 属性说明

除了上述参数外，还导出以下属性：

* `names` - 匹配条件后的存储桶名称列表。
* `buckets` - 匹配条件后的存储桶列表。每个元素包含以下属性：
  * `name` - 存储桶名称。
  * `acl` - 存储桶访问控制列表。可能的值包括：`private`（私有）、`public-read`（公共读）和 `public-read-write`（公共读写）。
  * `extranet_endpoint` - 从外部网络访问存储桶时使用的互联网域名。
  * `intranet_endpoint` - 从同一区域内的 ECS 实例访问存储桶时使用的内网域名。
  * `location` - 存储桶所在的数据中心区域。
  * `owner` - 存储桶的所有者信息。
  * `storage_class` - 对象存储类型。可能的值包括：`Standard`（标准存储）、`IA`（低频访问存储）和 `Archive`（归档存储）。
  * `creation_date` - 存储桶的创建日期，格式为 ISO 8601 标准。
  * `id` - 存储桶的唯一标识符（ID）。