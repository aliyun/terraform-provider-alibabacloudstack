---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_user_group"
sidebar_current: "docs-alibabacloudstack-resource-ascm-user_group"
description: |-
  Provides a Ascm user group resource.
---

# alibabacloudstack\_ascm_user_group

Provides a Ascm user group resource.

## Example Usage

```
resource "alibabacloudstack_ascm_organization" "default" {
  name = "Dummy_Test_1"
}

resource "alibabacloudstack_ascm_user_group" "default" {
   group_name = "test"
   organization_id = alibabacloudstack_ascm_organization.default.org_id
   role_in_ids =   []string{"2", "6"}
}

output "org" {
  value = alibabacloudstack_ascm_user_group.default.*
}
```
## Argument Reference

The following arguments are supported:

* `group_name` - (Required) group name. 
* `organization_id` - (Required) User Organization ID.
* `role_in_ids` - (Deprecated). Field 'role_in_ids' has been deprecated. New field 'role_ids' instead.
* `role_ids` - (Optional) ascm role id.

## Attributes Reference

The following attributes are exported:

* `id` - Login Name of the user group.
* `user_group_id` - ID of the user group.