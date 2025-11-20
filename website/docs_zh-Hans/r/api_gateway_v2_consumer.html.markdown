---
subcategory: "API Gateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_api_gateway_v2_consumer"
sidebar_current: "docs-Alibabacloudstack-api-gateway-api_gateway_v2_consumer"
description: |-
  管理API网关V2版本的消费者
---

# alibabacloudstack_api_gateway_v2_consumer

管理API网关V2版本的消费者，用于配置不同认证方式的应用访问凭证。

## 示例用法

### 基础用法

```hcl

variable "name" {
  default = "tf-testacc32955"
}

data "alibabacloudstack_api_gateway_v2_instance_types" "default" {
  sorted_by = "CPU"
}

resource "alibabacloudstack_api_gateway_v2_instance" "default" {
  instance_name      = var.name
  node_number        = "1"
  instance_class     = data.alibabacloudstack_api_gateway_v2_instance_types.default.instance_types[0].id
  broker_engine_type = "SCG"
  deploy_mode        = "custom"
}

// Basic认证
resource "alibabacloudstack_api_gateway_v2_consumer" "default" {
  auth_type      = 1
  app_name       = "${var.name}"
  description    = "${var.name}"
  key            = "root"
  password       = "admin"
  gw_instance_id = "${alibabacloudstack_api_gateway_v2_instance.default.id}"
  groups         = ["test"]
}
```

### OAuth2.0认证
```hcl
resource "alibabacloudstack_api_gateway_v2_consumer" "default" {
  gw_instance_id = alibabacloudstack_api_gateway_v2_instance.default.id
  groups = [
    "test"
  ]
  auth_type   = "2"
  app_name    = var.name
  description = var.name
  oauth2_payload = {
    "authorization_code": true,
    "client_credentials": true,
    "implicit_grant": true,
    "password_grant": true,
    "token_expiration": 3600,
    "refresh_token_expiration": 30,
    "pkce": true,
    "scopes": "api_gateway",
    "client_id": "a94231617c1a94cd900",
  }
}
```

### JWT认证
```hcl
resource "alibabacloudstack_api_gateway_v2_consumer" "default" {
  auth_type      = 3
  app_name       = "${var.name}"
  description    = "${var.name}"
  key            = "aaaaaaaaaaaaaaaaaa"
  expire_time    = 86400000
  app_secret     = "bbbbbbbbbbbbbbbbbbbbb"
  gw_instance_id = "${alibabacloudstack_api_gateway_v2_instance.default.id}"
  groups         = ["test"]
  payload = {
    issuer  = "http://127.0.0.1:8000/test"
    subject = "test.app"
  }
}
```

### API Key认证
```hcl
resource "alibabacloudstack_api_gateway_v2_consumer" "default" {
  gw_instance_id = alibabacloudstack_api_gateway_v2_instance.default.id
  groups = [
    "test"
  ]
  auth_type   = "5"
  app_name    = var.name
  description = var.name
  key         = "a94231617c1a94cd900"
}
```

### API Gateway应用认证
```hcl
resource "alibabacloudstack_api_gateway_v2_consumer" "default" {
  auth_type      = 6
  app_name       = "${var.name}"
  description    = "${var.name}"
  app_secret     = "vvvvvvvvvvvvvvvvvvvvv"
	app_code       = "aaaaaaaaaaa"
	key            = "aaaaaaaaaaaaaaaaaa"
  gw_instance_id = "${alibabacloudstack_api_gateway_v2_instance.default.id}"
  groups         = ["test"]
}
```

### CSB认证
```hcl
resource "alibabacloudstack_api_gateway_v2_consumer" "default" {
  auth_type      = 7
  app_name       = "${var.name}"
  description    = "${var.name}"
  key            = "aaaaaaaaaaaaaaaaaa"
  gw_instance_id = "${alibabacloudstack_api_gateway_v2_instance.default.id}"
  groups         = ["test"]
}
```

## 参数说明

支持以下参数：

* `app_name` - (必填) 应用名称，长度为1-128个字符。
* `auth_type` - (必填, 变更时重建) 认证方式类型，取值范围：
  * `1`：Basic认证
  * `2`：OAuth2.0认证
  * `3`：JWT认证
  * `5`：API Key认证
  * `6`：API Gateway应用认证
  * `7`：CSB认证
* `gw_instance_id` - (必填, 变更时重建) 网关实例ID，格式为`i-xxx`。
* `app_code` - (可选) API Gateway认证方式的应用代码。
* `app_secret` - (可选) 应用密钥，用于API Key、JWT等认证方式。
* `description` - (可选) 应用描述信息，长度为1-256个字符。
* `expire_time` - (可选) 令牌过期时间（毫秒），默认值根据认证方式不同而不同。
* `groups` - (可选) 应用所属的分组列表，用于权限控制。
* `key` - (可选) API Key认证方式的密钥，或Basic认证方式的用户名。
* `oauth2_payload` - (可选) OAuth2.0认证方式的配置信息，包含以下子参数：
  * `authorization_code` - (可选) 是否启用授权码模式。
  * `client_credentials` - (可选) 是否启用客户端凭证模式。
  * `implicit_grant` - (可选) 是否启用隐式授权模式。
  * `password_grant` - (可选) 是否启用密码凭证模式。
  * `token_expiration` - (可选) 访问令牌有效期（秒）。
  * `refresh_token_expiration` - (可选) 刷新令牌有效期（天）。
  * `pkce` - (可选) 是否启用PKCE扩展。
  * `scopes` - (可选) 授权范围，多个范围用逗号分隔。
  * `client_id` - (可选) 客户端ID。
  * `client_secret` - (可选) 客户端密钥。
  * `redirect_uris` - (可选) 重定向URI，多个URI用逗号分隔。
* `password` - (可选) Basic认证方式的密码。
* `payload` - (可选) JWT认证方式的负载信息，格式为键值对。
* `cascade_link_ids` - (可选) 级联链接的ID列表, 设置该参数以创建源消费者.

## 属性说明

以下属性会从API中导出：

* `id` - 资源ID，格式为`{gwInstanceId:appId}`。
* `access_key` - CSB认证方式的访问密钥。
* `app_id` - 应用ID，由系统生成的唯一标识。
* `auth_type_name` - 认证方式名称，如"API_KEY"、"BASIC"等。
* `enable` - 应用是否启用，`true`表示启用，`false`表示禁用。
* `request_header` - 请求头信息，用于OAuth2.0认证。
* `secret_key` - CSB认证方式的密钥。
* `token` - 生成的访问令牌，不同认证方式生成的令牌格式不同。
* `use_white_list` - 是否启用白名单，`true`表示启用，`false`表示禁用。