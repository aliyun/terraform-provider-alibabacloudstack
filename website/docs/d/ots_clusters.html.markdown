---
subcategory: "Table Store (OTS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ots_clusters"
description: |-
  Provides a datasource to retrieve the list of available Tablestore clusters.
---

# alibabacloudstack_ots_clusters

This data source retrieves the list of public Tablestore clusters available under your Alibaba Cloud account. You can filter clusters by name to obtain specific cluster information.

For more information about Tablestore and how to use it, see [What is Tablestore](https://www.alibabacloud.com/help/product/27278.htm).

## Example Usage

### Retrieve all public clusters
```hcl
data "alibabacloudstack_ots_clusters" "all" {}

output "cluster_names" {
  value = data.alibabacloudstack_ots_clusters.all.names
}
```

### Filter clusters by name using a regular expression
```hcl
data "alibabacloudstack_ots_clusters" "filtered" {
  name_regex = "^cn-"
}

output "matching_clusters" {
  value = data.alibabacloudstack_ots_clusters.filtered.clusters
}
```

### Specify a list of exact cluster names
```hcl
data "alibabacloudstack_ots_clusters" "specific" {
  names = ["cn-hangzhou", "cn-shanghai"]
}

output "selected_clusters" {
  value = data.alibabacloudstack_ots_clusters.specific.clusters
}
```

## Argument Reference

The following arguments are supported:

* `names` - (Optional) A list of cluster names to query. Only clusters whose names are in this list will be returned. Must contain at least one element.
* `name_regex` - (Optional) A regular expression used to match cluster names. Only clusters with names matching this regex will be returned. Must be a valid regular expression.

> **NOTE:** Both `names` and `name_regex` can be specified simultaneously. In this case, only clusters that satisfy both conditions (i.e., the intersection) will be returned.

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `clusters` - A list of cluster details. Each item contains:
  * `cluster_name` - The unique identifier of the cluster (e.g., `cn-hangzhou`).
  * `cluster_type` - The type of the cluster (e.g., `Public` for public clusters).
  * `alias_name` - The alias or display name of the cluster (typically a human-readable region name).
  * `support_replica` - Whether replica functionality is supported (boolean).

* `names` - A list of names of all matched clusters (string list), in the same order as the `clusters` list.

## Notes

- This data source returns only **public clusters** and does not include private clusters.
- The query results depend on your account’s permissions and service activation status in the corresponding regions.
- The `name_regex` uses Go’s regular expression syntax. Ensure it is valid; otherwise, Terraform will return an error.