---
subcategory: "云企业网"
layout: "alibabacloudstack"
page_title: "阿里云栈: alibabacloudstack_cen_transit_router_route_table"
sidebar_current: "docs-Alibabacloudstack-cen-transit-router-route-table"
description: |-
  提供一个 CEN 转发路由器路由表资源。
---

# alibabacloudstack\_cen\_transit\_router\_route\_table

提供一个 CEN 转发路由器路由表资源。

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
  * `transit_router_id` - (必填, ForceNew) - 转发路由器ID。修改此参数会强制重新创建资源。
  * `transit_router_route_table_name` - (可选, 可回读) - 路由表名称。此属性由 API 返回，可手动设置。
  * `transit_router_route_table_description` - (可选, 可回读) - 路由表描述。此属性由 API 返回，可手动设置。
  * `tags` - (可选) - 转发路由器路由表的标签。
    
    * `tag_key` - (可选) - 标签键。
    
    * `tag_value` - (可选) - 标签值。

## 属性引用

除了上述列出的参数外，还导出以下属性：
  * `id` - 资源的唯一标识，格式为 `<transit_router_id>:<transit_router_route_table_id>`。
  * `create_time` - 路由表的创建时间。
  * `transit_router_route_table_id` - 路由表ID。
  * `transit_router_route_table_type` - 路由表类型。
  * `status` - 路由表的状态。

## Import

转发路由器路由表可以使用 `transit_router_id` 和 `transit_router_route_table_id` 组合导入，格式为 `<transit_router_id>:<transit_router_route_table_id>`，例如：

```
$ terraform import alibabacloudstack_cen_transit_router_route_table.example tr-12345678:trtb-87654321
```