---
subcategory: "文件存储 NAS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_nas_namespaces"
sidebar_current: "docs-Alibabacloudstack-datasource-nas-namespaces"
description: |-
  查询阿里云NAS命名空间信息
---

# alibabacloudstack_nas_namespaces

NAS统一命名空间列表数据源，用于查询阿里云NAS命名空间信息。

## 示例用法

```hcl

variable "name" {
  default = "tf-testnasfs3606"
}

data "alibabacloudstack_nas_zones" "default" {
}

resource "alibabacloudstack_nas_namespace" "default" {
  zone_id       = data.alibabacloudstack_nas_zones.default.zones.0.zone_id
  cluster_id    = data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.cluster_id
  description   = var.name
  storage_type  = data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.instance_types.0.storage_type
  protocol_type = data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.instance_types.0.protocol_type
  encrypt_type  = "0"
}



data "alibabacloudstack_nas_namespaces" "default" {
  ids = [
    "${alibabacloudstack_nas_namespace.default.id}"
  ]
}
```

## 参数说明

以下参数支持作为过滤条件：aaaaaaa

`- file_system_type` (字符串, 可选)：文件系统类型。默认值：`standard`。

- `ids` (列表, 可选)：命名空间ID列表，用于过滤结果。

- `name_regex` (字符串, 可选)：正则表达式，用于通过命名空间描述过滤结果。

- `protocol_type` (字符串, 可选)：协议类型。可能的值：`NFS`，`SMB`。

- `storage_type` (字符串, 可选)：存储类型。可能的值：`Performance`，`Capacity`。

- `zone_id` (字符串, 可选)：命名空间所属的可用区ID。

## 属性说明

以下属性被导出：

- `ids` (列表)：匹配的命名空间ID列表。

- `namespaces` (列表)：匹配的命名空间列表。每个元素包含以下属性：

  - `id` (字符串)：命名空间ID，等同于`nas_namespace_id`。

  - `create_time` (字符串)：命名空间创建时间，遵循ISO 8601标准，格式为`yyyy-MM-ddTHH:mm:ssZ`。

  - `description` (字符串)：命名空间描述。

  - `encrypt_type` (整数)：命名空间加密类型。取值：`0`（不加密），`1`（加密）。

  - `file_system_type` (字符串)：文件系统类型。默认值：`standard`，通用型NAS。

  - `mount_target_count` (整数)：挂载点数量。

  - `nas_namespace_id` (字符串)：命名空间ID。

  - `protocol_type` (字符串)：协议类型。可能的值：`NFS`（NFS文件传输协议），`SMB`（SMB文件传输协议）。

  - `storage_type` (字符串)：存储类型。可能的值：`Performance`（性能型），`Capacity`（容量型）。

  - `status` (字符串)：命名空间状态。可能的值：`Initializing`（初始化中），`Normal`（正常）。

  - `zone_id` (字符串)：可用区ID。