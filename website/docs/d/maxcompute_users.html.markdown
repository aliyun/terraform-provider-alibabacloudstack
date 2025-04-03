---
subcategory: "MaxCompute"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_maxcompute_users"
sidebar_current: "docs-alibabacloudstack-datasource-maxcompute-users"
description: |-
  Provides a datasource of Max Compute Users owned by an Alibaba Cloud account.
---

# alibabacloudstack_maxcompute_users

This data source provides Max Compute Users available to the user.[What is User](https://www.alibabacloud.com/help/en/maxcompute/latest/users)

## Example Usage

```hcl
data "alibabacloudstack_maxcompute_users" "example" {
  name_regex = "example-user"
}

output "users" {
  value = data.alibabacloudstack_maxcompute_users.example.users
}
```

## Argument Reference
The following arguments are supported:

* `ids` - (Optional) A list of user IDs to filter the results.
* `name_regex` - (Optional) A regex pattern to filter users by name.


## Attributes Reference
The following attributes are exported:

* `ids` - A list of IDs of the users.
* `users` - A list of users. Each element contains the following attributes:
    * `id` - The unique identifier of the user.
    * `user_id` - The user ID.
    * `user_pk` - The primary key of the user.
    * `user_name` - The name of the user.
    * `user_type` - The type of the user.
    * `organization_id` - The ID of the organization.
    * `organization_name` - The name of the organization.
    * `description` - The description of the user.
