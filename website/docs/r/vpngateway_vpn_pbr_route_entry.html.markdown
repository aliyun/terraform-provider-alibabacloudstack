---
subcategory: "VPNGateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpngateway_vpnpbrrouteentry"
sidebar_current: "docs-Alibabacloudstack-vpngateway-vpnpbrrouteentry"
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
  * `next_hop` - (Required, ForceNew) - The next hop of the destination route entry.
  * `overlay_mode` - (Optional, ForceNew) - Tunnel protocol. Value: **Ipsec**(IPsec tunnel protocol).
  * `publish_vpc` - (Required) - Whether to publish a policy route to a VPC. Value:-**true**: The publish policy is routed to the VPC.-**false**: does not publish the policy route to the VPC.
  * `route_dest` - (Required, ForceNew) - The destination CIDR block of the destination route.
  * `route_source` - (Required, ForceNew) - The source CIDR block of the policy route.
  * `vpn_gateway_id` - (Required, ForceNew) - The ID of the VPN Gateway.
  * `weight` - (Required) - The weight of the destination route.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `create_time` - The time when the VPN route was created.
  * `status` - The status of the VPN destination route.
