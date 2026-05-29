---
subcategory: "文件存储 NAS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_nas_dir_quota"
sidebar_current: "docs-Alibabacloudstack-datasource-nas-dir-quota"
description: |-
  查询NAS目录配额信息
---

# alibabacloudstack_nas_dir_quota

查询阿里云NAS文件系统的目录配额信息。该数据源用于检索指定文件系统中已配置配额的目录信息，包括目录路径、inode号以及各用户的配额限制和实际使用情况。

## 示例用法

```hcl

variable "name" {
  default = "tf-testnasdirquotas2499023223732943598"
}


data "alibabacloudstack_nas_zones" "default" {
}

resource "alibabacloudstack_nas_file_system" "default" {
  protocol_type = data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.instance_types.0.protocol_type
  storage_type  = data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.instance_types.0.storage_type
  encrypt_type  = "0"
  zone_id       = data.alibabacloudstack_nas_zones.default.zones.0.zone_id
  cluster_id    = data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.cluster_id
  description   = var.name
}


resource "alibabacloudstack_nas_dir_quota" "default" {
  file_system_id = alibabacloudstack_nas_file_system.default.id
  path           = "/"
  quotas {
    quota_type       = "Enforcement"
    user_type        = "Uid"
    user_id          = "500"
    size_limit       = 100
    file_count_limit = 10000
  }
}


data "alibabacloudstack_nas_dir_quotas" "default" {
  file_system_id = alibabacloudstack_nas_dir_quota.default.file_system_id
}
```

## 参数说明
以下参数用于过滤查询结果：

* `file_system_id` (必填)：文件系统ID，用于指定要查询的NAS文件系统。
* `path` (可选)：目录在文件系统中的绝对路径。不填写时会返回文件系统中全部已设置了配额的目录。
* `ids` (可选)：根据ID列表过滤结果。ID格式为{FileSystemId:Path}。
* `name_regex` (可选)：根据正则表达式过滤结果，匹配的是目录路径。

## 属性说明
以下属性被导出：

* `id` (字符串)：资源ID，格式为{FileSystemId:Path}。

以下属性在`dir_quotas`列表中导出：

* `dir_inode` (字符串)：目录的inode号。
* `file_system_id` (字符串)：文件系统ID。
* `path` (字符串)：目录在文件系统中的绝对路径。
* `user_quotas` (列表)：用户配额信息列表。

以下属性在`user_quotas`列表中导出：

* `file_count_limit` (整数)：目录下用户的文件数目限制。-1表示无限制。
* `file_count_real` (整数)：目录下用户的实际文件数目。
* `quota_type` (字符串)：配额类型，包括统计型（Accounting）和限制型（Enforcement）。
* `size_limit` (整数)：目录下用户的文件总容量限制，单位为GB。-1表示无限制。
* `size_real` (整数)：目录下用户的实际文件总容量，单位为GB。
* `user_id` (字符串)：要限制的uid或gid，取决于UserType的值。当UserType为AllUsers时，此值为空。
* `user_type` (字符串)：指定UserId的类型，包括Uid、Gid、AllUsers三种类型。