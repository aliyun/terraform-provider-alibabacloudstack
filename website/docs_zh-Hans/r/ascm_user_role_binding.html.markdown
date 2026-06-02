---
subcategory: "企业控制台(ASCM)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_user_role_binding"
sidebar_current: "docs-Alibabacloudstack-resource-ascm-user-role-binding"
description: |-
  编排绑定ASCM用户和角色
---

# alibabacloudstack_ascm_user_role_binding

-> **弃用说明：** 该资源在未来版本中可能会被移除。`ascm_user` 资源已包含对应的功能，建议优先使用 `ascm_user`。

提供 ASCM 用户角色绑定资源。

## 示例用法

### 使用 `role_id`（推荐）

```hcl
resource "alibabacloudstack_ascm_user" "default" {
  cellphone_number   = "13900000000"
  email              = "test@example.com"
  display_name       = "测试用户"
  organization_id    = data.alibabacloudstack_account.current.organization_id
  mobile_nation_code = "91"
  login_name         = "testUser"
  login_policy_id    = 1
}

resource "alibabacloudstack_ascm_user_role_binding" "default" {
  role_id    = 5
  login_name = alibabacloudstack_ascm_user.default.login_name
}
```

### 使用 `role_ids`（自 v3.21.0 起废弃）

```hcl
resource "alibabacloudstack_ascm_user_role_binding" "default" {
  role_ids   = [5, 6, 7]
  login_name = alibabacloudstack_ascm_user.default.login_name
}
```

## 参数说明

以下参数为支持的配置项：

* `login_name` - （必填，变更时强制重建）要绑定角色的用户登录名。修改此参数会强制重新创建资源。
* `role_id` - （可选，变更时强制重建）要绑定到用户的单个角色 ID。该参数与 `role_ids` 互斥。修改此参数会强制重新创建资源。

-> **注意：** 使用 `role_id` 时，资源 ID 格式为 `login_name:role_id`。

* `role_ids` - （可选，Computed，已废弃）要绑定到用户的角色 ID 集合。自版本 3.21.0 起，该参数不再支持配置，建议使用 `role_id` 替代。该参数与 `role_id` 互斥。角色 ID 数量最大为 10 个。

-> **注意：** 使用 `role_ids` 时，资源 ID 格式为 `login_name`。使用 `role_ids` 创建的资源将不支持删除操作。

## 属性说明

以下属性为资源创建后导出的内容：

* `id` - 资源 ID。格式取决于使用 `role_id` 还是 `role_ids` 创建。
* `login_name` - 用户的登录名。
* `role_ids` - 已绑定到用户的角色 ID 列表。

## Import

ASCM 用户角色绑定可以通过资源 ID 导入。ID 格式取决于资源的创建方式：

- 使用 `role_id` 创建时：`login_name:role_id`，例如：

```shell
$ terraform import alibabacloudstack_ascm_user_role_binding.example testUser:5
```

- 使用 `role_ids` 创建时（已废弃）：`login_name`，例如：

```shell
$ terraform import alibabacloudstack_ascm_user_role_binding.example testUser
```
