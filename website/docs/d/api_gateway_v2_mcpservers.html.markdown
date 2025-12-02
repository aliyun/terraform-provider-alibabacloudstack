---
subcategory: "API Gateway V2"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_api_gateway_v2_mcpserver"
sidebar_current: "docs-Alibabacloudstack-datasource-api-gateway-v2-mcpserver"
description: |-
  Queries the MCP server configuration of API Gateway v2
---

# alibabacloudstack_api_gateway_v2_mcpserver

Queries the MCP server configuration information of API Gateway v2. MCP (Microservice Control Plane) is a core component of API Gateway for managing microservices, supporting three types of service access: OPEN_API, DATABASE, and DIRECT_ROUTE.

## Example Usage

```hcl
variable "name" {
  default = "tfAccmcp_4910"
}

resource "alibabacloudstack_api_gateway_v2_mcpserver" "default" {
  name           = var.name
  description    = var.name
  type           = "OPEN_API"
  service        = "kubernetes.default.svc.cluster.local"
  domains        = ["testtf.com"]
  consumer_auth  = true
  gw_instance_id = "i-a94231617cc664b8d00"
}

data "alibabacloudstack_api_gateway_v2_mcpservers" "default" {
  gw_instance_id = alibabacloudstack_api_gateway_v2_mcpserver.default.gw_instance_id
  name_regex     = alibabacloudstack_api_gateway_v2_mcpserver.default.name
}
```

## Argument Reference

The following arguments are supported as filter criteria:

* `gw_instance_id` (String, Required) - The ID of the API Gateway instance, used to specify the gateway instance to query.
* `ids` (List, Optional) - A list of MCP server IDs in the format "{gwInstanceId}:{name}", used to precisely match specific MCP servers.
* `name_regex` (String, Optional) - A regular expression for the MCP server name, used to fuzzy match MCP servers that meet the criteria.

## Attributes Reference

The following attributes are exported:

* `id` (String) - The unique identifier of the MCP server, in the format "{gwInstanceId}:{name}".
* `consumer_auth_info` (List) - Configuration for consumer authentication information. Each element contains the following attributes:
  * `allowed_consumers` (List) - A list of consumers allowed to access.
  * `enable` (Boolean) - Whether consumer authentication is enabled.
  * `type` (String) - The authentication type, such as "API_KEY".
* `db_type` (String) - The database type, valid only when type is "DATABASE", such as "MYSQL".
* `description` (String) - The description of the MCP server.
* `direct_route_config` (List) - Configuration for direct route. Each element contains the following attributes:
  * `path` (String) - The route path.
  * `transport_type` (String) - The transport type, such as "sse".
* `domains` (List) - A list of domain names associated with the MCP server.
* `dsn` (String) - The database connection string, valid only when type is "DATABASE".
* `gw_instance_id` (String) - The gateway instance ID, same as the gw_instance_id in the query parameters.
* `name` (String) - The name of the MCP server.
* `raw_configurations` (String) - The raw YAML format configuration content of the MCP server.
* `service` (String) - The backend service address.
* `services` (List) - A list of backend services. Each element contains the following attributes:
  * `name` (String) - The name of the upstream service.
  * `port` (Integer) - The port of the upstream service.
  * `version` (String) - The service version.
  * `weight` (Integer) - The service weight, used for load balancing.
* `type` (String) - The type of the MCP server, possible values include "OPEN_API", "DATABASE", and "DIRECT_ROUTE".
* `upstream_path_prefix` (String) - The path prefix for upstream requests.