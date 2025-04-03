---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_user"
sidebar_current: "docs-alibabacloudstack-resource-ascm-user"
description: |-
  编排ASCM用户
---

# alibabacloudstack_ascm_user

使用Provider配置的凭证在指定的资源集下编排ASCM用户。

## 示例用法

```
resource "alibabacloudstack_ascm_organization" "default" {
  name = "Dummy_Test_1"
}

resource "alibabacloudstack_ascm_user" "default" {
   cellphone_number = "892399537"
   email = "test@gmail.com"
   display_name = "C2C-DEL3"
   organization_id = alibabacloudstack_ascm_organization.default.org_id
   mobile_nation_code = "91"
   login_name = "C2C_alibabacloudstack_C2C"
}

output "org" {
  value = alibabacloudstack_ascm_user.default.*
}
```

## 参数参考

支持以下参数：

* `login_name` - (必填) 用户登录名。 
* `cellphone_number` - (必填) 用户的手机号码。
* `display_name` - (必填) 用户的显示名称。
* `email` - (必填) 用户的电子邮件地址。
* `mobile_nation_code` - (必填) 用户所属的移动国家代码。
* `organization_id` - (必填) 用户组织ID。
* `login_policy_id` - (可选) 用户登录策略ID。
* `role_ids` - 用户拥有的角色列表。
* `telephone_number` - (可选) 用户的电话号码。

## 属性参考

导出以下属性：

* `id` - 用户的登录名。
* `user_id` - 用户的ID。
* `init_password` - 用户的初始密码。
* `role_ids` - 用户拥有的角色列表。