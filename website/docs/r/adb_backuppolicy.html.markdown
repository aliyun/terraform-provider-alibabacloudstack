---
subcategory: "ADB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_adb_backup_policy"
sidebar_current: "docs-Alibabacloudstack-adb-backuppolicy"
description: |- 
  Provides a adb Backup Policy resource.
---

# alibabacloudstack_adb_backup_policy
-> **NOTE:** Alias name has: `alibabacloudstack_adb_backuppolicy`

Provides a [ADB](https://www.alibabacloud.com/help/product/92664.htm) cluster backup policy resource and used to configure cluster backup policy.

-> Each DB cluster has a backup policy and it will be set default values when destroying the resource.

## Example Usage

```hcl
variable "name" {
  default = "adbClusterconfig"
}

variable "creation" {
  default = "ADB"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = var.creation
}

resource "alibabacloudstack_vpc" "default" {
  name       = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = alibabacloudstack_vpc.default.id
  cidr_block        = "172.16.0.0/24"
  zone_id           = data.alibabacloudstack_zones.default.zones[0].id
  vswitch_name      = var.name
}

resource "alibabacloudstack_adb_db_cluster" "default" {
  db_cluster_version  = "3.0"
  db_cluster_category = "Basic"
  db_node_class       = "C8"
  db_node_count       = 2
  db_node_storage     = 200
  mode                = "reserver"
  pay_type            = "PostPaid"
  description         = var.name
  vswitch_id          = alibabacloudstack_vswitch.default.id
  cluster_type        = "analyticdb"
  cpu_type            = "intel"
}

resource "alibabacloudstack_adb_backup_policy" "policy" {
  db_cluster_id           = alibabacloudstack_adb_db_cluster.default.id
  preferred_backup_period = ["Tuesday", "Thursday", "Saturday"]
  preferred_backup_time   = "10:00Z-11:00Z"
}
```

### Removing `alibabacloudstack_adb_backup_policy` from your configuration

The `alibabacloudstack_adb_backup_policy` resource allows you to manage your ADB cluster's backup policy, but Terraform cannot destroy it. Removing this resource from your configuration will remove it from your state file and management, but will not destroy the cluster policy. You can resume managing the cluster via the ADB Console.

## Argument Reference

The following arguments are supported:

* `db_cluster_id` - (Required, ForceNew) The ID of the ADB cluster that needs to have its backup policy configured.
* `preferred_backup_period` - (Required) The days on which the ADB cluster backup should occur. Valid values include: `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday`, `Sunday`.
* `preferred_backup_time` - (Required) The time window during which the ADB cluster backup should occur, in the format `HH:mmZ-HH:mmZ`. The interval between start and end times is one hour. Note that the time is specified in UTC.
* `backup_retention_period` - (Optional) The number of days for which data backup files are retained. This value is fixed at 7 days and cannot be modified.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The current backup policy resource ID. It is the same as `db_cluster_id`.
* `backup_retention_period` - The number of days for which data backup files are retained. This value is fixed at 7 days and cannot be modified.