---
subcategory: "ApiGateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_api_gateway_app_attachment"
sidebar_current: "docs-alibabacloudstack-resource-api-gateway-app-attachment"
description: |- 
  Provides a Alibabacloudstack Api Gateway App Attachment Resource.
---

# alibabacloudstack_api_gateway_app_attachment
-> **NOTE:** Alias name has: `alibabacloudstack_apigateway_app`

Provides an app attachment resource. It is used for authorizing a specific API to an app accessing. 

For information about Api Gateway App attachment and how to use it, see [Add specified API access authorities](https://help.aliyun.com/apsara/enterprise/v_3_14_0_20210519/apigateway/apsara-developer-guide/authorize-the-app-to-use-multiple-apis-1.html?spm=a2c4g.14484438.10001.187)

-> **NOTE:** Terraform will auto build app attachment while it uses `alibabacloudstack_api_gateway_app_attachment` to build.

## Example Usage

Basic Usage

```hcl
variable "name" {
  default = "tf_testAccApp_3907459"
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

## Argument Reference

The following arguments are supported:

* `app_id` - (Required, ForceNew) The ID of the app that applies to the authorization.
* `api_id` - (Required, ForceNew) The ID of the API that the app applies to access.
* `group_id` - (Required, ForceNew) The ID of the API group that the API belongs to.
* `stage_name` - (Required, ForceNew) The stage that the app applies to access. Valid values include `"RELEASE"`, `"TEST"`, `"PRE"`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the app attachment of API Gateway, formatted as `<group_id>:<api_id>:<app_id>:<stage_name>`.

### Detailed Field Descriptions

#### `app_id`

This is the unique identifier of the application that requires authorization to access the specified API. This field is mandatory and cannot be modified after creation (`ForceNew`).

#### `api_id`

This is the unique identifier of the API that the application is authorized to access. This field is mandatory and cannot be modified after creation (`ForceNew`).

#### `group_id`

This is the unique identifier of the API group to which the specified API belongs. This field is mandatory and cannot be modified after creation (`ForceNew`).

#### `stage_name`

This specifies the stage where the application is authorized to access the API. Valid values are `"RELEASE"`, `"TEST"`, and `"PRE"`. This field is mandatory and cannot be modified after creation (`ForceNew`).