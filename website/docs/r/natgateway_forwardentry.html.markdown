---
subcategory: "NATGateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_natgateway_forwardentry"
sidebar_current: "docs-Alibabacloudstack-natgateway-forwardentry"
description: |- 
  Provides a natgateway Forwardentry resource.
---

# alibabacloudstack_natgateway_forwardentry
-> **NOTE:** Alias name has: `alibabacloudstack_forward_entry`

Provides a natgateway Forwardentry resource.

## Example Usage

```hcl
variable "name" {
  default = "tf-testAccForwardEntryConfig17430"
}

variable "number" {
  default = "2"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vpc" "default" {
  cidr_block = "172.16.0.0/12"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = alibabacloudstack_vpc.default.id
  cidr_block        = "172.16.0.0/21"
  availability_zone = data.alibabacloudstack_zones.default.zones.0.id
}

resource "alibabacloudstack_nat_gateway" "default" {
  vpc_id        = alibabacloudstack_vswitch.default.vpc_id
  specification = "Small"
}

resource "alibabacloudstack_eip" "default" {
  count = var.number
}

resource "alibabacloudstack_eip_association" "default" {
  count          = var.number
  allocation_id  = alibabacloudstack_eip.default[count.index].id
  instance_id    = alibabacloudstack_nat_gateway.default.id
}

resource "alibabacloudstack_forward_entry" "default" {
  name              = var.name
  forward_table_id = alibabacloudstack_nat_gateway.default.forward_table_ids
  external_ip      = alibabacloudstack_eip.default[0].ip_address
  external_port    = "80"
  ip_protocol      = "tcp"
  internal_ip     = "172.16.0.4"
  internal_port   = "8080"
}
```

## Argument Reference

The following arguments are supported:

* `forward_table_id` - (Required, ForceNew) The ID of the DNAT table to which the DNAT entry belongs.
* `external_ip` - (Required) The public IP address in the DNAT entry. The public IP address is used by the ECS instance to receive requests from the Internet.
* `external_port` - (Required) The external port in the DNAT entry. The external port is used by the ECS instance to receive requests from the Internet. Valid values are integers between 1 and 65535 or "any".
* `ip_protocol` - (Required) The type of the protocol. Valid values are `tcp`, `udp`, or `any`.
* `name` - (Optional) The name of the DNAT entry. This field can be used interchangeably with `forward_entry_name`.
* `forward_entry_name` - (Optional) The name of the DNAT entry. If not provided, the `name` field will be used.
* `internal_ip` - (Required) The private IP address that is mapped to the public IP address in the DNAT entry. It must be a valid private IP within the VPC.
* `internal_port` - (Required) The internal port that is mapped to the external port in the DNAT entry. Valid values are integers between 1 and 65535 or "any".

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the DNAT entry. The value is formatted as `<forward_table_id>:<forward_entry_id>`.
* `forward_entry_id` - The unique identifier for the DNAT entry on the server.
* `forward_entry_name` - The name of the DNAT entry. If not explicitly set, it defaults to the value of the `name` field.