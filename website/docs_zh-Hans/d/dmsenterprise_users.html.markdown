---
subcategory: "DMSEnterprise"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dmsenterprise_users"
sidebar_current: "docs-Alibabacloudstack-datasource-dmsenterprise-users"
description: |- 
  查询企业版数据库管理用户列表。
---

# alibabacloudstack_dmsenterprise_users
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_dms_enterprise_users`

根据指定过滤条件列出当前凭证权限可以访问的企业版数据库管理用户列表。

## 示例用法

```terraform
# 创建组织
resource "alibabacloudstack_ascm_organization" "default" {
  name = "Test_binder"
  parent_id = "1"
}

# 创建用户
resource "alibabacloudstack_ascm_user" "user" {
  cellphone_number = "13900000000"
  email = "test@gmail.com"
  display_name = "C2C-DELTA"
  organization_id = alibabacloudstack_ascm_organization.default.org_id
  mobile_nation_code = "91"
  login_name = "tf_testAccDmsEnterpriseUsersDataSource_3194885"
  login_policy_id = 1
}

# 创建企业版数据库管理用户
resource "alibabacloudstack_dms_enterprise_user" "default" {
  uid = alibabacloudstack_ascm_user.user.user_id
  user_name = alibabacloudstack_ascm_user.user.login_name
  mobile = "15910799999"
  role_names = ["DBA"]
}

# 查询企业版数据库管理用户
data "alibabacloudstack_dms_enterprise_users" "default" {
  ids = [alibabacloudstack_dms_enterprise_user.default.uid]
  name_regex = "user-.*"
  role = "USER"
  status = "NORMAL"
  tid = "1234567890"
  output_file = "users_output.txt"
}

output "first_user_id" {
  value = data.alibabacloudstack_dms_enterprise_users.default.users.0.id
}
```

## 参数参考

以下参数是支持的：

* `name_regex` - (选填, 变更时重建) 用于通过企业版数据库管理用户的昵称过滤结果的正则表达式字符串。
* `role` - (选填, 变更时重建) 要查询的用户的角色。例如：`USER`表示普通用户。
* `search_key` - (选填, 变更时重建) 用于查询用户的关键词。可以是用户名、昵称或其他相关信息。
* `status` - (选填, 变更时重建) 用户的状态。例如：`NORMAL`表示正常状态，`DISABLED`表示禁用状态。
* `tid` - (选填, 变更时重建) 企业版数据库管理中的租户ID。这是系统右上角显示的租户ID。更多信息，请参见[查看租户信息](~~181330~~)。
* `ids` - (选填, 变更时重建) 企业版数据库管理用户ID（UID）列表。可以通过此参数指定需要查询的具体用户ID。

## 属性说明

除了上述参数外，还导出以下属性：

* `ids` - 企业版数据库管理用户ID（UID）列表。此列表包含所有匹配条件的用户ID。
* `names` - 企业版数据库管理用户名字列表。此列表包含所有匹配条件的用户昵称。
* `users` - 企业版数据库管理用户列表。每个元素包含以下属性：
  * `mobile` - 用户的钉钉号或手机号。用于标识用户的联系方式。
  * `nick_name` - 用户的昵称。通常用于显示用户的友好名称。
  * `user_name` - 用户的登录名或昵称。与`nick_name`类似，但可能用于不同的场景。
  * `parent_uid` - 如果用户对应于资源访问管理（RAM）用户，则为父账户的阿里云唯一ID（UID）。此字段用于标识RAM用户的关系。
  * `role_ids` - 用户扮演的角色ID列表。每个角色ID对应一个特定的角色权限。
  * `role_names` - 用户扮演的角色名称列表。每个角色名称描述了用户在系统中的角色。
  * `status` - 用户的状态。例如：`NORMAL`表示正常状态，`DISABLED`表示禁用状态。
  * `id` - 用户的阿里云唯一ID（UID）。此字段是用户的唯一标识符。
  * `uid` - `id`的别名。可以用于替代`id`字段。
  * `user_id` - 用户的ID。此字段可能与`id`或`uid`相同，具体取决于系统实现。
