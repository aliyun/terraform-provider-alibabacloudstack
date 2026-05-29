---
subcategory: "ApsaraDB for Redis"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_kvstore_parameter_group"
sidebar_current: "docs-Alibabacloudstack-kvstore-parameter_group"
description: |-
  Manages parameter templates for KVStore (Redis).
---

# alibabacloudstack_kvstore_parameter_group

Manages parameter templates for KVStore, used to create and manage parameter configuration sets for database engines such as Redis.

## Example Usage

### Basic Usage

```hcl
variable "name" {
  default = "tf-kvparmgroup31846"
}

resource "alibabacloudstack_kvstore_parameter_group" "default" {
  character_type       = "logic"
  engine               = "Redis"
  parameter_group_name = var.name
  engine_version       = "7.0"
  parameter_group_desc = var.name
  parameters {
    param_name = "resp_version"
    value      = "3"
  }
  parameters {
    param_name = "rt_threshold_ms"
    value      = "400"
  }
  parameters {
    param_name = "#no_loose_check-whitelist-always"
    value      = "yes"
  }
}
```

## Argument Reference

The following arguments are supported:

* `character_type` - (Required, ForceNew) The character type of the parameter template, e.g., `logic` for logical parameter templates. `normal` for physical parameter templates.
* `engine_version` - (Required, ForceNew) The database engine version, e.g., "7.0".
* `parameter_group_desc` - (Required, ForceNew) The description of the parameter template, with a length between 1 and 256 characters.
* `parameter_group_name` - (Required, ForceNew) The name of the parameter template, with a length between 1 and 128 characters.
* `engine` - (Optional, ForceNew) The type of the database engine. Default is "Redis".
* `parameters` - (Optional, ForceNew) A list of custom parameters. Each element contains:
  * `param_name` - (Required, ForceNew) The name of the parameter.
  * `value` - (Required, ForceNew) The value of the parameter.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the parameter template.
* `create_time` - The creation time of the parameter template, in ISO 8601 standard format (e.g., "2026-01-20T07:27Z").
* `is_dynamic` - Whether the parameter template is dynamic (1 for yes, 0 for no).
* `type` - The type identifier of the parameter template.