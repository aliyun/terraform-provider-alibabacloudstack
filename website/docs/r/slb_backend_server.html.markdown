---
subcategory: "Server Load Balancer (SLB)"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_slb_backend_server"
sidebar_current: "docs-apsarastack-resource-slb-backend-server"
description: |-
  Provides an Application Load Balancer Backend Server resource.
---

# apsarastack\_slb\_backend\_server

Add a group of backend servers (ECS or ENI instance) to the Server Load Balancer or remove them from it.

## Example Usage

```
variable "name" {
  default = "slbbackendservertest"
}
data "apsarastack_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}
data "apsarastack_instance_types" "default" {
  availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
  cpu_core_count    = 1
  memory_size       = 2
}
data "apsarastack_images" "default" {
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}

resource "apsarastack_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}

resource "apsarastack_vswitch" "default" {
  vpc_id            = "${apsarastack_vpc.default.id}"
  cidr_block        = "172.16.0.0/16"
  availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
  name              = "${var.name}"
}

resource "apsarastack_security_group" "default" {
  name   = "${var.name}"
  vpc_id = "${apsarastack_vpc.default.id}"
}

resource "apsarastack_instance" "default" {
  image_id = "${data.apsarastack_images.default.images.0.id}"
  instance_type = "${data.apsarastack_instance_types.default.instance_types.0.id}"
  instance_name = "${var.name}"
  count = "2"
  security_groups = "${apsarastack_security_group.default.*.id}"
  internet_max_bandwidth_out = "10"
  availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
  system_disk_category = "cloud_efficiency"
  vswitch_id = "${apsarastack_vswitch.default.id}"
}

resource "apsarastack_slb" "default" {
  name       = "${var.name}"
  vswitch_id = "${apsarastack_vswitch.default.id}"
}

resource "apsarastack_slb_backend_server" "default" {
  	load_balancer_id = "${apsarastack_slb.default.id}"
  	
	backend_servers {
      server_id = "${apsarastack_instance.default.0.id}"
      weight     = 100
    }

    backend_servers {
      server_id = "${apsarastack_instance.default.1.id}"
      weight     = 100
    }
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_id` - (Required) ID of the load balancer.
* `backend_servers` - (Required) A list of instances to added backend server in the SLB. It contains two sub-fields as `Block server` follows.
* `delete_protection_validation` - (Optional) Checking DeleteProtection of SLB instance before deleting. If true, this resource will not be deleted when its SLB instance enabled DeleteProtection. Default to false.

## Block servers

The servers mapping supports the following:

* `server_id` - (Required) A list backend server ID (ECS instance ID).
* `weight` - (Optional) Weight of the backend server. Valid value range: [0-100]. 

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource and the value same as load balancer id.
