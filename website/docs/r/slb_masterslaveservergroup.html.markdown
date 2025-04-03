---
subcategory: "SLB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_slb_masterslaveservergroup"
sidebar_current: "docs-Alibabacloudstack-slb-masterslaveservergroup"
description: |- 
  Provides a slb Masterslaveservergroup resource.
---

# alibabacloudstack_slb_masterslaveservergroup
-> **NOTE:** Alias name has: `alibabacloudstack_slb_master_slave_server_group`

Provides a slb Masterslaveservergroup resource.

## Example Usage

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

## Argument Reference

The following arguments are supported:

* `load_balancer_id` - (Required, ForceNew) The ID of the SLB instance.
* `name` - (Optional, ForceNew) The name of the master-slave server group. It must be unique within the specified Load Balancer.
* `master_slave_server_group_name` - (Optional, ForceNew) The name of the primary/secondary server group. If not provided, it defaults to the value of `name`.
* `servers` - (Optional, ForceNew) A list of ECS instances to be added as backend servers in the master-slave server group. Only two ECS instances can be supported in one resource. Each server contains the following sub-fields:
  * `server_id` - (Required) The ID of the ECS instance to be added as a backend server.
  * `port` - (Required) The port used by the backend server. Valid value range: [1-65535].
  * `weight` - (Optional) Weight of the backend server. Valid value range: [0-100]. Default to 100.
  * `server_type` - (Optional) The type of the backend server. Valid values: `Master`, `Slave`. Defaults to `Master`.
* `delete_protection_validation` - (Optional) Checking DeleteProtection of SLB instance before deleting. If set to `true`, this resource will not be deleted when its SLB instance has enabled DeleteProtection. Default to `false`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the master-slave server group.
* `name` - The name of the master-slave server group.
* `master_slave_server_group_name` - The name of the primary/secondary server group.