---
subcategory: "ASCM"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_ascm_ram_policy_for_role"
sidebar_current: "docs-apsarastack-resource-ascm-ram-policy-for-role"
description: |-
    Provides Ascm ram policy for role resource.
---

# apsarastack\_ascm_ram_policy_for_role

Provides a Ascm ram policy for role.

## Example Usage

```
resource "apsarastack_ascm_ram_policy" "default" {
  name = "Testpolicy"
  description = "Testing Complete"
  policy_document = "{\"Statement\":[{\"Action\":\"ecs:*\",\"Effect\":\"Allow\",\"Resource\":\"*\"}],\"Version\":\"1\"}"

}

resource "apsarastack_ascm_ram_role" "default" {
  role_name = "TestRole"
  description = "TestingRole"
  organization_visibility = "organizationVisibility.global"
}

resource "apsarastack_ascm_ram_policy_for_role" "default" {
  ram_policy_id = apsarastack_ascm_ram_policy.default.ram_id
  role_id = apsarastack_ascm_ram_role.default.role_id
}
output "ramrolebinder" {
  value = apsarastack_ascm_ram_policy_for_role.default.*
}

```
## Argument Reference

The following arguments are supported:

* `ram_policy_id` - (Required) ID of the ram_policy_id which will be used to bind.
* `role_id` - (Required, ForceNew) ID of the role which will be used to bind.

