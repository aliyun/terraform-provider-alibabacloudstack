---
subcategory: "PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_backuppolicy"
sidebar_current: "docs-Alibabacloudstack-polardb-backuppolicy"
description: |- 
  编排polardb备份规则
---

# alibabacloudstack_polardb_backuppolicy

使用Provider配置的凭证在指定的资源集编排polardb备份规则。

## 示例用法
```hcl
resource "alibabacloudstack_polardb_backuppolicy" "example" {
  db_instance_id             = "your_db_instance_id"
  backup_policy_mode         = "DataBackupPolicy"
  backup_retention_period    = 30
  compress_type              = 1
  enable_backup_log          = 1
  high_space_usage_protection = "Enable"
  local_log_retention_hours  = 24
  local_log_retention_space  = 20
  log_backup_frequency       = "LogInterval"
  log_backup_local_retention_number = 60
  log_backup_retention_period = 30
  preferred_backup_period    = "Monday,Wednesday,Friday"
  preferred_backup_time      = "02:00Z-03:00Z"
  released_keep_policy       = "Lastest"
}
```

## 参数说明

支持以下参数：
  * `backup_log` - (选填) - 日志备份开关，取值：**Enable | Disabled**
  * `backup_policy_mode` - (选填) - 备份类型：
    * **DataBackupPolicy**：数据备份
    * **LogBackupPolicy**：日志备份
  * `backup_retention_period` - (选填) - 数据备份保留天数，取值范围：7~730。说明：当 `backup_policy_mode` 为 `DataBackupPolicy` 时，该参数必传。仅在 `backup_policy_mode` 参数为 `DataBackupPolicy` 时生效。
  * `compress_type` - (选填) - 备份压缩方式，取值：
    * **0**：不压缩
    * **1**：zlib 压缩
    * **2**：并行 zlib 压缩
    * **4**：quicklz 压缩，开启了库表恢复
    * **8**：MySQL8.0 quicklz 压缩但尚未支持库表恢复
  * `db_instance_id` - (必填) - PolarDB 实例 ID。
  * `enable_backup_log` - (选填) - 是否开启日志备份，取值：
    * **1**：表示开启
    * **0**：表示关闭
  * `high_space_usage_protection` - (选填) - 实例使用空间大于80%，或者剩余空间小于5GB时，是否强制清理 Binlog：
    * **Disable**：不清理
    * **Enable**：清理
  * `local_log_retention_hours` - (选填) - 日志备份本地保留小时数。
  * `local_log_retention_space` - (选填) - 本地日志最大循环空间使用率，超出后，则从最早的 Binlog 开始清理，直到空间使用率低于该比例。取值：0~50。默认不修改。说明：当 `backup_policy_mode` 为 `LogBackupPolicy` 时，该参数必传。仅在 `backup_policy_mode` 参数为 `LogBackupPolicy` 时生效。
  * `log_backup_frequency` - (选填) - 日志备份频率，取值：
    * **LogInterval**：每30分钟备份一次；
    * 默认与数据备份周期 `PreferredBackupPeriod` 一致。> 参数 `LogBackupFrequency` 仅适用于 SQL Server。
  * `log_backup_local_retention_number` - (选填) - 本地 Binlog 保留个数。默认为 60。取值：6~100。说明：仅在 `backup_policy_mode` 参数为 `LogBackupPolicy` 时生效。如果实例类型为 MySQL，可以传入 **-1**，即不限制本地 Binlog 的保留个数。
  * `log_backup_retention_period` - (选填) - 日志备份保留天数。取值：7~730，且不大于数据备份保留天数。说明：当开启日志备份时，可设置日志备份文件的保留天数，目前仅支持 MySQL 和 PostgreSQL 实例设置该值。`backup_policy_mode` 参数为 `DataBackupPolicy` 或 `LogBackupPolicy` 时都适用。
  * `preferred_backup_period` - (选填) - 数据备份周期，多个取值用英文逗号(,)隔开，取值：
    * **Monday**：周一
    * **Tuesday**：周二
    * **Wednesday**：周三
    * **Thursday**：周四
    * **Friday**：周五
    * **Saturday**：周六
    * **Sunday**：周日
  * `preferred_backup_time` - (选填) - 数据备份时间，格式：<i>HH:mm</i>Z-<i>HH:mm</i>Z(UTC 时间)。
  * `released_keep_policy` - (选填) - 已删除实例的归档备份保留策略。取值：
    * **None**：不保留
    * **Lastest**：保留最后一个
    * **All**：全部保留

## 属性说明

除了上述所有参数外，还导出了以下属性：
  * `id` - 资源的唯一标识符。
  * `status` - 当前备份策略的状态。
  * `last_backup_time` - 最近一次备份的时间戳。
  * `next_backup_time` - 下一次计划备份的时间戳。