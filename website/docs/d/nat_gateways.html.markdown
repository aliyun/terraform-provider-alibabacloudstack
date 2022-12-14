---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_nat_gateways"
sidebar_current: "docs-alibabacloudstack-datasource-nat-gateways"
description: |-
    Provides a list of Nat Gateways owned by an Alibabacloudstack Cloud account.
---

# alibabacloudstack\_nat\_gateways

This data source provides a list of Nat Gateways owned by an Alibabacloudstack Cloud account.



## Example Usage

```
variable "name" {
  default = "natGatewaysDatasource"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vpc" "foo" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alibabacloudstack_nat_gateway" "foo" {
  vpc_id        = "${alibabacloudstack_vpc.foo.id}"
  specification = "Small"
  name          = "${var.name}"
}

data "alibabacloudstack_nat_gateways" "foo" {
  vpc_id     = "${alibabacloudstack_vpc.foo.id}"
  name_regex = "${alibabacloudstack_nat_gateway.foo.name}"
  ids        = ["${alibabacloudstack_nat_gateway.foo.id}"]
}

output "nat_gateways" {
  value = "${data.alibabacloudstack_nat_gateways.foo.gateways}"
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of NAT gateways IDs.
* `name_regex` - (Optional) A regex string to filter nat gateways by name.
* `vpc_id` - (Optional) The ID of the VPC.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - (Optional) A list of Nat gateways IDs.
* `names` - A list of Nat gateways names.
* `gateways` - A list of Nat gateways. Each element contains the following attributes:
  * `id` - The ID of the NAT gateway.
  * `name` - Name of the NAT gateway.
  * `description` - The description of the NAT gateway.
  * `creation_time` - Time of creation.
  * `spec` - The specification of the NAT gateway.
  * `status` - The status of the NAT gateway.
  * `snat_table_id` - The snat table id.
  * `forward_table_id` - The forward table id. 

