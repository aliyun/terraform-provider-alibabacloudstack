---
subcategory: "VPNGateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpngateway-sslvpnserver"
sidebar_current: "docs-Alibabacloudstack-vpngateway-sslvpnserver"
description: |-
  提供一个 VPNGateway 的 SSL-VPN 服务器资源。
---

# alibabacloudstack\_vpngateway\_sslvpnserver

提供一个 VPNGateway 的 SSL-VPN 服务器资源。

## 示例用法

```hcl
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
## 参数说明
以下参数支持配置：

* `cipher` - (可选) - SSL-VPN 连接中使用的加密算法。
* `client_ip_pool` - (必需) - 客户端 IP 地址段。
* `compress` - (可选) - 是否启用数据压缩。
* `local_subnet` - (必需) - 本地子网 CIDR 段。
* `port` - (可选) - SSL-VPN 服务器使用的端口。
* `proto` - (可选) - SSL-VPN 服务器使用的协议。
* `ssl_vpn_server_name` - (必需) - SSL-VPN 服务器的名称。
* `vpn_gateway_id` - (必需, ForceNew) - VPN 网关的 ID。
## 属性说明
除上述参数外，还将导出以下属性：

* `cipher` - SSL-VPN 连接中使用的加密算法。
* `connections` - 当前连接数。
* `create_time` - SSL-VPN 服务器创建时间。
* `internet_ip` - 公网 IP 地址。
* `max_connections` - 最大连接数。
* `port` - SSL-VPN 服务器使用的端口。
* `proto` - SSL-VPN 服务器使用的协议。
* `ssl_vpn_server_id` - SSL-VPN 服务器的 ID。