---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_diskattachment"
sidebar_current: "docs-Alibabacloudstack-ecs-diskattachment"
description: |- 
  Provides a ecs Diskattachment resource.
---

# alibabacloudstack_ecs_diskattachment
-> **NOTE:** Alias name has: `alibabacloudstack_disk_attachment`

Provides a ecs Diskattachment resource.

## Example Usage

```hcl
variable "name" {
	default = "tf-testAccEcsDiskAttachmentConfig"
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
	cidr_block = "172.16.0.0/24"
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
	cidr_ip = "172.16.0.0/24"
}

data "alibabacloudstack_images" "default" {
	name_regex  = "^ubuntu_"
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
	zone_id    = data.alibabacloudstack_zones.default.zones.0.id
	is_outdated          = false
	lifecycle {
		ignore_changes = [
			instance_type
		]
	}
}

resource "alibabacloudstack_ecs_disk" "default" {
	availability_zone = data.alibabacloudstack_zones.default.zones[0].id
	size = "20"
	name = "${var.name}"
	category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"

	tags = {
		Name = "TerraformTest-disk"
	}
}

resource "alibabacloudstack_ecs_diskattachment" "default" {
	disk_id = "${alibabacloudstack_ecs_disk.default.id}"
	instance_id = "${alibabacloudstack_ecs_instance.default.id}"
	device_name = "/dev/xvdb"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the target ECS instance.
* `disk_id` - (Required, ForceNew) The ID of the cloud disk to be mounted. The cloud disk (`DiskId`) and instance (`InstanceId`) must be in the same zone. Supports mounting data disks and system disks. For related constraints, please refer to the interface description section above.
* `device_name` - (Optional, ForceNew) The device name exposed to the instance. It will be allocated automatically by the system according to the default order from `/dev/xvdb` to `/dev/xvdz`. If specified, it must match one of the available device names on the instance.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `instance_id` - ID of the Instance.
* `disk_id` - ID of the Disk.
* `device_name` - The device name exposed to the instance.