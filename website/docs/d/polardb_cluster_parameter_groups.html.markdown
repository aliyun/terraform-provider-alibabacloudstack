---
subcategory: "PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_cluster_parameter_groups"
description: |-
  Provides a list of PolarDB cluster parameter groups owned by an Alibaba Cloud Stack account.
---

# alibabacloudstack\_polardb\_cluster\_parameter\_groups

This data source provides a list of PolarDB cluster-level parameter groups in an Alibaba Cloud Stack account according to the specified filters. The list is returned by the `DescribeParameterGroups` action of the PolarDB `2017-08-01` API.

## Example Usage

```hcl
data "alibabacloudstack_polardb_cluster_parameter_groups" "default" {
  engine         = "MySQL"
  engine_version = "8.0"
}

# Filter by name regex
data "alibabacloudstack_polardb_cluster_parameter_groups" "example" {
  name_regex = "^my-cluster-parameter-group"
}

# Filter by IDs
data "alibabacloudstack_polardb_cluster_parameter_groups" "ids" {
  ids = ["pg-xxxxxxxxxxxxx"]
}

output "first_parameter_group_id" {
  value = data.alibabacloudstack_polardb_cluster_parameter_groups.default.groups.0.parameter_group_id
}
```

## Argument Reference

The following arguments are supported:

  * `ids` - (Optional) A list of cluster parameter group IDs to filter results.
  * `name_regex` - (Optional) A regex string to filter results by cluster parameter group name. Must be a valid regular expression.
  * `engine` - (Optional) The database engine type. For PolarDB cluster parameter groups, the value is `MySQL`.
  * `engine_version` - (Optional) The database engine version.
  * `parameter_group_type` - (Optional) The cluster parameter group type.

-> **Note:** The provider lists all cluster parameter groups returned by `DescribeParameterGroups` and then applies the `ids` and `name_regex` filters locally. `engine`, `engine_version` and `parameter_group_type` are declared as optional inputs but the engine information is exposed on each element of `groups`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

  * `ids` - A list of cluster parameter group IDs.
  * `names` - A list of cluster parameter group names.
  * `groups` - A list of cluster parameter groups. Each element contains the following attributes:
    * `id` - The ID of the cluster parameter group.
    * `parameter_group_id` - The ID of the cluster parameter group.
    * `parameter_group_name` - The name of the cluster parameter group.
    * `parameter_group_desc` - The description of the cluster parameter group.
    * `engine` - The database engine type.
    * `engine_version` - The database engine version.
    * `parameter_group_type` - The type of the cluster parameter group.
    * `force_restart` - Indicates whether a restart is required for the parameters to take effect.
    * `param_counts` - The number of parameters in the cluster parameter group.
    * `created` - The time when the cluster parameter group was created.
    * `modified` - The time when the cluster parameter group was last modified.
    * `parameters` - A list of parameters in the cluster parameter group. Each element contains:
      * `param_name` - The name of the parameter.
      * `param_value` - The value of the parameter.
