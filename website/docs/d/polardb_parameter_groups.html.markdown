---
subcategory: "PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_parameter_groups"
description: |-
  Provides a list of PolarDB parameter groups owned by an Alibaba Cloud Stack account.
---

# alibabacloudstack\_polardb\_parameter\_groups

This data source provides a list of PolarDB parameter groups in an Alibaba Cloud Stack account according to the specified filters.

## Example Usage

```hcl
data "alibabacloudstack_polardb_parameter_groups" "default" {
  engine         = "MySQL"
  engine_version = "8.0"
}

# Filter by name regex
data "alibabacloudstack_polardb_parameter_groups" "example" {
  name_regex = "^my-parameter-group"
}

# Filter by IDs
data "alibabacloudstack_polardb_parameter_groups" "ids" {
  ids = ["pg-xxxxxxxxxxxxx"]
}
```

## Argument Reference

The following arguments are supported:

  * `ids` - (Optional) A list of parameter group IDs to filter results.
  * `name_regex` - (Optional) A regex string to filter results by parameter group name.
  * `engine` - (Optional) The database engine type. Valid values: `MySQL`, `PostgreSQL`, `Oracle`.
  * `engine_version` - (Optional) The database engine version.
  * `parameter_group_type` - (Optional) The parameter group type.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

  * `ids` - A list of parameter group IDs.
  * `names` - A list of parameter group names.
  * `groups` - A list of parameter groups. Each element contains the following attributes:
    * `id` - The ID of the parameter group.
    * `parameter_group_id` - The ID of the parameter group.
    * `parameter_group_name` - The name of the parameter group.
    * `parameter_group_desc` - The description of the parameter group.
    * `engine` - The database engine type.
    * `engine_version` - The database engine version.
    * `parameter_group_type` - The parameter group type.
    * `force_restart` - Indicates whether a restart is required for the parameter to take effect.
    * `param_counts` - The number of parameters in the parameter group.
    * `created` - The time when the parameter group was created.
    * `modified` - The time when the parameter group was last modified.
    * `parameters` - A list of parameters in the parameter group. Each element contains:
      * `param_name` - The name of the parameter.
      * `param_value` - The value of the parameter.
