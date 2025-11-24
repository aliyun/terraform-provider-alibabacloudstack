---
subcategory: "API Gateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_api_gateway_v2_signature"
sidebar_current: "docs-Alibabacloudstack-api-gateway-v2-signature"
description: |-
  Manages signature schemes for API Gateway v2 version, used to configure the signature verification mechanism for API requests.
---

# alibabacloudstack_api_gateway_v2_signature

Manages signature schemes for API Gateway v2 version, used to configure the signature verification mechanism for API requests.

## Example Usage

### Basic Usage

```hcl

variable "signature_name" {
  default = "test"
}

variable "signature_algorithm" {
  default = "HmacSM3"
}

resource "alibabacloudstack_api_gateway_v2_instance" "default" {
  instance_name      = "testtf-apigw-instance"
  node_number        = "1"
  instance_class     = "mini"
  broker_engine_type = "SCG"
  deploy_mode        = "custom"
}

data "alibabacloudstack_api_gateway_v2_instances" "default" {
  name_regex = alibabacloudstack_api_gateway_v2_instance.default.instance_name
}

resource "alibabacloudstack_api_gateway_v2_signature" "default" {
  sig_scheme_name = var.signature_name
  sig_alg         = var.signature_algorithm
  gw_instance_id  = alibabacloudstack_api_gateway_v2_instance.default.id
}

```

## Argument Reference

The following arguments are supported:

* `gw_instance_id` - (Required, ForceNew) The ID of the API Gateway instance. Specifies the gateway instance to which the signature scheme belongs.
* `sig_alg` - (Required, ForceNew) The signature algorithm. Valid values: `HmacSHA256`, `HmacSHA1`, `HmacSM3`.
* `sig_scheme_name` - (Required) The name of the signature scheme. A unique name that identifies the signature scheme. The length limit is determined by the API Gateway service.
* `status` - (Optional) The status of the signature scheme. Valid values: `0` (disabled), `1` (enabled). Default is empty, which means to keep the current status.
* `cascade_link_ids` - (Optional) The list of cascade link IDs, used for CSB authentication. Set this attribute to create the source signature.


## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in the format `{gwInstanceId}:{sigSchemeId}`.
* `create_time` - The creation time of the signature scheme in the format `YYYY-MM-DD HH:mm:ss`.
* `secret_key` - The secret key of the signature scheme, used to generate and verify signatures.
* `sig_scheme_id` - The unique identifier ID of the signature scheme.
* `update_time` - The last update time of the signature scheme in the format `YYYY-MM-DD HH:mm:ss`.