---
subcategory: "CEN"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cen_transit_router_connect_peers"
sidebar_current: "docs-alibabacloudstack-datasource-cen-transit-router-connect-peers"
description: |-
  Provides a list of CEN Transit Router Connect Peers to be used by the alibabacloudstack_cen_transit_router_connect_peer resource.
---

# alibabacloudstack_cen_transit_router_connect_peers

This data source provides a list of CEN Transit Router Connect Peers in an Alibaba Cloud account according to the specified filters.

## Example Usage

```hcl
variable "name" {
  default = "tf-testAccTransitRouterConnectPeer48958"
}

resource "alibabacloudstack_cen_instance" "default" {
  description = "tf-testaccceninstance48958"
  cen_instance_name = "tf-testaccceninstance48958"
  transit_router_cidrs {
    cidr = "172.16.0.0/16"
  }
}

resource "alibabacloudstack_express_connect_virtual_border_router" "default" {
  physical_connection_id = "pc-xxxxxxxxxxxxxxxxxxx"
  vlan_id =                    99
  local_gateway_ip =           "10.0.0.1"
  peer_gateway_ip =            "10.0.0.2"
  peering_subnet_mask =        "255.255.255.252"
  virtual_border_router_name = "${var.name}"
}

resource "alibabacloudstack_cen_transit_router_vbr_attachment" "default" {
  vbr_id = "${alibabacloudstack_express_connect_virtual_border_router.default.id}"
  cen_id = "${alibabacloudstack_cen_instance.default.id}"
  transit_router_id = "${alibabacloudstack_cen_instance.default.transit_router_id}"
}

resource "alibabacloudstack_cen_transit_router_connect_attachment" "default" {
  cen_id = "${alibabacloudstack_cen_transit_router_vbr_attachment.default.cen_id}"
  transit_router_id = "${alibabacloudstack_cen_instance.default.transit_router_id}"
  transit_router_attachment_name = "${var.name}"
  depends_on = ["alibabacloudstack_cen_transit_router_vbr_attachment.default"]
}

resource "alibabacloudstack_cen_transit_router_connect_peer" "default" {
  name                  = "${var.name}"
  local_ip              = "172.16.0.1"
  peer_ip               = "172.16.0.11"
  cen_id                = "${alibabacloudstack_cen_transit_router_connect_attachment.default.cen_id}"
  connect_attachment_id = "${alibabacloudstack_cen_transit_router_connect_attachment.default.transit_router_attachment_id}"
}

data "alibabacloudstack_cen_transit_router_connect_peers" "example" {
  cen_id = "${alibabacloudstack_cen_transit_router_connect_peer.default.cen_id}"
  connect_attachment_id = "${alibabacloudstack_cen_transit_router_connect_peer.default.connect_attachment_id}"
}

output "first_connect_peer_id" {
  value = data.alibabacloudstack_cen_transit_router_connect_peers.example.peers.0.id
}
```

## Argument Reference

The following arguments are supported:

* `cen_id` - (Required, ForceNew) The ID of the CEN instance.
* `connect_attachment_id` - (Required, ForceNew) The ID of the Transit Router Connect attachment.
* `ids` - (Optional) A list of Transit Router Connect Peer IDs to filter results by.
* `name_regex` - (Optional) A regex string to filter results by the Transit Router Connect Peer name.
* `peer_name` - (Optional, ForceNew) The name of the Transit Router Connect Peer.

## Attributes Reference

The following attributes are exported:

* `peers` - A list of Transit Router Connect Peers. Each element contains the following attributes:
  * `id` - The ID of the Transit Router Connect Peer. The format is `{cen_id}:{connect_attachment_id}:{peer_id}`.
  * `cen_id` - The ID of the CEN instance.
  * `connect_attachment_id` - The ID of the Transit Router Connect attachment.
  * `peer_id` - The ID of the Transit Router Connect Peer.
  * `name` - The name of the Transit Router Connect Peer.
  * `local_ip` - The local IP address of the Transit Router Connect Peer.
  * `peer_ip` - The peer IP address of the Transit Router Connect Peer.
  * `region_id` - The region ID of the Transit Router Connect Peer.
  * `status` - The status of the Transit Router Connect Peer.
  * `creation_time` - The creation time of the Transit Router Connect Peer.