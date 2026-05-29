---
subcategory: "Cloud-Native Distributed Database PolarDB-X 2.0"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardbx_backup"
sidebar_current: "docs-Alibabacloudstack-resource-polardbx-backup"
description: |-
  Provides a PolarDB-X Backup resource.
---

# alibabacloudstack\_polardbx\_backup

Provides a PolarDB-X backup resource. This resource allows you to create and manage backups for PolarDB-X instances.

-> **Note:** This resource does not support update operations. Modifying any argument will force a new resource to be created.

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

  * `instance_id` - (Required, ForceNew) The ID of the PolarDB-X instance. Modifying this parameter will force a new resource to be created.
  * `backup_type` - (Optional, ForceNew) The type of backup. Valid values: `0` (physical backup). Default to `0`. Modifying this parameter will force a new resource to be created.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

  * `id` - The ID of the backup. The format is `<instance_id>:<backup_set_id>`.
  * `backup_model` - The backup mode. Valid values: `0` (physical backup), `1` (logical backup).
  * `backup_set_size` - The size of the backup set, in bytes.
  * `backup_set_id` - The ID of the backup set.
  * `status` - The status of the backup. Valid values: `0` (creating), `1` (success), `2` (failed).
  * `end_time` - The end time of the backup in UTC format.
  * `begin_time` - The start time of the backup in UTC format.

## Import

PolarDB-X Backup can be imported using the instance ID and backup set ID separated by a colon, e.g.

```
$ terraform import alibabacloudstack_polardbx_backup.example pc-xxxxxxxxxxxxx:1234567890
```
