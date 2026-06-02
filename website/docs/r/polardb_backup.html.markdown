---
subcategory: "PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_backup"
description: |-
  Provides a polardb Backup resource.
---

# alibabacloudstack\_polardb\_backup

Provides a polardb Backup resource.

## Example Usage
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

## Argument Reference

The following arguments are supported:
  * `backup_method` - (Optional) - The data backup method. Only Snapshot backup is supported. The value is fixed to **Snapshot * *.
  * `db_instance_id` - (Required) - The ID of the PolarDB Instance.
  * `backup_type` - (Optional) - The backup type. Only full backup is supported. The value is fixed to **FullBackup * *.
  * `backup_strategy` - (Optional) - Backup strategy.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `backup_method` - The data backup method. Only Snapshot backup is supported. The value is fixed to **Snapshot * *.
  * `backup_intranet_download_url` - The internal network download address of the backup.
  * `backup_mode` - Backup mode, the value range is as follows:* **Automated**: Automatic System Backup* **Manual**: Manual backup
  * `backup_size` - The size of the backup.
  * `backup_id` - The backup ID.
  * `slave_status` - The status of the backup.
  * `host_instance_id` - The ID of the ECS instance to which the backup belongs.
  * `backup_db_names` - The name of the database for which the backup is created.
  * `store_status` - StoreStatus
  * `backup_download_url` - The download address of the backup.
  * `backup_end_time` - The end time of this backup (UTC time).
  * `backup_start_time` - The backup start time (UTC time).
  * `backup_type` - The backup type. Only full backup is supported. The value is fixed to **FullBackup * *.
  * `backup_strategy` - The backup strategy.
  * `meta_status` - The status of the backup.
  * `backup_scale` - The backup scale.
  * `backup_status` - The backup status.
  * `backup_location` - The backup location. 
