---
subcategory: "VPC"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_route_entries"
sidebar_current: "docs-apsarastack-datasource-route-entries"
description: |-
    Provides a list of Route Entries owned by an Apsarastack Cloud account.
---

# apsarastack\_route\_entries

This data source provides a list of Route Entries owned by an Apsarastack Cloud account.


## Example Usage

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
  default = "tf-testAccRouteEntryConfig"
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

resource "apsarastack_route_entry" "foo" {
  route_table_id        = "${apsarastack_vpc.foo.route_table_id}"
  destination_cidrblock = "172.11.1.1/32"
  nexthop_type          = "Instance"
  nexthop_id            = "${apsarastack_instance.foo.id}"
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

  vswitch_id         = "${apsarastack_vswitch.foo.id}"
  allocate_public_ip = true

  # series III
  instance_type              = "${data.apsarastack_instance_types.default.instance_types.0.id}"
  internet_max_bandwidth_out = 5

  system_disk_category = "cloud_efficiency"
  image_id             = "${data.apsarastack_images.default.images.0.id}"
  instance_name        = "${var.name}"
}

data "apsarastack_route_entries" "foo" {
  route_table_id = "${apsarastack_route_entry.foo.route_table_id}"
}

output "route_entries" {
 data = data.apsarastack_route_entries.foo
}

```

## Argument Reference

The following arguments are supported:

* `route_table_id` - (Required, ForceNew) The ID of the router table to which the route entry belongs.
* `instance_id` - (Optional) The instance ID of the next hop.
* `type` - (Optional) The type of the route entry.
* `cidr_block` - (Optional) The destination CIDR block of the route entry.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `entries` - A list of Route Entries. Each element contains the following attributes:
  * `type` - The type of the route entry.
  * `next_hop_type` - The type of the next hop.
  * `status` - The status of the route entry.
  * `instance_id` - The instance ID of the next hop.
  * `route_table_id` - The ID of the router table to which the route entry belongs.
  * `cidr_block` - The destination CIDR block of the route entry.

