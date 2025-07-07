---
subcategory: "ExpressConnect"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_expressconnect_bgppeer"
sidebar_current: "docs-Alibabacloudstack-expressconnect-bgppeer"
description: |-
  Provides a expressconnect Bgppeer resource.
---

# alibabacloudstack\_expressconnect\_bgppeer

Provides a expressconnect Bgppeer resource.

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
	auth_key =       "YoPcOh&7Tt"
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
  * `auth_key` - (选填) - BGP组的认证密钥。
  * `bfd_multi_hop` - (选填) - 反射次数
  * `bgp_group_id` - (必填) - BGP组的ID。
  * `bgp_peer_name` - (选填) - BGP邻居的名称。
  * `description` - (选填) - BGP组的描述。 
  * `enable_bfd` - (选填) - 是否开启了BFD协议。
  * `ip_version` - (选填) - IP版本
  * `peer_ip_address` - (选填) - BGP邻居的IP地址。
  * `region_id` - (选填) - BGP组所属的地域ID。
  * `router_id` - (选填) - 路由器的ID。
  * `status` - (选填) - BGP邻居的状态

## 属性参考

除了上述所有参数外，还导出了以下属性：
  * `auth_key` - BGP组的认证密钥。
  * `bgp_peer_id` - BGP邻居的ID。
  * `bgp_peer_name` - BGP邻居的名称。
  * `bgp_status` - BGP的连接状态，包含以下状态：* creating：创建中。* working：使用中。* modifying：修改中。* deleting：删除中。* deleted：已删除。
  * `description` - BGP组的描述。 
  * `hold` - 保持时间。
  * `ip_version` - IP版本
  * `is_fake` - 是否启用了Fake AS号。
  * `keepalive` - 保活时间。
  * `local_asn` - 本地ASN号
  * `peer_asn` - BGP邻居的ASN。
  * `region_id` - BGP组所属的地域ID。
  * `route_limit` - 路由限制。
  * `status` - BGP邻居的状态
