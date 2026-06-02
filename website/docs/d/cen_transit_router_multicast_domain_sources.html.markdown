---
subcategory: "Cloud Enterprise Network"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cen_transit_router_multicast_domain_sources"
description: |-
  Provides a list of cen transitroutermulticastdomainsources owned by an alibabacloudstack account.
---

# alibabacloudstack\_cen\_transitroutermulticastdomainsources

This data source provides a list of cen transit router multicast domain sources in an alibabacloudstack account according to the specified filters.

## Example Usage

```hcl
variable "name" {
  default = "tf-testaccmulticast_domain_source"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details              = true
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name   = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  name       = "${var.name}_vsw"
  vpc_id     = alibabacloudstack_vpc_vpc.default.id
  cidr_block = "172.16.1.0/24"
  zone_id    = data.alibabacloudstack_zones.default.zones.0.id
}

resource "alibabacloudstack_cen_instance" "default" {
  cen_instance_name       = var.name
  description             = var.name
  transit_router_name     = var.name
  transit_router_id = var.name
}

resource "alibabacloudstack_cen_transit_router_vpc_attachment" "default" {
  transit_router_attachment_name        = var.name
  transit_router_attachment_description = var.name
  cen_id                                = alibabacloudstack_cen_instance.default.id
  vpc_id                                = alibabacloudstack_vpc_vswitch.default.vpc_id
  transit_router_id                     = alibabacloudstack_cen_instance.default.transit_router_id
  zone_mappings {
    vswitch_id = alibabacloudstack_vpc_vswitch.default.id
    zone_id    = alibabacloudstack_vpc_vswitch.default.zone_id
  }
}

resource "alibabacloudstack_cen_transit_router_multicast_domain" "default" {
  transit_router_multicast_domain_description = var.name
  transit_router_multicast_domain_name        = var.name
  transit_router_id                           = alibabacloudstack_cen_instance.default.transit_router_id
}

resource "alibabacloudstack_cen_transit_router_multicast_domain_association" "default" {
  transit_router_attachment_id       = alibabacloudstack_cen_transit_router_vpc_attachment.default.transit_router_attachment_id
  transit_router_multicast_domain_id = alibabacloudstack_cen_transit_router_multicast_domain.default.transit_router_multicast_domain_id
  vswitch_id                         = alibabacloudstack_vpc_vswitch.default.id
}

data "alibabacloudstack_cen_transit_router_multicast_domain_sources" "default" {
  transit_router_multicast_domain_id = alibabacloudstack_cen_transit_router_multicast_domain.default.transit_router_multicast_domain_id
}
```

## Argument Reference

The following arguments are supported:

* `transit_router_multicast_domain_id` - (Required) The ID of the transit router multicast domain.
* `ids` - (Optional, Computed) A list of IDs to filter results.
* `transit_router_attachment_id` - (Optional) The ID of the transit router attachment.
* `vswitch_id` - (Optional) The ID of the VSwitch.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `transit_router_multicast_groups` - A list of transit router multicast groups. Each element contains the following attributes:
  * `id` - The ID of the multicast group source, formatted as `<group_ip_address>:<transit_router_multicast_domain_id>:<resource_type>:<resource_data>`.
  * `group_ip_address` - The multicast IP address.
  * `network_interface_id` - The ID of the network interface.
  * `transit_router_multicast_domain_id` - The ID of the transit router multicast domain.
  * `status` - The status of the multicast group source.
  * `transit_router_attachment_id` - The ID of the transit router attachment.
  * `vswitch_id` - The ID of the VSwitch.
  * `resource_type` - The type of the resource. Valid values: `VPC`, `Connect`.
  * `resource_id` - The ID of the resource.
  * `source_type` - The source type.
  * `group_source` - Whether it is a group source.
  * `group_member` - Whether it is a group member.
* `ids` - A list of transit router multicast domain IDs.
