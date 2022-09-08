---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_snat_entries"
sidebar_current: "docs-alibabacloudstack-datasource-snat-entries"
description: |-
    Provides a list of Snat Entries owned by an Apsara Stack Cloud account.
---

# alibabacloudstack\_snat\_entries

This data source provides a list of Snat Entries owned by an Apsara Stack Cloud account.

## Example Usage

```
variable "name" {
  default = "snat-entry-example-name"
}
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vpc" "foo" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alibabacloudstack_vswitch" "foo" {
  vpc_id            = "${alibabacloudstack_vpc.foo.id}"
  cidr_block        = "172.16.0.0/21"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  name              = "${var.name}"
}

resource "alibabacloudstack_nat_gateway" "foo" {
  vpc_id        = "${alibabacloudstack_vpc.foo.id}"
  specification = "Small"
  name          = "${var.name}"
}

resource "alibabacloudstack_eip" "foo" {
  name = "${var.name}"
}

resource "alibabacloudstack_eip_association" "foo" {
  allocation_id = "${alibabacloudstack_eip.foo.id}"
  instance_id   = "${alibabacloudstack_nat_gateway.foo.id}"
}

resource "alibabacloudstack_snat_entry" "foo" {
  snat_table_id     = "${alibabacloudstack_nat_gateway.foo.snat_table_ids}"
  source_vswitch_id = "${alibabacloudstack_vswitch.foo.id}"
  snat_ip           = "${alibabacloudstack_eip.foo.ip_address}"
}

data "alibabacloudstack_snat_entries" "foo" {
  snat_table_id = "${alibabacloudstack_snat_entry.foo.snat_table_id}"
}

output "snat_entries" {
  value = "${data.alibabacloudstack_snat_entries.foo.entries}"
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

