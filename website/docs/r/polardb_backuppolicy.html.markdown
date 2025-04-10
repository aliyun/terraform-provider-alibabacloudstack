---
subcategory: "PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_backuppolicy"
sidebar_current: "docs-Alibabacloudstack-polardb-backuppolicy"
description: |-
  Provides a polardb Backuppolicy resource.
---

# alibabacloudstack_polardb_backuppolicy

Provides a polardb Backuppolicy resource.

## Example Usage
```
resource "alibabacloudstack_polardb_backuppolicy" "example" {
  db_instance_id             = "your_db_instance_id"
  backup_policy_mode         = "DataBackupPolicy"
  backup_retention_period    = 30
  compress_type              = 1
  enable_backup_log          = 1
  high_space_usage_protection = "Enable"
  local_log_retention_hours  = 24
  local_log_retention_space  = 20
  log_backup_retention_period = 30
  preferred_backup_period    = "Monday,Wednesday,Friday"
  preferred_backup_time      = "02:00Z-03:00Z"
  released_keep_policy       = "Lastest"
}
```

## Argument Reference

The following arguments are supported:
  * `backup_log` - (Optional) - Log backup switch. Value: **Enable | Disabled **
  * `backup_policy_mode` - (Optional) - Backup type:* **DataBackupPolicy**: Data backup* **LogBackupPolicy**: log backup
  * `backup_retention_period` - (Optional) - Data backup retention days, value: 7~730.DescriptionWhen BackupPolicyMode is databackuppolicymode, this parameter must be set.Takes effect only when the BackupPolicyMode parameter is databackuppolicyy.
  * `compress_type` - (Optional) - Backup compression method, value:* **0**: No compression* **1**:zlib compression* **2**: parallel zlib compression* **4**:quicklz compression, enabling database and table recovery* **8**:MySQL8.0 quicklz compression but library table recovery is not supported.
  * `db_instance_id` - (Required) - zhe PolarDB instance ID
  * `enable_backup_log` - (Optional) - Whether to enable log Backup. Valid values:* **1**: indicates enabled* **0**: indicates closed
  * `high_space_usage_protection` - (Optional) - If the instance usage space is greater than 80% or the remaining space is less than 5GB, whether to force Binlog cleanup:* **Disable**: do not clean up* **Enable**: clean up
  * `local_log_retention_hours` - (Optional) - Log backup local retention hours.
  * `local_log_retention_space` - (Optional) - The maximum circular space usage of local logs. If the maximum circular space usage is exceeded, the earliest Binlog is cleared until the space usage is lower than this ratio. Value: 0~50. Default is not modified.DescriptionWhen BackupPolicyMode is set to logbackuppolicymode, this parameter must be passed.Takes effect only when the BackupPolicyMode parameter is set to LogBackupPolicy.
  * `log_backup_frequency` - (Optional) - Log backup frequency, value:* **LogInterval**: Back up every 30 minutes;* The default data backup cycle is the same as the data backup cycle **PreferredBackupPeriod.> parameter **LogBackupFrequency** is only applicable to SQL Server.
  * `log_backup_local_retention_number` - (Optional) - The number of local binlogs retained. The default is 60. Value: 6~100.DescriptionTakes effect only when the BackupPolicyMode parameter is set to LogBackupPolicy.If the instance type is MySQL, you can pass in **-1**, that is, the number of reserved local binlogs is not limited.
  * `log_backup_retention_period` - (Optional) - The number of days for which the log backup is retained. Valid values: 7 to 730. The log backup retention period cannot be longer than the data backup retention period.NoteIf you enable the log backup feature, you can specify the log backup retention period. This parameter is supported for instances that run MySQL and PostgreSQL.This parameter takes effect only when BackupPolicyMode is set to DataBackupPolicy or LogBackupPolicy.
  * `preferred_backup_period` - (Optional) - Data backup cycle. Separate multiple values with commas (,). Valid values:* **Monday**: Monday* **Tuesday**: Tuesday* **Wednesday**: Wednesday* **Thursday**: Thursday* **Friday**: Friday* **Saturday**: Saturday* **Sunday**: Sunday
  * `preferred_backup_time` - (Optional) - Data backup time, format: <I> HH:mm</I> Z-<I> HH:mm</I> Z(UTC time).
  * `released_keep_policy` - (Optional) - Archive backup retention policy for deleted instances. Value:* **None**: not reserved* **Lastest**: Keep the last one* **All**: All reserved

## Attributes Reference

The following attributes are exported in addition to the arguments listed above: 

  * `backup_log` - (Computed) Log backup switch. Value: **Enable | Disabled **