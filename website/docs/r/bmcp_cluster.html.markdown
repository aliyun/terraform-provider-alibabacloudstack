---
subcategory: "Bare Metal Computing Platform (BMCP)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_bmcp_cluster"
sidebar_current: "docs-Alibabacloudstack-bmcp-cluster"
description: |-
  Provides a BMCP (Bare Metal Compute Platform) Cluster resource.
---

# alibabacloudstack_bmcp_cluster

Provides a BMCP (Bare Metal Compute Platform) Cluster resource.

## Example Usage

### Basic Usage

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

## Argument Reference

The following arguments are supported:

* `cluster_name` - (Required, ForceNew) The name of the BMCP cluster.
* `vpc_id` - (Required, ForceNew) The ID of the VPC where the cluster is located.
* `evpc_id` - (Required, ForceNew) The ID of the EVPC (Enterprise VPC) associated with the cluster.
* `zone_id` - (Required, ForceNew) The ID of the availability zone where the cluster is located.
* `standard_vswitch_id` - (Required, ForceNew) The ID of the standard VSwitch for the cluster.
* `machine_type` - (Required, ForceNew) The machine type of the cluster nodes. You can query available machine types using the `alibabacloudstack_bmcp_machinetypes` data source.
* `node_count` - (Required) The number of nodes in the cluster.
* `vswitch_id` - (Required, ForceNew) The ID of the VSwitch for the cluster nodes.
* `password` - (Optional, ForceNew, Sensitive) The password for the cluster nodes. Either `password` or `activation_code` must be provided.
* `activation_code` - (Optional, ForceNew, Sensitive) The activation code for the cluster. Either `password` or `activation_code` must be provided.
* `switch_method` - (Optional, ForceNew) The switch method for the cluster. Valid values: `CHSW` (default), `ROCE`.
* `cluster_arch_type` - (Optional, ForceNew) The architecture type of the cluster. Valid values: `standard` (default), `high_performance`.
* `is_install_yundun_aegis` - (Optional, ForceNew) Whether to install Yundun Aegis security agent. Default: `true`.
* `is_create_cpfs_cluster` - (Optional, ForceNew) Whether to create a CPFS cluster. Default: `false`.
* `enable_ipv6` - (Optional, ForceNew) Whether to enable IPv6. Default: `false`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the cluster (same as `cluster_id`).
* `cluster_id` - The unique identifier of the cluster.
* `status` - The status of the cluster (e.g., `active`, `creating`, `deleting`).
* `region_id` - The ID of the region where the cluster is located.
* `cpu_count` - The total number of CPUs in the cluster.
* `mem_count` - The total memory size of the cluster (in GB).
* `flops_count` - The total floating point operations per second (FLOPS) of the cluster.
* `video_memory` - The total video memory of the cluster (in GB).
* `create_time` - The creation time of the cluster.
* `update_time` - The last update time of the cluster.
* `gpu` - The list of GPU information in the cluster. Each element contains:
  * `gpu_num` - The number of GPUs.
  * `gpu_model` - The model of the GPU.
* `machine_type_list` - The list of machine type details. Each element contains:
  * `machine_type` - The machine type name.
  * `node_count` - The number of nodes of this machine type.
  * `vswitch_id` - The VSwitch ID for this machine type.
  * `arch` - The CPU architecture.
  * `gpu` - The GPU model.
  * `cpu_count` - The number of CPUs per node.
  * `mem_count` - The memory size per node (in GB).
  * `gpu_num` - The number of GPUs per node.
  * `video_mem_count` - The video memory per node (in GB).
  * `flops_count` - The FLOPS per node.
  * `image` - The image used.
  * `image_type` - The type of the image.
  * `use_origin_image` - Whether to use the original image.
  * `specification` - The specification string.

## Import

BMCP Cluster can be imported using the id, e.g.

```bash
$ terraform import alibabacloudstack_bmcp_cluster.default bmcp-xxx
```

Note: When importing, some fields such as `password`, `is_create_cpfs_cluster`, `is_install_yundun_aegis`, `standard_vswitch_id`, `cluster_arch_type`, `machine_type`, and `vswitch_id` will be ignored during state verification.
