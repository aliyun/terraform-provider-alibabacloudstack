---
subcategory: "Network Attached Storage (NAS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_nas_mounttarget"
sidebar_current: "docs-Alibabacloudstack-nas-mounttarget"
description: |- 
  编排文件存储（NAS）挂载点
---

# alibabacloudstack_nas_mounttarget
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_nas_mount_target`

使用Provider配置的凭证在指定的资源集编排文件存储（NAS）挂载点。

## 示例用法

```terraform
variable "name" {
    default = "tf-testaccnasmount_target81413"
}

resource "alibabacloudstack_vpc" "default" {
    cidr_block = "172.16.0.0/16"
    name       = "${var.name}"
}

data "alibabacloudstack_zones" "default" {
    available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vswitch" "default" {
    vpc_id            = "${alibabacloudstack_vpc.default.id}"
    cidr_block        = "172.16.0.0/21"
    availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
    name             = "${var.name}"
}

variable "storage_type" {
    default = "Capacity"
}

data "alibabacloudstack_nas_protocols" "default" {
    type = "${var.storage_type}"
}

resource "alibabacloudstack_nas_file_system" "default" {
    description  = "${var.name}"
    storage_type = "${var.storage_type}"
    protocol_type = "${data.alibabacloudstack_nas_protocols.default.protocols.0}"
}

resource "alibabacloudstack_nas_access_group" "default" {
    access_group_name = "tf-testAccNasConfig-resource-test86663"
    access_group_type = "Vpc"
    description       = "tf-testAccNasConfig"
}

resource "alibabacloudstack_nas_mounttarget" "default" {
    file_system_id    = "${alibabacloudstack_nas_file_system.default.id}"
    vswitch_id        = "${alibabacloudstack_vswitch.default.id}"
    access_group_name = "${alibabacloudstack_nas_access_group.default.access_group_name}"
    status            = "Active"
    security_group_id = "sg-xxxxxxxxx" # Example Security Group ID
}
```

## 参数说明

支持以下参数：

* `access_group_name` - (必填) 应用于此挂载目标的权限组名称。权限组定义了访问控制规则，用于限制客户端对文件系统的访问。
* `file_system_id` - (必填，强制更新) 文件系统的ID。这是要挂载的目标文件系统。
* `vswitch_id` - (可选，强制更新) 挂载目标所在的VPC的交换机ID。如果未指定，则默认使用经典网络。注意：当前仅支持在中国大陆区域的经典网络中创建挂载点。
* `status` - (可选) 挂载目标的状态。有效值为：`Active` 和 `Inactive`。默认值为`Active`。在挂载文件系统之前，请确保挂载目标处于活动状态。
* `security_group_id` - (可选, 强制更新, v1.95.0+可用) 安全组ID。用于限制访问挂载目标的流量。通过设置安全组，可以进一步增强文件系统的安全性。

## 属性说明

除了上述所有参数外，还导出了以下属性：

* `id` - 此资源的ID。它格式化为 `<file_system_id>:<mount_target_domain>`。在版本1.95.0之前，该值为 `<mount_target_domain>`。
* `status` - 挂载点的当前状态，包括 `Active` 和 `Inactive`。只有当状态为 `Active` 时，才能使用文件系统。此属性为只读，表示挂载点的实际状态。