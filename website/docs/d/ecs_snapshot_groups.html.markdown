---
subcategory: "Elastic Compute Service"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_snapshot_groups"
description: |-
  Provides a list of ecs snapshotgroups owned by an alibabacloudstack account.
---

# alibabacloudstack\_ecs\_snapshot_groups

This data source provides a list of ecs snapshotgroups in an alibabacloudstack account according to the specified filters.

## Example Usage
```
variable "name" {
  default = "tf-testAccEcsSnapshotGroupsDataSource-6474718"
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
  	description =                   "${var.name}"
	instance_id =                   "${alibabacloudstack_ecs_instance.default.id}"
	instant_access =                "true"
	instant_access_retention_days = "7"
	snapshot_group_name =           "${var.name}"
	disk_ids = [
		"${data.alibabacloudstack_ecs_disks.disks.disks.0.id}",
		"${data.alibabacloudstack_ecs_disks.disks.disks.1.id}",
		"${data.alibabacloudstack_ecs_disks.disks.disks.2.id}",
	]
}

 

data "alibabacloudstack_ecs_snapshot_groups" "default" {
  description_regex = "${alibabacloudstack_ecs_snapshot_group.default.description}"
}
```

## Argument Reference

The following arguments are supported:
  * `ids` - (Optional) - The ids of the snapshotgroups.
  * `name_regex` - (Optional) - A regex string to filter resulting snapshotgroups by name.
  * `description_regex` - (Optional) -  A regex string to filter resulting snapshotgroups by description.
  * `instance_id` - (Optional) - instance id
  * `snapshot_group_name` - (Optional) - name

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `snapshot_groups` - A list of snapshotgroups. Each element contains the following attributes:
    * `id` - The ID of the snapshotgroup.
    * `create_time` - The creation time of the resource
    * `description` - snapshot group description
    * `snapshots` - The snapshot list in the snapshotgroup.
    * `instance_id` - instance id
    * `snapshot_group_id` - The first ID of the resource
    * `snapshot_group_name` - name
    * `status` - The status of the resource
