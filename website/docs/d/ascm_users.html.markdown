---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_users"
sidebar_current: "docs-alibabacloudstack-datasource-ascm-users"
description: |-
    Provides a list of users to the user.
---

# alibabacloudstack_ascm_users

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
* `login_name` - (Optional) User login name for filtering.
* `role_id` - (Optional) Filter the results by the specified user role ID.
* `login_policy_id` - (Optional) Filter the results by the specified user login policy ID.
* `current_page` - (Optional) Current page number for pagination.
* `page_size` - (Optional) Number of items per page for pagination.
* `status` - (Optional) Filter the results by the specified user status.
* `mark` - (Optional) Marker for filtering users.

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
  * `default_role_id` - ID of the default role used by the user when logging in.
  * `cellphone_num` - Cellphone number of the user.
  * `default` - Indicates if the user is the default user.
  * `deleted` - Indicates if the user has been deleted.
  * `enable_ding_talk` - Indicates if DingTalk notifications are enabled for the user.
  * `enable_email` - Indicates if email notifications are enabled for the user.
  * `enable_short_message` - Indicates if short message notifications are enabled for the user.
  * `last_login_time` - Timestamp of the last login time for the user.
  * `parent_pk` - Parent primary key of the user.
  * `primary_key` - Primary key of the user.
  * `status` - Status of the user.
* `role_ids` - A list of all user owned roles.