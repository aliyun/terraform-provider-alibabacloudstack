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

## 参数参考

以下参数被支持：

* `user_group_id` - (必填) 用户组的 ID。
* `role_ids` - (可选) 用户角色 ID。

## 属性参考

导出以下属性：

* `user_group_id` - 用户组的 ID。
* `role_ids` - 用户角色 ID。