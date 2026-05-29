---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cen_transit_router_multicast_domain"
sidebar_current: "docs-Alibabacloudstack-resource-cen-transit-router-multicast-domain"
description: |-
  Provides a cen Transitroutermulticastdomain resource.
---

# alibabacloudstack\_cen\_transitroutermulticastdomain

Provides a cen Transitroutermulticastdomain resource.

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
```

## Argument Reference

The following arguments are supported:
  * `transit_router_id` - (Required, ForceNew) - the id of transit router.
  * `transit_router_multicast_domain_description` - (Optional) - the description of the multicast domain.
  * `transit_router_multicast_domain_name` - (Optional) - the name of the multicast domain.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `status` - the status of the multicast domain.
  * `transit_router_multicast_domain_id` - the id of the multicast domain.
