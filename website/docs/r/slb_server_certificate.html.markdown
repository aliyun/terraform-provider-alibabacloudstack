---
subcategory: "Server Load Balancer (SLB)"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_slb_server_certificate"
sidebar_current: "docs-apsarastack-resource-slb-server-certificate"
description: |-
  Provides a Load Banlancer Server Certificate resource.
---

# apsarastack\_slb\_server\_certificate

A Load Balancer Server Certificate is an ssl Certificate used by the listener of the protocol https.


## Example Usage

* using server_certificate/private content as string example

```
# create a server certificate
resource "apsarastack_slb_server_certificate" "foo" {
  name               = "slbservercertificate"
  server_certificate = "-----BEGIN CERTIFICATE-----\nMIIDRjCCAq+gAwIBAgI+OuMs******XTtI90EAxEG/bJJyOm5LqoiA=\n-----END CERTIFICATE-----"
  private_key        = "-----BEGIN RSA PRIVATE KEY-----\nMIICXAIBAAKBgQDO0knDrlNdiys******ErVpjsckAaOW/JDG5PCSwkaMxk=\n-----END RSA PRIVATE KEY-----"
}
```

* using server_certificate/private file example

```
# create a server certificate
resource "apsarastack_slb_server_certificate" "foo" {
  name               = "slbservercertificate"
  server_certificate = file("${path.module}/server_certificate.pem")
  private_key        = file("${path.module}/private_key.pem")
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) Name of the Server Certificate.
* `server_certificate` - (Optional, ForceNew) the content of the ssl certificate. where `apsarastack_certificate_id` is null, it is required, otherwise it is ignored.
* `private_key` - (Optional, ForceNew) the content of privat key of the ssl certificate specified by `server_certificate`. where `apsarastack_certificate_id` is null, it is required, otherwise it is ignored.

## Attributes Reference

The following attributes are exported:

* `id` - The Id of Server Certificate (SSL Certificate).

