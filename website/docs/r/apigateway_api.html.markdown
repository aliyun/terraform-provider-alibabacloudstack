---
subcategory: "ApiGateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_apigateway_api"
sidebar_current: "docs-Alibabacloudstack-apigateway-api"
description: |- 
  Provides a Alibabacloudstack Api Gateway Api Resource.
---

# alibabacloudstack_apigateway_api
-> **NOTE:** Alias name has: `alibabacloudstack_api_gateway_api`

Provides an API resource for Alibaba Cloud API Gateway. When you create an API, you must enter the basic information about the API and define the API request information, backend service configuration, and response information.

For information about Api Gateway Api and how to use it, see [Create an API](https://help.aliyun.com/apsara/enterprise/v_3_14_0_20210519/apigateway/apsara-developer-guide/create-an-api-dev1.html?spm=a2c4g.14484438.10001.154)

-> **NOTE:** Terraform will auto build API while it uses `alibabacloudstack_api_gateway_api` to build API.

## Example Usage

Basic Usage

```hcl
variable "name" {
  default = "tf_testAccApiGatewayApi_7875290"
}

variable "apigateway_group_description_test" {
  default = "tf_testAcc_api group description"
}

resource "alibabacloudstack_api_gateway_group" "default" {
  name        = "${var.name}"
  description = "${var.apigateway_group_description_test}"
}

resource "alibabacloudstack_api_gateway_api" "default" {
  name           = "${alibabacloudstack_api_gateway_group.default.name}"
  group_id       = "${alibabacloudstack_api_gateway_group.default.id}"
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

## Argument Reference

The following arguments are supported:

* `group_id` - (Required, ForceNew) The ID of the API Gateway group that this API belongs to.
* `name` - (Optional) The name of the API. If not specified, the name will be set to the value of `api_name`.
* `api_name` - (Optional) The name of the API. If not specified, the name will be set to the value of `name`.
* `description` - (Required) The description of the API.
* `auth_type` - (Required) The authorization type of the API. Valid values include `APP`, `ANONYMOUS`.
* `request_config` - (Required, List) Configuration for the API request.

  * `protocol` - (Required) The protocol used by the API. Valid values are `HTTP`, `HTTPS`, or `HTTP,HTTPS`.
  * `method` - (Required) The HTTP method used by the API. Valid values include `GET`, `POST`, `PUT`, etc.
  * `path` - (Required) The request path of the API.
  * `mode` - (Required) The mode of parameter mapping between request parameters and service parameters. Valid values are `MAPPING` and `PASSTHROUGH`.
  * `body_format` - (Optional) The body format of the API. Valid values are `STREAM` and `FORM`.

* `service_type` - (Required) The type of backend service. Valid values include `HTTP`, `VPC`, `MOCK`, `FunctionCompute`.
* `stage_names` - (Optional, List) Stages that the API needs to be deployed to. Valid values include `RELEASE`, `PRE`, `TEST`.
* `force_nonce_check` - (Optional, Boolean) Whether to enable nonce check to prevent API replay attacks. Default value is `false`.
* `fc_service_config` - (Optional, List) Configuration for Function Compute backend service when `service_type` is set to `FunctionCompute`.

  * `region` - (Required) The region where the Function Compute service is located.
  * `function_name` - (Required) The name of the Function Compute function.
  * `service_name` - (Required) The name of the Function Compute service.
  * `arn_role` - (Optional) The ARN role attached to the Function Compute service. This governs both who/what can invoke your Function, as well as what resources your Function has access to. See [User Permissions](https://www.alibabacloud.com/help/doc-detail/52885.htm) for more details.
  * `timeout` - (Required) Backend service time-out time; unit: millisecond.

* `request_parameters` - (Optional, List) Configuration for the request parameters of the API.

  * `name` - (Required) The name of the request parameter.
  * `type` - (Required) The type of the parameter. Valid values include `STRING`, `INT`, `BOOLEAN`, `LONG`, `FLOAT`, `DOUBLE`.
  * `required` - (Required) Whether the parameter is required. Valid values are `REQUIRED` and `OPTIONAL`.
  * `in` - (Required) The location of the parameter in the request. Valid values are `BODY`, `HEAD`, `QUERY`, `PATH`.
  * `in_service` - (Required) The location of the parameter in the backend service. Valid values are `BODY`, `HEAD`, `QUERY`, `PATH`.
  * `name_service` - (Required) The name of the parameter in the backend service.
  * `description` - (Optional) The description of the parameter.
  * `default_value` - (Optional) The default value of the parameter.

* `system_parameters` - (Optional, List) Configuration for system parameters of the API.

  * `name` - (Required) The name of the system parameter.
  * `in` - (Required) The location of the parameter in the request. Valid values are `HEAD`, `QUERY`.
  * `name_service` - (Required) The name of the parameter in the backend service.

* `http_vpc_service_config` - (Optional, List) Configuration for HTTP-VPC backend service when `service_type` is set to `HTTP-VPC`.

  * `name` - (Required) The name of the VPC instance.
  * `path` - (Required) The path of the backend service.
  * `method` - (Required) The HTTP method of the backend service.
  * `timeout` - (Required) Backend service time-out time; unit: millisecond.
  * `aone_name` - (Optional) The name of Aone.

* `http_service_config` - (Optional, List) Configuration for HTTP backend service when `service_type` is set to `HTTP`.

  * `address` - (Required) The address of the backend service.
  * `path` - (Required) The path of the backend service.
  * `method` - (Required) The HTTP method of the backend service.
  * `timeout` - (Required) Backend service time-out time; unit: millisecond.
  * `aone_name` - (Optional) The name of Aone.

* `constant_parameters` - (Optional, List) Configuration for constant parameters of the API.

  * `name` - (Required) The name of the constant parameter.
  * `in` - (Required) The location of the parameter in the request. Valid values are `HEAD`, `QUERY`.
  * `value` - (Required) The value of the constant parameter.
  * `description` - (Optional) The description of the constant parameter.

* `mock_service_config` - (Optional, List) Configuration for mock backend service when `service_type` is set to `MOCK`.

  * `result` - (Required) The result returned by the mock service.
  * `aone_name` - (Optional) The name of Aone.
  
* `body_format` - (Optional) Specifies the body format of the API request. Valid values are `STREAM` and `FORM`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the API resource in the API Gateway.
* `api_id` - The unique identifier of the API within the API Gateway.
* `name` - The name of the API.
* `api_name` - The name of the API.
* `force_nonce_check` - Whether nonce check is enabled for the API.
* `description` - The description of the API.