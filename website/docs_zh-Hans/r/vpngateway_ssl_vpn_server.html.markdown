---
subcategory: "VPN网关"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpngateway_sslvpnserver"
sidebar_current: "docs-Alibabacloudstack-vpngateway-sslvpnserver"
description: |-
  提供一个 VPNGateway 的 SSL-VPN 服务器资源。
---

# alibabacloudstack\_vpngateway\_sslvpnserver

提供一个 VPNGateway 的 SSL-VPN 服务器资源。

-> **Note:** 该资源也可以使用以下别名引用：
- `alibabacloudstack_vpngateway_ssl_vpn_server`
- `alibabacloudstack_vpngateway_ssl_vpnserver`

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

* `vpn_gateway_id` - (必需, ForceNew) VPN 网关的 ID。修改此参数会强制重新创建资源。
* `client_ip_pool` - (必需) 客户端 IP 地址段。VPN 网关从该地址段为 SSL-VPN 客户端分配 IP 地址。
* `local_subnet` - (必需) 本地子网 CIDR 段。客户端通过 SSL-VPN 连接需要访问的地址段。
* `ssl_vpn_server_name` - (必需) SSL-VPN 服务器的名称。长度为 2~100 个字符，不能以 `http://` 或 `https://` 开头。
* `cipher` - (可选) SSL-VPN 连接中使用的加密算法。取值：`AES-128-CBC`（默认值）、`AES-192-CBC`、`AES-256-CBC`、`none`。
* `proto` - (可选) SSL-VPN 服务器使用的协议。取值：`TCP`（默认值）、`UDP`。
* `port` - (可选) SSL-VPN 服务器使用的端口。取值范围：1~65535，默认值：1194。不支持以下端口：22、2222、22222、9000、9001、9002、7505、80、443、53、68、123、4510、4560、500、4500。
* `compress` - (可选) 是否对 SSL-VPN 连接启用数据压缩。取值：`true`、`false`（默认值）。
## 属性说明
除上述参数外，还将导出以下属性：

* `id` - SSL-VPN 服务器的 ID。
* `ssl_vpn_server_id` - SSL-VPN 服务器的 ID。
* `cipher` - SSL-VPN 连接中使用的加密算法。
* `proto` - SSL-VPN 服务器使用的协议。
* `port` - SSL-VPN 服务器使用的端口。
* `connections` - 当前连接数。
* `max_connections` - 最大连接数。
* `create_time` - SSL-VPN 服务器创建时间。
* `internet_ip` - SSL-VPN 服务器的公网 IP 地址。

## Import

SSL-VPN 服务器可以使用 `ssl_vpn_server_id` 导入，例如：

```
$ terraform import alibabacloudstack_vpngateway_ssl_vpn_server.example vss-12345678
```