subcategory: "Quick BI"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_quick_bi_users"
sidebar_current: "docs-alibabacloudstack-datasource-quick-bi-users"
description: |-
  Provides a list of Quick BI Users to the user.
---

# alibabacloudstack_quick_bi_users

This data source provides the Quick BI Users of the current Alibaba Cloud user.


## Example Usage

Basic Usage

```terraform
data "alibabacloudstack_quick_bi_users" "ids" {
  ids = ["example_id"]
}
output "quick_bi_user_id_1" {
  value = data.alibabacloudstack_quick_bi_users.ids.users.0.id
}
```

## Argument Reference

The following arguments are supported:

* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `ids` - (Optional, ForceNew)  A list of User IDs.
* `keyword` - (Optional, ForceNew) The keywords of the nicknames or usernames of the members of the organization.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `users` - A list of Quick BI Users. Each element contains the following attributes:
    * `account_id` - Alibaba Cloud account ID.
    * `account_name` - An Alibaba Cloud account, Alibaba Cloud name.
    * `admin_user` - Whether it is the administrator. Valid values: `true` and `false`.
    * `auth_admin_user` - Whether this is a permissions administrator. Valid values: `true` and `false`.
    * `email` - The email of the user.
    * `id` - The ID of the User.
    * `nick_name` - The nickname of the user.
    * `phone` - The phone number of the user.
    * `user_id` - The ID of the User.
    * `user_type` - The members of the organization of the type of role separately. Valid values: `Analyst`, `Developer` and `Visitor`.