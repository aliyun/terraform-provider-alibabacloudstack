---
subcategory: "VPNGateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpngateway_sslvpnserver"
sidebar_current: "docs-Alibabacloudstack-vpngateway-sslvpnserver"
description: |-
  Provides a vpngateway Sslvpnserver resource.
---

# alibabacloudstack\_vpngateway\_sslvpnserver

Provides a vpngateway Sslvpnserver resource.

## Example Usage
```
variable "name" {
    default = "tf-testacc-vpn_sslvpnserver83448"
}


data "alibabacloudstack_zones" default {
  available_resource_creation = "VSwitch"
  enable_details = true
}


resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  name = "${var.name}_vsw"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block = "172.16.1.0/24"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
}


resource "alibabacloudstack_vpn_gateway" "default" {
 name                 = "${var.name}"
 vpc_id               = "${alibabacloudstack_vpc_vpc.default.id}"
 bandwidth            = 10
 instance_charge_type = "PostPaid"
 enable_ssl           = true
 enable_ipsec		  = true
 vswitch_id			  = "${alibabacloudstack_vpc_vswitch.default.id}"
}




resource "alibabacloudstack_vpngateway_ssl_vpnserver" "default" {
  local_subnet = "192.168.1.0/24"
  ssl_vpn_server_name = "tf-testacc-vpn_sslvpnserver83448"
  vpn_gateway_id = "${alibabacloudstack_vpn_gateway.default.id}"
  proto = "TCP"
  port = "1193"
  cipher = "AES-128-CBC"
  compress = "true"
  client_ip_pool = "10.8.0.0/24"
}
```

## Argument Reference

The following arguments are supported:
  * `cipher` - (Optional) - The encryption algorithm that is used in the SSL-VPN connection.
  * `client_ip_pool` - (Required) - The CIDR block of the client.
  * `compress` - (Optional) - Specifies whether to enable data compression
  * `local_subnet` - (Required) - The local CIDR block.
  * `port` - (Optional) - The port that is used by the SSL-VPN server.
  * `proto` - (Optional) - The protocol that is used by the SSL-VPN server.
  * `ssl_vpn_server_name` - (Required) - The name of the SSL-VPN server.
  * `vpn_gateway_id` - (Required, ForceNew) - The ID of the VPN gateway.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `cipher` - The encryption algorithm that is used in the SSL-VPN connection.
  * `connections` - The total number of current connections.
  * `create_time` - The time when the SSL-VPN server was created.
  * `internet_ip` - The public IP address.
  * `max_connections` - The maximum number of connections.
  * `port` - The port that is used by the SSL-VPN server.
  * `proto` - The protocol that is used by the SSL-VPN server.
  * `ssl_vpn_server_id` - The ID of the SSL-VPN server.
