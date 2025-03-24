---
subcategory: "MaxCompute"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_maxcompute_clusters"
sidebar_current: "docs-alibabacloudstack-datasource-maxcompute-clusters"
description: |-
  Provides a datasource of Max Compute Clusters owned by an Alibaba Cloud account.
---

# alibabacloudstack_maxcompute_clusters

This data source provides Max Compute clusters available to the user.[What is Cluster](https://www.alibabacloud.com/help/en/maxcompute)


## Example Usage

```hcl
data "alibabacloudstack_maxcompute_clusters" "example" {
  name_regex = "example-cluster"
}

output "clusters" {
  value = data.alibabacloudstack_maxcompute_clusters.example.clusters
}
```

## Argument Reference
The following arguments are supported:

* `ids` - (Optional) A list of cluster IDs to filter the results.
* `name_regex` - (Optional) A regex pattern to filter clusters by name.

## Attributes Reference
The following attributes are exported:

* `ids` - A list of IDs of the clusters.
* `clusters` - A list of clusters. Each element contains the following attributes:
    * `cluster` - The name of the cluster.
    * `core_arch` - The core architecture of the cluster.
    * `project` - The project associated with the cluster.
    * `region` - The region where the cluster is located.
