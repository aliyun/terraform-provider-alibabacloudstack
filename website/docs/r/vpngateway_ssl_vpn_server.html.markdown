---
subcategory: "Virtual Private Cloud (VPC)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpngateway_sslvpnserver"
sidebar_current: "docs-Alibabacloudstack-vpngateway-sslvpnserver"
description: |-
  Provides a vpngateway Sslvpnserver resource.
---

# alibabacloudstack\_vpngateway\_sslvpnserver

Provides a vpngateway Sslvpnserver resource.

-> **Note:** This resource can also be referred to by the following aliases:
- `alibabacloudstack_vpngateway_ssl_vpn_server`
- `alibabacloudstack_vpngateway_ssl_vpnserver`

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

* `vpn_gateway_id` - (Required, ForceNew) The ID of the VPN gateway. Modifying this parameter will force a new resource to be created.
* `client_ip_pool` - (Required) The client CIDR block. This is the CIDR block from which the VPN gateway allocates IP addresses to clients connecting via SSL-VPN.
* `local_subnet` - (Required) The local CIDR block. This is the CIDR block that clients need to access through the SSL-VPN connection.
* `ssl_vpn_server_name` - (Required) The name of the SSL-VPN server. The length is 2~100 characters. It cannot start with `http://` or `https://`.
* `cipher` - (Optional) The encryption algorithm used in the SSL-VPN connection. Valid values: `AES-128-CBC` (default), `AES-192-CBC`, `AES-256-CBC`, `none`.
* `proto` - (Optional) The protocol used by the SSL-VPN server. Valid values: `TCP` (default), `UDP`.
* `port` - (Optional) The port used by the SSL-VPN server. Valid values: 1~65535. Default value: 1194. The following ports are not supported: 22, 2222, 22222, 9000, 9001, 9002, 7505, 80, 443, 53, 68, 123, 4510, 4560, 500, 4500.
* `compress` - (Optional) Specifies whether to enable data compression on the SSL-VPN connection. Valid values: `true`, `false` (default).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the SSL-VPN server.
* `ssl_vpn_server_id` - The ID of the SSL-VPN server.
* `cipher` - The encryption algorithm used in the SSL-VPN connection.
* `proto` - The protocol used by the SSL-VPN server.
* `port` - The port used by the SSL-VPN server.
* `connections` - The total number of current connections.
* `max_connections` - The maximum number of connections.
* `create_time` - The time when the SSL-VPN server was created.
* `internet_ip` - The public IP address of the SSL-VPN server.

## Import

SSL-VPN Server can be imported using the `ssl_vpn_server_id`, e.g.

```
$ terraform import alibabacloudstack_vpngateway_ssl_vpn_server.example vss-12345678
```
