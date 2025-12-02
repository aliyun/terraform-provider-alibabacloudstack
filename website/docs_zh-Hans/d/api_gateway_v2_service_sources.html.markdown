---
subcategory: "API Gateway V2"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_api_gateway_v2_service_sources"
sidebar_current: "docs-alibabacloudstack-datasource-api-gateway-v2-service-sources"
description: |-
    获取 Alibaba Cloud API 网关 V2 服务来源列表
---

# alibabacloudstack_api_gateway_v2_service_sources

> 查询阿里云API网关v2版本的服务来源

## 示例用法

```hcl

variable "name" {
  default = "tf-testacc-serviceSource19103"
}

resource "alibabacloudstack_api_gateway_v2_instance" "default" {
  instance_name      = var.name
  node_number        = "1"
  instance_class     = "mini"
  broker_engine_type = "SCG"
  deploy_mode        = "custom"
}

resource "alibabacloudstack_api_gateway_v2_service_source" "default" {
  source_name      = var.name
  source_type      = "1"
  instance_id      = alibabacloudstack_api_gateway_v2_instance.default.id
  description      = var.name
  check_type       = "1"
  nacos_access_key = "root"
  nacos_secret_key = "12345"
  nacos_registry   = "127.0.0.1:8000"
}

data "alibabacloudstack_api_gateway_v2_service_sources" "default" {
  instance_id = alibabacloudstack_api_gateway_v2_service_source.default.instance_id

  name_regex = "tf-testacc-serviceSource*"
}

```

## 参数说明
以下参数支持过滤查询结果：

- `instance_id` (字符串, 必填)：API网关实例ID，用于指定要查询的服务来源所属的API网关实例。

- `ids` (列表, 可选)：服务来源ID列表，用于过滤特定的服务来源。ID格式为{instance_id:source_id}。

- `name_regex` (字符串, 可选)：服务来源名称的正则表达式，用于过滤名称匹配的服务来源。

## 属性说明
以下属性被导出：

- `id` (字符串)：数据源ID，由服务来源ID列表生成的哈希值。

- `names` (列表)：匹配的服务来源名称列表。

- `sources` (列表)：匹配的服务来源列表，每个服务来源包含以下属性：
  - `create_time` (字符串)：服务来源的创建时间。
  - `description` (字符串)：服务来源的描述信息。
  - `id` (字符串)：服务来源ID，格式为{instance_id:source_id}。
  - `instance_id` (字符串)：API网关实例ID。
  - `source_id` (字符串)：服务来源的唯一标识ID。
  - `source_name` (字符串)：服务来源的名称。
  - `source_type` (字符串)：服务来源的类型编码（1: NACOS, 2: 微服务空间, 3: Eureka, 5: Database）。
  - `source_type_name` (字符串)：服务来源的类型名称（如"NACOS"）。
  - `update_time` (字符串)：服务来源的最后更新时间。