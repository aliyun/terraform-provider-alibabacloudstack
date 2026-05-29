---
subcategory: "高速通道"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_expressconnect_bgp_peer"
sidebar_current: "docs-Alibabacloudstack-expressconnect-bgppeer"
description: |-
  Provides a expressconnect Bgppeer resource.
---

# alibabacloudstack\_expressconnect\_bgppeer

使用Provider配置的凭证在指定的资源集下编排高速通道虚拟边界路由器下的BGP邻居。

## 示例用法
```
variable "name" {
  default = "tf-testaccexpressconnect-bgp-peer1321"
}

resource "alibabacloudstack_express_connect_virtual_border_router" "default" {
	physical_connection_id =     ""
	vlan_id =                    1321
	local_gateway_ip =           "10.0.0.1"
	peer_gateway_ip =            "10.0.0.2"
	peering_subnet_mask =        "255.255.255.252"
	virtual_border_router_name = "${var.name}"
	description =                "TestAccAlibabacloudStackExpressconnectBgpPeer_basic0"
}

resource "alibabacloudstack_expressconnect_bgp_group" "default" {
	bgp_group_name = "${var.name}"
	description =    "${var.name}"
	local_asn =      65534
	peer_asn =       10
	router_id =      "${alibabacloudstack_express_connect_virtual_border_router.default.id}"
	auth_key =       "<YOUR PASSWORD>"
}

resource "alibabacloudstack_expressconnect_bgp_peer" "default" {
  enable_bfd = "true"
  peer_ip_address = "192.168.0.1"
  bfd_multi_hop = "10"
  bgp_group_id = "${alibabacloudstack_expressconnect_bgp_group.default.id}"
  router_id = "${alibabacloudstack_express_connect_virtual_border_router.default.id}"
}
```

## 参数参考

支持以下参数：

  * `bgp_group_id` - (必填) BGP 组的 ID。
  * `router_id` - (选填) 虚拟边界路由器（VBR）的 ID。
  * `peer_ip_address` - (选填) BGP 邻居的 IP 地址。
  * `enable_bfd` - (选填) 是否启用 BFD（双向转发检测）。取值：`true`、`false`。
  * `bfd_multi_hop` - (选填) BFD 多跳的跳数。
  * `bgp_peer_name` - (选填， Available in 1.0.0+) BGP 邻居的名称。
  * `description` - (选填， Available in 1.0.0+) BGP 邻居的描述。
  * `status` - (选填，已废弃) BGP 邻居的状态。该参数不生效。
  * `auth_key` - (选填， Available in 1.0.0+) BGP 邻居的认证密钥。

## 属性参考

除了上述所有参数外，还导出了以下属性：

  * `id` - BGP 邻居的 ID。
  * `bgp_peer_id` - BGP 邻居的 ID。
  * `bgp_status` - BGP 邻居的状态。取值：`Idle`、`Connect`、`Active`、`OpenSent`、`OpenConfirm`、`Established`。
  * `local_asn` - 本地自治系统号。
  * `peer_asn` - BGP 邻居的自治系统号。
  * `ip_version` - IP 版本。取值：`IPv4`、`IPv6`。
  * `is_fake` - 是否启用了次要 ASN。
  * `hold` - BGP Hold 时间。
  * `keepalive` - BGP Keepalive 时间。
  * `route_limit` - BGP 邻居可学习的路由的最大数量。
  * `bgp_peer_name` - BGP 邻居的名称。
  * `description` - BGP 邻居的描述。
  * `auth_key` - BGP 邻居的认证密钥。

## Import

高速通道 BGP 邻居可以使用 BGP 邻居 ID 导入，例如：

```
$ terraform import alibabacloudstack_expressconnect_bgp_peer.example bgppeer-12345678
```
