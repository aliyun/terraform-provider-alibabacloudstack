---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_snapshotgroup"
sidebar_current: "docs-Alibabacloudstack-ecs-snapshotgroup"
description: |-
  Provides a ecs Snapshotgroup resource.
---

# alibabacloudstack\_ecs\_snapshotgroup

Provides a ecs Snapshotgroup resource.

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
  * `create_time` - (Optional) - The creation time of the resource
  * `description` - (Optional) - snapshot group description
  * `exclude_disk_ids` - (Optional) - The disk IDs excluded from the snapshot group.
  * `disk_ids` - (Optional) - The disk IDs included in the snapshot group.
  * `instance_id` - (Optional) - instance id
  * `instant_access` - (Optional) - Whether to enable snapshot speed is available. Value range:-true: on.-false: closed.The default value is false.
  * `instant_access_retention_days` - (Optional) - Set the usage time available for snapshot speed. Unit: days, value range: 1~65535.This parameter takes effect only when 'InstantAccess = true. Automatically turn off the snapshot speed function after expiration.Default value: null, indicating the same snapshot release time.
  * `snapshot_group_name` - (Optional) - name

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `create_time` - The creation time of the resource
  * `disk_ids` - The disk IDs included in the snapshot group.
  * `snapshot_group_id` - The first ID of the resource
  * `status` - The status of the resource
