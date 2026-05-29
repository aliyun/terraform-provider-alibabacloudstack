---
subcategory: "Server Load Balancer (SLB)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_slb_cacertificate"
sidebar_current: "docs-Alibabacloudstack-resource-slb-cacertificate"
description: |- 
  Provides a slb Cacertificate resource.
---

# alibabacloudstack_slb_cacertificate
-> **NOTE:** Alias name has: `alibabacloudstack_slb_ca_certificate`

Provides a slb Cacertificate resource.

## Example Usage

### Using CA Certificate Content

```hcl
variable "name" {
    default = "tf-testaccslbca_certificate67317"
}

resource "alibabacloudstack_slb_cacertificate" "default" {
  name              = var.name
  ca_certificate    = "-----BEGIN CERTIFICATE-----\nMIIDRjCCAq+gAwIBAgI+OuMs******XTtI90EAxEG/bJJyOm5LqoiA=\n-----END CERTIFICATE-----"
}
```

### Using CA Certificate File

```hcl
resource "alibabacloudstack_slb_cacertificate" "file_example" {
  name           = "tf-testaccslbca_certificate_file"
  ca_certificate = file("${path.module}/ca_certificate.pem")
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional, Deprecated) The name of the CA Certificate. **This field is deprecated**, please use `ca_certificate_name` instead. Conflicts with `ca_certificate_name`.
* `ca_certificate_name` - (Optional) The name of the CA certificate. Conflicts with `name`.
* `ca_certificate` - (Required, ForceNew) The content of the CA certificate in PEM format. This field is immutable; modifying it will force the creation of a new resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the CA Certificate, which is the `CACertificateId` returned by the API.
* `name` - The name of the CA Certificate.
* `ca_certificate_name` - The name of the CA Certificate.

## Import

SLB CA Certificate can be imported using the `CACertificateId`, e.g.

```
$ terraform import alibabacloudstack_slb_cacertificate.example CACertificateId
```

-> **NOTE:** The `ca_certificate` argument cannot be imported as it does not echo back from the API. Use `ImportStateVerifyIgnore` in tests or expect a non-empty plan after import.