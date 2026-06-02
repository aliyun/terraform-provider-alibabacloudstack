---
subcategory: "Cloud Enterprise Network"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cen_transit_router_route_table_associations"
description: |-
  Provides a list of cen transitrouterroutetableassociations owned by an alibabacloudstack account.
---

# alibabacloudstack\_cen\_transitrouterroutetableassociations

This data source provides a list of cen transitrouterroutetableassociations in an alibabacloudstack account according to the specified filters.

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

data "alibabacloudstack_cen_transit_router_route_table_associations" "default" {
  transit_router_route_table_id = "${alibabacloudstack_cen_transit_router_route_table_association.default.transit_router_route_table_id}"
}
```

## Argument Reference

The following arguments are supported:
  * `ids` - (Optional) - the ids of the cen transit router route table associations.
  * `transit_router_route_table_id` - (Required) - the id of the cen transit router route table.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `transit_router_route_associations` - the list of the cen transit router route table associations.
