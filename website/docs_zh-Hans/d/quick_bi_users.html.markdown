---
subcategory: "Quick BI"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_quick_bi_users"
sidebar_current: "docs-alibabacloudstack-datasource-quick-bi-users"
description: |-
  查询Quick BI 用户列表。
---

# alibabacloudstack_quick_bi_users

根据指定过滤条件列出当前凭证权限可以访问的Quick BI用户列表。


## 示例用法

### 基础用法

```terraform
data "alibabacloudstack_quick_bi_users" "ids" {
  ids = ["example_id"]
}
output "quick_bi_user_id_1" {
  value = data.alibabacloudstack_quick_bi_users.ids.users.0.id
}
```

## 参数说明

以下参数受支持：

* `enable_details` - (可选) 默认为 `false`。将其设置为 `true` 可以输出更多关于资源属性的详细信息。
* `ids` - (可选，强制更新) 用户 ID 列表。
* `keyword` - (可选，强制更新) 组织成员昵称或用户名的关键词。

## 属性说明

除了上述参数外，还导出以下属性：

* `users` - Quick BI 用户列表。每个元素包含以下属性：
    * `account_id` - 阿里云账户 ID。
    * `account_name` - 阿里云账户名称。
    * `admin_user` - 是否为管理员。有效值：`true` 和 `false`。
    * `auth_admin_user` - 是否为权限管理员。有效值：`true` 和 `false`。
    * `email` - 用户的电子邮件地址。
    * `id` - 用户的 ID。
    * `nick_name` - 用户的昵称。
    * `phone` - 用户的电话号码。
    * `user_id` - 用户的 ID。
    * `user_type` - 组织成员的角色类型。有效值：`Analyst`（分析师）、`Developer`（开发者）和 `Visitor`（访客）。