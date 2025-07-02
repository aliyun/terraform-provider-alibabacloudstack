---
subcategory: "VPNGateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpngateway_vpnpbrrouteentry"
sidebar_current: "docs-Alibabacloudstack-vpngateway-vpnpbrrouteentry"
description: |-
  Provides a vpngateway Vpnpbrrouteentry resource.
---

# alibabacloudstack\_vpngateway\_vpnpbrrouteentry

Provides a vpngateway Vpnpbrrouteentry resource.

## 示例用法
```
data "alibabacloudstack_zones" "default"{
}

resource "alibabacloudstack_vpc" "default" {
 name  = "tf-testaccVpngatewayVpnpbrrouteentrybasic11558"
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
  route_dest = "10.0.0.0/24"
  route_source = "192.168.0.0/24"
  next_hop = "${alibabacloudstack_vpn_connection.default.id}"
  weight = "100"
  publish_vpc = "false"
}
```

## 参数参考

支持以下参数：
  * `next_hop` - (必填, 变更时重建) - 目的路由的下一跳。
  * `overlay_mode` - (选填, 变更时重建) - 隧道协议，取值：**Ipsec**（IPsec隧道协议）。
  * `publish_vpc` - (必填) - 是否发布策略路由到VPC，取值：- **true**：发布策略路由到VPC。- **false**：不发布策略路由到VPC。
  * `route_dest` - (必填, 变更时重建) - 目的路由的目标网段。
  * `route_source` - (必填, 变更时重建) - 策略路由的源网段。
  * `vpn_gateway_id` - (必填, 变更时重建) - VPN网关的ID。
  * `weight` - (必填) - 目的路由的权重值，取值：**0**|**100**。

## 属性参考

除了上述所有参数外，还导出了以下属性：
  * `create_time` - VPN目的路由的创建时间。
  * `status` - VPN目的路由的状态。
