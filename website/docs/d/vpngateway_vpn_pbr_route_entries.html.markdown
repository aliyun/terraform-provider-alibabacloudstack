---
subcategory: "VPNGateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpngateway_vpn_pbr_route_entries"
sidebar_current: "docs-Alibabacloudstack-datasource-vpngateway-vpn-pbr-route-entries"
description: |-
  Provides a list of vpngateway vpnpbrrouteentries owned by an alibabacloudstack account.
---

# alibabacloudstack\_vpngateway\_vpnpbrrouteentries

This data source provides a list of vpngateway vpnpbrrouteentries in an alibabacloudstack account according to the specified filters.

## Example Usage
```
data "alibabacloudstack_zones" "default"{
}

resource "alibabacloudstack_vpc" "default" {
 name  = "tf-testAccVpngatewayVpnPbrRouteEntriesDataSource-6376584"
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
	route_dest =     "10.0.0.0/24"
	route_source =   "192.168.0.0/24"
	next_hop =      "${alibabacloudstack_vpn_connection.default.id}"
	weight =       "100"
	publish_vpc =    "false"
}
 

data "alibabacloudstack_vpngateway_vpn_pbr_route_entrys" "default" {
  vpn_gateway_id = [
                     "${alibabacloudstack_vpngateway_vpn_pbr_route_entry.default.vpn_gateway_id}"
                   ]
}
```

## Argument Reference

The following arguments are supported:
  * `ids` - (Optional) - A list of vpngateway vpnpbrrouteentry IDs.
  * `vpn_gateway_id` - (Required) - The ID of the VPN Gateway.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `vpn_pbr_route_entries` - A list of vpngateway vpnpbrrouteentrys.
    * `id` - TThe ID of the vpngateway vpnpbrrouteentry.
    * `create_time` - The time when the VPN route was created.
    * `description` - Description information of the policy route.
    * `next_hop` - The next hop of the destination route entry.
    * `route_dest` - The destination CIDR block of the destination route.
    * `route_source` - The source CIDR block of the policy route.
    * `status` - The status of the VPN destination route.
    * `vpn_gateway_id` - The ID of the VPN Gateway.
    * `weight` - The weight of the destination route.
