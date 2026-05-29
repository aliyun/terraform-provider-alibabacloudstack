---
subcategory: "企业级分布式应用服务 EDAS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_edas_instance_cluster_attachment"
sidebar_current: "docs-Alibabacloudstack-edas-instance-cluster-attachment"
description: |-
  将ECS实例导入到EDAS集群
---

# alibabacloudstack_edas_instance_cluster_attachment

将ECS实例导入到EDAS集群，实现应用部署环境的统一管理。

> **注意：** 该资源也可以使用以下别名引用：
> - `alibabacloudstack_edas_instanceclusterattachment`

## 示例用法

### 基础用法

```hcl

variable "name" {
  default = "tf16"
}

variable "logical_id" {
  default = ":tf16"
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
  instance_ids = [alibabacloudstack_ecs_instance.default.id]
}
```

## 参数说明

支持以下参数：

* `cluster_id` - (必填, 变更时重建) EDAS集群ID。需为有效的EDAS集群标识符，格式为UUID字符串。
* `instance_ids` - (必填, 变更时重建) 要导入到EDAS集群的ECS实例ID列表。

## 属性说明

导出以下属性：

* `id` - 资源ID，格式为`cluster_id:instance_id1,instance_id2,...`（多个实例ID用逗号分隔）。
* `status_map` - 实例在集群中的状态映射。键为实例ID，值为状态：`1`（运行中）、`0`（转换中）、`-1`（失败）、`-2`（离线）。
* `ecu_map` - 实例与ECU（弹性计算单元）的映射关系。键为实例ID，值为ECU ID。
* `cluster_member_ids` - 每个实例关联的集群成员ID映射。键为实例ID，值为集群成员ID。

## Import

EDAS Instance Cluster Attachment 可以使用 cluster_id 和 instance_ids 导入，格式为 `cluster_id:instance_id1,instance_id2,...`，例如：

```
$ terraform import alibabacloudstack_edas_instance_cluster_attachment.example cluster-abc123:i-xxx001,i-xxx002
```