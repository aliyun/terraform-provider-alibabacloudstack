---
subcategory: "VPN网关"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpngateway_vpn_pbr_route_entries"
sidebar_current: "docs-Alibabacloudstack-datasource-vpngateway-vpnpbrrouteentries"
description: |-
  提供阿里云账号下拥有的vpngateway vpnpbrrouteentries列表。
---

# alibabacloudstack\_vpngateway\_vpnpbrrouteentries

此数据源提供根据指定过滤条件列出的阿里云账号下的vpngateway vpnpbrrouteentries资源列表。

## 示例用法
```
data "alibabacloudstack_zones" "default"{
}

resource "alibabacloudstack_vpc" "default" {
 name  = "tf-testAccVpngatewayVpnPbrRouteEntriesDataSource-6376584"
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
	route_dest =     "10.0.0.0/24"
	route_source =   "192.168.0.0/24"
	next_hop =      "${alibabacloudstack_vpn_connection.default.id}"
	weight =       "100"
	publish_vpc =    "false"
}
 

data "alibabacloudstack_vpngateway_vpn_pbr_route_entrys" "default" {
  vpn_gateway_id = [
                     "${alibabacloudstack_vpngateway_vpn_pbr_route_entry.default.vpn_gateway_id}"
                   ]
}
```

## 参数参考
以下参数是支持的：
  * `ids` - (选填) - VPN目的路由的ID列表.
  * `vpn_gateway_id` - (必填) - VPN网关的ID。

## Attributes Reference
除了上述参数外，还导出以下属性：
  * `vpn_pbr_route_entries` - VPN目的路由策略列表。
    * `id` - VPN目的路由策略的ID。
    * `create_time` - VPN目的路由的创建时间。
    * `description` - 策略路由的描述信息。
    * `next_hop` - 目的路由策略的下一跳。
    * `route_dest` - 目的路由v的目标网段。
    * `route_source` - 策略路由的源网段。
    * `status` - VPN目的路由策略的状态。
    * `vpn_gateway_id` - VPN网关的ID。
    * `weight` - 目的路由的权重值，取值：**0**|**100**。
