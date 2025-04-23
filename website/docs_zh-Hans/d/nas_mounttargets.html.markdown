---
subcategory: "Network Attached Storage (NAS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_nas_mounttargets"
sidebar_current: "docs-Alibabacloudstack-datasource-nas-mounttargets"
description: |- 
  查询文件存储（NAS）挂载点
---

# alibabacloudstack_nas_mounttargets
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_nas_mount_targets`

根据指定过滤条件列出当前凭证权限可以访问的文件存储（NAS）挂载点列表。

## 示例用法

```terraform
variable "name" {
    default = "tf-testAccCheck-nasmount459950"
}

variable "storage_type" {
    default = "Capacity"
}

data "alibabacloudstack_zones" default {
    available_resource_creation = "VSwitch"
    enable_details = true
}

resource "alibabacloudstack_vpc_vpc" "default" {
    vpc_name = "${var.name}_vpc"
    cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
    name = "${var.name}_vsw"
    vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
    cidr_block = "172.16.0.0/24"
    zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

data "alibabacloudstack_nas_protocols" "default" {
    type = "${var.storage_type}"
}

resource "alibabacloudstack_nas_file_system" "default" {
    description = "${var.name}"
    storage_type = "${var.storage_type}"
    protocol_type = "${data.alibabacloudstack_nas_protocols.default.protocols.0}"
}

resource "alibabacloudstack_nas_access_group" "default" {
    access_group_name = "${var.name}"
    access_group_type = "Vpc"
    description = "${var.name}"
}

resource "alibabacloudstack_nas_mount_target" "default" {
    file_system_id = "${alibabacloudstack_nas_file_system.default.id}"
    access_group_name = "${alibabacloudstack_nas_access_group.default.access_group_name}"
    vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
}

data "alibabacloudstack_nas_mount_targets" "default" {
    file_system_id = "${alibabacloudstack_nas_mount_target.default.file_system_id}"
    access_group_name = "${alibabacloudstack_nas_mount_target.default.access_group_name}"
    network_type = "VPC"
    vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
    vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
    status = "Active"
}

output "first_mount_target_domain" {
    value = data.alibabacloudstack_nas_mount_targets.default.targets.0.mount_target_domain
}
```

## 参数参考

以下参数是支持的：

* `access_group_name` - (选填, 变更时重建) 权限组名称。
* `mount_target_domain` - (选填, 强制新建, 已在 1.53.0+ 版本中弃用) 挂载点域名。此字段自提供商版本 1.53.0 起已弃用，请改用 `ids`。
* `type` - (选填, 强制新建, 已在 1.95.0+ 版本中弃用) 字段 `type` 自提供商版本 1.95.0 起已弃用，请改用 `network_type`。
* `network_type` - (选填, 强制新建, 自 1.95.0+ 版本起可用) 网络类型。有效值包括 `VPC` 等。
* `vpc_id` - (选填, 变更时重建) VPC ID。
* `vswitch_id` - (选填, 变更时重建) 交换机 ID。
* `file_system_id` - (必填, 变更时重建) 文件系统 ID。
* `ids` - (选填, 强制新建, 自 1.53.0+ 版本起可用) 挂载目标域名 ID 列表。
* `status` - (选填, 强制新建, 自 1.95.0+ 版本起可用) 挂载点当前状态，包括：`Active`、`Inactive` 和 `Pending`。当状态为 `Active` 时才可以进行文件系统挂载使用。

## 属性参考

除了上述参数外，还导出以下属性：

* `targets` - NAS 挂载目标列表。每个元素包含以下属性：
  * `access_group_name` - 权限组名称。
  * `id` - 挂载目标的唯一标识符。
  * `mount_target_domain` - 挂载点域名。
  * `network_type` - 网络类型。
  * `type` - (已在 1.95.0+ 版本中弃用) 此字段自提供商版本 1.95.0 起已弃用，请改用 `network_type`。
  * `status` - 挂载点当前状态，包括：`Active`、`Inactive` 和 `Pending`。当状态为 `Active` 时才可以进行文件系统挂载使用。
  * `vpc_id` - VPC ID。
  * `vswitch_id` - 交换机 ID。