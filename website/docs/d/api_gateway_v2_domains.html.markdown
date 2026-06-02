---
subcategory: "API Gateway"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_api_gateway_v2_domains"
sidebar_current: "docs-Alibabacloudstack-datasource-api-gateway-v2-domains"
description: |-
  Provides a list of Api Gateway V2 Domains to the user.
---

# alibabacloudstack\_api\_gateway\_v2\_domains

This data source provides the Api Gateway V2 Domains of the current Alibaba Cloud user.

## Example Usage

```hcl
variable "name" {
  default = "tf-example"
}

data "alibabacloudstack_api_gateway_v2_domains" "example" {
  instance_id = "example-instance-id"
  domain      = "example.com"
  protocol    = "HTTPS"
}

output "first_api_gateway_v2_domain_id" {
  value = data.alibabacloudstack_api_gateway_v2_domains.example.domains.0.id
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) The ID of the API Gateway instance.
* `domain` - (Optional) The custom domain name for the API Gateway.
* `domain_regex` - (Optional) A regex string to filter results by domain name.
* `ids` - (Optional) A list of domain IDs.
* `protocol` - (Optional) The protocol used by the domain. Valid values: `HTTP`, `HTTPS`.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of domain IDs.
* `domains` - A list of Api Gateway V2 Domains. Each element contains the following attributes:
  * `id` - The ID of the domain. It formats as `<instance_id>:<domain_id>`.
  * `instance_id` - The ID of the API Gateway instance.
  * `domain` - The custom domain name for the API Gateway.
  * `domain_id` - The ID of the domain.
  * `protocol` - The protocol used by the domain.
  * `certificate_id` - The ID of the certificate.
  * `ca_certificate_id` - The ID of the CA certificate for client certificate authentication.
  * `client_auth` - Whether to enable client certificate authentication. Valid values: `0` (disabled), `1` (enabled).
  * `subject_dn` - The subject DN of the client certificate.
  * `issuer_dn` - The issuer DN of the client certificate.
