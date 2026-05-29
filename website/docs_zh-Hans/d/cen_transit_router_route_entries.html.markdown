---
subcategory: "云企业网"
layout: "alibabacloudstack"
page_title: "阿里云栈: alibabacloudstack_cen_transit_router_route_entries"
sidebar_current: "docs-Alibabacloudstack-datasource-cen-cen-transit-router-route-entries"
description: |-
  提供阿里云栈账户拥有的cen cen_transit_router_route_entries列表。
---

# alibabacloudstack\_cen\cen_transit_router_route_entries

该数据源根据指定的过滤条件提供阿里云栈账户中的路由条目列表。

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

data "alibabacloudstack_cen_transit_router_route_entries" "default" {
transit_router_route_table_id="${alibabacloudstack_cen_transit_router_route_table.default.transit_router_route_table_id}"
name_regex = "${alibabacloudstack_cen_transit_router_route_entry.default.transit_router_route_entry_name}"
```


## 参数引用

支持以下参数：
  * `ids` - (可选) - 路由条目的ID列表。
  * `transit_router_route_table_id` - (必填) - 路由表ID
  * `name_regex` - (可选) - 路由条目的名称。
  * `description_regex` - (可选) - 路由条目的描述。

## 属性引用

除了上述列出的参数外，还导出以下属性：
  * `transit_router_route_entries` - 路由条目列表。
    * `id` - 路由条目的ID。
    * `create_time` - 路由条目的创建时间。
    * `status` - 路由条目的状态。
    * `transit_router_route_entry_description` - 路由条目的描述。
    * `transit_router_route_entry_destination_cidr_block` - 路由条目的目标CIDR块。
    * `transit_router_route_entry_id` - 路由条目的ID。
    * `transit_router_route_entry_name` - 路由条目的名称。
    * `transit_router_route_entry_next_hop_id` - 路由条目的下一跳ID。
    * `transit_router_route_entry_next_hop_type` - 路由条目的下一跳类型。
    * `transit_router_route_entry_type` - 路由条目的类型。
    * `operational_mode` - 路由条目的操作模式。
    * `transit_router_route_entry_status` - 路由条目的状态。