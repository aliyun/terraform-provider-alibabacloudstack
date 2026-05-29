---
subcategory: "裸金属计算平台"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_bmcp_clusters"
sidebar_current: "docs-Alibabacloudstack-datasource-bmcp-clusters"
description: |-
  提供指定 alibabacloudstack 账户拥有的 BMCP（裸金属计算平台）集群列表。
---

# alibabacloudstack_bmcp_clusters

此数据源根据指定的过滤器提供 AlibabacloudStack 账户中的 BMCP 集群列表。

## 使用示例

### 查询所有集群

```hcl
data "alibabacloudstack_bmcp_clusters" "all" {
}
```

### 按集群名称正则过滤

```hcl
data "alibabacloudstack_bmcp_clusters" "name_filter" {
  cluster_name_regex = "my-cluster"
}
```

### 按集群 ID 正则过滤

```hcl
data "alibabacloudstack_bmcp_clusters" "id_filter" {
  cluster_id_regex = "bmcp-"
}
```

### 按状态正则过滤

```hcl
data "alibabacloudstack_bmcp_clusters" "status_filter" {
  status_regex = "active"
}
```

### 按区域 ID 正则过滤

```hcl
data "alibabacloudstack_bmcp_clusters" "region_filter" {
  region_id_regex = "cn-"
}
```

### 完整示例（包含资源创建）

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

# 创建后查询所有集群
data "alibabacloudstack_bmcp_clusters" "all" {
  depends_on = [alibabacloudstack_bmcp_cluster.default]
}

# 使用创建的集群名称进行过滤
data "alibabacloudstack_bmcp_clusters" "filtered" {
  cluster_name_regex = data.alibabacloudstack_bmcp_clusters.all.clusters.0.cluster_name
}
```

## 参数参考

支持以下参数：

* `cluster_name_regex` - (可选) 用于按集群名称过滤结果的正则表达式字符串。
* `cluster_id_regex` - (可选) 用于按集群 ID 过滤结果的正则表达式字符串。
* `status_regex` - (可选) 用于按状态过滤结果的正则表达式字符串。
* `region_id_regex` - (可选) 用于按区域 ID 过滤结果的正则表达式字符串。

## 属性参考

除上述列出的参数外，还导出以下属性：

* `ids` - 集群 ID 列表。
* `clusters` - 集群列表。每个元素包含以下属性：
  * `id` - 集群的 ID (与 `cluster_id` 相同)。
  * `cluster_id` - 集群的唯一标识符。
  * `cluster_name` - 集群的名称。
  * `status` - 集群的状态。
  * `region_id` - 集群所在区域的 ID。
  * `zone` - 集群所在的可用区。
  * `vpc_id` - 集群所在 VPC 的 ID。
  * `evpc_id` - 与集群关联的 EVPC 的 ID。
  * `organization` - 集群的组织。
  * `node_count` - 集群中的节点数量。
  * `cpu_count` - 集群中的 CPU 总数。
  * `mem_count` - 集群的总内存大小 (GB)。
  * `flops_count` - 集群的总 FLOPS。
  * `video_memory` - 集群的总显存 (GB)。
  * `cluster_arch_type` - 集群的架构类型。
  * `switch_method` - 集群的交换方法。
  * `enable_ipv6` - 是否启用 IPv6。
  * `create_time` - 集群的创建时间。
  * `update_time` - 集群的最后更新时间。
  * `gpu` - GPU 信息列表。每个元素包含：
    * `gpu_num` - GPU 数量。
    * `gpu_model` - GPU 型号。
