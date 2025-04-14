---
subcategory: "ACK"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ack_clusters"
sidebar_current: "docs-Alibabacloudstack-datasource-ack-clusters"
description: |- 
  查询ack集群
---

# alibabacloudstack_ack_clusters
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_cs_kubernetes_clusters`

根据指定过滤条件列出当前凭证权限可以访问的ack集群列表。

## 示例用法

```hcl
variable "name" {
  default = "tf-testacckubernetes-3795850"
}

data "alibabacloudstack_zones" default {
  available_resource_creation = "VSwitch"
}

data "alibabacloudstack_instance_types" "default_m" {
  availability_zone     = data.alibabacloudstack_zones.default.zones.0.id
  cpu_core_count       = 2
  memory_size          = 4
  kubernetes_node_role = "Master"
}

data "alibabacloudstack_instance_types" "default_w" {
  availability_zone     = data.alibabacloudstack_zones.default.zones.0.id
  cpu_core_count       = 2
  memory_size          = 4
  kubernetes_node_role = "Worker"
}

resource "alibabacloudstack_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.1.0.0/21"
}

resource "alibabacloudstack_vswitch" "default" {
  name             = var.name
  vpc_id           = alibabacloudstack_vpc.default.id
  cidr_block       = "10.1.1.0/24"
  availability_zone = data.alibabacloudstack_zones.default.zones.0.id
}

resource "alibabacloudstack_cs_kubernetes" "default" {
  master_vswitch_ids = [alibabacloudstack_vswitch.default.id]
  new_nat_gateway    = true
  enable_ssh         = true
  password          = "inputYourCodeHere"
  master_disk_category = "cloud_efficiency"
  node_cidr_mask     = 24
  vpc_id             = alibabacloudstack_vpc.default.id
  worker_disk_category = "cloud_efficiency"
  worker_instance_types = ["ecs.e4.large"]
  master_count       = 3
  service_cidr       = "172.21.0.0/20"
  os_type            = "linux"
  name               = var.name
  master_instance_types = ["ecs.e4.large", "ecs.e4.large", "ecs.e4.large"]
  platform           = "CentOS"
  version            = "1.18.8-aliyun.1"
  worker_data_disk_category = "cloud_efficiency"
  worker_data_disk_size = 100
  proxy_mode         = "ipvs"
  master_disk_size   = 45
  worker_vswitch_ids = [alibabacloudstack_vswitch.default.id]
  slb_internet_enabled = true
  timeout_mins       = 25
  worker_disk_size   = 30
  num_of_nodes       = 3
  pod_cidr           = "172.20.0.0/16"
  delete_protection  = false
}

data "alibabacloudstack_ack_clusters" "example" {
  name_regex = "my-first-ack"

  output_file = "clusters_output.txt"
}

output "ack_cluster_ids" {
  value = data.alibabacloudstack_ack_clusters.example.ids
}

output "ack_cluster_names" {
  value = data.alibabacloudstack_ack_clusters.example.names
}
```

## 参数说明

以下参数是支持的：

* `ids` - (可选) 集群 ID 列表，用于过滤结果。如果不指定，将考虑所有集群。
* `name_regex` - (可选) 用于通过集群名称过滤结果的正则表达式字符串。
* `state` - (可选) 按集群状态过滤结果。有效值包括 `Running`、`Creating`、`Updating` 等。
* `enable_details` - (可选) 布尔值，默认为 `false`。将此参数设置为 `true` 将返回有关每个集群的更多详细信息，例如 `master_nodes`、`worker_nodes` 和 `connections`。
* `kube_config` - (可选) 布尔值，如果要获取 `ids` 中指定的集群的 kubeconfig，请将其设置为 `true`。

## 属性说明

除了上述参数外，还导出以下属性：

* `names` - 匹配的 ACK 集群的名称列表。
* `ids` - 匹配的 ACK 集群的 ID 列表。
* `clusters` - 匹配的 ACK 集群列表。每个元素包含以下属性：
  * `id` - ACK 集群的 ID。
  * `name` - ACK 集群的名称。
  * `availability_zone` - 集群所在的可用区。
  * `slb_internet_enabled` - 指示是否创建了面向互联网的 API 服务器负载均衡器。
  * `security_group_id` - 与集群的工作节点关联的安全组 ID。
  * `nat_gateway_id` - 用于启动 Kubernetes 集群的 NAT 网关 ID。
  * `vpc_id` - 集群所在的 VPC ID。
  * `vswitch_ids` - 集群所在的交换机 ID 列表。
  * `master_instance_types` - 主节点的实例类型列表。
  * `worker_instance_types` - 工作节点的实例类型列表。
  * `worker_numbers` - 集群中的工作节点数量。
  * `pod_cidr` - 使用 Flannel 时的 Pod 网络 CIDR 块。
  * `cluster_network_type` - 集群使用的网络类型，如 `flannel` 或 `terway`。
  * `node_cidr_mask` - 每个节点上 Pod 使用的网络掩码。
  * `log_config` - 包含一个元素的列表，其中包含有关关联日志存储的信息。它包括：
    * `type` - 收集日志的类型。
    * `project` - 日志服务项目名称。
  * `image_id` - 节点使用的镜像 ID。
  * `master_disk_size` - 主节点的系统盘大小。
  * `state` - 集群的当前状态。
  * `master_disk_category` - 主节点的系统盘类别。
  * `worker_disk_size` - 工作节点的系统盘大小。
  * `worker_disk_category` - 工作节点的系统盘类别。
  * `master_nodes` - 集群中的主节点列表。每个元素包括：
    * `id` - 主节点的 ID。
    * `name` - 主节点的名称。
    * `private_ip` - 主节点的私有 IP 地址。
  * `worker_nodes` - 集群中的工作节点列表。每个元素包括：
    * `id` - 工作节点的 ID。
    * `name` - 工作节点的名称。
    * `private_ip` - 工作节点的私有 IP 地址。
  * `connections` - Kubernetes 集群的连接信息映射。它包括：
    * `api_server_internet` - API 服务器的互联网端点。
    * `api_server_intranet` - API 服务器的内网端点。
    * `master_public_ip` - SSH 访问主节点的公共 IP 地址。
    * `service_domain` - 在集群内部访问服务的域。