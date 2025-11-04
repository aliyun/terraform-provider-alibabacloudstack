---
subcategory: "API Gateway V2"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_api_gateway_v2_certificates"
sidebar_current: "docs-alibabacloudstack-datasource-api-gateway-v2-certificates"
description: |-
  Provides a list of Api Gateway V2 Certificates to the user.
---

# alibabacloudstack\_api\_gateway\_v2\_certificates

This data source provides the Api Gateway V2 Certificates of the current Alibaba Cloud user.

## Example Usage

```hcl
data "alibabacloudstack_api_gateway_v2_certificates" "example" {
  instance_id = "example-instance-id"
}

output "first_certificate_id" {
  value = data.alibabacloudstack_api_gateway_v2_certificates.example.certificates.0.id
}
```

```hcl
data "alibabacloudstack_api_gateway_v2_certificates" "filtered" {
  instance_id  = "example-instance-id"
  cert_type    = "0"
  name_regex   = "^example"
}

output "certificate_names" {
  value = [for cert in data.alibabacloudstack_api_gateway_v2_certificates.filtered.certificates : cert.certificate_name]
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the API Gateway instance.
* `cert_type` - (Optional, ForceNew) The certificate type. Valid values: `0` (Server certificate), `1` (CA certificate).
* `name_regex` - (Optional, ForceNew) A regex string to filter results by certificate name.
* `sni` - (Optional, ForceNew) The SNI (Server Name Indication) of the certificate.
* `ids` - (Optional, ForceNew) A list of certificate IDs.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the data source.
* `ids` - A list of certificate IDs.
* `certificates` - A list of certificates. Each element contains the following attributes:
  * `id` - The resource ID. The value is `<instance_id>:<certificate_id>`.
  * `certificate_id` - The ID of the certificate.
  * `certificate_name` - The name of the certificate.
  * `cert_type` - The certificate type.
  * `expire_time` - The expiration time of the certificate.
  * `create_time` - The creation time of the certificate.
  * `update_time` - The last update time of the certificate.
  * `snis` - The SNI list of the certificate.