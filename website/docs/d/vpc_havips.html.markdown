---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpc_havips"
sidebar_current: "docs-Alibabacloudstack-datasource-vpc-havips"
description: |-
  Provides a list of vpc havips owned by an alibabacloudstack account.
---

# alibabacloudstack\_vpc\_havips

This data source provides a list of vpc havips in an alibabacloudstack account according to the specified filters.

## Example Usage
```
variable "name" {
  default = "tf-testAccVpcHaVipsDataSource-7785094"
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

resource "alibabacloudstack_ecs_networkinterface" "default" {
  	count                	= 2
	network_interface_name 	= "${var.name}_eni"
    vswitch_id 				= "${alibabacloudstack_vpc_vswitch.default.id}"
	security_groups      	= [alibabacloudstack_ecs_securitygroup.default.id]
}

resource "alibabacloudstack_vpc_ha_vip" "default" {
  ha_vip_name = "${var.name}"
  description = "${var.name}"
  ip_address  = "172.16.1.88"
  vswitch_id  = "${alibabacloudstack_vpc_vswitch.default.id}"
  vpc_id      = "${alibabacloudstack_vpc_vpc.default.id}"
  associated_instance_type = "NetworkInterface"
  associated_instances = ["${alibabacloudstack_ecs_networkinterface.default.0.id}", "${alibabacloudstack_ecs_networkinterface.default[1].id}",]
}


data "alibabacloudstack_vpc_ha_vips" "default" {
  description_regex = "${alibabacloudstack_vpc_ha_vip.default.description}"
}
```

## Argument Reference

The following arguments are supported:
  * `ids` - (Optional) - The IDs of the vpc havips.
  * `filter` - (Optional) - Filter condition.
    * `key` - (Optional) - The key of the filter condition..
    * `value` - (Optional) - The value of the filter condition..

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `filter` -  Filter condition.
  * `ha_vips` - A list of vpc havip.
    * `id` - The ID of the resource
    * `associated_eip_addresses` - EIP bound to HaVip
    * `associated_instance_type` - The type of the instance that is bound to the HaVip. Value:-**EcsInstance**: ECS instance.-**NetworkInterface**: ENI instance.
    * `associated_instances` - An ECS instance that is bound to HaVip
    * `create_time` - The creation time of the resource
    * `description` - The description of the HaVip instance. The length is 2 to 256 characters.
    * `ha_vip_id` - The ID of the resource
    * `ha_vip_name` - The name of the HaVip instance
    * `ip_address` - The ip address of the HaVip. If not filled, the default will be assigned one from the vswitch.
    * `master_instance_id` - The primary instance ID bound to HaVip
    * `status` - The status of this resource instance.
    * `tags` - The tags of HaVip.
    * `vswitch_id` - The vswitch ID to which the HaVip instance belongs
    * `vpc_id` - The VPC ID to which the HaVip instance belongs
