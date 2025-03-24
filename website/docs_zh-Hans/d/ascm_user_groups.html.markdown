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

## 参数参考
支持以下参数：

* `ids` - (可选) 用于过滤结果的用户组 ID 列表。
* `name_regex` - (可选) 用于按名称过滤用户组的正则表达式模式。
* `organization_id` - (可选) 用于过滤用户组的组织 ID。

## 属性参考
导出以下属性：

* `ids` - 用户组的 ID 列表。
* `names` - 用户组的名称列表。
* `organization_id` - 组织的 ID。
* `role_ids` - 与用户组关联的角色 ID 列表。
* `groups` - 用户组列表。每个元素包含以下属性：
    * `id` - 用户组的唯一标识符。
    * `group_name` - 用户组的名称。
    * `organization_id` - 组织的 ID。
    * `user_group_id` - 用户组的唯一标识符。
    * `users` - 用户组中的用户名列表。
    * `role_ids` - 与用户组关联的角色 ID 列表。