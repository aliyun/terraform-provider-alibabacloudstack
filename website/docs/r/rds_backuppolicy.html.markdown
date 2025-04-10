---
subcategory: "RDS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_rds_backuppolicy"
sidebar_current: "docs-Alibabacloudstack-rds-backuppolicy"
description: |- 
  Provides a rds Backuppolicy resource.
---

# alibabacloudstack_rds_backuppolicy
-> **NOTE:** Alias name has: `alibabacloudstack_db_backup_policy`

Provides a rds Backuppolicy resource.

## Example Usage

```hcl
variable "creation" {
  default = "Rds"
}

variable "name" {
  default = "dbbackuppolicybasic"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "${var.creation}"
}

resource "alibabacloudstack_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = "${alibabacloudstack_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  name              = "${var.name}"
}

resource "alibabacloudstack_db_instance" "instance" {
  engine           = "MySQL"
  engine_version   = "5.6"
  instance_type    = "rds.mysql.s1.small"
  instance_storage = "10"
  vswitch_id       = "${alibabacloudstack_vswitch.default.id}"
  instance_name    = "${var.name}"
}

resource "alibabacloudstack_rds_backuppolicy" "policy" {
  instance_id                  = "${alibabacloudstack_db_instance.instance.id}"
  preferred_backup_time        = "02:00Z-03:00Z"
  backup_retention_period      = 14
  enable_backup_log            = true
  log_backup_retention_period  = 14
  local_log_retention_hours    = 72
  local_log_retention_space    = 20
  high_space_usage_protection  = "Enable"
  log_backup_frequency         = "LogInterval"
  compress_type                = 1
  archive_backup_retention_period = 90
  archive_backup_keep_count    = 5
  archive_backup_keep_policy    = "ByMonth"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the RDS instance for which the backup policy is being configured.
* `preferred_backup_time` - (Optional) Data backup time, in the format of `<I>HH:mm</I>Z-<I>HH:mm</I>Z` (UTC time). Default to `"02:00Z-03:00Z"`.
* `backup_retention_period` - (Optional) Data backup retention days. Valid values: [7-730]. Default to `7`. For MySQL instances with local disks, there is no upper limit.
* `enable_backup_log` - (Optional) Whether to enable log backup. Valid values:
  * **true**: Indicates enabled.
  * **false**: Indicates disabled. Default to `true`. Note: The 'Basic Edition' category RDS instance does not support setting log backup.
* `log_backup_retention_period` - (Optional) The number of days for which the log backup is retained. Valid values: [7-730]. Default to `7`. It cannot be larger than `backup_retention_period`. This parameter is supported for instances that run MySQL and PostgreSQL. This parameter takes effect only when `enable_backup_log` is set to `true`.
* `local_log_retention_hours` - (Optional) Log backup local retention hours. Valid values: [0-168] (0~7*24). Default is `72`.
* `local_log_retention_space` - (Optional) The maximum circular space usage of local logs. If the maximum circular space usage is exceeded, the earliest Binlog is cleared until the space usage is lower than this ratio. Valid values: [0-50]. Default is `20`.
* `high_space_usage_protection` - (Optional) If the instance usage space is greater than 80% or the remaining space is less than 5GB, whether to force Binlog cleanup:
  * **Enable**: Clean up.
  * **Disable**: Do not clean up. Default to `Enable`.
* `log_backup_frequency` - (Optional) Log backup frequency. Valid values:
  * **LogInterval**: Backup every 30 minutes.
  * Default data backup cycle is the same as the data backup cycle `preferred_backup_time`. This parameter is only applicable to SQL Server.
* `compress_type` - (Optional) Backup compression method. Valid values:
  * **0**: No compression.
  * **1**: zlib compression.
  * **2**: Parallel zlib compression.
  * **4**: quicklz compression, enabling database and table recovery.
  * **8**: MySQL8.0 quicklz compression but library table recovery is not supported. Default to `1`.
* `archive_backup_retention_period` - (Optional) Number of days to keep archived backups. Valid values: [30-1095]. Default to `0`, indicating that the archive backup is not enabled. Takes effect only when `enable_backup_log` is set to `true` and the instance is a MySQL local disk.
* `archive_backup_keep_count` - (Optional) Number of archived backups retained. Valid values:
  * When `archive_backup_keep_policy` is set to `ByMonth`, the value is from `1` to `31`.
  * When `archive_backup_keep_policy` is set to `ByWeek`, the value is from `1` to `7`.
  * When `archive_backup_keep_policy` is set to `KeepAll`, this parameter is not required. Default to `1`.
* `archive_backup_keep_policy` - (Optional) The retention period for archived backups. Valid values:
  * **ByMonth**: Month.
  * **ByWeek**: Week.
  * **KeepAll**: Keep all. Default to `KeepAll`.
* `preferred_backup_period` - (Optional) Preferred backup period.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `enable_backup_log` - Whether to enable log backup. Valid values:
  * **true**: Indicates enabled.
  * **false**: Indicates disabled.
* `log_backup_retention_period` - The number of days for which the log backup is retained. Valid values: [7-730]. The log backup retention period cannot be longer than the data backup retention period. This parameter takes effect only when `enable_backup_log` is set to `true`.
* `local_log_retention_hours` - Log backup local retention hours.
* `local_log_retention_space` - The maximum circular space usage of local logs. If the maximum circular space usage is exceeded, the earliest Binlog is cleared until the space usage is lower than this ratio. Value: [0-50]. Default is not modified.
* `log_backup_frequency` - Log backup frequency. Valid values:
  * **LogInterval**: Back up every 30 minutes.
  * The default data backup cycle is the same as the data backup cycle `preferred_backup_time`. This parameter is only applicable to SQL Server.
* `compress_type` - Backup compression method. Valid values:
  * **0**: No compression.
  * **1**: zlib compression.
  * **2**: Parallel zlib compression.
  * **4**: quicklz compression, enabling database and table recovery.
  * **8**: MySQL8.0 quicklz compression but library table recovery is not supported.
* `archive_backup_retention_period` - Number of days to keep archived backups. The default value is `0`, indicating that the archive backup is not enabled. Value: [30-1095]. Takes effect only when `enable_backup_log` is set to `true`.
* `archive_backup_keep_count` - Number of archived backups retained. Default is `1`. Value:
  * When `archive_backup_keep_policy` is set to `ByMonth`, the value is from `1` to `31`.
  * When `archive_backup_keep_policy` is set to `ByWeek`, the value is from `1` to `7`.
  * When `archive_backup_keep_policy` is set to `KeepAll`, this parameter is not required.
* `archive_backup_keep_policy` - The retention period for archived backups. The number of backups that can be saved in this period is determined by `archive_backup_keep_count`. Default is `KeepAll`. Valid values:
  * **ByMonth**: Month.
  * **ByWeek**: Week.
  * **KeepAll**: Keep all.
* `preferred_backup_period` - Preferred backup period.