---
subcategory: "对象存储 OSS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_oss_single_tunnel"
sidebar_current: "docs-Alibabacloudstack-oss-single_tunnel"
description: |-
  创建OSS单隧道资源，用于在OSS服务和VPC网络之间建立安全连接通道。
---

# alibabacloudstack_oss_single_tunnel

该资源用于创建OSS单隧道（Single Tunnel），实现OSS服务与指定VPC网络的安全连接。创建后资源不可修改，任何参数变更均需重建资源。

## 示例用法

### 基础用法

```hcl

variable "name" {
  default = "tf-oss-single-tunnel-25672"
}

data "alibabacloudstack_oss_clusters" "default" {
}


data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details              = true
}


resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name   = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
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
  zone_id      = data.alibabacloudstack_zones.default.zones.0.id
  lifecycle {
    ignore_changes = [
      tags
    ]
  }
}
resource "alibabacloudstack_oss_single_tunnel" "default" {
  shared     = "0"
  label      = var.name
  vpc_id     = alibabacloudstack_vpc_vpc.default.id
  vswitch_id = alibabacloudstack_vpc_vswitch.default.id
  cluster    = data.alibabacloudstack_oss_clusters.default.clusters.0.id
}
```

## 参数说明

支持以下参数，按类型排序（必填 → 变更时重建 → 可选 → 过时）：

* `cluster` - (必填, 变更时重建) OSS集群名称。需与目标OSS集群完全匹配，格式为字符串。
* `label` - (必填, 变更时重建) 隧道标签标识。用于资源分类管理，长度1-128字符，不可包含`http://`或`https://`前缀。
* `shared` - (必填, 变更时重建) 共享标识。`0`表示私有隧道（仅本账号可用），`1`表示共享隧道（可跨账号访问），值以字符串形式传递。
* `vpc_id` - (必填, 变更时重建) 目标VPC网络ID。需为当前地域下有效的VPC资源ID，格式如`vpc-xxx`。
* `vswitch_id` - (必填, 变更时重建) 目标虚拟交换机ID。需为`vpc_id`关联的VSwitch资源ID，格式如`vsw-xxx`。

## 属性说明

以下属性导出为资源状态（仅出参属性）：

* `id` - 资源唯一标识，格式为`{cluster}:{vpc_id}:{vip}`。
* `vip` - 分配的VPC内网IP地址。由系统自动分配，用于OSS服务端点访问，格式如`172.16.1.146`。