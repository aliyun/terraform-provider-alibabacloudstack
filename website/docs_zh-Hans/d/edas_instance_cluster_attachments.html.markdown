---
subcategory: "企业级分布式应用服务 EDAS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_edas_instance_cluster_attachments"
sidebar_current: "docs-Alibabacloudstack-datasource-edas-instance-cluster-attachments"
description: |-
  查询指定EDAS集群中已导入的ECS实例列表
---

# alibabacloudstack_edas_cluster_members

当前数据源用于查询阿里云EDAS（Enterprise Distributed Application Service）集群中已导入的ECS实例成员列表。通过指定集群ID，可检索该集群关联的所有实例信息，包括实例状态、ECU标识及时间戳等。

## 示例用法

```hcl

variable "name" {
  default = "tf6972"
}

variable "logical_id" {
  default = ":tf6972"
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


resource "alibabacloudstack_ecs_securitygroup" "default" {
  name   = "${var.name}_sg"
  vpc_id = alibabacloudstack_vpc_vpc.default.id
}

resource "alibabacloudstack_security_group_rule" "default" {
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "22/22"
  priority          = 1
  security_group_id = alibabacloudstack_ecs_securitygroup.default.id
  cidr_ip           = "192.168.0.0/16"
}


data "alibabacloudstack_images" "default" {
  name_regex = "^ubuntu_"
  //name_regex  = "arm_centos_7_6_20G_20211110.raw"
  //name_regex  = "^arm_centos_7"
  most_recent = true
  owners      = "system"
}


data "alibabacloudstack_instance_types" "all" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  sorted_by         = "CPU"
}

data "alibabacloudstack_instance_types" "default" {
  count = 8 # Traverse 1-8 core CPU configurations

  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  cpu_core_count    = count.index + 1 # 1-8
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

resource "alibabacloudstack_ecs_instance" "default" {
  image_id             = data.alibabacloudstack_images.default.images.0.id
  instance_type        = local.default_instance_type_id
  system_disk_category = data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0
  system_disk_size     = 20
  system_disk_name     = "test_sys_disk"
  security_groups      = [alibabacloudstack_ecs_securitygroup.default.id]
  instance_name        = "${var.name}_ecs"
  vswitch_id           = alibabacloudstack_vpc_vswitch.default.id
  zone_id              = data.alibabacloudstack_zones.default.zones.0.id
  is_outdated          = false
  enable_ipv6          = true
  ipv6_address_count   = 1
  lifecycle {
    ignore_changes = [
      instance_type,
      system_disk_category
    ]
  }
}



resource "alibabacloudstack_edas_namespace" "default" {
  description          = var.name
  namespace_logical_id = var.logical_id
  namespace_name       = var.name
}

resource "alibabacloudstack_edas_cluster" "default" {
  cluster_name      = var.name
  logical_region_id = alibabacloudstack_edas_namespace.default.namespace_logical_id
  network_mode      = "2"
  cluster_type      = "2"
  vpc_id            = alibabacloudstack_vpc_vpc.default.id
}

resource "alibabacloudstack_edas_instance_cluster_attachment" "default" {
  cluster_id  = alibabacloudstack_edas_cluster.default.id
  instance_ids =[ alibabacloudstack_ecs_instance.default.id]
}



data "alibabacloudstack_edas_instance_cluster_attachments" "default" {
  cluster_id = alibabacloudstack_edas_instance_cluster_attachment.default.cluster_id
  ids = [
    "${alibabacloudstack_edas_instance_cluster_attachment.default.id}"
  ]
}
```

## 参数说明
以下参数用于配置数据源查询条件：

- `cluster_id` (必填, 变更时重建)：EDAS集群的唯一标识符。 (必填, 变更时重建)
- `ids` (列表, 选填)：用于过滤结果的集群成员ID列表，每个ID格式为`ClusterId:InstanceId`。 (可选)

## 属性说明
以下属性从查询结果中导出：

- `id` (字符串)：集群成员的唯一标识符，格式为`ClusterId:InstanceId`。
- `cluster_id` (字符串)：EDAS集群的ID。
- `create_time` (整数)：集群成员创建的时间戳（Unix时间戳格式）。
- `ecu_id` (字符串)：ECU（Elastic Compute Unit）的唯一标识符。
- `ecs_id` (字符串)：ECS实例的ID（与`instance_id`值相同）。
- `instance_id` (字符串)：ECS实例的ID（即集群成员实例标识）。
- `status` (整数)：集群成员的当前状态（具体状态值参考EDAS文档）。
- `update_time` (整数)：集群成员最后更新的时间戳（Unix时间戳格式）。