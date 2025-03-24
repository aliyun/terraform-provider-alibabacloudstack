---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_route_entries"
sidebar_current: "docs-alibabacloudstack-datasource-route-entries"
description: |-
    Provides a list of Route Entries owned by an Alibabacloudstack Cloud account.
---

# alibabacloudstack_route_entries

This data source provides a list of Route Entries owned by an Alibabacloudstack Cloud account.


## Example Usage

```
data "alibabacloudstack_zones" "default" {
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

variable "name" {
  default = "tf-testAccRouteEntryConfig"
}
resource "alibabacloudstack_vpc" "foo" {
  name       = "${var.name}"
  cidr_block = "10.1.0.0/21"
}

resource "alibabacloudstack_vswitch" "foo" {
  vpc_id            = "${alibabacloudstack_vpc.foo.id}"
  cidr_block        = "10.1.1.0/24"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  name              = "${var.name}"
}

resource "alibabacloudstack_route_entry" "foo" {
  route_table_id        = "${alibabacloudstack_vpc.foo.route_table_id}"
  destination_cidrblock = "172.11.1.1/32"
  nexthop_type          = "Instance"
  nexthop_id            = "${alibabacloudstack_instance.foo.id}"
}

resource "alibabacloudstack_security_group" "tf_test_foo" {
  name        = "${var.name}"
  description = "foo"
  vpc_id      = "${alibabacloudstack_vpc.foo.id}"
}

resource "alibabacloudstack_security_group_rule" "ingress" {
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "22/22"
  priority          = 1
  security_group_id = "${alibabacloudstack_security_group.tf_test_foo.id}"
  cidr_ip           = "0.0.0.0/0"
}

resource "alibabacloudstack_instance" "foo" {
  security_groups = ["${alibabacloudstack_security_group.tf_test_foo.id}"]

  vswitch_id         = "${alibabacloudstack_vswitch.foo.id}"
  allocate_public_ip = true

  # series III
  instance_type              = "${data.alibabacloudstack_instance_types.default.instance_types.0.id}"
  internet_max_bandwidth_out = 5

  system_disk_category = "cloud_efficiency"
  image_id             = "${data.alibabacloudstack_images.default.images.0.id}"
  instance_name        = "${var.name}"
}

data "alibabacloudstack_route_entries" "foo" {
  route_table_id = "${alibabacloudstack_route_entry.foo.route_table_id}"
}

output "route_entries" {
 data = data.alibabacloudstack_route_entries.foo
}

```

## Argument Reference

The following arguments are supported:

* `route_table_id` - (Required, ForceNew) The ID of the router table to which the route entry belongs.
* `instance_id` - (Optional) The instance ID of the next hop.
* `type` - (Optional) The type of the route entry.
* `cidr_block` - (Optional) The destination CIDR block of the route entry.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `entries` - A list of Route Entries. Each element contains the following attributes:
  * `type` - The type of the route entry.
  * `next_hop_type` - The type of the next hop.
  * `status` - The status of the route entry.
  * `instance_id` - The instance ID of the next hop.
  * `route_table_id` - The ID of the router table to which the route entry belongs.
  * `cidr_block` - The destination CIDR block of the route entry.