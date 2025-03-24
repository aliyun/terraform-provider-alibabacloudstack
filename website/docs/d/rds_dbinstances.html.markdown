---
subcategory: "RDS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_rds_dbinstances"
sidebar_current: "docs-Alibabacloudstack-datasource-rds-dbinstances"
description: |- 
  Provides a list of rds dbinstances owned by an alibabacloudstack account.
---

# alibabacloudstack_rds_dbinstances
-> **NOTE:** Alias name has: `alibabacloudstack_db_instances`

This data source provides a list of RDS DB instances in an AlibabacloudStack account according to the specified filters.

## Example Usage

```hcl
data "alibabacloudstack_rds_dbinstances" "db_instances_ds" {
  name_regex = "data-\\d+"
  status     = "Running"
  engine     = "MySQL"
  tags       = {
    "type" = "database",
    "size" = "tiny"
  }
}

output "first_db_instance_id" {
  value = "${data.alibabacloudstack_rds_dbinstances.db_instances_ds.instances.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to filter results by instance name.
* `ids` - (Optional) A list of RDS instance IDs.
* `engine` - (Optional) Database type. Valid values include: `MySQL`, `PostgreSQL`, `SQLServer`, and `MariaDB`. If no value is specified, all types are returned.
* `status` - (Optional) The status of the resource. For example: `Running`, `Stopped`, etc.
* `db_type` - (Optional) Type of the database instance. Valid values:
  * `Primary`: Primary instance.
  * `Readonly`: Read-only instance.
  * `Guard`: Disaster recovery instance.
  * `Temp`: Temporary instance.
* `vpc_id` - (Optional) The ID of the VPC to which the instance belongs.
* `vswitch_id` - (Optional) The ID of the VSwitch to which the instance belongs.
* `connection_mode` - (Optional) The access mode of the instance. Valid values:
  * `Standard`: Standard access mode.
  * `Safe`: High security access mode (Database proxy mode).
  > **Note**: SQL Server 2012, 2016, and 2017 only support the standard access mode.
* `tags` - (Optional) A map of tags assigned to the DB instances.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of RDS instance IDs.
* `names` - A list of RDS instance names.
* `instances` - A list of RDS instances. Each element contains the following attributes:
  * `id` - The ID of the RDS instance.
  * `uid` - Alias of `id`.
  * `name` - The name of the RDS instance.
  * `db_type` - Type of the database instance. Valid values: `Primary`, `Readonly`, `Guard`, and `Temp`.
  * `charge_type` - The payment type of the resource. Valid values: `PrePaid` or `PostPaid`.
  * `region_id` - The region ID to which the instance belongs.
  * `create_time` - The creation time of the instance.
  * `expire_time` - Expiration time of the instance. Format: `yyyy-MM-ddTHH:mm:ssZ` (UTC time). Pay-as-you-go instances have no expiration time.
  * `status` - The status of the instance.
  * `engine` - Database type. Valid values: `MySQL`, `PostgreSQL`, `SQLServer`, and `MariaDB`.
  * `engine_version` - Database version.
  * `net_type` - Network type. Valid values: `Internet` or `Intranet`.
  * `connection_mode` - The access mode of the instance. Valid values: `Standard` or `Safe`.
  * `instance_type` - Sizing of the RDS instance.
  * `availability_zone` - Availability zone.
  * `master_instance_id` - The ID of the primary instance. If this parameter is not returned, the current instance is a primary instance.
  * `guard_instance_id` - If a disaster recovery instance is attached to the current instance, the ID of the disaster recovery instance applies.
  * `temp_instance_id` - If a temporary instance is attached to the current instance, the ID of the temporary instance applies.
  * `readonly_instance_ids` - A list of IDs of read-only instances attached to the primary instance.
  * `vpc_id` - The ID of the VPC to which the instance belongs.
  * `vswitch_id` - The ID of the VSwitch to which the instance belongs.
  * `port` - Connection port for the RDS database.
  * `connection_string` - Connection string for the RDS database.
  * `instance_storage` - User-defined storage space for the RDS instance.