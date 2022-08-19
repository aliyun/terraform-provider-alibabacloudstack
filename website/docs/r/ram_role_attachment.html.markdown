---
subcategory: "RAM"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_ram_role_attachment"
sidebar_current: "docs-apsarastack-resource-ram-role-attachment"
description: |-
  Provides a RAM role attachment resource to bind role for several ECS instances.
---

# apsarastack\_ram\_role\_attachment

Provides a RAM role attachment resource to bind role for several ECS instances.

## Example Usage

```
data "apsarastack_ascm_ram_service_roles" "role" {
  product = "ecs"
}

resource "apsarastack_ram_role_attachment" "attach" {
  role_name    = data.apsarastack_ascm_ram_service_roles.role.roles.0.name
  instance_ids = ["i-23jkek3dkhsdby8kba"]
}

output "attach" {
  value = apsarastack_ram_role_attachment.attach.*
}
```

## Argument Reference

The following arguments are supported:

* `role_name` - (Required, ForceNew) The name of role used to bind. This name can have a string of 1 to 64 characters, must contain only alphanumeric characters or hyphens, such as "-", "_", and must not begin with a hyphen.
* `instance_ids` - (Required, ForceNew) The list of ECS instance's IDs.

## Attributes Reference

The following attributes are exported:

* `role_name` - The name of the role.
* `instance_ids` The list of ECS instance's IDs.
