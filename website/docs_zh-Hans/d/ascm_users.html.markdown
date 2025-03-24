---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_users"
sidebar_current: "docs-alibabacloudstack-datasource-ascm-users"
description: |-
    查询ascm用户
---

# alibabacloudstack_ascm_users

根据指定过滤条件列出当前凭证权限可以访问的用户列表。

## 示例用法

```
resource "alibabacloudstack_ascm_organization" "default" {
  name = "Dummy_Test_1"
}

resource "alibabacloudstack_ascm_user" "default" {
  cellphone_number = "899999537"
   email = "test@gmail.com"
   display_name = "C2C-DEL3"
   organization_id = alibabacloudstack_ascm_organization.default.org_id
   mobile_nation_code = "91"
   login_name = "C2C_alibabacloudstack_C2C"

}
output "org" {
  value = alibabacloudstack_ascm_user.default.*
}
data "alibabacloudstack_ascm_users" "users" {
 ids = [alibabacloudstack_ascm_user.user.user_id]
}
output "users" {
 value = data.alibabacloudstack_ascm_users.users.*
}
```

## 参数参考

以下是支持的参数：

* `ids` - (可选) 用户ID列表。
* `name_regex` - (可选) 用于通过用户登录名过滤结果的正则表达式字符串。
* `organization_id` - (可选) 通过指定的用户组织ID过滤结果。
* `login_name` - (可选) 用于过滤的用户登录名。
* `role_id` - (可选) 通过指定的用户角色ID过滤结果。
* `login_policy_id` - (可选) 通过指定的用户登录策略ID过滤结果。
* `current_page` - (可选) 分页的当前页码。
* `page_size` - (可选) 分页的每页项目数。
* `status` - (可选) 通过指定的用户状态过滤结果。
* `mark` - (可选) 用于过滤用户的标记。

## 属性参考

除了上述列出的参数外，还导出以下属性：

* `users` - 用户列表。每个元素包含以下属性：
  * `id` - 用户的ID。
  * `login_name` - 用户登录名。
  * `cell_phone_number` - 用户的手机号码。
  * `display_name` - 用户的显示名称。
  * `email` - 用户的电子邮件地址。
  * `mobile_nation_code` - 用户所属的移动国家代码。
  * `organization_id` - 用户的组织ID。
  * `login_policy_id` - 用户的登录策略ID。
  * `role_ids` - 用户拥有的角色列表。
  * `default_role_id` - 用户登录时使用的默认角色ID。
  * `cellphone_num` - 用户的手机号码。
  * `default` - 指示该用户是否为默认用户。
  * `deleted` - 指示该用户是否已被删除。
  * `enable_ding_talk` - 指示是否为用户启用了DingTalk通知。
  * `enable_email` - 指示是否为用户启用了电子邮件通知。
  * `enable_short_message` - 指示是否为用户启用了短消息通知。
  * `last_login_time` - 用户上次登录的时间戳。
  * `parent_pk` - 用户的父主键。
  * `primary_key` - 用户的主键。
  * `status` - 用户的状态。
* `role_ids` - 用户拥有的所有角色列表。
