---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cen_transit_router_route_table_propagation"
sidebar_current: "docs-Alibabacloudstack-cen-transitrouterroutetablepropagation"
description: |-
  Provides a cen Transitrouterroutetablepropagation resource.
---

# alibabacloudstack\_cen\_transitrouterroutetablepropagation

Provides a cen Transitrouterroutetablepropagation resource.

## Example Usage
```
variable "name" {
	default = "tf-testaccrouter_table_propagation68221"
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



resource "alibabacloudstack_cen_transit_router_route_table_propagation" "default" {
  transit_router_route_table_id = "${data.alibabacloudstack_cen_transit_router_route_tables.default.transit_router_route_tables.0.id}"
  transit_router_attachment_id = "${alibabacloudstack_cen_transit_router_vpc_attachment.default.transit_router_attachment_id}"
}
```

## Argument Reference

The following arguments are supported:
  * `transit_router_attachment_id` - (Required) - the id of the transit router attachment.
  * `transit_router_route_table_id` - (Required, ForceNew) - the id of the transit router route table.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `resource_id` - the id of the resource.
  * `resource_type` - the type of the resource.
  * `status` - the status of the Propagation.

## Import

CEN Transit Router Route Table Propagation can be imported using the transit_router_route_table_id and transit_router_attachment_id separated by a colon, e.g.

```
$ terraform import alibabacloudstack_cen_transit_router_route_table_propagation.example <transit_router_route_table_id>:<transit_router_attachment_id>
```
