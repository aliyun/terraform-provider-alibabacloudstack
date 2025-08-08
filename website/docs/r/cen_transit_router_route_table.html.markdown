---
subcategory: "CEN"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cen_transitrouterroutetable"
sidebar_current: "docs-Alibabacloudstack-cen-transitrouterroutetable"
description: |-
  Provides a cen Transitrouterroutetable resource.
---

# alibabacloudstack\_cen\_transitrouterroutetable

Provides a cen Transitrouterroutetable resource.

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
```

## Argument Reference

The following arguments are supported:
  * `transit_router_id` - (Required, ForceNew) - the transit router id.
  * `transit_router_route_table_name` - (Optional) - the transit router route table name.
  * `transit_router_route_table_description` - (Optional) - the transit router route table description.
  * `tags` - (Optional) - The tag of the resource
    
    * `tag_key` - (Optional) - the tag key.
    
    * `tag_value` - (Optional) - the tag value.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `transit_router_route_table_name` - the transit router route table name.
  * `transit_router_route_table_description` - the transit router route table description.
  * `create_time` - the create time of the transit router route table.
  * `transit_router_route_table_id` - the transit router route table id.
  * `transit_router_route_table_type` - the transit router route table type.
  * `status` - the status of the transit router route table.
