---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cen_transit_router_multicast_domain_member"
sidebar_current: "docs-Alibabacloudstack-datasource-cen-transit-router-multicast-domain-member"
description: |-
  Provides a AlibabacloudStack CEN transit router multicast domain member resource.
---

# alibabacloudstack_cen_transit_router_multicast_domain_member

Provides a CEN transit router multicast domain member resource.

For information about CEN transit router multicast domain member and how to use it, see [What is Transit Router Multicast Domain Member](https://www.alibabacloud.com/help/doc-detail/).

## Example Usage

Basic Usage

```hcl
variable "name" {
  default = "tf-testAccRouteTable"
}

data "alibabacloudstack_zones" default {
  available_resource_creation = "VSwitch"
  enable_details = true
}


resource "alibabacloudstack_vpc" "example" {
  vpc_name       = var.name
  cidr_block     = "10.0.0.0/8"
}

resource "alibabacloudstack_vswitch" "example" {
  vpc_id       = alibabacloudstack_vpc.example.id
  cidr_block   = "10.1.0.0/16"
  zone_id      = "${data.alibabacloudstack_zones.default.zones.0.id}"
  vswitch_name = var.name
}

resource "alibabacloudstack_cen_instance" "example" {
  cen_instance_name = var.name
  description       = var.name
}

resource "alibabacloudstack_cen_transit_router" "example" {
  cen_id = alibabacloudstack_cen_instance.example.id
}

resource "alibabacloudstack_cen_transit_router_multicast_domain" "example" {
  cen_id                                = alibabacloudstack_cen_instance.example.id
  transit_router_id                     = alibabacloudstack_cen_transit_router.example.transit_router_id
  transit_router_multicast_domain_name  = var.name
}

resource "alibabacloudstack_network_interface" "example" {
  vswitch_id = alibabacloudstack_vswitch.example.id
}

resource "alibabacloudstack_cen_transit_router_multicast_domain_member" "example" {
  group_ip_address                      = "224.0.0.1"
  network_interface_id                  = alibabacloudstack_network_interface.example.id
  transit_router_multicast_domain_id    = alibabacloudstack_cen_transit_router_multicast_domain.example.id
  vswitch_id                            = alibabacloudstack_vswitch.example.id
  resource_type                         = "VPC"
}
```

## Argument Reference

The following arguments are supported:

* `group_ip_address` - (Required, ForceNew) The multicast IP address. 
* `transit_router_multicast_domain_id` - (Required, ForceNew) Forwarding router multicast domain ID.
* `resource_type` - (Required, ForceNew) Resource type. Valid values: `VPC`, `Connect`.
* `vswitch_id` - (Optional, ForceNew) The ID of the switch to which the multicast member belongs. Required when `resource_type` is `VPC`.
* `network_interface_id` - (Optional, ForceNew) The ID of the network interface. Required when `resource_type` is `VPC`.
* `connect_peer_id` - (Optional, ForceNew) Connection peer ID. Required when `resource_type` is `Connect`.
* `connect_attachment_id` - (Optional, ForceNew, Computed) Connection attachment ID. Required when `resource_type` is `Connect`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource, formatted as `<group_ip_address>:<transit_router_multicast_domain_id>:<resource_type>:<key>`.
* `status` - The status of the multicast member.
* `vswitch_id` - The ID of the switch to which the multicast member belongs.
* `network_interface_id` - The ID of the network interface.
* `connect_peer_id` - Connection peer ID.
* `connect_attachment_id` - Connection attachment ID.

## Import

CEN transit router multicast domain member can be imported using the id, e.g.

```bash
$ terraform import alibabacloudstack_cen_transit_router_multicast_domain_member.default 224.0.0.1:tr-mcast-domain-1234567890abcdef0:VPC:eni-1234567890abcdef0
```