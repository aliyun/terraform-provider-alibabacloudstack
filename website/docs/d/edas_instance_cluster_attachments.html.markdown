---
subcategory: "Enterprise Distributed Application Service (EDAS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_edas_instance_cluster_attachments"
description: |-
  Queries the list of ECS instances imported into a specified EDAS cluster
---

# alibabacloudstack_edas_cluster_members

This data source retrieves the list of ECS instance members imported into an Alibaba Cloud EDAS (Enterprise Distributed Application Service) cluster. By specifying the cluster ID, you can obtain all associated instance information, including instance status, ECU (Elastic Compute Unit) identifier, and timestamps.

## Example Usage

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
  name_regex  = "^ubuntu_"
  //name_regex = "arm_centos_7_6_20G_20211110.raw"  # ARM CentOS image example
  //name_regex = "^arm_centos_7"                    # ARM CentOS regex pattern
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

## Argument Reference

The following arguments are supported:

* `cluster_id` (Required) - The unique identifier of the EDAS cluster. (Required, Forces new resource when changed)
* `ids` (Optional) - A list of cluster member IDs used to filter results. Each ID is in the format `ClusterId:InstanceId`. (Optional)

## Attributes Reference

The following attributes are exported:

* `id` (String) - The unique identifier of the cluster member, in the format `ClusterId:InstanceId`.
* `cluster_id` (String) - The ID of the EDAS cluster.
* `create_time` (Integer) - The timestamp when the cluster member was created, in Unix timestamp format.
* `ecu_id` (String) - The unique identifier of the ECU (Elastic Compute Unit).
* `ecs_id` (String) - The ID of the ECS instance (same as `instance_id`).
* `instance_id` (String) - The ID of the ECS instance (i.e., the cluster member instance identifier).
* `status` (Integer) - The current status of the cluster member (refer to EDAS documentation for specific status values).
* `update_time` (Integer) - The timestamp of the last update to the cluster member, in Unix timestamp format.