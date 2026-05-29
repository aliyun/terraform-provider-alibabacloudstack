---
subcategory: "企业控制台(ASCM)"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_ascm_roles"
sidebar_current: "docs-alibabacloudstack-datasource-ascm-roles"
description: |-
    查询角色列表

> **注意：** 此数据源也可以使用以下别名引用：
> - `alibabacloudstack_ascm_ram_roles`

---

# alibabacloudstack_ascm_roles

根据指定过滤条件列出当前凭证权限可以访问的的角色列表。

## 示例用法

``` hcl
resource "alibabacloudstack_ascm_ram_role" "default" {
  role_name = "DELTA1"
  description = "Testing Complete"
  organization_visibility = "organizationVisibility.global"
}

data "alibabacloudstack_ascm_roles" "default" {
  id = alibabacloudstack_ascm_ram_role.default.role_id
  name_regex = alibabacloudstack_ascm_ram_role.default.role_name
  role_type = "ROLETYPE_RAM"
}

output "roles" {
  value = data.alibabacloudstack_ascm_roles.default.*
}
```

## 参数说明

支持以下参数：

* `id` - （可选，已废弃）用于通过角色 ID 过滤结果。该字段已废弃，将在 3.21.0 版本中移除，请使用 `ids` 代替。
* `ids` - （可选）角色 ID 列表。用于通过角色 ID 过滤结果。
* `name_regex` - （可选）用于通过角色名称过滤结果的正则表达式字符串。
* `description` - （可选）角色的描述信息。
* `role_type` - （可选）角色类型，例如 `"ROLETYPE_RAM"` 表示 RAM 类型的角色。
* `output_file` - （可选，已废弃）数据源的输出文件。该字段已废弃，将在 3.19.0 版本中移除。如需将内容写入文件，请使用 `local_file` 提供者。

## 属性说明

除了上述列出的参数外，还导出以下属性：

* `roles` - 角色列表。每个元素包含以下属性：
    * `id` - 角色的ID。
    * `name` - 角色名称。
    * `description` - 关于角色的描述。
    * `role_level` - 角色级别，表示角色的优先级或范围。
    * `role_type` - 角色类型，例如 `"ROLETYPE_RAM"`。
    * `ram_role` - 是否为 RAM 授权角色（布尔值）。
    * `role_range` - 角色的具体范围，例如全局范围或特定组织范围。
    * `user_count` - 使用该角色的用户数量。
    * `enable` - 是否启用该角色（布尔值）。
    * `default` - 是否为默认角色（布尔值）。
    * `active` - 角色的状态，表示是否处于活动状态（布尔值）。
    * `owner_organization_id` - 角色所属的组织所有者的ID。
    * `code` - 角色代码，通常用于标识角色的唯一性。
    * `assume_role_policy_document` - RAM 授权角色的信任策略文档。
    * `organization_visibility` - 角色的组织可见性范围。

## 导入

此数据源为只读数据源，不支持导入。
