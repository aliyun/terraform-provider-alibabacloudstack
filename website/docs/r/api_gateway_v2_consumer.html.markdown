---
subcategory: "API Gateway V2"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_api_gateway_v2_consumer"
sidebar_current: "docs-Alibabacloudstack-resource-api-gateway-v2-consumer"
description: |-
  Manage API Gateway V2 consumers for configuring access credentials for applications with different authentication methods.
---

# alibabacloudstack_api_gateway_v2_consumer

Manages API Gateway V2 consumers for configuring access credentials for applications with different authentication methods.

## Example Usage

### Basic Usage

```hcl

variable "name" {
  default = "tf-testacc32955"
}

data "alibabacloudstack_api_gateway_v2_instance_types" "default" {
  sorted_by = "CPU"
}

resource "alibabacloudstack_api_gateway_v2_instance" "default" {
  instance_name      = var.name
  node_number        = "1"
  instance_class     = data.alibabacloudstack_api_gateway_v2_instance_types.default.instance_types[0].id
  broker_engine_type = "SCG"
  deploy_mode        = "custom"
}



// Basic authentication
resource "alibabacloudstack_api_gateway_v2_consumer" "default" {
  auth_type      = 1
  app_name       = "${var.name}"
  description    = "${var.name}"
  key            = "root"
  password       = "admin"
  gw_instance_id = "${alibabacloudstack_api_gateway_v2_instance.default.id}"
  groups         = ["test"]
}
```

### OAuth2.0 authentication
```hcl
resource "alibabacloudstack_api_gateway_v2_consumer" "default" {
  gw_instance_id = alibabacloudstack_api_gateway_v2_instance.default.id
  groups = [
    "test"
  ]
  auth_type   = "2"
  app_name    = var.name
  description = var.name
  oauth2_payload = {
    "authorization_code": true,
    "client_credentials": true,
    "implicit_grant": true,
    "password_grant": true,
    "token_expiration": 3600,
    "refresh_token_expiration": 30,
    "pkce": true,
    "scopes": "api_gateway",
    "client_id": "a94231617c1a94cd900",
  }
}
```

### JWT authentication
```hcl
resource "alibabacloudstack_api_gateway_v2_consumer" "default" {
  auth_type      = 3
  app_name       = "${var.name}"
  description    = "${var.name}"
  key            = "aaaaaaaaaaaaaaaaaa"
  expire_time    = 86400000
  app_secret     = "bbbbbbbbbbbbbbbbbbbbb"
  gw_instance_id = "${alibabacloudstack_api_gateway_v2_instance.default.id}"
  groups         = ["test"]
  payload = {
    issuer  = "http://127.0.0.1:8000/test"
    subject = "test.app"
  }
}
```

### API Key authentication
```hcl
resource "alibabacloudstack_api_gateway_v2_consumer" "default" {
  gw_instance_id = alibabacloudstack_api_gateway_v2_instance.default.id
  groups = [
    "test"
  ]
  auth_type   = "5"
  app_name    = var.name
  description = var.name
  key         = "a94231617c1a94cd900"
}
```

### API Gateway application authentication
```hcl
resource "alibabacloudstack_api_gateway_v2_consumer" "default" {
  auth_type      = 6
  app_name       = "${var.name}"
  description    = "${var.name}"
  app_secret     = "vvvvvvvvvvvvvvvvvvvvv"
	app_code       = "aaaaaaaaaaa"
	key            = "aaaaaaaaaaaaaaaaaa"
  gw_instance_id = "${alibabacloudstack_api_gateway_v2_instance.default.id}"
  groups         = ["test"]
}
```

### CSB authentication
```hcl
resource "alibabacloudstack_api_gateway_v2_consumer" "default" {
  auth_type      = 7
  app_name       = "${var.name}"
  description    = "${var.name}"
  key            = "aaaaaaaaaaaaaaaaaa"
  gw_instance_id = "${alibabacloudstack_api_gateway_v2_instance.default.id}"
  groups         = ["test"]
}
```

## Argument Reference

The following arguments are supported:

* `app_name` - (Required) The name of the application. It must be 1 to 128 characters in length.
* `auth_type` - (Required, ForceNew) The authentication type. Valid values:
  * `1`: Basic authentication
  * `2`: OAuth 2.0 authentication
  * `3`: JWT authentication
  * `5`: API Key authentication
  * `6`: API Gateway application authentication
  * `7`: CSB authentication
* `gw_instance_id` - (Required, ForceNew) The ID of the gateway instance. The format is `i-xxx`.
* `app_code` - (Optional) The application code for API Gateway authentication.
* `app_secret` - (Optional, Computed) The application secret, used for authentication methods such as API Key and JWT.
* `cascade_link_ids` - (Optional, ForceNew) The list of cascade link IDs. Set this attribute to create the source consumer.
* `description` - (Optional) The description of the application. It must be 1 to 256 characters in length.
* `expire_time` - (Optional, Computed) The token expiration time in milliseconds. The default value varies depending on the authentication method.
* `groups` - (Optional) The list of groups to which the application belongs, used for permission control.
* `key` - (Optional) The key for API Key authentication, or the username for Basic authentication.
* `oauth2_payload` - (Optional) The configuration information for OAuth 2.0 authentication, including the following sub-parameters:
  * `authorization_code` - (Optional) Whether to enable the authorization code grant type.
  * `client_credentials` - (Optional) Whether to enable the client credentials grant type.
  * `implicit_grant` - (Optional) Whether to enable the implicit grant type.
  * `password_grant` - (Optional) Whether to enable the password grant type.
  * `token_expiration` - (Optional) The access token expiration time in seconds.
  * `refresh_token_expiration` - (Optional) The refresh token expiration time in days.
  * `pkce` - (Optional) Whether to enable the PKCE extension.
  * `scopes` - (Optional) The authorization scopes, separated by commas.
  * `client_id` - (Optional) The client ID.
  * `client_secret` - (Optional) The client secret.
  * `redirect_uris` - (Optional) The redirect URIs, separated by commas.
* `password` - (Optional) The password for Basic authentication.
* `payload` - (Optional) The payload information for JWT authentication, in key-value pair format (map of strings).

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in the format of `{prefix}:{gwInstanceId}:{appId}`, where `prefix` is either `app` or `sourceApp`.
* `access_key` - The access key for CSB authentication.
* `app_id` - The application ID, which is a unique identifier generated by the system.
* `auth_type_name` - The authentication type name, such as "API_KEY", "BASIC", etc.
* `enable` - Whether the application is enabled. `true` indicates enabled, and `false` indicates disabled.
* `request_header` - The request header information for OAuth 2.0 authentication.
* `secret_key` - The secret key for CSB authentication.
* `token` - The generated access token. The token format varies depending on the authentication method.
* `use_white_list` - Whether the whitelist is enabled. `true` indicates enabled, and `false` indicates disabled.

## Import

API Gateway V2 Consumer can be imported using the `{prefix}:{gwInstanceId}:{appId}`, e.g.

```
$ terraform import alibabacloudstack_api_gateway_v2_consumer.example app:i-xxxxxxxxxxxx:xxxxxxxxxx
```