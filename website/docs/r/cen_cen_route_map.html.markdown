---
subcategory: "CEN"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cen_cenroutemap"
sidebar_current: "docs-Alibabacloudstack-cen-cenroutemap"
description: |-
  Provides a cen Cenroutemap resource.
---

# alibabacloudstack\_cen\_cenroutemap

Provides a cen Cenroutemap resource.

## Example Usage
```
variable "name" {
	default = "tf-testaccroute_map95277"
}

resource "alibabacloudstack_cen_instance" "default" {
	cen_instance_name = "${var.name}"
	description = "${var.name}"
	transit_router_name = "${var.name}"
	transit_router_description = "${var.name}"
}

data "alibabacloudstack_cen_transit_router_route_tables" "default" {
	transit_router_route_id="${alibabacloudstack_cen_instance.default.transit_router_id}"
}



resource "alibabacloudstack_cen_route_map" "default" {
  cen_id = "${alibabacloudstack_cen_instance.default.cen_id}"
  transit_router_route_table_id = "${data.alibabacloudstack_cen_transit_router_route_tables.default.transit_router_route_tables.0.id}"
  priority = "3"
  transmit_direction = "RegionIn"
  map_result = "Deny"
}
```

## Argument Reference

The following arguments are supported:
  * `cen_id` - (Required, ForceNew) - cen instance id
  * `transit_router_route_table_id` - (Required, ForceNew) - cen route table id
  * `priority` - (Required) - priority
  * `transmit_direction` - (Required) - transmit direction("RegionIn", "RegionOut")
  * `map_result` - (Required) - The policy behavior after all matching conditions are passed. Support the following behaviors:("Permit", "Deny")
Allow: Allow routing through the matched route. Allow to modify routing properties.
Reject: Refuse to pass through the matched route.
  * `as_path_match_mode` - (Optional) - as path match mode("Include", "Complete")
  * `cidr_match_mode` - (Optional) - cidr match mode(Include", "Complete")
  * `community_match_mode` - (Optional) - community match mode("Include", "Complete")
  * `community_operate_mode` - (Optional) - community operate mode("Additive", "Replace")
  * `description` - (Optional) - description
  * `destination_instance_ids_reverse_match` - (Optional) - destination instance ids reverse match
  * `match_address_type` - (Optional) - match address type("IPv6", "IPv4")
  * `next_priority` - (Optional) - Priority of the next routing policy associated with it
  * `preference` - (Optional) - Set routing priority. The value range is 1-100, and the default priority for routing is 50. The lower the value, the higher the priority.
  * `source_instance_ids_reverse_match` - (Optional) - source instance ids reverse match

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `description` - description of the route map.
  * `route_map_id` - the id of route map.
  * `status` - the status of route map.
