---
subcategory: "VPNGateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpngateway_vpn_pbr_route_entry"
sidebar_current: "docs-Alibabacloudstack-resource-vpngateway-vpn-pbr-route-entry"
description: |-
  Provides a vpngateway Vpnpbrrouteentry resource.
---

# alibabacloudstack\_vpngateway\_vpnpbrrouteentry

Provides a vpngateway Vpnpbrrouteentry resource.

## Example Usage
```
data "alibabacloudstack_zones" "default"{
}

resource "alibabacloudstack_vpc" "default" {
 name  = "tf-testaccVpngatewayVpnpbrrouteentrybasic11558"
 cidr_block = "10.1.0.0/21"
}
resource "alibabacloudstack_vswitch" "default" {
 name			   = "${alibabacloudstack_vpc.default.name}"
 vpc_id            = "${alibabacloudstack_vpc.default.id}"
 cidr_block        = "10.1.1.0/24"
 availability_zone = "${data.alibabacloudstack_zones.default.ids.0}"
}
resource "alibabacloudstack_vpn_gateway" "default" {
 name                 = "${alibabacloudstack_vpc.default.name}"
 vpc_id               = "${alibabacloudstack_vpc.default.id}"
 bandwidth            = 10
 instance_charge_type = "PostPaid"
 enable_ssl           = true
 enable_ipsec		  = true
 vswitch_id			  = "${alibabacloudstack_vswitch.default.id}"
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


resource "alibabacloudstack_vpngateway_vpn_pbr_route_entry" "default" {
  vpn_gateway_id = "${alibabacloudstack_vpn_gateway.default.id}"
  route_dest = "10.0.0.0/24"
  route_source = "192.168.0.0/24"
  next_hop = "${alibabacloudstack_vpn_connection.default.id}"
  weight = "100"
  publish_vpc = "false"
}
```

## Argument Reference

The following arguments are supported:
  * `vpn_gateway_id` - (Required, ForceNew) The ID of the VPN Gateway.
  * `route_source` - (Required, ForceNew) The source CIDR block of the policy route.
  * `route_dest` - (Required, ForceNew) The destination CIDR block of the destination route.
  * `next_hop` - (Required, ForceNew) The next hop of the destination route entry.
  * `weight` - (Required) The weight of the destination route. Valid values: `0` and `100`.
  * `publish_vpc` - (Required) Whether to publish the policy route to a VPC. Valid values: `true` and `false`.
  * `overlay_mode` - (Optional) The tunnel protocol. Default value: `Ipsec`. Valid value: `Ipsec` (IPsec tunnel protocol).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `id` - The ID of the resource. The format is `<vpn_gateway_id>_<route_source>_<route_dest>_<next_hop>`.
  * `create_time` - The time when the VPN PBR route entry was created.
  * `status` - The status of the VPN PBR route entry.

## Import

VPN PBR Route Entry can be imported using the composite ID, e.g.

```
$ terraform import alibabacloudstack_vpngateway_vpn_pbr_route_entry.example <vpn_gateway_id>_<route_source>_<route_dest>_<next_hop>
```
