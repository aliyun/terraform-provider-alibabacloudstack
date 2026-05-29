---
subcategory: "专有网络 VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpc_ha_vips"
sidebar_current: "docs-Alibabacloudstack-datasource-vpc-ha-vips"
description: |-
  提供阿里云账号下拥有的vpc havips列表。
---

# alibabacloudstack\_vpc\_havips

此数据源提供根据指定过滤条件列出的阿里云账号下的vpc havips资源列表。

## 示例用法
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

## 参数参考
以下参数是支持的：
  * `ids` - (选填) - 高可用虚拟IP实例的ID 列表。
  * `filter` - (选填) - 过滤 参数。
    * `key` - (选填) - 过滤参数得字段名。
    * `value` - (选填) - 过滤参数得值。

## Attributes Reference
除了上述参数外，还导出以下属性：
  * `filter` - 过滤 参数。
  * `ha_vips` - 高可用虚拟IP实例列表
    * `id` - 高可用虚拟IP实例的ID
    * `associated_eip_addresses` - 与HaVip绑定的EIP
    * `associated_instance_type` - 与HaVip绑定的实例类型。取值：- **EcsInstance**：云服务器ECS实例。- **NetworkInterface**：弹性网卡实例。
    * `associated_instances` - 与HaVip绑定的ECS实例
    * `create_time` - 创建时间
    * `description` - HaVip实例的描述，长度为2到256个字符。
    * `ha_vip_id` - HaVip实例的ID
    * `ha_vip_name` - HaVip实例的名称
    * `ip_address` - HaVip的私网IP地址
    * `master_instance_id` - 与HaVip绑定的主实例ID
    * `status` - 代表资源实例的状态
    * `tags` - 资源标签。
    * `vswitch_id` - HaVip实例所属的VSwitch ID
    * `vpc_id` - HaVip实例所属的VPC ID
