---
subcategory: "PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_cluster_instance_types"
description: |-
  Provides a list of PolarDB cluster instance types.
---

# alibabacloudstack_polardb_cluster_instance_types

This data source provides a list of PolarDB cluster instance types in an Alibaba Cloud Stack environment.

## Example Usage

```hcl
data "alibabacloudstack_polardb_cluster_instance_types" "default" {
  db_type     = "MySQL"
  db_version  = "8.0"
  sorted_by   = "CPU"
  cpu_type    = "hygon"
  sub_category = "normal_exclusive"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of instance type IDs.
* `db_version` - (Optional) The database engine version.
* `cpu` - (Optional) The number of CPU cores.
* `cpu_type` - (Optional, ForceNew) The CPU type. Valid values: `intel`, `arm64`, `hygon`.
* `memory` - (Optional) The memory size in GB.
* `sorted_by` - (Optional, ForceNew) The sorting method. Valid values: `CPU`
* `sub_category` - (Optional, ForceNew) The sub category. Valid values: `normal_general`, `normal_exclusive`.
* `db_type` - (Optional, ForceNew) The database engine type. Valid values: `MySQL`, `PostgreSQL`, `Oracle`.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of instance type IDs.
* `instance_types` - A list of instance types. Each element contains the following attributes:
  * `id` - The ID of the instance type.
  * `proxy_mem` - The proxy memory.
  * `sub_category` - The sub category.
  * `cpu_type` - The CPU type.
  * `memory` - The memory size in GB.
  * `gmt_modify` - The modification time.
  * `spec_from` - The specification source.
  * `proxy_cpu` - The proxy CPU.
  * `product_type` - The product type.
  * `db_node_class` - The DB node class.
  * `product` - The product.
  * `db_version` - The database version.
  * `proxy_type` - The proxy type.
  * `db_node_num` - The number of DB nodes.
  * `db_type` - The database type.
  * `cpu` - The number of CPU cores.
  * `uni_key` - The unique key.
  * `proxy_class` - The proxy class.
  * `gmt_create` - The creation time.
  * `region_id` - The region ID.
  * `status` - The status.