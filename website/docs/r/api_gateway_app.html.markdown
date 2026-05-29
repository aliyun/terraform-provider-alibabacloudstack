---
subcategory: "API Gateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_api_gateway_app"
sidebar_current: "docs-Alibabacloudstack-resource-api-gateway-app"
description: |-
  Provides a Alibabacloudstack Api Gateway App Resource.
---

# alibabacloudstack_api_gateway_app

Provides an API Gateway App resource. This resource is used to create and manage applications for API Gateway, which are used for API authentication and authorization.

For information about API Gateway App and how to use it, see [Application Management](https://help.aliyun.com/zh/api-gateway/traditional-api-gateway/developer-reference/api-cloudapi-2016-07-14-dir-applications/).

-> **NOTE:** Terraform will automatically build the app while it uses `alibabacloudstack_api_gateway_app` to build.

## Example Usage

Basic Usage

```hcl
variable "name" {
  default = "tf_testAccApp_example"
}

variable "description" {
  default = "tf_testAcc api gateway description"
}

resource "alibabacloudstack_api_gateway_app" "default" {
  name        = var.name
  description = var.description
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the API Gateway app. It must start with a letter or Chinese character, and can contain letters, digits, and underscores. The length must be between 4 and 26 characters.
* `description` - (Optional) The description of the app. The maximum length is 180 characters.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the app. It is the unique identifier assigned by API Gateway.

## Import

API Gateway App can be imported using the app ID, e.g.

```
$ terraform import alibabacloudstack_api_gateway_app.example 12345678
```
