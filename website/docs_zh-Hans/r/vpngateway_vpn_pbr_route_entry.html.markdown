---
subcategory: "VPN网关"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpngateway_vpn_pbr_route_entry"
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
  * `vpn_gateway_id` - (必填, 变更时重建) VPN网关的ID。
  * `route_source` - (必填, 变更时重建) 策略路由的源网段。
  * `route_dest` - (必填, 变更时重建) 目的路由的目标网段。
  * `next_hop` - (必填, 变更时重建) 目的路由的下一跳。
  * `weight` - (必填) 目的路由的权重值，取值：`0` 或 `100`。
  * `publish_vpc` - (必填) 是否发布策略路由到VPC，取值：`true`（发布）或 `false`（不发布）。
  * `overlay_mode` - (选填) 隧道协议，默认值：`Ipsec`，取值：`Ipsec`（IPsec隧道协议）。

## 属性参考

除了上述所有参数外，还导出了以下属性：
  * `id` - 资源的ID，格式为 `<vpn_gateway_id>_<route_source>_<route_dest>_<next_hop>`。
  * `create_time` - VPN PBR路由条目的创建时间。
  * `status` - VPN PBR路由条目的状态。

## Import

VPN PBR路由条目可以使用组合ID进行导入，例如：

```
$ terraform import alibabacloudstack_vpngateway_vpn_pbr_route_entry.example <vpn_gateway_id>_<route_source>_<route_dest>_<next_hop>
```
