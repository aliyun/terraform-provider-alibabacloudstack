---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_custom_role"
sidebar_current: "docs-alibabacloudstack-resource-ascm-custom-role"
description: |-
  编排Ascm自定义角色。
---

# alibabacloudstack_ascm_custom_role

使用Provider配置的凭证在指定的资源集下编排Ascm自定义角色。

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

## 参数参考

支持以下参数：

* `role_name` - (必填) 自定义角色名称。
* `organization_visibility` - (必填) 组织可见性。有效值为 - "organizationVisibility.organization", "organizationVisibility.orgAndSubOrgs" 和 "organizationVisibility.global"。
* `description` - (可选) 自定义角色的描述。注意 - 它不应包含任何空格。
* `role_range` - (必填) 自定义角色的角色范围。
* `privileges` - (必填) 分配给该角色的权限。

## 属性参考

导出以下属性：

* `id` - 自定义角色名称和用户ID。
* `role_id` - 自定义角色的ID。