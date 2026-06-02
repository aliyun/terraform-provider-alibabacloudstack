---
subcategory: "Cloud Enterprise Network"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cen_transit_router_route_table_association"
sidebar_current: "docs-Alibabacloudstack-resource-cen-transit-router-route-table-association"
description: |-
  Provides a cen Transitrouterroutetableassociation resource.
---

# alibabacloudstack\_cen\_transitrouterroutetableassociation

Provides a cen Transitrouterroutetableassociation resource.

## Example Usage
```
variable "name" {
	default = "tf-testaccrouter_table_association32777"
}


data "alibabacloudstack_zones" default {
  available_resource_creation = "VSwitch"
  enable_details = true
}


resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  name = "${var.name}_vsw"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block = "172.16.1.0/24"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
}



resource "alibabacloudstack_cen_instance" "default" {
	cen_instance_name = "${var.name}"
	description = "${var.name}"
	transit_router_name = "${var.name}"
	transit_router_description = "${var.name}"
}

resource "alibabacloudstack_cen_transit_router_vpc_attachment" "default" {
	transit_router_attachment_name = "${var.name}"
	transit_router_attachment_description = "${var.name}"
	cen_id = "${alibabacloudstack_cen_instance.default.id}"
	vpc_id = "${alibabacloudstack_vpc_vswitch.default.vpc_id}"
	transit_router_id = "${alibabacloudstack_cen_instance.default.transit_router_id}"
	zone_mappings {
			vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
			 zone_id = "${alibabacloudstack_vpc_vswitch.default.zone_id}"
		}
}

data "alibabacloudstack_cen_transit_router_route_tables" "default" {
	transit_router_route_id="${alibabacloudstack_cen_instance.default.transit_router_id}"
}



resource "alibabacloudstack_cen_transit_router_route_table_association" "default" {
  transit_router_route_table_id = "${data.alibabacloudstack_cen_transit_router_route_tables.default.transit_router_route_tables.0.id}"
  transit_router_attachment_id = "${alibabacloudstack_cen_transit_router_vpc_attachment.default.transit_router_attachment_id}"
}
```

## Argument Reference

The following arguments are supported:
  * `transit_router_attachment_id` - (Required) - the transit router attachment id.
  * `transit_router_route_table_id` - (Required, ForceNew) - the transit router route table id.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `resource_id` - the resource id.
  * `resource_type` - the resource type.
  * `status` - the transit route table association status.
