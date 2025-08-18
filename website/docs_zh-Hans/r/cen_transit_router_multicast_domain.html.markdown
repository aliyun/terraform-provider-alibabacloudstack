---
subcategory: "CEN"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cen_transitroutermulticastdomains"
sidebar_current: "docs-Alibabacloudstack-datasource-cen-transitroutermulticastdomains"
description: |-
  提供阿里云账户拥有的 cen transitroutermulticastdomains 列表。
---

# alibabacloudstack\_cen\_transitroutermulticastdomains

该数据源根据指定的过滤条件提供阿里云账户中的 cen transitroutermulticastdomains 列表。

## 示例用法
```
variable "name" {
	default = "tf-testaccmulticast_domain37103"
}

resource "alibabacloudstack_cen_instance" "default" {
	cen_instance_name = "${var.name}"
	description = "${var.name}"
	transit_router_name = "${var.name}"
	transit_router_description = "${var.name}"
}


resource "alibabacloudstack_cen_transit_router_multicast_domain" "default" {
  transit_router_multicast_domain_description = "tf-testaccmulticast_domain37103"
  transit_router_multicast_domain_name = "tf-testaccmulticast_domain37103"
  transit_router_id = "${alibabacloudstack_cen_instance.default.transit_router_id}"
}
```


## 参数参考

支持以下参数：
  * `ids` - (可选) - cen 转发路由器组播域的 ID 列表。
  * `transit_router_route_id` - (必填) - 转发路由器路由 ID。
  * `name_regex` - (可选) - cen 转发路由器组播域的名称正则表达式。
  * `description_regex` - (可选) - cen 转发路由器组播域的描述正则表达式。

## 属性参考

除了上述参数外，还导出以下属性：
  * `transit_router_multicast_domains` - cen 转发路由器组播域列表。
    * `id` - cen 转发路由器组播域的 ID。
    * `transit_router_multicast_domain_name` - cen 转发路由器组播域的名称。
    * `transit_router_multicast_domain_description` - cen 转发路由器组播域的描述。
    * `status` - cen 转发路由器组播域的状态。
    * `transit_router_id` - 转发路由器的 ID。