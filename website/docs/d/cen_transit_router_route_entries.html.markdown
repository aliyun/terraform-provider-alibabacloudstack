---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cen_transit_router_route_entries"
sidebar_current: "docs-Alibabacloudstack-datasource-cen-transit-router-route-entries"
description: |-
  Provides a list of cen cen_transit_router_route_entries owned by an alibabacloudstack account.
---

# alibabacloudstack\_cen\cen_transit_router_route_entries

This data source provides a list of cen transitrouterrouteentries in an alibabacloudstack account according to the specified filters.

## Example Usage
```

variable "name" {
	default = "tf-testaccrouter_entry58946"
}

resource "alibabacloudstack_cen_instance" "default" {
	cen_instance_name = "${var.name}"
	description = "${var.name}"
	transit_router_name = "${var.name}"
	transit_router_description = "${var.name}"
}
resource "alibabacloudstack_cen_transit_router_route_table" "default" {
	transit_router_route_table_description = "${var.name}"
	transit_router_route_table_name = "${var.name}"
	transit_router_id = "${alibabacloudstack_cen_instance.default.transit_router_id}"
}


resource "alibabacloudstack_cen_transit_router_route_entry" "default" {
  transit_router_route_entry_next_hop_type = "BlackHole"
  transit_router_route_table_id = "${alibabacloudstack_cen_transit_router_route_table.default.transit_router_route_table_id}"
  transit_router_route_entry_description = "tf-testaccrouter_entry58946"
  transit_router_route_entry_destination_cidr_block = "10.10.10.1/32"
  transit_router_route_entry_name = "tf-testaccrouter_entry58946"
}

data "alibabacloudstack_cen_transit_router_route_entries" "default" {
transit_router_route_table_id="${alibabacloudstack_cen_transit_router_route_table.default.transit_router_route_table_id}"
name_regex = "${alibabacloudstack_cen_transit_router_route_entry.default.transit_router_route_entry_name}"
```

## Argument Reference

The following arguments are supported:
  * `ids` - (Optional) - the ids of the cen transitrouterrouteentries.
  * `transit_router_route_table_id` - (Required) - the router table id
  * `name_regex` - (Optional) - the name of the cen transitrouterrouteentries.
  * `description_regex` - (Optional) - the description of the cen transitrouterrouteentries.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `transit_router_route_entries` - the list of cen transitrouterrouteentries.
    * `id` - the id of the cen transitrouterrouteentries.
    * `create_time` - the create time of the cen transitrouterrouteentries.
    * `status` - the status of the cen transitrouterrouteentries.
    * `transit_router_route_entry_description` - the description of the cen transitrouterrouteentries.
    * `transit_router_route_entry_destination_cidr_block` - the destination cidr block of the cen transitrouterrouteentries.
    * `transit_router_route_entry_id` - the id of the cen transitrouterrouteentries.
    * `transit_router_route_entry_name` - the name of the cen transitrouterrouteentries.
    * `transit_router_route_entry_next_hop_id` - the next hop id of the cen transitrouterrouteentries.
    * `transit_router_route_entry_next_hop_type` - the next hop type of the cen transitrouterrouteentries.
    * `transit_router_route_entry_type` - the type of the cen transitrouterrouteentries.
    * `operational_mode` - the operational mode of the cen transitrouterrouteentries.
    * `transit_router_route_entry_status` - the status of the cen transitrouterrouteentries.
