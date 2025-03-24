---
subcategory: "DBS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dbs_backupplan"
sidebar_current: "docs-Alibabacloudstack-dbs-backupplan"
description: |- 
  Provides a dbs Backupplan resource.
---

# alibabacloudstack_dbs_backupplan
-> **NOTE:** Alias name has: `alibabacloudstack_dbs_backup_plan`

Provides a dbs Backupplan resource.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-testaccdbsbackupplan74295"
}

resource "alibabacloudstack_dbs_backup_plan" "default" {
  backup_method     = "logical"
  database_type     = "MySQL"
  instance_class    = "large"
  backup_plan_name  = var.name
  database_region   = "cn-hangzhou"
  storage_region    = "cn-hangzhou"
  instance_type     = "RDS"
  from_app          = "OpenAPI"
}
```

## Argument Reference

The following arguments are supported:

* `backup_plan_id` - (ForceNew) The unique identifier for the backup plan. This ID is automatically generated upon creation and cannot be modified afterward.
* `backup_method` - (Required) The backup method to be used. Valid values include:
  * `logical`: Logical backup.
  * `physical`: Physical backup.
* `database_type` - (Required) The type of database being backed up. Valid values include:
  * `MySQL`
  * `MSSQL`
  * `Oracle`
  * `MongoDB`
  * `Redis`
* `instance_class` - (Required) The class of the backup instance. Valid values include:
  * `small`: Small instance.
  * `large`: Large instance.
* `backup_plan_name` - (Optional) The name of the backup plan. If not specified, Terraform will auto-generate a name.
* `database_region` - (Optional) The region where the source database resides. For example, `cn-hangzhou`.
* `storage_region` - (Optional) The region where the backup data will be stored. It can be the same or different from the `database_region`.
* `instance_type` - (Optional) The type of the database instance. Valid values include:
  * `RDS`: Relational Database Service.
  * `PolarDB`: Polar Database.
  * `DDS`: Document Database Service.
  * `Kvstore`: Key-Value Store.
  * `Other`: Other types.
* `from_app` - (Optional) Indicates the source of the request. The default value is `OpenAPI`. Manual setting is generally unnecessary.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `backup_plan_id` - The unique identifier for the backup plan.
* `backup_plan_name` - The name of the backup plan.


### Explanation of Changes Made:

1. **Example Usage**: Added variables like `backup_plan_name`, `database_region`, `storage_region`, `instance_type`, and `from_app` to make the example more comprehensive and align with the detailed fields provided in the second document.

2. **Argument Reference**:
   - Expanded descriptions for each argument to clarify their purpose and valid values.
   - Included all optional fields (`database_region`, `storage_region`, `instance_type`, `from_app`) with proper explanations.
   - Provided valid options for `backup_method`, `database_type`, `instance_class`, and `instance_type`.

3. **Attributes Reference**: Clearly stated that `backup_plan_id` and `backup_plan_name` are exported attributes in addition to the arguments listed above.