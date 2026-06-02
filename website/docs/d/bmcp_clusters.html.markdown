---
subcategory: "Bare Metal Computing Platform (BMCP)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_bmcp_clusters"
description: |-
  Provides a list of BMCP (Bare Metal Compute Platform) clusters owned by an alibabacloudstack account.
---

# alibabacloudstack_bmcp_clusters

This data source provides a list of BMCP clusters in an AlibabacloudStack account according to the specified filters.

## Example Usage

### Query all clusters

```hcl
data "alibabacloudstack_bmcp_clusters" "all" {
}
```

### Filter by cluster name regex

```hcl
data "alibabacloudstack_bmcp_clusters" "name_filter" {
  cluster_name_regex = "my-cluster"
}
```

### Filter by cluster ID regex

```hcl
data "alibabacloudstack_bmcp_clusters" "id_filter" {
  cluster_id_regex = "bmcp-"
}
```

### Filter by status regex

```hcl
data "alibabacloudstack_bmcp_clusters" "status_filter" {
  status_regex = "active"
}
```

### Filter by region ID regex

```hcl
data "alibabacloudstack_bmcp_clusters" "region_filter" {
  region_id_regex = "cn-"
}
```

### Complete example with resource creation

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

# Query all clusters after creation
data "alibabacloudstack_bmcp_clusters" "all" {
  depends_on = [alibabacloudstack_bmcp_cluster.default]
}

# Filter by the created cluster name
data "alibabacloudstack_bmcp_clusters" "filtered" {
  cluster_name_regex = data.alibabacloudstack_bmcp_clusters.all.clusters.0.cluster_name
}
```

## Argument Reference

The following arguments are supported:

* `cluster_name_regex` - (Optional) A regex string to filter resulting clusters by cluster name.
* `cluster_id_regex` - (Optional) A regex string to filter resulting clusters by cluster ID.
* `status_regex` - (Optional) A regex string to filter resulting clusters by status.
* `region_id_regex` - (Optional) A regex string to filter resulting clusters by region ID.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of cluster IDs.
* `clusters` - A list of clusters. Each element contains the following attributes:
  * `id` - The ID of the cluster (same as `cluster_id`).
  * `cluster_id` - The unique identifier of the cluster.
  * `cluster_name` - The name of the cluster.
  * `status` - The status of the cluster.
  * `region_id` - The ID of the region where the cluster is located.
  * `zone` - The zone where the cluster is located.
  * `vpc_id` - The ID of the VPC where the cluster is located.
  * `evpc_id` - The ID of the EVPC associated with the cluster.
  * `organization` - The organization of the cluster.
  * `node_count` - The number of nodes in the cluster.
  * `cpu_count` - The total number of CPUs in the cluster.
  * `mem_count` - The total memory size of the cluster (in GB).
  * `flops_count` - The total FLOPS of the cluster.
  * `video_memory` - The total video memory of the cluster (in GB).
  * `cluster_arch_type` - The architecture type of the cluster.
  * `switch_method` - The switch method of the cluster.
  * `enable_ipv6` - Whether IPv6 is enabled.
  * `create_time` - The creation time of the cluster.
  * `update_time` - The last update time of the cluster.
  * `gpu` - The list of GPU information. Each element contains:
    * `gpu_num` - The number of GPUs.
    * `gpu_model` - The model of the GPU.
