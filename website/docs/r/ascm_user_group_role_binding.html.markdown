---
subcategory: "Apsara Stack Cloud Management"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_user_group_role_binding"
sidebar_current: "docs-alibabacloudstack-resource-ascm-user-role-binding"
description: |-
  Provides Ascm User Group Role Binding.
---

# alibabacloudstack_ascm_user_group_role_binding

-> **Deprecated:** This resource is deprecated because ascm_user_group already includes corresponding functions.

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

* `user_group_id` - (Required, ForceNew) ID of the user group. Modifying this parameter will force the resource to be recreated.
* `role_ids` - (Optional) List of user role IDs to bind to the user group.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID, which is the same as `user_group_id`.
* `user_group_id` - ID of the user group.
* `role_ids` - List of user role IDs that are actually bound to the user group.

## Import

ASCM User Group Role Binding can be imported using the user group ID, e.g.

```
$ terraform import alibabacloudstack_ascm_user_group_role_binding.example 12345
```