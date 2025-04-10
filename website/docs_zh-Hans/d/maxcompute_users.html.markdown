---
subcategory: "MaxCompute"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_maxcompute_users"
sidebar_current: "docs-alibabacloudstack-datasource-maxcompute-users"
description: |-
  查询Max Compute 用户
---

# alibabacloudstack_maxcompute_users

根据指定过滤条件列出当前凭证权限可以访问的Max Compute 用户列表。[什么是用户](https://www.alibabacloud.com/help/en/maxcompute/latest/users)

## 示例用法

```hcl
data "alibabacloudstack_maxcompute_users" "example" {
  name_regex = "example-user"
}

output "users" {
  value = data.alibabacloudstack_maxcompute_users.example.users
}
```

## 参数说明
支持以下参数：

* `ids` - (可选) 用于过滤结果的用户 ID 列表。
* `name_regex` - (可选) 按名称过滤用户的正则表达式模式。
* `user_name` - (必填) 用户的名称。


## 属性说明
导出以下属性：

* `ids` - 用户的唯一标识符列表。
* `users` - 用户列表。每个元素包含以下属性：
    * `id` - 用户的唯一标识符。
    * `user_id` - 用户 ID。
    * `user_pk` - 用户的主要键。
    * `user_name` - 用户的名称。
    * `user_type` - 用户类型。
    * `organization_id` - 组织 ID。
    * `organization_name` - 组织名称。
    * `description` - 用户的描述。