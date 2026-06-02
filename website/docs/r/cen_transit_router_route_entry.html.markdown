---
subcategory: "Cloud Enterprise Network"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cen_transit_router_route_entry"
description: |-
  Provides a cen Transitrouterrouteentry resource.
---

# alibabacloudstack\_cen\_transitrouterrouteentry

Provides a cen Transitrouterrouteentry resource.

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
```

## Argument Reference

The following arguments are supported:
  * `transit_router_route_entry_description` - (Optional) - the description of the transit router route entry.
  * `transit_router_route_entry_destination_cidr_block` - (Required, ForceNew) - the destination CIDR block of the transit router route entry.
  * `transit_router_route_entry_name` - (Optional) - the name of the transit router route entry.
  * `transit_router_route_entry_next_hop_id` - (Optional) - the next hop ID of the transit router route entry.
  * `transit_router_route_entry_next_hop_type` - (Required, ForceNew) - the next hop type of the transit router route entry.
  * `transit_router_route_table_id` - (Required, ForceNew) - the transit router route table ID of the transit router route entry.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `create_time` - the creation time of the transit router route entry.
  * `status` - the status of the transit router route entry.
  * `transit_router_route_entry_id` - the transit router route entry ID of the transit router route entry.
  * `transit_router_route_entry_type` -the type of the transit router route entry.
  * `operational_mode` - the operational mode of the transit router route entry.
