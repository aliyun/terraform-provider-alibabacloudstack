---
subcategory: "CEN"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cen_cenroutemaps"
sidebar_current: "docs-Alibabacloudstack-datasource-cen-cenroutemaps"
description: |-
  Provides a list of cen cenroutemaps owned by an alibabacloudstack account.
---

# alibabacloudstack\_cen\_cenroutemaps

This data source provides a list of cen cenroutemaps in an alibabacloudstack account according to the specified filters.

## Example Usage
```

variable "name" {
  default = "tf-testAccRouteMapsDatasource123"
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
	description = "tf_testAccCenRouteMap"
}

data "alibabacloudstack_cen_route_maps" "default" {
	cen_id="${alibabacloudstack_cen_instance.default.cen_id}"
	transit_router_route_table_id="${data.alibabacloudstack_cen_transit_router_route_tables.default.transit_router_route_tables.0.id}"
}
```

## Argument Reference

The following arguments are supported:
  * `ids` - (Optional) - the route map id list
  * `transit_router_route_table_id` - (Required) - the route table id
  * `cen_id` - (Required) - the cen id
  * `description_regex` - (Optional) - the route map description regex

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `route_maps` - the route map list
    * `id` - the route map id
    * `route_map_id` - the route map id
    * `cen_id` - the cen id
    * `transit_router_route_table_id` - the route table id
    * `priority` - the priority
    * `map_result` - the policy behavior after all matching conditions are passed.
    * `transmit_direction` - the transmit direction
    * `status` - the status
    * `description` - the description
