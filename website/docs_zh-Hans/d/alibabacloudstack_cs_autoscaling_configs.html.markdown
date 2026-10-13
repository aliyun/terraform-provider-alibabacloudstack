---
subcategory: "容器服务 Kubernetes 版"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cs_autoscaling_configs"
sidebar_current: "docs-Alibabacloudstack-datasource-cs-autoscaling-configs"
description: |-
  查询ACK集群自动扩缩容配置列表
---

# alibabacloudstack_cs_autoscaling_configs

该数据源提供容器服务Kubernetes版（ACK）集群的自动扩缩容配置列表。

## 示例用法

```hcl
variable "name" {
  default = "tf-testAccAutoscalingConfigs"
}

variable "existed_k8s_cluster_id" {
  default = ""
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details              = true
}

resource "random_password" "password" {
  count            = 1
  length           = 12
  special          = true
  override_special = "!@#$^&*()_"
  min_lower        = 1
  min_upper        = 1
  min_numeric      = 1
}

data "alibabacloudstack_instance_types" "all" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  sorted_by         = "CPU"
}

data "alibabacloudstack_instance_types" "default" {
  count = 8

  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  cpu_core_count    = count.index + 1
  sorted_by         = "Memory"
}

locals {
  filtered_default = [for d in data.alibabacloudstack_instance_types.default : d if length(d.ids) > 0]
  fallback_all     = length(data.alibabacloudstack_instance_types.all.ids) > 0 ? data.alibabacloudstack_instance_types.all.ids : []

  default_instance_type_id = coalesce(
    try(local.filtered_default[0].ids[0], null),
    try(local.fallback_all[0], null),
    "no-available-instance-type"
  )
}

data "alibabacloudstack_cs_kubernetes_clusters" "default" {
  ids = var.existed_k8s_cluster_id == "" ? [] : [var.existed_k8s_cluster_id]
}

locals {
  create_count = length(data.alibabacloudstack_cs_kubernetes_clusters.default.ids) > 0 ? 0 : 1
}

resource "alibabacloudstack_vpc_vpc" "default" {
  count      = local.create_count
  vpc_name   = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
  lifecycle {
    ignore_changes = [
      secondary_cidr_blocks,
    ]
  }
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  count        = local.create_count
  vswitch_name = "${var.name}_vsw"
  vpc_id       = alibabacloudstack_vpc_vpc.default.0.id
  cidr_block   = "172.16.1.0/24"
  zone_id      = data.alibabacloudstack_zones.default.zones.0.id
  lifecycle {
    ignore_changes = [
      tags
    ]
  }
}

resource "alibabacloudstack_ecs_securitygroup" "default" {
  count   = local.create_count
  name    = "${var.name}_sg"
  vpc_id  = alibabacloudstack_vpc_vpc.default.0.id
}

resource "alibabacloudstack_security_group_rule" "default" {
  count             = local.create_count
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "22/22"
  priority          = 1
  security_group_id = alibabacloudstack_ecs_securitygroup.default.0.id
  cidr_ip           = "192.168.0.0/16"
}

resource "alibabacloudstack_cs_kubernetes" "default" {
  count                = local.create_count
  name                 = var.name
  version              = "1.34.1-aliyun.1"
  os_type              = "linux"
  platform             = "AliyunLinux"
  num_of_nodes         = "1"
  master_count         = "3"
  master_vswitch_ids   = [alibabacloudstack_vpc_vswitch.default.0.id, alibabacloudstack_vpc_vswitch.default.0.id, alibabacloudstack_vpc_vswitch.default.0.id]
  master_instance_types = [local.default_instance_type_id, local.default_instance_type_id, local.default_instance_type_id]
  master_disk_category = data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0
  vpc_id               = alibabacloudstack_vpc_vpc.default.0.id
  worker_instance_types = [local.default_instance_type_id]
  worker_vswitch_ids   = [alibabacloudstack_vpc_vswitch.default.0.id]
  worker_disk_category = data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0
  password             = random_password.password.0.result
  pod_cidr             = "172.20.0.0/16"
  service_cidr         = "172.21.0.0/20"
  worker_disk_size     = "40"
  master_disk_size     = "40"
  slb_internet_enabled = "true"
  new_nat_gateway      = "false"
  security_group_id    = alibabacloudstack_ecs_securitygroup.default.0.id
  runtime {
    name    = "containerd"
    version = "2.1.5"
  }
}

locals {
  k8s_cluster_id = length(data.alibabacloudstack_cs_kubernetes_clusters.default.ids) > 0 ? data.alibabacloudstack_cs_kubernetes_clusters.default.ids.0 : alibabacloudstack_cs_kubernetes.default.0.id
}

data "alibabacloudstack_cs_autoscaling_configs" "default" {
  ids = [local.k8s_cluster_id]
}
```

## 参数参考

以下参数被支持：

* `ids` - (可选) 用于过滤自动扩缩容配置的集群ID列表。如果不提供，将查询所有集群。

## 属性参考

以下属性被导出：

* `id` - 返回的集群ID列表的哈希值。
* `ids` - 具有自动扩缩容配置的集群ID列表。
* `configs` - 自动扩缩容配置对象列表。每个元素包含以下属性：
  * `id` - 集群ID。
  * `cluster_id` - 集群ID。
  * `cool_down_duration` - 缩容活动后的冷却时间。
  * `unneeded_duration` - 节点在未使用状态下持续多长时间后才可被缩容。
  * `utilization_threshold` - 节点利用率阈值，低于此值的节点将被考虑缩容。
  * `gpu_utilization_threshold` - GPU利用率阈值，低于此值的GPU节点将被考虑缩容。
