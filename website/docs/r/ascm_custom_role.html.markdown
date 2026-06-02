---
subcategory: "Apsara Stack Cloud Management"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_custom_role"
description: |-
  Provides an ASCM custom role resource.
---

# alibabacloudstack_ascm_custom_role

Provides an ASCM custom role resource.

-> **NOTE:** The Update operation is not implemented for this resource. Once created, the resource cannot be modified. Changing any arguments will not trigger an update.

## Example Usage

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
## Argument Reference

The following arguments are supported:

* `role_name` - (Required, ForceNew) The name of the custom role. The length is between 2 and 128 characters.
* `organization_visibility` - (Required, ForceNew) The organization visibility of the custom role. Valid values: `organizationVisibility.organization`, `organizationVisibility.orgAndSubOrgs`, `organizationVisibility.global`.
* `role_range` - (Required, ForceNew) The range of the custom role. Valid values: `roleRange.allOrganizations`, `roleRange.currentOrganization`.
* `privileges` - (Required, ForceNew) A list of privileges assigned to the custom role. At least one privilege must be specified. Each privilege is represented as a string.
* `description` - (Optional, ForceNew) The description of the custom role.

-> **NOTE:** Since the Update operation is not supported, all arguments are effectively ForceNew. Modifying any argument will force the resource to be recreated.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the custom role. The format is `<role_name>:<role_id>`.
* `role_id` - The internal ID of the custom role.
* `role_name` - The name of the custom role.

## Import

ASCM custom role can be imported using the role name and role ID separated by a colon, e.g.

```
$ terraform import alibabacloudstack_ascm_custom_role.example my-custom-role:12345
```