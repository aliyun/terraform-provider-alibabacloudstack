---
subcategory: "VPC"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_snat_entries"
sidebar_current: "docs-apsarastack-datasource-snat-entries"
description: |-
    Provides a list of Snat Entries owned by an Apsara Stack Cloud account.
---

# apsarastack\_snat\_entries

This data source provides a list of Snat Entries owned by an Apsara Stack Cloud account.

## Example Usage

```
variable "name" {
  default = "snat-entry-example-name"
}
data "apsarastack_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "apsarastack_vpc" "foo" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "apsarastack_vswitch" "foo" {
  vpc_id            = "${apsarastack_vpc.foo.id}"
  cidr_block        = "172.16.0.0/21"
  availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
  name              = "${var.name}"
}

resource "apsarastack_nat_gateway" "foo" {
  vpc_id        = "${apsarastack_vpc.foo.id}"
  specification = "Small"
  name          = "${var.name}"
}

resource "apsarastack_eip" "foo" {
  name = "${var.name}"
}

resource "apsarastack_eip_association" "foo" {
  allocation_id = "${apsarastack_eip.foo.id}"
  instance_id   = "${apsarastack_nat_gateway.foo.id}"
}

resource "apsarastack_snat_entry" "foo" {
  snat_table_id     = "${apsarastack_nat_gateway.foo.snat_table_ids}"
  source_vswitch_id = "${apsarastack_vswitch.foo.id}"
  snat_ip           = "${apsarastack_eip.foo.ip_address}"
}

data "apsarastack_snat_entries" "foo" {
  snat_table_id = "${apsarastack_snat_entry.foo.snat_table_id}"
}

output "snat_entries" {
  value = "${data.apsarastack_snat_entries.foo.entries}"
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of Snat Entries IDs.
* `source_cidr` - (Optional) The source CIDR block of the Snat Entry.
* `snat_table_id` - (Required) The ID of the Snat table.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - (Optional) A list of Snat Entries IDs.
* `entries` - A list of Snat Entries. Each element contains the following attributes:
  * `id` - The ID of the Snat Entry.
  * `snat_ip` - The public IP of the Snat Entry.
  * `source_cidr` - The source CIDR block of the Snat Entry.
  * `status` - The status of the Snat Entry.

