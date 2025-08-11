---
subcategory: "PolarDBX"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardbx_backup_policies"
sidebar_current: "docs-Alibabacloudstack-datasource-polardbx-backup-policies"
description: |-
  Provides a list of PolarDBX backup policies owned by an Alibaba Cloud Stack account.
---

# alibabacloudstack\_polardbx\_backup\_policies

This data source provides a list of PolarDBX backup policies in an Alibaba Cloud Stack account according to the specified filters.

## Example Usage

```hcl
variable "name" {
  default = "tf-testAccPolardbxInstancesDataSource-3625795"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details              = true
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name   = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  name       = "${var.name}_vsw"
  vpc_id     = alibabacloudstack_vpc_vpc.default.id
  cidr_block = "172.16.1.0/24"
  zone_id    = data.alibabacloudstack_zones.default.zones.0.id
}

resource "alibabacloudstack_polardbx_instance" "default" {
  description      = "testtf1111"
  series           = "enterprise"
  topology_type    = "1azone"
  zone_id          = data.alibabacloudstack_zones.default.zones.0.id
  engine_version   = "5.7"
  storage          = "50"
  network_type     = "vpc"
  vpc_id           = alibabacloudstack_vpc_vpc.default.id
  vswitch_id       = alibabacloudstack_vpc_vswitch.default.id
  cn_node_class    = "polarx.x4.medium.2e"
  cn_node_count    = "2"
  dn_node_class    = "mysql.n4.medium.25"
  dn_node_count    = "2"
}

data "alibabacloudstack_polardbx_backup_policies" "default" {
  db_instance_id = alibabacloudstack_polardbx_instance.default.id
}
```

## Argument Reference
The following arguments are supported:

* `db_instance_id` - (Required, ForceNew) The ID of the PolarDBX instance.


## Attributes Reference
The following attributes are exported in addition to the arguments listed above:

* `backup_period` - Backup cycle. You can select multiple days in a week. Valid values: Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday. Multiple values are separated by commas.
* `backup_set_retention` - The retention period of the backup set. Valid values: 7 to 730 days.
* `backup_plan_begin` - The start time of the backup. The time format is HH:mmZ, for example, 03:00Z.
* `remove_log_retention` - The retention period of transaction logs. Valid values: 7 to 730 days.
* `cold_data_backup_interval` - The interval for cold data backup.
* `local_log_retention_number` - The number of transaction logs retained locally. Valid values: 6 to 100.
* `cold_data_backup_retention` - The retention period of cold data backups.
* `force_clean_on_high_space_usage` - Specifies whether to forcibly delete backups when the storage space usage is high. Valid values: 0 (no), 1 (yes).
* `backup_way` - Backup method.
* `local_log_retention` - The retention period of transaction logs that are retained locally. Unit: days.
* `backup_type` - Backup type.
* `log_local_retention_space` - The maximum space occupied by transaction logs that are retained locally. Unit: %. Valid values: 0 to 100.