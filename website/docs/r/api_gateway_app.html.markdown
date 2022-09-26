---
subcategory: "API Gateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_api_gateway_app"
sidebar_current: "docs-alibabacloudstack-resource-api-gateway-app"
description: |-
  Provides a Alibabacloudstack Api Gateway App Resource.
---

# alibabacloudstack_api_gateway_app

Provides an app resource.It must create an app before calling a third-party API because the app is the identity used to call the third-party API.

For information about Api Gateway App and how to use it, see [Create An APP](https://help.aliyun.com/apsara/enterprise/v_3_14_0_20210519/apigateway/apsara-developer-guide/create-an-application-1.html?spm=a2c4g.14484438.10001.177)

-> **NOTE:** Terraform will auto build api app while it uses `alibabacloudstack_api_gateway_app` to build api app.

## Example Usage

Basic Usage

```
variable "name" {
  default = "tf_testAccApp_2887078"
}
variable "description" {
  default = "tf_testAcc api gateway description"
}
resource "alibabacloudstack_api_gateway_app" "default" {
  name = "${var.name}"
  description = "${var.description}"
}

```
## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the app. 
* `description` - (Optional) The description of the app. Defaults to null.
* `tags` - (Optional, Available in v1.55.3+) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the app of api gateway.

## Import

Api gateway app can be imported using the id, e.g.

```
$ terraform import alicloud_api_gateway_app.example "7379660"
```
