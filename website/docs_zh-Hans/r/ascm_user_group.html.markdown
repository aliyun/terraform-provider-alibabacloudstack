---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_user_group"
sidebar_current: "docs-alibabacloudstack-resource-ascm-user-group"
description: |-
  编排ASCM用户组
---

# alibabacloudstack_ascm_user_group

使用Provider配置的凭证在指定的资源集下编排ASCM用户组。

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

* `group_name` - (Required) 用户组名称。
* `organization_id` - (Required) 用户组织ID。
* `role_in_ids` - (Deprecated). 字段 'role_in_ids' 已被弃用，改为使用新字段 'role_ids'。
* `role_ids` - (Optional) ASCM角色ID。
* `group_name` - (Required) 用户组名称。
* `role_in_ids` - (Optional) ASCM角色ID。

## Attributes Reference

The following attributes are exported:

* `id` - 用户组的登录名。
* `user_group_id` - 用户组的ID。
* `organization_id` - 用户组织ID。
* `role_ids` - ASCM角色ID。