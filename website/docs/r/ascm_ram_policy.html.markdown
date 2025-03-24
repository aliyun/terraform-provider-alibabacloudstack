---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_ram_policy"
sidebar_current: "docs-alibabacloudstack-resource-ascm-ram-policy"
description: |-
  Provides Ascm ram policy.
---

# alibabacloudstack_ascm_ram_policy

Provides Ascm ram policy.

## Example Usage

```
resource "alibabacloudstack_ascm_ram_policy" "default" {
  name = "TestPolicy"
  description = "Testing"
  policy_document = "{\"Statement\":[{\"Action\":\"ecs:*\",\"Effect\":\"Allow\",\"Resource\":\"*\"}],\"Version\":\"1\"}"
}
output "rampolicy" {
  value = alibabacloudstack_ascm_ram_policy.default.*
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required) Ram Policy name. 
* `policy_document` - (Required) Policy document of the policy.
* `description` - (Optional) Description for the ram policy.
* `name` - (Required)  Ram Policy name with length between 3 and 64 characters.
* `policy_document` - (Required)  Policy document of the policy.

## Attributes Reference

The following attributes are exported:

* `id` - Ram policy Name of the user.
* `ram_id` - The ID of the ram policy.  Exported attribute indicating the unique identifier for the RAM policy.