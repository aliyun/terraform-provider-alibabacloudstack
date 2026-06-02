---
subcategory: "AnalyticDB for PostgreSQL"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_gpdb_instance_types"
sidebar_current: "docs-Alibabacloudstack-datasource-gpdb-instance-types"
description: |-
  Provides a list of GPDB Instance Types.
---

# alibabacloudstack_gpdb_instance_types

This data source provides the GPDB (Greenplum Database) Instance Types available in ApsaraStack.

-> **NOTE:** Available in ApsaraStack.

## Example Usage

```hcl
data "alibabacloudstack_gpdb_instance_types" "example" {
  cpu    = 4
  memory = 16
}

output "instance_types" {
  value = data.alibabacloudstack_gpdb_instance_types.example.instance_types
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew) A list of instance type IDs to filter results.
* `engine_version` - (Optional) The engine version to filter instance types.
* `cpu` - (Optional) The number of CPUs to filter instance types.
* `memory` - (Optional) The memory size in GB to filter instance types.
* `status` - (Optional) The status to filter instance types.
* `sorted_by` - (Optional, ForceNew) The field to sort by. Valid values: `CPU`, `Memory`.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of instance type IDs.
* `instance_types` - A list of GPDB Instance Types. Each element contains the following attributes:
  * `id` - The ID of the instance type.
  * `cpu` - The number of CPUs.
  * `memory` - The memory size in GB.
  * `engine` - The database engine.
  * `engine_version` - The engine version.
  * `connections` - The maximum number of connections.
  * `storage_type` - The storage type.
  * `storage_min` - The minimum storage size in GB.
  * `storage_max` - The maximum storage size in GB.
  * `specification` - The specification code.
  * `specification_label` - The specification label.
  * `db_instance_mode` - The database instance mode.
  * `db_instance_mode_label` - The database instance mode label.
  * `node` - The node configuration.
  * `region_id` - The region ID.
  * `status` - The status of the instance type.
  * `product` - The product name.
  * `storage` - The storage configuration.
  * `cpu_label` - The CPU label.
  * `memory_label` - The memory label.
  * `storage_label` - The storage label.
  * `engine_version_label` - The engine version label.
  * `gmt_create` - The creation timestamp.
  * `gmt_modify` - The modification timestamp.
  * `spec_from` - The specification source.
