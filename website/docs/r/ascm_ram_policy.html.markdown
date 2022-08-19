---
subcategory: "ASCM"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_ascm_ram_policy"
sidebar_current: "docs-apsarastack-resource-ascm-ram-policy"
description: |-
  Provides Ascm ram policy.
---

# apsarastack\_ascm_ram_policy

Provides Ascm ram policy.

## Example Usage

```
resource "apsarastack_ascm_ram_policy" "default" {
  name = "TestPolicy"
  description = "Testing"
  policy_document = "{\"Statement\":[{\"Action\":\"ecs:*\",\"Effect\":\"Allow\",\"Resource\":\"*\"}],\"Version\":\"1\"}"
}
output "rampolicy" {
  value = apsarastack_ascm_ram_policy.default.*
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required) Ram Policy name. 
* `policy_document` - (Required) Policy document of the policy.
* `description` - (Optional) Description for the ram policy.

## Attributes Reference

The following attributes are exported:

* `id` - Ram policy Name of the user.
* `ram_id` - The ID of the ram policy.