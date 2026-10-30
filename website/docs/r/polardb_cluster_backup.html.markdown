---
subcategory: "PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_cluster_backup"
sidebar_current: "docs-Alibabacloudstack-resource-polardb-cluster-backup"
description: |-
  Provides a PolarDB cluster edition manual backup resource.
---

# alibabacloudstack\_polardb\_cluster\_backup

Provides a PolarDB **cluster edition** manual backup resource, which triggers a full backup of a
`alibabacloudstack_polardb_cluster_instance` and waits until the backup set is registered.

The cluster edition and the primary/standby edition are two different products on the POP gateway:
this resource calls `polardb / 2017-08-01` with `DBClusterId`, while
[`alibabacloudstack_polardb_backup`](/docs/r/polardb_backup.html.markdown) calls
`polardb / 2024-01-30` with `DBInstanceId`. Use this resource for cluster ids starting with `pc-`.

## Example Usage
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

## Argument Reference

The following arguments are supported:

  * `db_cluster_id` - (Required, ForceNew) The ID of the PolarDB cluster.

`CreateBackup` of the cluster edition accepts no other parameter: the backup method, mode and type
are reported back by the service as attributes.

## Attributes Reference

The following attributes are exported:

  * `backup_job_id` - The ID of the backup job returned by `CreateBackup`. It belongs to another ID
    space than `backup_id`, so it cannot be used to look the backup set up.
  * `backup_id` - The ID of the backup set, used together with `db_cluster_id` as the resource ID.
  * `backup_method` - The backup method, for example `Physical`.
  * `backup_mode` - The backup mode, `Manual` for a backup triggered by this resource.
  * `backup_type` - The backup type, for example `FullBackup`.
  * `backup_status` - The status of the backup set, `Success` once the job has finished.
  * `backups_level` - The level of the backup set, for example `Level-1`.
  * `store_status` - Whether the backup set is available for restore, for example `Enabled`.
  * `backup_set_size` - The size of the backup set, in bytes.
  * `is_avail` - Whether the backup set is available.
  * `consistent_time` - The consistent point in time of the backup set, in epoch seconds.
  * `backup_start_time` - The time at which the backup started.
  * `backup_end_time` - The time at which the backup finished.

## Timeouts

The `timeouts` block allows you to specify the following action:

  * `create` - (Defaults to 20 minutes) Used when the backup job is triggered and the backup set is
    being registered.

## Import

PolarDB cluster backups can be imported with `terraform import` using the `db_cluster_id:backup_id`
format, separated by a colon:

```
$ terraform import alibabacloudstack_polardb_cluster_backup.default pc-qj0jt53b6trb0443y:1793345712
```

## Notes

  * Destroying this resource does not delete the backup set. The lifecycle of backup sets is managed
    by the backup policy of the cluster, see
    [`alibabacloudstack_polardb_cluster_backup_policy`](/docs/r/polardb_cluster_backup_policy.html.markdown).
  * `backup_job_id` is only known for a backup created by Terraform, so it is empty on an imported
    backup set.
