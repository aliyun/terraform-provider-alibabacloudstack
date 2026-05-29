---
subcategory: "Cloud-Native Distributed Database PolarDB-X 2.0"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardbx_instance_types"
sidebar_current: "docs-Alibabacloudstack-datasource-polardbx-instance-types"
description: |-
  Query PolarDB-X database instance specifications.
---

# alibabacloudstack_polardbx_instance_types

Lists the PolarDB-X database instance specifications that can be accessed with the current credential permissions based on specified filter conditions.

## Example Usage

```hcl
data "alibabacloudstack_polardbx_instance_types" "example" {
  engine_version = "8.0"
  cpu            = 4
  memory         = 16
}

output "instance_types" {
  value = data.alibabacloudstack_polardbx_instance_types.example.instance_types
}
```
## Argument Reference
The following arguments are supported:

* `ids` - (Optional) - Used to filter instance types by ID list
* `spec_series` - (Optional) - Specification series, supports SHARE, SINGLE.
* `engine_version` - (Optional) - Database engine version (supports 5.7, 8.0)
* `cpu` - (Optional) - Number of CPU cores
* `cpu_type` - (Optional) - CPU architecture
* `spec_type` - (Optional) - Specification type, supports DN (Data Node), CN (Compute Node)
* `memory` - (Optional) - Memory size (GB)
* `sorted_by` - (Optional) - Sorting method (supports sorting by CPU or Memory)
* `series` - (Optional) - Instance series, supports enterprise, standard

## Attributes Reference
In addition to the arguments listed above, the following attributes are exported:

* `instance_types` - Instance type detailed information list, containing the following attributes
  * `id` - Instance type unique identifier
  * `cpu` - Number of CPU cores
  * `memory` - Memory size (GB)
