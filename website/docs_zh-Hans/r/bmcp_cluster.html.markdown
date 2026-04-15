---
subcategory: "裸金属计算平台 (BMCP)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_bmcp_cluster"
sidebar_current: "docs-Alibabacloudstack-bmcp-cluster"
description: |-
  提供 BMCP（裸金属计算平台）集群资源。
---

# alibabacloudstack_bmcp_cluster

提供 BMCP（裸金属计算平台）集群资源。

## 使用示例

### 基础用法

```hcl
variable "name" {
  default = "tf-testacc-bmcp-cluster"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vpc" "default" {
  name       = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alibabacloudstack_vswitch" "default" {
  name              = var.name
  vpc_id            = alibabacloudstack_vpc.default.id
  cidr_block        = "192.168.40.0/24"
  is_cgw            = true
  availability_zone = data.alibabacloudstack_zones.default.zones.0.id
}

resource "alibabacloudstack_vswitch" "standard" {
  name              = "${var.name}-std"
  vpc_id            = alibabacloudstack_vpc.default.id
  cidr_block        = "192.168.50.0/24"
  availability_zone = data.alibabacloudstack_zones.default.zones.0.id
}

resource "alibabacloudstack_evpc_evpc" "default" {
  evpc_name   = var.name
  description = var.name
}

data "alibabacloudstack_bmcp_machinetypes" "all" {
  min_standard_instance_count = 1
}

resource "alibabacloudstack_bmcp_cluster" "default" {
  cluster_name        = var.name
  vpc_id              = alibabacloudstack_vpc.default.id
  evpc_id             = alibabacloudstack_evpc_evpc.default.id
  zone_id             = data.alibabacloudstack_zones.default.zones.0.id
  standard_vswitch_id = alibabacloudstack_vswitch.standard.id
  password            = "Test1234!"
  machine_type        = data.alibabacloudstack_bmcp_machinetypes.all.machinetypes.0.name
  node_count          = 1
  vswitch_id          = alibabacloudstack_vswitch.default.id
}
```

## 参数参考

支持以下参数：

* `cluster_name` - (必填，强制新建) BMCP 集群的名称。
* `vpc_id` - (必填，强制新建) 集群所在 VPC 的 ID。
* `evpc_id` - (必填，强制新建) 与集群关联的企业 VPC (EVPC) 的 ID。
* `zone_id` - (必填，强制新建) 集群所在可用区的 ID。
* `standard_vswitch_id` - (必填，强制新建) 集群标准交换机的 ID。
* `machine_type` - (必填，强制新建) 集群节点的机器类型。可使用 `alibabacloudstack_bmcp_machinetypes` 数据源查询可用的机器类型。
* `node_count` - (必填) 集群中的节点数量。
* `vswitch_id` - (必填，强制新建) 集群节点所用交换机的 ID。
* `password` - (可选，强制新建，敏感) 集群节点的密码。必须提供 `password` 或 `activation_code` 之一。
* `activation_code` - (可选，强制新建，敏感) 集群的激活码。必须提供 `password` 或 `activation_code` 之一。
* `switch_method` - (可选，强制新建) 集群的交换方法。有效值：`CHSW` (默认)，`ROCE`。
* `cluster_arch_type` - (可选，强制新建) 集群的架构类型。有效值：`standard` (默认)，`high_performance`。
* `is_install_yundun_aegis` - (可选，强制新建) 是否安装云盾安骑士安全代理。默认值：`true`。
* `is_create_cpfs_cluster` - (可选，强制新建) 是否创建 CPFS 集群。默认值：`false`。
* `enable_ipv6` - (可选，强制新建) 是否启用 IPv6。默认值：`false`。

## 属性参考

除上述所有参数外，还导出以下属性：

* `id` - 集群的 ID (与 `cluster_id` 相同)。
* `cluster_id` - 集群的唯一标识符。
* `status` - 集群的状态 (例如：`active`、`creating`、`deleting`)。
* `region_id` - 集群所在区域的 ID。
* `cpu_count` - 集群中的 CPU 总数。
* `mem_count` - 集群的总内存大小 (GB)。
* `flops_count` - 集群的总浮点运算次数 (FLOPS)。
* `video_memory` - 集群的总显存 (GB)。
* `create_time` - 集群的创建时间。
* `update_time` - 集群的最后更新时间。
* `gpu` - 集群中的 GPU 信息列表。每个元素包含：
  * `gpu_num` - GPU 数量。
  * `gpu_model` - GPU 型号。
* `machine_type_list` - 机器类型详细信息列表。每个元素包含：
  * `machine_type` - 机器类型名称。
  * `node_count` - 该机器类型的节点数量。
  * `vswitch_id` - 该机器类型的交换机 ID。
  * `arch` - CPU 架构。
  * `gpu` - GPU 型号。
  * `cpu_count` - 每个节点的 CPU 数量。
  * `mem_count` - 每个节点的内存大小 (GB)。
  * `gpu_num` - 每个节点的 GPU 数量。
  * `video_mem_count` - 每个节点的显存 (GB)。
  * `flops_count` - 每个节点的 FLOPS。
  * `image` - 使用的镜像。
  * `image_type` - 镜像类型。
  * `use_origin_image` - 是否使用原始镜像。
  * `specification` - 规格字符串。

## 导入

BMCP 集群可以使用其 ID 导入，例如：

```bash
$ terraform import alibabacloudstack_bmcp_cluster.default bmcp-xxx
```

注意：导入时，某些字段如 `password`、`is_create_cpfs_cluster`、`is_install_yundun_aegis`、`standard_vswitch_id`、`cluster_arch_type`、`machine_type` 和 `vswitch_id` 将在状态验证期间被忽略。
