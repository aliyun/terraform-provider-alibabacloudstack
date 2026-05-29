---
subcategory: "API Gateway V2"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_api_gateway_v2_route_group"
sidebar_current: "docs-Alibabacloudstack-resource-api-gateway-v2-route-group"
description: |-
  Create and manage API Gateway V2 route groups
---

# alibabacloudstack_api_gateway_v2_route_group

> API Gateway V2 Route Group

Create and manage route groups in the specified API Gateway instance using credentials configured in the Provider. Route groups are used to define API access paths and associated domains.

## Example Usage

### Basic Usage

```hcl

variable "name" {
  default = "tf-testacc-routegroup1529"
}

resource "alibabacloudstack_api_gateway_v2_instance" "default" {
  instance_name      = var.name
  node_number        = "1"
  instance_class     = "mini"
  broker_engine_type = "SCG"
  deploy_mode        = "custom"
}

resource "alibabacloudstack_api_gateway_v2_domain" "domain0" {
  domain      = "${var.name}1.com"
  instance_id = alibabacloudstack_api_gateway_v2_instance.default.id
  protocol    = "HTTP"
  client_auth = "0"
}

resource "alibabacloudstack_api_gateway_v2_domain" "domain1" {
  domain      = "${var.name}2.com"
  instance_id = alibabacloudstack_api_gateway_v2_instance.default.id
  protocol    = "HTTP"
  client_auth = "0"
}

resource "alibabacloudstack_api_gateway_v2_domain" "domain2" {
  domain      = "${var.name}3.com"
  instance_id = alibabacloudstack_api_gateway_v2_instance.default.id
  protocol    = "HTTP"
  client_auth = "0"
}



resource "alibabacloudstack_api_gateway_v2_route_group" "default" {
  description = var.name
  domain_ids = [
    "${alibabacloudstack_api_gateway_v2_domain.domain0.domain_id}",
    "${alibabacloudstack_api_gateway_v2_domain.domain1.domain_id}"
  ]
  instance_id = alibabacloudstack_api_gateway_v2_instance.default.id
  name        = var.name
  base_path   = "/test"
}
```

## Argument Reference

The following arguments are supported:

* `base_path` - (Required) The base path of the route group. Must start with a slash (/), such as "/api". This path will serve as the common prefix for all APIs under this group.
* `instance_id` - (Required, ForceNew) The ID of the API Gateway instance. Specifies the gateway instance to which the route group belongs
* `name` - (Required) The name of the route group. The name length is limited to 1-128 characters and cannot contain special characters.
* `description` - (Optional) The description of the route group. The description length is limited to 1-256 characters.
* `domain_ids` - (Optional) A list of associated domain IDs. Each domain ID corresponds to a configured custom domain used to access APIs under this route group.

## Attributes Reference

The following attributes are exported from the API Gateway service:

* `id` - The resource ID, formatted as "instance_id:group_id".
* `create_time` - The creation time of the route group, formatted as "YYYY-MM-DD HH:MM:SS".
* `domains` - A list of domains. Each domain contains the following attributes:
  * `protocol` - The protocol type (e.g., "HTTP" or "HTTPS").
  * `create_time` - The creation time of the domain, formatted as "YYYY-MM-DD HH:MM:SS".
  * `domain` - The full domain name (e.g., "api.example.com").
  * `domain_id` - The unique identifier ID of the domain.
* `editable` - Whether the route group is editable (boolean value). Returns `true` when the route group is in an editable state.
* `group_id` - The unique identifier ID of the route group.
* `update_time` - The last update time of the route group, formatted as "YYYY-MM-DD HH:MM:SS".

## Import

API Gateway V2 Route Group can be imported using the instance_id and group_id separated by a colon, e.g.

```
$ terraform import alibabacloudstack_api_gateway_v2_route_group.example gw-instance-123456:group-789
```