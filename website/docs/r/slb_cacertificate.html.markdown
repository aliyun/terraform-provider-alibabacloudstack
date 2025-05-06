---
subcategory: "Server Load Balancer (SLB)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_slb_cacertificate"
sidebar_current: "docs-Alibabacloudstack-slb-cacertificate"
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

* `name` - (Optional) The name of the CA Certificate. This can be used to identify the certificate.
* `ca_certificate_name` - (Optional) The name of the CA Certificate, which serves as an identifier for the certificate.
* `ca_certificate` - (Required, ForceNew) The content of the CA certificate in PEM format. This field is immutable and cannot be updated after creation.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the CA Certificate, which uniquely identifies the resource.
* `name` - The name of the CA Certificate as provided during creation.
* `ca_certificate_name` - The name of the CA Certificate, which is useful for referencing the certificate in other resources or configurations.