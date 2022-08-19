---
subcategory: "VPC"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_route_entry"
sidebar_current: "docs-apsarastack-resource-route-entry"
description: |-
  Provides a Apsarastack Route Entry resource.
---

# apsarastack\_route\_entry

Provides a route entry resource. A route entry represents a route item of one VPC route table.

## Example Usage

Basic Usage

```
data "apsarastack_zones" "default" {
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

variable "name" {
  default = "RouteEntryConfig"
}
resource "apsarastack_vpc" "foo" {
  name       = "${var.name}"
  cidr_block = "10.1.0.0/21"
}

resource "apsarastack_vswitch" "foo" {
  vpc_id            = "${apsarastack_vpc.foo.id}"
  cidr_block        = "10.1.1.0/24"
  availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
  name              = "${var.name}"
}

resource "apsarastack_security_group" "tf_test_foo" {
  name        = "${var.name}"
  description = "foo"
  vpc_id      = "${apsarastack_vpc.foo.id}"
}

resource "apsarastack_security_group_rule" "ingress" {
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "22/22"
  priority          = 1
  security_group_id = "${apsarastack_security_group.tf_test_foo.id}"
  cidr_ip           = "0.0.0.0/0"
}

resource "apsarastack_instance" "foo" {
  security_groups = ["${apsarastack_security_group.tf_test_foo.id}"]

  vswitch_id = "${apsarastack_vswitch.foo.id}"
  instance_type              = "${data.apsarastack_instance_types.default.instance_types.0.id}"
  internet_max_bandwidth_out = 5
  system_disk_category = "cloud_efficiency"
  image_id             = "${data.apsarastack_images.default.images.0.id}"
  instance_name        = "${var.name}"
}
resource "apsarastack_route_entry" "foo" {
  route_table_id        = "${apsarastack_vpc.foo.route_table_id}"
  destination_cidrblock = "172.11.1.1/32"
  nexthop_type          = "Instance"
  nexthop_id            = "${apsarastack_instance.foo.id}"
}
```


## Argument Reference

The following arguments are supported:

* `route_table_id` - (Required, ForceNew) The ID of the route table.
* `destination_cidrblock` - (ForceNew) The RouteEntry's target network segment.
* `nexthop_type` - (ForceNew) The next hop type. Available values:
    - `Instance` (Default): Route the traffic destined for the destination CIDR block to an ECS instance in the VPC.
    - `RouterInterface`: Route the traffic destined for the destination CIDR block to a router interface.
    - `VpnGateway`: Route the traffic destined for the destination CIDR block to a VPN Gateway.
    - `HaVip`: Route the traffic destined for the destination CIDR block to an HAVIP.
    - `NetworkInterface`: Route the traffic destined for the destination CIDR block to an NetworkInterface.
    - `NatGateway`: Route the traffic destined for the destination CIDR block to an Nat Gateway.

* `nexthop_id` - (ForceNew) The route entry's next hop. ECS instance ID or VPC router interface ID.

## Attributes Reference

The following attributes are exported:

* `id` - The route entry id,it formats of `<route_table_id:router_id:destination_cidrblock:nexthop_type:nexthop_id>`.
* `route_table_id` - The ID of the route table.
* `destination_cidrblock` - The RouteEntry's target network segment.
* `nexthop_type` - The next hop type.
* `nexthop_id` - The route entry's next hop.


