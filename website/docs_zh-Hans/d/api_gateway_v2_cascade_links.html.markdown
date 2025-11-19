---
subcategory: "API Gateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_api_gateway_v2_cascade_link"
sidebar_current: "docs-Alibabacloudstack-datasource-api_gateway_v2_cascade_link"
description: |-
  查询API网关v2版本的级联链路
---

# alibabacloudstack_api_gateway_v2_cascade_link

> API网关v2版本级联链路数据源，用于查询和管理API网关级联链路资源

## 示例用法

```hcl

variable "name" {
  default = "tf-testAccApiGwV23033395573387097040"
}

resource "alibabacloudstack_api_gateway_v2_instance" "source" {
  instance_name      = "${var.name}-source"
  node_number        = 1
  instance_class     = "mini"
  broker_engine_type = "SCG"
  deploy_mode        = "custom"
}

resource "alibabacloudstack_api_gateway_v2_instance" "cascade" {
  instance_name      = "${var.name}-cascade"
  node_number        = 1
  instance_class     = "mini"
  broker_engine_type = "SCG"
  deploy_mode        = "custom"
}

resource "alibabacloudstack_api_gateway_v2_cascade_instance" "default" {
  instance_name       = var.name
  cascade_instance_id = alibabacloudstack_api_gateway_v2_instance.cascade.id
}

resource "alibabacloudstack_api_gateway_v2_cascade_link" "default" {
  source_instance_id      = alibabacloudstack_api_gateway_v2_instance.source.id
  source_instance_address = "10.17.94.180"
  cascade_instance_id     = alibabacloudstack_api_gateway_v2_cascade_instance.default.id
  link_name               = var.name
}

data "alibabacloudstack_api_gateway_v2_cascade_links" "default" {
  name_regex = alibabacloudstack_api_gateway_v2_cascade_link.default.link_name
}

```

## 参数说明

以下参数支持过滤查询结果：

* `cascade_instance_name` (可选)：级联实例名称，用于过滤特定级联实例的链路。

* `ids` (可选)：级联链路ID列表，用于精确匹配指定ID的级联链路。

* `name_regex` (可选)：级联链路名称的正则表达式，用于按名称模式过滤级联链路。

* `source_instance_name` (可选)：源实例名称，用于过滤特定源实例的链路。

## 属性说明

以下属性被导出：

* `id` (字符串)：级联链路的唯一标识符，等同于link_id。

* `cascade_instance_id` (字符串)：级联实例ID。

* `cascade_instance_name` (字符串)：级联实例名称。

* `cascade_service_id` (字符串)：级联服务ID。

* `link_id` (字符串)：级联链路ID。

* `link_name` (字符串)：级联链路名称。

* `source_instance_address` (字符串)：源实例地址。

* `source_instance_id` (字符串)：源实例ID。

* `source_instance_name` (字符串)：源实例名称。