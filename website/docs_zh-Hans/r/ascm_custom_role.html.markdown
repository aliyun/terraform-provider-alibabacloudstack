---
subcategory: "企业控制台(ASCM)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_custom_role"
sidebar_current: "docs-Alibabacloudstack-resource-ascm-custom-role"
description: |-
  编排 ASCM 自定义角色资源。
---

# alibabacloudstack_ascm_custom_role

编排 ASCM 自定义角色资源。

-> **注意：** 该资源不支持更新操作。创建完成后，无法修改任何参数。更改任何参数都将导致资源重新创建。

## 示例用法

```
resource "alibabacloudstack_ascm_custom_role" "ramrole" {
  role_name = "Test_DELTA_Custom"
  description = "TestingComplete"
  organization_visibility = "organizationVisibility.global"
  role_range = "roleRange.allOrganizations"
  privileges = [
          "PRIG_SYS_BILLING_CLOUDPRODUCTBILL_READ",
          "PRIG_SYS_BILLING_ORGRSBILL_READ",
          "PRIG_SYS_BILLING_BILL_EXPORT",
          "PRIG_SYS_BILLING_BILL_MODIFY",
          "PRIG_SYS_CHANGEOWN_READ",
          "PRIG_SYS_CHANGEOWN_ORGANIZATION",
          "PRIG_SYS_CHANGEOWN_RESOURCESET",
          "PRIG_SYS_CHANGEOWN_USER",
          "PRIG_SYS_CHANGEOWN_RESOURCE",
          "PRIG_SYS_CHARGING_PRICE_READ",
          "PRIG_SYS_CHARGING_PRICE_OPERATE",
          "PRIG_SYS_CHARGING_PRICE_CREATE_DELETE",
          "PRIG_SYS_DOWNLOAD_CENTER_TASK_READ",
          "PRIG_SYS_DOWNLOAD_CENTER_TASK_CREATE",
          "PRIG_SYS_DOWNLOAD_CENTER_TASK_DELETE",
          "PRIG_SYS_DOWNLOAD_CENTER_REPORT_DOWNLOAD",
          "PRIG_SYS_LOGINPOLICY_READ",
          "PRIG_SYS_LOGINPOLICY_CREATE_DELETE",
          "PRIG_SYS_LOGINPOLICY_OPERATE",
          "PRIG_SYS_MENU_MANAGE",
          "PRIG_SYS_METERING_READ",
          "PRIG_SYS_METERING_EXPORT",
          "PRIG_SYS_MSGCENTER",
          "PRIG_SYS_OPLOG_READ",
          "PRIG_SYS_OPLOG_OPERATE",
          "PRIG_SYS_ORG_READ",
          "PRIG_SYS_ORG_CREATE_DELETE",
          "PRIG_SYS_ORG_OPERATE",
          "PRIG_SYS_ORG_AK_READ",
          "PRIG_SYS_QUOTA_READ",
          "PRIG_SYS_QUOTA_OPERATE",
          "PRIG_SYS_RESOURCESET_READ",
          "PRIG_SYS_RESOURCESET_CREATE_DELETE",
          "PRIG_SYS_RESOURCESET_OPERATE",
          "PRIG_SYS_ROLE_READ",
          "PRIG_SYS_ROLE_CREATE_DELETE",
          "PRIG_SYS_ROLE_OPERATE",
          "PRIG_SYS_SYSCONF",
          "PRIG_SYS_USER_READ",
          "PRIG_SYS_USER_CREATE_DELETE",
          "PRIG_SYS_USER_OPERATE",
          "PRIG_SYS_USERGROUP_READ",
          "PRIG_SYS_USERGROUP_CREATE_DELETE",
          "PRIG_SYS_USERGROUP_OPERATE"
          ]
}
output "Custom_role" {
  value = alibabacloudstack_ascm_custom_role.ramrole.*
}
```

## 参数说明

支持以下参数：

* `role_name` -（必填，变更后重建）自定义角色名称。长度介于 2 到 128 个字符之间。
* `organization_visibility` -（必填，变更后重建）自定义角色的组织可见性。有效值：
  - `organizationVisibility.organization`：仅当前组织可见。
  - `organizationVisibility.orgAndSubOrgs`：当前组织及其子组织可见。
  - `organizationVisibility.global`：全局可见。
* `role_range` -（必填，变更后重建）自定义角色的作用范围。有效值：
  - `roleRange.allOrganizations`：所有组织。
  - `roleRange.currentOrganization`：当前组织。
* `privileges` -（必填，变更后重建）分配给该自定义角色的权限列表。至少需要指定一个权限，每个权限以字符串形式表示。
* `description` -（可选，变更后重建）自定义角色的描述信息。

-> **注意：** 由于不支持更新操作，所有参数实际上都是 变更后重建。修改任何参数都会强制重新创建资源。

## 属性说明

导出以下属性：

* `id` - 自定义角色的唯一标识符，格式为 `<role_name>:<role_id>`。
* `role_id` - 自定义角色的内部 ID。
* `role_name` - 自定义角色的名称。

## Import

ASCM 自定义角色可以使用角色名称和角色 ID（以冒号分隔）进行导入，例如：

```
$ terraform import alibabacloudstack_ascm_custom_role.example my-custom-role:12345
```