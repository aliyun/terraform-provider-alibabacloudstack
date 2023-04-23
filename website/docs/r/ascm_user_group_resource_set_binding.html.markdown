---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_user_group_resource_set_binding"
sidebar_current: "docs-alibabacloudstack-resource-ascm-user-group-resource-set-binding"
description: |-
  Provides Ascm User Role Binding.
---

# alibabacloudstack\_ascm_user_group_resource_set_binding

## Example Usage

```
resource "alibabacloudstack_ascm_organization" "default" {
 name = "Test_binder"
 parent_id = "1"
}

resource "alibabacloudstack_ascm_user_group" "default" {
 group_name =      "%s"
 organization_id = alibabacloudstack_ascm_organization.default.org_id
}


resource "alibabacloudstack_ascm_resource_group" "default" {
  organization_id = alibabacloudstack_ascm_organization.default.org_id
  name = "alibabacloudstack-terraform-resourceGroup"
}

resource "alibabacloudstack_ascm_user_group_resource_set_binding" "default" {
  resource_set_id = alibabacloudstack_ascm_resource_group.default.rg_id
  user_group_id = alibabacloudstack_ascm_user_group.default.user_group_id
  ascm_role_id = "2"
}

output "binder" {
  value = alibabacloudstack_ascm_user_group_resource_set_binding.default.*
}
```
## Argument Reference

The following arguments are supported:

* `resource_set_id` - (Required) List of resource group id.
* `user_group_id` - (Required) user group id.
* `ascm_role_id` - (Optional) ascm role id.

## Attributes Reference

The following attributes are exported:

* `resource_set_id` - (Required) List of resource group id.
* `user_group_id` - (Required) user group id.
