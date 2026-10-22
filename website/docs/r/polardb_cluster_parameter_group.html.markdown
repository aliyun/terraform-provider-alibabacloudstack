---
subcategory: "PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_cluster_parameter_group"
description: |-
  Provides a PolarDB Cluster Parameter Group resource.
---

# alibabacloudstack_polardb_cluster_parameter_group

Provides a PolarDB cluster-level parameter group (parameter template) resource, which is applied to PolarDB clusters of the cluster edition.

The resource calls the `CreateParameterGroup`, `DescribeParameterGroups` and `DeleteParameterGroup` actions of the PolarDB `2017-08-01` API. The parameter group ID is returned by the `CreateParameterGroup` action and used as the Terraform resource ID.

-> **Note:** Some resources may have limits on creating parameters. For more details, see [PolarDB product documentation](https://help.aliyun.com/product/polardb.html).

-> **Note:** `parameters` are sent to the API only when the parameter group is created. Modifying the map afterwards does not update the parameter group in place, and the provider does not read the parameter values back, so a diff on `parameters` after creation may require recreating the resource.

## Example Usage

```hcl
variable "name" {
  default = "tf-testaccclusterparamgroup"
}

resource "alibabacloudstack_polardb_cluster_parameter_group" "default" {
  engine               = "MySQL"
  engine_version       = "8.0"
  parameter_group_name = var.name
  parameter_group_desc = "Terraform Cluster Parameter Group"
  parameters = {
    "wait_timeout" = "3600"
  }
}
```

## Argument Reference

The following arguments are supported:

  * `engine` - (Required, ForceNew) The database engine type of the cluster parameter group. Only `MySQL` is supported for PolarDB cluster parameter groups, and the engine is validated by the API.
  * `engine_version` - (Required, ForceNew) The database engine version. Valid values depend on the engine type. For MySQL: `5.6`, `5.7`, `8.0`.
  * `parameter_group_name` - (Required) The name of the parameter group.
  * `parameter_group_desc` - (Optional, Computed) The description of the parameter group.
  * `parameters` - (Optional) A map of parameters to configure in the parameter group. The key is the parameter name and the value is the parameter value. Both the key and the value are strings.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

  * `id` - The ID of the parameter group. Same as `parameter_group_id`.
  * `parameter_group_id` - The ID of the parameter group, returned by the `CreateParameterGroup` action.
  * `parameter_group_type` - The type of the parameter group.
  * `param_counts` - The number of parameters in the parameter group.
  * `force_restart` - Indicates whether the cluster needs to be restarted after the parameters are modified.
  * `created` - The time when the parameter group was created.
  * `modified` - The time when the parameter group was last modified.

## Import

PolarDB Cluster Parameter Group can be imported using the parameter group ID, e.g.

```
$ terraform import alibabacloudstack_polardb_cluster_parameter_group.example pg-12345678
```

-> **Note:** `parameters` is not returned when the parameter group is read, so it is not restored on import.
