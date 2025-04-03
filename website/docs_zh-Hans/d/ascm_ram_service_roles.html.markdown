---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_ram_service_roles"
sidebar_current: "docs-alibabacloudstack-datasource-ascm-ram-service-roles"
description: |-
    查询RAM服务角色列表。
---

# alibabacloudstack_ascm_ram_service_roles

根据指定过滤条件列出当前凭证权限可以访问的RAM角色列表。

## 示例用法

```
data "alibabacloudstack_ascm_ram_service_roles" "role" {
  product = "ECS"
}
output "role" {
  value = data.alibabacloudstack_ascm_ram_service_roles.role.*
}
```

## 参数参考

支持以下参数：

* `ids` - (可选) RAM 角色 ID 列表。
* `product` - (可选) 按其产品过滤结果的正则表达式字符串。有效值 - "ECS"。
* `description` - (可选) 关于 RAM 角色的描述。

## 属性参考

除了上述列出的参数外，还导出以下属性：

* `roles` - 角色列表。每个元素包含以下属性：
    * `id` - 角色的 ID。
    * `name` - 角色名称。
    * `description` - 关于角色的描述。
    * `role_type` - 角色类型。
    * `product` - 角色类型。
    * `organization_name` - 组织名称。
    * `aliyun_user_id` - 阿里云用户 ID。