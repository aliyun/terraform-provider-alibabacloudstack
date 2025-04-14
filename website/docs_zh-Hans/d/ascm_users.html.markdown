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

## 参数说明

以下是支持的参数：

* `ids` - (可选) 用户ID列表，用于精确匹配用户。
* `name_regex` - (可选) 一个正则表达式字符串，用于通过用户登录名过滤结果。
* `organization_id` - (可选) 通过指定的用户组织ID过滤结果，帮助定位特定组织内的用户。
* `login_name` - (可选) 用户登录名，用于过滤特定登录名的用户。
* `role_id` - (可选) 通过指定的用户角色ID过滤结果，筛选具有特定角色的用户。
* `login_policy_id` - (可选) 通过指定的用户登录策略ID过滤结果，筛选使用特定登录策略的用户。
* `current_page` - (可选) 分页的当前页码，默认为第一页，用于分页查询。
* `page_size` - (可选) 每页显示的项目数，默认值通常为10或20，具体取决于系统配置。
* `status` - (可选) 通过指定的用户状态过滤结果，例如“启用”或“禁用”。
* `mark` - (可选) 用于过滤用户的标记，可以是自定义标签或其他标识符。

## 属性说明

除了上述列出的参数外，还导出以下属性：

* `users` - 用户列表。每个元素包含以下属性：
  * `id` - 用户的唯一标识符，用于区分不同用户。
  * `login_name` - 用户登录时使用的名称，通常是用户名。
  * `cell_phone_number` - 用户的手机号码，用于接收通知或验证。
  * `display_name` - 用户的显示名称，通常用于界面展示。
  * `email` - 用户的电子邮件地址，用于接收系统通知或重置密码。
  * `mobile_nation_code` - 用户所属国家/地区的移动国家代码，例如中国的代码为“86”。
  * `organization_id` - 用户所属的组织ID，用于组织级别的权限管理。
  * `login_policy_id` - 用户应用的登录策略ID，控制登录行为和规则。
  * `role_ids` - 用户拥有的角色列表，表示用户在系统中的权限范围。
  * `default_role_id` - 用户登录时默认使用的角色ID，决定初始权限。
  * `cellphone_num` - 用户的手机号码，与`cell_phone_number`类似，可能是冗余字段。
  * `default` - 布尔值，指示该用户是否为默认用户（例如管理员账户）。
  * `deleted` - 布尔值，指示该用户是否已被删除，逻辑删除而非物理删除。
  * `enable_ding_talk` - 布尔值，指示是否为用户启用了DingTalk通知功能。
  * `enable_email` - 布尔值，指示是否为用户启用了电子邮件通知功能。
  * `enable_short_message` - 布尔值，指示是否为用户启用了短信通知功能。
  * `last_login_time` - 时间戳，记录用户上次登录的时间，用于审计或安全分析。
  * `parent_pk` - 用户的父主键，可能用于层级结构的用户管理。
  * `primary_key` - 用户的主键，唯一标识用户记录。
  * `status` - 用户的状态，例如“启用”、“禁用”或“锁定”。
* `role_ids` - 用户拥有的所有角色列表，提供更全面的角色信息，便于权限管理。