---
subcategory: "Express Connect"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_expressconnect_bgp_peers"
sidebar_current: "docs-Alibabacloudstack-datasource-expressconnect-bgppeers"
description: |-
  Provides a list of expressconnect bgppeers owned by an alibabacloudstack account.
---

# alibabacloudstack\_expressconnect\_bgppeers

This data source provides a list of expressconnect bgppeers in an alibabacloudstack account according to the specified filters.

## Example Usage
```
resource "alibabacloudstack_expressconnect_bgp_peer" "default" {
  bgp_group_id =   "${alibabacloudstack_expressconnect_bgp_group.default.id}"
  router_id =      "${alibabacloudstack_express_connect_virtual_border_router.default.id}"
  enable_bfd =      "true"
  peer_ip_address = "192.168.0.1"
  bfd_multi_hop =   "10"
}

data "alibabacloudstack_expressconnect_bgp_peers" "default" {
  
}
```

## Argument Reference

The following arguments are supported:
  * `ids` - (Optional) - A list of BGP peer IDs.
  * `router_id` - (Optional) - The ID of the router.
  * `region_id` - (Optional) - The ID of the region to which the BGP group belongs.
  * `bgp_peer_id` - (Optional) - The ID of the BGP neighbor.
  * `bgp_group_id` - (Optional) - The ID of the BGP group.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `bgp_peers` - A list of BGP peers.
    * `id` - The ID of the BGP peer.
    * `auth_key` - The authentication key of the BGP group.
    * `bfd_multi_hop` - Number of reflexes
    * `bgp_group_id` - The ID of the BGP group.
    * `bgp_peer_id` - The ID of the BGP neighbor.
    * `bgp_peer_name` - The name of the BGP neighbor.
    * `bgp_status` - The connection status of BGP, including the following states:* creating: creating.* working: in use.* modifying: Under Modification.* deleting: deleting.* deleted: deleted.
    * `description` - Description of the BGP group.
    * `enable_bfd` - Whether the BFD protocol is enabled.
    * `hold` - Hold time.
    * `ip_version` - IP version
    * `is_fake` - Whether the Fake AS number is enabled.
    * `keepalive` - Live time.
    * `local_asn` - Local ASN number
    * `peer_asn` - ASN of BGP neighbor.
    * `peer_ip_address` - The IP address of the BGP neighbor.
    * `region_id` - The ID of the region to which the BGP group belongs.
    * `route_limit` - Routing restrictions.
    * `router_id` - The ID of the router.
    * `status` - Status of BGP neighbors
