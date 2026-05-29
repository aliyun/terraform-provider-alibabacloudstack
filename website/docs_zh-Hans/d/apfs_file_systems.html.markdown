---
subcategory: "阿里云并行文件系统(APFS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_apfs_file_systems"
sidebar_current: "docs-Alibabacloudstack-datasource-apfs-file-systems"
description: |-
  查询APFS文件系统信息
---

# alibabacloudstack_apfs_file_systems

查询APFS（Alibaba Parallel File System）文件系统信息的数据源，用于检索已存在的APFS文件系统资源。

## 示例用法

```hcl

variable "name" {
  default = "tf-testacc486230"
}

data "alibabacloudstack_zones" "default" {
  enable_details = true
}

resource "alibabacloudstack_apfs_file_system" "default" {
  zone_id      = data.alibabacloudstack_zones.default.zones.0.id
  cluster_id   = "EfsStorageCluster-A-20251125-0131"
  storage_type = "parastor_advance_300"
  volume_size  = 1024
  description  = var.name
}

data "alibabacloudstack_apfs_file_systems" "default" {
  name_regex = alibabacloudstack_apfs_file_system.default.description
}
```

## 参数说明

以下参数支持过滤查询结果：

* `file_system_type` (可选)：文件系统的类型。默认值为"efs"。
* `ids` (可选)：文件系统ID列表，用于按ID过滤结果。
* `name_regex` (可选)：用于按文件系统描述过滤结果的正则表达式。
* `status` (可选)：文件系统的状态，例如Pending、Running等。

## 属性说明

以下属性被导出：

* `id` (字符串)：数据源ID，由匹配的文件系统ID哈希生成。
* `bandwidth` (整数)：文件系统的带宽限制。
* `capacity` (整数)：文件系统的总容量（GB）。
* `charge_type` (字符串)：文件系统的计费方式。
* `create_time` (字符串)：文件系统的创建时间。
* `description` (字符串)：文件系统的描述信息。
* `encrypt_type` (整数)：文件系统的加密类型。
* `file_system_id` (字符串)：文件系统的ID。
* `file_systems` (列表)：匹配的文件系统列表。每个元素包含以下属性：
  * `bandwidth` (整数)：文件系统的带宽限制。
  * `capacity` (整数)：文件系统的总容量（GB）。
  * `charge_type` (字符串)：文件系统的计费方式。
  * `create_time` (字符串)：文件系统的创建时间。
  * `description` (字符串)：文件系统的描述信息。
  * `encrypt_type` (整数)：文件系统的加密类型。
  * `file_system_id` (字符串)：文件系统的ID。
  * `id` (字符串)：文件系统的ID（与file_system_id相同）。
  * `location` (字符串)：文件系统所在的集群位置。
  * `metered_size` (整数)：文件系统的计量使用量（GB）。
  * `protocol_type` (字符串)：文件系统使用的协议类型。
  * `quota_size` (整数)：分配给文件系统的配额大小（GB）。
  * `region_id` (字符串)：文件系统所在的区域ID。
  * `status` (字符串)：文件系统的当前状态。
  * `storage_type` (字符串)：文件系统的存储类型。
  * `volume_size` (整数)：文件系统的卷大小（GB）。
  * `zone_id` (字符串)：文件系统所在的可用区ID。
* `location` (字符串)：文件系统所在的集群位置。
* `metered_size` (整数)：文件系统的计量使用量（GB）。
* `protocol_type` (字符串)：文件系统使用的协议类型。
* `quota_size` (整数)：分配给文件系统的配额大小（GB）。
* `region_id` (字符串)：文件系统所在的区域ID。
* `status` (字符串)：文件系统的当前状态。
* `storage_type` (字符串)：文件系统的存储类型。
* `volume_size` (整数)：文件系统的卷大小（GB）。
* `zone_id` (字符串)：文件系统所在的可用区ID。