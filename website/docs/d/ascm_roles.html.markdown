---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_ascm_roles"
sidebar_current: "docs-alibabacloudstack-datasource-ascm-roles"
description: |-
    Provides a list of roles to the user.
---

# alibabacloudstack_ascm_roles

This data source provides the roles of the current Apsara Stack Cloud user.

## Example Usage

```
resource "alibabacloudstack_ascm_ram_role" "default" {
  role_name = "DELTA1"
  description = "Testing Complete"
  organization_visibility = "organizationVisibility.global"
}

data "alibabacloudstack_ascm_roles" "default" {
  id = alibabacloudstack_ascm_ram_role.default.role_id
  name_regex = alibabacloudstack_ascm_ram_role.default.role_name
  role_type = "ROLETYPE_RAM"
}

output "roles" {
  value = data.alibabacloudstack_ascm_roles.default.*
}


```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) It is used to filter results by role ID.
* `name_regex` - (Optional) A regex string to filter results by role name.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `roles` - A list of roles. Each element contains the following attributes:
    * `id` - ID of the role.
    * `name` - role name.
    * `description` - Description about the role.
    * `role_level` - role level.
    * `role_type` - types of role.
    * `ram_role` - ram authorized role.
    * `role_range` - specific range for a role.
    * `user_count` - user count.
    * `enable` - enable.
    * `default` - default role.
    * `active` - Role status.
    * `owner_organization_id` - ID of the owner organization where role belongs.
    * `code` - role code.