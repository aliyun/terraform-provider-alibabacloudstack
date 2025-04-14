---
subcategory: "SLB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_slb_masterslaveservergroup"
sidebar_current: "docs-Alibabacloudstack-slb-masterslaveservergroup"
description: |- 
  编排负载均衡(SLB)主备服务器组
---

# alibabacloudstack_slb_masterslaveservergroup
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_slb_master_slave_server_group`

使用Provider配置的凭证在指定的资源集编排负载均衡(SLB)主备服务器组。

## 示例用法

```hcl
variable "name" {
	default = "tf-testAccSlbMasterSlaveServerGroupVpc1592616"
}

data "alibabacloudstack_instance_types" "new" {
	availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
	eni_amount = 2
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

resource "alibabacloudstack_ecs_instance" "new" {
	image_id             = "${data.alibabacloudstack_images.default.images.0.id}"
	instance_type        = "${data.alibabacloudstack_instance_types.new.instance_types[0].id}"
	system_disk_category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
	system_disk_size     = 40
	system_disk_name     = "test_sys_diskv2"
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

resource "alibabacloudstack_network_interface" "default" {
	count = 1
	name = "${var.name}"
	vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
	security_groups = [ "${alibabacloudstack_ecs_securitygroup.default.id}" ]
}

resource "alibabacloudstack_network_interface_attachment" "default" {
	count = 1
	instance_id = "${alibabacloudstack_ecs_instance.new.id}"
	network_interface_id = "${element(alibabacloudstack_network_interface.default.*.id, count.index)}"
}

resource "alibabacloudstack_slb" "default" {
	name = "${var.name}"
	vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
}

resource "alibabacloudstack_slb_master_slave_server_group" "default" {
  name = "${var.name}"
  load_balancer_id = "${alibabacloudstack_slb.default.id}"

  servers {
    server_id = "${alibabacloudstack_ecs_instance.default.id}"
    port      = "100"
    weight    = "100"
    server_type = "Master"
  }

  servers {
    server_id = "${alibabacloudstack_ecs_instance.new.id}"
    port      = "100"
    weight    = "100"
    server_type = "Slave"
  }
}
```

## 参数说明

支持以下参数：

* `load_balancer_id` - (必填, 变更时重建) - 关联的负载均衡实例ID。
* `name` - (选填, 变更时重建) - 主备服务器组的名称。如果未提供，默认为`master_slave_server_group_name`的值。
* `master_slave_server_group_name` - (选填, 变更时重建) - 主备服务器组的名称。如果未提供，默认为`name`的值。
* `servers` - (选填, 变更时重建) - 要添加到主备服务器组中的ECS实例列表。一个资源中仅支持两个ECS实例。每个服务器包含以下子字段：
  * `server_id` - (必填) 要作为后端服务器添加的ECS实例的ID。
  * `port` - (必填) 后端服务器使用的端口。有效值范围：[1-65535]。
  * `weight` - (可选) 后端服务器的权重。有效值范围：[0-100]。默认为100。
  * `server_type` - (可选) 后端服务器的类型。有效值：`Master`，`Slave`。默认为`Master`。
* `delete_protection_validation` - (可选) 删除前检查负载均衡实例的删除保护。如果设置为`true`，当其负载均衡实例启用了删除保护时，此资源将不会被删除。默认为`false`。

## 属性说明

除了上述所有参数外，还导出了以下属性：

* `id` - 主备服务器组的ID。
* `name` - 主备服务器组的名称。
* `master_slave_server_group_name` - 主备服务器组的名称。