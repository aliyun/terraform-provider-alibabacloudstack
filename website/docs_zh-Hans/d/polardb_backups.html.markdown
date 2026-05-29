---
subcategory: "云原生数据库 PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_backups"
sidebar_current: "docs-Alibabacloudstack-datasource-polardb-backups"
description: |-
  提供阿里云账号下拥有的polardb backups列表。
---

# alibabacloudstack\_polardb\_backups

此数据源提供根据指定过滤条件列出的阿里云账号下的polardb backups资源列表。

## 示例用法
```
variable "name" {
  default = "tf-testAccPolardbBackups15301"
}

data "alibabacloudstack_backups" default {
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
  backup_method  Physical
}

data "alibabacloudstack_polardb_backups" "default" {
  db_instance_id = "${alibabacloudstack_polardb_dbinstance.default.id}"
}
```

## 参数参考
以下参数是支持的：
  * `ids` - (选填) - 备份ID列表
  * `db_instance_id` - (必填, 强制新建) - 备份所属数据库实例ID

## Attributes Reference
除了上述参数外，还导出以下属性：
  * `backups` - 备份列表
    * `id` - 备份ID
    * `backup_method` - 数据备份方法，仅支持快照备份，取值固定为**Snapshot** 。
    * `backup_intranet_download_url` - 备份内网下载地址
    * `backup_mode` - 备份模式，取值范围如下： * **Automated**：系统自动备份* **Manual**：手动备份
    * `backup_size` - 备份大小
    * `backup_id` - 备份ID。
    * `slave_status` - 从节点状态
    * `host_instance_id` - 该备份所属实例的 ID。
    * `backup_db_names` - 备份所属的数据库名称。
    * `store_status` - 存储状态
    * `db_instance_id` - 备份所属的数据库实例名称。
    * `backup_download_url` - 备份下载地址
    * `backup_end_time` - 本次备份结束时间（UTC时间）。 
    * `backup_start_time` - 本次备份开始时间（UTC时间）。 
    * `backup_type` - 备份类型，仅支持全量备份，取值固定为**FullBackup**。
    * `backup_strategy` - 备份策略
    * `meta_status` - 备份数据状态
    * `backup_scale` - 备份规模
    * `backup_status` - 备份状态
    * `backup_location` - 备份位置
