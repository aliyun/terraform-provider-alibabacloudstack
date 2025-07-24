---
subcategory: "PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_instance_types"
sidebar_current: "docs-Alibabacloudstack-datasource-polardb-instance-types"
description: |-
  Provides a list of PolarDB instance types owned by an Alibaba Cloud Stack account.
---

# alibabacloudstack_polardb_instance_types

This data source provides the PolarDB instance types available in Alibaba Cloud Stack.

## Example Usage

```hcl
data "alibabacloudstack_polardb_instance_types" "example" {
  engine        = "MySQL"
  engine_version = "8.0"
  cpu           = 4
  memory        = 8
}

output "instance_types" {
  value = data.alibabacloudstack_polardb_instance_types.example.instance_types
}
```
## Argument Reference
The following arguments are supported:

  * `ids` - (Optional) A list of instance type IDs to filter results.
  * `engine` - (Optional) The database engine. Valid values: PolarDB_PG, MySQL, PolarDB_PPAS.
  * `engine_version` - (Optional) The version of the database engine. Required when engine is specified.
  * `cpu` - (Optional) The number of CPU cores.
  * `cpu_type` - (Optional) The CPU architecture. Valid values: intel, arm64.
  * `memory` - (Optional) The memory size in GB.
  * `sorted_by` - (Optional) Sorting method. Valid values: CPU, Memory.
  * `series` - (Optional) The series of the instance. Valid values: dual_ha, read_only.

## Attributes Reference
The following attributes are exported:

  * `instance_types` - A list of instance types. Each element contains the following attributes:
    * `id` - The unique identifier of the instance type.
    * `cpu` - The number of CPU cores.
    * `memory` - The memory size in GB.
    * `engine` - The database engine.
    * `engine_version` - The version of the database engine.
    * `cpu_type` - The CPU architecture.
    * `series` - The series of the instance.
    * `connections` - The maximum number of concurrent connections.
    * `storage_type` - The type of storage.
    * `storage_min` - The minimum storage capacity.
    * `storage_max` - The maximum storage capacity.