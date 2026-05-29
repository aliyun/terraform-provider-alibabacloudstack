---
subcategory: "云原生数据库 PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_backup"
sidebar_current: "docs-Alibabacloudstack-polardb-backup"
description: |-
  Provides a polardb Backup resource.
---

# alibabacloudstack\_polardb\_backup

Provides a polardb Backup resource.

## 示例用法
```
variable "name" {
  default = "tf-testAccPolardbBackupBasic_55"
}

data "alibabacloudstack_zones" default {
  available_resource_creation = "VSwitch"
  enable_details = true
}

resource "alibabacloudstack_polardb_dbinstance" "default" {
  instance_storage = "5"
  instance_name = "${var.name}"
  storage_type = "local_ssd"
  engine = "MySQL"
  engine_version = "5.7"
  instance_type = "rds.mysql.t1.small"
}

resource "alibabacloudstack_polardb_backup" "default" {
  db_instance_id = "${alibabacloudstack_polardb_dbinstance.default.id}"
  backup_method = "Physical"
}
```

## 参数参考

支持以下参数：
  * `backup_method` - (选填) - 数据备份方法，仅支持快照备份，取值固定为**Snapshot** 。
  * `db_instance_id` - (必填) - 备份数据库实例ID
  * `backup_type` - (选填) - 备份类型，仅支持全量备份，取值固定为**FullBackup**。
  * `backup_strategy` - (选填) - 备份策略

## 属性参考

除了上述所有参数外，还导出了以下属性：
  * `backup_method` - 数据备份方法，仅支持快照备份，取值固定为**Snapshot** 。
  * `backup_intranet_download_url` - 备份内网下载地址
  * `backup_mode` - 备份模式，取值范围如下： * **Automated**：系统自动备份* **Manual**：手动备份
  * `backup_size` - 备份大小
  * `backup_id` - 备份ID。
  * `slave_status` -  从节点状态
  * `host_instance_id` - 该备份所属实例的 ID。
  * `backup_db_names` - 备份所属的数据库名称。
  * `store_status` - 存储状态
  * `backup_download_url` - 备份下载地址
  * `backup_end_time` - 本次备份结束时间（UTC时间）。 
  * `backup_start_time` - 本次备份开始时间（UTC时间）。 
  * `backup_type` - 备份类型，仅支持全量备份，取值固定为**FullBackup**。
  * `backup_strategy` - 备份策略
  * `meta_status` - 备份数据状态
  * `backup_scale` - 备份规模
  * `backup_status` - 备份状态
  * `backup_location` - 备份位置