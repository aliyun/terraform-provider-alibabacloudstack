---
subcategory: "PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_parameter_group"
description: |-
  Provides a PolarDB Parameter Group resource.
---

# alibabacloudstack_polardb_parameter_group

Provides a PolarDB Parameter Group resource.

For information about PolarDB Parameter Group and how to use it, see [What is Parameter Group](https://help.aliyun.com/document_detail/polardb-parameter-group.html).

-> **Note:** Some resources may have limits on creating parameters. For more details, see [PolarDB product documentation](https://help.aliyun.com/product/polardb.html).

## Example Usage

```hcl
variable "name" {
  default = "tf-testaccparamgroup"
}

resource "alibabacloudstack_polardb_parameter_group" "default" {
  engine               = "MySQL"
  engine_version       = "8.0"
  parameter_group_name = var.name
  parameter_group_desc = "Terraform Parameter Group"
  parameters = {
    "loose_multi_blocks_ddl_count" = "11"
  }
}
```

## Argument Reference

The following arguments are supported:

  * `engine` - (Required, ForceNew) The database engine type of the parameter group. Valid values: `MySQL`, `PostgreSQL`, `PPAS`.
  * `engine_version` - (Required, ForceNew) The database engine version. Valid values depend on the engine type. For MySQL: `5.6`, `5.7`, `8.0`.
  * `parameter_group_name` - (Required) The name of the parameter group. The name must be unique within a region.
  * `parameter_group_desc` - (Optional, Computed) The description of the parameter group.
  * `parameters` - (Optional) A map of parameters to configure in the parameter group. The key is the parameter name and the value is the parameter value.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

  * `id` - The ID of the parameter group. Same as `parameter_group_id`.
  * `parameter_group_id` - The ID of the parameter group.
  * `parameter_group_type` - The type of the parameter group.
  * `param_counts` - The number of parameters in the parameter group.
  * `force_restart` - Indicates whether the instance needs to be restarted after the parameters are modified.
  * `created` - The time when the parameter group was created.
  * `modified` - The time when the parameter group was last modified.

## Import

PolarDB Parameter Group can be imported using the parameter group ID, e.g.

```
$ terraform import alibabacloudstack_polardb_parameter_group.example pg-12345678
```
