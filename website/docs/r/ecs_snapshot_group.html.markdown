---
subcategory: "Elastic Compute Service"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_snapshot_group"
sidebar_current: "docs-Alibabacloudstack-resource-ecs-snapshot-group"
description: |-
  Provides a ECS Snapshot Group resource.
---

# alibabacloudstack\_ecs\_snapshot\_group

> **Note:** This resource can also be referred to by the following aliases:
> - `alibabacloudstack_ecs_snapshotgroup`

Provides a ECS Snapshot Group resource. A snapshot-consistent group contains snapshots of one or more disks that belong to the same ECS instance or multiple instances within the same zone.

## Example Usage
```
variable "name" {
  default = "tf-testaccebs-snapshot-group41329"
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

resource "alibabacloudstack_ecs_securitygroup" "default" {
  name   = "${var.name}_sg"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
}

resource "alibabacloudstack_security_group_rule" "default" {
  	type = "ingress"
  	ip_protocol = "tcp"
  	nic_type = "intranet"
  	policy = "accept"
  	port_range = "22/22"
  	priority = 1
  	security_group_id = "${alibabacloudstack_ecs_securitygroup.default.id}"
  	cidr_ip = "192.168.0.0/16"
}

data "alibabacloudstack_images" "default" {
  name_regex  = "^ubuntu_"
  //name_regex  = "arm_centos_7_6_20G_20211110.raw"
  //name_regex  = "^arm_centos_7"
  most_recent = true
  owners      = "system"
}

data "alibabacloudstack_instance_types" "all" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
}

data "alibabacloudstack_instance_types" "any_n4" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  instance_type_family = "ecs.n4"
  sorted_by            = "Memory"
}

data "alibabacloudstack_instance_types" "default" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  cpu_core_count       = 1
  memory_size          = 1
  instance_type_family = "ecs.n4"
  sorted_by            = "Memory"
}

locals {
	default_instance_type_id = try(element(sort(length(data.alibabacloudstack_instance_types.default.instance_types) > 0 ? data.alibabacloudstack_instance_types.default.ids : data.alibabacloudstack_instance_types.any_n4.ids), 0), sort(data.alibabacloudstack_instance_types.all.ids)[0])
}


resource "alibabacloudstack_ecs_instance" "default" {
  image_id             = "${data.alibabacloudstack_images.default.images.0.id}"
  instance_type        = "${local.default_instance_type_id}"
  system_disk_category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
  system_disk_size     = 20
  system_disk_name     = "test_sys_disk"
  security_groups      = [alibabacloudstack_ecs_securitygroup.default.id]
  instance_name        = "${var.name}_ecs"
  vswitch_id           = alibabacloudstack_vpc_vswitch.default.id
  zone_id    		   = data.alibabacloudstack_zones.default.zones.0.id
  is_outdated          = false
  data_disks {
      name                 = "disk1"
      category             = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
      size                 = 20
      delete_with_instance = true
    }
  data_disks {
      name                 = "disk2"
      category             = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
      size                 = 20
	  delete_with_instance = true
    }
  lifecycle {
    ignore_changes = [
      instance_type
    ]
  }
}

data "alibabacloudstack_ecs_disks" "disks" {
	instance_id = "${alibabacloudstack_ecs_instance.default.id}"
}

resource "alibabacloudstack_ecs_snapshot_group" "default" {
  instant_access_retention_days = "7"
  snapshot_group_name = "${var.name}"
  disk_ids = [
               "${data.alibabacloudstack_ecs_disks.disks.disks.0.id}",
               "${data.alibabacloudstack_ecs_disks.disks.disks.1.id}",
               "${data.alibabacloudstack_ecs_disks.disks.disks.2.id}"
             ]
  description = "${var.name}"
  instance_id = "${alibabacloudstack_ecs_instance.default.id}"
  instant_access = "true"
}
```

## Argument Reference

The following arguments are supported:

  * `disk_ids` - (Optional) The IDs of the cloud disks for which you want to create a snapshot-consistent group. You can specify the IDs of cloud disks that are attached to multiple instances within the same zone. Valid values of N: 1 to 16. A snapshot-consistent group can contain snapshots of up to 16 cloud disks whose total disk size does not exceed 32 TiB. This parameter cannot be used together with `exclude_disk_ids`. If you set `instance_id`, this parameter can only specify disks attached to the specified instance.
  * `exclude_disk_ids` - (Optional) The IDs of the cloud disks for which you do not want to create snapshots. After you specify the IDs of cloud disks, the snapshot-consistent group that you create does not contain the snapshots of the specified cloud disks. Valid values of N: 1 to 16. Default value: empty, which means snapshots are created for all disks of the instance. This parameter cannot be used together with `disk_ids`.
  * `instance_id` - (Optional) The ID of the ECS instance. You can set this parameter to create a snapshot-consistent group for disks on a specific ECS instance.
  * `snapshot_group_name` - (Optional) The name of the snapshot-consistent group. The name must be 2 to 128 characters in length, and can contain letters, digits, periods (.), underscores (_), hyphens (-), and colons (:). It must start with a letter or Chinese character and cannot start with `http://` or `https://`.
  * `description` - (Optional) The description of the snapshot-consistent group. The description must be 2 to 256 characters in length and cannot start with `http://` or `https://`.
  * `instant_access` - (Optional, Available in 3.18+) Specifies whether to enable the instant access feature for the snapshot. Valid values: `true`, `false`. Default value: `false`.
  * `instant_access_retention_days` - (Optional, Available in 3.18+) The retention period for the instant access feature. Unit: days. Valid values: 1 to 65535. This parameter takes effect only when `instant_access` is set to `true`. After the retention period expires, the instant access feature is automatically disabled. Default value: empty, which means the instant access feature is disabled when the snapshot is released.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

  * `id` - The ID of the snapshot-consistent group.
  * `snapshot_group_id` - The ID of the snapshot-consistent group.
  * `status` - The status of the snapshot-consistent group. Valid values: `progressing`, `accomplished`, `failed`.
  * `create_time` - The time when the snapshot-consistent group was created.
  * `disk_ids` - The IDs of the cloud disks included in the snapshot-consistent group.

## Import

ECS Snapshot Group can be imported using the snapshot group ID, e.g.

```
$ terraform import alibabacloudstack_ecs_snapshot_group.example ssg-j6ciyh3k52qp7ovm****
```
