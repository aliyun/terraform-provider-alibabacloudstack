---
subcategory: "ApiGateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_apigateway_apigroup"
sidebar_current: "docs-Alibabacloudstack-apigateway-apigroup"
description: |- 
  Provides a apigateway Apigroup resource.
---

# alibabacloudstack_apigateway_apigroup
-> **NOTE:** Alias name has: `alibabacloudstack_api_gateway_group`

Provides an API Gateway Group resource. To create an API, you must first create a group which is a basic attribute of the API.

For information about Api Gateway Group and how to use it, see [Create An Api Group](https://help.aliyun.com/apsara/enterprise/v_3_14_0_20210519/apigateway/apsara-developer-guide/creates-an-api-group--1.html?spm=a2c4g.14484438.10001.139)

-> **NOTE:** Terraform will auto build api group while it uses `alibabacloudstack_api_gateway_group` to build api group.

## Example Usage

Basic Usage

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

## Argument Reference

The following arguments are supported:

* `api_group_name` - (Required) The name of the API Gateway Group. It must be unique within the AlibabaCloudStack account.
* `description` - (Required) Description of the API Gateway Group, which cannot exceed 180 characters.
* `custom_trace_config` - (Optional) Custom trace configuration for the API Gateway Group. It should be a JSON string or a JSON object that can be serialized into a string. The structure includes:
  * `parameterLocation` - (Required) Specifies where the trace ID is located, such as in the query string (`QUERY`) or headers (`HEADER`).
  * `parameterName` - (Required) The name of the trace ID parameter.
* `compatible_flags` - (Optional) Compatibility flags for the API Gateway Group. For example, `supportSSE` enables Server-Sent Events support.
* `user_log_config` - (Optional) User log configuration for the API Gateway Group. It should be a JSON string or a JSON object that can be serialized into a string. The structure includes:
  * `requestBody` - (Optional) Whether to log the request body.
  * `responseBody` - (Optional) Whether to log the response body.
  * `queryString` - (Optional) Query string parameters to log.
  * `requestHeaders` - (Optional) Request headers to log.
  * `responseHeaders` - (Optional) Response headers to log.
  * `jwtClaims` - (Optional) JWT claims to log.
* `passthrough_headers` - (Optional) Headers that should be passed through to the backend service. This is a comma-separated list of header names.
* `name` - (Optional) Deprecated field representing the name of the API Gateway Group. Use `api_group_name` instead.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the API Gateway Group.
* `sub_domain` - Second-level domain name automatically assigned to the API Gateway Group. This domain is used for API call testing.
* `vpc_domain` - Second-level VPC domain name automatically assigned to the API Gateway Group. This domain is used for internal network communication.
* `api_group_name` - The name of the API Gateway Group.
* `sub_domain` - Second-level domain name bound to the API Gateway Group, which is used for API call testing.
* `vpc_domain` - Inner sub domain.
* `name` - Deprecated field representing the name of the API Gateway Group. Use `api_group_name` instead.

## Import

API Gateway Group can be imported using the id, e.g.

```shell
$ terraform import alibabacloudstack_apigateway_apigroup.example "ab2351f2ce904edaa8d92a0510832b91"
```