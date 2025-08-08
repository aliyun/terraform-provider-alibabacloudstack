---
subcategory: "CEN"
layout: "alibabacloudstack"
page_title: "阿里云栈: alibabacloudstack_cen_transitrouterroutetables"
sidebar_current: "docs-Alibabacloudstack-datasource-cen-transitrouterroutetables"
description: |-
  提供阿里云栈账户拥有的cen transitrouterroutetables列表。
---

# alibabacloudstack\_cen\_transitrouterroutetables

该数据源根据指定的过滤条件提供阿里云栈账户中的cen transitrouterroutetables列表。

## 使用示例
```

variable "name" {
	default = "tf-testaccrouter_table70689"
}

resource "alibabacloudstack_cen_instance" "default" {
	cen_instance_name = "${var.name}"
	description = "${var.name}"
	transit_router_name = "${var.name}"
	transit_router_description = "${var.name}"
}


resource "alibabacloudstack_cen_transit_router_route_table" "default" {
  transit_router_route_table_description = "tf-testaccrouter_table70689"
  transit_router_route_table_name = "tf-testaccrouter_table70689"
  transit_router_id = "${alibabacloudstack_cen_instance.default.transit_router_id}"
}

data "alibabacloudstack_cen_transit_router_route_tables" "default" {
	transit_router_route_id="${alibabacloudstack_cen_transit_router_route_table.default.transit_router_id}"

```


## 参数引用

支持以下参数：
  * `ids` - (可选) - 路由表的ID列表。
  * `transit_router_route_id` - (必填) - 路由表的ID。
  * `name_regex` - (可选) - 路由表的名称。
  * `description_regex` - (可选) - 路由表的描述。

## 属性引用

除了上述列出的参数外，还导出以下属性：
  * `transit_router_route_tables` - 路由表列表。
    * `id` - 路由表的ID。
    * `transit_router_route_table_name` - 路由表的名称。
    * `transit_router_route_table_description` - 路由表的描述。
    * `create_time` - 路由表的创建时间。
    * `transit_router_route_table_type` - 路由表的类型。
    * `status` - 路由表的状态。