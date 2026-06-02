---
subcategory: "API Gateway"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_api_gateway_v2_routes"
description: |-
  Provides a list of Api Gateway V2 Routes to be used by the Terraform engine.
---

# alibabacloudstack\_api\_gateway\_v2\_routes

This data source provides a list of API Gateway V2 Routes in an Alibaba Cloud API Gateway V2 instance according to the specified filters.

## Example Usage

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

## Argument Reference

The following arguments are supported:

* `gw_instance_id` - (Required) The ID of the API Gateway V2 instance.
* `is_source_route` - (Optional) Whether to query source routes. Default to `false`.
* `ids` - (Optional) A list of route IDs.
* `name_regex` - (Optional) A regex string to filter routes by name.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of route IDs.
* `routes` - A list of routes. Each element contains the following attributes:
  * `id` - The ID of the route, formatted as `{prefix}:{gwInstanceId}:{routeId}`.
  * `route_id` - The ID of the route.
  * `route_name` - The name of the route.
  * `group_id` - The ID of the route group.
  * `group_name` - The name of the route group.
  * `service_name` - The name of the service.
  * `route_path` - A list of route paths.
  * `service_id` - The ID of the service.
  * `enable_status` - Whether the route is enabled.
  * `create_time` - The creation time of the route.