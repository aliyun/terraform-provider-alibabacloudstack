---
subcategory: "AnalyticDB for PostgreSQL"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_gpdb_backup_policy"
sidebar_current: "docs-alibabacloudstack-resource-gpdb-backup_policy"
description: |-
  Configures the backup policy for an AnalyticDB PostgreSQL instance.
---

# alibabacloudstack_gpdb_backup_policy

Configures the backup policy for an AnalyticDB PostgreSQL instance.

## Example Usage

### Basic Usage

```hcl
variable "name" {
  default = "tftest2982"
}

data "alibabacloudstack_zones" "gpdb" {
  available_resource_creation = "Gpdb"
}

data "alibabacloudstack_gpdb_instance_types" "default" {
  engine_version = "6.0"
}

data "alibabacloudstack_gpdb_instances" "default" {
  ids = [""]
}

resource "alibabacloudstack_gpdb_instance" "default" {
  count                    = length(data.alibabacloudstack_gpdb_instances.default.ids) > 0 ? 0 : 1
  engine                   = "gpdb"
  engine_version           = data.alibabacloudstack_gpdb_instance_types.default.instance_types.0.engine_version
  instance_class           = data.alibabacloudstack_gpdb_instance_types.default.instance_types.0.id
  db_instance_mode         = data.alibabacloudstack_gpdb_instance_types.default.instance_types.0.db_instance_mode
  db_instance_storage_type = "local_ssd"
  availability_zone        = data.alibabacloudstack_zones.gpdb.zones.0.id
  description              = var.name
  seg_node_num             = "2"
  cpu_type                 = "Intel"
}

locals {
  gpdb_instance_id = length(data.alibabacloudstack_gpdb_instances.default.ids) > 0 ? data.alibabacloudstack_gpdb_instances.default.ids.0 : alibabacloudstack_gpdb_instance.default.0.id
}

resource "alibabacloudstack_gpdb_backup_policy" "default" {
  preferred_backup_time   = "02:00Z-03:00Z"
  preferred_backup_period = "Monday,Wednesday,Friday"
  backup_retention_period = "7"
  enable_recovery_point   = "false"
  db_instance_id          = local.gpdb_instance_id
}
```

## Argument Reference

The following arguments are supported:

* `db_instance_id` - (Required, Forces new resource) The ID of the instance.

* `backup_retention_period` - (Optional) The retention period of data backups in days. Default value: 7. Maximum value: 7. Valid values: 1 to 7.
* `enable_recovery_point` - (Optional) Specifies whether to enable automatic recovery points. Valid values:
  * `true`: Enabled.
  * `false`: Disabled.
* `preferred_backup_period` - (Optional) The backup period. Multiple values must be separated by commas (,). Valid values:
  * `Monday`: Monday.
  * `Tuesday`: Tuesday.
  * `Wednesday`: Wednesday.
  * `Thursday`: Thursday.
  * `Friday`: Friday.
  * `Saturday`: Saturday.
  * `Sunday`: Sunday.
* `preferred_backup_time` - (Optional) The backup time. Format: HH:mmZ-HH:mmZ (UTC time).
* `recovery_point_period` - (Optional) The frequency of recovery points. Valid values:
  * `1`: Hourly.
  * `2`: Every two hours.
  * `4`: Every four hours.
  * `8`: Every eight hours.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID, which is the instance ID.
* `backup_retention_period` - The retention period of data backups in days.
* `enable_recovery_point` - Indicates whether automatic recovery points are enabled.
* `preferred_backup_period` - The backup period.
* `preferred_backup_time` - The backup time.
* `recovery_point_period` - The frequency of recovery points.

## Import

GPDB Backup Policy can be imported using the DBInstanceId, e.g.

```
$ terraform import alibabacloudstack_gpdb_backup_policy.example pgm-xxxxxxxxx
```