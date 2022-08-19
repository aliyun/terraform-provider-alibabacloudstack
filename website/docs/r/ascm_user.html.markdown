---
subcategory: "ASCM"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_ascm_user"
sidebar_current: "docs-apsarastack-resource-ascm-user"
description: |-
  Provides a Ascm user resource.
---

# apsarastack\_ascm_user

Provides a Ascm user resource.

## Example Usage

```
resource "apsarastack_ascm_organization" "default" {
  name = "Dummy_Test_1"
}

resource "apsarastack_ascm_user" "default" {
   cellphone_number = "892399537"
   email = "test@gmail.com"
   display_name = "C2C-DEL3"
   organization_id = apsarastack_ascm_organization.default.org_id
   mobile_nation_code = "91"
   login_name = "C2C_apsarastack_C2C"
}

output "org" {
  value = apsarastack_ascm_user.default.*
}
```
## Argument Reference

The following arguments are supported:

* `login_name` - (Required) User login name. 
* `cell_phone_number` - (Required) Cellphone Number of a user.
* `display_name` - (Required) Display name of a user.
* `email` - (Required) Email ID of a user.
* `mobile_nation_code` - (Required) Mobile Nation Code of a user, where user belongs to.
* `organization_id` - (Required) User Organization ID.
* `login_policy_id` - (Optional) User login policy ID.

## Attributes Reference

The following attributes are exported:

* `id` - Login Name of the user.
* `user_id` - The ID of the user.