---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_user_role_binding"
sidebar_current: "docs-alibabacloudstack-resource-ascm-user-role-binding"
description: |-
  Provides Ascm User Role Binding.
---

# alibabacloudstack\_ascm_user_role_binding

## Example Usage

```
resource "alibabacloudstack_ascm_user_role_binding" "default" {
  role_id = 5
  login_name = "testUser"
}

output "binder" {
  value = alibabacloudstack_ascm_user_role_binding.default.*
}
```
## Argument Reference

The following arguments are supported:

* `role_ids` - (Required) ID list of the role which will be used to bind with user.
* `login_name` - (Required) Name of the User.

## Attributes Reference

The following attributes are exported:

* `id` - Name of the User.
* `login_name` - Name of the User.
* `role_ids` - User Role Id list.
