---
subcategory: "Cloud-Native Distributed Database PolarDB-X 2.0"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardbx_backups"
sidebar_current: "docs-Alibabacloudstack-datasource-polardbx-backups"
description: |-
  Provides a list of polardbx backups owned by an alibabacloudstack account.
---

# alibabacloudstack\_polardbx\_backups

This data source provides a list of polardbx backups in an alibabacloudstack account according to the specified filters.

## Example Usage
```
variable "name" {
  default = "tf-testAccPolardbxInstancesDataSource-3625795"
}

data "alibabacloudstack_zones" default {
  available_resource_creation = "VSwitch"
  enable_details = true
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  name = "${var.name}_vsw"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block = "172.16.1.0/24"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

resource "alibabacloudstack_polardbx_instance" "default" {
  description = "testtf1111"
	series = "enterprise"
	topology_type = "1azone"
	zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
	engine_version = "5.7"
	storage = "50"
	network_type = "vpc"
	vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
	vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
	cn_node_class = "polarx.x4.medium.2e"
	cn_node_count = "2"
	dn_node_class = "mysql.n4.medium.25"
	dn_node_count = "2"
}

data "alibabacloudstack_polardbx_backups" "default" {
  db_instance_id = "${alibabacloudstack_polardbx_instance.default.id}"
}
```

## Argument Reference

The following arguments are supported:
  * `ids` - (Optional) - A list of backup IDs.
  * `db_instance_id` - (Required, ForceNew) - The ID of the ApsaraDB for PolarDBX cluster for which backups are to be queried.
  * `end_time` - (Optional) - The end time of this backup (UTC time).
  * `start_time` - (Optional) - The backup start time (UTC time).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `backups` - The list of PolarDBX backups.
    * `id` - The ID of the backup.
    * `backup_method` - The data backup method. Only Snapshot backup is supported. The value is fixed to **Snapshot * *.
    * `backup_model` - Backup mode
    * `backup_set_size` - The backup size.
    * `backup_type` - The backup type.
    * `backup_set_id` - The backup set ID.
    * `status` - The status of the backup.
    * `end_time` - The end time of this backup (UTC time).
    * `begin_time` - The backup start time (UTC time).