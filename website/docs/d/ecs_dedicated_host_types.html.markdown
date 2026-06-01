---
subcategory: "Elastic Compute Service"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_dedicated_host_types"
sidebar_current: "docs-Alibabacloudstack-datasource-ecs-dedicated-host-types"
description: |-
  Provides a list of available dedicated host types.
---

# alibabacloudstack\_ecs\_dedicated\_host\_types

-> **NOTE:** Alias names include: `alibabacloudstack_ecs_dedicatedhost_types`

This data source provides a list of available dedicated host types in an Apsara Stack cloud account according to the specified filters.

## Example Usage

```hcl
# Query all available dedicated host types
data "alibabacloudstack_ecs_dedicated_host_types" "default" {
}

# Output the first dedicated host type ID
output "first_ddh_type_id" {
  value = "${data.alibabacloudstack_ecs_dedicated_host_types.default.ddh_types.0.id}"
}

# Filter dedicated host types by availability zone
data "alibabacloudstack_ecs_dedicated_host_types" "filtered" {
  availability_zone = "cn-hangzhou-a"
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Optional, ForceNew) The ID of the availability zone. Used to filter dedicated host types available in the specified zone.
* `ids` - (Optional, ForceNew, Computed) A list of dedicated host type IDs to filter results.

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `ids` - A list of dedicated host type IDs.
* `ddh_types` - A list of dedicated host types. Each element contains the following attributes:
  * `id` - The ID of the dedicated host type, e.g., `ddh.g5`, `ddh.c5`.
  * `availability_zones` - A list of availability zone IDs where this type is available.
