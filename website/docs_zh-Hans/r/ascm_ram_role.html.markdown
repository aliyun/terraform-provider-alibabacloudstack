---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_ram_role"
sidebar_current: "docs-alibabacloudstack-resource-ascm-ram-role"
description: |-
  编排ASCM的RAM角色
---

# alibabacloudstack_ascm_ram_role

使用Provider配置的凭证在指定的资源集下编排ASCM的RAM角色。

## 示例用法

```
resource "alibabacloudstack_ascm_ram_role" "default" {
  role_name = "TestingRamRole"
  description = "TestingRam"
  organization_visibility = "organizationVisibility.global"
}
output "ramrole" {
  value = alibabacloudstack_ascm_ram_role.default.*
}
```

## 参数参考

以下参数是支持的：

* `role_name` - (必填) 角色名称。 
* `organization_visibility` - (必填) 组织可见性。有效值为 - "organizationVisibility.organization", "organizationVisibility.orgAndSubOrgs" 和 "organizationVisibility.global"。
* `description` - (可选) RAM角色的描述。注意 - 它不应包含任何空格。
* `role_range` - (必填) 角色管理的权限范围
  * `roleRange.orgAndSubOrgs` - 组织和级联的下属组织
  * `roleRange.allOrganizations` - 所有组织
  * `roleRange.userGroup` - 组织下的资源组集

## 属性参考

以下属性是导出的：

* `id` - 用户的RAM角色名称。
* `role_id` - RAM角色的ID。