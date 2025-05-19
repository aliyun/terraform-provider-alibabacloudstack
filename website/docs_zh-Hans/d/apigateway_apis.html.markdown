---
subcategory: "ApiGateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_apigateway_apis"
sidebar_current: "docs-Alibabacloudstack-datasource-apigateway-apis"
description: |- 
  查询接口网关中的接口
---

# alibabacloudstack_apigateway_apis
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_api_gateway_apis`

根据指定过滤条件列出当前凭证权限可以访问的接口网关下的接口列表。

## 示例用法

```hcl
variable "name" {
  default = "tf_testAccApisDataSource_2707186"
}

variable "apigateway_group_description_test" {
  default = "tf_testAcc_api group description"
}

resource "alibabacloudstack_api_gateway_group" "default" {
  name        = var.name
  description = var.apigateway_group_description_test
}

resource "alibabacloudstack_api_gateway_api" "default" {
  name         = var.name
  group_id     = alibabacloudstack_api_gateway_group.default.id
  description  = "tf_testAcc_api description"
  auth_type    = "APP"

  request_config {
    protocol = "HTTP"
    method   = "GET"
    path     = "/test/path"
    mode     = "MAPPING"
  }

  service_type = "HTTP"

  http_service_config {
    address   = "http://apigateway-backend.alicloudapi.com:8080"
    method    = "GET"
    path      = "/web/cloudapi"
    timeout   = 20
    aone_name = "cloudapi-openapi"
  }

  request_parameters {
    name          = "testparam"
    type          = "STRING"
    required      = "OPTIONAL"
    in            = "QUERY"
    in_service    = "QUERY"
    name_service  = "testparams"
  }
}

data "alibabacloudstack_apigateway_apis" "default" {
  ids = [
    alibabacloudstack_api_gateway_api.default.api_id
  ]

  output_file = "output_api_list.txt"
}

output "first_api_id" {
  value = data.alibabacloudstack_apigateway_apis.default.apis[0].id
}
```

## 参数说明

以下参数是支持的：

* `group_id` - (可选) API 所属组的 ID。  
* `api_id` - (可选，已废弃) 特定 API 的 ID。此字段从版本 1.52.2 开始已被废弃，并由 `ids` 字段取代。  
* `ids` - (可选) 用于过滤结果的 API ID 列表。  
* `name_regex` - (可选) 用于按名称筛选 API Gateway API 的正则表达式字符串。  

> **注意**：`api_id` 字段已废弃，请使用 `ids` 字段替代。

## 属性说明

除了上述参数外，还导出以下属性：

* `names` - 匹配条件的 API 名称列表。  
* `apis` - 匹配条件的 API 列表。每个元素包含以下属性：
  * `id` - API 的唯一标识符。  
  * `name` - API 的名称。  
  * `region_id` - API 所在的区域 ID。  
  * `group_id` - API 所属组的 ID。  
  * `group_name` - API 所属组的名称。  
  * `description` - API 的描述信息。  
