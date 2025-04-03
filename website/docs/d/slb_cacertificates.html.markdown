---
subcategory: "SLB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_slb_cacertificates"
sidebar_current: "docs-Alibabacloudstack-datasource-slb-cacertificates"
description: |- 
  Provides a list of slb cacertificates owned by an alibabacloudstack account.
---

# alibabacloudstack_slb_cacertificates
-> **NOTE:** Alias name has: `alibabacloudstack_slb_ca_certificates`

This data source provides a list of slb cacertificates in an alibabacloudstack account according to the specified filters.

## Example Usage

```hcl
data "alibabacloudstack_slb_cacertificates" "sample_ds" {
  name_regex = "example_cert"
  ids        = ["cert-id-1", "cert-id-2"]
}

output "first_slb_ca_certificate_id" {
  value = "${data.alibabacloudstack_slb_cacertificates.sample_ds.certificates.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of CA certificate IDs to filter results.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by CA certificate name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `names` - A list of SLB CA certificate names.
* `certificates` - A list of SLB CA certificates. Each element contains the following attributes:
  * `id` - The ID of the CA certificate.
  * `name` - The name of the CA certificate.
  * `fingerprint` - The fingerprint of the CA certificate.
  * `created_timestamp` - The timestamp when the CA certificate was created.
  * `region_id` - The region ID where the CA certificate is located.