---
subcategory: "VPNGateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpngateway_sslvpnclientcert"
description: |-
  Provides a vpngateway Sslvpnclientcert resource.
---

# alibabacloudstack\_vpngateway\_sslvpnclientcert

Provides a vpngateway Sslvpnclientcert resource.

## Example Usage
```
variable "name" {
    default = "tf-testaccvpn_gateway_sslVpnClient17243"
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
	client_ip_pool =  "10.8.0.0/24"

	local_subnet =  "192.168.1.0/24"

	ssl_vpn_server_name =  "${var.name}"
	vpn_gateway_id =     "${alibabacloudstack_vpn_gateway.default.id}"
	}





resource "alibabacloudstack_vpngateway_sslvpnclientcert" "default" {
  ssl_vpn_client_cert_name = "tf-testaccvpn_gateway_sslVpnClient17243"
  ssl_vpn_server_id = "${alibabacloudstack_vpngateway_ssl_vpnserver.default.id}"
  vpn_gateway_id = "${alibabacloudstack_vpn_gateway.default.id}"
}
```

## Argument Reference

The following arguments are supported:
  * `ssl_vpn_client_cert_name` - (Required) - The name of the client certificate.
  * `vpn_gateway_id` - (Required, ForceNew) - The ID of the VPN gateway. 
  * `ssl_vpn_server_id` - (Required, ForceNew) - The ID of the SSL server.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `ca_cert` - The CA certificate.
  * `client_cert` - The client certificate.
  * `client_config` - The client configuration.
  * `client_key` - The client key.
  * `create_time` - The time when the SSL client certificate was created.
  * `end_time` - The time when the SSL client certificate expires.
  * `ssl_vpn_client_cert_id` - The ID of the SSL client certificate.
  * `status` - The status of the client certificate.
