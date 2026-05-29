---
subcategory: "Apsara Stack Cloud Management"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_ascm_roles"
sidebar_current: "docs-alibabacloudstack-datasource-ascm-roles"
description: |-
    Provides a list of roles to the user.
---

# alibabacloudstack_ascm_roles

This data source provides the roles of the current Apsara Stack Cloud user.

> **NOTE:** This data source can also be referred to by the following alias:
> - `alibabacloudstack_ascm_ram_roles`

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

* `id` - (Optional, Deprecated) It is used to filter results by role ID. This field is deprecated and will be removed in version 3.21.0. Please use `ids` instead.
* `ids` - (Optional, Available in v1.68.0+) A list of role IDs. The field is used to filter results by role IDs.
* `name_regex` - (Optional) A regex string to filter results by role name.
* `description` - (Optional) Description about the role.
* `role_type` - (Optional) Types of role.
* `output_file` - (Optional, Deprecated) The output file for the data source. This field is deprecated and will be removed in version 3.19.0. To write content to a file, use the `local_file` provider instead.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `roles` - A list of roles. Each element contains the following attributes:
    * `id` - ID of the role.
    * `name` - Role name.
    * `description` - Description about the role.
    * `role_level` - Role level.
    * `role_type` - Types of role.
    * `ram_role` - RAM authorized role.
    * `role_range` - Specific range for a role.
    * `user_count` - User count.
    * `enable` - Enable status.
    * `default` - Default role.
    * `active` - Role status.
    * `owner_organization_id` - ID of the owner organization where role belongs.
    * `code` - Role code.
    * `assume_role_policy_document` - The assume role policy document for RAM-authorized roles.
    * `organization_visibility` - The organization visibility scope for the role.

## Import

This data source does not support importing because it is a read-only data source.