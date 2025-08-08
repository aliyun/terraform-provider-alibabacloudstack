---
subcategory: "CEN"
layout: "alibabacloudstack"
page_title: "阿里云栈: alibabacloudstack_cen_transitrouterroutetable"
sidebar_current: "docs-Alibabacloudstack-cen-transitrouterroutetable"
description: |-
  提供一个cen Transitrouterroutetable资源。
---

# alibabacloudstack\_cen\_transitrouterroutetable

提供一个cen Transitrouterroutetable资源。

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
```


## 参数引用

支持以下参数：
  * `transit_router_id` - (必填, ForceNew) - 路由器ID。
  * `transit_router_route_table_name` - (可选) - 路由表名称。
  * `transit_router_route_table_description` - (可选) - 路由表描述。
  * `tags` - (可选) - 资源的标签
    
    * `tag_key` - (可选) - 标签键。
    
    * `tag_value` - (可选) - 标签值。

## 属性引用

除了上述列出的参数外，还导出以下属性：
  * `transit_router_route_table_name` - 路由表名称。
  * `transit_router_route_table_description` - 路由表描述。
  * `create_time` - 路由表的创建时间。
  * `transit_router_route_table_id` - 路由表ID。
  * `transit_router_route_table_type` - 路由表类型。
  * `status` - 路由表的状态。