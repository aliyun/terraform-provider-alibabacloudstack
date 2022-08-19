---
subcategory: "VPC"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_snat_entry"
sidebar_current: "docs-apsarastack-resource-vpc"
description: |-
  Provides a Apsarastack snat resource.
---

# apsarastack\_snat_entry

Provides a snat resource.

## Example Usage

Basic Usage

```
variable "name" {
  default = "snat-entry-example-name"
}
data "apsarastack_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "apsarastack_vpc" "vpc" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "apsarastack_vswitch" "vswitch" {
  vpc_id            = "${apsarastack_vpc.vpc.id}"
  cidr_block        = "172.16.0.0/21"
  availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
  name              = "${var.name}"
}

resource "apsarastack_nat_gateway" "default" {
  vpc_id        = "${apsarastack_vswitch.vswitch.vpc_id}"
  specification = "Small"
  name          = "${var.name}"
}

resource "apsarastack_eip" "default" {
  count = 2
  name  = "${var.name}"
}

resource "apsarastack_eip_association" "default" {
  count         = 2
  allocation_id = "${element(apsarastack_eip.default.*.id, count.index)}"
  instance_id   = "${apsarastack_nat_gateway.default.id}"
}

resource "apsarastack_snat_entry" "default" {
  depends_on        = [apsarastack_eip_association.default]
  snat_table_id     = "${apsarastack_nat_gateway.default.snat_table_ids}"
  source_vswitch_id = "${apsarastack_vswitch.vswitch.id}"
  snat_ip           = "${join(",", apsarastack_eip.default.*.ip_address)}"
}
```

## Argument Reference

The following arguments are supported:

* `snat_table_id` - (Required, ForceNew) The value can get from `apsarastack_nat_gateway` Attributes "snat_table_ids".
* `source_vswitch_id` - (Optional, ForceNew) The vswitch ID.
* `source_cidr` - (Optional, ForceNew) The private network segment of Ecs. This parameter and the `source_vswitch_id` parameter are mutually exclusive and cannot appear at the same time.
* `snat_ip` - (Required) The SNAT ip address, the ip must along bandwidth package public ip which `apsarastack_nat_gateway` argument `bandwidth_packages`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the snat entry. The value formats as `<snat_table_id>:<snat_entry_id>`
* `snat_entry_id` - The id of the snat entry on the server.


