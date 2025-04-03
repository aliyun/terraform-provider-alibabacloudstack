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

* `name_regex` - (选填, 变更时重建） 用于通过企业版数据库管理用户的昵称过滤结果的正则表达式字符串。
* `role` - (选填, 变更时重建） 要查询的用户的角色。
* `search_key` - (选填, 变更时重建） 用于查询用户的关键词。
* `status` - (选填, 变更时重建） 用户的状态。例如：`NORMAL`表示正常状态。
* `tid` - (选填, 变更时重建） 企业版数据库管理中的租户ID。这是系统右上角显示的租户ID。更多信息，请参见[查看租户信息](~~181330~~)。
* `ids` - (选填, 变更时重建） 企业版数据库管理用户ID（UID）列表。

## 属性参考

除了上述参数外，还导出以下属性：

* `ids` - 企业版数据库管理用户ID（UID）列表。
* `names` - 企业版数据库管理用户名字列表。
* `users` - 企业版数据库管理用户列表。每个元素包含以下属性：
  * `mobile` - 用户的钉钉号或手机号。
  * `nick_name` - 用户的昵称。
  * `user_name` - 用户的昵称。
  * `parent_uid` - 如果用户对应于资源访问管理（RAM）用户，则为父账户的阿里云唯一ID（UID）。
  * `role_ids` - 用户扮演的角色ID列表。
  * `role_names` - 用户扮演的角色名称列表。
  * `status` - 用户的状态。
  * `id` - 用户的阿里云唯一ID（UID）。
  * `uid` - `id`的别名。
  * `user_id` - 用户的ID。