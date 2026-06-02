---
subcategory: "API Gateway"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_api_gateway_v2_certificate"
description: |-
  Provides a AlibabacloudStack API Gateway V2 Certificate resource.
---

# alibabacloudstack\_api\_gateway\_v2\_certificate

Provides a API Gateway V2 Certificate resource.


## Example Usage

Basic Usage

```hcl
resource "alibabacloudstack_api_gateway_v2_certificate" "example" {
  cert_type        = "0"
  instance_id      = "instance-id"
  certificates     = "-----BEGIN CERTIFICATE-----\n**********\n-----END CERTIFICATE-----"
  private_key      = "-----BEGIN PRIVATE KEY-----\n**********\n-----END PRIVATE KEY-----"
  certificate_name = "example-cert"
}
```

## Argument Reference

The following arguments are supported:

* `cert_type` - (Required, ForceNew) The certificate type. Valid values: `0` (Server certificate), `1` (CA certificate).
* `instance_id` - (Required, ForceNew) The ID of the API Gateway instance.
* `certificates` - (Required) The certificate content. When `cert_type` is `0`, this parameter represents the server certificate. When `cert_type` is `1`, this parameter represents the CA certificate.
* `private_key` - (Optional) The private key. This parameter is required when `cert_type` is `0`.
* `certificate_name` - (Required) The name of the certificate.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID. The value is the certificate ID.
* `certificate_id` - The ID of the certificate.
* `expire_time` - The expiration time of the certificate.
* `create_time` - The creation time of the certificate.
* `update_time` - The last update time of the certificate.

## Import

API Gateway V2 Certificate can be imported using the id, e.g.

```shell
$ terraform import alibabacloudstack_api_gateway_v2_certificate.example <instance_id:certificate_id>
```
