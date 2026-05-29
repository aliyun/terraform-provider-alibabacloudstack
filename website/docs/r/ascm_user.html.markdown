---
subcategory: "Application"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_user"
sidebar_current: "docs-alibabacloudstack-resource-ascm-user"
description: |-
  Provides a ASCM user resource.
---

# alibabacloudstack_ascm_user

Provides a ASCM user resource.

-> **Note:** The resource also supports the following alias: `apsarastack_ascm_user`.

## Example Usage

```terraform
resource "alibabacloudstack_ascm_logon_policy" "default" {
  name        = "default-policy"
  description = "Default login policy"
  rule        = "ALLOW"
}

resource "alibabacloudstack_ascm_user" "default" {
  login_name         = "example-user"
  display_name       = "Example User"
  email              = "example@example.com"
  cellphone_number   = "13800138000"
  mobile_nation_code = "86"
  login_policy_id    = alibabacloudstack_ascm_logon_policy.default.policy_id
  role_ids           = ["1"]
}
```
## Argument Reference

The following arguments are supported:

* `login_name` - (Required, ForceNew) User login name. This parameter uniquely identifies the user. Modifying this parameter will force a recreation of the resource.
* `display_name` - (Required) Display name of the user.
* `email` - (Required) Email address of the user.
* `cellphone_number` - (Required) Mobile phone number of the user.
* `mobile_nation_code` - (Required) Mobile nation code of the user, such as `86` for China.
* `login_policy_id` - (Required) Login policy ID to associate with the user.
* `role_ids` - (Optional) A list of role IDs to assign to the user. Must contain at least one role ID if specified.
* `telephone_number` - (Optional) Telephone number of the user.
* `organization_id` - (Deprecated, ForceNew) Organization ID to which the user belongs. This field is deprecated since provider version 1.0.32. The user will be created under the organization configured in the provider.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID, which is the login name of the user.
* `user_id` - The internal ID of the user.
* `user_uid` - The unique primary key (UID) of the user.
* `init_password` - The initial password generated for the user. This is only available after creation and cannot be set manually.
* `role_ids` - The list of role IDs assigned to the user.
* `organization_id` - Organization ID to which the user belongs. Deprecated since provider version 1.0.32.

## Import

ASCM User can be imported using the `login_name`, e.g.

```shell
$ terraform import alibabacloudstack_ascm_user.example example-user
```