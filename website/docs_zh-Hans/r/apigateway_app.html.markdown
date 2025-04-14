---
subcategory: "ApiGateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_apigateway_app"
sidebar_current: "docs-Alibabacloudstack-apigateway-app"
description: |- 
  编排API网关下的应用
---

# alibabacloudstack_apigateway_app
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_api_gateway_app_attachment`

使用Provider配置的凭证在指定的资源集下编排API网关下的应用。

## 示例用法

以下是一个完整的示例，展示了如何创建一个API网关应用并将其绑定到特定的API。

```hcl
variable "name" {
  default = "tf_testAccApp_8232031"
}

resource "alibabacloudstack_api_gateway_group" "default" {
  name        = var.name
  description = "tf_testAccApiGroup Description"
}

resource "alibabacloudstack_api_gateway_api" "default" {
  name        = var.name
  group_id    = alibabacloudstack_api_gateway_group.default.id
  description = "description"
  auth_type   = "APP"

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
    timeout   = 22
    aone_name = "cloudapi-openapi"
  }

  request_parameters {
    name         = "aa"
    type         = "STRING"
    required     = "OPTIONAL"
    in           = "QUERY"
    in_service   = "QUERY"
    name_service = "testparams"
  }
}

resource "alibabacloudstack_api_gateway_app" "default" {
  name        = var.name
  description = "tf_testAccApiAPP Description"
}

resource "alibabacloudstack_api_gateway_app_attachment" "default" {
  app_id     = alibabacloudstack_api_gateway_app.default.id
  api_id     = alibabacloudstack_api_gateway_api.default.api_id
  group_id   = alibabacloudstack_api_gateway_group.default.id
  stage_name = "PRE"
}
```

## 参数说明

支持以下参数：

* `app_id` - (必填，变更时重建) 需要授权的应用的ID。此字段是必填项，并且创建后无法修改(`ForceNew`)。
* `group_id` - (必填，变更时重建) API所属的API组的ID。此字段是必填项，并且创建后无法修改(`ForceNew`)。
* `api_id` - (必填，变更时重建) 应用需要访问的API的ID。此字段是必填项，并且创建后无法修改(`ForceNew`)。
* `stage_name` - (必填，变更时重建) 应用需要访问的阶段。有效值包括 `"RELEASE"`，`"TEST"` 和 `"PRE"`。此字段是必填项，并且创建后无法修改(`ForceNew`)。
* `name` - (必填) 资源名称。
* `description` - (可选) 资源描述。

### 详细字段描述

#### `app_id`
这是需要授权以访问指定API的应用程序的唯一标识符。此字段是必填项，并且创建后无法修改(`ForceNew`)。

#### `group_id`
这是指定API所属的API组的唯一标识符。此字段是必填项，并且创建后无法修改(`ForceNew`)。

#### `api_id`
这是应用程序被授权访问的API的唯一标识符。此字段是必填项，并且创建后无法修改(`ForceNew`)。

#### `stage_name`
这指定了应用程序被授权访问API的阶段。有效值为 `"RELEASE"`，`"TEST"` 和 `"PRE"`。此字段是必填项，并且创建后无法修改(`ForceNew`)。

## 属性说明

除了上述所有参数外，还导出了以下属性：

* `id` - API网关的应用绑定ID，格式为 `<group_id>:<api_id>:<app_id>:<stage_name>`。