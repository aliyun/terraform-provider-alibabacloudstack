---
subcategory: "NATGateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_natgateway_snatentries"
sidebar_current: "docs-Alibabacloudstack-datasource-natgateway-snatentries"
description: |- 
  Provides a list of natgateway snatentries owned by an alibabacloudstack account.
---

# alibabacloudstack_natgateway_snatentries
-> **NOTE:** Alias name has: `alibabacloudstack_snat_entries`

This data source provides a list of natgateway snatentries in an alibabacloudstack account according to the specified filters.

## Example Usage

```hcl
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

data "alibabacloudstack_natgateway_snatentries" "foo" {
  snat_table_id = "${alibabacloudstack_snat_entry.foo.snat_table_id}"
  source_cidr   = "172.16.0.0/21"
}

output "snat_entries" {
  value = "${data.alibabacloudstack_natgateway_snatentries.foo.entries}"
}
```

## Argument Reference

The following arguments are supported:

* `snat_table_id` - (Required, ForceNew) The SNAT table ID to which the SNAT entry belongs.
* `source_cidr` - (Optional) The source network segment of a SNAT entry.
* `ids` - (Optional) A list of Snat Entry IDs. This can be used to filter the results.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Snat Entry IDs.
* `entries` - A list of Snat Entries. Each element contains the following attributes:
  * `id` - The ID of the Snat Entry.
  * `snat_ip` - The public IP of the SNAT entry.
  * `source_cidr` - The source network segment of the SNAT entry.
  * `status` - The status of the SNAT entry. Possible values include `available`, `pending`, and `inactive`.