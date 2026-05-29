---
subcategory: "API Gateway V2"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_api_gateway_v2_route"
sidebar_current: "docs-Alibabacloudstack-resource-api-gateway-v2-route"
description: |-
  Manages API Gateway v2 routes.
---

# alibabacloudstack_api_gateway_v2_route

Manages API Gateway v2 routes.

## Example Usage

### Basic Usage

```hcl
variable "name" {
  default = "tftestaccapiroute"
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
    port   = 80
    weight = 100
    enable = "true"
  }
}

resource "alibabacloudstack_api_gateway_v2_domain" "default" {
  domain      = "${var.name}.com"
  instance_id = alibabacloudstack_api_gateway_v2_instance.default.id
  protocol    = "HTTP"
}

resource "alibabacloudstack_api_gateway_v2_route" "default" {
  methods = [
    "GET",
    "POST",
    "PUT",
    "DELETE"
  ]
  header {
    key   = "header"
    value = "aaaaa"
  }

  domain_ids = [
    "${alibabacloudstack_api_gateway_v2_domain.default.domain_id}"
  ]
  route_name = var.name
  route_path = [
    "/testtc",
    "/test/aaa"
  ]
  cookie {
    key   = "cookie"
    value = "bbbbb"
  }

  query_param {
    key   = "query"
    value = "ccccc"
  }

  strip_prefix = 2
  service_ids {
    service_id = alibabacloudstack_api_gateway_v2_service.default.service_id
    weight     = 100
  }

  gw_instance_id = alibabacloudstack_api_gateway_v2_instance.default.id
  order          = 100
}
```

## Argument Reference

The following arguments are supported:

* `gw_instance_id` - (Required, ForceNew) The ID of the API Gateway instance.
* `route_name` - (Required, ForceNew) The name of the route.
* `group_id` - (Optional, ForceNew) The group ID. Default value is "DEFAULT".
* `cookie` - (Optional) A list of cookie matching rules. Each rule contains the following attributes:
  * `key` - (Required) The name of the cookie.
  * `value` - (Required) The value of the cookie.
* `domain_ids` - (Optional, Computed) A list of domain IDs to bind.
* `enable_status` - (Optional, Computed) Whether to enable the route.
* `header` - (Optional) A list of request header matching rules. Each rule contains the following attributes:
  * `key` - (Required) The name of the request header.
  * `value` - (Required) The value of the request header.
* `methods` - (Optional) A list of supported HTTP methods.
* `order` - (Optional) The priority of the route.
* `path` - (Optional) The path matching rule. Contains the following attributes:
  * `match_type` - (Optional) The matching type.
  * `match_value` - (Optional) The matching value.
  * `case_sensitive` - (Optional) Whether the matching is case-sensitive.
* `query_param` - (Optional) A list of query parameter matching rules. Each rule contains the following attributes:
  * `key` - (Required) The name of the query parameter.
  * `value` - (Required) The value of the query parameter.
* `route_path` - (Optional) A list of route paths.
* `service_id` - (Optional) The backend service ID (for single backend service mode).
* `service_ids` - (Optional) A list of backend service IDs (for multiple backend service mode). Each service contains the following attributes:
  * `service_id` - (Optional) The backend service ID.
  * `weight` - (Optional) The weight, ranging from 1 to 100.
* `strip_prefix` - (Optional) The length of the prefix to strip. When the value is greater than 0, the prefix stripping feature is automatically enabled.
* `cascade_link_ids` - (Optional, ForceNew) The list of cascade link IDs, used for CSB authentication. Set this attribute to create the source route. 

-> **NOTE:** When `cascade_link_ids` is not empty, only the `service_id` attribute is supported for setting the server.

-> **NOTE:**  `service_ids` and `service_id` are mutually exclusive; only one of them can be set.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `route_id` - The ID of the route.

## Import

API Gateway V2 Route can be imported using the resource ID in the format `<type>:<gw_instance_id>:<route_id>`, e.g.

```
$ terraform import alibabacloudstack_api_gateway_v2_route.example route:gw-12345678:route-87654321
```

-> **NOTE:** If the route is a source route (created with `cascade_link_ids`), the type prefix should be `sourceRoute` instead of `route`.