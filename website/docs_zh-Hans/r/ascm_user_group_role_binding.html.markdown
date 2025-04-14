---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_user_group_role_binding"
sidebar_current: "docs-alibabacloudstack-resource-ascm-user-role-binding"
description: |-
  编排绑定ASCM用户组和角色
---

# alibabacloudstack_ascm_user_group_role_binding

使用Provider配置的凭证在指定的资源集下编排绑定ASCM用户组和角色。

## 示例用法

```
resource "alibabacloudstack_ascm_organization" "default" {
 name = "Test_binder"
 parent_id = "1"
}

resource "alibabacloudstack_ascm_user_group" "default" {
 group_name =      "%s"
 organization_id = alibabacloudstack_ascm_organization.default.org_id
}

resource "alibabacloudstack_ascm_user_group_role_binding" "default" {
  role_ids = [5,]
  user_group_id = alibabacloudstack_ascm_user_group.default.user_group_id
}

output "binder" {
  value = alibabacloudstack_ascm_user_group_role_binding.default.*
}
```

## 参数说明

以下参数被支持：

* `user_group_id` - (必填) 用户组的 ID。该参数用于指定需要绑定角色的用户组。
* `role_ids` - (可选) 用户角色 ID 列表。该参数用于指定需要绑定到用户组的角色 ID。如果未提供，则默认不绑定任何角色。

## 属性说明

导出以下属性：

* `user_group_id` - (必填) 绑定成功后的用户组 ID，表示该用户组已成功与角色关联。
* `role_ids` - (必填) 已成功绑定到用户组的角色 ID 列表。该属性返回实际绑定的角色 ID 列表，可能与输入参数略有不同，具体取决于系统配置和权限限制。