---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_user"
sidebar_current: "docs-alibabacloudstack-resource-ascm-user"
description: |-
  Provides a Ascm user resource.
---

# alibabacloudstack_ascm_user

Provides a Ascm user resource.

## Example Usage

```
resource "alibabacloudstack_ascm_organization" "default" {
  name = "Dummy_Test_1"
}

resource "alibabacloudstack_ascm_user" "default" {
   cellphone_number = "892399537"
   email = "test@gmail.com"
   display_name = "C2C-DEL3"
   organization_id = alibabacloudstack_ascm_organization.default.org_id
   mobile_nation_code = "91"
   login_name = "C2C_alibabacloudstack_C2C"
}

output "org" {
  value = alibabacloudstack_ascm_user.default.*
}
```
## Argument Reference

The following arguments are supported:

* `login_name` - (Required) User login name. 
* `cellphone_number` - (Required) Cellphone Number of a user.
* `display_name` - (Required) Display name of a user.
* `email` - (Required) Email ID of a user.
* `mobile_nation_code` - (Required) Mobile Nation Code of a user, where user belongs to.
* `organization_id` - (Required) User Organization ID.
* `login_policy_id` - (Optional) User login policy ID.
* `role_ids` - A list of the user owned roles.
* `telephone_number` - (Optional) Telephone number of a user.
* `init_password` - (Optional) Init password for the user.

## Attributes Reference

The following attributes are exported:

* `id` - Login Name of the user.
* `user_id` - The ID of the user.
* `init_password` - Init Password of the user.
* `role_ids` - A list of the user owned roles.
* `organization_id` - Organization ID to which the current user belongs. Field 'organization_id' has been deprecated from provider version 1.0.32.