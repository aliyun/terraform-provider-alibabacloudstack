---
subcategory: "云企业网"
layout: "alibabacloudstack"
page_title: "阿里云栈: alibabacloudstack_cen_transit_router_route_entry"
sidebar_current: "docs-Alibabacloudstack-resource-cen-transit-router-route-entry"
description: |-
  提供一个cen Transitrouterrouteentry资源。
---

# alibabacloudstack\_cen\_transitrouterrouteentry

提供一个cen Transitrouterrouteentry资源。

## 使用示例
```
variable "name" {
	default = "tf-testaccrouter_entry58946"
}

resource "alibabacloudstack_cen_instance" "default" {
	cen_instance_name = "${var.name}"
	description = "${var.name}"
	transit_router_name = "${var.name}"
	transit_router_description = "${var.name}"
}
resource "alibabacloudstack_cen_transit_router_route_table" "default" {
	transit_router_route_table_description = "${var.name}"
	transit_router_route_table_name = "${var.name}"
	transit_router_id = "${alibabacloudstack_cen_instance.default.transit_router_id}"
}


resource "alibabacloudstack_cen_transit_router_route_entry" "default" {
  transit_router_route_entry_next_hop_type = "BlackHole"
  transit_router_route_table_id = "${alibabacloudstack_cen_transit_router_route_table.default.transit_router_route_table_id}"
  transit_router_route_entry_description = "tf-testaccrouter_entry58946"
  transit_router_route_entry_destination_cidr_block = "10.10.10.1/32"
  transit_router_route_entry_name = "tf-testaccrouter_entry58946"
}
```


## 参数引用

支持以下参数：
  * `transit_router_route_entry_description` - (可选) - 路由条目的描述。
  * `transit_router_route_entry_destination_cidr_block` - (必填, 变更后重建) - 路由条目的目标CIDR块。
  * `transit_router_route_entry_name` - (可选) - 路由条目的名称。
  * `transit_router_route_entry_next_hop_id` - (可选) - 路由条目的下一跳ID。
  * `transit_router_route_entry_next_hop_type` - (必填, 变更后重建) - 路由条目的下一跳类型。
  * `transit_router_route_table_id` - (必填, 变更后重建) - 路由条目所属的路由表ID。

## 属性引用

除了上述列出的参数外，还导出以下属性：
  * `create_time` - 路由条目的创建时间。
  * `status` - 路由条目的状态。
  * `transit_router_route_entry_id` - 路由条目的ID。
  * `transit_router_route_entry_type` - 路由条目的类型。
  * `operational_mode` - 路由条目的操作模式。