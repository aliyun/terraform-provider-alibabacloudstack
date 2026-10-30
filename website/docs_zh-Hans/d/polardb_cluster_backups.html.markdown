---
subcategory: "云原生数据库 PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_cluster_backups"
sidebar_current: "docs-Alibabacloudstack-datasource-polardb-cluster-backups"
description: |-
  查询 PolarDB 集群版实例的备份集列表。
---

# alibabacloudstack_polardb_cluster_backups

查询 PolarDB **集群版**实例的备份集列表。

集群版的 `DescribeBackups` 强制要求同时提供 `StartTime` 与 `EndTime`，因此未指定时默认查询最近
7 天，且默认窗口的右端向后补一天，避免刚完成数秒的备份集因分钟截断而漏查。时间窗口两端均使用
`2006-01-02T15:04Z` 格式，带秒的值会被服务拒绝。`end_time` 可以是未来时间，而未来的 `start_time`
会以 `InvalidStartTime.Malformed` 被拒绝；起始晚于结束的窗口不会匹配任何备份集。

## 示例用法

```hcl
data "alibabacloudstack_polardb_cluster_backups" "default" {
  db_cluster_id = "${alibabacloudstack_polardb_cluster_instance.default.id}"
}

```

## 参数参考

以下参数支持设置：

  * `ids` - （可选，变更需重建，计算）按备份集 ID 列表过滤查询结果。
  * `db_cluster_id` - （必填）PolarDB 集群 ID。
  * `start_time` - （可选）查询时间窗口的起始时间。
  * `end_time` - （可选）查询时间窗口的结束时间。

## 属性参考

除以上参数外，还导出以下属性：

  * `ids` - 查询到的备份集 ID 列表。
  * `backups` - 备份集列表，每个元素包含以下属性：
    * `id` - 备份集 ID。
    * `backup_id` - 备份集 ID。
    * `db_cluster_id` - 备份集所属的 PolarDB 集群 ID。
    * `backup_method` - 备份方式，例如 `Physical`。
    * `backup_mode` - 备份模式，`Automated` 或 `Manual`。
    * `backup_type` - 备份类型，例如 `FullBackup`。
    * `backup_status` - 备份状态，备份作业完成后为 `Success`。
    * `backups_level` - 备份集级别，例如 `Level-1`。
    * `store_status` - 备份集是否可用于恢复，例如 `Enabled`。
    * `backup_set_size` - 备份集大小，单位为字节。
    * `is_avail` - 备份集是否可用。
    * `consistent_time` - 备份集的一致性时间点，Unix 秒级时间戳。
    * `backup_start_time` - 备份开始时间。
    * `backup_end_time` - 备份结束时间。
