---
subcategory: "Express Connect"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_expressconnect_bgppeer"
sidebar_current: "docs-Alibabacloudstack-expressconnect-bgppeer"
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
  * `auth_key` - (Optional) - The authentication key of the BGP group.
  * `bfd_multi_hop` - (Optional) - Number of reflexes
  * `bgp_group_id` - (Required) - The ID of the BGP group.
  * `bgp_peer_name` - (Optional) - The name of the BGP neighbor.
  * `description` - (Optional) - Description of the BGP group.
  * `enable_bfd` - (Optional) - Whether the BFD protocol is enabled.
  * `ip_version` - (Optional) - IP version
  * `peer_ip_address` - (Optional) - The IP address of the BGP neighbor.
  * `region_id` - (Optional) - The ID of the region to which the BGP group belongs.
  * `router_id` - (Optional) - The ID of the router.
  * `status` - (Optional) - Status of BGP neighbors

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `auth_key` - The authentication key of the BGP group.
  * `bgp_peer_id` - The ID of the BGP neighbor.
  * `bgp_peer_name` - The name of the BGP neighbor.
  * `bgp_status` - The connection status of BGP, including the following states:* creating: creating.* working: in use.* modifying: Under Modification.* deleting: deleting.* deleted: deleted.
  * `description` - Description of the BGP group.
  * `hold` - Hold time.
  * `ip_version` - IP version
  * `is_fake` - Whether the Fake AS number is enabled.
  * `keepalive` - Live time.
  * `local_asn` - Local ASN number
  * `peer_asn` - ASN of BGP neighbor.
  * `region_id` - The ID of the region to which the BGP group belongs.
  * `route_limit` - Routing restrictions.
  * `status` - Status of BGP neighbors
