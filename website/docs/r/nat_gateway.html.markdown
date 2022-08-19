---
subcategory: "VPC"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_nat_gateway"
sidebar_current: "docs-apsarastack-resource-nat-gateway"
description: |-
  Provides a resource to create a VPC NAT Gateway.
---

# apsarastack\_nat\_gateway

Provides a resource to create a VPC NAT Gateway.


## Example Usage

Basic usage

```
variable "name" {
  default = "natGatewayExampleName"
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
  vpc_id = "${apsarastack_vswitch.default.vpc_id}"
  name   = "${var.name}"
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required, ForceNew) The VPC ID.
* `specification` - (Optional) The specification of the nat gateway. Valid values are `Small`, `Middle` and `Large`. Default to `Small`. 
* `name` - (Optional) Name of the nat gateway. The value can have a string of 2 to 128 characters, must contain only alphanumeric characters or hyphens, such as "-",".","_", and must not begin or end with a hyphen, and must not begin with http:// or https://. Defaults to null.
* `description` - (Optional) Description of the nat gateway, This description can have a string of 2 to 256 characters, It cannot begin with http:// or https://. Defaults to null.
* `bandwidth_packages` - (Optional) A list of bandwidth packages for the nat gateway. Only support nat gateway created before 00:00 on November 4, 2017.

## Block bandwidth packages
The bandwidth package mapping supports the following:

* `ip_count` - (Required) The IP number of the current bandwidth package. Its value range from 1 to 50.
* `bandwidth` - (Required) The bandwidth value of the current bandwidth package. Its value range from 5 to 5000.
* `zone` - (Optional) The AZ for the current bandwidth. If this value is not specified, Terraform will set a random AZ.
* `public_ip_addresses` - (Computer) The public ip for bandwidth package. the public ip count equal `ip_count`, multi ip would complex with ",", such as "10.0.0.1,10.0.0.2".

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the nat gateway.
* `name` - The name of the nat gateway.
* `description` - The description of the nat gateway.
* `specification` - The specification of the nat gateway.
* `vpc_id` - The VPC ID for the nat gateway.
* `bandwidth_package_ids` - A list ID of the bandwidth packages, and split them with commas.
* `snat_table_ids` - The nat gateway will auto create a snap and forward item, the `snat_table_ids` is the created one.
* `forward_table_ids` - The nat gateway will auto create a snap and forward item, the `forward_table_ids` is the created one.


