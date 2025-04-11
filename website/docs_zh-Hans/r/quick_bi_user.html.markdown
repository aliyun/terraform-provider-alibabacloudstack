---
subcategory: "Quick BI"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_quick_bi_user"
sidebar_current: "docs-alibabacloudstack-resource-quick-bi-user"
description: |-
  编排Quick BI 用户
---

# alibabacloudstack_quick_bi_user

使用Provider配置的凭证在指定的资源集编排Quick BI 用户。

有关 Quick BI 用户及其使用方法的信息，请参见 [什么是用户](https://www.alibabacloud.com/help/doc-detail/33813.htm)。

## 示例用法

### 基础用法

```terraform
resource "alibabacloudstack_quick_bi_user" "example" {
  account_name    = "example_value"
  admin_user      = false
  auth_admin_user = false
  nick_name       = "example_value"
  user_type       = "Analyst"
}

```

## 参数说明

支持以下参数：

* `account_id` - (可选，ForceNew) 阿里巴巴云账户 ID。如果未提供，则使用当前账户。
* `account_name` - (必填) 阿里巴巴云账户名称。这是用户的登录名。
* `admin_user` - (必填) 是否为管理员用户。有效值：`true` 和 `false`。如果设置为 `true`，则该用户将拥有管理员权限。
* `auth_admin_user` - (必填) 是否为权限管理员用户。有效值：`true` 和 `false`。如果设置为 `true`，则该用户将能够管理权限。
* `nick_name` - (必填，ForceNew) 用户的昵称。这是显示在系统中的用户名。
* `user_type` - (必填) 用户的角色类型。有效值：
  * `Analyst` - 分析师角色，主要用于数据分析和报表制作。
  * `Developer` - 开发者角色，主要用于开发和集成。
  * `Visitor` - 访客角色，仅具有查看权限。

## 属性说明

导出以下属性：

* `id` - Terraform 中的 User 资源 ID。此 ID 是用户在系统中的唯一标识符。

## 导入

Quick BI 用户可以使用 id 导入，例如

```bash
$ terraform import alibabacloudstack_quick_bi_user.example <id>
```