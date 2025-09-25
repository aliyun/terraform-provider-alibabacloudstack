---
subcategory: "CEN"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cen_transit_router_connect_peer"
sidebar_current: "docs-alibabacloudstack-resource-cen-transit-router-connect-peer"
description: |-
  Provides a Alibaba Cloud CEN transit router connect peer resource.
---

# alibabacloudstack_cen_transit_router_connect_peer

Provides a CEN transit router connect peer resource.

For information about CEN transit router connect peer and how to use it, see [What is Transit Router Connect Peer](https://www.alibabacloud.com/help/doc-detail/65872.html).

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
```

## Argument Reference

The following arguments are supported:

* `cen_id` - (Required, ForceNew) The ID of the CEN instance.
* `connect_attachment_id` - (Required, ForceNew) The ID of the transit router connect attachment.
* `local_ip` - (Optional, ForceNew) The local IP address of the transit router connect peer.Must be within the Transit Router VPC network segment of the instance.
* `name` - (Optional, ForceNew) The name of the transit router connect peer.
* `peer_ip` - (Required, ForceNew) The peer IP address of the transit router connect peer.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource. The value formats as `<cen_id>:<connect_attachment_id>:<peer_id>`.
* `peer_id` - The ID of the transit router connect peer.
* `region_id` - The region ID of the transit router connect peer.

## Import

CEN transit router connect peer can be imported using the id, e.g.

```bash
$ terraform import alibabacloudstack_cen_transit_router_connect_peer.example cen-2cudwl7et3716ozh****:tr-attach-2cudwl7et3716ozh****:tr-cp-2cudwl7et3716ozh****
```