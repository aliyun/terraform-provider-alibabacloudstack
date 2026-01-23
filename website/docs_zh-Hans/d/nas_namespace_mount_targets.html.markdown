---
subcategory: "Network Attached Storage (NAS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_nas_namespace_mount_target"
sidebar_current: "docs-Alibabacloudstack-resource-nas-namespace-mount-target"
description: |-
  管理NAS统一命名空间挂载目标，用于配置文件系统访问权限和网络连接。
---

# alibabacloudstack_nas_namespace_mount_target

管理NAS统一命名空间挂载目标，用于配置文件系统访问权限和网络连接。挂载目标是访问NAS文件系统的入口点，需指定VPC网络、交换机及权限组等参数。

## 示例用法

```hcl

variable "name" {
  default = "tf-testnasfs9865"
}

data "alibabacloudstack_nas_zones" "default" {
}


data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details              = true
}


resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name   = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
  tags = {
    common_test = "terraform"
    filter      = var.name
  }
  enable_ipv6 = true
  lifecycle {
    ignore_changes = [
      secondary_cidr_blocks,
      tags
    ]
  }
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  vswitch_name = "${var.name}_vsw"
  vpc_id       = alibabacloudstack_vpc_vpc.default.id
  cidr_block   = "172.16.1.0/24"
  enable_ipv6  = true
  zone_id      = data.alibabacloudstack_zones.default.zones.0.id
  lifecycle {
    ignore_changes = [
      tags
    ]
  }
}



resource "alibabacloudstack_nas_namespace" "default" {
  zone_id       = data.alibabacloudstack_nas_zones.default.zones.0.zone_id
  cluster_id    = data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.cluster_id
  description   = var.name
  storage_type  = data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.instance_types.0.storage_type
  protocol_type = data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.instance_types.0.protocol_type
  encrypt_type  = "0"
}

resource "alibabacloudstack_nas_accessgroup" "default" {
  access_group_name = var.name
  access_group_type = "Vpc"
}

resource "alibabacloudstack_nas_namespace_mount_target" "default" {
  nas_namespace_id  = alibabacloudstack_nas_namespace.default.id
  access_group_name = alibabacloudstack_nas_accessgroup.default.access_group_name
  vswitch_id        = alibabacloudstack_vpc_vswitch.default.id
  network_type      = "Vpc"
}



data "alibabacloudstack_nas_namespace_mount_targets" "default" {
  nas_namespace_id = alibabacloudstack_nas_namespace_mount_target.default.nas_namespace_id
  ids = [
    "${alibabacloudstack_nas_namespace_mount_target.default.id}"
  ]
}
```

## 参数说明

以下参数支持：

* `access_group_name` (字符串) - (必填) 权限组名称。用于控制客户端访问权限的权限组标识，需提前创建。

* `nas_namespace_id` (字符串) - (必填, 变更时重建) 命名空间ID。挂载目标关联的统一命名空间ID，格式为`1NSxxx`。

* `network_type` (字符串) - (必填, 变更时重建) 网络类型。当前仅支持`Vpc`，表示挂载目标部署在VPC网络中。

* `vpc_id` (字符串) - (必填, 变更时重建) VPC ID。挂载目标所属的专有网络ID，需与命名空间在同一地域。

* `vswitch_id` (字符串) - (必填, 变更时重建) 交换机ID。挂载目标部署的虚拟交换机ID，需与VPC匹配且具有足够IP资源。

* `status` (字符串) - (可选) 挂载目标状态。可选值为`Active`（启用）或`Inactive`（停用），默认为`Active`。状态变更会立即生效。

## 属性说明

以下属性被导出：

* `id` (字符串)：资源ID，格式为`{nas_namespace_id}:{mount_target_domain}`，由系统生成。

* `mount_target_domain` (字符串)：挂载目标域名。系统自动生成的访问域名，用于客户端挂载文件系统，格式为`ns-{nas_namespace_id}-xxx.region.nas.inter.envxxx.shuguang.com`。