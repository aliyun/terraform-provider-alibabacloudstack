---
subcategory: "VPNGateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpngateway_sslvpnclientcert"
sidebar_current: "docs-Alibabacloudstack-vpngateway-sslvpnclientcert"
description: |-
  提供一个 VPNGateway 的 SSL VPN 客户端证书资源。
---

# alibabacloudstack\_vpngateway\_sslvpnclientcert

提供一个 VPNGateway 的 SSL VPN 客户端证书资源。

## 示例用法

```hcl
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

## 参数说明
以下参数支持配置：

* `ssl_vpn_client_cert_name` - (必需) - 客户端证书的名称。
* `vpn_gateway_id` - (必需, ForceNew) - VPN 网关的 ID。
* `ssl_vpn_server_id` - (必需, ForceNew) - SSL 服务器的 ID。
## 属性说明
除上述参数外，还将导出以下属性：

* `ca_cert` - CA 证书内容。
* `client_cert` - 客户端证书内容。
* `client_config` - 客户端配置信息。
* `client_key` - 客户端私钥。
* `create_time` - SSL 客户端证书创建时间。
* `end_time` - SSL 客户端证书过期时间。
* `ssl_vpn_client_cert_id` - SSL 客户端证书的 ID。
* `status` - 客户端证书的状态。