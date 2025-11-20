---
subcategory: "API Gateway V2"
page_title: "AlibabacloudStack: alibabacloudstack_api_gateway_v2_signatures"
sidebar_current: "docs-alibabacloudstack-datasource-api-gateway-v2-signatures"
description: |-
  Provides a list of API Gateway V2 Signatures.
---

# alibabacloudstack\_api\_gateway\_v2\_signatures

This data source provides a list of API Gateway V2 Signatures in an Alibaba Cloud account according to the specified filters.

## Example Usage

```terraform
data "alibabacloudstack_api_gateway_v2_signatures" "example" {
  gw_instance_id = "api-gateway-v2-instance-id"
}

output "first_signature_id" {
  value = data.alibabacloudstack_api_gateway_v2_signatures.example.signatures.0.id
}
```

## Argument Reference

The following arguments are supported:

* `gw_instance_id` - (Required) The ID of the API Gateway V2 instance.
* `ids` - (Optional) A list of signature IDs.
* `name_regex` - (Optional) A regex string to filter results by signature name.
* `names` - (Optional) A list of signature names to filter signatures by name.
* `is_source_signature` - (Optional) Whether the signature is a source signature. Default value: `false`.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of signature IDs.
* `names` - A list of signature names.
* `signatures` - A list of signatures. Each element contains the following attributes:
  * `id` - The ID of the resource, formatted as `{gw_instance_id}:{sig_scheme_id}`.
  * `gw_instance_id` - The ID of the API Gateway V2 instance.
  * `sig_scheme_id` - The ID of the signature scheme.
  * `sig_scheme_name` - The name of the signature scheme.
  * `sig_alg` - The signature algorithm.
  * `secret_key` - The secret key.
  * `create_time` - The creation time of the signature.
  * `update_time` - The last update time of the signature.
  * `status` - The status of the signature.