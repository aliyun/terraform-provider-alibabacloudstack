---
subcategory: "VPNGateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpngateway_vpnrouteentry"
sidebar_current: "docs-Alibabacloudstack-vpngateway-vpnrouteentry"
description: |- 
  编排VPN网关VPN路由表
---

# alibabacloudstack_vpngateway_vpnrouteentry
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_vpn_route_entry`

使用Provider配置的凭证在指定的资源集编排VPN网关VPN路由表。

## 示例用法

```hcl
variable "route_dests" {
  default = ["10.1.0.0/24", "10.1.0.0/32"]
}

data "alibabacloudstack_zones" "default" {}

resource "alibabacloudstack_vpc" "default" {
  name        = "tf-testaccvpnRouteEntrybasic11064"
  cidr_block  = "10.1.0.0/21"
}

resource "alibabacloudstack_vswitch" "default" {
  name               = "${alibabacloudstack_vpc.default.name}"
  vpc_id             = "${alibabacloudstack_vpc.default.id}"
  cidr_block         = "10.1.1.0/24"
  availability_zone  = "${data.alibabacloudstack_zones.default.ids.0}"
}

resource "alibabacloudstack_vpn_gateway" "default" {
  name                 = "${alibabacloudstack_vpc.default.name}"
  vpc_id               = "${alibabacloudstack_vpc.default.id}"
  bandwidth            = 10
  instance_charge_type = "PostPaid"
  enable_ssl           = false
  vswitch_id           = "${alibabacloudstack_vswitch.default.id}"
}

resource "alibabacloudstack_vpn_customer_gateway" "default" {
  name       = "${alibabacloudstack_vpc.default.name}"
  ip_address = "192.168.1.1"
}

resource "alibabacloudstack_vpn_connection" "default" {
  name                = "${alibabacloudstack_vpc.default.name}"
  customer_gateway_id = "${alibabacloudstack_vpn_customer_gateway.default.id}"
  vpn_gateway_id      = "${alibabacloudstack_vpn_gateway.default.id}"
  local_subnet        = ["192.168.2.0/24"]
  remote_subnet       = ["192.168.3.0/24"]
}

resource "alibabacloudstack_vpn_route_entry" "default" {
  weight          = "100"
  publish_vpc     = "false"
  vpn_gateway_id  = "${alibabacloudstack_vpn_gateway.default.id}"
  route_dest      = "10.0.0.0/24"
  next_hop        = "${alibabacloudstack_vpn_connection.default.id}"
}
```

## 参数说明

支持以下参数：

* `vpn_gateway_id` - (必填, 变更时重建) - VPN网关的ID。这是创建路由条目时必须指定的参数，用于标识与该路由条目关联的VPN网关。
* `next_hop` - (必填, 变更时重建) - 目的路由的下一跳。通常情况下，这将是`alibabacloudstack_vpn_connection`资源的ID，表示数据包将通过此连接转发。
* `route_dest` - (必填, 变更时重建) - 目的路由的目标网段。这是一个CIDR格式的字符串，定义了该路由条目适用的网络范围。
* `weight` - (必填) - 目的路由的权重值，取值为 **0** 或 **100**。较高的权重(100)表示优先级更高的路由。
* `publish_vpc` - (必填) - 是否将目的路由发布到VPC。如果设置为`true`，则该路由条目将被传播到关联的VPC路由表中；如果设置为`false`，则不会传播。

## 属性说明

除了上述所有参数外，还导出了以下属性：

* `id` - VPN路由条目的唯一标识符。它由`VpnGatewayId`、`NextHop`和`RouteDest`组合而成，用于唯一标识一个路由条目。
