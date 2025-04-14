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

## 参数说明

以下参数被支持：

* `group_name` - (Required) 用户组名称。
* `organization_id` - (Required) 用户组织ID。
* `role_in_ids` - (Deprecated) 字段 'role_in_ids' 已被弃用，建议改为使用新字段 'role_ids'。
* `role_ids` - (Optional) ASCM角色ID列表。用于指定该用户组所关联的角色。
* `group_name` - (Required) 用户组名称（重复定义，实际只需一处即可）。
* `role_in_ids` - (Optional) 已弃用的ASCM角色ID列表。

## 属性说明

以下属性会被导出：

* `id` - 用户组的登录名，通常作为唯一标识符。
* `user_group_id` - 用户组的ID，用于唯一标识该用户组。
* `organization_id` - 用户所属的组织ID。
* `role_ids` - 与该用户组关联的ASCM角色ID列表。