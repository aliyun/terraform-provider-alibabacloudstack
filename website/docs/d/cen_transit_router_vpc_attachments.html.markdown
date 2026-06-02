---
subcategory: "Cloud Enterprise Network"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cen_transit_router_vpc_attachments"
sidebar_current: "docs-Alibabacloudstack-datasource-cen-transit-router-vpc-attachments"
description: |-
  Provides a list of cen transitroutervpcattachments owned by an alibabacloudstack account.
---

# alibabacloudstack\_cen\_transitroutervpcattachments

This data source provides a list of cen transitroutervpcattachments in an alibabacloudstack account according to the specified filters.

## Example Usage
```
variable "name" {
  default = "tf-testAccRouterVpcAttachmentsDatasource18484"
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

data "alibabacloudstack_cen_transit_router_vpc_attachments" "default" {
	cen_id = "${alibabacloudstack_cen_instance.default.id}"
	vpc_id = "${alibabacloudstack_vpc_vswitch.default.vpc_id}"
	transit_router_id = "${alibabacloudstack_cen_instance.default.transit_router_id}"
	name_regex = "${alibabacloudstack_cen_transit_router_vpc_attachment.default.transit_router_attachment_name}"
}
```

## Argument Reference

The following arguments are supported:
  * `ids` - (Optional) - the ids of the vpc attachment
  * `cen_id` - (Required) - the cen instance id
  * `vpc_id` - (Required) - the vpc id
  * `tags` - (Optional) - The tag of the resource
    
    * `tag_key` - (Optional) - tag key
    
    * `tag_value` - (Optional) - tag value
  * `transit_router_id` - (Required) - the transit router id
  * `name_regex` - (Optional) - the name regex of the vpc attachment
  * `description_regex` - (Optional) - the description regex of the vpc attachment

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `transitrouterattachments` - the list of vpc attachment
