---
subcategory: "API Gateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_api_gateway_v2_route_group"
sidebar_current: "docs-Alibabacloudstack-datasource-api_gateway_v2_route_group"
description: |-
  查询API网关V2版本的路由分组
---

# alibabacloudstack_api_gateway_v2_route_group

查询API网关V2版本的路由分组信息。

## 示例用法

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

## 参数说明
以下参数用于过滤查询结果：

- `instance_id` (字符串, 必填)：API网关实例ID，用于指定要查询的API网关实例。

- `ids` (列表, 可选)：路由分组ID列表，用于过滤结果，格式为{instance_id}:{group_id}。

- `name_regex` (字符串, 可选)：路由分组名称的正则表达式，用于过滤结果。

## 属性说明
以下属性被导出：

- `id` (字符串)：路由分组的唯一标识符，格式为{instance_id}:{group_id}。

- `base_path` (字符串)：路由分组的基础路径，所有API的公共前缀。

- `create_time` (字符串)：路由分组的创建时间，格式为"YYYY-MM-DD HH:MM:SS"。

- `description` (字符串)：路由分组的描述信息。

- `domains` (列表)：路由分组关联的域名列表。
  - `create_time` (字符串)：域名创建时间，格式为"YYYY-MM-DD HH:MM:SS"。
  - `domain` (字符串)：域名。
  - `domain_id` (字符串)：域名ID。
  - `protocol` (字符串)：域名协议（HTTP/HTTPS）。

- `editable` (布尔)：是否可编辑，true表示可编辑，false表示不可编辑。

- `group_id` (字符串)：路由分组ID。

- `instance_id` (字符串)：API网关实例ID。

- `name` (字符串)：路由分组名称。

- `update_time` (字符串)：路由分组的更新时间，格式为"YYYY-MM-DD HH:MM:SS"。