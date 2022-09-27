---
subcategory: "Database Backup(DBS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dbs_backup_plan"
sidebar_current: "docs-alibabacloudstack-resource-dbs-backup-plan"
description: |-
  Provides a Alibabacloudstack DBS Backup Plan resource.
---

# alibabacloudstack\_dbs\_backup\_plan

Provides a DBS Backup Plan resource.

For information about DBS Backup Plan and how to use it, see [What is Backup Plan](https://help.aliyun.com/apsara/enterprise/v_3_16_2_20220708/cbs/enterprise-developer-guide/CreateBackupPlan.html?spm=a2c4g.14484438.10004.1#doc-api-Dbs-CreateBackupPlan).

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-testacccn-wulan-env180-d01dbsbackupplan45232"
}


resource "alibabacloudstack_dbs_backup_plan" "default" {
  backup_method = "logical"
  database_type = "MySQL"
  instance_class = "large"
}
```

## Argument Reference

The following arguments are supported:

* `backup_method` - (Required, ForceNew) Backup method. Valid values: `logical`, `physical`.
* `backup_plan_name` - (Required, ForceNew) The name of the resource.
* `database_type` - (Required, ForceNew) Database type. Valid values: `MySQL`, `MSSQL`, `Oracle`, `MongoDB`, `Redis`.
* `instance_class` - (Required, ForceNew) The instance class. Valid values: `large`, `small`.
* `instance_type` - (Optional) The instance type. Valid values: `RDS`, `PolarDB`, `DDS`, `Kvstore`, `Other`.
* `database_region` - (Optional) The region of the database.
* `storage_region` - (Optional) The storage region.
* `from_app` - (Optional) It is used to remark the request source. The default value is OpenAPI, and manual setting is unnecessary.

## Attributes Reference

The following attributes are exported:

* `backup_plan_id` - The resource ID in terraform of Backup Plan.

## Import

DBS Backup Plan can be imported using the id, e.g.

```
$ terraform import alibabacloudstack_dbs_backup_plan.example <id>
```