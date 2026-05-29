---
subcategory: "文件存储 NAS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_nas_namespace_filesystem_attachment"
sidebar_current: "docs-Alibabacloudstack-resource-nas-namespace-filesystem-attachment"
description: |-
  编排NAS统一命名空间文件存储映射
---

# alibabacloudstack_nas_namespace_filesystem_attachment

创建、管理和配置NAS统一命名空间文件存储映射。

## 示例用法

### 基础用法

```hcl

variable "name" {
  default = "tf-testAccNasnpsFsAttachment26768"
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
  mapped_path      = var.name
  file_system_id   = alibabacloudstack_nas_file_system.default.id
}
```

## 参数说明

支持以下参数：

* `file_system_id` - (必填, 变更时重建) 文件系统ID。需要映射到命名空间的文件系统标识符。
* `nas_namespace_id` - (必填, 变更时重建) 命名空间ID。统一命名空间的标识符。
* `mapped_path` - (必填) 映射路径。文件系统在命名空间中的挂载路径。

## 属性说明

以下属性会从API响应中导出：

* `id` - 资源ID，格式为`{NasNamespaceId:FileSystemId}`。
* `create_time` - 创建时间。文件系统映射到命名空间的时间，格式为ISO 8601标准时间。
* `file_system_type` - 文件系统类型。可能的值包括：
  * `standard`：通用型NAS
* `storage_type` - 存储类型。可能的值包括：
  * `Capacity`：容量型
  * `Performance`：性能型