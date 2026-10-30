---
subcategory: "PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_cluster_backups"
sidebar_current: "docs-Alibabacloudstack-datasource-polardb-cluster-backups"
description: |-
  This data source provides the backup sets of a PolarDB cluster edition instance.
---

# alibabacloudstack\_polardb\_cluster\_backups

This data source provides the backup sets of a PolarDB **cluster edition** instance.

By default the backup sets created in the last 7 days are returned, because `DescribeBackups` of the
cluster edition makes both `StartTime` and `EndTime` mandatory; the default end bound is padded by one
day so that a backup set finishing seconds after its own minute is never missed. Both ends of the
window use the `2006-01-02T15:04Z` format; a value carrying seconds is rejected by the service.
`end_time` may be a future time, while a future `start_time` is rejected with
`InvalidStartTime.Malformed`; a window whose start is later than its end simply matches nothing.

## Example Usage
```hcl
data "alibabacloudstack_polardb_cluster_backups" "default" {
  db_cluster_id = "${alibabacloudstack_polardb_cluster_instance.default.id}"
}

```

## Argument Reference

The following arguments are supported:

  * `ids` - (Optional, ForceNew, Computed) A list of backup set IDs to filter the search by.
  * `db_cluster_id` - (Required) The ID of the PolarDB cluster.
  * `start_time` - (Optional) The beginning of the time window in which the backup sets were created.
  * `end_time` - (Optional) The end of the time window in which the backup sets were created.

## Attributes Reference

The following attributes are exported:

  * `ids` - A list of the backup set IDs found.
  * `backups` - A list of backup sets. Each element contains the following attributes:
    * `id` - The ID of the backup set.
    * `backup_id` - The ID of the backup set.
    * `db_cluster_id` - The ID of the PolarDB cluster the backup set belongs to.
    * `backup_method` - The backup method, for example `Physical`.
    * `backup_mode` - The backup mode, `Automated` or `Manual`.
    * `backup_type` - The backup type, for example `FullBackup`.
    * `backup_status` - The status of the backup set, `Success` once the job has finished.
    * `backups_level` - The level of the backup set, for example `Level-1`.
    * `store_status` - Whether the backup set is available for restore, for example `Enabled`.
    * `backup_set_size` - The size of the backup set, in bytes.
    * `is_avail` - Whether the backup set is available.
    * `consistent_time` - The consistent point in time of the backup set, in epoch seconds.
    * `backup_start_time` - The time at which the backup started.
    * `backup_end_time` - The time at which the backup finished.
