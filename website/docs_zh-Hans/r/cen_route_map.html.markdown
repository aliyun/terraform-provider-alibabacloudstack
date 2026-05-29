---
subcategory: "云企业网"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cen_route_map"
sidebar_current: "docs-Alibabacloudstack-resource-cen-route-map"
description: |-
  提供一个CEN路由映射资源�?
---

# alibabacloudstack\_cen\_route_map

提供一个CEN路由映射资源�?

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

支持以下参数�?
  * `cen_id` - (必填, ForceNew) - CEN实例ID
  * `transit_router_route_table_id` - (必填, ForceNew) - CEN路由表ID
  * `priority` - (必填) - 策略优先�?
  * `transmit_direction` - (必填) - 应用方向("RegionIn", "RegionOut")
  * `map_result` - (必填) - 所有匹配条件通过后的策略行为。支持以下行�?("Permit", "Deny")
Allow: 允许匹配的路由通过。允许修改路由属性�?
Reject: 拒绝匹配的路由通过�?
  * `as_path_match_mode` - (可�? - AS路径匹配模式("Include", "Complete")
  * `cidr_match_mode` - (可�? - CIDR匹配模式("Include", "Complete")
  * `community_match_mode` - (可�? - Community匹配模式("Include", "Complete")
  * `community_operate_mode` - (可�? - Community操作模式("Additive", "Replace")
  * `description` - (可�? - 描述
  * `destination_instance_ids_reverse_match` - (可�? - 目标实例ID反向匹配
  * `match_address_type` - (可�? - 匹配地址类型("IPv6", "IPv4")
  * `next_priority` - (可�? - 与之关联的下一路由策略的优先级
  * `preference` - (可�? - 设置路由优先级。取值范围为1-100，默认路由优先级�?0。数值越小优先级越高�?
  * `source_instance_ids_reverse_match` - (可�? - 源实例ID反向匹配

## 属性引�?

除了上述参数外，还导出以下属性：
  * `description` - 路由策略的描述�?
  * `route_map_id` - 路由策略的ID�?
  * `status` - 路由策略的状态