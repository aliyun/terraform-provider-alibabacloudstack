---
subcategory: "RDS"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_db_instances"
sidebar_current: "docs-apsarastack-datasource-db-instances"
description: |-
    Provides a collection of RDS instances according to the specified filters.
---

# apsarastack\_db\_instances

The `apsarastack_db_instances` data source provides a collection of RDS instances available in ApsaraStack account.
Filters support regular expression for the instance name, searches by tags, and other filters which are listed below.

## Example Usage

```
data "apsarastack_db_instances" "db_instances_ds" {
  name_regex = "data-\\d+"
  status     = "Running"
  tags       = {
    "type" = "database",
    "size" = "tiny"
  }

}

output "first_db_instance_id" {
  value = "${data.apsarastack_db_instances.db_instances_ds.instances.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to filter results by instance name.
* `ids` - (Optional) A list of RDS instance IDs. 
* `engine` - (Optional) Database type. Options are `MySQL`, `SQLServer`, `PostgreSQL` and `PPAS`. If no value is specified, all types are returned.
* `status` - (Optional) Status of the instance.
* `db_type` - (Optional) `Primary` for primary instance, `Readonly` for read-only instance, `Guard` for disaster recovery instance, and `Temp` for temporary instance.
* `vpc_id` - (Optional) Used to retrieve instances belong to specified VPC.
* `vswitch_id` - (Optional) Used to retrieve instances belong to specified `vswitch` resources.
* `connection_mode` - (Optional) `Standard` for standard access mode and `Safe` for high security access mode.
* `tags` - (Optional) A map of tags assigned to the DB instances. 
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of RDS instance IDs. 
* `names` - A list of RDS instance names. 
* `instances` - A list of RDS instances. Each element contains the following attributes:
  * `id` - The ID of the RDS instance.
  * `name` - The name of the RDS instance.
  * `db_type` - `Primary` for primary instance, `Readonly` for read-only instance, `Guard` for disaster recovery instance, and `Temp` for temporary instance.
  * `region_id` - Region ID the instance belongs to.
  * `create_time` - Creation time of the instance.
  * `expire_time` - Expiration time. Pay-As-You-Go instances never expire.
  * `status` - Status of the instance.
  * `engine` - Database type. Options are `MySQL`, `SQLServer`, `PostgreSQL` and `PPAS`. If no value is specified, all types are returned.
  * `engine_version` - Database version.
  * `net_type` - `Internet` for public network or `Intranet` for private network.
  * `connection_mode` - `Standard` for standard access mode and `Safe` for high security access mode.
  * `instance_type` - Sizing of the RDS instance.
  * `availability_zone` - Availability zone.
  * `master_instance_id` - ID of the primary instance. If this parameter is not returned, the current instance is a primary instance.
  * `guard_instance_id` - If a disaster recovery instance is attached to the current instance, the ID of the disaster recovery instance applies.
  * `temp_instance_id` - If a temporary instance is attached to the current instance, the ID of the temporary instance applies.
  * `readonly_instance_ids` - A list of IDs of read-only instances attached to the primary instance.
  * `vpc_id` - ID of the VPC the instance belongs to.
  * `vswitch_id` - ID of the VSwitch the instance belongs to.
  * `port` - () RDS database connection port.
  * `connection_string` - RDS database connection string.
  * `instance_storage` -  User-defined DB instance storage space.
