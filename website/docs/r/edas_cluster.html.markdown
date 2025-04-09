---
subcategory: "EDAS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_edas_cluster"
sidebar_current: "docs-Alibabacloudstack-edas-cluster"
description: |- 
  使用Provider配置的凭证在指定的资源集下编排企业级分布式应用服务（Edas） Cluster resource.
---

# alibabacloudstack_edas_cluster

使用Provider配置的凭证在指定的资源集下编排企业级分布式应用服务（Edas） Cluster resource.

## Example Usage

Basic Usage

```hcl
resource "alibabacloudstack_edas_cluster" "default" {
  cluster_name      = var.cluster_name
  cluster_type      = var.cluster_type
  network_mode      = var.network_mode
  logical_region_id = var.logical_region_id
  vpc_id            = var.vpc_id
}
```

## Argument Reference

The following arguments are supported:

* `cluster_name` - (Required, ForceNew) The name of the cluster that you want to create. It must be unique within the Alibaba Cloud account and cannot be modified after creation.
* `cluster_type` - (Required, ForceNew) The type of the cluster that you want to create. Valid values:
  * `1`: Swarm cluster.
  * `2`: ECS cluster.
  * `3`: Kubernetes cluster.
* `network_mode` - (Required, ForceNew) The network type of the cluster that you want to create. Valid values:
  * `1`: Classic network.
  * `2`: VPC.
* `logical_region_id` - (Optional, ForceNew) The ID of the logical region where the cluster is located. You can call the `ListUserDefineRegion` operation to query the logical region ID.
* `vpc_id` - (Optional, ForceNew) The ID of the Virtual Private Cloud (VPC) for the cluster. This parameter is required if `network_mode` is set to `2` (VPC).
* `region_id` - (Optional, ForceNew) The ID of the region where the cluster is located.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the EDAS cluster. It is formulated as `<cluster_id>`.


### Explanation of Changes

1. **Example Usage**: Added a complete HCL configuration example demonstrating how to use the `alibabacloudstack_edas_cluster` resource with required parameters.
2. **Argument Reference**:
   - Expanded descriptions for each argument to clarify their purpose and constraints.
   - Included valid values for `cluster_type` and `network_mode` to make it easier for users to understand the options available.
   - Renamed `region_id` to `logical_region_id` for consistency with the first document.
3. **Attributes Reference**:
   - Simplified the description of the `id` attribute to focus on its formulation and purpose.
   - Removed redundant information and ensured clarity in the exported attributes section. 

This updated documentation provides a comprehensive guide for using the `alibabacloudstack_edas_cluster` resource effectively.