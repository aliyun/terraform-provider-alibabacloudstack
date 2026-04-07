---
subcategory: "BMCP"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_bmcp_nodes"
sidebar_current: "docs-Alibabacloudstack-datasource-bmcp-nodes"
description: |-
  查询裸金属计算平台(BMCP)节点
---

# alibabacloudstack_bmcp_nodes

根据指定过滤条件列出当前凭证权限可以访问的BMCP节点列表。

## 示例用法

```hcl
data "alibabacloudstack_bmcp_nodes" "default" {
}
```

### 按节点名称正则表达式过滤

```hcl
data "alibabacloudstack_bmcp_nodes" "node_name_filter" {
  node_name_regex = "node"
}
```

### 按节点ID正则表达式过滤

```hcl
data "alibabacloudstack_bmcp_nodes" "node_id_filter" {
  node_id_regex = "node"
}
```

### 按集群ID正则表达式过滤

```hcl
data "alibabacloudstack_bmcp_nodes" "cluster_id_filter" {
  cluster_id_regex = "cluster"
}
```

### 按集群名称正则表达式过滤

```hcl
data "alibabacloudstack_bmcp_nodes" "cluster_name_filter" {
  cluster_name_regex = "cluster"
}
```

### 按SN正则表达式过滤

```hcl
data "alibabacloudstack_bmcp_nodes" "sn_filter" {
  sn_regex = "TC"
}
```

### 按专有网络IP正则表达式过滤

```hcl
data "alibabacloudstack_bmcp_nodes" "vpc_ip_filter" {
  vpc_ip_regex = "172"
}
```

### 按带外IP正则表达式过滤

```hcl
data "alibabacloudstack_bmcp_nodes" "out_of_band_ip_filter" {
  out_of_band_ip_regex = "10"
}
```

## 参数说明

以下参数是支持的：

* `node_name_regex` - (选填) 用于按节点名称过滤结果的正则表达式字符串。
* `node_id_regex` - (选填) 用于按节点ID过滤结果的正则表达式字符串。
* `cluster_id_regex` - (选填) 用于按集群ID过滤结果的正则表达式字符串。
* `cluster_name_regex` - (选填) 用于按集群名称过滤结果的正则表达式字符串。
* `sn_regex` - (选填) 用于按SN（序列号）过滤结果的正则表达式字符串。
* `vpc_ip_regex` - (选填) 用于按专有网络IP过滤结果的正则表达式字符串。
* `out_of_band_ip_regex` - (选填) 用于按带外IP过滤结果的正则表达式字符串。

## 属性说明

除了上述参数外，还导出以下属性：

* `ids` - 节点ID列表。
* `nodes` - 节点列表。每个元素包含以下属性：
  * `id` - 节点的ID。
  * `node_id` - 节点ID。
  * `node_name` - 节点的名称。
  * `cluster_id` - 集群ID。
  * `cluster_name` - 集群名称。
  * `sn` - 节点的序列号。
  * `vpc_ip` - 节点的专有网络IP地址。
  * `out_of_band_ip` - 节点的带外IP地址。
  * `key_pair_name` - 与节点关联的密钥对名称。
  * `status` - 节点的状态。
  * `machine_type` - 节点的机型。
  * `machine_type_name` - 节点的机型名称。
  * `cpu_arch` - 节点的CPU架构。
  * `cpu_number` - CPU数量。
  * `memory` - 内存大小（GB）。
  * `disk` - 磁盘大小（GB）。
  * `gpu_num` - GPU数量。
  * `gpu_model` - GPU型号。
  * `region_id` - 地域ID。
  * `create_time` - 节点的创建时间。
  * `update_time` - 节点的更新时间。
