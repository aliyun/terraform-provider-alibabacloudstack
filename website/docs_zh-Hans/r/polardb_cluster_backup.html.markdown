---
subcategory: "云原生数据库 PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_cluster_backup"
sidebar_current: "docs-Alibabacloudstack-resource-polardb-cluster-backup"
description: |-
  提供 PolarDB 集群版手动备份资源。
---

# alibabacloudstack_polardb_cluster_backup

为 PolarDB **集群版**实例触发一次手动全量备份，并等待备份集生成完成。

集群版与主备版在 POP 网关上是两个不同的产品：本资源调用 `polardb / 2017-08-01`，以 `DBClusterId`
定位实例；而 [`alibabacloudstack_polardb_backup`](/docs_zh-Hans/r/polardb_backup.html.markdown)
调用 `polardb / 2024-01-30`，以 `DBInstanceId` 定位实例。`pc-` 开头的集群实例请使用本资源。

## 示例用法

```hcl
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details              = true
}

resource "alibabacloudstack_polardb_cluster_instance" "default" {
  db_cluster_description   = "tf-example"
  db_type                  = "PostgreSQL"
  db_version               = "16"
  db_node_class            = "polar.pg.x4.small.c"
  zone_id                  = "${data.alibabacloudstack_zones.default.zones.0.id}"
  vswitch_id               = "${alibabacloudstack_vpc_vswitch.default.id}"
  storage_type             = "ESSDPL1"
  storage_space            = 20
  sub_category             = "normal_exclusive"
}

resource "alibabacloudstack_polardb_cluster_backup" "default" {
  db_cluster_id = "${alibabacloudstack_polardb_cluster_instance.default.id}"
}
```

## 参数参考

以下参数支持设置：

  * `db_cluster_id` - （必填，变更需重建）PolarDB 集群 ID。

集群版的 `CreateBackup` 不接受其他入参，备份方式、备份模式与备份类型均由服务返回，体现在属性中。

## 属性参考

除以上参数外，还导出以下属性：

  * `backup_job_id` - `CreateBackup` 返回的备份任务 ID。它与 `backup_id` 不属于同一套 ID 体系，无法用于查询备份集。
  * `backup_id` - 备份集 ID，与 `db_cluster_id` 共同组成资源 ID。
  * `backup_method` - 备份方式，例如 `Physical`。
  * `backup_mode` - 备份模式，由本资源触发的备份为 `Manual`。
  * `backup_type` - 备份类型，例如 `FullBackup`。
  * `backup_status` - 备份状态，备份作业完成后为 `Success`。
  * `backups_level` - 备份集级别，例如 `Level-1`。
  * `store_status` - 备份集是否可用于恢复，例如 `Enabled`。
  * `backup_set_size` - 备份集大小，单位为字节。
  * `is_avail` - 备份集是否可用。
  * `consistent_time` - 备份集的一致性时间点，Unix 秒级时间戳。
  * `backup_start_time` - 备份开始时间。
  * `backup_end_time` - 备份结束时间。

## 超时配置

`timeouts` 块支持以下操作：

  * `create` - （默认 20 分钟）用于触发备份作业并等待备份集生成。

## 导入

PolarDB 集群备份可使用 `db_cluster_id:backup_id` 格式（以冒号分隔）导入：

```
$ terraform import alibabacloudstack_polardb_cluster_backup.default pc-abc0abc0abc0abc00:1793345712
```

## 注意事项

  * 销毁本资源不会删除备份集。备份集的生命周期由集群备份策略管理，参见
    [`alibabacloudstack_polardb_cluster_backup_policy`](/docs_zh-Hans/r/polardb_cluster_backup_policy.html.markdown)。
  * `backup_job_id` 仅对由 Terraform 创建的备份有效，导入的备份集该属性为空。
