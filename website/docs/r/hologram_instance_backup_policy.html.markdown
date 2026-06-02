---
subcategory: "Hologres"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_hologram_instance_backup_policy"
description: |-
  Provides a Hologram Instance Backup Policy resource.
---

# alibabacloudstack_hologram_instance_backup_policy

Provides a Hologram Instance Backup Policy resource.

## Example Usage

```hcl
variable "name" {
  default = "tf-testacc"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
}

data "alibabacloudstack_hologram_clusters" "default" {
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
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

resource "alibabacloudstack_hologram_instance" "example" {
  zone_id =  "${data.alibabacloudstack_zones.default.zones.0.id}"
  instance_name = "${var.name}"
  compute_type = "Standard"
  cpu = "intel"
  node = 2
  cluster = "${data.alibabacloudstack_hologram_clusters.default.clusters.0.id}"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
  vswitch_id = "${data.alibabacloudstack_vpc_vswitch.default.id}"
}

resource "alibabacloudstack_hologram_instance_backup_policy" "example" {
  instance_id        = "${alibabacloudstack_hologram_instance.example.id}"
  hour               = 2
  data_keep_quantity = 7
  week               = ["0", "2", "4", "6"]
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the Hologram instance.
* `hour` - (Required) The hour of the day to perform backup. Valid values: 0-23.
* `data_keep_quantity` - (Required) The number of days to keep backup data. Valid values: 1-31.
* `week` - (Required) The days of the week to perform backups. Valid values: "0" (Sunday), "1" (Monday), "2" (Tuesday), "3" (Wednesday), "4" (Thursday), "5" (Friday), "6" (Saturday).
* `enabled` - (Optional) Whether to enable the backup policy. Default to `true`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource, same as `instance_id`.

## Import

Hologram Instance Backup Policy can be imported using the id (instance_id), e.g.

```shell
$ terraform import alibabacloudstack_hologram_instance_backup_policy.example example-instance-id
```