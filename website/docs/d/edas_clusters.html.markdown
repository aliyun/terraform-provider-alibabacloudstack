---
subcategory: "EDAS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_edas_clusters"
sidebar_current: "docs-Alibabacloudstack-datasource-edas-clusters"
description: |- 
  Provides a list of edas clusters owned by an alibabacloudstack account.
---

# alibabacloudstack_edas_clusters

This data source provides a list of EDAS clusters in an Alibaba Cloud account according to the specified filters.

## Example Usage

```hcl
data "alibabacloudstack_edas_clusters" "clusters" {
  logical_region_id = "cn-shenzhen:xxx"
  ids               = ["addfs-dfsasd"]
  name_regex        = "test-cluster-*"
  output_file       = "clusters.txt"
}

output "first_cluster_name" {
  value = data.alibabacloudstack_edas_clusters.clusters.clusters[0].cluster_name
}
```

## Argument Reference

The following arguments are supported:

* `logical_region_id` - (Required, ForceNew) The ID of the namespace in EDAS. This is used to specify the logical region where the clusters reside.
* `ids` - (Optional) A list of cluster IDs to filter the results by specific cluster IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by the cluster name.


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of cluster names.
* `ids` - A list of cluster IDs.
* `clusters` - A list of clusters. Each cluster contains the following attributes:
  * `cluster_id` - The ID of the cluster that was queried.
  * `cluster_name` - The name of the cluster that was queried.
  * `cluster_type` - The type of the cluster that was queried. Valid values:
    * `1`: Swarm cluster.
    * `2`: ECS cluster.
    * `3`: Kubernetes cluster.
  * `create_time` - The time when the cluster was created.
  * `update_time` - The time when the cluster was last updated.
  * `cpu` - The total number of CPUs in the cluster.
  * `cpu_used` - The number of used CPUs in the cluster.
  * `mem` - The total amount of memory in the cluster. Unit: MB.
  * `mem_used` - The amount of used memory in the cluster. Unit: MB.
  * `network_mode` - The type of the network where the cluster is located. Valid values:
    * `1`: Classic network.
    * `2`: VPC.
  * `node_num` - The number of Elastic Compute Service (ECS) instances deployed to the cluster.
  * `vpc_id` - The ID of the virtual private cloud (VPC) where the cluster is located.
  * `region_id` - The ID of the logical zone where the cluster is located.
  * `logical_region_id` - The ID of the logical region where the cluster resides.
