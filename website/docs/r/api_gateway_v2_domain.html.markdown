---
subcategory: "API Gateway V2"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_api_gateway_v2_domain"
sidebar_current: "docs-Alibabacloudstack-resource-api-gateway-v2-domain"
description: |-
  Provides a AlibabacloudStack API Gateway V2 Domain resource.
---

# alibabacloudstack\_api\_gateway\_v2\_domain

Provides a API Gateway V2 Domain resource.

## Example Usage

Basic Usage

```hcl
variable "name" {
  default = "tf-example"
}

data "alibabacloudstack_api_gateway_v2_instance_types" "default" {
  sorted_by = "CPU"
}

resource "alibabacloudstack_api_gateway_v2_instance" "default" {
  instance_name  = var.name
  node_number    = "1"
  instance_class = data.alibabacloudstack_api_gateway_v2_instance_types.default.instance_types.0.id
  broker_engine_type = "SCG"
  deploy_mode    = "custom"
}

resource "alibabacloudstack_api_gateway_v2_certificate" "default" {
  certificate_name = var.name
  instance_id      = alibabacloudstack_api_gateway_v2_instance.default.id
  cert_type        = "0"
  certificates     = "-----BEGIN CERTIFICATE-----\n***********\n-----END CERTIFICATE-----"
  private_key      = "-----BEGIN RSA PRIVATE KEY-----\n***********\n-----END RSA PRIVATE KEY-----"
}

resource "alibabacloudstack_api_gateway_v2_domain" "default" {
  domain         = "${var.name}.com"
  instance_id    = alibabacloudstack_api_gateway_v2_instance.default.id
  protocol       = "HTTPS"
  certificate_id = alibabacloudstack_api_gateway_v2_certificate.default.certificate_id
  client_auth    = "0"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the API Gateway instance.
* `domain` - (Required, ForceNew) The custom domain name for the API Gateway.
* `protocol` - (Required) The protocol used by the domain. Valid values: `HTTP`, `HTTPS`.
* `certificate_id` - (Optional) The ID of the certificate. Required when protocol is `HTTPS`.
* `ca_certificate_id` - (Optional) The ID of the CA certificate for client certificate authentication. Required when `client_auth` is `1`.
* `client_auth` - (Optional) Whether to enable client certificate authentication. Valid values: `0` (disabled), `1` (enabled). Default: `0`.
* `subject_dn` - (Optional) The subject DN of the client certificate.
* `issuer_dn` - (Optional) The issuer DN of the client certificate.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the domain. It formats as `<instance_id>:<domain_id>`.
* `domain_id` - The ID of the domain.

## Import

API Gateway V2 Domain can be imported using the id, e.g.

```shell
$ terraform import alibabacloudstack_api_gateway_v2_domain.example <instance_id>:<domain_id>
```