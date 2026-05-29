---
subcategory: "并行文件存储 CPFS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cpfs_filesystem"
sidebar_current: "docs-alibabacloudstack-resource-cpfs-filesystem"
description: |-
  提供CPFS文件系统资源
---

# alibabacloudstack_cpfs_filesystem

提供CPFS文件系统资源。

> **注意:** 该资源也可以使用以下别名引用：
> - `alibabacloudstack_cpfs_file_system`

## 示例用法

基础用法

```terraform
variable "name" {
	default = "tfacctest124"
}

data "alibabacloudstack_nas_zones" "default" {
	file_system_type = "bmcpfs"
}

resource "alibabacloudstack_cpfs_filesystem" "example" {
	storage_type = "${data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.instance_types.0.storage_type}"
	zone_id = "${data.alibabacloudstack_nas_zones.default.zones.0.zone_id}"
	cluster_id = "${data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.cluster_id}"
	description = "${var.name}"
	capacity = 20480
}
```

## 参数说明

支持以下参数：

* `storage_type` - (必填，变更时重建) CPFS文件系统的存储类型。有效值：`advance_100` (每TiB 100MB/s 基准I/O带宽)、`advance_200` (每TiB 200MB/s 基准I/O带宽)、`economic` (经济型)。
* `protocol_type` - (选填，变更时重建) 文件传输协议类型。默认值：`CPFS`。有效值：`CPFS`。
* `description` - (选填) 文件系统描述。长度为2至256个字符。
* `encrypt_type` - (选填，变更时重建，已弃用) **该字段已弃用，将在 3.21.0 版本中移除。** 文件系统是否加密。有效值：`0`（默认，不加密）、`1`（加密）、`2`（使用KMS密钥加密，需配合 `kms_key_id`）。
* `file_system_type` - (选填，变更时重建) 文件系统类型。默认值：`bmcpfs`。有效值：`bmcpfs`。
* `capacity` - (选填) 文件系统的容量。单位：GiB。有效值：20480 至 1740800。当 `file_system_type` 为 `bmcpfs` 时必填。
* `zone_id` - (必填，变更时重建) 可用区ID。
* `cluster_id` - (必填，变更时重建) 集群ID。该参数为专有云特有，对应API返回的 `Location` 字段。
* `kms_key_id` - (选填，已弃用) **该字段已弃用，将在 3.21.0 版本中移除。** KMS密钥ID。当 `encrypt_type` 为 `2` 时必填。

## 属性说明

导出以下属性：

* `id` - CPFS文件系统的ID。
* `capacity` - (可回读) 文件系统的容量。
* `zone_id` - (可回读) 可用区ID。
* `kms_key_id` - (可回读，已弃用) KMS密钥ID。
* `storage_type` - (可回读) 存储类型。
* `protocol_type` - (可回读) 协议类型。
* `file_system_type` - (可回读) 文件系统类型。
* `encrypt_type` - (可回读，已弃用) 加密类型。
* `description` - (可回读) 文件系统描述。

## 导入

可以使用ID导入CPFS文件系统，例如：

```bash
$ terraform import alibabacloudstack_cpfs_filesystem.example 1337849c59
```