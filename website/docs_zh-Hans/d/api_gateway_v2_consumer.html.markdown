---
subcategory: "API 网关（API Gateway）V2 版"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_api_gateway_v2_consumer"
sidebar_current: "docs-Alibabacloudstack-datasource-api-gateway-v2-consumer"
description: |-
  查询API网关v2的消费者信息
---

# alibabacloudstack_api_gateway_v2_consumer

查询API网关v2的消费者信息。API网关v2消费者是用于访问API的认证主体，支持多种认证方式，包括API Key、Basic Auth、JWT、OAuth2.0等。

## 示例用法

```hcl

variable "name" {
  default = "tf-testacc2388364"
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

resource "alibabacloudstack_api_gateway_v2_consumer" "default" {
  auth_type      = "1"
  app_name       = var.name
  description    = var.name
  key            = "root"
  password       = "admin"
  gw_instance_id = alibabacloudstack_api_gateway_v2_instance.default.id
  groups         = ["test"]
}

data "alibabacloudstack_api_gateway_v2_consumers" "default" {
  gw_instance_id = alibabacloudstack_api_gateway_v2_consumer.default.gw_instance_id
  appid          = alibabacloudstack_api_gateway_v2_consumer.default.app_id
}

```

## 参数说明
以下参数支持过滤查询结果：

* `gw_instance_id` (必填)：网关实例ID，用于指定要查询的API网关实例。

* `appid` (可选)：应用ID，用于精确过滤特定应用ID的消费者。

* `ids` (可选)：消费者ID列表，格式为`{gwInstanceId:appId}`，用于过滤指定ID的消费者。

* `name_regex` (可选)：应用名称正则表达式，用于通过正则表达式匹配应用名称来过滤消费者。

* `is_source_consumer` (可选)：是否为源消费者，true表示是源消费者，false表示不是源消费者。

## 属性说明
以下属性被导出：

* `id` (字符串)：消费者ID，格式为`{gwInstanceId:appId}`，是资源的唯一标识符。

* `access_key` (字符串)：访问密钥，当认证类型为CSB时返回。

* `app_code` (字符串)：应用代码，当认证类型为appgw-app时返回。

* `app_id` (字符串)：应用ID，API网关系统生成的唯一标识。

* `app_name` (字符串)：应用名称，创建消费者时指定的名称。

* `app_secret` (字符串)：应用密钥，用于API认证的密钥信息。

* `auth_type` (整数)：认证类型编码，5=API_KEY, 1=BASIC, 3=JWT, 2=OAuth2.0, 6=API_GATEWAY, 7=CSB_AUTH。

* `auth_type_name` (字符串)：认证类型名称，如"API_KEY"、"BASIC"、"JWT"等。

* `description` (字符串)：消费者描述信息。

* `enable` (布尔)：是否启用该消费者，true表示启用，false表示禁用。

* `expire_time` (整数)：过期时间（毫秒），表示令牌的有效期。

* `groups` (列表)：所属分组列表，消费者所属的API分组。

* `key` (字符串)：认证密钥，不同认证类型含义不同，如API Key认证中的key、Basic认证中的用户名等。

* `oauth2_payload` (列表)：OAuth2.0认证负载信息，包含以下子属性：
  * `authorization_code` (布尔)：是否启用授权码模式
  * `client_credentials` (布尔)：是否启用客户端凭证模式
  * `implicit_grant` (布尔)：是否启用隐式授权模式
  * `password_grant` (布尔)：是否启用密码凭证模式
  * `token_expiration` (整数)：令牌过期时间（秒）
  * `refresh_token_expiration` (整数)：刷新令牌过期时间（天）
  * `pkce` (布尔)：是否启用PKCE（Proof Key for Code Exchange）
  * `scopes` (字符串)：作用域列表，逗号分隔
  * `client_id` (字符串)：OAuth2.0客户端ID
  * `client_secret` (字符串)：OAuth2.0客户端密钥
  * `redirect_uris` (字符串)：重定向URI列表，逗号分隔

* `password` (字符串)：密码，当认证类型为Basic时返回。

* `request_header` (字符串)：请求头，OAuth2.0认证中使用的请求头名称。

* `secret_key` (字符串)：密钥，当认证类型为CSB时返回。

* `token` (字符串)：令牌，不同认证类型生成的令牌信息。

* `use_white_list` (布尔)：是否启用白名单，true表示启用，false表示禁用。