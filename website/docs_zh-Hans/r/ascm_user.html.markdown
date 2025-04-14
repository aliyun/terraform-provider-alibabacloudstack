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

## 参数说明

支持以下参数：

* `login_name` - (必填) 用户登录名。
* `cellphone_number` - (必填) 用户的手机号码。
* `display_name` - (必填) 用户的显示名称。
* `email` - (必填) 用户的电子邮件地址。
* `mobile_nation_code` - (必填) 用户所属的移动国家代码（例如：中国为“86”）。
* `organization_id` - (必填) 用户所属的组织ID。
* `login_policy_id` - (可选) 用户登录策略ID。用于定义用户的登录策略。
* `role_ids` - (可选) 用户拥有的角色ID列表。可以通过此参数为用户分配多个角色。
* `telephone_number` - (可选) 用户的固定电话号码。
* `init_password` - (可选) 用户的初始密码。如果未提供，系统将自动生成一个初始密码。

## 属性说明

导出以下属性：

* `id` - 用户的登录名（与`login_name`相同）。
* `user_id` - 用户的唯一标识符（UUID）。
* `init_password` - 用户的初始密码。如果在创建时未指定，则返回系统生成的初始密码。
* `role_ids` - 用户当前拥有的角色ID列表。
* `organization_id` - 用户所属的组织ID。注意：从1.0.32版本开始，该字段已被废弃，建议使用其他方式获取组织信息。