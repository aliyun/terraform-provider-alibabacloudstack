---
subcategory: "VPC"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_forward_entry"
sidebar_current: "docs-apsarastack-resource-vpc"
description: |-
  Provides a Apsarastack forward resource.
---

# apsarastack\_forward\_entry

Provides a forward resource.

## Example Usage

Basic Usage

```
variable "name" {
  default = "forward-entry-example-name"
}

data "apsarastack_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "apsarastack_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "apsarastack_vswitch" "default" {
  vpc_id            = "${apsarastack_vpc.default.id}"
  cidr_block        = "172.16.0.0/21"
  availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
  name              = "${var.name}"
}

resource "apsarastack_nat_gateway" "default" {
  vpc_id        = "${apsarastack_vpc.default.id}"
  specification = "Small"
  name          = "${var.name}"
}

resource "apsarastack_eip" "default" {
  name = "${var.name}"
}

resource "apsarastack_eip_association" "default" {
  allocation_id = "${apsarastack_eip.default.id}"
  instance_id   = "${apsarastack_nat_gateway.default.id}"
}

resource "apsarastack_forward_entry" "default" {
  forward_table_id = "${apsarastack_nat_gateway.default.forward_table_ids}"
  external_ip      = "${apsarastack_eip.default.ip_address}"
  external_port    = "80"
  ip_protocol      = "tcp"
  internal_ip      = "172.16.0.3"
  internal_port    = "8080"
}

```
## Argument Reference

The following arguments are supported:

* `forward_table_id` - (Required, ForceNew) The value can get from `apsarastack_nat_gateway` Attributes "forward_table_ids".
* `external_ip` - (Required, ForceNew) The external ip address, the ip must along bandwidth package public ip which `apsarastack_nat_gateway` argument `bandwidth_packages`.
* `external_port` - (Required) The external port, valid value is 1~65535|any.
* `ip_protocol` - (Required) The ip protocal, valid value is tcp|udp|any.
* `internal_ip` - (Required) The internal ip, must a private ip.
* `internal_port` - (Required) The internal port, valid value is 1~65535|any.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the forward entry. The value formats as `<forward_table_id>:<forward_entry_id>`
* `forward_entry_id` - The id of the forward entry on the server.
