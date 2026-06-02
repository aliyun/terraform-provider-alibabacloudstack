---
subcategory: "Virtual Private Cloud (VPC)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpc_ha_vip"
description: |-
  Provides a vpc Havip resource.
---

# alibabacloudstack\_vpc\_havip

Provides a vpc Havip resource.

## Example Usage
```
variable "name" {
  default = "tf-testAccVpcHavipBasic"
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

resource "alibabacloudstack_ecs_instance" "default0" {
  image_id             = "${data.alibabacloudstack_images.default.images.0.id}"
  instance_type        = "${local.default_instance_type_id}"
  system_disk_category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
  system_disk_size     = 20
  system_disk_name     = "test_sys_disk"
  security_groups      = [alibabacloudstack_ecs_securitygroup.default.id]
  instance_name        = "${var.name}_ecs"
  vswitch_id           = alibabacloudstack_vpc_vswitch.default.id
  zone_id    		   = data.alibabacloudstack_zones.default.zones.0.id
  lifecycle {
    ignore_changes = [
      instance_type,
	  system_disk_category
    ]
  }
}

resource "alibabacloudstack_ecs_instance" "default1" {
  image_id             = "${data.alibabacloudstack_images.default.images.0.id}"
  instance_type        = "${local.default_instance_type_id}"
  system_disk_category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
  system_disk_size     = 20
  system_disk_name     = "test_sys_disk"
  security_groups      = [alibabacloudstack_ecs_securitygroup.default.id]
  instance_name        = "${var.name}_ecs1"
  vswitch_id           = alibabacloudstack_vpc_vswitch.default.id
  zone_id    		   = data.alibabacloudstack_zones.default.zones.0.id
  lifecycle {
    ignore_changes = [
      instance_type,
	  system_disk_category
    ]
  }
}

resource "alibabacloudstack_ecs_networkinterface" "default" {
  	count                	= 2
	network_interface_name 	= "${var.name}_eni"
    vswitch_id 				= "${alibabacloudstack_vpc_vswitch.default.id}"
	security_groups      	= [alibabacloudstack_ecs_securitygroup.default.id]
}

resource "alibabacloudstack_vpc_ha_vip" "default" {
  associated_instances = [
                           "${alibabacloudstack_ecs_instance.default0.id}",
                           "${alibabacloudstack_ecs_instance.default1.id}"
                         ]
  ha_vip_name = "${var.name}"
  description = "${var.name}"
  ip_address = "172.16.1.88"
  vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
  associated_instance_type = "EcsInstance"
}
```

## Argument Reference

The following arguments are supported:
  * `associated_instance_type` - (Optional) - The type of the instance that is bound to the HaVip. Value:-**EcsInstance**: ECS instance.-**NetworkInterface**: ENI instance.
  * `associated_instances` - (Optional) - An ECS instance that is bound to HaVip
  * `description` - (Optional) - The description of the HaVip instance. The length is 2 to 256 characters.
  * `ha_vip_name` - (Optional) - The name of the HaVip instance
  * `ip_address` - (Optional, ForceNew) - The ip address of the HaVip. If not filled, the default will be assigned one from the vswitch.
  * `vswitch_id` - (Required, ForceNew) - The vswitch ID to which the HaVip instance belongs
  * `vpc_id` - (Optional, ForceNew) - The VPC ID to which the HaVip instance belongs

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `associated_eip_addresses` - EIP bound to HaVip
  * `associated_instances` - An ECS instance that is bound to HaVip
  * `create_time` - The creation time of the resource
  * `ha_vip_id` - The ID of the resource
  * `ip_address` - The ip address of the HaVip. If not filled, the default will be assigned one from the vswitch.
  * `status` - The status of this resource instance.
  * `vpc_id` - The VPC ID to which the HaVip instance belongs
