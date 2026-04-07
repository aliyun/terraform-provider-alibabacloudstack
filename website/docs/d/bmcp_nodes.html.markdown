---
subcategory: "BMCP"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_bmcp_nodes"
sidebar_current: "docs-Alibabacloudstack-datasource-bmcp-nodes"
description: |-
  Provides a list of BMCP (Bare Metal Compute Platform) nodes owned by an alibabacloudstack account.
---

# alibabacloudstack_bmcp_nodes

This data source provides a list of BMCP nodes in an AlibabacloudStack account according to the specified filters.

## Example Usage

```hcl
data "alibabacloudstack_bmcp_nodes" "default" {
}
```

### Filter by node name regex

```hcl
data "alibabacloudstack_bmcp_nodes" "node_name_filter" {
  node_name_regex = "node"
}
```

### Filter by node ID regex

```hcl
data "alibabacloudstack_bmcp_nodes" "node_id_filter" {
  node_id_regex = "node"
}
```

### Filter by cluster ID regex

```hcl
data "alibabacloudstack_bmcp_nodes" "cluster_id_filter" {
  cluster_id_regex = "cluster"
}
```

### Filter by cluster name regex

```hcl
data "alibabacloudstack_bmcp_nodes" "cluster_name_filter" {
  cluster_name_regex = "cluster"
}
```

### Filter by SN regex

```hcl
data "alibabacloudstack_bmcp_nodes" "sn_filter" {
  sn_regex = "TC"
}
```

### Filter by VPC IP regex

```hcl
data "alibabacloudstack_bmcp_nodes" "vpc_ip_filter" {
  vpc_ip_regex = "172"
}
```

### Filter by out-of-band IP regex

```hcl
data "alibabacloudstack_bmcp_nodes" "out_of_band_ip_filter" {
  out_of_band_ip_regex = "10"
}
```

## Argument Reference

The following arguments are supported:

* `node_name_regex` - (Optional) A regex string to filter resulting nodes by node name.
* `node_id_regex` - (Optional) A regex string to filter resulting nodes by node ID.
* `cluster_id_regex` - (Optional) A regex string to filter resulting nodes by cluster ID.
* `cluster_name_regex` - (Optional) A regex string to filter resulting nodes by cluster name.
* `sn_regex` - (Optional) A regex string to filter resulting nodes by SN (Serial Number).
* `vpc_ip_regex` - (Optional) A regex string to filter resulting nodes by VPC IP.
* `out_of_band_ip_regex` - (Optional) A regex string to filter resulting nodes by out-of-band IP.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of node IDs.
* `nodes` - A list of nodes. Each element contains the following attributes:
  * `id` - The ID of the node.
  * `node_id` - The node ID.
  * `node_name` - The name of the node.
  * `cluster_id` - The cluster ID.
  * `cluster_name` - The cluster name.
  * `sn` - The serial number of the node.
  * `vpc_ip` - The VPC IP address of the node.
  * `out_of_band_ip` - The out-of-band IP address of the node.
  * `key_pair_name` - The key pair name associated with the node.
  * `status` - The status of the node.
  * `machine_type` - The machine type of the node.
  * `machine_type_name` - The machine type name of the node.
  * `cpu_arch` - The CPU architecture of the node.
  * `cpu_number` - The number of CPUs.
  * `memory` - The memory size in GB.
  * `disk` - The disk size in GB.
  * `gpu_num` - The number of GPUs.
  * `gpu_model` - The GPU model.
  * `region_id` - The region ID.
  * `create_time` - The creation time of the node.
  * `update_time` - The update time of the node.
