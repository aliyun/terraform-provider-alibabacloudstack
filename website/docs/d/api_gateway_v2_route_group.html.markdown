---
subcategory: "API Gateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_api_gateway_v2_route_group"
sidebar_current: "docs-Alibabacloudstack-datasource-api-gateway-v2-route-group"
description: |-
  Query the route groups of API Gateway V2
---

# alibabacloudstack_api_gateway_v2_route_group

Query the route group information of API Gateway V2.

## Example Usage

```hcl

variable "name" {
  default = "tf-testacc-routegroup16079"
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

resource "alibabacloudstack_api_gateway_v2_route_group" "default" {
  name        = var.name
  base_path   = "/test"
  description = "test description"
  instance_id = alibabacloudstack_api_gateway_v2_instance.default.id
  domain_ids  = [alibabacloudstack_api_gateway_v2_domain.domain0.domain_id]
}

data "alibabacloudstack_api_gateway_v2_route_groups" "default" {
  instance_id = alibabacloudstack_api_gateway_v2_instance.default.id
  name_regex  = "tf-testacc-routegroup.*"
}

```

## Argument Reference

The following arguments are supported:

- `instance_id` (String, Required): The ID of the API Gateway instance, used to specify the API Gateway instance to query.

- `ids` (List, Optional): A list of route group IDs for filtering results, in the format {instance_id}:{group_id}.

- `name_regex` (String, Optional): A regular expression for the route group name to filter results.

## Attributes Reference

The following attributes are exported:

- `id` (String): The unique identifier of the route group, in the format {instance_id}:{group_id}.

- `base_path` (String): The base path of the route group, which is the common prefix for all APIs.

- `create_time` (String): The creation time of the route group, in the format "YYYY-MM-DD HH:MM:SS".

- `description` (String): The description of the route group.

- `domains` (List): A list of domains associated with the route group.
  - `create_time` (String): The creation time of the domain, in the format "YYYY-MM-DD HH:MM:SS".
  - `domain` (String): The domain name.
  - `domain_id` (String): The ID of the domain.
  - `protocol` (String): The protocol of the domain (HTTP/HTTPS).

- `editable` (Boolean): Whether the route group is editable. `true` means editable, `false` means not editable.

- `group_id` (String): The ID of the route group.

- `instance_id` (String): The ID of the API Gateway instance.

- `name` (String): The name of the route group.

- `update_time` (String): The update time of the route group, in the format "YYYY-MM-DD HH:MM:SS".