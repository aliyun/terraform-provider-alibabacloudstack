---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cen_cenroutemap"
sidebar_current: "docs-Alibabacloudstack-cen-cenroutemap"
description: |-
  提供一个CEN路由映射资源。
---

# alibabacloudstack\_cen\_cenroutemap

提供一个CEN路由映射资源。

## 示例用法
```
variable "name" {
	default = "tf-testaccroute_map95277"
}

resource "alibabacloudstack_cen_instance" "default" {
	cen_instance_name = "${var.name}"
	description = "${var.name}"
	transit_router_name = "${var.name}"
	transit_router_description = "${var.name}"
}

data "alibabacloudstack_cen_transit_router_route_tables" "default" {
	transit_router_route_id="${alibabacloudstack_cen_instance.default.transit_router_id}"
}



resource "alibabacloudstack_cen_route_map" "default" {
  cen_id = "${alibabacloudstack_cen_instance.default.cen_id}"
  transit_router_route_table_id = "${data.alibabacloudstack_cen_transit_router_route_tables.default.transit_router_route_tables.0.id}"
  priority = "3"
  transmit_direction = "RegionIn"
  map_result = "Deny"
}
```


## 参数引用

支持以下参数：
  * `cen_id` - (必填, ForceNew) - CEN实例ID
  * `transit_router_route_table_id` - (必填, ForceNew) - CEN路由表ID
  * `priority` - (必填) - 策略优先级
  * `transmit_direction` - (必填) - 应用方向("RegionIn", "RegionOut")
  * `map_result` - (必填) - 所有匹配条件通过后的策略行为。支持以下行为:("Permit", "Deny")
Allow: 允许匹配的路由通过。允许修改路由属性。
Reject: 拒绝匹配的路由通过。
  * `as_path_match_mode` - (可选) - AS路径匹配模式("Include", "Complete")
  * `cidr_match_mode` - (可选) - CIDR匹配模式("Include", "Complete")
  * `community_match_mode` - (可选) - Community匹配模式("Include", "Complete")
  * `community_operate_mode` - (可选) - Community操作模式("Additive", "Replace")
  * `description` - (可选) - 描述
  * `destination_instance_ids_reverse_match` - (可选) - 目标实例ID反向匹配
  * `match_address_type` - (可选) - 匹配地址类型("IPv6", "IPv4")
  * `next_priority` - (可选) - 与之关联的下一路由策略的优先级
  * `preference` - (可选) - 设置路由优先级。取值范围为1-100，默认路由优先级为50。数值越小优先级越高。
  * `source_instance_ids_reverse_match` - (可选) - 源实例ID反向匹配

## 属性引用

除了上述参数外，还导出以下属性：
  * `description` - 路由策略的描述。
  * `route_map_id` - 路由策略的ID。
  * `status` - 路由策略的状态。