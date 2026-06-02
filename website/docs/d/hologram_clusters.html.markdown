---
subcategory: "Hologres"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_hologram_clusters"
description: |-
  Provides a list of Hologram Clusters to the user.
---

# alibabacloudstack_hologram_clusters

This data source provides a list of Hologram Clusters according to the specified filters.

## Example Usage

```hcl
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
}

data "alibabacloudstack_hologram_clusters" "default" {
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
  compute_type = "Standard"
  cpu = "intel"
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required) The zone ID to which the clusters belong.
* `compute_type` - (Optional) The compute type of the clusters. Valid values: `Standard`, `Follower`. Default to `Standard`.
* `cpu` - (Optional) The CPU brand. Default to `intel`.
* `name_regex` - (Optional) A regex string to filter results by cluster name.
* `ids` - (Optional) A list of cluster IDs.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of cluster IDs.
* `clusters` - A list of Hologram clusters. Each element contains the following attributes:
  * `id` - The ID of the cluster.
  * `type` - The type of the cluster.
  * `cpu` - The CPU brand of the cluster.
  * `cpu_arch` - The CPU architecture of the cluster.
  * `performance` - The performance level of the cluster.
  * `support_replica` - Whether the cluster supports replica.
  * `cluster` - The name of the cluster.