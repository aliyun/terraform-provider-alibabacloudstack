---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_user_groups"
sidebar_current: "docs-alibabacloudstack-datasource-ascm-user-groups"
description: |-
    查询ascm用户组
---

# alibabacloudstack_ascm_user_groups

根据指定过滤条件列出当前凭证权限可以访问的用户组列表。

## 示例用法

```hcl
data "alibabacloudstack_ascm_user_groups" "example" {
  name_regex = "example-group"
}

output "user_groups" {
  value = data.alibabacloudstack_ascm_user_groups.example.groups
}
```

## 参数说明
支持以下参数：

* `ids` - (可选) 用于过滤结果的用户组 ID 列表。如果提供此参数，仅返回匹配这些 ID 的用户组。
* `name_regex` - (可选) 用于按名称过滤用户组的正则表达式模式。通过该参数，您可以筛选出符合特定名称模式的用户组。
* `organization_id` - (可选) 组织 ID，用于限定查询范围为特定组织下的用户组。如果提供此参数，仅返回属于该组织的用户组。

## 属性说明
导出以下属性：

* `ids` - 用户组的 ID 列表。每个 ID 唯一标识一个用户组。
* `names` - 用户组的名称列表。每个名称对应一个用户组。
* `organization_id` - 组织的 ID。表示这些用户组所属的组织。
* `role_ids` - 与用户组关联的角色 ID 列表。这些角色定义了用户组的权限范围。
* `groups` - 用户组列表。每个元素包含以下属性：
    * `id` - 用户组的唯一标识符，用于唯一标识一个用户组。
    * `group_name` - 用户组的名称，表示用户组的友好名称。
    * `organization_id` - 组织的 ID，表示该用户组所属的组织。
    * `user_group_id` - 用户组的唯一标识符，与 `id` 相同，用于唯一标识一个用户组。
    * `users` - 用户组中的用户名列表，表示属于该用户组的所有用户。
    * `role_ids` - 与用户组关联的角色 ID 列表，这些角色定义了用户组的权限范围。
