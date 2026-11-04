---
subcategory: "企业控制台(ASCM)"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_account_ascm_custom_role_arn"
sidebar_current: "docs-Alibabacloudstack-datasource-account-ascm-custom-role-arn"
description: |-
    查询企业控制台为当前登录账号持有的指定角色物化的 RAM 角色 ARN
---

# alibabacloudstack_account_ascm_custom_role_arn

本数据源查询企业控制台（ASCM）为**当前凭证所属账号**持有的某个 ASCM 角色物化的 RAM 角色 ARN。

~> **说明：** ASCM 角色本身（`roleType = ROLETYPE_ASCM`，即 `alibabacloudstack_ascm_custom_role`、
`alibabacloudstack_ascm_ram_role` 创建出来的角色）是管控面角色，无法被 STS 扮演。对于每一个
「组织 + ASCM 角色 + ASCM 用户」的绑定关系，ASCM 会在 RAM/STS 后端维护一个独立角色，
只有该角色的 ARN 才能用于 `AssumeRole`，或用于日志库加密 `arn` 这类资源参数。物化出来的角色名遵循
`ascm-role-<组织 ID>-<ASCM 角色 ID>-<ASCM 用户 ID>` 命名规律。

~> **说明：** 组织和 ASCM 用户不是本数据源的入参。ASCM 后端只按调用者身份解析这两个值，并拒绝读取属于其他用户
的角色，因此本数据源回答的问题始终是「我登录用的这个账号，在该 ASCM 角色下对应哪个 RAM 角色」。若要查询其他
身份，请使用配置了该用户凭证的 alias provider。

## 示例用法

```
data "alibabacloudstack_account_ascm_custom_role_arn" "operation_admin" {
  ascm_role_id = 1
}

output "role_arn" {
  value = data.alibabacloudstack_account_ascm_custom_role_arn.operation_admin.arn
}
```

## 参数说明

以下参数受支持：

* `ascm_role_id` - (必填) ASCM 管控面侧角色的 ID。当前账号必须持有该角色，否则读取报 `ascm.auth.EntityNotExist.User.Role`。

## 属性说明

除了上述参数外，还导出以下属性：

* `id` - 物化出来的 RAM 角色的 ARN。
* `arn` - ASCM 物化出来的 RAM 角色 ARN，格式为 `acs:ram::<账号 ID>:role/<角色名>`，即传给 STS `AssumeRole` 的取值。
* `role_name` - 物化出来的 RAM 角色名称。
* `ram_role_id` - 物化出来的 RAM 角色在 RAM 侧的 ID。与 `ascm_role_id` 不同，后者是 ASCM 管控面侧的 ID。
* `assume_role_policy_document` - 物化出来的 RAM 角色的信任策略文档，JSON 格式。其中 `Principal` 决定了哪些账号或用户可以扮演该角色。
