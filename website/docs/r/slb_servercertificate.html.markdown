---
subcategory: "SLB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_slb_servercertificate"
sidebar_current: "docs-Alibabacloudstack-slb-servercertificate"
description: |- 
  Provides a slb Servercertificate resource.
---

# alibabacloudstack_slb_servercertificate
-> **NOTE:** Alias name has: `alibabacloudstack_slb_server_certificate`

Provides a slb Servercertificate resource.

## Example Usage

### Example 1: Using server_certificate/private_key as string content

```hcl
# Create a server certificate using string content
resource "alibabacloudstack_slb_servercertificate" "foo" {
  name                = "slbservercertificate"
  server_certificate  = "-----BEGIN CERTIFICATE-----\nMIIDRjCCAq+gAwIBAgI+OuMs******XTtI90EAxEG/bJJyOm5LqoiA=\n-----END CERTIFICATE-----"
  private_key         = "-----BEGIN RSA PRIVATE KEY-----\nMIICXAIBAAKBgQDO0knDrlNdiys******ErVpjsckAaOW/JDG5PCSwkaMxk=\n-----END RSA PRIVATE KEY-----"
}
```

### Example 2: Using server_certificate/private_key from files

```hcl
# Create a server certificate using file content
resource "alibabacloudstack_slb_servercertificate" "foo" {
  name                = "slbservercertificate"
  server_certificate  = file("${path.module}/server_certificate.pem")
  private_key         = file("${path.module}/private_key.pem")
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the server certificate. If not provided, Terraform will auto-generate a unique name.
* `server_certificate_name` - (Optional) The name of the server certificate. This can be used to identify the certificate in the SLB service.
* `server_certificate` - (Required, ForceNew) The public key certificate to be uploaded. This parameter is required if you do not use an Alibaba Cloud-managed certificate.
* `private_key` - (Required, ForceNew) The private key corresponding to the public key certificate specified in `server_certificate`. This parameter is required if you do not use an Alibaba Cloud-managed certificate.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the server certificate (SSL Certificate).
* `name` - The name of the server certificate.
* `server_certificate_name` - The name of the server certificate as specified during creation.
* `server_certificate` -  Represents the public key certificate that was uploaded.
* `private_key` -  Represents the private key corresponding to the uploaded public key certificate.