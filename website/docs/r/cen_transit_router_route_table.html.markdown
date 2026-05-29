---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cen_transit_router_route_table"
sidebar_current: "docs-Alibabacloudstack-resource-cen-transit-router-route-table"
description: |-
  Provides a CEN Transit Router Route Table resource.
---

# alibabacloudstack\_cen\_transit\_router\_route\_table

Provides a CEN Transit Router Route Table resource.

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
  * `transit_router_id` - (Required, ForceNew) - The ID of the transit router. Modifying this parameter will force a new resource to be created.
  * `transit_router_route_table_name` - (Optional, Readable) - The name of the route table. This attribute is returned by the API and can be manually set.
  * `transit_router_route_table_description` - (Optional, Readable) - The description of the route table. This attribute is returned by the API and can be manually set.
  * `tags` - (Optional) - The tags of the transit router route table.
    
    * `tag_key` - (Optional) - The key of the tag.
    
    * `tag_value` - (Optional) - The value of the tag.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `id` - The unique identifier of the resource, formatted as `<transit_router_id>:<transit_router_route_table_id>`.
  * `create_time` - The creation time of the route table.
  * `transit_router_route_table_id` - The ID of the transit router route table.
  * `transit_router_route_table_type` - The type of the transit router route table.
  * `status` - The status of the transit router route table.

## Import

Transit Router Route Table can be imported using the `transit_router_id` and `transit_router_route_table_id` combination, formatted as `<transit_router_id>:<transit_router_route_table_id>`, e.g.

```
$ terraform import alibabacloudstack_cen_transit_router_route_table.example tr-12345678:trtb-87654321
```
