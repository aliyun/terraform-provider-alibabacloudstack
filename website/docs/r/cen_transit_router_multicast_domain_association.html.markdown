---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cen_transitroutermulticastdomainassociation"
sidebar_current: "docs-Alibabacloudstack-cen-transitroutermulticastdomainassociation"
description: |-
  Provides a cen Transitroutermulticastdomainassociation resource.
---

# alibabacloudstack\_cen\_transitroutermulticastdomainassociation

Provides a cen Transitroutermulticastdomainassociation resource.

## Example Usage
```
variable "name" {
			  default = "tf-testaccmulticast_domain_association93302"
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

resource "alibabacloudstack_cen_transit_router_multicast_domain" "default" {
	transit_router_multicast_domain_description = "${var.name}"
	transit_router_multicast_domain_name = "${var.name}"
	transit_router_id = "${alibabacloudstack_cen_instance.default.transit_router_id}"
}
			

resource "alibabacloudstack_cen_transit_router_multicast_domain_association" "default" {
  transit_router_attachment_id = "${alibabacloudstack_cen_transit_router_vpc_attachment.default.transit_router_attachment_id}"
  transit_router_multicast_domain_id = "${alibabacloudstack_cen_transit_router_multicast_domain.default.transit_router_multicast_domain_id}"
  vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
}
```

## Argument Reference

The following arguments are supported:
  * `transit_router_attachment_id` - (Required, ForceNew) - the network instance id.
  * `transit_router_multicast_domain_id` - (Required, ForceNew) - the transit router multicast domain id.
  * `vswitch_id` - (Required) - the vswitch id.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `status` - the status of the resource.
