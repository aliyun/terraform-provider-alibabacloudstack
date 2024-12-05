---
subcategory: "Server Load Balancer (SLB)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_slb_master_slave_server_group"
sidebar_current: "docs-alibabacloudstack-resource-slb-master-slave-server-group"
description: |-
  Provides a Load Balancer Master Slave Server Group resource.
---

# alibabacloudstack\_slb\_master\_slave\_server\_group

A master slave server group contains two ECS instances. The master slave server group can help you to define multiple listening dimension.

-> **NOTE:** One ECS instance can be added into multiple master slave server groups.

-> **NOTE:** One master slave server group can only add two ECS instances, which are master server and slave server.

-> **NOTE:** One master slave server group can be attached with tcp/udp listeners in one load balancer.

-> **NOTE:** One Classic and Internet load balancer, its master slave server group can add Classic and VPC ECS instances.

-> **NOTE:** One Classic and Intranet load balancer, its master slave server group can only add Classic ECS instances.

-> **NOTE:** One VPC load balancer, its master slave server group can only add the same VPC ECS instances.
## Example Usage

```
data "alibabacloudstack_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alibabacloudstack_instance_types" "default" {
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  eni_amount        = 2
}

data "alibabacloudstack_images" "image" {
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}

variable "name" {
  default = "tf-testAccSlbMasterSlaveServerGroupVpc"
}

variable "number" {
  default = "1"
}

resource "alibabacloudstack_vpc" "main" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vswitch" "main" {
  vpc_id            = "${alibabacloudstack_vpc.main.id}"
  cidr_block        = "172.16.0.0/16"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  name              = "${var.name}"
}

resource "alibabacloudstack_security_group" "group" {
  name   = "${var.name}"
  vpc_id = "${alibabacloudstack_vpc.main.id}"
}

resource "alibabacloudstack_instance" "instance" {
  image_id                   = "${data.alibabacloudstack_images.image.images.0.id}"
  instance_type              = "${data.alibabacloudstack_instance_types.default.instance_types.0.id}"
  instance_name              = "${var.name}"
  count                      = "2"
  security_groups            = ["${alibabacloudstack_security_group.group.id}"]
  internet_max_bandwidth_out = "10"
  availability_zone          = "${data.alibabacloudstack_zones.default.zones.0.id}"
  system_disk_category       = "cloud_efficiency"
  vswitch_id                 = "${alibabacloudstack_vswitch.main.id}"
}

resource "alibabacloudstack_slb" "instance" {
  name          = "${var.name}"
  vswitch_id    = "${alibabacloudstack_vswitch.main.id}"
}

resource "alibabacloudstack_network_interface" "default" {
  count           = "${var.number}"
  name            = "${var.name}"
  vswitch_id      = "${alibabacloudstack_vswitch.main.id}"
  security_groups = ["${alibabacloudstack_security_group.group.id}"]
}

resource "alibabacloudstack_network_interface_attachment" "default" {
  count                = "${var.number}"
  instance_id          = "${alibabacloudstack_instance.instance.0.id}"
  network_interface_id = "${element(alibabacloudstack_network_interface.default.*.id, count.index)}"
}

resource "alibabacloudstack_slb_master_slave_server_group" "group" {
  load_balancer_id = "${alibabacloudstack_slb.instance.id}"
  name             = "${var.name}"

  servers {
    server_id = "${alibabacloudstack_instance.instance.0.id}"
    port       = 100
    weight     = 100
    server_type = "Master"
  }

  servers {
    server_id = "${alibabacloudstack_instance.instance.1.id}"
    port       = 100
    weight     = 100
    server_type = "Slave"
  }
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_id` - (Required, ForceNew) The Load Balancer ID which is used to launch a new master slave server group.
* `name` - (Optional, ForceNew) Name of the master slave server group. 
* `delete_protection_validation` - (Optional) Checking DeleteProtection of SLB instance before deleting. If true, this resource will not be deleted when its SLB instance enabled DeleteProtection. Default to false.
* `servers` - (Optional, ForceNew) A list of ECS instances to be added. Only two ECS instances can be supported in one resource. It contains six sub-fields as `Block server` follows.The servers mapping supports the following:
  * `server_id` - (Required) backend server ID (ECS instance ID).
  * `port` - (Required) The port used by the backend server. Valid value range: [1-65535].
  * `weight` - (Optional) Weight of the backend server. Valid value range: [0-100]. Default to 100.
  * `server_type` - (Optional) The server type of the backend server. Valid value Master, Slave.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the master slave server group.

