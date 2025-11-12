---
subcategory: "API Gateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_api_gateway_v2_mcpserver"
sidebar_current: "docs-Alibabacloudstack-datasource-api-gateway-v2-mcpserver"
description: |-
  查询API网关v2版本的MCP服务器配置
---

# alibabacloudstack_api_gateway_v2_mcpserver

查询API网关v2版本的MCP服务器配置信息。MCP（Microservice Control Plane）是API网关用于管理微服务的核心组件，支持OPEN_API、DATABASE和DIRECT_ROUTE三种类型的服务接入。

## 示例用法

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

## 参数说明

以下参数支持作为过滤条件：

* `gw_instance_id` (字符串, 必填) - API网关实例ID，用于指定要查询的网关实例。

* `ids` (列表, 可选) - MCP服务器ID列表，格式为"{gwInstanceId}:{name}"，用于精确匹配特定的MCP服务器。

* `name_regex` (字符串, 可选) - MCP服务器名称的正则表达式，用于模糊匹配符合条件的MCP服务器。

## 属性说明

以下属性被导出：

* `id` (字符串) - MCP服务器的唯一标识符，格式为"{gwInstanceId}:{name}"。

* `consumer_auth_info` (列表) - 消费者认证信息配置。
  * `allowed_consumers` (列表) - 允许访问的消费者列表。
  * `enable` (布尔) - 是否启用消费者认证。
  * `type` (字符串) - 认证类型，如"API_KEY"。

* `db_type` (字符串) - 数据库类型，仅当type为"DATABASE"时有效，如"MYSQL"。

* `description` (字符串) - MCP服务器的描述信息。

* `direct_route_config` (列表) - 直连路由配置。
  * `path` (字符串) - 路由路径。
  * `transport_type` (字符串) - 传输类型，如"sse"。

* `domains` (列表) - 与MCP服务器关联的域名列表。

* `dsn` (字符串) - 数据库连接字符串，仅当type为"DATABASE"时有效。

* `gw_instance_id` (字符串) - 网关实例ID，与查询参数中的gw_instance_id相同。

* `name` (字符串) - MCP服务器的名称。

* `raw_configurations` (字符串) - MCP服务器的原始YAML格式配置内容。

* `service` (字符串) - 后端服务地址。

* `services` (列表) - 后端服务列表。
  * `name` (字符串) - 上游服务名称。
  * `port` (整数) - 上游服务端口。
  * `version` (字符串) - 服务版本。
  * `weight` (整数) - 服务权重，用于负载均衡。

* `type` (字符串) - MCP服务器类型，可能的值包括"OPEN_API"、"DATABASE"和"DIRECT_ROUTE"。

* `upstream_path_prefix` (字符串) - 上游请求的路径前缀。