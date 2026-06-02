---
subcategory: "Express Connect"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_expressconnect_bgp_peer"
description: |-
  Provides a expressconnect Bgppeer resource.
---

# alibabacloudstack\_expressconnect\_bgppeer

Provides a expressconnect Bgppeer resource.

## Example Usage
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

## Argument Reference

The following arguments are supported:

  * `bgp_group_id` - (Required) The ID of the BGP group.
  * `router_id` - (Optional) The ID of the VBR (Virtual Border Router).
  * `peer_ip_address` - (Optional) The IP address of the BGP peer.
  * `enable_bfd` - (Optional) Whether to enable BFD (Bidirectional Forwarding Detection). Valid values: `true`, `false`.
  * `bfd_multi_hop` - (Optional) The number of BFD multi-hop hops.
  * `bgp_peer_name` - (Optional) The name of the BGP peer.
  * `description` - (Optional) The description of the BGP peer.
  * `status` - (Optional, Deprecated) The status of the BGP peer. This parameter is not effective.
  * `auth_key` - (Optional) The authentication key of the BGP peer.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

  * `id` - The ID of the BGP peer.
  * `bgp_peer_id` - The ID of the BGP peer.
  * `bgp_status` - The status of the BGP peer. Valid values: `Idle`, `Connect`, `Active`, `OpenSent`, `OpenConfirm`, `Established`.
  * `local_asn` - The local autonomous system number.
  * `peer_asn` - The autonomous system number of the BGP peer.
  * `ip_version` - The IP version. Valid values: `IPv4`, `IPv6`.
  * `is_fake` - Whether the secondary ASN is enabled.
  * `hold` - The BGP Hold time.
  * `keepalive` - The BGP Keepalive time.
  * `route_limit` - The maximum number of routes that the BGP peer can learn.
  * `bgp_peer_name` - The name of the BGP peer.
  * `description` - The description of the BGP peer.
  * `auth_key` - The authentication key of the BGP peer.

## Import

Express Connect BGP Peer can be imported using the BGP peer ID, e.g.

```
$ terraform import alibabacloudstack_expressconnect_bgp_peer.example bgppeer-12345678
```
