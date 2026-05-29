---
subcategory: "文件存储 NAS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_nas_namespace_group"
sidebar_current: "docs-Alibabacloudstack-nas-namespace-group"
description: |-
  编排NAS跨域挂载编排
---

# alibabacloudstack_nas_namespace_group

使用Provider配置的凭证在指定的资源集编排NAS跨域挂载编排资源。

## 示例用法

### 基础用法

```hcl

variable "name" {
  default = "tf-AccNasaccgop35376"
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
  network_type      = "Vpc"
  access_group_name = alibabacloudstack_nas_accessgroup.default.access_group_name
  nas_namespace_id  = alibabacloudstack_nas_namespace.default.id
  vswitch_id        = alibabacloudstack_vpc_vswitch.default.id
}



resource "alibabacloudstack_nas_namespace_group" "default" {
  network_type        = "Vpc"
  mapped_path         = var.name
  nas_namespace_id    = alibabacloudstack_nas_namespace.default.id
  mount_target_domain = alibabacloudstack_nas_namespace_mount_target.default.mount_target_domain
}
```

## 参数说明

支持以下参数：

* `mapped_path` - (必填, 变更时重建) 跨域挂载编排映射路径。一旦映射，路径将永久绑定到此挂载点。
* `mount_target_domain` - (必填, 变更时重建) 命名空间挂载点域名。
* `nas_namespace_id` - (必填, 变更时重建) 命名空间ID。
* `network_type` - (必填, 变更时重建) 网络类型。取值：`Vpc`（专有网络）或`Classic`（经典网络）。一旦设置，无法更改，除非清除列表。

## 属性说明

以下属性导出为资源属性：

* `id` - 资源ID，值为挂载点域名（`mount_target_domain`）。
* `create_time` - 资源创建时间，格式为ISO 8601标准时间。
* `member_id` - 跨域挂载编排成员ID，用于删除操作。
* `status` - 资源状态，表示当前挂载编排的启用状态。