---
subcategory: "VPNGateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpngateway_vpnrouteentry"
sidebar_current: "docs-Alibabacloudstack-vpngateway-vpnrouteentry"
description: |- 
  Provides a vpngateway Vpnrouteentry resource.
---

# alibabacloudstack_vpngateway_vpnrouteentry
-> **NOTE:** Alias name has: `alibabacloudstack_vpn_route_entry`

Provides a vpngateway Vpnrouteentry resource.

## Example Usage

```hcl
variable "route_dests" {
  default = ["10.1.0.0/24", "10.1.0.0/32"]
}

data "alibabacloudstack_zones" "default" {}

resource "alibabacloudstack_vpc" "default" {
  name        = "tf-testaccvpnRouteEntrybasic11064"
  cidr_block  = "10.1.0.0/21"
}

resource "alibabacloudstack_vswitch" "default" {
  name               = "${alibabacloudstack_vpc.default.name}"
  vpc_id             = "${alibabacloudstack_vpc.default.id}"
  cidr_block         = "10.1.1.0/24"
  availability_zone  = "${data.alibabacloudstack_zones.default.ids.0}"
}

resource "alibabacloudstack_vpn_gateway" "default" {
  name                 = "${alibabacloudstack_vpc.default.name}"
  vpc_id               = "${alibabacloudstack_vpc.default.id}"
  bandwidth            = 10
  instance_charge_type = "PostPaid"
  enable_ssl           = false
  vswitch_id           = "${alibabacloudstack_vswitch.default.id}"
}

resource "alibabacloudstack_vpn_connection" "default" {
  name                = "${alibabacloudstack_vpc.default.name}"
  customer_gateway_id = "${alibabacloudstack_vpn_customer_gateway.default.id}"
  vpn_gateway_id      = "${alibabacloudstack_vpn_gateway.default.id}"
  local_subnet        = ["192.168.2.0/24"]
  remote_subnet       = ["192.168.3.0/24"]
}

resource "alibabacloudstack_vpn_customer_gateway" "default" {
  name       = "${alibabacloudstack_vpc.default.name}"
  ip_address = "192.168.1.1"
}

resource "alibabacloudstack_vpn_route_entry" "default" {
  weight          = "100"
  publish_vpc     = "false"
  vpn_gateway_id  = "${alibabacloudstack_vpn_gateway.default.id}"
  route_dest      = "10.0.0.0/24"
  next_hop        = "${alibabacloudstack_vpn_connection.default.id}"
}
```

## Argument Reference

The following arguments are supported:

* `vpn_gateway_id` - (Required, ForceNew) The ID of the VPN Gateway.
* `next_hop` - (Required, ForceNew) The next hop of the destination route entry. This is typically the ID of a `alibabacloudstack_vpn_connection`.
* `route_dest` - (Required, ForceNew) The destination CIDR block of the destination route. This defines the network segment that the route applies to.
* `weight` - (Required) The weight of the destination route. It can be set to either `0` or `100`. A higher weight indicates a preferred route.
* `publish_vpc` - (Required) Whether to publish the destination route to the associated VPC. Set this to `true` if you want the route to be propagated to the VPC routing tables.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the VPN Route Entry. It is composed of `VpnGatewayId`, `NextHop`, and `RouteDest`.