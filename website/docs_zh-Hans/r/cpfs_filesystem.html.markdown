---
subcategory: "File Storage CPFS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cpfs_filesystem"
sidebar_current: "docs-alibabacloudstack-resource-cpfs-filesystem"
description: |-
  提供CPFS文件系统资源
---

# alibabacloudstack_cpfs_filesystem

提供CPFS文件系统资源。

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

* `storage_type` - (必填，变更时重建) CPFS文件系统的存储类型。有效值：`advance_200` (每TiB 200MB/s 基准I/O带宽)。
* `protocol_type` - (选填，变更时重建) 文件传输协议类型。默认值：`CPFS`。有效值：`CPFS`。
* `description` - (选填) 文件系统描述。长度为2至256个字符。
* `encrypt_type` - (选填，变更时重建) 文件系统是否加密。有效值：
  * `0` (默认)：不加密。
  * `1`：加密。
* `file_system_type` - (选填，变更时重建) 文件系统类型。默认值：`bmcpfs`。有效值：`bmcpfs`。
* `capacity` - (选填) 文件系统的容量。单位：GiB。有效值：20480 至 1740800。
* `zone_id` - (必填，变更时重建) 可用区ID。
* `cluster_id` - (必填，变更时重建) 集群ID。
* `kms_key_id` - (选填) KMS密钥ID。

## 属性说明

导出以下属性：

* `id` - CPFS文件系统的ID。
* `capacity` - 文件系统的容量。
* `zone_id` - 可用区ID。
* `kms_key_id` - KMS密钥ID。

## 导入

可以使用ID导入CPFS文件系统，例如：

```bash
$ terraform import alibabacloudstack_cpfs_filesystem.example 1337849c59
```