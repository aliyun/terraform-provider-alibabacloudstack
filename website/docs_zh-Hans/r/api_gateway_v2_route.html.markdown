---
subcategory: "API Gateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_api_gateway_v2_route"
sidebar_current: "docs-Alibabacloudstack-api-gateway-api-gateway-v2-route"
description: |-
  管理API网关v2版本的路由
---

# alibabacloudstack_api_gateway_v2_route

管理API网关v2版本的路由。

## 示例用法

### 基础用法

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
    port   = "80"
    weight = "100"
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

  strip_prefix = "2"
  service_ids {
    service_id = alibabacloudstack_api_gateway_v2_service.default.service_id
    weight     = "100"
  }

  gw_instance_id = alibabacloudstack_api_gateway_v2_instance.default.id
  order          = "100"
}
```

## 参数说明

支持以下参数：

* `gw_instance_id` - (必填, 变更时重建) API网关实例ID。
* `route_name` - (必填, 变更时重建) 路由名称。
* `group_id` - (可选, 变更时重建) 分组ID。默认值为"DEFAULT"。
* `cookie` - (可选) Cookie匹配规则列表。每个规则包含以下属性：
  * `key` - (必填) Cookie名称。
  * `value` - (必填) Cookie值。
* `domain_ids` - (可选) 绑定的域名ID列表。
* `enable_status` - (可选) 是否启用路由。
* `header` - (可选) 请求头匹配规则列表。每个规则包含以下属性：
  * `key` - (必填) 请求头名称。
  * `value` - (必填) 请求头值。
* `methods` - (可选) 支持的HTTP方法列表。
* `order` - (可选) 路由优先级。
* `path` - (可选) 路径匹配规则。包含以下属性：
  * `match_type` - (可选) 匹配类型。
  * `match_value` - (可选) 匹配值。
  * `case_sensitive` - (可选) 是否区分大小写。
* `query_param` - (可选) 查询参数匹配规则列表。每个规则包含以下属性：
  * `key` - (必填) 查询参数名称。
  * `value` - (必填) 查询参数值。
* `route_path` - (可选) 路由路径列表。
* `service_id` - (可选) 后端服务ID（单后端服务模式）。
* `service_ids` - (可选) 后端服务ID列表（多后端服务模式）。每个服务包含以下属性：
  * `service_id` - (可选) 后端服务ID。
  * `weight` - (可选) 权重，取值范围1-100。
* `strip_prefix` - (可选) 需要剥离的前缀长度。当值大于0时，会自动启用前缀剥离功能。

注意：`service_ids`和`service_id`是互斥的，只能设置其中一个。

## 属性说明

除了上述所有参数外，还导出了以下属性：

* `route_id` - 路由ID。