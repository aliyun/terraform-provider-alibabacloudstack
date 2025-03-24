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

## 参数参考

支持以下参数：

* `role_ids` - (必填) 角色ID列表，将用于与用户绑定的角色。
* `login_name` - (必填) 用户的名称。

## 属性参考

导出以下属性：

* `id` - 用户的名称。
* `login_name` - 用户的名称。
* `role_ids` - 用户角色ID列表。