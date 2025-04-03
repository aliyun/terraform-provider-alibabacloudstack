---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_ascm_roles"
sidebar_current: "docs-alibabacloudstack-datasource-ascm-roles"
description: |-
    查询角色列表

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

## 参数参考

支持以下参数：

* `id` - (可选) 用于通过角色ID过滤结果。
* `name_regex` - (可选) 用于通过角色名称过滤结果的正则表达式字符串。

## 属性参考

除了上述列出的参数外，还导出以下属性：

* `roles` - 角色列表。每个元素包含以下属性：
    * `id` - 角色的ID。
    * `name` - 角色名称。
    * `description` - 关于角色的描述。
    * `role_level` - 角色级别。
    * `role_type` - 角色类型。
    * `ram_role` - RAM授权角色。
    * `role_range` - 角色的具体范围。
    * `user_count` - 用户数量。
    * `enable` - 是否启用。
    * `default` - 默认角色。
    * `active` - 角色状态。
    * `owner_organization_id` - 角色所属的组织所有者的ID。
    * `code` - 角色代码。