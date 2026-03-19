---
subcategory: "Network Attached Storage (NAS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_nas_namespace_mount_target"
sidebar_current: "docs-Alibabacloudstack-resource-nas-namespace-mount-target"
description: |-
  创建和管理NAS统一命名空间挂载点
---

# alibabacloudstack_nas_namespace_mount_target

> NAS统一命名空间挂载设置

## 示例用法

### 基础用法

```hcl

variable "name" {
  default = "tf-AccNasaccgop86369"
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

resource "alibabacloudstack_nas_accessgroup" "default1" {
  access_group_name = "${var.name}1"
  access_group_type = "Vpc"
}

resource "alibabacloudstack_nas_accessgroup" "default2" {
  access_group_name = "${var.name}2"
  access_group_type = "Vpc"
}



resource "alibabacloudstack_nas_namespace_mount_target" "default" {
  network_type      = "Vpc"
  access_group_name = alibabacloudstack_nas_accessgroup.default1.access_group_name
  nas_namespace_id  = alibabacloudstack_nas_namespace.default.id
  vswitch_id        = alibabacloudstack_vpc_vswitch.default.id
}
```

## 参数说明

支持以下参数：

* `nas_namespace_id` - (必填, 变更时重建) 命名空间ID。
* `network_type` - (必填, 变更时重建) 挂载点网络类型。取值：
  * `Vpc`：专有网络
  * `Classic`：经典网络
* `access_group_name` - (必填) 权限组名称。限制：
  * 长度为3~64个字符
  * 必须以大小写字母开头，可以包含英文字母、数字、下划线（_）或者短划线（-）
  * 新创建的权限组名称不能与两个默认的权限组相同（DEFAULT_VPC_GROUP_NAME和DEFAULT_CLASSIC_GROUP_NAME）
* `vswitch_id` - (可选, 变更时重建) 交换机ID。当网络类型是专有网络时，此字段必填且有意义。
* `status` - (可选) 挂载点状态。取值：
  * `Active`：可用
  * `Inactive`：不可用

## 属性说明

以下属性会从API响应中导出：

* `id` - 资源ID，格式为{NasNamespaceId:MountTargetDomain}。
* `mount_target_domain` - 挂载点域名。
* `vpc_id` - 专有网络ID。当网络类型是专有网络时，此字段有值。