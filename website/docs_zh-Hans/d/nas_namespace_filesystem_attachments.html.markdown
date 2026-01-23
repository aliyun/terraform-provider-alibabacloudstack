---
subcategory: "Network Attached Storage (NAS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_nas_namespace_filesystem_attachment"
sidebar_current: "docs-Alibabacloudstack-datasource-nas-namespace-filesystem-attachment"
description: |-
  查询阿里云NAS统一命名空间文件存储映射信息
---
# alibabacloudstack_nas_namespace_filesystem_attachment

> NAS统一命名空间文件存储映射数据源，用于查询命名空间中已映射的文件系统信息

## 示例用法

```hcl

variable "name" {
  default = "tf-testnasfs1942"
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


resource "alibabacloudstack_nas_namespace" "default" {
  zone_id       = data.alibabacloudstack_nas_zones.default.zones.0.zone_id
  cluster_id    = data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.cluster_id
  description   = var.name
  storage_type  = data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.instance_types.0.storage_type
  protocol_type = data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.instance_types.0.protocol_type
  encrypt_type  = "0"
}

resource "alibabacloudstack_nas_namespace_filesystem_attachment" "default" {
  nas_namespace_id = alibabacloudstack_nas_namespace.default.id
  file_system_id   = alibabacloudstack_nas_file_system.default.id
  mapped_path      = var.name
}



data "alibabacloudstack_nas_namespace_filesystem_attachments" "default" {
  nas_namespace_id = alibabacloudstack_nas_namespace_filesystem_attachment.default.nas_namespace_id
  ids = [
    "${alibabacloudstack_nas_namespace_filesystem_attachment.default.id}"
  ]
}
```

## 参数说明
以下参数支持：

* `nas_namespace_id` (字符串, 必填)：命名空间的ID，用于指定要查询的命名空间。

* `file_system_id` (字符串, 可选)：文件系统的ID，用于过滤特定文件系统的映射信息。

* `ids` (列表, 可选)：附件ID列表，格式为NasNamespaceId:FileSystemId，用于精确匹配特定附件。

* `mapped_path` (字符串, 可选)：映射路径，用于过滤特定映射路径的文件系统。

* `name_regex` (字符串, 可选)：正则表达式，用于通过映射路径过滤结果。

## 属性说明
以下属性被导出：

* `id` (字符串)：数据源ID，由附件ID的哈希值生成。

* `attachment_id` (字符串)：附件的ID，格式为NasNamespaceId:FileSystemId。

* `create_time` (字符串)：附件的创建时间，格式为ISO 8601标准时间。

* `file_system_type` (字符串)：文件系统的类型，例如standard表示通用型NAS。

* `storage_type` (字符串)：文件系统的存储类型，可能值包括Capacity（容量型）和Performance（性能型）。

* `id` (字符串, 已弃用)：字段'id'已从provider版本1.200.0弃用，请使用'attachment_id'代替。