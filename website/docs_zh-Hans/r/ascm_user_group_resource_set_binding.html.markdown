---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_user_group_resource_set_binding"
sidebar_current: "docs-alibabacloudstack-resource-ascm-user-group-resource-set-binding"
description: |-
  编排绑定ASCM用户组和资源集
---

# alibabacloudstack_ascm_user_group_resource_set_binding

使用Provider配置的凭证在指定的资源集下编排绑定ASCM用户组和资源集。

## 示例用法

```
resource "alibabacloudstack_ascm_organization" "default" {
 name = "Test_binder"
 parent_id = "1"
}

resource "alibabacloudstack_ascm_user_group" "default" {
 group_name =      "%s"
 organization_id = alibabacloudstack_ascm_organization.default.org_id
}


resource "alibabacloudstack_ascm_resource_group" "default" {
  organization_id = alibabacloudstack_ascm_organization.default.org_id
  name = "alibabacloudstack-terraform-resourceGroup"
}

resource "alibabacloudstack_ascm_user_group_resource_set_binding" "default" {
  resource_set_id = alibabacloudstack_ascm_resource_group.default.rg_id
  user_group_id = alibabacloudstack_ascm_user_group.default.user_group_id
  ascm_role_id = "2"
}

output "binder" {
  value = alibabacloudstack_ascm_user_group_resource_set_binding.default.*
}
```

## 参数参考

支持以下参数：

* `resource_set_id` - (必填) 资源集ID列表。
* `user_group_id` - (必填) 用户组ID。
* `ascm_role_id` - (可选) ASCM角色ID。

## 属性参考

导出以下属性：

* `resource_set_id` - 资源集ID列表。
* `user_group_id` - 用户组ID。