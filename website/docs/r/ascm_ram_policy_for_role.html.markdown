---
subcategory: "Apsara Stack Cloud Management"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_ram_policy_for_role"
sidebar_current: "docs-alibabacloudstack-resource-ascm-ram-policy-for-role"
description: |-
    Provides Ascm ram policy for role resource.
---

# alibabacloudstack_ascm_ram_policy_for_role

Provides a ASCM RAM policy for role resource to bind a RAM policy to a RAM role.

-> **Note:** This resource can also be referred to by the following alias: `apsarastack_ascm_ram_policy_for_role`.

## Example Usage

```
resource "alibabacloudstack_ascm_ram_policy" "default" {
  name = "Testpolicy"
  description = "Testing Complete"
  policy_document = "{\"Statement\":[{\"Action\":\"ecs:*\",\"Effect\":\"Allow\",\"Resource\":\"*\"}],\"Version\":\"1\"}"

}

resource "alibabacloudstack_ascm_ram_role" "default" {
  role_name = "TestRole"
  description = "TestingRole"
  organization_visibility = "organizationVisibility.global"
}

resource "alibabacloudstack_ascm_ram_policy_for_role" "default" {
  ram_policy_id = alibabacloudstack_ascm_ram_policy.default.ram_id
  role_id = alibabacloudstack_ascm_ram_role.default.role_id
}
```

## Argument Reference

The following arguments are supported:

* `ram_policy_id` - (Required, ForceNew) The ID of the RAM policy to bind to the role. Changing this forces a new resource to be created.
* `role_id` - (Required, ForceNew) The ID of the RAM role to which the policy will be bound. Changing this forces a new resource to be created.

## Attributes Reference

No attributes are currently defined for this resource.

## Import

ASCM RAM policy for role can be imported using the ram_policy_id and role_id separated by colon, e.g.

```
$ terraform import alibabacloudstack_ascm_ram_policy_for_role.example <ram_policy_id>:<role_id>
```