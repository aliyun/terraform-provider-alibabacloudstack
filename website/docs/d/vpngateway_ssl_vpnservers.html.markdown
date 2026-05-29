---
subcategory: "VPNGateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpngateway_ssl_vpnservers"
sidebar_current: "docs-Alibabacloudstack-datasource-vpngateway-ssl-vpn-servers"
description: |-
  Provides a list of vpngateway ssl vpn servers owned by an alibabacloudstack account.
---

# alibabacloudstack\_vpngateway\_ssl\_vpn\_servers

This data source provides a list of vpngateway sslvpnservers in an alibabacloudstack account according to the specified filters.

## Example Usage
```
variable "name" {
    default = "tf_testAccVpngatewaySslvpnserverDataSource_2773692"
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
ssl_vpn_server_name = "${var.name}"

vpn_gateway_id = "${alibabacloudstack_vpn_gateway.default.id}"

local_subnet = "192.168.1.0/24"

client_ip_pool = "10.8.0.0/24"
}




data "alibabacloudstack_vpngateway_ssl_vpnservers" "default" {
  ids = [
          "${alibabacloudstack_vpngateway_ssl_vpnserver.default.id}"
        ]
}
```

## Argument Reference

The following arguments are supported:
* `ids` - (Optional) A list of ssl vpn servers IDs. If specified, the data source will return matching ssl vpn servers.
* `name_regex` - (Optional, ForceNew) A regex string used to filter ssl vpn servers by their names. This allows you to match specific patterns in the names of the ssl vpn servers.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `ssl_vpn_servers` - the list of SSL VPN SERVERS
    * `id` - the ID of the SSL VPN SERVER
    * `cipher` - The encryption algorithm that is used in the SSL-VPN connection.
    * `client_ip_pool` - The CIDR block of the client.
    * `compress` - Specifies whether to enable data compression
    * `connections` - The total number of current connections.
    * `create_time` - The time when the SSL-VPN server was created.
    * `internet_ip` - The public IP address.
    * `local_subnet` - The local CIDR block.
    * `max_connections` - The maximum number of connections.
    * `port` - The port that is used by the SSL-VPN server.
    * `proto` - The protocol that is used by the SSL-VPN server.
    * `ssl_vpn_server_id` - The ID of the SSL-VPN server.
    * `ssl_vpn_server_name` - The name of the SSL-VPN server.
    * `vpn_gateway_id` - The ID of the VPN gateway.
