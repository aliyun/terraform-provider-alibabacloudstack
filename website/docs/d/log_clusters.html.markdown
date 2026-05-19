---
layout: "alicloud-doc"
page_title: "Data Source: alibabacloudstack_log_clusters"
subcategory: "Simple Log Service"
---

# alibabacloudstack_log_clusters

This data source provides the available SLS Clusters.

## Example Usage

```hcl
data "alibabacloudstack_log_clusters" "default" {
  name_regex = "^my-cluster.*"
}

output "first_cluster_name" {
  value = data.alibabacloudstack_log_clusters.default.clusters.0.name
}
```

## Argument Reference

The following arguments are supported:

- `ids` - (Optional, ForceNew) A list of cluster IDs.
- `name_regex` - (Optional, ForceNew) A regex string to filter results by cluster name.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

- `ids` - A list of cluster IDs.
- `clusters` - A list of clusters. Each element contains the following attributes:
  - `name` - The name of the cluster.
  - `data_server` - The data server address of the cluster.
  - `status` - The status of the cluster.
  - `zone` - The zone where the cluster is located.
  - `description` - The description of the cluster.
