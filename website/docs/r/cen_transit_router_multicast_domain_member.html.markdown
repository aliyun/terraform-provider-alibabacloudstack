---
subcategory: "Cloud Enterprise Network"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cen_transit_router_multicast_domain_member"
sidebar_current: "docs-Alibabacloudstack-resource-cen-transit-router-multicast-domain-member"
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

* `group_ip_address` - (Required, ForceNew) The multicast IP address. Valid values: 224.0.0.1 to 239.255.255.254. 224.0.0.0 to 224.0.0.127 are reserved addresses and cannot be used.
* `transit_router_multicast_domain_id` - (Required, ForceNew) The ID of the multicast domain to which the multicast member belongs.
* `resource_type` - (Required, ForceNew) The type of the resource. Valid values: `VPC`, `Connect`.
* `vswitch_id` - (Optional, ForceNew, Computed) The ID of the vSwitch. Required when `resource_type` is set to `VPC`.
* `network_interface_id` - (Optional, ForceNew) The ID of the elastic network interface (ENI). Required when `resource_type` is set to `VPC`.
* `connect_peer_id` - (Optional, ForceNew) The ID of the Connect peer. Required when `resource_type` is set to `Connect`.
* `connect_attachment_id` - (Optional, ForceNew, Computed) The ID of the Connect attachment. Required when `resource_type` is set to `Connect`.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource, formatted as `<group_ip_address>:<transit_router_multicast_domain_id>:<resource_type>:<key>`, where `<key>` is `network_interface_id` when `resource_type` is `VPC`, or `connect_peer_id` when `resource_type` is `Connect`.
* `status` - The status of the multicast member.

## Import

CEN transit router multicast domain member can be imported using the id, e.g.

```bash
$ terraform import alibabacloudstack_cen_transit_router_multicast_domain_member.example 224.0.0.1:tr-mcast-domain-1234567890abcdef0:VPC:eni-1234567890abcdef0
```