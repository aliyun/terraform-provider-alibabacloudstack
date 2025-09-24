---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cen_transit_router_multicast_domain_members"
sidebar_current: "docs-alibabacloudstack-datasource-cen-transit-router-multicast-domain-members"
description: |-
  Provides a list of CEN transit router multicast domain members to be used by an Alibaba Cloud Stack account.
---

# alibabacloudstack\_cen\_transit\_router\_multicast\_domain\_members

This data source provides a list of CEN transit router multicast domain members in an Alibaba Cloud Stack account according to the specified filters.

## Example Usage

```hcl
variable "name" {
  default = "tf-testAccMulticastDomainMember"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
}

resource "alibabacloudstack_vpc" "default" {
  vpc_name = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vswitch" "default" {
  name = "${var.name}_vsw"
  vpc_id = "${alibabacloudstack_vpc.default.id}"
  cidr_block = "172.16.1.0/24"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

resource "alibabacloudstack_cen_instance" "default" {
  cen_instance_name = "${var.name}"
  description = "${var.name}"
}

resource "alibabacloudstack_cen_transit_router" "default" {
  cen_id = "${alibabacloudstack_cen_instance.default.id}"
}

resource "alibabacloudstack_cen_transit_router_multicast_domain" "default" {
  transit_router_multicast_domain_name = "${var.name}"
  transit_router_id = "${alibabacloudstack_cen_transit_router.default.transit_router_id}"
  cen_id = "${alibabacloudstack_cen_instance.default.id}"
}

resource "alibabacloudstack_network_interface" "default" {
  vswitch_id = "${alibabacloudstack_vswitch.default.id}"
}

resource "alibabacloudstack_cen_transit_router_multicast_domain_member" "default" {
  group_ip_address = "224.0.0.1"
  network_interface_id = "${alibabacloudstack_network_interface.default.id}"
  transit_router_multicast_domain_id = "${alibabacloudstack_cen_transit_router_multicast_domain.default.id}"
  vswitch_id = "${alibabacloudstack_vswitch.default.id}"
}

data "alibabacloudstack_cen_transit_router_multicast_domain_members" "default" {
  transit_router_multicast_domain_id = "${alibabacloudstack_cen_transit_router_multicast_domain.default.id}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of multicast domain member IDs.
* `transit_router_multicast_domain_id` - (Required) The ID of the multicast domain to which the multicast member belongs.
* `transit_router_attachment_id` - (Optional) The ID of the transit router attachment.
* `vswitch_id` - (Optional) The ID of the VSwitch.

## Attributes Reference

The following attributes are exported:

* `transit_router_multicast_groups` - A list of multicast domain members. Each element contains the following attributes:
  * `id` - The ID of the multicast member.
  * `group_ip_address` - The multicast IP address.
  * `network_interface_id` - The ID of the network interface.
  * `status` - The status of the multicast member.
  * `transit_router_multicast_domain_id` - The ID of the multicast domain to which the multicast member belongs.
  * `transit_router_attachment_id` - The ID of the transit router attachment.
  * `vswitch_id` - The ID of the VSwitch.
  * `resource_type` - The type of the resource.
  * `member_type` - The type of the member.
  * `resource_id` - The ID of the resource.
  * `group_source` - Whether the group is a source group.
  * `group_member` - Whether the group is a member group.
