---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_hpcclusters"
sidebar_current: "docs-Alibabacloudstack-datasource-ecs-hpcclusters"
description: |- 
  Provides a list of ecs hpcclusters owned by an alibabacloudstack account.
---

# alibabacloudstack_ecs_hpcclusters
-> **NOTE:** Alias name has: `alibabacloudstack_ecs_hpc_clusters`

This data source provides a list of ECS HPC Clusters in an Alibabacloudstack account according to the specified filters.

## Example Usage

```terraform
data "alibabacloudstack_ecs_hpcclusters" "example" {
  ids        = ["hpc-bp1i09xxxxxxxx"]
  name_regex = "tf-testAcc"
}

output "first_ecs_hpc_cluster_id" {
  value = data.alibabacloudstack_ecs_hpcclusters.example.clusters.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew) A list of HPC Cluster IDs. 
* `name_regex` - (Optional, ForceNew) A regex string to filter results by HPC Cluster name.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of HPC Cluster names.
* `clusters` - A list of ECS HPC Clusters. Each element contains the following attributes:
    * `description` - The description of the ECS HPC Cluster.
    * `id` - The ID of the HPC Cluster. This is equivalent to `hpc_cluster_id`.
    * `hpc_cluster_id` - The unique identifier of the HPC Cluster.
    * `name` - The name of the ECS HPC Cluster.