---
subcategory: "自助式 BI（商业智能）工具"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_quick_bi_user_group_user"
sidebar_current: "docs-Alibabacloudstack-resource-quick-bi-user-group-user"
description: |-
  编排Quick BI 用户组成员。
---

# alibabacloudstack_quick_bi_user_group_user

使用Provider配置的凭证在指定的资源集编排Quick BI 用户组成员。


## 示例用法

### 基础用法

```terraform

resource "alibabacloudstack_quick_bi_user_group_user" "example" {
  user_group_id = alibabacloudstack_quick_bi_user_group.default.id
  account_id    = alibabacloudstack_quick_bi_user.default.id
}

resource "alibabacloudstack_quick_bi_user" "default" {
  nick_name       = "example_value"
  account_name    = "example_value"
  admin_user      = "false"
  auth_admin_user = "false"
  user_type       = "Developer"
}

resource "alibabacloudstack_quick_bi_user_group" "default" {
  user_group_name        = "example_value"
  user_group_description = "example_value"
}

```

## 参数说明

支持以下参数：

* `user_group_id` - (必填) 用户组 ID。修改此参数会强制重新创建资源。
* `account_id` - (必填) Quick BI 用户 ID。修改此参数会强制重新创建资源。

## 属性说明

导出以下属性：

* `id` - 资源 ID，格式为 `<user_group_id>:<account_id>`。

## 导入

Quick BI 用户组成员可以使用 user_group_id 和 account_id 组合导入，例如

```bash
$ terraform import alibabacloudstack_quick_bi_user_group_user.example <user_group_id>:<account_id>
```
