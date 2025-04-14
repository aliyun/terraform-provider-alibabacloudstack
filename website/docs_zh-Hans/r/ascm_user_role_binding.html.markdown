---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_user_role_binding"
sidebar_current: "docs-alibabacloudstack-resource-ascm-user-role-binding"
description: |-
  编排绑定ASCM用户和角色
---

# alibabacloudstack_ascm_user_role_binding

使用Provider配置的凭证在指定的资源集下编排绑定ASCM用户和角色。

## 示例用法

```
resource "alibabacloudstack_ascm_user_role_binding" "default" {
  role_ids = ["5"]
  login_name = "testUser"
}

output "binder" {
  value = alibabacloudstack_ascm_user_role_binding.default.*
}
```

## 参数说明

以下参数为支持的配置项：

* `role_ids` - (必填) 角色ID列表，用于将指定的角色与用户绑定。每个角色ID对应一个具体的权限角色。
* `login_name` - (必填) 用户的登录名称，表示需要绑定角色的用户。

## 属性说明

以下属性为资源创建后导出的内容：

* `id` - 用户的唯一标识符，通常与用户的登录名称一致。
* `login_name` - 用户的登录名称，与输入参数中的`login_name`相同。
* `role_ids` - 已成功绑定到用户的角色ID列表，表示该用户所拥有的角色权限集合。