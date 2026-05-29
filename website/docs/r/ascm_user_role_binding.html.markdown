---
subcategory: "Application"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_user_role_binding"
sidebar_current: "docs-alibabacloudstack-resource-ascm-user-role-binding"
description: |-
  Provides Ascm User Role Binding.
---

# alibabacloudstack_ascm_user_role_binding

-> **Deprecation Note:** This resource may be removed in future versions. The `ascm_user` resource already includes corresponding functions. Use `ascm_user` instead when possible.

Provides a ASCM User Role Binding resource.

## Example Usage

### Using `role_id` (Recommended)

```hcl
resource "alibabacloudstack_ascm_user" "default" {
  cellphone_number   = "13900000000"
  email              = "test@example.com"
  display_name       = "Test User"
  organization_id    = data.alibabacloudstack_account.current.organization_id
  mobile_nation_code = "91"
  login_name         = "testUser"
  login_policy_id    = 1
}

resource "alibabacloudstack_ascm_user_role_binding" "default" {
  role_id    = 5
  login_name = alibabacloudstack_ascm_user.default.login_name
}
```

### Using `role_ids` (Deprecated since v3.21.0)

```hcl
resource "alibabacloudstack_ascm_user_role_binding" "default" {
  role_ids   = [5, 6, 7]
  login_name = alibabacloudstack_ascm_user.default.login_name
}
```

## Argument Reference

The following arguments are supported:

* `login_name` - (Required, ForceNew) The login name of the user to bind roles to. Modifying this parameter will force a new resource to be created.
* `role_id` - (Optional, ForceNew) The ID of a single role to bind to the user. This parameter conflicts with `role_ids`. Modifying this parameter will force a new resource to be created.

-> **Note:** When using `role_id`, the resource ID format is `login_name:role_id`.

* `role_ids` - (Optional, Computed, Deprecated) A set of role IDs to bind to the user. This parameter is deprecated starting from version 3.21.0 and is no longer supported for configuration. Use `role_id` instead. This parameter conflicts with `role_id`. The maximum number of role IDs is 10.

-> **Note:** When using `role_ids`, the resource ID format is `login_name`. Resources created with `role_ids` will not support deletion.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID. Format depends on whether `role_id` or `role_ids` is used.
* `login_name` - The login name of the user.
* `role_ids` - The list of role IDs bound to the user.

## Import

ASC User Role Binding can be imported using the resource ID. The ID format depends on how the resource was created:

- When created with `role_id`: `login_name:role_id`, e.g.

```shell
$ terraform import alibabacloudstack_ascm_user_role_binding.example testUser:5
```

- When created with `role_ids` (deprecated): `login_name`, e.g.

```shell
$ terraform import alibabacloudstack_ascm_user_role_binding.example testUser
```
