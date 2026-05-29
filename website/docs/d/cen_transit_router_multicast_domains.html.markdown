---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cen_transit_router_multicast_domains"
sidebar_current: "docs-Alibabacloudstack-datasource-cen-transitroutermulticastdomains"
description: |-
  Provides a list of cen transitroutermulticastdomains owned by an alibabacloudstack account.
---

# alibabacloudstack\_cen\_transitroutermulticastdomains

This data source provides a list of cen transitroutermulticastdomains in an alibabacloudstack account according to the specified filters.

## Example Usage
```
variable "name" {
	default = "tf-testaccmulticast_domain37103"
}

resource "alibabacloudstack_cen_instance" "default" {
	cen_instance_name = "${var.name}"
	description = "${var.name}"
	transit_router_name = "${var.name}"
	transit_router_description = "${var.name}"
}


resource "alibabacloudstack_cen_transit_router_multicast_domain" "default" {
  transit_router_multicast_domain_description = "tf-testaccmulticast_domain37103"
  transit_router_multicast_domain_name = "tf-testaccmulticast_domain37103"
  transit_router_id = "${alibabacloudstack_cen_instance.default.transit_router_id}"
}

data "alibabacloudstack_cen_transit_router_multicast_domains" "default" {
	transit_router_route_id="${alibabacloudstack_cen_transit_router_multicast_domain.default.transit_router_id}"
}
```

## Argument Reference

The following arguments are supported:
  * `ids` - (Optional) - the ids of the cen transit router multicast domains.
  * `transit_router_route_id` - (Required) -the transit router route id.
  * `name_regex` - (Optional) - the name regex of the cen transit router multicast domains.
  * `description_regex` - (Optional) - the description regex of the cen transit router multicast domains.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `transit_router_multicast_domains` - the list of cen transit router multicast domains.
    * `id` - the id of the cen transit router multicast domain.
    * `transit_router_multicast_domain_name` - the name of the cen transit router multicast domain.
    * `transit_router_multicast_domain_description` -the description of the cen transit router multicast domain.
    * `status` - the status of the cen transit router multicast domain.
    * `transit_router_id` - the id of the cen transit router.
