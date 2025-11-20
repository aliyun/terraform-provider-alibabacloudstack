---
subcategory: "API 网关"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_api_gateway_v2_routes"
sidebar_current: "docs-alibabacloudstack-datasource-api-gateway-v2-routes"
description: |-
  获取 Alibaba Cloud API 网关 V2 路由列表
---

# alibabacloudstack\_api\_gateway\_v2\_routes

该数据源提供 Alibaba Cloud API 网关 V2 实例中符合条件的路由列表。

## 示例用法

```hcl
variable "name" {
  default = "tf-testacc-route-12345"
}

resource "alibabacloudstack_api_gateway_v2_instance" "default" {
  instance_name      = var.name
  node_number        = "1"
  instance_class     = "mini"
  broker_engine_type = "SCG"
  deploy_mode        = "custom"
}

resource "alibabacloudstack_api_gateway_v2_service" "default" {
  name              = var.name
  description       = var.name
  protocol          = "HTTP"
  upstream_type     = "1"
  load_balance_type = "1"
  gw_instance_id    = alibabacloudstack_api_gateway_v2_instance.default.id
  service_nodes {
    ip     = "127.0.0.1"
    port   = "80"
    weight = "100"
    enable = "true"
  }
}

resource "alibabacloudstack_api_gateway_v2_route" "default" {
  gw_instance_id = alibabacloudstack_api_gateway_v2_instance.default.id
  route_name     = var.name
  strip_prefix   = "2"
  order          = "100"
  route_path     = ["/test", "/api"]
  methods        = ["GET", "POST"]
  service_id     = alibabacloudstack_api_gateway_v2_service.default.service_id
}

data "alibabacloudstack_api_gateway_v2_routes" "default" {
  gw_instance_id = alibabacloudstack_api_gateway_v2_route.default.gw_instance_id
  name_regex     = "tf-testacc-route-*"
  ids            = [alibabacloudstack_api_gateway_v2_route.default.id]
}

output "route_info" {
  value = data.alibabacloudstack_api_gateway_v2_routes.default.routes
}
```

## 参数说明

以下参数支持入参配置：

* `gw_instance_id` - (必选) API 网关 V2 实例的 ID。
* `is_source_route` - (可选) 是否查询源路由。默认为 `false`。
* `ids` - (可选) 路由 ID 列表。
* `name_regex` - (可选) 用于根据名称过滤路由的正则表达式。

## 属性说明

以下属性会被导出：

* `ids` - 路由 ID 列表。
* `routes` - 路由列表。每个元素包含以下属性：
  * `id` - 路由 ID，格式为 `{prefix}:{gwInstanceId}:{routeId}`。
  * `route_id` - 路由的 ID。
  * `route_name` - 路由的名称。
  * `group_id` - 路由组的 ID。
  * `group_name` - 路由组的名称。
  * `service_name` - 服务的名称。
  * `route_path` - 路由路径列表。
  * `service_id` - 服务的 ID。
  * `enable_status` - 路由是否启用。
  * `create_time` - 路由的创建时间。