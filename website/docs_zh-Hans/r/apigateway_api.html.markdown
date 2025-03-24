---
subcategory: "ApiGateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_apigateway_api"
sidebar_current: "docs-Alibabacloudstack-apigateway-api"
description: |- 
  编排API网关下的API
---

# alibabacloudstack_apigateway_api
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_api_gateway_api`

使用Provider配置的凭证在指定的资源集下编排API网关下的API。

## 示例用法

以下是一个完整的示例，展示了如何使用 `alibabacloudstack_api_gateway_api` 资源来创建一个 API 网关中的 API。

```hcl
variable "name" {
  default = "tf_testAccApiGatewayApi_7875290"
}

variable "apigateway_group_description_test" {
  default = "tf_testAcc_api group description"
}

resource "alibabacloudstack_api_gateway_group" "default" {
  name        = var.name
  description = var.apigateway_group_description_test
}

resource "alibabacloudstack_api_gateway_api" "default" {
  name           = alibabacloudstack_api_gateway_group.default.name
  group_id       = alibabacloudstack_api_gateway_group.default.id
  description    = "tf_testAcc_api description"
  auth_type      = "APP"
  force_nonce_check = true

  request_config {
    protocol = "HTTP"
    method   = "GET"
    path     = "/test/path"
    mode     = "MAPPING"
  }

  service_type = "HTTP"

  http_service_config {
    address    = "http://apigateway-backend.alicloudapi.com:8080"
    method     = "GET"
    path       = "/web/cloudapi"
    timeout    = "20"
    aone_name  = "cloudapi-openapi"
  }

  request_parameters {
    name         = "testparam"
    type         = "STRING"
    required     = "OPTIONAL"
    in           = "QUERY"
    in_service   = "QUERY"
    name_service = "testparams"
  }
}
```

## 参数参考

支持以下参数：

* `group_id` - (必填, 变更时重建) - 所属分组的 ID。
* `name` - (选填) - API 的名称。如果未指定，则默认使用 `api_name` 的值。
* `api_name` - (选填) - API 的名称。如果未指定，则默认使用 `name` 的值。
* `description` - (必填) - API 的描述信息。
* `auth_type` - (必填) - API 的授权类型。有效值包括 `APP` 和 `ANONYMOUS`。
* `request_config` - (必填) - API 请求的配置。

  * `protocol` - (必填) - API 使用的协议。有效值为 `HTTP`、`HTTPS` 或 `HTTP,HTTPS`。
  * `method` - (必填) - API 使用的 HTTP 方法。有效值包括 `GET`、`POST`、`PUT` 等。
  * `path` - (必填) - API 的请求路径。
  * `mode` - (必填) - 请求参数和服务参数之间参数映射的模式。有效值为 `MAPPING` 和 `PASSTHROUGH`。
  * `body_format` - (选填) - API 的主体格式。有效值为 `STREAM` 和 `FORM`。

* `service_type` - (必填) - 后端服务的类型。有效值包括 `HTTP`、`VPC`、`MOCK`、`FunctionCompute`。
* `stage_names` - (选填) - 需要部署 API 的阶段。有效值包括 `RELEASE`、`PRE`、`TEST`。
* `force_nonce_check` - (选填) - 是否启用 nonce 检查以防止 API 回放攻击。默认值为 `false`。
* `fc_service_config` - (选填) - 当 `service_type` 设置为 `FunctionCompute` 时，Function Compute 后端服务的配置。

  * `region` - (必填) - Function Compute 服务所在的区域。
  * `function_name` - (必填) - Function Compute 函数的名称。
  * `service_name` - (必填) - Function Compute 服务的名称。
  * `arn_role` - (选填) - 附加到 Function Compute 服务的 ARN 角色。
  * `timeout` - (必填) - 后端服务超时时间(单位：毫秒)。

* `request_parameters` - (选填) - API 请求参数的配置。

  * `name` - (必填) - 请求参数的名称。
  * `type` - (必填) - 参数的类型。有效值包括 `STRING`、`INT`、`BOOLEAN`、`LONG`、`FLOAT`、`DOUBLE`。
  * `required` - (必填) - 参数是否必填。有效值为 `REQUIRED` 和 `OPTIONAL`。
  * `in` - (必填) - 参数在请求中的位置。有效值为 `BODY`、`HEAD`、`QUERY`、`PATH`。
  * `in_service` - (必填) - 参数在后端服务中的位置。有效值为 `BODY`、`HEAD`、`QUERY`、`PATH`。
  * `name_service` - (必填) - 参数在后端服务中的名称。
  * `description` - (选填) - 参数的描述。
  * `default_value` - (选填) - 参数的默认值。

* `system_parameters` - (选填) - API 系统参数的配置。

  * `name` - (必填) - 系统参数的名称。
  * `in` - (必填) - 参数在请求中的位置。有效值为 `HEAD`、`QUERY`。
  * `name_service` - (必填) - 参数在后端服务中的名称。

* `http_vpc_service_config` - (选填) - 当 `service_type` 设置为 `HTTP-VPC` 时，HTTP-VPC 后端服务的配置。

  * `name` - (必填) - VPC 实例的名称。
  * `path` - (必填) - 后端服务的路径。
  * `method` - (必填) - 后端服务的 HTTP 方法。
  * `timeout` - (必填) - 后端服务超时时间(单位：毫秒)。
  * `aone_name` - (选填) - Aone 的名称。

* `http_service_config` - (选填) - 当 `service_type` 设置为 `HTTP` 时，HTTP 后端服务的配置。

  * `address` - (必填) - 后端服务的地址。
  * `path` - (必填) - 后端服务的路径。
  * `method` - (必填) - 后端服务的 HTTP 方法。
  * `timeout` - (必填) - 后端服务超时时间(单位：毫秒)。
  * `aone_name` - (选填) - Aone 的名称。

* `constant_parameters` - (选填) - API 常量参数的配置。

  * `name` - (必填) - 常量参数的名称。
  * `in` - (必填) - 参数在请求中的位置。有效值为 `HEAD`、`QUERY`。
  * `value` - (必填) - 常量参数的值。
  * `description` - (选填) - 常量参数的描述。

* `mock_service_config` - (选填) - 当 `service_type` 设置为 `MOCK` 时，模拟后端服务的配置。

  * `result` - (必填) - 模拟服务返回的结果。
  * `aone_name` - (选填) - Aone 的名称。

## 属性参考

除了上述所有参数外，还导出了以下属性：

* `id` - API Gateway 中 API 资源的 ID。
* `api_id` - API Gateway 内 API 的唯一标识符。
* `name` - API 的名称。
* `api_name` - API 的名称。
* `force_nonce_check` - 是否为 API 启用了 nonce 检查。
