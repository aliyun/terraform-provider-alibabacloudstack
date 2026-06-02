---
subcategory: "企业控制台(ASCM)"
layout: "alibabacloudstack"
page_title: "阿里云专有云: alibabacloudstack_ascm_ram_service_role"
sidebar_current: "docs-Alibabacloudstack-resource-ascm-ram-service-role"
description: |-
  提供AscmRAM服务角色。
---

# alibabacloudstack_ascm_ram_service_role

提供AscmRAM服务角色。

~> **注意:** 此资源也可以使用以下别名引用：
- `alibabacloudstack_ascm_service_ram_role`

## 使用示例

```
resource "alibabacloudstack_ascm_organization" "default" {
  name = "Tf-testingresource-org2"
  parent_id = "1"
} 
 resource "alibabacloudstack_ascm_ram_service_role" "default" {
  organization_id = "${alibabacloudstack_ascm_organization.default.id}"
  product_name = "ECS"
}
```

## 参数引用

以下参数被支持：

* `product_name` - (必选, 变更后重建) 产品名称。
* `organization_id` - (必选, 变更后重建) 组织ID。

## 属性引用

以下属性被导出：

* `id` - 服务RAM角色的ID。
* `ram_roles` - RAM角色列表。
  * `id` - RAM角色的ID。
  * `product_name` - 产品名称。
  * `organization_id` - 组织ID。
  * `arn` - RAM角色的ARN。
  * `role_id` - 角色ID。
  * `role_name` - 角色名称。
  * `role_type` - 角色类型。
  * `region` - RAM角色所在区域。
  * `description` - RAM角色描述。
  * `aliyun_user_id` - 阿里云用户ID。
  * `assume_role_policy_document` - 角色策略文档。
  * `organization_name` - 组织名称。
  * `policies` - 附加到RAM角色的策略列表。
    * `policy_id` - 策略ID。
    * `region` - 策略所在区域。
    * `policy_name` - 策略名称。
    * `description` - 策略描述。
    * `policy_document` - 策略文档。
    * `policy_type` - 策略类型。
    * `default_version` - 策略默认版本。
    * `aliyun_user_id` - 阿里云用户ID。
    * `ram_group_id` - RAM用户组ID。
    * `ascm_ram_policy_id` - ASCM RAM策略ID。
    * `attach_date` - 策略附加日期。
    * `resource_set_id` - 资源集ID。
    * `privilege_id` - 权限ID。
    * `ram_role_id` - RAM角色ID。

## 导入

ASCM 服务 RAM 角色可以使用 organization_id 和 product_name 以 `<organization_id>:<product_name>` 格式导入，例如：

```
$ terraform import alibabacloudstack_ascm_ram_service_role.example 1:ECS
```