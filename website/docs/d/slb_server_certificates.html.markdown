---
subcategory: "Server Load Balancer (SLB)"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_slb_server_certificates"
sidebar_current: "docs-apsarastack-datasource-slb-server-certificates"
description: |-
    Provides a list of slb server certificates.
---
# apsarastack\_slb_server_certificates

This data source provides the server certificate list.

## Example Usage

```
data "apsarastack_slb_server_certificates" "sample_ds" {
}

output "first_slb_server_certificate_id" {
  value = "${data.apsarastack_slb_server_certificates.sample_ds.certificates.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of server certificates IDs to filter results.
* `name_regex` - (Optional) A regex string to filter results by server certificate name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of SLB server certificates IDs.
* `names` - A list of SLB server certificates names.
* `certificates` - A list of SLB server certificates. Each element contains the following attributes:
  * `id` - Server certificate ID.
  * `name` - Server certificate name.
  * `fingerprint` - Server certificate fingerprint.
  * `created_time` - Server certificate created time.
  * `created_timestamp` - Server certificate created timestamp.
