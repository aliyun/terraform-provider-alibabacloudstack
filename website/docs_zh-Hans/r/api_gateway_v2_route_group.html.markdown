---
subcategory: "API 网关"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_api_gateway_v2_route_group"
sidebar_current: "docs-Alibabacloudstack-api-gateway-v2-route-group"
description: |-
  创建和管理API网关V2版本路由分组
---

# alibabacloudstack_api_gateway_v2_route_group

> API网关V2版本路由分组

使用Provider配置的凭证在指定的API网关实例中创建和管理路由分组。路由分组用于定义API的访问路径和关联域名。

## 示例用法

### 基础用法

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

## 参数说明

支持以下参数：

* `base_path` - (必填) 路由分组的基础路径。必须以斜杠（/）开头，例如 "/api"。该路径将作为该分组下所有API的公共前缀。
* `instance_id` - (必填, 变更时重建) API网关实例ID。指定路由分组所属的网关实例。
* `name` - (必填) 路由分组的名称。名称长度限制为1-128个字符，不能包含特殊字符。
* `description` - (可选) 路由分组的描述信息。描述长度限制为1-256个字符。
* `domain_ids` - (可选) 关联的域名ID列表。每个域名ID对应一个已配置的自定义域名，用于访问该路由分组下的API。

## 属性说明

以下属性会从API网关服务端导出：

* `id` - 资源ID，格式为 "instance_id:group_id"。
* `create_time` - 路由分组的创建时间，格式为 "YYYY-MM-DD HH:MM:SS"。
* `domains` - 域名列表。每个域名包含以下属性：
  * `protocol` - 协议类型（例如 "HTTP" 或 "HTTPS"）。
  * `create_time` - 域名的创建时间，格式为 "YYYY-MM-DD HH:MM:SS"。
  * `domain` - 完整域名（例如 "api.example.com"）。
  * `domain_id` - 域名唯一标识ID。
* `editable` - 是否可编辑（布尔值）。当路由分组处于可编辑状态时返回 `true`。
* `group_id` - 路由分组唯一标识ID。
* `update_time` - 路由分组的最后更新时间，格式为 "YYYY-MM-DD HH:MM:SS"。

## Import

API网关V2路由分组可以使用 instance_id 和 group_id（用冒号分隔）进行导入，例如：

```
$ terraform import alibabacloudstack_api_gateway_v2_route_group.example gw-instance-123456:group-789
```