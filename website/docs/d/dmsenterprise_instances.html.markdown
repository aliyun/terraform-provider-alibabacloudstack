---
subcategory: "DMSEnterprise"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dmsenterprise_instances"
sidebar_current: "docs-Alibabacloudstack-datasource-dmsenterprise-instances"
description: |- 
  Provides a list of dmsenterprise instances owned by an AlibabaCloudStack account.
---

# alibabacloudstack_dmsenterprise_instances
-> **NOTE:** Alias name has: `alibabacloudstack_dms_enterprise_instances`

This data source provides a list of DMS Enterprise Instances in an Alibaba Cloud account according to the specified filters.

## Example Usage

```terraform
# Declare the data source
data "alibabacloudstack_dmsenterprise_instances" "example" {
  env_type       = "test"
  instance_source = "RDS"
  instance_type   = "mysql"
  net_type        = "VPC"
  status          = "NORMAL"
  name_regex      = "^my-instance-.*"
  output_file     = "dms_enterprise_instances.json"
}

output "first_database_instance_id" {
  value = "${data.alibabacloudstack_dmsenterprise_instances.example.instances.0.instance_id}"
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional, ForceNew) A regex string to filter results by the alias (`instance_alias`) of the DMS Enterprise Instance.
* `instance_alias_regex` - (Optional, ForceNew) A regex string to filter results by the alias (`instance_alias`) of the DMS Enterprise Instance.
* `env_type` - (Optional, ForceNew) The type of the environment to which the database instance belongs. For example, `prod`, `test`, or `dev`.
* `instance_source` - (Optional, ForceNew) The source of the database instance. For example, `RDS`, `ECS`, or `OnPremise`.
* `instance_type` - (Optional, ForceNew) The type of the database instance. For example, `mysql`, `sqlserver`, or `postgresql`.
* `net_type` - (Optional, ForceNew) The network type of the database instance. Valid values include `CLASSIC` and `VPC`.
* `search_key` - (Optional, ForceNew) The keyword used to query database instances.
* `status` - (Optional, ForceNew) Filter the results by status of the DMS Enterprise Instances. Valid values include `NORMAL`, `UNAVAILABLE`, `UNKNOWN`, `DELETED`, and `DISABLE`.
* `tid` - (Optional, ForceNew) The ID of the tenant in Data Management (DMS) Enterprise.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `ids` - A list of DMS Enterprise IDs (Each of them consists of `host:port`).
* `names` - A list of DMS Enterprise names.
* `instances` - A list of DMS Enterprise Instances. Each element contains the following attributes:
  * `id` - The unique identifier for the DMS enterprise instance, formatted as `<host>:<port>`.
  * `data_link_name` - The name of the data link for the database instance.
  * `database_password` - The logon password of the database instance.
  * `database_user` - The logon username of the database instance.
  * `dba_id` - The ID of the database administrator (DBA) of the database instance.
  * `dba_nick_name` - The nickname of the DBA.
  * `ddl_online` - Indicates whether the online data description language (DDL) service was enabled for the database instance.
  * `ecs_instance_id` - The ID of the Elastic Compute Service (ECS) instance to which the database instance belongs.
  * `ecs_region` - The region where the database instance resides.
  * `env_type` - The type of the environment to which the database instance belongs.
  * `export_timeout` - The timeout period for exporting the database instance.
  * `host` - The endpoint of the database instance.
  * `instance_alias` - The alias of the database instance.
  * `instance_name` - Alias of the key `instance_alias`.
  * `instance_id` - The ID of the database instance.
  * `instance_source` - The source of the database instance.
  * `instance_type` - The type of the database instance.
  * `port` - The connection port of the database instance.
  * `query_timeout` - The timeout period for querying the database instance.
  * `safe_rule_id` - The ID of the security rule for the database instance.
  * `sid` - The system ID (SID) of the database instance.
  * `status` - The status of the database instance.
  * `use_dsql` - Indicates whether cross-database query was enabled for the database instance.
  * `vpc_id` - The ID of the Virtual Private Cloud (VPC) to which the database instance belongs.