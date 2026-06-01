---
subcategory: "ApsaraMQ for RocketMQ"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ons_clusters"
sidebar_current: "docs-Alibabacloudstack-datasource-ons-clusters"
description: |-
    Retrieves a list of Message Queue clusters
---

# alibabacloudstack_ons_clusters

Retrieves a list of Message Queue clusters available in the current region based on the specified filtering conditions.

## Example Usage

```hcl
data "alibabacloudstack_ons_clusters" "example" {
  ids = ["cluster1"]
}

output "first_cluster_id" {
  value = data.alibabacloudstack_ons_clusters.example.clusters[0].id
}
```

## Arguments Reference

The following arguments are supported:

* `ids` - (Optional) A list of cluster IDs used to filter the results.
* `name_regex` - (Optional) A regex string used to filter the results by cluster name.

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `ids` - A list of cluster IDs.
* `clusters` - A list of clusters. Each element contains the following attributes:
  * `id` - The cluster ID.
  * `name` - The cluster name.
  * `cpu_brand` - The CPU brand of the cluster nodes.
  * `cpu_arch` - The CPU architecture of the cluster nodes.
