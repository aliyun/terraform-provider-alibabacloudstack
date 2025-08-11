---
subcategory: "PolarDBX"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardbx_backup"
sidebar_current: "docs-Alibabacloudstack-polardbx-backup"
description: |-
  Provides a polardbx Backup resource.
---

# alibabacloudstack\_polardbx\_backup

Provides a polardbx Backup resource.

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

resource "alibabacloudstack_polardbx_backup" "default" {
  instance_id = "${alibabacloudstack_polardbx_instance.default.id}"
}
```

## Argument Reference

The following arguments are supported:
  * `instance_id` - (Required) - The ID of the PolarDBX Instance.
  * `backup_type` - (Optional) - The backup type. Currently only supports "0".

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `backup_mode` - Backup mode. Currently only supports "0".
  * `backup_set_size` - The size of the backup.
  * `backup_set_id` - The backup set ID.
  * `status` - The status of the backup.
  * `end_time` - The end time of this backup (UTC time).
  * `begin_time` - The backup start time (UTC time).
