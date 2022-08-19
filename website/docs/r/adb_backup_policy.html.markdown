---
subcategory: "AnalyticDB for MySQL (ADB)"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_adb_backup_policy"
sidebar_current: "docs-apsarastack-resource-adb-backup-policy"
description: |-
  Provides a ADB backup policy resource.
---

# apsarastack\_adb\_backup\_policy

Provides a [ADB](https://www.alibabacloud.com/help/product/92664.htm) cluster backup policy resource and used to configure cluster backup policy.

-> Each DB cluster has a backup policy and it will be set default values when destroying the resource.

## Example Usage

```
variable "name" {
  default = "adbClusterconfig"
}

variable "creation" {
  default = "ADB"
}

data "apsarastack_zones" "default" {
  available_resource_creation = var.creation
}

resource "apsarastack_vpc" "default" {
  name       = var.name
  cidr_block = "172.16.0.0/16"
}

resource "apsarastack_vswitch" "default" {
  vpc_id            = apsarastack_vpc.default.id
  cidr_block        = "172.16.0.0/24"
  zone_id           = data.apsarastack_zones.default.zones[0].id
  vswitch_name      = var.name
}

resource "apsarastack_adb_db_cluster" "default" {
  db_cluster_version  = "3.0"
  db_cluster_category = "Basic"
  db_node_class       = "C8"
  db_node_count       = 2
  db_node_storage     = 200
  mode				  = "reserver"
  pay_type            = "PostPaid"
  description         = var.name
  vswitch_id          = apsarastack_vswitch.default.id
  cluster_type      = "analyticdb"
  cpu_type          = "intel"
}

resource "apsarastack_adb_backup_policy" "policy" {
  db_cluster_id           = apsarastack_adb_db_cluster.default.id
  preferred_backup_period = ["Tuesday", "Thursday", "Saturday"]
  preferred_backup_time   = "10:00Z-11:00Z"
}
```
### Removing apsarastack_adb_cluster from your configuration
 
The apsarastack_adb_backup_policy resource allows you to manage your adb cluster policy, but Terraform cannot destroy it. Removing this resource from your configuration will remove it from your statefile and management, but will not destroy the cluster policy. You can resume managing the cluster via the adb Console.
 
## Argument Reference

The following arguments are supported:

* `db_cluster_id` - (Required, ForceNew) The Id of cluster that can run database.
* `preferred_backup_period` - (Required) ADB Cluster backup period. Valid values: [Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday].
* `preferred_backup_time` - (Required) ADB Cluster backup time, in the format of HH:mmZ- HH:mmZ. Time setting interval is one hour. China time is 8 hours behind it.

## Attributes Reference

The following attributes are exported:

* `id` - The current backup policy resource ID. It is same as 'db_cluster_id'.
* `backup_retention_period` - Cluster backup retention days, Fixed for 7 days, not modified.

## Import

ADB backup policy can be imported using the id or cluster id, e.g.

```
$ terraform import apsarastack_adb_backup_policy.example "am-12345678"
```
