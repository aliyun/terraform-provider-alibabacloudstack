---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_ram_role"
sidebar_current: "docs-alibabacloudstack-resource-ascm-ram-role"
description: |-
  Provides Ascm ram role.
---

# alibabacloudstack\_ascm_ram_role

Provides Ascm ram role.

## Example Usage

```
resource "alibabacloudstack_ascm_ram_role" "default" {
  role_name = "TestingRamRole"
  description = "TestingRam"
  organization_visibility = "organizationVisibility.global"
}
output "ramrole" {
  value = alibabacloudstack_ascm_ram_role.default.*
}
```
## Argument Reference

The following arguments are supported:

* `role_name` - (Required) Role name. 
* `organization_visibility` - (Required) organization visibility. Valid Values are - "organizationVisibility.organization", "organizationVisibility.orgAndSubOrgs" and "organizationVisibility.global".
* `description` - (Optional) Description for the ram role. Note - It should not contain any spaces.
* `role_range` - (Required) Range of permissions for role management
  * roleRange.orgAndSubOrgs - Organization and cascading subordinate organizations
  * roleRange.allOrganizations - All Organizations
  * roleRange.userGroup - ResourceGroupSet under the Organization

## Attributes Reference

The following attributes are exported:

* `id` - Ram Role Name of the user.
* `role_id` - The ID of the ram role.