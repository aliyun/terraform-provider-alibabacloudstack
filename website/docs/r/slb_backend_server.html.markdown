---
subcategory: "Server Load Balancer (SLB)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_slb_backend_server"
sidebar_current: "docs-alibabacloudstack-resource-slb-backend-server"
description: |-
  Provides an Application Load Balancer Backend Server resource.
---

# alibabacloudstack\_slb\_backend\_server

Add a group of backend servers (ECS or ENI instance) to the Server Load Balancer or remove them from it.

## Example Usage

```
variable "name" {
  default = "slbbackendservertest"
}
data "alibabacloudstack_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}
data "alibabacloudstack_instance_types" "default" {
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  cpu_core_count    = 1
  memory_size       = 2
}
data "alibabacloudstack_images" "default" {
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}

resource "alibabacloudstack_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = "${alibabacloudstack_vpc.default.id}"
  cidr_block        = "172.16.0.0/16"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  name              = "${var.name}"
}

resource "alibabacloudstack_security_group" "default" {
  name   = "${var.name}"
  vpc_id = "${alibabacloudstack_vpc.default.id}"
}

resource "alibabacloudstack_instance" "default" {
  image_id = "${data.alibabacloudstack_images.default.images.0.id}"
  instance_type = "${data.alibabacloudstack_instance_types.default.instance_types.0.id}"
  instance_name = "${var.name}"
  count = "2"
  security_groups = "${alibabacloudstack_security_group.default.*.id}"
  internet_max_bandwidth_out = "10"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  system_disk_category = "cloud_efficiency"
  vswitch_id = "${alibabacloudstack_vswitch.default.id}"
}

resource "alibabacloudstack_slb" "default" {
  name       = "${var.name}"
  vswitch_id = "${alibabacloudstack_vswitch.default.id}"
}

resource "alibabacloudstack_slb_backend_server" "default" {
  	load_balancer_id = "${alibabacloudstack_slb.default.id}"
  	
	backend_servers {
      server_id = "${alibabacloudstack_instance.default.0.id}"
      weight     = 100
    }

    backend_servers {
      server_id = "${alibabacloudstack_instance.default.1.id}"
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
