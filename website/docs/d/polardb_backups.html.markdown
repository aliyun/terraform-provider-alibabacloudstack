---
subcategory: "PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_backups"
sidebar_current: "docs-Alibabacloudstack-datasource-polardb-backups"
description: |-
  Provides a list of polardb backups owned by an alibabacloudstack account.
---

# alibabacloudstack\_polardb\_backups

This data source provides a list of polardb backups in an alibabacloudstack account according to the specified filters.

## Example Usage
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

## Argument Reference

The following arguments are supported:
  * `ids` - (Optional) - A list of backup IDs.
  * `db_instance_id` - (Required, ForceNew) - The ID of the ApsaraDB for PolarDB cluster for which backups are to be queried.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `backups` - The list of PolarDB backups.
    * `id` - The ID of the backup.
    * `backup_method` - The data backup method. Only Snapshot backup is supported. The value is fixed to **Snapshot * *.
    * `backup_intranet_download_url` - The URL of the internal network of the backup.
    * `backup_mode` - Backup mode, the value range is as follows:* **Automated**: Automatic System Backup* **Manual**: Manual backup
    * `backup_size` - The backup size.
    * `backup_id` - The backup ID.
    * `slave_status` - The status of the backup.
    * `host_instance_id` - The ID of the instance to which the backup belongs.
    * `backup_db_names` - The database names of the backup.
    * `store_status` - StoreStatus
    * `db_instance_id` - The ID of the instance to which the backup belongs.
    * `backup_download_url` - The URL of the backup.
    * `backup_end_time` - The end time of this backup (UTC time).
    * `backup_start_time` - The backup start time (UTC time).
    * `backup_type` - The backup type. Only full backup is supported. The value is fixed to **FullBackup * *.
    * `backup_strategy` - The backup strategy.
    * `meta_status` - The status of the backup meta.
    * `backup_scale` - The backup scale.
    * `backup_status` - The backup status.
    * `backup_location` - The backup location.
