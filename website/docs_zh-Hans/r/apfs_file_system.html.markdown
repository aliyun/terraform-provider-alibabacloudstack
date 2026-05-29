---
subcategory: "阿里云并行文件系统(APFS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_apfs_file_system"
sidebar_current: "docs-Alibabacloudstack-resource-apfs-file-system"
description: |-
  管理APFS文件存储系统
---

# alibabacloudstack_apfs_file_system

使用Provider配置的凭证在指定的资源集管理APFS（Alibaba Parallel File System）文件存储系统。

-> **注意：** 该资源也可以使用以下别名引用：`apsarastack_apfs_file_system`。

## 示例用法

### 基础用法

```hcl
variable "name" {
  default = "tf-testacc47734"
}

data "alibabacloudstack_apfs_zones" "default" {
}

resource "alibabacloudstack_apfs_file_system" "default" {
  zone_id      = data.alibabacloudstack_apfs_zones.default.zones.0.zone_id
  cluster_id   = data.alibabacloudstack_apfs_zones.default.zones.0.clusters.0.cluster_id
  storage_type = data.alibabacloudstack_apfs_zones.default.zones.0.clusters.0.storage_type
  volume_size  = 256
  description  = var.name
}
```

## 参数说明

支持以下参数：

* `storage_type` - (必填, 变更时重建) 存储类型。必须指定有效的存储类型，例如：`parastor_advance_300`。
* `volume_size` - (必填) 文件系统容量，单位GB。
* `cluster_id` - (可选, 变更时重建) 集群ID。指定文件系统所属的集群。
* `file_system_type` - (可选, 变更时重建) 文件系统类型。默认值为`efs`，表示标准文件系统。
* `protocol_type` - (可选, 变更时重建) 协议类型。默认值为`EFS`，表示使用EFS协议。
* `zone_id` - (可选, 变更时重建) 可用区ID。指定文件系统创建的可用区。
* `description` - (可选) 文件系统的描述信息。长度限制为0-256个字符。

## 属性说明

以下属性导出：

* `id` - 文件系统的ID，格式如`efs-0187aa76`。
* `file_system_id` - 文件系统的ID（同`id`属性）。
* `create_time` - 文件系统的创建时间，格式为ISO8601（例如：`2025-12-11T07:12:37Z`）。
* `status` - 文件系统的当前状态。可能值包括`Pending`（创建中）、`Running`（运行中）、`Deleting`（删除中）等。
* `metered_size` - 已计量的存储大小，单位GB。表示当前实际使用的存储容量。
* `quota_size` - 文件系统的配额大小，单位GB。表示系统分配的最大存储配额。

## Import

APFS文件系统可以使用FileSystemId导入，例如：

```
$ terraform import alibabacloudstack_apfs_file_system.example efs-xxxxxxxx
```