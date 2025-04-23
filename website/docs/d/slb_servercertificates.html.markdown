---
subcategory: "Server Load Balancer (SLB)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_slb_servercertificates"
sidebar_current: "docs-Alibabacloudstack-datasource-slb-servercertificates"
description: |- 
  Provides a list of slb servercertificates owned by an alibabacloudstack account.
---

# alibabacloudstack_slb_servercertificates
-> **NOTE:** Alias name has: `alibabacloudstack_slb_server_certificates`

This data source provides a list of SLB server certificates in an AlibabacloudStack account according to the specified filters.

## Example Usage

```hcl
data "alibabacloudstack_slb_servercertificates" "example" {
  name_regex = "example_cert.*"

  ids = ["cert-12345678"]
}

output "first_certificate_id" {
  value = data.alibabacloudstack_slb_servercertificates.example.certificates.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of SLB server certificate IDs used to filter results.
* `name_regex` - (Optional, ForceNew) A regex string used to filter results by SLB server certificate name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `names` - A list of SLB server certificate names.
* `certificates` - A list of SLB server certificates. Each element contains the following attributes:
  * `id` - The ID of the SLB server certificate.
  * `name` - The name of the SLB server certificate.
  * `fingerprint` - The fingerprint of the server certificate.
  * `created_time` - The creation time of the server certificate in human-readable format.
  * `created_timestamp` - The creation timestamp of the server certificate in Unix epoch time.
