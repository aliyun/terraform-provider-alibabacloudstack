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

## 参数参考

支持以下参数：

* `account_id` - (可选，变更时重建) 阿里巴巴云账户 ID。 
* `account_name` - (必填) 一个阿里巴巴云账户，阿里巴巴云名称。
* `admin_user` - (必填) 是否为管理员。有效值：`true` 和 `false`。
* `auth_admin_user` - (必填) 此用户是否为权限管理员。有效值：`false`，`true`。
* `nick_name` - (必填，变更时重建) 用户的昵称。
* `user_type` - (必填) 组织成员的类型角色分别。有效值：`Analyst`，`Developer` 和 `Visitor`。

## 属性参考

导出以下属性：

* `id` - Terraform 中的 User 资源 ID。

## 导入

Quick BI 用户可以使用 id 导入，例如

```bash
$ terraform import alibabacloudstack_quick_bi_user.example <id>
```