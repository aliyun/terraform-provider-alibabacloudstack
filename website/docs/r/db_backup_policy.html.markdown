---
subcategory: "RDS"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_db_backup_policy"
sidebar_current: "docs-apsarastack-resource-db-backup-policy"
description: |-
  Provides an RDS backup policy resource.
---

# apsarastack\_db\_backup\_policy

Provides an RDS instance backup policy resource and used to configure instance backup policy.

-> **NOTE:** Each DB instance has a backup policy and it will be set default values when destroying the resource.

## Example Usage

```
variable "creation" {
  default = "Rds"
}

variable "name" {
  default = "dbbackuppolicybasic"
}

data "apsarastack_zones" "default" {
  available_resource_creation = "${var.creation}"
}

resource "apsarastack_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}

resource "apsarastack_vswitch" "default" {
  vpc_id            = "${apsarastack_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
  name              = "${var.name}"
}

resource "apsarastack_db_instance" "instance" {
  engine           = "MySQL"
  engine_version   = "5.6"
  instance_type    = "rds.mysql.s1.small"
  instance_storage = "10"
  vswitch_id       = "${apsarastack_vswitch.default.id}"
  instance_name    = "${var.name}"
}

resource "apsarastack_db_backup_policy" "policy" {
  instance_id = "${apsarastack_db_instance.instance.id}"
} 
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The Id of instance that can run database.
* `preferred_backup_period` - (Optional) DB Instance backup period. Please set at least two days to ensure backing up at least twice a week. Valid values: [Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday]. Default to ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"].
* `preferred_backup_time` - (Optional) DB instance backup time, in the format of HH:mmZ- HH:mmZ. Time setting interval is one hour. Default to "02:00Z-03:00Z". China time is 8 hours behind it.
* `backup_retention_period` - (Optional, available in 1.69.0+) Instance backup retention days. Valid values: [7-730]. Default to 7. But mysql local disk is unlimited.
* `enable_backup_log` - (Optional, available in 1.68.0+) Whether to backup instance log. Valid values are `true`, `false`, Default to `true`. Note: The 'Basic Edition' category Rds instance does not support setting log backup. [What is Basic Edition](https://www.alibabacloud.com/help/doc-detail/48980.htm).
* `log_backup_retention_period` - (Optional, available in 1.69.0+) Instance log backup retention days. Valid when the `enable_backup_log` is `1`. Valid values: [7-730]. Default to 7. It cannot be larger than `backup_retention_period`.
* `local_log_retention_hours` - (Optional, available in 1.69.0+) Instance log backup local retention hours. Valid when the `enable_backup_log` is `true`. Valid values: [0-7*24].
* `local_log_retention_space` - (Optional, available in 1.69.0+) Instance log backup local retention space. Valid when the `enable_backup_log` is `true`. Valid values: [5-50].
* `high_space_usage_protection` - (Optional, available in 1.69.0+) Instance high space usage protection policy. Valid when the `enable_backup_log` is `true`. Valid values are `Enable`, `Disable`.
* `log_backup_frequency` - (Optional, available in 1.69.0+) Instance log backup frequency. Valid when the instance engine is `SQLServer`. Valid values are `LogInterval`.
* `compress_type` - (Optional, available in 1.69.0+) The compress type of instance policy. Valid values are `1`, `4`, `8`.
* `archive_backup_retention_period` - (Optional, available in 1.69.0+) Instance archive backup retention days. Valid when the `enable_backup_log` is `true` and instance is mysql local disk. Valid values: [30-1095], and `archive_backup_retention_period` must larger than `backup_retention_period` 730.
* `archive_backup_keep_count` - (Optional, available in 1.69.0+) Instance archive backup keep count. Valid when the `enable_backup_log` is `true` and instance is mysql local disk. When `archive_backup_keep_policy` is `ByMonth` Valid values: [1-31]. When `archive_backup_keep_policy` is `ByWeek` Valid values: [1-7].
* `archive_backup_keep_policy` - (Optional, available in 1.69.0+) Instance archive backup keep policy. Valid when the `enable_backup_log` is `true` and instance is mysql local disk. Valid values are `ByMonth`, `Disable`, `KeepAll`.

-> **NOTE:** Currently, the SQLServer instance does not support to modify `log_backup_retention_period`.

## Attributes Reference

The following attributes are exported:

* `id` - The current backup policy resource ID. It is same as 'instance_id'.
