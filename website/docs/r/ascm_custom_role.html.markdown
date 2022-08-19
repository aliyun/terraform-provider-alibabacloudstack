---
subcategory: "ASCM"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_ascm_custom_role"
sidebar_current: "docs-apsarastack-resource-ascm-custom-role"
description: |-
  Provides Ascm custom role.
---

# apsarastack\_ascm_custom_role

Provides Ascm custom role.

## Example Usage

```
resource "apsarastack_ascm_custom_role" "ramrole" {
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
  value = apsarastack_ascm_custom_role.ramrole.*
}
```
## Argument Reference

The following arguments are supported:

* `role_name` - (Required) Custom Role name. 
* `organization_visibility` - (Required) organization visibility. Valid Values are - "organizationVisibility.organization", "organizationVisibility.orgAndSubOrgs" and "organizationVisibility.global".
* `description` - (Optional) Description for the custom role. Note - It should not contain any spaces.
* `role_range` - (Required) Role Range for the custom role.
* `privileges` - (Required) Privileges assign to that role. 

## Attributes Reference

The following attributes are exported:

* `id` - Custom Role Name and ID of the user.
* `role_id` - The ID of the custom role.