---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_user_group_role_binding"
sidebar_current: "docs-alibabacloudstack-resource-ascm-user-role-binding"
description: |-
  Provides Ascm User Role Binding.
---

# alibabacloudstack\_ascm_user_group_role_binding

## Example Usage

```
resource "alibabacloudstack_ascm_organization" "default" {
 name = "Test_binder"
 parent_id = "1"
}

resource "alibabacloudstack_ascm_user_group" "default" {
 group_name =      "%s"
 organization_id = alibabacloudstack_ascm_organization.default.org_id
}

resource "alibabacloudstack_ascm_user_group_role_binding" "default" {
  role_ids = [5,]
  user_group_id = alibabacloudstack_ascm_user_group.default.user_group_id
}

output "binder" {
  value = alibabacloudstack_ascm_user_group_role_binding.default.*
}
```
## Argument Reference

The following arguments are supported:

* `user_group_id` - (Required) ID of user group.
* `role_ids` - (Required) User Role Id.

## Attributes Reference

The following attributes are exported:

* `user_group_id` - (Required) ID of user group.
* `role_ids` - (Required) User Role Id.
