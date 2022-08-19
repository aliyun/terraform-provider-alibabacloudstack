---
subcategory: "Server Load Balancer (SLB)"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_slb_ca_certificate"
sidebar_current: "docs-apsarastack-resource-slb-ca-certificate"
description: |-
  Provides a Load Banlancer CA Certificate resource.
---

# apsarastack\_slb\_ca\_certificate

A Load Balancer CA Certificate is used by the listener of the protocol https.


## Example Usage

* using CA certificate content

```
# create a CA certificate
resource "apsarastack_slb_ca_certificate" "foo" {
  name           = "tf-testAccSlbCACertificate"
  ca_certificate = "-----BEGIN CERTIFICATE-----\nMIIDRjCCAq+gAwIBAgIJAJnI******90EAxEG/bJJyOm5LqoiA=\n-----END CERTIFICATE-----"
}
```

* using CA certificate file

```
resource "apsarastack_slb_ca_certificate" "foo-file" {
  name           = "tf-testAccSlbCACertificate"
  ca_certificate = file("${path.module}/ca_certificate.pem")
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) Name of the CA Certificate.
* `ca_certificate` - (Required, ForceNew) the content of the CA certificate.

## Attributes Reference

The following attributes are exported:

* `id` - The Id of CA Certificate .
