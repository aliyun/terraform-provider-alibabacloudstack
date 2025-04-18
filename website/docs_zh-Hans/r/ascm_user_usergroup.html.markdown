---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_usergroup_user"
sidebar_current: "docs-alibabacloudstack-resource-ascm-usergroup_user"
description: |-
  编排绑定ASCM用户和用户组
---

# alibabacloudstack_ascm_usergroup_user

使用Provider配置的凭证在指定的资源集下编排绑定ASCM用户和用户组。

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

## 参数参考

支持以下参数：

* `user_group_id` - (必填) 组ID。
* `login_names` - (可选) 用户登录名列表。

## 属性参考

导出以下属性：

* `id` - 用户组用户的登录名。