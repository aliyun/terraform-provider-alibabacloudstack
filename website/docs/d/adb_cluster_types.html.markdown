---
subcategory: "AnalyticDB for MySQL V3.0"
layout: "alibabacloudstack"
page_title: "ApsaraStack: alibabacloudstack_adb_cluster_types"
description: |-
  Provides a list of ADB cluster types.
---

# alibabacloudstack_adb_cluster_types

This data source provides the ADB cluster types available in ApsaraStack.


## Example Usage

```hcl
data "alibabacloudstack_adb_cluster_types" "example" {
  cpu_type = "Intel"
  cpu      = 16
  memory   = 64
}

output "adb_cluster_types" {
  value = data.alibabacloudstack_adb_cluster_types.example.instance_types
}
```

## Argument Reference

The following arguments are supported:

* `cpu_type` - (Optional) The CPU type of the ADB cluster.
* `cpu` - (Optional) The number of CPUs.
* `memory` - (Optional) The memory size in GB.
* `sorted_by` - (Optional, ForceNew) The field to sort by. Valid values: `CPU`, `Memory`.
* `status` - (Optional) The status of the cluster type.
* `cluster_type` - (Optional) The cluster type.
* `ids` - (Optional, ForceNew) A list of cluster type IDs.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of cluster type IDs.
* `instance_types` - A list of ADB cluster types. Each element contains the following attributes:
  * `id` - The ID of the cluster type.
  * `cpu_type` - The CPU type.
  * `cpu` - The number of CPUs.
  * `memory` - The memory size in GB.
  * `storage_type` - The storage type.
  * `storage_min` - The minimum storage size in GB.
  * `storage_max` - The maximum storage size in GB.
  * `mode` - The mode of the cluster.
  * `node_min` - The minimum number of nodes.
  * `node_max` - The maximum number of nodes.
  * `cluster_type` - The cluster type.
  * `status` - The status of the cluster type.
  * `series` - The series of the cluster.
  * `cluster_category` - The cluster category.
