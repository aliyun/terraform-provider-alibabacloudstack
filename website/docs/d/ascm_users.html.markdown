---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_users"
sidebar_current: "docs-alibabacloudstack-datasource-ascm-users"
description: |-
    Provides a list of users to the user.
---

# alibabacloudstack\_ascm_users

This data source provides the users of the current Apsara Stack Cloud user.

## Example Usage

```
resource "alibabacloudstack_ascm_organization" "default" {
  name = "Dummy_Test_1"
}

resource "alibabacloudstack_ascm_user" "default" {
  cellphone_number = "899999537"
   email = "test@gmail.com"
   display_name = "C2C-DEL3"
   organization_id = alibabacloudstack_ascm_organization.default.org_id
   mobile_nation_code = "91"
   login_name = "C2C_alibabacloudstack_C2C"

}
output "org" {
  value = alibabacloudstack_ascm_user.default.*
}
data "alibabacloudstack_ascm_users" "users" {
 ids = [alibabacloudstack_ascm_user.user.user_id]
}
output "users" {
 value = data.alibabacloudstack_ascm_users.users.*
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of users IDs.
* `name_regex` - (Optional) A regex string to filter results by user login name.
* `organization_id` - (Optional) Filter the results by the specified user Organization ID.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `users` - A list of users. Each element contains the following attributes:
  * `id` - ID of the user.
  * `login_name` - User login name.
  * `cell_phone_number` - Cellphone Number of a user.
  * `display_name` - Display name of a user.
  * `email` - Email ID of a user.
  * `mobile_nation_code` - Mobile Nation Code of a user, where user belongs to.
  * `organization_id` - User Organization ID.
  * `login_policy_id` - User login policy ID.
  * `role_ids` - A list of the user owned roles.
  * `default_role_id` - ID of the default role used by the user when logging in
* `role_ids` - A list of all user owned roles.
