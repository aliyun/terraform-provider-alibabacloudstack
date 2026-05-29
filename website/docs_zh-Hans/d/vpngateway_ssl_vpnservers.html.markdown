---
subcategory: "VPN网关"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpngateway_ssl_vpnservers"
sidebar_current: "docs-Alibabacloudstack-datasource-vpngateway-ssl-vpnservers"
description: |-
  提供当前 AlibabacloudStack 账户下所有 SSL VPN 服务器的列表。
---

# alibabacloudstack\_vpngateway\_sslvpnservers

该数据源用于根据指定过滤条件获取 VPNGateway 下的 SSL VPN 服务器列表。

## 示例用法

```hcl
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

## 参数说明

以下参数支持配置：
* `ids` - (可选) SSL VPN 服务器 ID 列表。若设置此参数，数据源将返回匹配的 SSL VPN 服务器。
* `name_regex` - (可选, ForceNew) 用于通过名称模式匹配 SSL VPN 服务器的正则表达式。可以使用该参数筛选出符合特定命名规则的 SSL VPN 服务器。

## 属性说明

除上述参数外，还将导出以下属性：
* `ssl_vpn_servers` - SSL VPN 服务器列表
  * `id` - SSL VPN 服务器的 ID
  * `cipher` - SSL-VPN 连接中使用的加密算法。
  * `client_ip_pool` - 客户端的 CIDR 块。
  * `compress` - 是否启用数据压缩。
  * `connections` - 当前连接总数。
  * `create_time` - SSL-VPN 服务器创建时间。
  * `internet_ip` - 公网 IP 地址。
  * `local_subnet` - 本地 CIDR 块。
  * `max_connections` - 最大连接数。
  * `port` - SSL-VPN 服务器使用的端口。
  * `proto` - SSL-VPN 服务器使用的协议。
  * `ssl_vpn_server_id` - SSL-VPN 服务器的 ID。
  * `ssl_vpn_server_name` - SSL-VPN 服务器的名称。
  * `vpn_gateway_id` - VPN 网关的 ID。