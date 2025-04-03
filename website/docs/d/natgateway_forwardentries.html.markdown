---
subcategory: "NATGateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_natgateway_forwardentries"
sidebar_current: "docs-Alibabacloudstack-datasource-natgateway-forwardentries"
description: |- 
  Provides a list of natgateway forwardentries owned by an AlibabaCloudStack account.
---

# alibabacloudstack_natgateway_forwardentries
-> **NOTE:** Alias name has: `alibabacloudstack_forward_entries`

This data source provides a list of NAT Gateway Forward Entries in an AlibabaCloudStack account according to the specified filters.

## Example Usage

```hcl
variable "name" {
  default = "natgateway-forward-entry-example-name"
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
  internal_ip     = "172.16.0.3"
  internal_port   = "8080"
}

data "alibabacloudstack_natgateway_forwardentries" "default" {
  forward_table_id = "${alibabacloudstack_forward_entry.default.forward_table_id}"
  external_ip      = "${alibabacloudstack_eip.default.ip_address}"
  internal_ip      = "172.16.0.3"
  name_regex       = "example.*"

  output_file = "forward_entries_output.txt"
}

output "natgateway_forward_entries" {
  value = "${data.alibabacloudstack_natgateway_forwardentries.default.entries}"
}
```

## Argument Reference

The following arguments are supported:

* `forward_table_id` - (Required, ForceNew) The ID of the DNAT table to which the DNAT entry belongs.
* `name_regex` - (Optional) A regex string used to filter results by the name of the Forward Entry.
* `external_ip` - (Optional) The public IP address in the DNAT entry. The public IP address is used by the ECS instance to receive requests from the Internet.
* `internal_ip` - (Optional) The private IP address that is mapped to the public IP address in the DNAT entry.
* `ids` - (Optional) A list of Forward Entry IDs.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Forward Entry IDs.
* `names` - A list of Forward Entry names.
* `entries` - A list of Forward Entries. Each element contains the following attributes:
  * `id` - The ID of the Forward Entry.
  * `external_ip` - The public IP address in the DNAT entry. The public IP address is used by the ECS instance to receive requests from the Internet.
  * `internal_ip` - The private IP address that is mapped to the public IP address in the DNAT entry.
  * `external_port` - The external port in the DNAT entry. The external port is used by the ECS instance to receive requests from the Internet.
  * `internal_port` - The internal port that is mapped to the external port in the DNAT entry.
  * `ip_protocol` - The type of the protocol.
  * `status` - The state of the DNAT entry.