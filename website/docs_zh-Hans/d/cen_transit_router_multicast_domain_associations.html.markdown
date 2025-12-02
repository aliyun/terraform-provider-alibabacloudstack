---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cen_transitroutermulticastdomainassociations"
sidebar_current: "docs-Alibabacloudstack-datasource-cen-transitroutermulticastdomainassociations"
description: |-
  提供阿里云账户拥有的 cen transitroutermulticastdomainassociations 列表。
---

# alibabacloudstack\_cen\_transitroutermulticastdomainassociations

此数据源根据指定的过滤条件提供阿里云账户中的 cen transitroutermulticastdomainassociations 列表。

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

data "alibabacloudstack_cen_transit_router_multicast_domain_associations" "default" {
  transit_router_multicast_domain_id = "${alibabacloudstack_cen_transit_router_multicast_domain.default.transit_router_multicast_domain_id}
}
```


## Argument Reference

支持以下参数：
  * `ids` - (可选) - cen 路由器组播域绑定交换机 ID 列表。
  * `transit_router_multicast_domain_id` - (必选) - cen 路由器组播域的 ID。
  * `vswitch_id_regex` - (可选) - 用于按交换机 ID 过滤结果的正则表达式字符串。

## Attributes Reference

除上述参数外，还导出以下属性：
  * `transit_router_multicast_associations` - cen 路由器组播域关联交换机列表。