---
subcategory: "ASCM"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_ascm_user_role_binding"
sidebar_current: "docs-apsarastack-resource-ascm-user-role-binding"
description: |-
  Provides Ascm User Role Binding.
---

# apsarastack\_ascm_user_role_binding

## Example Usage

```
resource "apsarastack_ascm_user_role_binding" "default" {
  role_id = 5
  login_name = "testUser"
}

output "binder" {
  value = apsarastack_ascm_user_role_binding.default.*
}
```
## Argument Reference

The following arguments are supported:

* `role_id` - (Required) ID of the role which will be used to bind with user.
* `login_name` - (Required) Name of the User.

## Attributes Reference

The following attributes are exported:

* `id` - Name of the User.
* `login_name` - Name of the User.
* `role_id` - User Role Id.
