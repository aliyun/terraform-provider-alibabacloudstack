---
subcategory: "RDS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_rds_backuppolicy"
sidebar_current: "docs-Alibabacloudstack-rds-backuppolicy"
description: |- 
  编排RDS备份策略资源
---

# alibabacloudstack_rds_backuppolicy
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_db_backup_policy`

使用Provider配置的凭证在指定的资源集编排RDS备份策略资源。

## 示例用法

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

## 参数参考

支持以下参数：

* `instance_id` - (必填, 变更时重建) - 要配置备份策略的 RDS 实例的 ID。
* `preferred_backup_time` - (选填) - 数据备份时间，格式：`HH:mmZ-HH:mmZ`(UTC 时间)。默认为 `"02:00Z-03:00Z"`。
* `backup_retention_period` - (选填) - 数据备份保留天数。取值范围：7~730。默认为 `7`。对于使用本地磁盘的 MySQL 实例，没有上限。
* `enable_backup_log` - (选填) - 是否开启日志备份：
  * **true**：表示已启用。
  * **false**：表示未启用。默认为 `true`。注意：“基础版”类别的 RDS 实例不支持设置日志备份。
* `log_backup_retention_period` - (选填) - 日志备份保留天数。取值范围：7~730，且不大于数据备份保留天数。仅适用于运行 MySQL 和 PostgreSQL 的实例。此参数仅在 `enable_backup_log` 设置为 `true` 时生效。
* `local_log_retention_hours` - (选填) - 日志备份本地保留小时数。取值范围：0~168(0~7*24)。默认为 `72`。
* `local_log_retention_space` - (选填) - 本地日志的最大循环空间使用率。如果最大循环空间使用率超过，则清除最早的 Binlog 直到空间使用率低于此比例。取值范围：0~50。默认为 `20`。
* `high_space_usage_protection` - (选填) - 如果实例使用空间大于 80% 或剩余空间小于 5GB，是否强制清理 Binlog：
  * **Enable**：清理。
  * **Disable**：不清除。默认为 `Enable`。
* `log_backup_frequency` - (选填) - 日志备份频率：
  * **LogInterval**：每 30 分钟备份一次。
  * 默认与数据备份周期 `preferred_backup_time` 一致。此参数仅适用于 SQL Server。
* `compress_type` - (选填) - 备份压缩方式：
  * **0**：无压缩。
  * **1**：zlib 压缩。
  * **2**：并行 zlib 压缩。
  * **4**：quicklz 压缩，启用数据库和表恢复。
  * **8**：MySQL8.0 quicklz 压缩但库表恢复不支持。默认为 `1`。
* `archive_backup_retention_period` - (选填) - 归档备份保留天数。取值范围：30~1095。默认为 `0`，表示未启用归档备份。仅在 `enable_backup_log` 设置为 `true` 且实例是 MySQL 本地磁盘时生效。
* `archive_backup_keep_count` - (选填) - 保留的归档备份数量：
  * 当 `archive_backup_keep_policy` 设置为 `ByMonth` 时，值为 `1` 到 `31`。
  * 当 `archive_backup_keep_policy` 设置为 `ByWeek` 时，值为 `1` 到 `7`。
  * 当 `archive_backup_keep_policy` 设置为 `KeepAll` 时，此参数不需要。默认为 `1`。
* `archive_backup_keep_policy` - (选填) - 归档备份的保留周期：
  * **ByMonth**：按月。
  * **ByWeek**：按周。
  * **KeepAll**：保留所有。默认为 `KeepAll`。
* `preferred_backup_period` - (选填) - 优选备份周期。

## 属性参考

除了上述参数外，还导出以下属性：

* `enable_backup_log` - 是否开启日志备份：
  * **true**：表示已启用。
  * **false**：表示未启用。
* `log_backup_retention_period` - 日志备份保留天数。取值范围：7~730，且不大于数据备份保留天数。仅在 `enable_backup_log` 设置为 `true` 时生效。
* `local_log_retention_hours` - 日志备份本地保留小时数。
* `local_log_retention_space` - 本地日志的最大循环空间使用率。如果最大循环空间使用率超过，则清除最早的 Binlog 直到空间使用率低于此比例。取值范围：0~50。默认为不修改。
* `log_backup_frequency` - 日志备份频率：
  * **LogInterval**：每 30 分钟备份一次。
  * 默认与数据备份周期 `preferred_backup_time` 一致。此参数仅适用于 SQL Server。
* `compress_type` - 备份压缩方式：
  * **0**：无压缩。
  * **1**：zlib 压缩。
  * **2**：并行 zlib 压缩。
  * **4**：quicklz 压缩，启用数据库和表恢复。
  * **8**：MySQL8.0 quicklz 压缩但库表恢复不支持。
* `archive_backup_retention_period` - 归档备份保留天数。默认值为 `0`，表示未启用归档备份。取值范围：30~1095。仅在 `enable_backup_log` 设置为 `true` 时生效。
* `archive_backup_keep_count` - 保留的归档备份数量。默认为 `1`。取值：
  * 当 `archive_backup_keep_policy` 设置为 `ByMonth` 时，值为 `1` 到 `31`。
  * 当 `archive_backup_keep_policy` 设置为 `ByWeek` 时，值为 `1` 到 `7`。
  * 当 `archive_backup_keep_policy` 设置为 `KeepAll` 时，此参数不需要。
* `archive_backup_keep_policy` - 归档备份的保留周期。保留在此周期内的备份数量由 `archive_backup_keep_count` 决定。默认为 `KeepAll`。取值：
  * **ByMonth**：按月。
  * **ByWeek**：按周。
  * **KeepAll**：保留所有。
* `preferred_backup_period` - 优选备份周期。