---
subcategory: "Relational Database Service(RDS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_rds_instance_types"
sidebar_current: "docs-alibabacloudstack-rds-instance-types"
description: |-
  Provides a list of Rds Instance types to be used by the alibabacloudstack_rds_instance resource.
---

# alibabacloudstack_rds_instance_types

This data source provides the Rds Instance types of AlibabacloudStack.

## Example Usage

```
data "alibabacloudstack_rds_instance_types" "default" {
  engine               = "MySQL"
  engine_version       = "5.7"
  sorted_by            = "CPU"
  series               = "dual_ha"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew) Specifies the ID range of instance specifications. If not specified, returns instance specifications across all availability zones.
* `engine` - (Optional, ForceNew) Filter the results to a specific rds engine. valid values: `PostgreSQL`, `MySQL`,`POLARDB`.
* `engine_version` - (Optional, ForceNew) Filter the results to a specific rds engine version.
* `series` - (Optional) Filter the results to a specific rds series.
* `cpu` - (Optional) Filter the results to a specific number of cpu cores.
* `cpu_type` - (Optional) Filter the results to a specific number of cpu type.
* `memory` - (Optional) Filter the results to a specific memory size in GB.
* `sorted_by` - (Optional, ForceNew) Sort mode, valid values: `CPU`, `Memory`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:


* `ids` - A list of instance series IDs matching all conditions.
* `names` - A list of instance series names matching all conditions.
* `instance_types` - A list of detailed instance series specifications. Each element contains the following attributes: 
  * `id` - Unique identifier of the instance series.
  * `name` - Name of the instance series.
  * `cpu` - Number of cpu cores.
  * `memory` - Memory size in GB.
  * `engine` - Type of rds engine.
  * `engine_version` - Egnine version for the rds instance type.
  * `cpu_type` - Cpu type for the rds instance type.
  * `series` - rds series id.
  * `connections` - Maximum Connections.
  * `storage_type` - Storage Type.
  * `storage_min` - Minimum Storage Capacity.
  * `storage_max` - Maximum Storage Capacity.
