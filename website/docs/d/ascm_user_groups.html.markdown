---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_user_groups"
sidebar_current: "docs-alibabacloudstack-datasource-ascm-user-groups"
description: |-
    Provides a list of users to the user groups.
---

# alibabacloudstack_ascm_user_groups

This data source provides a list of users to the user groups.

## Example Usage

```hcl
data "alibabacloudstack_ascm_user_groups" "example" {
  name_regex = "example-group"
}

output "user_groups" {
  value = data.alibabacloudstack_ascm_user_groups.example.groups
}
```

## Argument Reference
The following arguments are supported:

* `ids` - (Optional) A list of user group IDs to filter the results.
* `name_regex` - (Optional) A regex pattern to filter user groups by name.
* `organization_id` - (Optional) The ID of the organization to filter user groups.

## Attributes Reference
The following attributes are exported:

* `ids` - A list of IDs of the user groups.
* `names` - A list of names of the user groups.
* `organization_id` - The ID of the organization.
* `role_ids` - A list of role IDs associated with the user groups.
* `groups` - A list of user groups. Each element contains the following attributes:
    * `id` - The unique identifier of the user group.
    * `group_name` - The name of the user group.
    * `organization_id` - The ID of the organization.
    * `user_group_id` - The unique identifier of the user group.
    * `users` - A list of usernames in the user group.
    * `role_ids` - A list of role IDs associated with the user group.
