---
subcategory: "ApsaraDB for MongoDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_mongodb_instance_types"
sidebar_current: "docs-Alibabacloudstack-datasource-mongodb-instance-types"
description: |-
  Provides a list of Mongodb Instance types to be used by the alibabacloudstack_mongodb_instance resource.
---

# alibabacloudstack_mongodb_instance_types

This data source provides the Mongodb Instance types of AlibabacloudStack.

## Example Usage

```
data "alibabacloudstack_mongodb_instance_types" "default" {
  engine_version = "4.0"
  db_instnace_type = "replicate"
  sorted_by = "CPU"
}

data "alibabacloudstack_mongodb_instance_types" "shard" {
  engine_version = "4.0"
  db_instnace_type = "sharding"
  node_type = "shard"
  sorted_by = "CPU"
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew) Specifies the ID range of instance specifications. If not specified, returns instance specifications across all availability zones.
* `db_instnace_type` - (Required, ForceNew) Mongodb instance type, valid values: `sharding`, `replicate`.
* `node_type` - (Optional, ForceNew) Mongodb sharding instance node type, valid values: `configserver`, `shard`, `mongos`.
* `engine_version` - (Optional, ForceNew) Filter the results to a specific mongodb engine version.
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
  * `engine_version` - Egnine version for the mongodb instance type.
  * `cpu_type` - Cpu type for the mongodb instance type.
  * `series` - mongodb series id.
  * `connections` - Maximum Connections.
  * `storage_min` - Minimum Storage Capacity.
  * `storage_max` - Maximum Storage Capacity.
