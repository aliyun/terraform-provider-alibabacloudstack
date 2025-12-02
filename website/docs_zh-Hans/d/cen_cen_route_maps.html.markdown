---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cen_cenroutemaps"
sidebar_current: "docs-Alibabacloudstack-datasource-cen-cenroutemaps"
description: |-
  提供阿里云账户下拥有的cen路由映射列表。
---

# alibabacloudstack\_cen\_cenroutemaps

本数据源根据指定的过滤条件，提供阿里云账户下的cen路由映射列表。

## 示例用法
variable "name" { default = "tf-testAccRouteMapsDatasource123" }

resource "alibabacloudstack_cen_instance" "default" { cen_instance_name = "${var.name}" description = "${var.name}" transit_router_name = "${var.name}" transit_router_description = "${var.name}" }

data "alibabacloudstack_cen_transit_router_route_tables" "default" { transit_router_route_id="${alibabacloudstack_cen_instance.default.transit_router_id}" }

resource "alibabacloudstack_cen_route_map" "default" {

cen_id = "${alibabacloudstack_cen_instance.default.cen_id}"
transit_router_route_table_id = "${data.alibabacloudstack_cen_transit_router_route_tables.default.transit_router_route_tables.0.id}"
priority = "3"
transmit_direction = "RegionIn"
map_result = "Deny"
description = "tf_testAccCenRouteMap"
}

data "alibabacloudstack_cen_route_maps" "default" { cen_id="${alibabacloudstack_cen_instance.default.cen_id}" transit_router_route_table_id="${data.alibabacloudstack_cen_transit_router_route_tables.default.transit_router_route_tables.0.id}" }


## 参数引用

支持以下参数：
  * `ids` - (可选) - 路由策略ID列表
  * `transit_router_route_table_id` - (必填) - 路由表ID
  * `cen_id` - (必填) - cen ID
  * `description_regex` - (可选) - 路由策略映射描述正则表达式

## 属性引用

除了上述参数外，还导出以下属性：
  * `route_maps` - 路由映射列表
    * `id` - 路由策略ID
    * `route_map_id` - 路由策略ID
    * `cen_id` - cen ID
    * `transit_router_route_table_id` - 路由表ID
    * `priority` - 策略优先级
    * `map_result` - 所有匹配条件通过后的策略行为
    * `transmit_direction` - 应用方向
    * `status` - 路由策略状态
    * `description` - 路由策略描述