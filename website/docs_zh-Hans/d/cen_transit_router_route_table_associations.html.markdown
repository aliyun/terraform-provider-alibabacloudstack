---
subcategory: "云企业网"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cen_transit_router_route_table_associations"
sidebar_current: "docs-Alibabacloudstack-datasource-cen-transitrouterroutetableassociations"
description: |-
  提供阿里云账户拥有的cen transitrouterroutetableassociations列表。
---

# alibabacloudstack\_cen\_transitrouterroutetableassociations

该数据源根据指定的过滤条件提供阿里云账户中的cen transitrouterroutetableassociations列表。

## Example Usage
```
variable "name" {
	default = "tf-testaccrouter_table_association32777"
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



resource "alibabacloudstack_cen_transit_router_route_table_association" "default" {
  transit_router_route_table_id = "${data.alibabacloudstack_cen_transit_router_route_tables.default.transit_router_route_tables.0.id}"
  transit_router_attachment_id = "${alibabacloudstack_cen_transit_router_vpc_attachment.default.transit_router_attachment_id}"
}

data "alibabacloudstack_cen_transit_router_route_table_associations" "default" {
  transit_router_route_table_id = "${alibabacloudstack_cen_transit_router_route_table_association.default.transit_router_route_table_id}"
}
```


## Argument Reference

支持以下参数：
  * `ids` - (可选) - cen转发路由器路由表关联的ID列表。
  * `transit_router_route_table_id` - (必选) - cen转发路由器路由表的ID。

## Attributes Reference

除了上述参数外，还导出以下属性：
  * `transit_router_route_associations` - cen转发路由器路由表关联列表。