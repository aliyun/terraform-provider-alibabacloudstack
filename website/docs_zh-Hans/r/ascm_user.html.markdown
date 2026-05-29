---
subcategory: "应用"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_user"
sidebar_current: "docs-alibabacloudstack-resource-ascm-user"
description: |-
  编排 ASCM 用户资源
---

# alibabacloudstack_ascm_user

使用 Provider 配置的凭证在指定的资源集下编排 ASCM 用户。

-> **注意:** 该资源同时支持以下别名：`apsarastack_ascm_user`。

## 示例用法

```terraform
resource "alibabacloudstack_ascm_logon_policy" "default" {
  name        = "default-policy"
  description = "默认登录策略"
  rule        = "ALLOW"
}

resource "alibabacloudstack_ascm_user" "default" {
  login_name         = "example-user"
  display_name       = "示例用户"
  email              = "example@example.com"
  cellphone_number   = "13800138000"
  mobile_nation_code = "86"
  login_policy_id    = alibabacloudstack_ascm_logon_policy.default.policy_id
  role_ids           = ["1"]
}
```

## 参数说明

支持以下参数：

* `login_name` - (必填，ForceNew) 用户登录名。该参数唯一标识一个用户。修改此参数会强制重新创建资源。
* `display_name` - (必填) 用户的显示名称。
* `email` - (必填) 用户的电子邮件地址。
* `cellphone_number` - (必填) 用户的手机号码。
* `mobile_nation_code` - (必填) 用户所属的移动国家代码，例如中国为 `86`。
* `login_policy_id` - (必填) 与用户关联的登录策略 ID。
* `role_ids` - (可选) 要分配给用户的角色 ID 列表。如果指定，则必须至少包含一个角色 ID。
* `telephone_number` - (可选) 用户的固定电话号码。
* `organization_id` - (已废弃，ForceNew) 用户所属的组织 ID。该字段自 provider 版本 1.0.32 起已废弃。用户将创建在 Provider 配置的组织下。

## 属性说明

导出以下属性：

* `id` - 资源 ID，即用户的登录名。
* `user_id` - 用户的内部 ID。
* `user_uid` - 用户的唯一主键标识（UID）。
* `init_password` - 为用户生成的初始密码。仅在创建后可获取，无法手动设置。
* `role_ids` - 分配给用户的角色 ID 列表。
* `organization_id` - 用户所属的组织 ID。自 provider 版本 1.0.32 起已废弃。

## Import

ASCM 用户可以通过 `login_name` 进行导入，例如：

```shell
$ terraform import alibabacloudstack_ascm_user.example example-user
```