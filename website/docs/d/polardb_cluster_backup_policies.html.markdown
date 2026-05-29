---
subcategory: "PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_cluster_backup_policies"
sidebar_current: "docs-Alibabacloudstack-datasource-polardb-cluster-backup-policies"
description: |-
  Provides a PolarDB cluster backup policies data source.
---

# alibabacloudstack_polardb_cluster_backup_policies

This data source provides the backup policies of a PolarDB cluster.

## Example Usage

```hcl
variable "name" {
  default = "tfaccpolicy-datasource-test"
}

variable "db_type" {
  default = "MySQL"
}

variable "db_version" {
  default = "8.0"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
  lifecycle {
      ignore_changes = [
		secondary_cidr_blocks,
        tags
      ]
  }
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  vswitch_name = "${var.name}_vsw"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block = "172.16.1.0/24"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
  lifecycle {
      ignore_changes = [
        tags
      ]
  }
}

data "alibabacloudstack_polardb_cluster_instance_types" "default" {
  db_type = "${var.db_type}"
  db_version = "${var.db_version}"
  sorted_by = "CPU"
  cpu_type = "hygon"
  sub_category = "normal_exclusive"
}

resource "alibabacloudstack_polardb_cluster_instance" "default" {
	db_cluster_description 	=  "${var.name}"
	zone_id 				= "${data.alibabacloudstack_zones.default.zones.0.id}"
	db_type 				= "${var.db_type}"
	db_version 				= "${var.db_version}"
	storage_space 			= "20"
	vpc_id 					= "${alibabacloudstack_vpc_vpc.default.id}"
	vswitch_id				= "${alibabacloudstack_vpc_vswitch.default.id}"
	db_node_class 			= "${data.alibabacloudstack_polardb_cluster_instance_types.default.instance_types.0.id}"
	sub_category 			= "${data.alibabacloudstack_polardb_cluster_instance_types.default.instance_types.0.sub_category}"
	storage_type 			= "ESSDPL1"
}

data "alibabacloudstack_polardb_cluster_backup_policies" "policy" {
	db_cluster_id = "${alibabacloudstack_polardb_cluster_instance.default.id}"
}

```

## Argument Reference

The following arguments are supported:

* `db_cluster_id` - (Required) The ID of the PolarDB cluster.

## Attributes Reference

The following attributes are exported:

* `data_level1_backup_frequency` - The frequency of Level-1 data backups.
* `data_level1_backup_period` - The period of Level-1 data backups.
* `data_level1_backup_time` - The time of Level-1 data backups.
* `data_level1_backup_retention_period` - The retention period of Level-1 data backups.
* `data_level2_backup_retention_period` - The retention period of Level-2 data backups.
* `backup_retention_policy_on_cluster_deletion` - The backup retention policy when the cluster is deleted.
* `log_backup_retention_period` - The retention period of log backups.
* `preferred_backup_time` - The preferred backup time.
* `preferred_backup_period` - The preferred backup period.
* `backup_retention_period` - The backup retention period.
* `backup_frequency` - The backup frequency.
* `preferred_next_backup_time` - The preferred next backup time.
* `data_level2_backup_period` - The period of Level-2 data backups.
* `data_level2_backup_another_region_region` - The region for Level-2 backups in another region.
* `data_level2_backup_another_region_retention_period` - The retention period for Level-2 backups in another region.
* `log_backup_another_region_region` - The region for log backups in another region.
* `log_backup_another_region_retention_period` - The retention period for log backups in another region.
* `enable_backup_log` - Whether log backup is enabled.