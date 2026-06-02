---
subcategory: "Cloud Enterprise Network"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cen_transit_router_route_tables"
description: |-
  Provides a list of cen transitrouterroutetables owned by an alibabacloudstack account.
---

# alibabacloudstack\_cen\_transitrouterroutetables

This data source provides a list of cen transitrouterroutetables in an alibabacloudstack account according to the specified filters.

## Example Usage
```

variable "name" {
	default = "tf-testaccrouter_table70689"
}

resource "alibabacloudstack_cen_instance" "default" {
	cen_instance_name = "${var.name}"
	description = "${var.name}"
	transit_router_name = "${var.name}"
	transit_router_description = "${var.name}"
}


resource "alibabacloudstack_cen_transit_router_route_table" "default" {
  transit_router_route_table_description = "tf-testaccrouter_table70689"
  transit_router_route_table_name = "tf-testaccrouter_table70689"
  transit_router_id = "${alibabacloudstack_cen_instance.default.transit_router_id}"
}

data "alibabacloudstack_cen_transit_router_route_tables" "default" {
	transit_router_route_id="${alibabacloudstack_cen_transit_router_route_table.default.transit_router_id}"

```

## Argument Reference

The following arguments are supported:
  * `ids` - (Optional) - the ids of the cen transit router route tables.
  * `transit_router_route_id` - (Required) - the id of the cen transit router route table.
  * `name_regex` - (Optional) - the name of the cen transit router route table.
  * `description_regex` - (Optional) - the description of the cen transit router route table.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `transit_router_route_tables` - the list of the cen transit router route tables.
    * `id` - the id of the cen transit router route table.
    * `transit_router_route_table_name` - the name of the cen transit router route table.
    * `transit_router_route_table_description` - the description of the cen transit router route table.
    * `create_time` - the create time of the cen transit router route table.
    * `transit_router_route_table_type` - the type of the cen transit router route table.
    * `status` - the status of the cen transit router route table.
