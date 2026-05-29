---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cen_transit_router_vpc_attachment"
sidebar_current: "docs-Alibabacloudstack-resource-cen-transit-router-vpc-attachment"
description: |-
  Provides a cen Transitroutervpcattachment resource.
---

# alibabacloudstack\_cen\_transitroutervpcattachment

Provides a cen Transitroutervpcattachment resource.

## Example Usage
```
variable "name" {
	default = "tf-testaccrouter_vpc_attachment33915"
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




resource "alibabacloudstack_vpc_vswitch" "vswitchv2" {
	name = "${var.name}v2"
	vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
	cidr_block = "172.16.0.0/24"
	zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
  }

resource "alibabacloudstack_cen_instance" "default" {
	cen_instance_name = "${var.name}"
	description = "${var.name}"
	transit_router_name = "${var.name}"
	transit_router_description = "${var.name}"
}



resource "alibabacloudstack_cen_transit_router_vpc_attachment" "default" {
  zone_mappings {
    vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
    zone_id = "${alibabacloudstack_vpc_vswitch.default.zone_id}"
  }
  
  auto_create_vpc_route = "true"
  route_table_association_enabled = "true"
  route_table_propagation_enabled = "true"
  cen_id = "${alibabacloudstack_cen_instance.default.id}"
  vpc_id = "${alibabacloudstack_vpc_vswitch.default.vpc_id}"
  transit_router_id = "${alibabacloudstack_cen_instance.default.transit_router_id}"
}
```

## Argument Reference

The following arguments are supported:
  * `cen_id` - (Required) - the cen
  * `tags` - (Optional) - The tag of the resource
    
    * `tag_key` - (Optional) - tag key
    
    * `tag_value` - (Optional) - tag value
  * `transit_router_attachment_description` - (Optional) - the description of the vpc attachment
  * `transit_router_id` - (Required) -  the transit router id
  * `transit_router_attachment_name` - (Optional) - the name of the vpc attachment
  * `auto_create_vpc_route` - (Optional) - auto create vpc route entry
  * `route_table_propagation_enabled` - (Optional) - auto create route table propagation
  * `route_table_association_enabled` - (Optional) - auto create route table association
  * `vpc_id` - (Required) - the vpc id
  * `zone_mappings` - (Required) - ZoneMappingss
    
    * `vswitch_id` - (Required) - the vswitch id
    
    * `zone_id` - (Required) - the zone id
    

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `auto_publish_route_enabled` - the auto publish route enabled
  * `charge_type` - the charge type
  * `creation_time` - the creation time
  * `resource_type` - the resource type
  * `status` - the status of vpc attachment instance
  * `transit_router_attachment_id` - the id of vpc attachment instance
  * `vpc_owner_id` - the vpc owner id
