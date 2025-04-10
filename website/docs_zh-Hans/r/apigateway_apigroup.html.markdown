---
subcategory: "ApiGateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_apigateway_apigroup"
sidebar_current: "docs-Alibabacloudstack-apigateway-apigroup"
description: |- 
  编排API网关下的API组
---

# alibabacloudstack_apigateway_apigroup
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_api_gateway_group`

使用Provider配置的凭证在指定的资源集下编排API网关下的API组。

## 示例用法

### 基础用法

```hcl
variable "name" {
  default = "tf_testAccGroup_4663614"
}

variable "description" {
  default = "tf_testAcc api gateway description"
}

resource "alibabacloudstack_apigateway_apigroup" "default" {
  api_group_name     = var.name
  description        = var.description
  custom_trace_config = jsonencode({
    parameterLocation = "QUERY"
    parameterName     = "traceId"
  })
  compatible_flags   = "supportSSE"
  user_log_config = jsonencode({
    requestBody       = false
    responseBody      = false
    queryString      = ""
    requestHeaders   = ""
    responseHeaders  = ""
    jwtClaims        = ""
  })
  passthrough_headers = "eagleeye-rpcid,x-b3-traceid"
}
```

## 参数说明

支持以下参数：

* `api_group_name` - (必填) API 网关组的名称。在同一 AlibabaCloudStack 账户内必须唯一。
* `description` - (必填) API 网关组的描述，不能超过 180 个字符。
* `custom_trace_config` - (可选) API 网关组的自定义跟踪配置。它应该是一个 JSON 字符串或可以序列化为字符串的 JSON 对象。结构包括：
  * `parameterLocation` - (必填) 指定跟踪 ID 的位置，例如在查询字符串中 (`QUERY`) 或标头中 (`HEADER`)。
  * `parameterName` - (必填) 跟踪 ID 参数的名称。
* `compatible_flags` - (可选) API 网关组的兼容标志。例如，`supportSSE` 启用 Server-Sent Events 支持。
* `user_log_config` - (可选) API 网关组的用户日志配置。它应该是一个 JSON 字符串或可以序列化为字符串的 JSON 对象。结构包括：
  * `requestBody` - (可选) 是否记录请求正文。
  * `responseBody` - (可选) 是否记录响应正文。
  * `queryString` - (可选) 要记录的查询字符串参数。
  * `requestHeaders` - (可选) 要记录的请求头。
  * `responseHeaders` - (可选) 要记录的响应头。
  * `jwtClaims` - (可选) 要记录的 JWT 声明。
* `passthrough_headers` - (可选) 应传递到后端服务的标头。这是一个以逗号分隔的标头名称列表。
* `name` - (可选) 已废弃字段，表示 API 网关组的名称。请改用 `api_group_name`。

## 属性说明

除了上述所有参数外，还导出了以下属性：

* `id` - API 网关组的 ID。
* `sub_domain` - 系统为 API 分组分配的二级域名，用于 API 调用测试。
* `vpc_domain` - 内网二级域名，用于内部网络通信。
* `api_group_name` - API 网关组的名称。
* `name` - 已废弃字段，表示 API 网关组的名称。请改用 `api_group_name`。

## 导入

API 网关组可以通过 ID 导入，例如：

```shell
$ terraform import alibabacloudstack_apigateway_apigroup.example "ab2351f2ce904edaa8d92a0510832b91"
``` 