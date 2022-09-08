---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_forward_entries"
sidebar_current: "docs-alibabacloudstack-datasource-forward-entries"
description: |-
    Provides a list of Forward Entries owned by an Apsara Stack Cloud account.
---

# alibabacloudstack\_forward\_entries

This data source provides a list of Forward Entries owned by an Apsara Stack Cloud account.


## Example Usage

```
variable "name" {
  default = "forward-entry-config-example-name"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = "${alibabacloudstack_vpc.default.id}"
  cidr_block        = "172.16.0.0/21"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  name              = "${var.name}"
}

resource "alibabacloudstack_nat_gateway" "default" {
  vpc_id        = "${alibabacloudstack_vpc.default.id}"
  specification = "Small"
  name          = "${var.name}"
}

resource "alibabacloudstack_eip" "default" {
  name = "${var.name}"
}

resource "alibabacloudstack_eip_association" "default" {
  allocation_id = "${alibabacloudstack_eip.default.id}"
  instance_id   = "${alibabacloudstack_nat_gateway.default.id}"
}

resource "alibabacloudstack_forward_entry" "default" {
  forward_table_id = "${alibabacloudstack_nat_gateway.default.forward_table_ids}"
  external_ip      = "${alibabacloudstack_eip.default.ip_address}"
  external_port    = "80"
  ip_protocol      = "tcp"
  internal_ip      = "172.16.0.3"
  internal_port    = "8080"
}

data "alibabacloudstack_forward_entries" "default" {
  forward_table_id = "${alibabacloudstack_forward_entry.default.forward_table_id}"
}

output "forward_entries" {
  value = "${data.alibabacloudstack_forward_entries.default.entries}"
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of Forward Entries IDs.
* `name_regex` - (Optional) A regex string to filter results by forward entry name.
* `external_ip` - (Optional) The public IP address.
* `internal_ip` - (Optional) The private IP address.
* `forward_table_id` - (Required) The ID of the Forward table.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Forward Entries IDs.
* `names` - A list of Forward Entries names.
* `entries` - A list of Forward Entries. Each element contains the following attributes:
  * `id` - The ID of the Forward Entry.
  * `external_ip` - The public IP address.
  * `external_port` - The public port.
  * `ip_protocol` - The protocol type.
  * `internal_ip` - The private IP address.
  * `internal_port` - The private port.
  * `status` - The status of the Forward Entry.

