---
subcategory: "API Gateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_api_gateway_app_attachment"
sidebar_current: "docs-alibabacloudstack-resource-api-gateway-app-attachment"
description: |-
  Provides a Alibabacloudstack Api Gateway App Attachment Resource.
---

# alibabacloudstack_api_gateway_app_attachment
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_apigateway_app`

提供一个API网关应用绑定资源。它用于授权特定的API给应用访问。

关于API网关应用绑定及其使用方法，参见 [添加指定API访问权限](https://help.aliyun.com/apsara/enterprise/v_3_14_0_20210519/apigateway/apsara-developer-guide/authorize-the-app-to-use-multiple-apis-1.html?spm=a2c4g.14484438.10001.187)

-> **注意:** Terraform在使用`alibabacloudstack_api_gateway_app_attachment`构建时会自动创建应用绑定。

## 示例用法

### 基础用法

```

variable "name" {
  default = "tf_testAccApp_3907459"
}
resource "alibabacloudstack_api_gateway_group" "default" {
  name        = "${var.name}"
  description = "tf_testAccApiGroup Description"
}
resource "alibabacloudstack_api_gateway_api" "default" {
  name        = "${var.name}"
  group_id    = "${alibabacloudstack_api_gateway_group.default.id}"
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
  name        = "${var.name}"
  description = "tf_testAccApiAPP Description"
}

resource "alibabacloudstack_api_gateway_app_attachment" "default" {
  app_id = "${alibabacloudstack_api_gateway_app.default.id}"
  api_id = "${alibabacloudstack_api_gateway_api.default.api_id}"
  group_id = "${alibabacloudstack_api_gateway_group.default.id}"
  stage_name = "PRE"
}
```

## 参数说明

支持以下参数：

* `app_id` - (必填，变更时重建) 需要授权的应用ID。
* `api_id` - (必填，变更时重建) 授权应用访问的目标API ID。
* `group_id` - (必填，变更时重建) 目标API所属的分组ID。
* `stage_name` - (必填，变更时重建) 应用申请访问的阶段名称。通常包括开发环境（`TEST`）、预发布环境（`PRE`）和生产环境（`RELEASE`）。

## 属性说明

导出以下属性：

* `id` - API网关应用绑定的唯一标识符，格式为 `<group_id>:<api_id>:<app_id>:<stage_name>`。此ID可用于后续资源的引用或查询。