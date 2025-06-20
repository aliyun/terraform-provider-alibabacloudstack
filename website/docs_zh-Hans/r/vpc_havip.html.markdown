---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpc_havip"
sidebar_current: "docs-Alibabacloudstack-vpc-havip"
description: |-
  Provides a vpc Havip resource.
---

# alibabacloudstack\_vpc\_havip

Provides a vpc Havip resource.

## 示例用法
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

## 参数参考

支持以下参数：
  * `associated_instance_type` - (选填) - 与HaVip绑定的实例类型。取值：- **EcsInstance**：云服务器ECS实例。- **NetworkInterface**：弹性网卡实例。
  * `associated_instances` - (选填) - 与HaVip绑定的ECS实例
  * `description` - (选填) - HaVip实例的描述，长度为2到256个字符。
  * `ha_vip_name` - (选填) - HaVip实例的名称
  * `ip_address` - (选填, 强制新建) - HaVip的私网IP地址
  * `vswitch_id` - (必填, 强制新建) - HaVip实例所属的VSwitch ID
  * `vpc_id` - (选填, 强制新建) - HaVip实例所属的VPC ID

## 属性参考

除了上述所有参数外，还导出了以下属性：
  * `associated_eip_addresses` - 与HaVip绑定的EIP
  * `associated_instances` - 与HaVip绑定的ECS实例
  * `create_time` - 创建时间
  * `ha_vip_id` - HaVip实例的ID
  * `ip_address` - HaVip的私网IP地址
  * `status` - 代表资源实例的状态
  * `vpc_id` - HaVip实例所属的VPC ID
