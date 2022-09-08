---
subcategory: "Server Load Balancer (SLB)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_slb_ca_certificates"
sidebar_current: "docs-alibabacloudstack-datasource-slb-ca-certificates"
description: |-
    Provides a list of slb CA certificates.
---
# alibabacloudstack\_slb_ca_certificates

This data source provides the CA certificate list.

## Example Usage

```
data "alibabacloudstack_slb_ca_certificates" "sample_ds" {
}

output "first_slb_ca_certificate_id" {
  value = "${data.alibabacloudstack_slb_ca_certificates.sample_ds.certificates.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of ca certificates IDs to filter results.
* `name_regex` - (Optional) A regex string to filter results by ca certificate name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of SLB ca certificates IDs.
* `names` - A list of SLB ca certificates names.
* `certificates` - A list of SLB ca certificates. Each element contains the following attributes:
  * `id` - CA certificate ID.
  * `name` - CA certificate name.
  * `fingerprint` - CA certificate fingerprint.
  * `created_time` - CA certificate created time.
  * `created_timestamp` - CA certificate created timestamp.
  * `region_id` - The region Id of CA certificate.
