---
subcategory: "云企业网"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cen_transit_router_multicast_domain_association"
sidebar_current: "docs-Alibabacloudstack-resource-cen-transit-router-multicast-domain-association"
description: |-
  提供 cen Transitroutermulticastdomainassociation 资源。
---

# alibabacloudstack\_cen\_transitroutermulticastdomainassociation

提供 cen Transitroutermulticastdomainassociation 资源。

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

支持以下参数：
  * `transit_router_attachment_id` - (必选, 变更时强制重建) - 网络实例ID。
  * `transit_router_multicast_domain_id` - (必选, 变更时强制重建) - 路由器组播域ID。
  * `vswitch_id` - (必选) - 交换机ID。

## Attributes Reference

除上述参数外，还导出以下属性：
  * `status` - 资源的状态。

## Import

CEN 转发路由器组播域关联可以通过 transit_router_attachment_id、transit_router_multicast_domain_id 和 vswitch_id 使用冒号分隔来导入，例如：

```
$ terraform import alibabacloudstack_cen_transit_router_multicast_domain_association.example <transit_router_attachment_id>:<transit_router_multicast_domain_id>:<vswitch_id>
```