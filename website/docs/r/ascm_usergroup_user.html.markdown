---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_usergroup_user"
sidebar_current: "docs-alibabacloudstack-resource-ascm-usergroup_user"
description: |-
  Provides a Ascm usergroup_user resource.
---

# alibabacloudstack\_ascm_usergroup_user

Provides a Ascm usergroup_user resource.

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

resource "alibabacloudstack_ascm_user" "default" {
 cellphone_number = "13900000000"
 email = "test@gmail.com"
 display_name = "C2C-DELTA"
 organization_id = alibabacloudstack_ascm_organization.default.org_id
 mobile_nation_code = "91"
 login_name = "User_Role_Test%d"
 login_policy_id = 1
}

resource "alibabacloudstack_ascm_usergroup_user" "default" {
  login_names = ["${alibabacloudstack_ascm_user.default.login_name}"]
  user_group_id = alibabacloudstack_ascm_user_group.default.user_group_id
}

output "org" {
  value = alibabacloudstack_ascm_usergroup_user.default.*
}
```
## Argument Reference

The following arguments are supported:

* `user_group_id` - (Required) group name. 
* `login_names` - (Required) List of user login name.

## Attributes Reference

The following attributes are exported:

* `id` - Login Name of the usergroup_user.