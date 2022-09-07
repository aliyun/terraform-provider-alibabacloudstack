---
subcategory: "VPN"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpn_route_entry"
sidebar_current: "docs-alibabacloudstack-resource-vpn-route-entry"
description: |-
  Provides a Alibabacloudstack VPN Route Entry resource.
---

# alibabacloudstack\_vpn_route_entry

Provides a VPN Route Entry resource.

-> **NOTE:** Terraform will build vpn route entry instance while it uses `alibabacloudstack_vpn_route_entry` to build a VPN Route Entry resource.

## Example Usage

Basic Usage

```
data "alibabacloudstack_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vpc" "default" {
  name       = "tf_test"
  cidr_block = "10.1.0.0/21"
}

resource "alibabacloudstack_vswitch" "default" {
  name              = "tf_test"
  vpc_id            = alibabacloudstack_vpc.default.id
  cidr_block        = "10.1.0.0/24"
  zone_id           = data.alibabacloudstack_zones.default.zones[0].id
}

resource "alibabacloudstack_vpn_gateway" "default" {
  name                 = "tf_vpn_gateway_test"
  vpc_id               = alibabacloudstack_vpc.default.id
  bandwidth            = 10
  instance_charge_type = "PayByTraffic"
  enable_ssl           = false
  vswitch_id           = alibabacloudstack_vswitch.default.id
}

resource "alibabacloudstack_vpn_connection" "default" {
  name                = "tf_vpn_connection_test"
  customer_gateway_id = alibabacloudstack_vpn_customer_gateway.default.id
  vpn_gateway_id      = alibabacloudstack_vpn_gateway.default.id
  local_subnet        = ["192.168.2.0/24"]
  remote_subnet       = ["192.168.3.0/24"]
}

resource "alibabacloudstack_vpn_customer_gateway" "default" {
  name       = "tf_customer_gateway_test"
  ip_address = "192.168.1.1"
}

resource "alibabacloudstack_vpn_route_entry" "default" {
  vpn_gateway_id = alibabacloudstack_vpn_gateway.default.id
  route_dest     = "10.0.0.0/24"
  next_hop       = alibabacloudstack_vpn_connection.default.id
  weight         = 0
  publish_vpc    = false
}
```
## Argument Reference

The following arguments are supported:

* `vpn_gateway_id` - (Required, ForceNew) The id of the vpn gateway.
* `next_hop` - (Required, ForceNew) The next hop of the destination route.
* `publish_vpc` - (Required) Whether to issue the destination route to the VPC.
* `route_dest` - (Required, ForceNew) The destination network segment of the destination route.
* `weight` - (Required) The value should be 0 or 100.

## Attributes Reference

The following attributes are exported:

* `id` - The combination id of the vpn route entry.

## Import

VPN route entry can be imported using the id(VpnGatewayId +":"+ NextHop +":"+ RouteDest), e.g.

```
$ terraform import alibabacloudstack_vpn_route_entry.example vpn-abc123456:vco-abc123456:10.0.0.10/24
```
