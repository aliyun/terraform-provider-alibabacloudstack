---
subcategory: "企业控制台(ASCM)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_user_group_role_binding"
sidebar_current: "docs-alibabacloudstack-resource-ascm-user-role-binding"
description: |-
  编排绑定ASCM用户组和角色
---

# alibabacloudstack_ascm_user_group_role_binding

-> **已弃用：** 该资源已被弃用，因为 ascm_user_group 资源已包含对应功能。

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

* `user_group_id` - (必填，ForceNew) 用户组的 ID。修改此参数会强制重新创建资源。
* `role_ids` - (可选) 要绑定到用户组的角色 ID 列表。

## 属性说明

导出以下属性：

* `id` - 资源 ID，与 `user_group_id` 相同。
* `user_group_id` - 用户组 ID。
* `role_ids` - 实际绑定到用户组的角色 ID 列表。

## Import

ASCM 用户组角色绑定可以通过用户组 ID 导入，例如：

```
$ terraform import alibabacloudstack_ascm_user_group_role_binding.example 12345
```