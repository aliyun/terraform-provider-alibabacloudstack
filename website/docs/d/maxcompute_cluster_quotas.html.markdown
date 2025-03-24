---
subcategory: "MaxCompute"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_maxcompute_cluster_quotas"
sidebar_current: "docs-alibabacloudstack-datasource-maxcompute-cluster-quotas"
description: |-
  Provides a datasource of Max Compute cluster Quotas
---

# alibabacloudstack_maxcompute_cluster_quotas

This data source provides Max Compute cluster Quotas


## Example Usage

```hcl
data "alibabacloudstack_maxcompute_cluster_quotas" "example" {
  cluster = "example-cluster"
}

output "cluster_quotas" {
  value = data.alibabacloudstack_maxcompute_cluster_quotas.example
}
```

## Argument Reference
The following arguments are supported:

* `cluster` - (Required) The name of the MaxCompute cluster for which to retrieve quotas.

## Attributes Reference
The following attributes are exported:

* `cluster` - The name of the MaxCompute cluster.
* `cu_total` - The total number of Compute Units (CUs) allocated to the cluster.
* `disk_available` - The available disk space in the cluster.
* `cu_available` - The available number of Compute Units (CUs) in the cluster.
* `disk_total` - The total disk space allocated to the cluster.
