---
subcategory: "API Gateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_api_gateway_app_attachment"
sidebar_current: "docs-alibabacloudstack-resource-api-gateway-app-attachment"
description: |-
  Provides a Alibabacloudstack Api Gateway App Attachment Resource.
---

# alibabacloudstack_api_gateway_app_attachment
-> **NOTE:** Alias name has: `alibabacloudstack_apigateway_app`

Provides an app attachment resource.It is used for authorizing a specific api to an app accessing. 

For information about Api Gateway App attachment and how to use it, see [Add specified API access authorities](https://help.aliyun.com/apsara/enterprise/v_3_14_0_20210519/apigateway/apsara-developer-guide/authorize-the-app-to-use-multiple-apis-1.html?spm=a2c4g.14484438.10001.187)

-> **NOTE:** Terraform will auto build app attachment while it uses `alibabacloudstack_api_gateway_app_attachment` to build.

## Example Usage

Basic Usage

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

## Argument Reference

The following arguments are supported:

* `app_id` - (Required，ForceNew) The app that apply to the authorization.
* `api_id` - (Required，ForceNew) The api_id that app apply to access.
* `group_id` - (Required，ForceNew) The group that the api belongs to.
* `stage_name` - (Required，ForceNew) Stage that the app apply to access.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the app attachment of api gateway., formatted as `<group_id>:<api_id>:<app_id>:<stage_name>`.
