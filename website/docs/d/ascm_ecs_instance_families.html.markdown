---
subcategory: "ASCM"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_ascm_ecs_instance_families"
sidebar_current: "docs-apsarastack-datasource-ascm-ecs-instance-families"
description: |-
    Provides a list of ecs instance families to the user.
---

# apsarastack\_ascm_ecs_instance_families

This data source provides the ecs instance families of the current Apsara Stack Cloud user.

## Example Usage

```
data "apsarastack_ascm_ecs_instance_families" "default" {
  status = "Available"
  output_file = "ecs_instance"
}
output "ecs_instance" {
  value = data.apsarastack_ascm_ecs_instance_families.default.*
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of ecs instance family IDs.
* `status` - (Required) Filter the results by specifying the status of ecs instance families.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `families` - A list of ecs instance families. Each element contains the following attributes:
    * `instance_type_family_id` - ID of the ecs instance families.
    * `generation` - generation of ecs instance families.
